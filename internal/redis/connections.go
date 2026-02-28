package redis

import (
	"context"
	"fmt"
	"strconv"
)

// ConnectionScanner audits Redis for connection waste indicators.
type ConnectionScanner struct{}

func (s *ConnectionScanner) Name() string { return "connections" }

func (s *ConnectionScanner) Audit(ctx context.Context, client RedisClient, cfg AuditConfig) ([]Finding, error) {
	var findings []Finding

	clientsRaw, err := client.Info(ctx, "clients")
	if err != nil {
		return nil, fmt.Errorf("info clients: %w", err)
	}
	clientsInfo := ParseInfo(clientsRaw)

	statsRaw, err := client.Info(ctx, "stats")
	if err != nil {
		return nil, fmt.Errorf("info stats: %w", err)
	}
	statsInfo := ParseInfo(statsRaw)

	connectedClients, _ := strconv.ParseInt(clientsInfo["connected_clients"], 10, 64)
	blockedClients, _ := strconv.ParseInt(clientsInfo["blocked_clients"], 10, 64)
	rejectedConnections, _ := strconv.ParseInt(statsInfo["rejected_connections"], 10, 64)

	if rejectedConnections > 0 {
		findings = append(findings, Finding{
			ID:           FindingConnectionWaste,
			Severity:     SeverityLow,
			ResourceType: "Redis",
			ResourceID:   cfg.Addr,
			Message:      fmt.Sprintf("%d rejected connections detected", rejectedConnections),
			Metadata: map[string]any{
				"connected_clients":    connectedClients,
				"blocked_clients":      blockedClients,
				"rejected_connections": rejectedConnections,
			},
		})
	}

	return findings, nil
}
