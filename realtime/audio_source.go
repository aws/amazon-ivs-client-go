package realtime

import "unsafe"

// AudioSource is a PCM audio source. Push audio to it with PushBuffer. Close it
// when done.
type AudioSource struct {
	ptr unsafe.Pointer
}

// AudioSourceConfig describes every buffer pushed to an AudioSource. Pass it to
// NewAudioSource.
type AudioSourceConfig struct {
	// SamplesPerSec is the sample rate in Hz, between 8000 and 96000.
	SamplesPerSec uint32
	// NumChannels is 1 for mono or 2 for stereo.
	NumChannels int
	// Format is the sample format and layout.
	Format PCMFormat
}

// NewAudioSource creates an audio source. A SamplesPerSec or NumChannels outside
// its range returns ErrorCodeInvalidArgument.
func NewAudioSource(config AudioSourceConfig) (*AudioSource, error) {
	ptr, err := rawAudioSourceCreate(config.SamplesPerSec, config.NumChannels, config.Format)
	if err != nil {
		return nil, err
	}
	return &AudioSource{ptr: ptr}, nil
}

// AudioBuffer is one buffer of PCM to push to an AudioSource, in the format its
// AudioSourceConfig describes.
type AudioBuffer struct {
	// Data holds NumFrames frames of PCM.
	Data []byte
	// NumFrames counts frames, each with one sample per channel.
	NumFrames int
	// PTS is the frame index of the buffer's first frame, on the caller's own
	// timeline. The origin is arbitrary. Contiguous audio advances PTS by
	// NumFrames per push. A jump in PTS creates a gap, filled with silence.
	PTS int64
}

// PushBuffer pushes one buffer of PCM. It copies buffer.Data before returning,
// so the caller can reuse the slice right away.
//
// Push in real time, as the audio is produced. PushBuffer does not pace audio,
// so it drops audio pushed faster than real time. A source with no consumer
// discards the push.
//
// len(buffer.Data) must equal buffer.NumFrames * NumChannels * the format's
// bytes per sample, or PushBuffer returns
// ErrorCodeAudioSourceBufferSizeMismatch.
//
// Call PushBuffer from one goroutine at a time. Any goroutine may call it.
func (a *AudioSource) PushBuffer(buffer AudioBuffer) error {
	return rawAudioSourcePushBuffer(a.ptr, buffer.Data, buffer.NumFrames, buffer.PTS)
}

// Close releases the audio source.
//
// You can close a source that a Stage publishes: the stream stays published
// and delivers silence.
//
// Do not call Close concurrently with any other method on the same value.
func (a *AudioSource) Close() error {
	if a.ptr == nil {
		return nil
	}
	rawAudioSourceDestroy(a.ptr)
	a.ptr = nil
	return nil
}
