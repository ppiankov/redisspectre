package commands

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ppiankov/redisspectre/internal/analyzer"
	"github.com/ppiankov/redisspectre/internal/redis"
	"github.com/ppiankov/redisspectre/internal/report"
	"github.com/spf13/cobra"
)

var auditFlags struct {
	format     string
	outputFile string
	sampleSize int
	idleDays   int
	bigKeySize int64
	timeout    time.Duration
}

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Run full Redis audit",
	Long: `Audit a Redis instance for waste and hygiene issues: memory fragmentation,
idle keys, big keys, connection waste, eviction policy, persistence
configuration, and slow commands.

Requires connectivity to the Redis instance.`,
	RunE: runAudit,
}

func init() {
	auditCmd.Flags().StringVar(&auditFlags.format, "format", "text", "Output format: text, json, sarif, spectrehub")
	auditCmd.Flags().StringVarP(&auditFlags.outputFile, "output", "o", "", "Output file path (default: stdout)")
	auditCmd.Flags().IntVar(&auditFlags.sampleSize, "sample-size", 10000, "Number of keys to sample")
	auditCmd.Flags().IntVar(&auditFlags.idleDays, "idle-days", 30, "Key inactivity threshold (days)")
	auditCmd.Flags().Int64Var(&auditFlags.bigKeySize, "big-key-size", 10*1024*1024, "Big key threshold (bytes)")
	auditCmd.Flags().DurationVar(&auditFlags.timeout, "timeout", 5*time.Minute, "Audit timeout")

	rootCmd.AddCommand(auditCmd)
}

func runAudit(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	if auditFlags.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, auditFlags.timeout)
		defer cancel()
	}

	applyConfigDefaults()

	resolvedAddr := resolveAddr()
	resolvedPassword := resolvePassword()

	client, err := redis.NewClient(resolvedAddr, resolvedPassword, db)
	if err != nil {
		return enhanceError("create redis client", err)
	}
	defer func() { _ = client.Close() }()

	if err := client.Ping(ctx); err != nil {
		return enhanceError("connect to redis", err)
	}

	auditCfg := redis.AuditConfig{
		Addr:       resolvedAddr,
		DB:         db,
		SampleSize: auditFlags.sampleSize,
		IdleDays:   auditFlags.idleDays,
		BigKeySize: auditFlags.bigKeySize,
	}

	slog.Info("Starting audit", "addr", resolvedAddr, "db", db, "sample-size", auditFlags.sampleSize)

	multi := redis.NewMultiAuditor(redis.AllAuditors(), 4)
	result, err := multi.AuditAll(ctx, client, auditCfg)
	if err != nil {
		return enhanceError("audit redis", err)
	}

	analysis := analyzer.Analyze(result, analyzer.AnalyzerConfig{})

	data := report.Data{
		Tool:      "redisspectre",
		Version:   version,
		Timestamp: time.Now().UTC(),
		Target: report.Target{
			Type:    "redis",
			URIHash: computeTargetHash(resolvedAddr, db),
		},
		Config: report.ReportConfig{
			Addr:       resolvedAddr,
			DB:         db,
			SampleSize: auditFlags.sampleSize,
			IdleDays:   auditFlags.idleDays,
		},
		Findings: analysis.Findings,
		Summary:  analysis.Summary,
		Errors:   analysis.Errors,
	}

	reporter, err := selectReporter(auditFlags.format, auditFlags.outputFile)
	if err != nil {
		return err
	}

	if err := reporter.Generate(data); err != nil {
		return fmt.Errorf("generate report: %w", err)
	}

	if analysis.Summary.TotalFindings > 0 {
		slog.Info("Audit complete", "findings", analysis.Summary.TotalFindings)
	}

	return nil
}

func applyConfigDefaults() {
	if auditFlags.format == "text" && cfg.Format != "" {
		auditFlags.format = cfg.Format
	}
	if auditFlags.sampleSize == 10000 && cfg.SampleSize > 0 {
		auditFlags.sampleSize = cfg.SampleSize
	}
	if auditFlags.idleDays == 30 && cfg.IdleDays > 0 {
		auditFlags.idleDays = cfg.IdleDays
	}
	if auditFlags.bigKeySize == 10*1024*1024 && cfg.BigKeySize > 0 {
		auditFlags.bigKeySize = cfg.BigKeySize
	}
}
