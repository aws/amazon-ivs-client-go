package realtime

import (
	"runtime/cgo"
	"unsafe"
)

// StageHooks receives Stage events. Each hook is optional, and a nil hook skips
// its event. Every hook is called on one serialized thread, so state used only
// by these hooks needs no synchronization.
type StageHooks struct {
	OnConnectionStateChanged           func(ConnectionStateChangedEvent)
	OnParticipantJoined                func(ParticipantJoinedEvent)
	OnParticipantLeft                  func(ParticipantLeftEvent)
	OnParticipantMetadataUpdated       func(ParticipantMetadataUpdatedEvent)
	OnParticipantPublishStateChanged   func(ParticipantPublishStateChangedEvent)
	OnParticipantSubscribeStateChanged func(ParticipantSubscribeStateChangedEvent)
	OnStreamsAdded                     func(StreamsAddedEvent)
	OnStreamsRemoved                   func(StreamsRemovedEvent)
	OnStreamsMutedChanged              func(StreamsMutedChangedEvent)
	OnError                            func(ErrorEvent)
}

// ConnectionStateChangedEvent reports that the Stage connection state changed.
type ConnectionStateChangedEvent struct {
	State ConnectionState
	// Err is the error that caused the transition, or nil for an expected
	// transition.
	Err *Error
}

// ParticipantJoinedEvent reports that a participant joined the Stage.
type ParticipantJoinedEvent struct {
	Participant ParticipantInfo
}

// ParticipantLeftEvent reports that a participant left the Stage or stopped
// publishing.
type ParticipantLeftEvent struct {
	Participant ParticipantInfo
}

// ParticipantMetadataUpdatedEvent reports that a participant's UserID or
// Attributes changed. Participant is a new snapshot with the same ParticipantID
// and IsLocal as the prior one.
type ParticipantMetadataUpdatedEvent struct {
	Participant ParticipantInfo
}

// ParticipantPublishStateChangedEvent reports that a participant's publish
// state changed.
type ParticipantPublishStateChangedEvent struct {
	Participant ParticipantInfo
	State       PublishState
}

// ParticipantSubscribeStateChangedEvent reports that the local participant's
// subscription to a remote publisher changed. Participant is that publisher.
type ParticipantSubscribeStateChangedEvent struct {
	Participant ParticipantInfo
	State       SubscribeState
}

// StreamsAddedEvent reports that a participant added streams. For a remote
// participant, a subscription took effect. For the local participant, these are
// the streams the strategy asked to publish.
type StreamsAddedEvent struct {
	Participant ParticipantInfo
	Streams     []StageStream
}

// StreamsRemovedEvent reports that a participant removed streams: it
// unpublished, stopped sharing a stream type, or left the Stage.
type StreamsRemovedEvent struct {
	Participant ParticipantInfo
	Streams     []StageStream
}

// StreamsMutedChangedEvent reports that streams changed mute state.
type StreamsMutedChangedEvent struct {
	Participant ParticipantInfo
	Streams     []StageStream
}

// ErrorEvent reports an asynchronous error that leaves the Stage connected,
// such as a failed publish or subscribe. An error that changes the connection,
// such as a disconnect, arrives in ConnectionStateChangedEvent instead.
type ErrorEvent struct {
	// Err is never nil.
	Err *Error
}

// events resolves the anchor the SDK carries as user_data.
func events(userData unsafe.Pointer) *eventState {
	return cgo.Handle(uintptr(userData)).Value().(*eventState)
}

// The dispatch functions below are the Go end of each generated trampoline. They
// exist one per spec callback, so adding a callback to the spec fails the build
// here until it is handled.

func dispatchConnectionStateChanged(userData unsafe.Pointer, connectionState ConnectionState, err *Error) {
	if hook := events(userData).hooks.OnConnectionStateChanged; hook != nil {
		hook(ConnectionStateChangedEvent{State: connectionState, Err: err})
	}
}

func dispatchParticipantJoined(userData unsafe.Pointer, participant ParticipantInfo) {
	if hook := events(userData).hooks.OnParticipantJoined; hook != nil {
		hook(ParticipantJoinedEvent{Participant: participant})
	}
}

func dispatchParticipantLeft(userData unsafe.Pointer, participant ParticipantInfo) {
	if hook := events(userData).hooks.OnParticipantLeft; hook != nil {
		hook(ParticipantLeftEvent{Participant: participant})
	}
}

func dispatchParticipantMetadataUpdated(userData unsafe.Pointer, participant ParticipantInfo) {
	if hook := events(userData).hooks.OnParticipantMetadataUpdated; hook != nil {
		hook(ParticipantMetadataUpdatedEvent{Participant: participant})
	}
}

func dispatchParticipantPublishStateChanged(userData unsafe.Pointer, participant ParticipantInfo, publishState PublishState) {
	if hook := events(userData).hooks.OnParticipantPublishStateChanged; hook != nil {
		hook(ParticipantPublishStateChangedEvent{Participant: participant, State: publishState})
	}
}

func dispatchParticipantSubscribeStateChanged(userData unsafe.Pointer, participant ParticipantInfo, subscribeState SubscribeState) {
	if hook := events(userData).hooks.OnParticipantSubscribeStateChanged; hook != nil {
		hook(ParticipantSubscribeStateChangedEvent{Participant: participant, State: subscribeState})
	}
}

func dispatchStreamsAdded(userData unsafe.Pointer, participant ParticipantInfo, streams []unsafe.Pointer) {
	state := events(userData)
	added := state.stage.stageStreamsFromC(streams)
	if hook := state.hooks.OnStreamsAdded; hook != nil {
		hook(StreamsAddedEvent{Participant: participant, Streams: added})
	}
}

func dispatchStreamsRemoved(userData unsafe.Pointer, participant ParticipantInfo, streams []unsafe.Pointer) {
	state := events(userData)
	removed := state.stage.stageStreamsFromC(streams)
	// Retire after reporting, so the streams stay intact for the handler.
	if hook := state.hooks.OnStreamsRemoved; hook != nil {
		hook(StreamsRemovedEvent{Participant: participant, Streams: removed})
	}
	state.stage.afterStreamsRemoved(removed)
}

func dispatchStreamsMutedChanged(userData unsafe.Pointer, participant ParticipantInfo, streams []unsafe.Pointer) {
	state := events(userData)
	changed := state.stage.stageStreamsFromC(streams)
	if hook := state.hooks.OnStreamsMutedChanged; hook != nil {
		hook(StreamsMutedChangedEvent{Participant: participant, Streams: changed})
	}
}

func dispatchError(userData unsafe.Pointer, err *Error) {
	if hook := events(userData).hooks.OnError; hook != nil {
		hook(ErrorEvent{Err: err})
	}
}

// The strategy dispatchers resolve a different anchor: a strategy is registered
// on its own, and the SDK may consult it before any Stage exists.

func strategyFor(userData unsafe.Pointer) StageStrategy {
	return cgo.Handle(uintptr(userData)).Value().(StageStrategy)
}

func dispatchSubscriptionTypes(userData unsafe.Pointer, participant ParticipantInfo) StreamType {
	return strategyFor(userData).SubscriptionTypes(participant)
}

func dispatchShouldPublish(userData unsafe.Pointer, participant ParticipantInfo) bool {
	return strategyFor(userData).ShouldPublish(participant)
}

func dispatchStreamsToPublish(userData unsafe.Pointer, participant ParticipantInfo, streams unsafe.Pointer) {
	applyPublishStreams(streams, strategyFor(userData).StreamsToPublish(participant))
}

func dispatchAudio(userData unsafe.Pointer, data []byte, numFrames, bytesPerSample, numChannels int, samplesPerSec uint32, timestampMs int64) {
	cgo.Handle(uintptr(userData)).Value().(*pcmDelivery).deliver(PCMBuffer{
		Data:           data,
		NumFrames:      numFrames,
		BytesPerSample: bytesPerSample,
		NumChannels:    numChannels,
		SamplesPerSec:  samplesPerSec,
		TimestampMs:    timestampMs,
	})
}
