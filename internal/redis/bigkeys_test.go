package redis

import (
	"context"
	"testing"
)

func TestBigKeyScanner_Name(t *testing.T) {
	s := &BigKeyScanner{}
	if s.Name() != "big_keys" {
		t.Errorf("expected name 'big_keys', got %q", s.Name())
	}
}

func TestBigKeyScanner_FindsBigKeys(t *testing.T) {
	mock := newMockClient()
	mock.scanKeys = []string{"small", "big", "huge"}
	mock.memoryUsages = map[string]int64{
		"small": 1024,              // 1 KB
		"big":   15 * 1024 * 1024,  // 15 MB
		"huge":  100 * 1024 * 1024, // 100 MB
	}

	s := &BigKeyScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{BigKeySize: 10 * 1024 * 1024})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 2 {
		t.Fatalf("expected 2 big key findings, got %d", len(findings))
	}
	for _, f := range findings {
		if f.ID != FindingBigKey {
			t.Errorf("expected finding ID %q, got %q", FindingBigKey, f.ID)
		}
		if f.Severity != SeverityMedium {
			t.Errorf("expected severity medium, got %q", f.Severity)
		}
	}
}

func TestBigKeyScanner_NoBigKeys(t *testing.T) {
	mock := newMockClient()
	mock.scanKeys = []string{"key1", "key2"}
	mock.memoryUsages = map[string]int64{
		"key1": 1024,
		"key2": 2048,
	}

	s := &BigKeyScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}

func TestBigKeyScanner_EmptyDatabase(t *testing.T) {
	mock := newMockClient()

	s := &BigKeyScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}
