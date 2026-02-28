package commands

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"

	"github.com/ppiankov/redisspectre/internal/report"
)

// enhanceError wraps an error with context and suggestions for common Redis issues.
func enhanceError(action string, err error) error {
	msg := err.Error()

	var hint string
	switch {
	case strings.Contains(msg, "connection refused"):
		hint = "Cannot connect to Redis. Verify the server is running and --addr is correct"
	case strings.Contains(msg, "NOAUTH") || strings.Contains(msg, "ERR AUTH"):
		hint = "Authentication failed. Check --password or REDIS_PASSWORD environment variable"
	case strings.Contains(msg, "NOPERM") || strings.Contains(msg, "no permissions"):
		hint = "Insufficient permissions. redisspectre needs INFO, SCAN, OBJECT, MEMORY, SLOWLOG, and CONFIG GET access"
	case strings.Contains(msg, "timeout") || strings.Contains(msg, "deadline exceeded"):
		hint = "Operation timed out. Try increasing --timeout"
	case strings.Contains(msg, "EOF") || strings.Contains(msg, "broken pipe"):
		hint = "Connection lost. Verify Redis server stability and network connectivity"
	}

	if hint != "" {
		return fmt.Errorf("%s: %w\n  hint: %s", action, err, hint)
	}
	return fmt.Errorf("%s: %w", action, err)
}

// computeTargetHash generates a SHA256 hash for the target URI.
func computeTargetHash(addr string, db int) string {
	input := fmt.Sprintf("redis://%s/%d", addr, db)
	h := sha256.Sum256([]byte(input))
	return fmt.Sprintf("sha256:%x", h)
}

// selectReporter creates the appropriate reporter for the given format.
func selectReporter(format, outputFile string) (report.Reporter, error) {
	w := os.Stdout
	if outputFile != "" {
		f, err := os.Create(outputFile)
		if err != nil {
			return nil, fmt.Errorf("create output file: %w", err)
		}
		w = f
	}

	switch format {
	case "json":
		return &report.JSONReporter{Writer: w}, nil
	case "text":
		return &report.TextReporter{Writer: w}, nil
	case "sarif":
		return &report.SARIFReporter{Writer: w}, nil
	case "spectrehub":
		return &report.SpectreHubReporter{Writer: w}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s (use text, json, sarif, or spectrehub)", format)
	}
}
