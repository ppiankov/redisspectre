package report

import (
	"encoding/json"
	"fmt"

	"github.com/ppiankov/redisspectre/internal/redis"
)

const sarifSchema = "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/main/sarif-2.1/schema/sarif-schema-2.1.0.json"

type sarifReport struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name    string      `json:"name"`
	Version string      `json:"version"`
	Rules   []sarifRule `json:"rules"`
}

type sarifRule struct {
	ID               string            `json:"id"`
	ShortDescription sarifMessage      `json:"shortDescription"`
	DefaultConfig    sarifDefaultLevel `json:"defaultConfiguration"`
}

type sarifDefaultLevel struct {
	Level string `json:"level"`
}

type sarifMessage struct {
	Text string `json:"text"`
}

type sarifResult struct {
	RuleID    string         `json:"ruleId"`
	Level     string         `json:"level"`
	Message   sarifMessage   `json:"message"`
	Locations []sarifLoc     `json:"locations,omitempty"`
	Props     map[string]any `json:"properties,omitempty"`
}

type sarifLoc struct {
	PhysicalLocation sarifPhysical `json:"physicalLocation"`
}

type sarifPhysical struct {
	ArtifactLocation sarifArtifact `json:"artifactLocation"`
}

type sarifArtifact struct {
	URI string `json:"uri"`
}

// Generate writes SARIF v2.1.0 output.
func (r *SARIFReporter) Generate(data Data) error {
	rules := buildSARIFRules()
	results := make([]sarifResult, 0, len(data.Findings))

	for _, f := range data.Findings {
		uri := fmt.Sprintf("redis://%s/%s/%s", data.Target.URIHash, f.ResourceType, f.ResourceID)
		results = append(results, sarifResult{
			RuleID:  string(f.ID),
			Level:   sarifLevel(f.Severity),
			Message: sarifMessage{Text: f.Message},
			Locations: []sarifLoc{
				{
					PhysicalLocation: sarifPhysical{
						ArtifactLocation: sarifArtifact{URI: uri},
					},
				},
			},
			Props: f.Metadata,
		})
	}

	rpt := sarifReport{
		Schema:  sarifSchema,
		Version: "2.1.0",
		Runs: []sarifRun{
			{
				Tool: sarifTool{
					Driver: sarifDriver{
						Name:    data.Tool,
						Version: data.Version,
						Rules:   rules,
					},
				},
				Results: results,
			},
		},
	}

	enc := json.NewEncoder(r.Writer)
	enc.SetIndent("", "  ")
	if err := enc.Encode(rpt); err != nil {
		return fmt.Errorf("encode SARIF report: %w", err)
	}
	return nil
}

func sarifLevel(s redis.Severity) string {
	switch s {
	case redis.SeverityCritical, redis.SeverityHigh:
		return "error"
	case redis.SeverityMedium:
		return "warning"
	default:
		return "note"
	}
}

func buildSARIFRules() []sarifRule {
	return []sarifRule{
		{ID: string(redis.FindingHighFragmentation), ShortDescription: sarifMessage{Text: "High memory fragmentation"}, DefaultConfig: sarifDefaultLevel{Level: "error"}},
		{ID: string(redis.FindingIdleKey), ShortDescription: sarifMessage{Text: "Idle key"}, DefaultConfig: sarifDefaultLevel{Level: "warning"}},
		{ID: string(redis.FindingBigKey), ShortDescription: sarifMessage{Text: "Oversized key"}, DefaultConfig: sarifDefaultLevel{Level: "warning"}},
		{ID: string(redis.FindingConnectionWaste), ShortDescription: sarifMessage{Text: "Connection waste"}, DefaultConfig: sarifDefaultLevel{Level: "note"}},
		{ID: string(redis.FindingEvictionRisk), ShortDescription: sarifMessage{Text: "Eviction risk"}, DefaultConfig: sarifDefaultLevel{Level: "error"}},
		{ID: string(redis.FindingNoPersistence), ShortDescription: sarifMessage{Text: "No persistence configured"}, DefaultConfig: sarifDefaultLevel{Level: "error"}},
		{ID: string(redis.FindingSlowCommand), ShortDescription: sarifMessage{Text: "Slow command"}, DefaultConfig: sarifDefaultLevel{Level: "warning"}},
	}
}
