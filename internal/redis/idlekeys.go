package redis

import (
	"context"
	"fmt"
	"time"
)

// IdleKeyScanner audits Redis keys for inactivity using OBJECT IDLETIME.
type IdleKeyScanner struct{}

func (s *IdleKeyScanner) Name() string { return "idle_keys" }

func (s *IdleKeyScanner) Audit(ctx context.Context, client RedisClient, cfg AuditConfig) ([]Finding, error) {
	var findings []Finding

	idleDays := cfg.IdleDays
	if idleDays <= 0 {
		idleDays = 30
	}

	sampleSize := cfg.SampleSize
	if sampleSize <= 0 {
		sampleSize = 10000
	}

	threshold := time.Duration(idleDays) * 24 * time.Hour
	sampled := 0
	var cursor uint64

	for sampled < sampleSize {
		batchSize := int64(100)
		if remaining := sampleSize - sampled; remaining < int(batchSize) {
			batchSize = int64(remaining)
		}

		keys, nextCursor, err := client.Scan(ctx, cursor, "*", batchSize)
		if err != nil {
			return nil, fmt.Errorf("scan keys: %w", err)
		}

		for _, key := range keys {
			idle, err := client.ObjectIdleTime(ctx, key)
			if err != nil {
				continue
			}

			if idle >= threshold {
				idleDaysActual := int(idle.Hours() / 24)
				findings = append(findings, Finding{
					ID:           FindingIdleKey,
					Severity:     SeverityMedium,
					ResourceType: "Key",
					ResourceID:   key,
					Message:      fmt.Sprintf("key %q idle for %d days (threshold: %d days)", key, idleDaysActual, idleDays),
					Metadata: map[string]any{
						"key":            key,
						"idle_seconds":   int64(idle.Seconds()),
						"idle_days":      idleDaysActual,
						"threshold_days": idleDays,
					},
				})
			}
		}

		sampled += len(keys)
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return findings, nil
}
