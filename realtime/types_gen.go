package realtime

import (
	"fmt"
	"strings"
)

// ConnectionState is the connection state of a Stage.
type ConnectionState int32

const (
	ConnectionStateDisconnected ConnectionState = 0 // Not connected to a Stage.
	ConnectionStateConnecting   ConnectionState = 1 // Establishing a connection.
	ConnectionStateConnected    ConnectionState = 2 // Connected.
)

func (v ConnectionState) String() string {
	switch v {
	case ConnectionStateDisconnected:
		return "Disconnected"
	case ConnectionStateConnecting:
		return "Connecting"
	case ConnectionStateConnected:
		return "Connected"
	}
	return fmt.Sprintf("ConnectionState(%d)", int32(v))
}

// PublishState is the publish state of a participant.
type PublishState int32

const (
	PublishStateNotPublished      PublishState = 0 // Not publishing and not attempting to.
	PublishStateAttemptingPublish PublishState = 1 // Attempting to publish.
	PublishStatePublished         PublishState = 2 // Publishing.
)

func (v PublishState) String() string {
	switch v {
	case PublishStateNotPublished:
		return "NotPublished"
	case PublishStateAttemptingPublish:
		return "AttemptingPublish"
	case PublishStatePublished:
		return "Published"
	}
	return fmt.Sprintf("PublishState(%d)", int32(v))
}

// SubscribeState is the local participant's subscribe state with respect to a
// remote publisher.
type SubscribeState int32

const (
	SubscribeStateNotSubscribed       SubscribeState = 0 // Not subscribed.
	SubscribeStateAttemptingSubscribe SubscribeState = 1 // Establishing a subscription.
	SubscribeStateSubscribed          SubscribeState = 2 // Subscribed.
)

func (v SubscribeState) String() string {
	switch v {
	case SubscribeStateNotSubscribed:
		return "NotSubscribed"
	case SubscribeStateAttemptingSubscribe:
		return "AttemptingSubscribe"
	case SubscribeStateSubscribed:
		return "Subscribed"
	}
	return fmt.Sprintf("SubscribeState(%d)", int32(v))
}

// StreamType is the bitmask of stream types. Each stream, published or
// subscribed, reports exactly one bit. A strategy returns the set of stream
// types it wants to subscribe to.
type StreamType uint32

const (
	StreamTypeNone     StreamType = 0 // The empty set. No stream reports this.
	StreamTypeAudio    StreamType = 1 // Audio (bit 0).
	StreamTypeVideo    StreamType = 2 // Video (bit 1).
	StreamTypeSubtitle StreamType = 4 // Subtitles (bit 2).
)

func (v StreamType) String() string {
	if v == 0 {
		return "None"
	}
	var names []string
	if v&StreamTypeAudio != 0 {
		names = append(names, "Audio")
		v &^= StreamTypeAudio
	}
	if v&StreamTypeVideo != 0 {
		names = append(names, "Video")
		v &^= StreamTypeVideo
	}
	if v&StreamTypeSubtitle != 0 {
		names = append(names, "Subtitle")
		v &^= StreamTypeSubtitle
	}
	if v != 0 {
		names = append(names, fmt.Sprintf("0x%x", uint32(v)))
	}
	return strings.Join(names, "|")
}

// PCMFormat is the sample format of raw PCM. Interleaved formats store each
// frame's channels together; planar formats store one whole channel after
// another.
type PCMFormat int32

const (
	PCMFormatS16Interleaved PCMFormat = 0 // Signed 16-bit integer, interleaved.
	PCMFormatS16Planar      PCMFormat = 1 // Signed 16-bit integer, planar.
	PCMFormatF32Interleaved PCMFormat = 2 // 32-bit float, interleaved.
	PCMFormatF32Planar      PCMFormat = 3 // 32-bit float, planar.
	PCMFormatS32Interleaved PCMFormat = 4 // Signed 32-bit integer, interleaved.
	PCMFormatS32Planar      PCMFormat = 5 // Signed 32-bit integer, planar.
	PCMFormatF64Interleaved PCMFormat = 6 // 64-bit float, interleaved.
	PCMFormatF64Planar      PCMFormat = 7 // 64-bit float, planar.
)

func (v PCMFormat) String() string {
	switch v {
	case PCMFormatS16Interleaved:
		return "S16Interleaved"
	case PCMFormatS16Planar:
		return "S16Planar"
	case PCMFormatF32Interleaved:
		return "F32Interleaved"
	case PCMFormatF32Planar:
		return "F32Planar"
	case PCMFormatS32Interleaved:
		return "S32Interleaved"
	case PCMFormatS32Planar:
		return "S32Planar"
	case PCMFormatF64Interleaved:
		return "F64Interleaved"
	case PCMFormatF64Planar:
		return "F64Planar"
	}
	return fmt.Sprintf("PCMFormat(%d)", int32(v))
}

// SubtitleCue is a subtitle cue to publish.
type SubtitleCue struct {
	// The cue text, as UTF-8. Required; must be non-empty.
	Text string
	// The language of the cue as an ISO 639-1 code, for example "en". Required;
	// must be non-empty. Carried through to receivers unmodified and otherwise
	// unvalidated.
	LanguageCode string
	// The participant the cue describes, which is not necessarily the one that
	// authored it. Leave it unset to attribute the cue to the local participant.
	TargetParticipantID string
	// Optional human-readable label for the language, for example "English", which
	// provides a readable passthrough display receivers can use. Max 96 bytes.
	// Carried through to receivers unmodified and otherwise unvalidated.
	LanguageLabel string
	// Start time of the cue, in milliseconds from the NTP capture time of the media
	// it describes. 0 means unspecified. Real-time subtitles usually leave this
	// unspecified, because they render as they arrive. Set the start and the end
	// together: a cue with only one of them is published with neither.
	StartTimeMs int64
	// End time of the cue, on the same timeline as the start. 0 means unspecified.
	EndTimeMs int64
}

// ParticipantInfo is an immutable snapshot of a participant in a Stage.
type ParticipantInfo struct {
	// Unique participant ID within the Stage.
	ParticipantID string
	// User ID associated with the participant. May be empty.
	UserID string
	// True if this is the local participant (the token used to join).
	IsLocal bool
	// Free-form participant attributes as key/value metadata.
	Attributes map[string]string
}
