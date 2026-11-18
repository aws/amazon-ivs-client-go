package realtime

import "testing"

func TestStreamTypeString(t *testing.T) {
	tests := []struct {
		name  string
		value StreamType
		want  string
	}{
		{"empty set", 0, "None"},
		{"single flag", StreamTypeSubtitle, "Subtitle"},
		{"combined flags in declaration order", StreamTypeVideo | StreamTypeAudio, "Audio|Video"},
		{"unknown bit appended as hex", StreamTypeAudio | 8, "Audio|0x8"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.value.String(); got != test.want {
				t.Errorf("StreamType(%d).String() = %q, want %q", uint32(test.value), got, test.want)
			}
		})
	}
}
