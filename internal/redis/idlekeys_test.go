package redis

import (
	"context"
	"testing"
	"time"
)

func TestIdleKeyScanner_Name(t *testing.T) {
	s := &IdleKeyScanner{}
	if s.Name() != "idle_keys" {
		t.Errorf("expected name 'idle_keys', got %q", s.Name())
	}
}

func TestIdleKeyScanner_FindsIdleKeys(t *testing.T) {
	mock := newMockClient()
	mock.scanKeys = []string{"key1", "key2", "key3"}
	mock.idleTimes = map[string]time.Duration{
		"key1": 60 * 24 * time.Hour, // 60 days idle
		"key2": 5 * 24 * time.Hour,  // 5 days idle
		"key3": 45 * 24 * time.Hour, // 45 days idle
	}

	s := &IdleKeyScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{IdleDays: 30})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 2 {
		t.Fatalf("expected 2 idle key findings, got %d", len(findings))
	}
	for _, f := range findings {
		if f.ID != FindingIdleKey {
			t.Errorf("expected finding ID %q, got %q", FindingIdleKey, f.ID)
		}
		if f.Severity != SeverityMedium {
			t.Errorf("expected severity medium, got %q", f.Severity)
		}
	}
}

func TestIdleKeyScanner_NoIdleKeys(t *testing.T) {
	mock := newMockClient()
	mock.scanKeys = []string{"key1", "key2"}
	mock.idleTimes = map[string]time.Duration{
		"key1": 1 * time.Hour,
		"key2": 2 * time.Hour,
	}

	s := &IdleKeyScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{IdleDays: 30})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}

func TestIdleKeyScanner_EmptyDatabase(t *testing.T) {
	mock := newMockClient()
	mock.scanKeys = nil

	s := &IdleKeyScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}

func TestIdleKeyScanner_DefaultIdleDays(t *testing.T) {
	mock := newMockClient()
	mock.scanKeys = []string{"key1"}
	mock.idleTimes = map[string]time.Duration{
		"key1": 31 * 24 * time.Hour, // 31 days, default threshold is 30
	}

	s := &IdleKeyScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 1 {
		t.Errorf("expected 1 finding with default 30-day threshold, got %d", len(findings))
	}
}
