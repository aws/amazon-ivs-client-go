package realtime

import (
	"runtime/cgo"
	"sync/atomic"
)

// PCMBuffer carries one delivery of PCM audio. Values handed to a callback
// registered with Stage.SetPCMCallback or AudioStageStream.SetPCMCallback are owned by
// the receiver: Data is a private copy that stays valid and unchanged for as long
// as it is retained.
type PCMBuffer struct {
	// Data holds interleaved PCM samples,
	// len(Data) == NumFrames * NumChannels * BytesPerSample.
	Data []byte
	// NumFrames is the sample count per channel.
	NumFrames int
	// BytesPerSample is 2 for signed 16-bit, 4 for 32-bit float.
	BytesPerSample int
	// NumChannels is 1 for mono, 2 for stereo.
	NumChannels int
	// SamplesPerSec is the sample rate in Hz, e.g. 48000.
	SamplesPerSec uint32
	// TimestampMs is the playout timestamp in milliseconds. It shares a timeline
	// with Transcription.StartMs and Transcription.EndMs.
	TimestampMs int64
}

type pcmDelivery struct {
	cb     atomic.Pointer[func(PCMBuffer)]
	handle cgo.Handle
}

type pcmInstaller func(anchor cgo.Handle, install bool) error

// The order here is load-bearing. Each attempt takes a new anchor and retires the
// previous one only after the SDK accepts it, so a rejected call leaves the
// active registration and its anchor untouched.
func (d *pcmDelivery) set(install pcmInstaller, fn func(PCMBuffer)) error {
	if fn == nil && d.handle == 0 {
		return nil
	}
	handle := cgo.NewHandle(d)
	if err := install(handle, fn != nil); err != nil {
		handle.Delete()
		return err
	}
	if d.handle != 0 {
		d.handle.Delete()
	}
	d.handle = handle
	if fn == nil {
		d.cb.Store(nil)
	} else {
		d.cb.Store(&fn)
	}
	return nil
}

func (d *pcmDelivery) deliver(buf PCMBuffer) {
	cb := d.cb.Load()
	if cb == nil {
		return
	}
	buf.Data = append([]byte(nil), buf.Data...)
	(*cb)(buf)
}

func (d *pcmDelivery) release() {
	d.cb.Store(nil)
	if d.handle != 0 {
		d.handle.Delete()
		d.handle = 0
	}
}
