package analyzer

import (
	"github.com/ppiankov/redisspectre/internal/redis"
)

// Analyze filters findings by minimum severity and computes summary statistics.
func Analyze(result *redis.ScanResult, cfg AnalyzerConfig) *AnalysisResult {
	var filtered []redis.Finding
	for _, f := range result.Findings {
		if redis.MeetsSeverityMin(f.Severity, cfg.SeverityMin) {
			filtered = append(filtered, f)
		}
	}

	summary := Summary{
		TotalResourcesScanned: result.ResourcesScanned,
		TotalFindings:         len(filtered),
		BySeverity:            make(map[string]int),
		ByResourceType:        make(map[string]int),
		ByFindingID:           make(map[string]int),
	}

	for _, f := range filtered {
		summary.BySeverity[string(f.Severity)]++
		summary.ByResourceType[f.ResourceType]++
		summary.ByFindingID[string(f.ID)]++
	}

	return &AnalysisResult{
		Findings: filtered,
		Summary:  summary,
		Errors:   result.Errors,
	}
}
