package redis

import (
	"context"
	"testing"
	"time"
)

func TestSlowLogScanner_Name(t *testing.T) {
	s := &SlowLogScanner{}
	if s.Name() != "slowlog" {
		t.Errorf("expected name 'slowlog', got %q", s.Name())
	}
}

func TestSlowLogScanner_FindsSlowCommands(t *testing.T) {
	mock := newMockClient()
	mock.slowLog = []SlowLogEntry{
		{ID: 1, Time: time.Now(), Duration: 50 * time.Millisecond, Args: []string{"KEYS", "*"}},
		{ID: 2, Time: time.Now(), Duration: 5 * time.Millisecond, Args: []string{"GET", "key1"}},
		{ID: 3, Time: time.Now(), Duration: 25 * time.Millisecond, Args: []string{"HGETALL", "big-hash"}},
	}

	s := &SlowLogScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{Addr: "localhost:6379"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 2 {
		t.Fatalf("expected 2 slow command findings, got %d", len(findings))
	}
	for _, f := range findings {
		if f.ID != FindingSlowCommand {
			t.Errorf("expected finding ID %q, got %q", FindingSlowCommand, f.ID)
		}
		if f.Severity != SeverityMedium {
			t.Errorf("expected severity medium, got %q", f.Severity)
		}
	}
}

func TestSlowLogScanner_NoSlowCommands(t *testing.T) {
	mock := newMockClient()
	mock.slowLog = []SlowLogEntry{
		{ID: 1, Time: time.Now(), Duration: 1 * time.Millisecond, Args: []string{"GET", "key1"}},
		{ID: 2, Time: time.Now(), Duration: 5 * time.Millisecond, Args: []string{"SET", "key2", "val"}},
	}

	s := &SlowLogScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}

func TestSlowLogScanner_EmptyLog(t *testing.T) {
	mock := newMockClient()
	mock.slowLog = nil

	s := &SlowLogScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}

func TestSlowLogScanner_TruncatesLongCommands(t *testing.T) {
	longArgs := make([]string, 50)
	for i := range longArgs {
		longArgs[i] = "arg-with-some-content"
	}

	mock := newMockClient()
	mock.slowLog = []SlowLogEntry{
		{ID: 1, Time: time.Now(), Duration: 20 * time.Millisecond, Args: longArgs},
	}

	s := &SlowLogScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{Addr: "localhost:6379"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if len(findings[0].Message) > 200 {
		t.Errorf("expected truncated message, got length %d", len(findings[0].Message))
	}
}
