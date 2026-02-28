package redis

import (
	"context"
	"testing"
)

func TestConnectionScanner_Name(t *testing.T) {
	s := &ConnectionScanner{}
	if s.Name() != "connections" {
		t.Errorf("expected name 'connections', got %q", s.Name())
	}
}

func TestConnectionScanner_RejectedConnections(t *testing.T) {
	mock := newMockClient()
	mock.infoResponses["clients"] = "# Clients\nconnected_clients:50\nblocked_clients:3\n"
	mock.infoResponses["stats"] = "# Stats\nrejected_connections:15\n"

	s := &ConnectionScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{Addr: "localhost:6379"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].ID != FindingConnectionWaste {
		t.Errorf("expected finding ID %q, got %q", FindingConnectionWaste, findings[0].ID)
	}
	if findings[0].Severity != SeverityLow {
		t.Errorf("expected severity low, got %q", findings[0].Severity)
	}
}

func TestConnectionScanner_NoRejected(t *testing.T) {
	mock := newMockClient()
	mock.infoResponses["clients"] = "# Clients\nconnected_clients:10\nblocked_clients:0\n"
	mock.infoResponses["stats"] = "# Stats\nrejected_connections:0\n"

	s := &ConnectionScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}

func TestConnectionScanner_MissingStats(t *testing.T) {
	mock := newMockClient()
	mock.infoResponses["clients"] = "# Clients\nconnected_clients:10\n"
	mock.infoResponses["stats"] = ""

	s := &ConnectionScanner{}
	findings, err := s.Audit(context.Background(), mock, AuditConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings when stats missing, got %d", len(findings))
	}
}
