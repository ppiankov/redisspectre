package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/ppiankov/redisspectre/internal/analyzer"
	"github.com/ppiankov/redisspectre/internal/redis"
)

func TestJSONReporter_Generate(t *testing.T) {
	var buf bytes.Buffer
	r := &JSONReporter{Writer: &buf}

	data := Data{
		Tool:      "redisspectre",
		Version:   "0.1.0",
		Timestamp: time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC),
		Target:    Target{Type: "redis", URIHash: "sha256:abc123"},
		Findings: []redis.Finding{
			{ID: redis.FindingHighFragmentation, Severity: redis.SeverityHigh, ResourceType: "Redis", ResourceID: "localhost:6379"},
		},
		Summary: analyzer.Summary{TotalFindings: 1, BySeverity: map[string]int{"high": 1}},
	}

	if err := r.Generate(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"$schema": "spectre/v1"`) {
		t.Errorf("expected spectre/v1 schema in JSON output")
	}
	if !strings.Contains(output, `"tool": "redisspectre"`) {
		t.Errorf("expected tool name in JSON output")
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if envelope["$schema"] != "spectre/v1" {
		t.Errorf("expected $schema spectre/v1, got %v", envelope["$schema"])
	}
}

func TestJSONReporter_EmptyFindings(t *testing.T) {
	var buf bytes.Buffer
	r := &JSONReporter{Writer: &buf}

	data := Data{
		Tool:    "redisspectre",
		Summary: analyzer.Summary{BySeverity: map[string]int{}, ByResourceType: map[string]int{}, ByFindingID: map[string]int{}},
	}

	if err := r.Generate(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
}
