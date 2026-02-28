package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/ppiankov/redisspectre/internal/analyzer"
	"github.com/ppiankov/redisspectre/internal/redis"
)

func TestSARIFReporter_Generate(t *testing.T) {
	var buf bytes.Buffer
	r := &SARIFReporter{Writer: &buf}

	data := Data{
		Tool:    "redisspectre",
		Version: "0.1.0",
		Target:  Target{Type: "redis", URIHash: "sha256:abc123"},
		Findings: []redis.Finding{
			{ID: redis.FindingHighFragmentation, Severity: redis.SeverityHigh, ResourceType: "Redis", ResourceID: "localhost:6379", Message: "high fragmentation"},
			{ID: redis.FindingIdleKey, Severity: redis.SeverityMedium, ResourceType: "Key", ResourceID: "cache:old", Message: "idle key"},
		},
		Summary: analyzer.Summary{TotalFindings: 2},
	}

	if err := r.Generate(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "sarif-schema-2.1.0") {
		t.Errorf("expected SARIF schema reference")
	}
	if !strings.Contains(output, `"version": "2.1.0"`) {
		t.Errorf("expected SARIF version 2.1.0")
	}

	var rpt map[string]any
	if err := json.Unmarshal(buf.Bytes(), &rpt); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	runs := rpt["runs"].([]any)
	if len(runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(runs))
	}

	run := runs[0].(map[string]any)
	results := run["results"].([]any)
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestSARIFReporter_EmptyFindings(t *testing.T) {
	var buf bytes.Buffer
	r := &SARIFReporter{Writer: &buf}

	data := Data{
		Tool:    "redisspectre",
		Version: "0.1.0",
		Summary: analyzer.Summary{},
	}

	if err := r.Generate(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var rpt map[string]any
	if err := json.Unmarshal(buf.Bytes(), &rpt); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestSarifLevel(t *testing.T) {
	tests := []struct {
		severity redis.Severity
		want     string
	}{
		{redis.SeverityCritical, "error"},
		{redis.SeverityHigh, "error"},
		{redis.SeverityMedium, "warning"},
		{redis.SeverityLow, "note"},
	}
	for _, tt := range tests {
		if got := sarifLevel(tt.severity); got != tt.want {
			t.Errorf("sarifLevel(%q) = %q, want %q", tt.severity, got, tt.want)
		}
	}
}
