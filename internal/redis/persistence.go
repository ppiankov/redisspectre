package redis

import (
	"context"
	"fmt"
)

// PersistenceScanner audits Redis persistence configuration.
type PersistenceScanner struct{}

func (s *PersistenceScanner) Name() string { return "persistence" }

func (s *PersistenceScanner) Audit(ctx context.Context, client RedisClient, cfg AuditConfig) ([]Finding, error) {
	var findings []Finding

	saveConfig, err := client.ConfigGet(ctx, "save")
	if err != nil {
		return nil, fmt.Errorf("config get save: %w", err)
	}

	aofConfig, err := client.ConfigGet(ctx, "appendonly")
	if err != nil {
		return nil, fmt.Errorf("config get appendonly: %w", err)
	}

	saveValue := saveConfig["save"]
	appendOnly := aofConfig["appendonly"]

	rdbDisabled := saveValue == ""
	aofDisabled := appendOnly != "yes"

	if rdbDisabled && aofDisabled {
		findings = append(findings, Finding{
			ID:           FindingNoPersistence,
			Severity:     SeverityHigh,
			ResourceType: "Config",
			ResourceID:   cfg.Addr,
			Message:      "no persistence configured: both RDB and AOF are disabled",
			Metadata: map[string]any{
				"save_config": saveValue,
				"appendonly":  appendOnly,
			},
		})
	}

	return findings, nil
}
