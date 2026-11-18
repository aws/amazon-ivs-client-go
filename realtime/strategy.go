package realtime

import (
	"runtime/cgo"
	"unsafe"
)

// PublishStreams is the set of streams a Stage publishes. A zero PublishStreams
// publishes nothing. A Stage publishes at most one audio source, and any number
// of subtitle streams.
type PublishStreams struct {
	// Audio is the audio source to publish, or nil to publish no audio.
	Audio *AudioSource
	// SubtitleStreams are the subtitle streams to publish. Leave it empty to
	// publish no subtitles. Listing the same stream more than once publishes it
	// once.
	SubtitleStreams []*LocalSubtitleSource
}

// StageStrategy decides what a Stage subscribes to and what it publishes.
// Implement it and set it as StageConfig.Strategy, or use a ChannelTranscriber's
// Strategy.
//
// Every method is required, and every one is called on the same serialized
// thread as StageHooks, so state shared only among them needs no
// synchronization. Call Stage.RefreshStrategy whenever an answer changes.
type StageStrategy interface {
	// SubscriptionTypes returns the stream types to subscribe to from a remote
	// participant: StreamTypeNone to skip it, or StreamTypeAudio for its
	// audio. It is called for each remote participant when first seen, and
	// never for the local participant.
	SubscriptionTypes(participant ParticipantInfo) StreamType

	// ShouldPublish reports whether to publish as the local participant. It is
	// called on joining, and only if the participant token permits publishing.
	// Returning true leads to a StreamsToPublish call.
	ShouldPublish(participant ParticipantInfo) bool

	// StreamsToPublish returns what to publish as the local participant. It is
	// called after ShouldPublish returns true. The SDK compares each field
	// against what it is already publishing: a source already named is left
	// alone, and a null field stops publishing that kind of media. ShouldPublish
	// controls whether the participant publishes at all; a zero PublishStreams
	// here stops that kind of media, not the participant's publishing state.
	StreamsToPublish(participant ParticipantInfo) PublishStreams
}

// prebuiltStrategy is a StageStrategy whose policy is fixed and applied by the
// SDK itself, as ChannelTranscriber's is. That policy is applied directly, so the
// methods below are never called and answer neutrally to satisfy the interface.
type prebuiltStrategy struct {
	ptr unsafe.Pointer
}

func (prebuiltStrategy) SubscriptionTypes(ParticipantInfo) StreamType {
	return StreamTypeNone
}

func (prebuiltStrategy) ShouldPublish(ParticipantInfo) bool { return false }

func (prebuiltStrategy) StreamsToPublish(ParticipantInfo) PublishStreams { return PublishStreams{} }

// boundStrategy is the strategy a Stage was created with, and whether the Stage
// owns it: it owns what it builds for a StageStrategy implementation, and borrows
// a prebuiltStrategy's from whatever created that.
type boundStrategy struct {
	ptr      unsafe.Pointer
	goHandle cgo.Handle
	owned    bool
}

// newStrategy registers the callbacks that consult strategy, returning the SDK
// handle plus the anchor keeping strategy reachable for as long as the SDK may
// ask it anything.
func newStrategy(strategy StageStrategy) (unsafe.Pointer, cgo.Handle, error) {
	anchor := cgo.NewHandle(strategy)
	ptr, err := rawStageStrategyCreate(anchor)
	if err != nil {
		anchor.Delete()
		return nil, 0, err
	}
	return ptr, anchor, nil
}

func bindStrategy(strategy StageStrategy) (boundStrategy, error) {
	switch s := strategy.(type) {
	case nil:
		return boundStrategy{}, nil
	case prebuiltStrategy:
		return boundStrategy{ptr: s.ptr}, nil
	default:
		ptr, goHandle, err := newStrategy(s)
		if err != nil {
			return boundStrategy{}, err
		}
		return boundStrategy{ptr: ptr, goHandle: goHandle, owned: true}, nil
	}
}

func (b *boundStrategy) release() {
	if !b.owned {
		return
	}
	rawStageStrategyDestroy(b.ptr)
	b.ptr = nil
	b.goHandle.Delete()
	b.goHandle = 0
	b.owned = false
}

// ChannelTranscriber is a live-transcription helper. It owns a StageStrategy that
// subscribes to audio from one remote participant, and it accepts
// transcription results for delivery back to the Stage. Set its Strategy as
// StageConfig.Strategy, then submit results with WriteTranscription.
type ChannelTranscriber struct {
	ptr      unsafe.Pointer
	strategy prebuiltStrategy
}

// ChannelTranscriberConfig configures a ChannelTranscriber. Pass it to
// NewChannelTranscriber.
type ChannelTranscriberConfig struct {
	// ParticipantID is the remote participant to transcribe. It must be
	// non-empty.
	ParticipantID string
}

// NewChannelTranscriber creates a transcriber for one remote participant. Its
// Strategy subscribes to audio from config.ParticipantID and no other remote
// participant. An empty config.ParticipantID returns ErrorCodeInvalidArgument.
// Close the returned value when done, after closing any Stage created from its
// Strategy.
func NewChannelTranscriber(config ChannelTranscriberConfig) (*ChannelTranscriber, error) {
	ptr, err := rawChannelTranscriberCreate(config.ParticipantID)
	if err != nil {
		return nil, err
	}
	stratPtr, err := rawChannelTranscriberGetStrategy(ptr)
	if err != nil {
		rawChannelTranscriberDestroy(ptr)
		return nil, err
	}
	return &ChannelTranscriber{ptr: ptr, strategy: prebuiltStrategy{ptr: stratPtr}}, nil
}

// Strategy returns the StageStrategy owned by this transcriber, for use as
// StageConfig.Strategy. It is valid until the transcriber is closed.
func (t *ChannelTranscriber) Strategy() StageStrategy { return t.strategy }

// Transcription is one transcription result. Submit it with
// ChannelTranscriber.WriteTranscription.
type Transcription struct {
	// StartMs and EndMs bound the transcribed audio in milliseconds, on the
	// same timeline as PCMBuffer.TimestampMs.
	StartMs, EndMs int64
	// LanguageCode is opaque: the SDK carries it through unvalidated, so any
	// non-empty scheme the receiver understands works. Do not assume BCP-47 or
	// ISO 639; agree on a scheme with the receiver.
	LanguageCode string
	// Text is attributed to the remote participant whose audio the Stage is
	// subscribed to.
	Text string
	// UtteranceID aggregates the partials of one utterance across writes, so a
	// receiver can de-duplicate partials for the same utterance even as StartMs
	// shifts between them. Leave it empty if the caller has no such identifier;
	// the SDK then uses StartMs as the ID.
	//
	// Avoid IDs that reset or repeat, such as a per-connection sequence number.
	// If a transcription service assigns each result the next number starting
	// from 0, and a network drop makes your app reconnect, the new connection
	// also starts counting from 0. The receiver then sees the same ID twice for
	// two different utterances.
	UtteranceID string
}

// WriteTranscription submits a transcription result for delivery to the Stage.
// It is safe to call from any goroutine.
//
// EndMs must be greater than StartMs, or this returns
// ErrorCodeChannelTranscriptionInvalidTimeRange. StartMs must not lag the live
// audio position by more than 60 seconds, or this returns
// ErrorCodeChannelTranscriptionStale. If no remote participant has subscribed
// audio, this returns ErrorCodeChannelTranscriptionNoParticipant.
//
// Returns ErrorCodeChannelTranscriberNotAttached if this transcriber's Strategy has
// not been passed to a Stage yet, and ErrorCodeChannelTranscriptionCaptureTimeUnavailable
// if delivery timing could not be established for the participant's audio
// stream.
func (t *ChannelTranscriber) WriteTranscription(transcription Transcription) error {
	return rawChannelTranscriberWriteTranscription(t.ptr, transcription.StartMs, transcription.EndMs, transcription.LanguageCode, transcription.Text, transcription.UtteranceID)
}

// Close releases the transcriber and the strategy it owns. It is safe to call
// more than once and satisfies io.Closer. Close any Stage using this
// transcriber's Strategy first.
//
// Do not call Close concurrently with any other method on the same value.
func (t *ChannelTranscriber) Close() error {
	if t.ptr == nil {
		return nil
	}
	rawChannelTranscriberDestroy(t.ptr)
	t.ptr = nil
	t.strategy.ptr = nil
	return nil
}
