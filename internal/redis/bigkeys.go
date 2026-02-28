package redis

import (
	"context"
	"fmt"
)

const defaultBigKeySize = 10 * 1024 * 1024 // 10 MB

// BigKeyScanner audits Redis keys for excessive memory usage.
type BigKeyScanner struct{}

func (s *BigKeyScanner) Name() string { return "big_keys" }

func (s *BigKeyScanner) Audit(ctx context.Context, client RedisClient, cfg AuditConfig) ([]Finding, error) {
	var findings []Finding

	bigKeySize := cfg.BigKeySize
	if bigKeySize <= 0 {
		bigKeySize = defaultBigKeySize
	}

	sampleSize := cfg.SampleSize
	if sampleSize <= 0 {
		sampleSize = 10000
	}

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
			mem, err := client.MemoryUsage(ctx, key)
			if err != nil {
				continue
			}

			if mem > bigKeySize {
				findings = append(findings, Finding{
					ID:           FindingBigKey,
					Severity:     SeverityMedium,
					ResourceType: "Key",
					ResourceID:   key,
					Message:      fmt.Sprintf("key %q uses %s (threshold: %s)", key, FormatBytes(mem), FormatBytes(bigKeySize)),
					Metadata: map[string]any{
						"key":             key,
						"size_bytes":      mem,
						"size_human":      FormatBytes(mem),
						"threshold_bytes": bigKeySize,
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
