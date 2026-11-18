//go:build cgo

package realtime

/*
#include <aws_ivs_realtime.h>
#include <stdlib.h>
#include <stdint.h>

// cHandleToPtr widens a cgo.Handle into the void* user_data this ABI carries.
// The cast is done in C so that go vet's unsafeptr check stays enabled for
// the Go side of the package.
static inline void* cHandleToPtr(uintptr_t handle) { return (void*)handle; }

// Aliases carrying const in the type name. A Go signature cannot spell const,
// so the trampolines below name these instead, which is what makes the
// prototypes cgo generates for them match the SDK's callback types.
typedef aws_ivs_realtime_stage_stream* aws_ivs_realtime_stage_stream_ptr;
typedef const aws_ivs_realtime_error const_aws_ivs_realtime_error;
typedef const aws_ivs_realtime_participant_info const_aws_ivs_realtime_participant_info;
typedef const aws_ivs_realtime_stage_stream_ptr const_aws_ivs_realtime_stage_stream_ptr;
typedef const uint8_t const_uint8_t;

// The Go trampolines, declared so this file can take their addresses.
extern void goCueSentCb(const_aws_ivs_realtime_error*, void*);
extern void goAudioCb(const_uint8_t*, size_t, size_t, size_t, size_t, uint32_t, int64_t, void*);
extern void goConnectionStateChangedCb(aws_ivs_realtime_connection_state, const_aws_ivs_realtime_error*, void*);
extern void goParticipantJoinedCb(const_aws_ivs_realtime_participant_info*, void*);
extern void goParticipantLeftCb(const_aws_ivs_realtime_participant_info*, void*);
extern void goParticipantMetadataUpdatedCb(const_aws_ivs_realtime_participant_info*, void*);
extern void goParticipantPublishStateChangedCb(const_aws_ivs_realtime_participant_info*, aws_ivs_realtime_publish_state, void*);
extern void goParticipantSubscribeStateChangedCb(const_aws_ivs_realtime_participant_info*, aws_ivs_realtime_subscribe_state, void*);
extern void goStreamsAddedCb(const_aws_ivs_realtime_participant_info*, const_aws_ivs_realtime_stage_stream_ptr*, size_t, void*);
extern void goStreamsRemovedCb(const_aws_ivs_realtime_participant_info*, const_aws_ivs_realtime_stage_stream_ptr*, size_t, void*);
extern void goStreamsMutedChangedCb(const_aws_ivs_realtime_participant_info*, const_aws_ivs_realtime_stage_stream_ptr*, size_t, void*);
extern void goErrorCb(const_aws_ivs_realtime_error*, void*);
extern aws_ivs_realtime_stream_type goSubscriptionTypesCb(const_aws_ivs_realtime_participant_info*, void*);
extern bool goShouldPublishCb(const_aws_ivs_realtime_participant_info*, void*);
extern void goStreamsToPublishCb(const_aws_ivs_realtime_participant_info*, aws_ivs_realtime_publish_streams*, void*);

// The address of each trampoline, as the SDK's callback type. Taking it
// requires C: a Go func value is not a C function pointer.
static inline aws_ivs_realtime_cue_sent_cb cGetCueSentCb(void) { return goCueSentCb; }
static inline aws_ivs_realtime_audio_cb cGetAudioCb(void) { return goAudioCb; }
static inline aws_ivs_realtime_connection_state_changed_cb cGetConnectionStateChangedCb(void) { return goConnectionStateChangedCb; }
static inline aws_ivs_realtime_participant_joined_cb cGetParticipantJoinedCb(void) { return goParticipantJoinedCb; }
static inline aws_ivs_realtime_participant_left_cb cGetParticipantLeftCb(void) { return goParticipantLeftCb; }
static inline aws_ivs_realtime_participant_metadata_updated_cb cGetParticipantMetadataUpdatedCb(void) { return goParticipantMetadataUpdatedCb; }
static inline aws_ivs_realtime_participant_publish_state_changed_cb cGetParticipantPublishStateChangedCb(void) { return goParticipantPublishStateChangedCb; }
static inline aws_ivs_realtime_participant_subscribe_state_changed_cb cGetParticipantSubscribeStateChangedCb(void) { return goParticipantSubscribeStateChangedCb; }
static inline aws_ivs_realtime_streams_added_cb cGetStreamsAddedCb(void) { return goStreamsAddedCb; }
static inline aws_ivs_realtime_streams_removed_cb cGetStreamsRemovedCb(void) { return goStreamsRemovedCb; }
static inline aws_ivs_realtime_streams_muted_changed_cb cGetStreamsMutedChangedCb(void) { return goStreamsMutedChangedCb; }
static inline aws_ivs_realtime_error_cb cGetErrorCb(void) { return goErrorCb; }
static inline aws_ivs_realtime_subscription_types_cb cGetSubscriptionTypesCb(void) { return goSubscriptionTypesCb; }
static inline aws_ivs_realtime_should_publish_cb cGetShouldPublishCb(void) { return goShouldPublishCb; }
static inline aws_ivs_realtime_streams_to_publish_cb cGetStreamsToPublishCb(void) { return goStreamsToPublishCb; }
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

func rawAudioSourceCreate(samplesPerSec uint32, numChannels int, format PCMFormat) (unsafe.Pointer, error) {
	var out *C.aws_ivs_realtime_audio_source
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_audio_source_create(C.uint32_t(samplesPerSec), C.size_t(numChannels), C.aws_ivs_realtime_pcm_format(format), &out, &cErr)) {
		return nil, takeError(cErr)
	}
	return unsafe.Pointer(out), nil
}

func rawAudioSourcePushBuffer(ptr unsafe.Pointer, data []byte, numFrames int, pts int64) error {
	var dataPtr *C.uint8_t
	if len(data) > 0 {
		dataPtr = (*C.uint8_t)(unsafe.Pointer(&data[0]))
	}
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_audio_source_push_buffer((*C.aws_ivs_realtime_audio_source)(ptr), dataPtr, C.size_t(len(data)), C.size_t(numFrames), C.int64_t(pts), &cErr)) {
		return takeError(cErr)
	}
	return nil
}

func rawAudioSourceDestroy(ptr unsafe.Pointer) {
	C.aws_ivs_realtime_audio_source_destroy((*C.aws_ivs_realtime_audio_source)(ptr))
}

func rawChannelTranscriberCreate(participantID string) (unsafe.Pointer, error) {
	cParticipantID := C.CString(participantID)
	defer C.free(unsafe.Pointer(cParticipantID))
	var out *C.aws_ivs_realtime_channel_transcriber
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_channel_transcriber_create(cParticipantID, &out, &cErr)) {
		return nil, takeError(cErr)
	}
	return unsafe.Pointer(out), nil
}

func rawChannelTranscriberGetStrategy(ptr unsafe.Pointer) (unsafe.Pointer, error) {
	var out *C.aws_ivs_realtime_stage_strategy
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_channel_transcriber_get_strategy((*C.aws_ivs_realtime_channel_transcriber)(ptr), &out, &cErr)) {
		return nil, takeError(cErr)
	}
	return unsafe.Pointer(out), nil
}

func rawChannelTranscriberWriteTranscription(ptr unsafe.Pointer, startMs int64, endMs int64, languageCode string, text string, utteranceID string) error {
	cLanguageCode := C.CString(languageCode)
	defer C.free(unsafe.Pointer(cLanguageCode))
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	cUtteranceID := C.CString(utteranceID)
	defer C.free(unsafe.Pointer(cUtteranceID))
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_channel_transcriber_write_transcription((*C.aws_ivs_realtime_channel_transcriber)(ptr), C.int64_t(startMs), C.int64_t(endMs), cLanguageCode, cText, cUtteranceID, &cErr)) {
		return takeError(cErr)
	}
	return nil
}

func rawChannelTranscriberDestroy(ptr unsafe.Pointer) {
	C.aws_ivs_realtime_channel_transcriber_destroy((*C.aws_ivs_realtime_channel_transcriber)(ptr))
}

func rawVersion() (string, error) {
	var out *C.char
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_get_version(&out, &cErr)) {
		return "", takeError(cErr)
	}
	return C.GoString(out), nil
}

func rawLocalSubtitleStreamCreate() (unsafe.Pointer, error) {
	var out *C.aws_ivs_realtime_local_subtitle_stream
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_local_subtitle_stream_create(&out, &cErr)) {
		return nil, takeError(cErr)
	}
	return unsafe.Pointer(out), nil
}

func rawLocalSubtitleStreamGetStreamID(ptr unsafe.Pointer) (string, error) {
	var out *C.char
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_local_subtitle_stream_get_stream_id((*C.aws_ivs_realtime_local_subtitle_stream)(ptr), &out, &cErr)) {
		return "", takeError(cErr)
	}
	return C.GoString(out), nil
}

func rawLocalSubtitleStreamSendCue(ptr unsafe.Pointer, cue unsafe.Pointer, handle cgo.Handle, install bool) error {
	callback := C.cGetCueSentCb()
	if !install {
		callback = nil
	}
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_local_subtitle_stream_send_cue((*C.aws_ivs_realtime_local_subtitle_stream)(ptr), (*C.aws_ivs_realtime_subtitle_cue)(cue), callback, C.cHandleToPtr(C.uintptr_t(handle)), &cErr)) {
		return takeError(cErr)
	}
	return nil
}

func rawLocalSubtitleStreamDestroy(ptr unsafe.Pointer) {
	C.aws_ivs_realtime_local_subtitle_stream_destroy((*C.aws_ivs_realtime_local_subtitle_stream)(ptr))
}

func rawStageCreate(token string, stageStrategy unsafe.Pointer, callbacks cgo.Handle) (unsafe.Pointer, error) {
	cToken := C.CString(token)
	defer C.free(unsafe.Pointer(cToken))
	var cCallbacks C.aws_ivs_realtime_stage_event_callbacks
	cCallbacks.size = C.size_t(C.sizeof_aws_ivs_realtime_stage_event_callbacks)
	if callbacks != 0 {
		cCallbacks.user_data = C.cHandleToPtr(C.uintptr_t(callbacks))
		cCallbacks.on_connection_state_changed = C.cGetConnectionStateChangedCb()
		cCallbacks.on_participant_joined = C.cGetParticipantJoinedCb()
		cCallbacks.on_participant_left = C.cGetParticipantLeftCb()
		cCallbacks.on_participant_metadata_updated = C.cGetParticipantMetadataUpdatedCb()
		cCallbacks.on_participant_publish_state_changed = C.cGetParticipantPublishStateChangedCb()
		cCallbacks.on_participant_subscribe_state_changed = C.cGetParticipantSubscribeStateChangedCb()
		cCallbacks.on_streams_added = C.cGetStreamsAddedCb()
		cCallbacks.on_streams_removed = C.cGetStreamsRemovedCb()
		cCallbacks.on_streams_muted_changed = C.cGetStreamsMutedChangedCb()
		cCallbacks.on_error = C.cGetErrorCb()
	}
	var out *C.aws_ivs_realtime_stage
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_stage_create(cToken, (*C.aws_ivs_realtime_stage_strategy)(stageStrategy), &cCallbacks, &out, &cErr)) {
		return nil, takeError(cErr)
	}
	return unsafe.Pointer(out), nil
}

func rawStageSetAudioCallback(ptr unsafe.Pointer, handle cgo.Handle, install bool) error {
	callback := C.cGetAudioCb()
	if !install {
		callback = nil
	}
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_stage_set_audio_callback((*C.aws_ivs_realtime_stage)(ptr), callback, C.cHandleToPtr(C.uintptr_t(handle)), &cErr)) {
		return takeError(cErr)
	}
	return nil
}

func rawStageJoin(ptr unsafe.Pointer) error {
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_stage_join((*C.aws_ivs_realtime_stage)(ptr), &cErr)) {
		return takeError(cErr)
	}
	return nil
}

func rawStageRefreshStrategy(ptr unsafe.Pointer) error {
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_stage_refresh_strategy((*C.aws_ivs_realtime_stage)(ptr), &cErr)) {
		return takeError(cErr)
	}
	return nil
}

func rawStageLeave(ptr unsafe.Pointer) error {
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_stage_leave((*C.aws_ivs_realtime_stage)(ptr), &cErr)) {
		return takeError(cErr)
	}
	return nil
}

func rawStageDestroy(ptr unsafe.Pointer) {
	C.aws_ivs_realtime_stage_destroy((*C.aws_ivs_realtime_stage)(ptr))
}

func rawPublishStreamsSetAudio(ptr unsafe.Pointer, audio unsafe.Pointer) error {
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_publish_streams_set_audio((*C.aws_ivs_realtime_publish_streams)(ptr), (*C.aws_ivs_realtime_audio_source)(audio), &cErr)) {
		return takeError(cErr)
	}
	return nil
}

func rawPublishStreamsAddSubtitleStream(ptr unsafe.Pointer, subtitleStream unsafe.Pointer) error {
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_publish_streams_add_subtitle_stream((*C.aws_ivs_realtime_publish_streams)(ptr), (*C.aws_ivs_realtime_local_subtitle_stream)(subtitleStream), &cErr)) {
		return takeError(cErr)
	}
	return nil
}

func rawStageStrategyCreate(callbacks cgo.Handle) (unsafe.Pointer, error) {
	var cCallbacks C.aws_ivs_realtime_stage_strategy_callbacks
	cCallbacks.size = C.size_t(C.sizeof_aws_ivs_realtime_stage_strategy_callbacks)
	if callbacks != 0 {
		cCallbacks.user_data = C.cHandleToPtr(C.uintptr_t(callbacks))
		cCallbacks.subscription_types = C.cGetSubscriptionTypesCb()
		cCallbacks.should_publish = C.cGetShouldPublishCb()
		cCallbacks.streams_to_publish = C.cGetStreamsToPublishCb()
	}
	var out *C.aws_ivs_realtime_stage_strategy
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_stage_strategy_create(&cCallbacks, &out, &cErr)) {
		return nil, takeError(cErr)
	}
	return unsafe.Pointer(out), nil
}

func rawStageStrategyDestroy(ptr unsafe.Pointer) {
	C.aws_ivs_realtime_stage_strategy_destroy((*C.aws_ivs_realtime_stage_strategy)(ptr))
}

func rawStageStreamClone(source unsafe.Pointer) (unsafe.Pointer, error) {
	var out *C.aws_ivs_realtime_stage_stream
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_stage_stream_clone((*C.aws_ivs_realtime_stage_stream)(source), &out, &cErr)) {
		return nil, takeError(cErr)
	}
	return unsafe.Pointer(out), nil
}

func rawStageStreamGetStreamID(ptr unsafe.Pointer) (string, error) {
	var out *C.char
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_stage_stream_get_stream_id((*C.aws_ivs_realtime_stage_stream)(ptr), &out, &cErr)) {
		return "", takeError(cErr)
	}
	return C.GoString(out), nil
}

func rawStageStreamGetStreamType(ptr unsafe.Pointer) (StreamType, error) {
	var out C.aws_ivs_realtime_stream_type
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_stage_stream_get_stream_type((*C.aws_ivs_realtime_stage_stream)(ptr), &out, &cErr)) {
		return StreamType(0), takeError(cErr)
	}
	return StreamType(out), nil
}

func rawStageStreamGetIsMuted(ptr unsafe.Pointer) (bool, error) {
	var out C.bool
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_stage_stream_get_is_muted((*C.aws_ivs_realtime_stage_stream)(ptr), &out, &cErr)) {
		return false, takeError(cErr)
	}
	return bool(out), nil
}

func rawStageStreamSetAudioCallback(ptr unsafe.Pointer, handle cgo.Handle, install bool) error {
	callback := C.cGetAudioCb()
	if !install {
		callback = nil
	}
	var cErr *C.aws_ivs_realtime_error
	if !bool(C.aws_ivs_realtime_stage_stream_set_audio_callback((*C.aws_ivs_realtime_stage_stream)(ptr), callback, C.cHandleToPtr(C.uintptr_t(handle)), &cErr)) {
		return takeError(cErr)
	}
	return nil
}

func rawStageStreamDestroy(ptr unsafe.Pointer) {
	C.aws_ivs_realtime_stage_stream_destroy((*C.aws_ivs_realtime_stage_stream)(ptr))
}

//export goCueSentCb
func goCueSentCb(
	error *C.const_aws_ivs_realtime_error,
	userData unsafe.Pointer,
) {
	dispatchCueSent(userData, unmarshalError(error))
}

//export goAudioCb
func goAudioCb(
	data *C.const_uint8_t,
	dataSize C.size_t,
	numFrames C.size_t,
	bytesPerSample C.size_t,
	numChannels C.size_t,
	samplesPerSec C.uint32_t,
	timestampMs C.int64_t,
	userData unsafe.Pointer,
) {
	dispatchAudio(userData, unsafe.Slice((*byte)(unsafe.Pointer(data)), int(dataSize)), int(numFrames), int(bytesPerSample), int(numChannels), uint32(samplesPerSec), int64(timestampMs))
}

//export goConnectionStateChangedCb
func goConnectionStateChangedCb(
	state C.aws_ivs_realtime_connection_state,
	error *C.const_aws_ivs_realtime_error,
	userData unsafe.Pointer,
) {
	dispatchConnectionStateChanged(userData, ConnectionState(state), unmarshalError(error))
}

//export goParticipantJoinedCb
func goParticipantJoinedCb(
	participant *C.const_aws_ivs_realtime_participant_info,
	userData unsafe.Pointer,
) {
	dispatchParticipantJoined(userData, unmarshalParticipantInfo(participant))
}

//export goParticipantLeftCb
func goParticipantLeftCb(
	participant *C.const_aws_ivs_realtime_participant_info,
	userData unsafe.Pointer,
) {
	dispatchParticipantLeft(userData, unmarshalParticipantInfo(participant))
}

//export goParticipantMetadataUpdatedCb
func goParticipantMetadataUpdatedCb(
	participant *C.const_aws_ivs_realtime_participant_info,
	userData unsafe.Pointer,
) {
	dispatchParticipantMetadataUpdated(userData, unmarshalParticipantInfo(participant))
}

//export goParticipantPublishStateChangedCb
func goParticipantPublishStateChangedCb(
	participant *C.const_aws_ivs_realtime_participant_info,
	publishState C.aws_ivs_realtime_publish_state,
	userData unsafe.Pointer,
) {
	dispatchParticipantPublishStateChanged(userData, unmarshalParticipantInfo(participant), PublishState(publishState))
}

//export goParticipantSubscribeStateChangedCb
func goParticipantSubscribeStateChangedCb(
	participant *C.const_aws_ivs_realtime_participant_info,
	subscribeState C.aws_ivs_realtime_subscribe_state,
	userData unsafe.Pointer,
) {
	dispatchParticipantSubscribeStateChanged(userData, unmarshalParticipantInfo(participant), SubscribeState(subscribeState))
}

//export goStreamsAddedCb
func goStreamsAddedCb(
	participant *C.const_aws_ivs_realtime_participant_info,
	streams *C.const_aws_ivs_realtime_stage_stream_ptr,
	streamsCount C.size_t,
	userData unsafe.Pointer,
) {
	dispatchStreamsAdded(userData, unmarshalParticipantInfo(participant), stageStreamHandles(streams, streamsCount))
}

//export goStreamsRemovedCb
func goStreamsRemovedCb(
	participant *C.const_aws_ivs_realtime_participant_info,
	streams *C.const_aws_ivs_realtime_stage_stream_ptr,
	streamsCount C.size_t,
	userData unsafe.Pointer,
) {
	dispatchStreamsRemoved(userData, unmarshalParticipantInfo(participant), stageStreamHandles(streams, streamsCount))
}

//export goStreamsMutedChangedCb
func goStreamsMutedChangedCb(
	participant *C.const_aws_ivs_realtime_participant_info,
	streams *C.const_aws_ivs_realtime_stage_stream_ptr,
	streamsCount C.size_t,
	userData unsafe.Pointer,
) {
	dispatchStreamsMutedChanged(userData, unmarshalParticipantInfo(participant), stageStreamHandles(streams, streamsCount))
}

//export goErrorCb
func goErrorCb(
	error *C.const_aws_ivs_realtime_error,
	userData unsafe.Pointer,
) {
	dispatchError(userData, unmarshalError(error))
}

//export goSubscriptionTypesCb
func goSubscriptionTypesCb(
	participant *C.const_aws_ivs_realtime_participant_info,
	userData unsafe.Pointer,
) C.aws_ivs_realtime_stream_type {
	return C.aws_ivs_realtime_stream_type(dispatchSubscriptionTypes(userData, unmarshalParticipantInfo(participant)))
}

//export goShouldPublishCb
func goShouldPublishCb(
	participant *C.const_aws_ivs_realtime_participant_info,
	userData unsafe.Pointer,
) C.bool {
	return C.bool(dispatchShouldPublish(userData, unmarshalParticipantInfo(participant)))
}

//export goStreamsToPublishCb
func goStreamsToPublishCb(
	participant *C.const_aws_ivs_realtime_participant_info,
	streams *C.aws_ivs_realtime_publish_streams,
	userData unsafe.Pointer,
) {
	dispatchStreamsToPublish(userData, unmarshalParticipantInfo(participant), unsafe.Pointer(streams))
}

// stageStreamHandles copies a borrowed handle array into a Go slice. The
// handles stay borrowed: they are valid only until the callback returns.
func stageStreamHandles(items *C.const_aws_ivs_realtime_stage_stream_ptr, count C.size_t) []unsafe.Pointer {
	if count == 0 {
		return nil
	}
	cItems := unsafe.Slice(items, int(count))
	result := make([]unsafe.Pointer, int(count))
	for i := range result {
		result[i] = unsafe.Pointer(cItems[i])
	}
	return result
}

// unmarshalError builds a Error from the SDK's callback argument, which is
// valid only until the callback returns, so every string is copied.
func unmarshalError(raw *C.aws_ivs_realtime_error) *Error {
	if raw == nil {
		return nil
	}
	return &Error{
		Code:         ErrorCode(raw.code),
		Type:         ErrorType(raw._type),
		PlatformCode: int32(raw.platform_code),
		IsFatal:      bool(raw.is_fatal),
		Source:       C.GoString(raw.source),
		Message:      C.GoString(raw.message),
	}
}

// unmarshalParticipantInfo builds a ParticipantInfo from the SDK's callback
// argument, which is valid only until the callback returns, so every string is
// copied.
func unmarshalParticipantInfo(raw *C.aws_ivs_realtime_participant_info) ParticipantInfo {
	var attributes map[string]string
	if raw.attributes_count > 0 {
		attributes = make(map[string]string, raw.attributes_count)
		keys := unsafe.Slice(raw.attributes_keys, int(raw.attributes_count))
		values := unsafe.Slice(raw.attributes_values, int(raw.attributes_count))
		for i := range keys {
			attributes[C.GoString(keys[i])] = C.GoString(values[i])
		}
	}
	return ParticipantInfo{
		ParticipantID: C.GoString(raw.participant_id),
		UserID:        C.GoString(raw.user_id),
		IsLocal:       bool(raw.is_local),
		Attributes:    attributes,
	}
}

// takeError copies an error the SDK returned through out_error, then frees it.
func takeError(raw *C.aws_ivs_realtime_error) error {
	defer C.aws_ivs_realtime_error_destroy(raw)
	return unmarshalError(raw)
}
