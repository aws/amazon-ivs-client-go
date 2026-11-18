package realtime

import (
	"runtime"
	"runtime/cgo"
	"unsafe"
)

// StageStream is one stream from one participant.
//
// The SDK emits a StageStream to StageHooks.OnStreamsAdded. Later
// events about the same stream emit the same value. Use the == operator to
// compare two streams.
//
// You can retain a StageStream after the callback returns. After the SDK reports
// the stream in OnStreamsRemoved, the stream becomes inert. It stays safe to use
// and StreamID keeps its value, but it delivers no more content and accepts no
// new registration.
type StageStream interface {
	// StreamID returns the identifier of the stream. The identifier is unique
	// within the Stage. The identifier does not change while the stream exists.
	StreamID() string
}

// AudioStageStream is a remote audio stream.
type AudioStageStream interface {
	StageStream

	// IsMuted returns true if the stream is muted on the Stage. A muted stream
	// sends no media to subscribers. The value does not describe local
	// playback.
	//
	// StageHooks.OnStreamsMutedChanged reports only a change of state.
	// A stream can be muted when the SDK first reports it. Use IsMuted to read
	// the current state.
	IsMuted() bool

	// SetPCMCallback registers fn to receive the audio of this stream. The SDK
	// calls fn with the audio of this participant only.
	//
	// To stop delivery, set fn to nil. To replace the callback, call
	// SetPCMCallback again.
	//
	// Register the callback in StageHooks.OnStreamsAdded. The
	// registration stays active while the stream exists. The SDK reports the
	// end of a stream in OnStreamsRemoved, and then SetPCMCallback returns
	// ErrorCodeStageStreamEnded. A replacement stream has a different
	// identifier and needs a new registration. On a stream that is not audio,
	// SetPCMCallback returns ErrorCodeStageStreamNotAudio. You do not need to retain the stream to
	// continue to receive audio.
	//
	// The SDK calls fn on a real-time audio thread. Every stream callback and
	// Stage.SetPCMCallback share this one thread, and a slow callback delays
	// the audio of every participant. Return from fn quickly and do not call
	// SDK methods from fn. To do more work, send the buffer to a goroutine. You
	// own the PCMBuffer and can retain it.
	SetPCMCallback(fn func(PCMBuffer)) error
}

// VideoStageStream is a remote video stream.
type VideoStageStream interface {
	StageStream

	// IsMuted returns true if the stream is muted on the Stage. A muted stream
	// sends no media to subscribers. The value does not describe local
	// playback.
	//
	// StageHooks.OnStreamsMutedChanged reports only a change of state.
	// A stream can be muted when the SDK first reports it. Use IsMuted to read
	// the current state.
	IsMuted() bool
}

// SubtitleStageStream is a subtitle stream, local or remote. A LocalSubtitleSource
// a Stage publishes appears here too, on the local participant.
type SubtitleStageStream interface {
	StageStream
}

type stageStream struct {
	ptr      unsafe.Pointer
	streamID string
}

func (s *stageStream) StreamID() string {
	return s.streamID
}

func (s *stageStream) IsMuted() bool {
	muted, _ := rawStageStreamGetIsMuted(s.ptr)
	return muted
}

type audioStageStream struct {
	*stageStream
	audio pcmDelivery
}

func (s *audioStageStream) SetPCMCallback(fn func(PCMBuffer)) error {
	return s.audio.set(func(anchor cgo.Handle, install bool) error {
		return rawStageStreamSetAudioCallback(s.ptr, anchor, install)
	}, fn)
}

type videoStageStream struct {
	*stageStream
}

type subtitleStageStream struct {
	*stageStream
}

type unsupportedStageStream struct {
	*stageStream
}

// Unsynchronized because the SDK delivers every event on one serialized thread,
// which is the only place these entries are added or removed.
type streamRegistry struct {
	byID map[string]StageStream
}

func (s *Stage) stageStreamsFromC(borrowed []unsafe.Pointer) []StageStream {
	if len(borrowed) == 0 {
		return nil
	}
	streams := make([]StageStream, 0, len(borrowed))
	for _, handle := range borrowed {
		id, err := rawStageStreamGetStreamID(handle)
		if err != nil {
			continue
		}
		if stream := s.resolveStream(id, handle); stream != nil {
			streams = append(streams, stream)
		}
	}
	return streams
}

func (s *Stage) resolveStream(id string, borrowed unsafe.Pointer) StageStream {
	if existing, ok := s.streams.byID[id]; ok {
		return existing
	}
	stream := newStageStream(borrowed, id)
	if stream == nil {
		return nil
	}
	if s.streams.byID == nil {
		s.streams.byID = make(map[string]StageStream)
	}
	s.streams.byID[id] = stream
	return stream
}

func (s *Stage) afterStreamsRemoved(streams []StageStream) {
	for _, stream := range streams {
		releaseStageStream(stream)
		delete(s.streams.byID, stream.StreamID())
	}
}

func (r *streamRegistry) release() {
	for id, stream := range r.byID {
		releaseStageStream(stream)
		delete(r.byID, id)
	}
}

func releaseStageStream(stream StageStream) {
	if audio, ok := stream.(*audioStageStream); ok {
		audio.audio.release()
	}
}

func newStageStream(borrowed unsafe.Pointer, id string) StageStream {
	ptr, err := rawStageStreamClone(borrowed)
	if err != nil {
		return nil
	}
	streamType, err := rawStageStreamGetStreamType(ptr)
	if err != nil {
		rawStageStreamDestroy(ptr)
		return nil
	}
	stream := &stageStream{ptr: ptr, streamID: id}
	runtime.AddCleanup(stream, rawStageStreamDestroy, ptr)
	switch streamType {
	case StreamTypeAudio:
		return &audioStageStream{stageStream: stream}
	case StreamTypeVideo:
		return &videoStageStream{stageStream: stream}
	case StreamTypeSubtitle:
		return &subtitleStageStream{stageStream: stream}
	}
	return &unsupportedStageStream{stageStream: stream}
}
