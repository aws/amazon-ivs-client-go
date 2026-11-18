package realtime

import (
	"bytes"
	"testing"
)

// These tests drive pcmDelivery directly, storing the callback without going
// through SetPCMCallback, because registering with the SDK needs a live Stage and
// a valid token. Everything below the registration -- the ownership boundary and
// release -- is exercisable on its own.

func TestDeliverHandsCallbackAnIndependentCopy(t *testing.T) {
	var delivery pcmDelivery
	var got PCMBuffer
	fn := func(buf PCMBuffer) { got = buf }
	delivery.cb.Store(&fn)

	// Stands in for the SDK buffer the trampoline aliases.
	source := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	delivery.deliver(PCMBuffer{
		Data:           source,
		NumFrames:      4,
		BytesPerSample: 2,
		NumChannels:    1,
		SamplesPerSec:  48000,
		TimestampMs:    1234,
	})

	// The SDK reuses its buffer once the callback returns; the caller's copy must
	// not follow it.
	for i := range source {
		source[i] = 0xff
	}

	if want := []byte{1, 2, 3, 4, 5, 6, 7, 8}; !bytes.Equal(got.Data, want) {
		t.Errorf("Data = %v, want %v (delivered buffer must not alias SDK memory)", got.Data, want)
	}
	if got.NumFrames != 4 || got.BytesPerSample != 2 || got.NumChannels != 1 {
		t.Errorf("layout not preserved: %+v", got)
	}
	if got.SamplesPerSec != 48000 || got.TimestampMs != 1234 {
		t.Errorf("metadata not preserved: %+v", got)
	}
}

func TestDeliverWithoutCallbackIsNoOp(t *testing.T) {
	var delivery pcmDelivery
	// A delivery scheduled before a clear may still arrive afterwards, which must
	// be harmless rather than a nil-func panic.
	delivery.deliver(PCMBuffer{Data: []byte{1, 2}, NumFrames: 1, BytesPerSample: 2, NumChannels: 1})
}

func TestDeliverPreservesEmptyBuffer(t *testing.T) {
	var delivery pcmDelivery
	var called bool
	fn := func(buf PCMBuffer) {
		called = true
		if len(buf.Data) != 0 {
			t.Errorf("Data = %v, want empty", buf.Data)
		}
	}
	delivery.cb.Store(&fn)

	delivery.deliver(PCMBuffer{NumFrames: 0})
	if !called {
		t.Error("callback was not invoked for an empty buffer")
	}
}

func TestSetNilOnUnregisteredDeliveryIsNoOp(t *testing.T) {
	var delivery pcmDelivery
	// Clearing a callback that was never set must not try to reach the SDK: a Stage
	// would be required, and teardown clears unconditionally.
	if err := delivery.set(nil, nil); err != nil {
		t.Errorf("set(nil) = %v, want nil", err)
	}
	if delivery.handle != 0 {
		t.Error("clearing an unregistered delivery allocated an anchor")
	}
}

func TestReleaseIsIdempotent(t *testing.T) {
	var delivery pcmDelivery
	fn := func(PCMBuffer) {}
	delivery.cb.Store(&fn)

	delivery.release()
	delivery.release()

	if delivery.cb.Load() != nil {
		t.Error("release did not drop the callback")
	}
	// Delivery after release must be inert: Stage.Close releases while the SDK
	// may still be winding down.
	delivery.deliver(PCMBuffer{Data: []byte{1, 2}})
}
