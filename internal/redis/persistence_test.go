package redis

import (
	"context"
	"testing"
)

func TestPersistenceScanner_Name(t *testing.T) {
	s := &PersistenceScanner{}
	if s.Name() != "persistence" {
		t.Errorf("expected name 'persistence', got %q", s.Name())
	}
}

func TestPersistenceScanner_NoPersistence(t *testing.T) {
	mock := newMockClient()
	mock.configValues["save"] = map[string]string{"save": ""}
	mock.configValues["appendonly"] = map[string]string{"appendonly": "no"}

	s := &PersistenceScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{Addr: "localhost:6379"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].ID != FindingNoPersistence {
		t.Errorf("expected finding ID %q, got %q", FindingNoPersistence, findings[0].ID)
	}
	if findings[0].Severity != SeverityHigh {
		t.Errorf("expected severity high, got %q", findings[0].Severity)
	}
}

func TestPersistenceScanner_RDBEnabled(t *testing.T) {
	mock := newMockClient()
	mock.configValues["save"] = map[string]string{"save": "3600 1 300 100"}
	mock.configValues["appendonly"] = map[string]string{"appendonly": "no"}

	s := &PersistenceScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings (RDB enabled), got %d", len(findings))
	}
}

func TestPersistenceScanner_AOFEnabled(t *testing.T) {
	mock := newMockClient()
	mock.configValues["save"] = map[string]string{"save": ""}
	mock.configValues["appendonly"] = map[string]string{"appendonly": "yes"}

	s := &PersistenceScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings (AOF enabled), got %d", len(findings))
	}
}

func TestPersistenceScanner_BothEnabled(t *testing.T) {
	mock := newMockClient()
	mock.configValues["save"] = map[string]string{"save": "3600 1"}
	mock.configValues["appendonly"] = map[string]string{"appendonly": "yes"}

	s := &PersistenceScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings (both enabled), got %d", len(findings))
	}
}
