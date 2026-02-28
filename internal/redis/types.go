package redis

// Severity levels for findings.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
)

// SeverityRank returns a numeric rank for sorting (higher = more severe).
func SeverityRank(s Severity) int {
	switch s {
	case SeverityCritical:
		return 4
	case SeverityHigh:
		return 3
	case SeverityMedium:
		return 2
	case SeverityLow:
		return 1
	default:
		return 0
	}
}

// MeetsSeverityMin returns true if s meets or exceeds the minimum severity.
func MeetsSeverityMin(s, min Severity) bool {
	return SeverityRank(s) >= SeverityRank(min)
}

// ParseSeverity converts a string to Severity, defaulting to low.
func ParseSeverity(s string) Severity {
	switch s {
	case "critical":
		return SeverityCritical
	case "high":
		return SeverityHigh
	case "medium":
		return SeverityMedium
	case "low":
		return SeverityLow
	default:
		return SeverityLow
	}
}

// FindingID identifies the type of issue detected.
type FindingID string

const (
	FindingHighFragmentation FindingID = "HIGH_FRAGMENTATION"
	FindingIdleKey           FindingID = "IDLE_KEY"
	FindingBigKey            FindingID = "BIG_KEY"
	FindingConnectionWaste   FindingID = "CONNECTION_WASTE"
	FindingEvictionRisk      FindingID = "EVICTION_RISK"
	FindingNoPersistence     FindingID = "NO_PERSISTENCE"
	FindingSlowCommand       FindingID = "SLOW_COMMAND"
)

// Finding represents a single audit issue.
type Finding struct {
	ID           FindingID      `json:"id"`
	Severity     Severity       `json:"severity"`
	ResourceType string         `json:"resource_type"`
	ResourceID   string         `json:"resource_id"`
	Message      string         `json:"message"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

// ScanResult holds all findings from scanning a Redis instance.
type ScanResult struct {
	Findings         []Finding `json:"findings"`
	Errors           []string  `json:"errors,omitempty"`
	ResourcesScanned int       `json:"resources_scanned"`
}

// AuditConfig holds parameters that control auditing behavior.
type AuditConfig struct {
	Addr       string
	DB         int
	SampleSize int
	IdleDays   int
	BigKeySize int64
}
