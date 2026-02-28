package analyzer

import (
	"testing"

	"github.com/ppiankov/redisspectre/internal/redis"
)

func TestAnalyze_FiltersBySeverity(t *testing.T) {
	result := &redis.ScanResult{
		Findings: []redis.Finding{
			{ID: redis.FindingHighFragmentation, Severity: redis.SeverityHigh, ResourceType: "Redis"},
			{ID: redis.FindingIdleKey, Severity: redis.SeverityMedium, ResourceType: "Key"},
			{ID: redis.FindingConnectionWaste, Severity: redis.SeverityLow, ResourceType: "Redis"},
			{ID: redis.FindingEvictionRisk, Severity: redis.SeverityCritical, ResourceType: "Config"},
		},
		ResourcesScanned: 100,
	}

	analysis := Analyze(result, AnalyzerConfig{SeverityMin: redis.SeverityHigh})

	if len(analysis.Findings) != 2 {
		t.Errorf("expected 2 findings (high+critical), got %d", len(analysis.Findings))
	}
	if analysis.Summary.TotalFindings != 2 {
		t.Errorf("expected total findings 2, got %d", analysis.Summary.TotalFindings)
	}
	if analysis.Summary.TotalResourcesScanned != 100 {
		t.Errorf("expected resources scanned 100, got %d", analysis.Summary.TotalResourcesScanned)
	}
}

func TestAnalyze_AllPass(t *testing.T) {
	result := &redis.ScanResult{
		Findings: []redis.Finding{
			{ID: redis.FindingHighFragmentation, Severity: redis.SeverityHigh, ResourceType: "Redis"},
			{ID: redis.FindingIdleKey, Severity: redis.SeverityMedium, ResourceType: "Key"},
		},
	}

	analysis := Analyze(result, AnalyzerConfig{SeverityMin: redis.SeverityLow})

	if len(analysis.Findings) != 2 {
		t.Errorf("expected 2 findings with low min, got %d", len(analysis.Findings))
	}
}

func TestAnalyze_SummaryBreakdown(t *testing.T) {
	result := &redis.ScanResult{
		Findings: []redis.Finding{
			{ID: redis.FindingIdleKey, Severity: redis.SeverityMedium, ResourceType: "Key"},
			{ID: redis.FindingIdleKey, Severity: redis.SeverityMedium, ResourceType: "Key"},
			{ID: redis.FindingBigKey, Severity: redis.SeverityMedium, ResourceType: "Key"},
		},
	}

	analysis := Analyze(result, AnalyzerConfig{})

	if analysis.Summary.BySeverity["medium"] != 3 {
		t.Errorf("expected 3 medium findings, got %d", analysis.Summary.BySeverity["medium"])
	}
	if analysis.Summary.ByResourceType["Key"] != 3 {
		t.Errorf("expected 3 Key findings, got %d", analysis.Summary.ByResourceType["Key"])
	}
	if analysis.Summary.ByFindingID["IDLE_KEY"] != 2 {
		t.Errorf("expected 2 IDLE_KEY findings, got %d", analysis.Summary.ByFindingID["IDLE_KEY"])
	}
	if analysis.Summary.ByFindingID["BIG_KEY"] != 1 {
		t.Errorf("expected 1 BIG_KEY finding, got %d", analysis.Summary.ByFindingID["BIG_KEY"])
	}
}

func TestAnalyze_Empty(t *testing.T) {
	result := &redis.ScanResult{}
	analysis := Analyze(result, AnalyzerConfig{})

	if len(analysis.Findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(analysis.Findings))
	}
	if analysis.Summary.TotalFindings != 0 {
		t.Errorf("expected total 0, got %d", analysis.Summary.TotalFindings)
	}
}

func TestAnalyze_PreservesErrors(t *testing.T) {
	result := &redis.ScanResult{
		Errors: []string{"memory: connection refused"},
	}
	analysis := Analyze(result, AnalyzerConfig{})

	if len(analysis.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(analysis.Errors))
	}
}
