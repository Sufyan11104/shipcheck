package report

import (
	"github.com/Sufyan11104/shipcheck/internal/engine"
	"github.com/Sufyan11104/shipcheck/internal/rules"
	"github.com/Sufyan11104/shipcheck/internal/scanner"
)

const (
	FormatText     = "text"
	FormatJSON     = "json"
	FormatMarkdown = "markdown"
	FormatSARIF    = "sarif"
)

// AuditReport is the stable report model used by all output formats.
type AuditReport struct {
	Path               string          `json:"path"`
	GitRepository      bool            `json:"gitRepository"`
	FilesScanned       int64           `json:"filesScanned"`
	DirectoriesScanned int64           `json:"directoriesScanned"`
	Score              int             `json:"score"`
	PassedCount        int             `json:"passedCount"`
	WarningCount       int             `json:"warningCount"`
	FailedCount        int             `json:"failedCount"`
	SkippedCount       int             `json:"skippedCount"`
	Findings           []ReportFinding `json:"findings"`
}

// ReportFinding is the JSON/Markdown representation of a rules finding.
type ReportFinding struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Category    string         `json:"category"`
	Severity    rules.Severity `json:"severity"`
	Status      rules.Status   `json:"status"`
	Message     string         `json:"message"`
	Remediation string         `json:"remediation"`
	Path        string         `json:"path,omitempty"`
	Evidence    string         `json:"evidence,omitempty"`
}

// NewAuditReport builds a report from scanner and engine results.
func NewAuditReport(result *scanner.ScanResult, findings []rules.Finding, score int) AuditReport {
	passed, warned, failed, skipped := engine.SummarizeFindingsWithSkipped(findings)

	return AuditReport{
		Path:               result.Path,
		GitRepository:      result.IsGitRepository,
		FilesScanned:       result.FileCount,
		DirectoriesScanned: result.DirectoryCount,
		Score:              score,
		PassedCount:        passed,
		WarningCount:       warned,
		FailedCount:        failed,
		SkippedCount:       skipped,
		Findings:           toReportFindings(findings),
	}
}

func toReportFindings(findings []rules.Finding) []ReportFinding {
	reportFindings := make([]ReportFinding, 0, len(findings))
	for _, finding := range findings {
		reportFindings = append(reportFindings, ReportFinding{
			ID:          finding.ID,
			Title:       finding.Title,
			Category:    finding.Category,
			Severity:    finding.Severity,
			Status:      finding.Status,
			Message:     finding.Message,
			Remediation: finding.Remediation,
			Path:        finding.Path,
			Evidence:    finding.Evidence,
		})
	}

	return reportFindings
}
