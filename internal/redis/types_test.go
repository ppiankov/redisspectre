package redis

import "testing"

func TestSeverityRank(t *testing.T) {
	tests := []struct {
		severity Severity
		want     int
	}{
		{SeverityCritical, 4},
		{SeverityHigh, 3},
		{SeverityMedium, 2},
		{SeverityLow, 1},
		{Severity("unknown"), 0},
	}
	for _, tt := range tests {
		if got := SeverityRank(tt.severity); got != tt.want {
			t.Errorf("SeverityRank(%q) = %d, want %d", tt.severity, got, tt.want)
		}
	}
}

func TestMeetsSeverityMin(t *testing.T) {
	tests := []struct {
		s, min Severity
		want   bool
	}{
		{SeverityCritical, SeverityLow, true},
		{SeverityHigh, SeverityHigh, true},
		{SeverityLow, SeverityHigh, false},
		{SeverityMedium, SeverityMedium, true},
		{SeverityLow, SeverityCritical, false},
	}
	for _, tt := range tests {
		if got := MeetsSeverityMin(tt.s, tt.min); got != tt.want {
			t.Errorf("MeetsSeverityMin(%q, %q) = %v, want %v", tt.s, tt.min, got, tt.want)
		}
	}
}

func TestParseSeverity(t *testing.T) {
	tests := []struct {
		input string
		want  Severity
	}{
		{"critical", SeverityCritical},
		{"high", SeverityHigh},
		{"medium", SeverityMedium},
		{"low", SeverityLow},
		{"unknown", SeverityLow},
		{"", SeverityLow},
	}
	for _, tt := range tests {
		if got := ParseSeverity(tt.input); got != tt.want {
			t.Errorf("ParseSeverity(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
