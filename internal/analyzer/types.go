package analyzer

import (
	"github.com/ppiankov/redisspectre/internal/redis"
)

// Summary holds aggregated statistics about audit findings.
type Summary struct {
	TotalResourcesScanned int            `json:"total_resources_scanned"`
	TotalFindings         int            `json:"total_findings"`
	BySeverity            map[string]int `json:"by_severity"`
	ByResourceType        map[string]int `json:"by_resource_type"`
	ByFindingID           map[string]int `json:"by_finding_id"`
}

// AnalysisResult holds filtered findings and computed summary.
type AnalysisResult struct {
	Findings []redis.Finding `json:"findings"`
	Summary  Summary         `json:"summary"`
	Errors   []string        `json:"errors,omitempty"`
}

// AnalyzerConfig controls analysis behavior.
type AnalyzerConfig struct {
	SeverityMin redis.Severity
}
