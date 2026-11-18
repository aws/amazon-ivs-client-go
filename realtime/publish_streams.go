//go:build cgo

package realtime

import "unsafe"

// applyPublishStreams names what a strategy wants published on the SDK's
// collector. A stream left unset is not published.
func applyPublishStreams(streams unsafe.Pointer, value PublishStreams) {
	if value.Audio != nil {
		rawPublishStreamsSetAudio(streams, value.Audio.ptr)
	}
	for _, subtitles := range value.SubtitleStreams {
		if subtitles != nil && subtitles.ptr != nil {
			rawPublishStreamsAddSubtitleStream(streams, subtitles.ptr)
		}
	}
}
