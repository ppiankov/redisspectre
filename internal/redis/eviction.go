package redis

import (
	"context"
	"fmt"
	"strconv"
)

// EvictionScanner audits Redis eviction policy configuration.
type EvictionScanner struct{}

func (s *EvictionScanner) Name() string { return "eviction" }

func (s *EvictionScanner) Audit(ctx context.Context, client RedisClient, cfg AuditConfig) ([]Finding, error) {
	var findings []Finding

	policyConfig, err := client.ConfigGet(ctx, "maxmemory-policy")
	if err != nil {
		return nil, fmt.Errorf("config get maxmemory-policy: %w", err)
	}
	policy := policyConfig["maxmemory-policy"]

	raw, err := client.Info(ctx, "memory")
	if err != nil {
		return nil, fmt.Errorf("info memory: %w", err)
	}
	info := ParseInfo(raw)

	usedMemory, _ := strconv.ParseInt(info["used_memory"], 10, 64)
	maxMemory, _ := strconv.ParseInt(info["maxmemory"], 10, 64)

	if policy == "noeviction" && maxMemory > 0 {
		usagePercent := float64(usedMemory) / float64(maxMemory) * 100

		if usagePercent > 80 {
			findings = append(findings, Finding{
				ID:           FindingEvictionRisk,
				Severity:     SeverityCritical,
				ResourceType: "Config",
				ResourceID:   cfg.Addr,
				Message:      fmt.Sprintf("noeviction policy with %.1f%% memory usage (used: %s, max: %s)", usagePercent, FormatBytes(usedMemory), FormatBytes(maxMemory)),
				Metadata: map[string]any{
					"policy":        policy,
					"used_memory":   usedMemory,
					"maxmemory":     maxMemory,
					"usage_percent": usagePercent,
				},
			})
		}
	}

	return findings, nil
}
