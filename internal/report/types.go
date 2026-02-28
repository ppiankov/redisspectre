package report

import (
	"io"
	"time"

	"github.com/ppiankov/redisspectre/internal/analyzer"
	"github.com/ppiankov/redisspectre/internal/redis"
)

// Reporter is the interface for output formatters.
type Reporter interface {
	Generate(data Data) error
}

// Data holds all information needed to generate a report.
type Data struct {
	Tool      string           `json:"tool"`
	Version   string           `json:"version"`
	Timestamp time.Time        `json:"timestamp"`
	Target    Target           `json:"target"`
	Config    ReportConfig     `json:"config"`
	Findings  []redis.Finding  `json:"findings"`
	Summary   analyzer.Summary `json:"summary"`
	Errors    []string         `json:"errors,omitempty"`
}

// Target identifies what was audited.
type Target struct {
	Type    string `json:"type"`
	URIHash string `json:"uri_hash"`
}

// ReportConfig captures the audit configuration used.
type ReportConfig struct {
	Addr       string `json:"addr"`
	DB         int    `json:"db"`
	SampleSize int    `json:"sample_size"`
	IdleDays   int    `json:"idle_days"`
}

// TextReporter generates human-readable terminal output.
type TextReporter struct {
	Writer io.Writer
}

// JSONReporter generates spectre/v1 envelope JSON output.
type JSONReporter struct {
	Writer io.Writer
}

// SpectreHubReporter generates SpectreHub envelope JSON output.
type SpectreHubReporter struct {
	Writer io.Writer
}

// SARIFReporter generates SARIF v2.1.0 output.
type SARIFReporter struct {
	Writer io.Writer
}
