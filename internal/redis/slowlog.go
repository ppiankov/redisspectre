package redis

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const slowThreshold = 10 * time.Millisecond

// SlowLogScanner audits Redis slow log for slow commands.
type SlowLogScanner struct{}

func (s *SlowLogScanner) Name() string { return "slowlog" }

func (s *SlowLogScanner) Audit(ctx context.Context, client RedisClient, cfg AuditConfig) ([]Finding, error) {
	var findings []Finding

	entries, err := client.SlowLogGet(ctx, 128)
	if err != nil {
		return nil, fmt.Errorf("slowlog get: %w", err)
	}

	for _, entry := range entries {
		if entry.Duration >= slowThreshold {
			command := strings.Join(entry.Args, " ")
			if len(command) > 100 {
				command = command[:100] + "..."
			}

			findings = append(findings, Finding{
				ID:           FindingSlowCommand,
				Severity:     SeverityMedium,
				ResourceType: "Redis",
				ResourceID:   cfg.Addr,
				Message:      fmt.Sprintf("slow command: %s (%.1f ms)", command, float64(entry.Duration.Microseconds())/1000),
				Metadata: map[string]any{
					"command":     strings.Join(entry.Args, " "),
					"duration_ms": float64(entry.Duration.Microseconds()) / 1000,
					"timestamp":   entry.Time.UTC().Format(time.RFC3339),
					"slowlog_id":  entry.ID,
				},
			})
		}
	}

	return findings, nil
}
