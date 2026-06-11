package rules

// Severity levels for findings
type Severity string

const (
	SeverityInfo   Severity = "info"
	SeverityLow    Severity = "low"
	SeverityMedium Severity = "medium"
	SeverityHigh   Severity = "high"
)

// Status of a finding
type Status string

const (
	StatusPass Status = "pass"
	StatusWarn Status = "warn"
	StatusFail Status = "fail"
	StatusSkip Status = "skip"
)

// Finding represents a single audit finding
type Finding struct {
	ID          string
	Title       string
	Category    string
	Severity    Severity
	Status      Status
	Message     string
	Remediation string
	Path        string // optional
}

// CheckResult groups findings by status
type CheckResult struct {
	Findings []Finding
	Passed   int
	Warned   int
	Failed   int
	Skipped  int
	Score    int
}
