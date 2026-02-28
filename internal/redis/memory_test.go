package redis

import (
	"context"
	"testing"
)

func TestMemoryScanner_Name(t *testing.T) {
	s := &MemoryScanner{}
	if s.Name() != "memory" {
		t.Errorf("expected name 'memory', got %q", s.Name())
	}
}

func TestMemoryScanner_HighFragmentation(t *testing.T) {
	mock := newMockClient()
	mock.infoResponses["memory"] = "# Memory\nmem_fragmentation_ratio:2.50\nused_memory:1048576\nused_memory_rss:2621440\n"

	s := &MemoryScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{Addr: "localhost:6379"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].ID != FindingHighFragmentation {
		t.Errorf("expected finding ID %q, got %q", FindingHighFragmentation, findings[0].ID)
	}
	if findings[0].Severity != SeverityHigh {
		t.Errorf("expected severity high, got %q", findings[0].Severity)
	}
}

func TestMemoryScanner_NormalFragmentation(t *testing.T) {
	mock := newMockClient()
	mock.infoResponses["memory"] = "# Memory\nmem_fragmentation_ratio:1.10\nused_memory:1048576\nused_memory_rss:1153434\n"

	s := &MemoryScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}

func TestMemoryScanner_MissingRatio(t *testing.T) {
	mock := newMockClient()
	mock.infoResponses["memory"] = "# Memory\nused_memory:1048576\n"

	s := &MemoryScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for missing ratio, got %d", len(findings))
	}
}
