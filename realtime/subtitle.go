//go:build cgo

package realtime

/*
#include <aws_ivs_realtime.h>
#include <stdlib.h>
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

// LocalSubtitleSource publishes subtitle cues. Create one, list it in
// PublishStreams.SubtitleStreams, then send cues with SendCue. Close it when
// done.
//
// A subtitle stream carries no media. A Stage can publish any number of them at
// once, and publishing one does not make the participant an audio or video
// publisher.
type LocalSubtitleSource struct {
	ptr unsafe.Pointer
}

// NewLocalSubtitleSource creates a subtitle source. It publishes nothing until a
// StageStrategy lists it in PublishStreams.SubtitleStreams.
func NewLocalSubtitleSource() (*LocalSubtitleSource, error) {
	ptr, err := rawLocalSubtitleStreamCreate()
	if err != nil {
		return nil, err
	}
	return &LocalSubtitleSource{ptr: ptr}, nil
}

// StreamID returns the identifier of this stream. It is unique within the Stage
// and does not change while the stream exists. It matches the StreamID of the
// SubtitleStageStream the SDK reports for this stream in
// StageHooks.OnStreamsAdded.
func (s *LocalSubtitleSource) StreamID() string {
	id, _ := rawLocalSubtitleStreamGetStreamID(s.ptr)
	return id
}

// SendCue sends one subtitle cue. onSent may be nil to send without a callback.
//
// The error is ErrorCodeSubtitleStreamNotPublishing when the local participant
// is not publishing this stream, either because it is not publishing at all or
// because the StageStrategy does not list this stream in
// PublishStreams.SubtitleStreams.
// This is the only error returned synchronously; it does not call onSent.
//
// A nil error means the SDK accepted the cue for processing, not that it was
// sent or reached the transport. Everything past that point, success or
// failure, is reported on onSent: the SDK calls it exactly once and does not
// retain it afterward, on the same serialized thread as StageHooks
// and StageStrategy.
//
// Since outcomes are reported asynchronously, onSent may be invoked after the
// local participant stops publishing the stream, including after the stream is
// reported removed through StageHooks. Keep whatever it touches valid
// until Close returns.
//
// Any goroutine may call SendCue.
func (s *LocalSubtitleSource) SendCue(cue SubtitleCue, onSent func(err *Error)) error {
	var raw C.aws_ivs_realtime_subtitle_cue

	cText := C.CString(cue.Text)
	defer C.free(unsafe.Pointer(cText))
	raw.text = cText

	cLanguage := C.CString(cue.LanguageCode)
	defer C.free(unsafe.Pointer(cLanguage))
	raw.language_code = cLanguage

	if cue.TargetParticipantID != "" {
		cTarget := C.CString(cue.TargetParticipantID)
		defer C.free(unsafe.Pointer(cTarget))
		raw.target_participant_id = cTarget
	}

	if cue.LanguageLabel != "" {
		cLabel := C.CString(cue.LanguageLabel)
		defer C.free(unsafe.Pointer(cLabel))
		raw.language_label = cLabel
	}

	raw.start_time_ms = C.int64_t(cue.StartTimeMs)
	raw.end_time_ms = C.int64_t(cue.EndTimeMs)

	// Anchors onSent for dispatchCueSent to resolve. The SDK invokes it at most once and only on
	// the success path below, so a synchronous failure here must delete it itself.
	var handle cgo.Handle
	if onSent != nil {
		handle = cgo.NewHandle(onSent)
	}
	err := rawLocalSubtitleStreamSendCue(s.ptr, unsafe.Pointer(&raw), handle, onSent != nil)
	if err != nil {
		if handle != 0 {
			handle.Delete()
		}
		return err
	}
	return nil
}

// Close releases the stream. A Stage publishing it stops publishing it; the Stage
// keeps publishing everything else. It is safe to call more than once and
// satisfies io.Closer.
//
// Any cue still awaiting an outcome has its onSent invoked before Close returns,
// so no callback is lost with the stream. Those callbacks run on the same
// serialized thread as StageHooks and StageStrategy, so Close must not
// be called from within a callback.
//
// Do not call Close concurrently with any other method on the same value.
func (s *LocalSubtitleSource) Close() error {
	if s.ptr == nil {
		return nil
	}
	rawLocalSubtitleStreamDestroy(s.ptr)
	s.ptr = nil
	return nil
}

func dispatchCueSent(userData unsafe.Pointer, err *Error) {
	handle := cgo.Handle(uintptr(userData))
	defer handle.Delete()
	onSent, ok := handle.Value().(func(*Error))
	if !ok {
		return
	}
	onSent(err)
}
