package realtime

import (
	"runtime/cgo"
	"unsafe"
)

// Stage is a connection to a single IVS Stage. Create one with NewStage and
// release it with Close. A Stage must not be used concurrently from multiple
// goroutines.
type Stage struct {
	ptr unsafe.Pointer
	// events keeps the Go event-callback state reachable for as long as the SDK
	// may deliver events.
	events cgo.Handle
	// strategy is the strategy this Stage was created with, and whether the Stage
	// owns it. See strategy.go.
	strategy boundStrategy
	// audio owns the mixed PCM delivery state, including its own anchor. See
	// audio.go.
	audio pcmDelivery
	// streams gives each stream the SDK reports one Go value for as long as it
	// lives. See stage_stream.go.
	streams streamRegistry
}

// Version returns the version of the SDK in use.
func Version() string {
	version, _ := rawVersion()
	return version
}

// StageConfig configures a Stage. Pass it to NewStage.
type StageConfig struct {
	// Token is the IVS Stage participant token the Stage joins with.
	Token string
	// Strategy decides what the Stage subscribes to and publishes: your own
	// StageStrategy implementation or a ChannelTranscriber's Strategy. It is
	// required. A Strategy borrowed from a ChannelTranscriber outlives the
	// Stage, so close the Stage before the transcriber.
	Strategy StageStrategy
	// Hooks receives Stage events. The zero value receives none.
	Hooks StageHooks
}

// NewStage creates a Stage. It does not connect; call Join for that.
//
// A nil config.Strategy returns ErrorCodeInvalidArgument.
func NewStage(config StageConfig) (*Stage, error) {
	bound, err := bindStrategy(config.Strategy)
	if err != nil {
		return nil, err
	}

	state := &eventState{hooks: config.Hooks}
	anchor := cgo.NewHandle(state)

	ptr, err := rawStageCreate(config.Token, bound.ptr, anchor)
	if err != nil {
		anchor.Delete()
		bound.release()
		return nil, err
	}

	stage := &Stage{ptr: ptr, events: anchor, strategy: bound}
	state.stage = stage
	return stage, nil
}

// eventState is what the SDK holds as user_data for a Stage's events. A Go
// pointer cannot cross into C, so it travels as a cgo.Handle that each trampoline
// resolves back to this value.
type eventState struct {
	stage *Stage
	hooks StageHooks
}

// SetPCMCallback registers fn to receive mixed PCM audio: every remote
// participant the Stage is subscribed to, combined into one stream. Pass nil to
// stop delivery.
//
// Call this before Join or after Leave; on a joined Stage it returns
// ErrorCodeStageSetAudioCallbackWhileJoined and changes nothing. A registration
// survives Leave, so a Stage that rejoins keeps delivering to fn.
//
// fn runs on an audio thread and must return promptly. Do not call SDK methods from
// fn; hand the buffer to a goroutine for anything more than trivial work. The
// PCMBuffer it receives is owned by the caller and may be retained.
func (s *Stage) SetPCMCallback(fn func(PCMBuffer)) error {
	return s.audio.set(func(anchor cgo.Handle, install bool) error {
		return rawStageSetAudioCallback(s.ptr, anchor, install)
	}, fn)
}

// Join connects to the Stage its token was issued for. It returns once the
// attempt is under way, not once it succeeds: progress arrives through
// StageHooks.OnConnectionStateChanged.
//
// Join returns ErrorCodeStageInvalidToken when the token cannot be parsed, and
// ErrorCodeStageExpiredToken when the token has expired.
func (s *Stage) Join() error {
	return rawStageJoin(s.ptr)
}

// Leave disconnects from the Stage, and does nothing if not joined. The Stage
// remains usable and may be joined again.
func (s *Stage) Leave() error {
	return rawStageLeave(s.ptr)
}

// RefreshStrategy re-runs the strategy. It is consulted again for every
// participant currently in the Stage, and subscriptions are updated to match its
// new answers. Safe to call from an event callback.
func (s *Stage) RefreshStrategy() error {
	return rawStageRefreshStrategy(s.ptr)
}

// Close disconnects and releases the Stage. It is safe to call more than once, in
// any order relative to other Close calls, and satisfies io.Closer.
//
// Do not call Close concurrently with any other method on the same value.
func (s *Stage) Close() error {
	if s.ptr == nil {
		return nil
	}
	// Destroy first: it returns only once the SDK has stopped delivering
	// callbacks, which is what makes releasing the anchors below safe.
	rawStageDestroy(s.ptr)
	s.ptr = nil
	if s.events != 0 {
		s.events.Delete()
		s.events = 0
	}
	s.strategy.release()
	s.audio.release()
	s.streams.release()
	return nil
}
