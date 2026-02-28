package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/ppiankov/redisspectre/internal/analyzer"
	"github.com/ppiankov/redisspectre/internal/redis"
)

func TestSpectreHubReporter_Generate(t *testing.T) {
	var buf bytes.Buffer
	r := &SpectreHubReporter{Writer: &buf}

	data := Data{
		Tool:    "redisspectre",
		Version: "0.1.0",
		Target:  Target{Type: "redis", URIHash: "sha256:abc123"},
		Findings: []redis.Finding{
			{ID: redis.FindingBigKey, Severity: redis.SeverityMedium, ResourceType: "Key", ResourceID: "big-hash"},
		},
		Summary: analyzer.Summary{TotalFindings: 1, BySeverity: map[string]int{"medium": 1}},
	}

	if err := r.Generate(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if envelope["schema"] != "spectre/v1" {
		t.Errorf("expected schema 'spectre/v1', got %v", envelope["schema"])
	}
	if envelope["tool"] != "redisspectre" {
		t.Errorf("expected tool 'redisspectre', got %v", envelope["tool"])
	}
}

func TestSpectreHubReporter_EmptyFindings(t *testing.T) {
	var buf bytes.Buffer
	r := &SpectreHubReporter{Writer: &buf}

	data := Data{
		Tool:    "redisspectre",
		Summary: analyzer.Summary{BySeverity: map[string]int{}, ByResourceType: map[string]int{}, ByFindingID: map[string]int{}},
	}

	if err := r.Generate(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}
