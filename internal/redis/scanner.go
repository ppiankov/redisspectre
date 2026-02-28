package redis

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"golang.org/x/sync/errgroup"
)

// Auditor is the interface each resource-type auditor implements.
type Auditor interface {
	Audit(ctx context.Context, client RedisClient, cfg AuditConfig) ([]Finding, error)
	Name() string
}

// MultiAuditor orchestrates running multiple auditors in parallel.
type MultiAuditor struct {
	auditors    []Auditor
	concurrency int
}

// NewMultiAuditor creates an auditor that runs the specified auditors in parallel.
func NewMultiAuditor(auditors []Auditor, concurrency int) *MultiAuditor {
	if concurrency <= 0 {
		concurrency = 4
	}
	return &MultiAuditor{
		auditors:    auditors,
		concurrency: concurrency,
	}
}

// AuditAll runs all auditors and returns combined results.
func (m *MultiAuditor) AuditAll(ctx context.Context, client RedisClient, cfg AuditConfig) (*ScanResult, error) {
	var (
		mu       sync.Mutex
		combined ScanResult
	)

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(m.concurrency)

	for _, auditor := range m.auditors {
		a := auditor
		g.Go(func() error {
			slog.Debug("Running auditor", "name", a.Name())

			findings, err := a.Audit(ctx, client, cfg)
			if err != nil {
				mu.Lock()
				combined.Errors = append(combined.Errors, fmt.Sprintf("%s: %v", a.Name(), err))
				mu.Unlock()
				slog.Warn("Auditor failed", "name", a.Name(), "error", err)
				return nil
			}

			mu.Lock()
			combined.Findings = append(combined.Findings, findings...)
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &combined, nil
}

// AllAuditors returns the full set of Redis auditors.
func AllAuditors() []Auditor {
	return []Auditor{
		&MemoryScanner{},
		&IdleKeyScanner{},
		&BigKeyScanner{},
		&ConnectionScanner{},
		&EvictionScanner{},
		&PersistenceScanner{},
		&SlowLogScanner{},
	}
}
