package report

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/ppiankov/redisspectre/internal/analyzer"
	"github.com/ppiankov/redisspectre/internal/redis"
)

func TestTextReporter_NoFindings(t *testing.T) {
	var buf bytes.Buffer
	r := &TextReporter{Writer: &buf}

	data := Data{
		Tool:      "redisspectre",
		Version:   "0.1.0",
		Timestamp: time.Now(),
		Summary:   analyzer.Summary{TotalFindings: 0},
	}

	if err := r.Generate(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "No issues found") {
		t.Errorf("expected 'No issues found' in output, got: %s", output)
	}
	if !strings.Contains(output, "redisspectre") {
		t.Errorf("expected 'redisspectre' in header")
	}
}

func TestTextReporter_WithFindings(t *testing.T) {
	var buf bytes.Buffer
	r := &TextReporter{Writer: &buf}

	data := Data{
		Tool:    "redisspectre",
		Version: "0.1.0",
		Findings: []redis.Finding{
			{ID: redis.FindingHighFragmentation, Severity: redis.SeverityHigh, ResourceType: "Redis", ResourceID: "localhost:6379", Message: "fragmentation ratio 2.5"},
			{ID: redis.FindingIdleKey, Severity: redis.SeverityMedium, ResourceType: "Key", ResourceID: "cache:old", Message: "idle for 45 days"},
		},
		Summary: analyzer.Summary{
			TotalFindings:         2,
			TotalResourcesScanned: 50,
			BySeverity:            map[string]int{"high": 1, "medium": 1},
			ByResourceType:        map[string]int{"Redis": 1, "Key": 1},
		},
	}

	if err := r.Generate(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Found 2 issues") {
		t.Errorf("expected 'Found 2 issues' in output")
	}
	if !strings.Contains(output, "HIGH_FRAGMENTATION") {
		t.Errorf("expected finding ID in output")
	}
	if !strings.Contains(output, "high=1") {
		t.Errorf("expected severity breakdown in output")
	}
}

func TestTextReporter_WithErrors(t *testing.T) {
	var buf bytes.Buffer
	r := &TextReporter{Writer: &buf}

	data := Data{
		Summary: analyzer.Summary{TotalFindings: 0},
		Errors:  []string{"memory: timeout"},
	}

	if err := r.Generate(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Warnings (1)") {
		t.Errorf("expected warnings section in output")
	}
	if !strings.Contains(output, "memory: timeout") {
		t.Errorf("expected error message in output")
	}
}
