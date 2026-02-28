package redis

import (
	"context"
	"testing"
)

func TestMultiAuditorEmpty(t *testing.T) {
	multi := NewMultiAuditor(nil, 4)
	result, err := multi.AuditAll(context.Background(), newMockClient(), AuditConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(result.Findings))
	}
}

func TestMultiAuditorCombinesFindings(t *testing.T) {
	mock := newMockClient()
	mock.infoResponses["memory"] = "# Memory\nmem_fragmentation_ratio:2.5\nused_memory:1000\nused_memory_rss:2500\nmaxmemory:0\n"
	mock.configValues["maxmemory-policy"] = map[string]string{"maxmemory-policy": "allkeys-lru"}
	mock.configValues["save"] = map[string]string{"save": ""}
	mock.configValues["appendonly"] = map[string]string{"appendonly": "no"}
	mock.infoResponses["clients"] = "# Clients\nconnected_clients:5\nblocked_clients:0\n"
	mock.infoResponses["stats"] = "# Stats\nrejected_connections:0\n"

	auditors := []Auditor{
		&MemoryScanner{},
		&PersistenceScanner{},
	}

	multi := NewMultiAuditor(auditors, 2)
	result, err := multi.AuditAll(context.Background(), mock, AuditConfig{Addr: "localhost:6379"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Findings) != 2 {
		t.Errorf("expected 2 findings (fragmentation + no persistence), got %d", len(result.Findings))
	}
}

func TestMultiAuditorHandlesErrors(t *testing.T) {
	mock := newMockClient()
	mock.infoErr = context.DeadlineExceeded

	multi := NewMultiAuditor([]Auditor{&MemoryScanner{}}, 1)
	result, err := multi.AuditAll(context.Background(), mock, AuditConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(result.Errors))
	}
}

func TestAllAuditors(t *testing.T) {
	auditors := AllAuditors()
	if len(auditors) != 7 {
		t.Errorf("expected 7 auditors, got %d", len(auditors))
	}
}

func TestNewMultiAuditorDefaultConcurrency(t *testing.T) {
	multi := NewMultiAuditor(nil, 0)
	if multi.concurrency != 4 {
		t.Errorf("expected default concurrency 4, got %d", multi.concurrency)
	}
}
