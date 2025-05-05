package interfaces

// RuleInput represents the data passed into each rule.
// This can be extended to support different agent types (MRs, alerts, incidents, etc.)
type RuleInput struct {
	Source    string            // e.g., "gitlab", "alertmanager", "pagerduty"
	Subject   string            // Title, summary, or key identifier
	Content   string            // Description, message body, etc.
	Tags      []string          // Labels, environments, categories
	Files     []string          // List of touched files (MRs, deploys)
	Diff      string            // Raw code or config diff, if available
	Metadata  map[string]string // Agent-specific details (e.g., "author", "alertSeverity")
	Raw       any                   // Optional strongly typed struct, like MergeRequestMetadata
}

// RuleResult represents the outcome of a rule evaluation.
type RuleResult struct {
	RuleID     string
	Applicable bool
	Passed     bool
	Message    string
	Severity   SeverityLevel
}

// SeverityLevel defines rule result severity for decisioning.
type SeverityLevel string

const (
	SeverityInfo     SeverityLevel = "info"
	SeverityLow      SeverityLevel = "low"
	SeverityMedium   SeverityLevel = "medium"
	SeverityHigh     SeverityLevel = "high"
	SeverityCritical SeverityLevel = "critical"
)