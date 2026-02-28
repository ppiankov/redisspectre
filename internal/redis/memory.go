package redis

import (
	"context"
	"fmt"
	"strconv"
)

// MemoryScanner audits Redis memory for high fragmentation.
type MemoryScanner struct{}

func (s *MemoryScanner) Name() string { return "memory" }

func (s *MemoryScanner) Audit(ctx context.Context, client RedisClient, cfg AuditConfig) ([]Finding, error) {
	var findings []Finding

	raw, err := client.Info(ctx, "memory")
	if err != nil {
		return nil, fmt.Errorf("info memory: %w", err)
	}

	info := ParseInfo(raw)

	fragRatioStr := info["mem_fragmentation_ratio"]
	if fragRatioStr == "" {
		return findings, nil
	}

	fragRatio, err := strconv.ParseFloat(fragRatioStr, 64)
	if err != nil {
		return nil, fmt.Errorf("parse fragmentation ratio: %w", err)
	}

	if fragRatio > 1.5 {
		usedMemory, _ := strconv.ParseInt(info["used_memory"], 10, 64)
		usedMemoryRSS, _ := strconv.ParseInt(info["used_memory_rss"], 10, 64)

		findings = append(findings, Finding{
			ID:           FindingHighFragmentation,
			Severity:     SeverityHigh,
			ResourceType: "Redis",
			ResourceID:   cfg.Addr,
			Message:      fmt.Sprintf("memory fragmentation ratio %.2f exceeds threshold 1.5", fragRatio),
			Metadata: map[string]any{
				"fragmentation_ratio":   fragRatio,
				"used_memory":           usedMemory,
				"used_memory_human":     FormatBytes(usedMemory),
				"used_memory_rss":       usedMemoryRSS,
				"used_memory_rss_human": FormatBytes(usedMemoryRSS),
			},
		})
	}

	return findings, nil
}
