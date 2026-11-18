// Package realtime is the Go binding of the IVS Linux Broadcast SDK. It joins
// a Stage as a headless participant, subscribes to remote participants' audio,
// and delivers that audio to your program as Go values.
//
// A session has four steps: build a [StageStrategy] deciding what to subscribe to,
// create a Stage with it, register a PCM callback, and join. The package
// example below shows the flow.
//
// # Threading
//
// Stage events arrive through [StageHooks] on one serialized thread, so hook
// state needs no locking unless something outside the hooks reads it. PCM audio
// arrives on a separate audio thread, so a PCM callback must return promptly;
// buffers it receives are owned by the caller and may be retained freely.
//
// No method waits on the network, so none takes a context.Context. Close
// waits for in-flight callbacks to finish. Cancellation belongs in
// whatever goroutine you hand work to.
//
// # Errors
//
// Every error the SDK returns or delivers is an [*Error]. A method's doc names
// the error codes a caller can handle specifically. That list is not
// exhaustive. For any other error, [Error.IsFatal] says whether the Stage can
// still be used; [Error.Code] and [Error.Message] identify the failure when
// reporting it.
//
// Use [errors.As] to reach one, including through wrapping:
//
//	var sdkErr *realtime.Error
//	if errors.As(err, &sdkErr) && sdkErr.Code == realtime.ErrorCodeStageExpiredToken {
//		...
//	}
//
// # Lifetimes
//
// Everything you create is released with Close, which is idempotent and may be
// called in any order. A transcriber outlives the Stage built from its
// Strategy: close the Stage first, then the transcriber.
//
// Close is not safe to call concurrently with any other method on the same
// value, including methods documented as safe from any goroutine. Make sure
// every other call has returned before you call Close.
//
// # Requirements
//
// Linux only, on x86_64 or arm64. The package links against
// libaws-ivs-realtime, so it requires cgo and does not build with
// CGO_ENABLED=0. See the README for how to obtain and place the shared library.
package realtime
