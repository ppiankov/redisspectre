package redis

import (
	"context"
	"testing"
)

func TestEvictionScanner_Name(t *testing.T) {
	s := &EvictionScanner{}
	if s.Name() != "eviction" {
		t.Errorf("expected name 'eviction', got %q", s.Name())
	}
}

func TestEvictionScanner_CriticalRisk(t *testing.T) {
	mock := newMockClient()
	mock.configValues["maxmemory-policy"] = map[string]string{"maxmemory-policy": "noeviction"}
	mock.infoResponses["memory"] = "# Memory\nused_memory:900000\nmaxmemory:1000000\n"

	s := &EvictionScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{Addr: "localhost:6379"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].ID != FindingEvictionRisk {
		t.Errorf("expected finding ID %q, got %q", FindingEvictionRisk, findings[0].ID)
	}
	if findings[0].Severity != SeverityCritical {
		t.Errorf("expected severity critical, got %q", findings[0].Severity)
	}
}

func TestEvictionScanner_NoRiskBelowThreshold(t *testing.T) {
	mock := newMockClient()
	mock.configValues["maxmemory-policy"] = map[string]string{"maxmemory-policy": "noeviction"}
	mock.infoResponses["memory"] = "# Memory\nused_memory:500000\nmaxmemory:1000000\n"

	s := &EvictionScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings (50%% usage), got %d", len(findings))
	}
}

func TestEvictionScanner_NoRiskWithEvictionPolicy(t *testing.T) {
	mock := newMockClient()
	mock.configValues["maxmemory-policy"] = map[string]string{"maxmemory-policy": "allkeys-lru"}
	mock.infoResponses["memory"] = "# Memory\nused_memory:950000\nmaxmemory:1000000\n"

	s := &EvictionScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings (eviction policy set), got %d", len(findings))
	}
}

func TestEvictionScanner_NoMaxmemory(t *testing.T) {
	mock := newMockClient()
	mock.configValues["maxmemory-policy"] = map[string]string{"maxmemory-policy": "noeviction"}
	mock.infoResponses["memory"] = "# Memory\nused_memory:950000\nmaxmemory:0\n"

	s := &EvictionScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings (maxmemory=0 means unlimited), got %d", len(findings))
	}
}
