package report

import (
	"fmt"
	"io"
	"os"

	"github.com/Sufyan11104/shipcheck/internal/engine"
	"github.com/Sufyan11104/shipcheck/internal/rules"
	"github.com/Sufyan11104/shipcheck/internal/scanner"
)

// PrintTextReport prints a text-based deployment readiness report
func PrintTextReport(result *scanner.ScanResult) error {
	if result.Error != nil {
		return result.Error
	}

	eng := engine.NewEngine(result.Path)
	findings, score := eng.RunChecks(result.IsGitRepository)
	auditReport := NewAuditReport(result, findings, score)

	return RenderText(os.Stdout, auditReport)
}

// RenderText writes a text-based deployment readiness report.
func RenderText(w io.Writer, auditReport AuditReport) error {
	gitRepo := "no"
	if auditReport.GitRepository {
		gitRepo = "yes"
	}

	report := fmt.Sprintf(`ShipCheck Deployment Readiness Report
Path: %s
Git repository: %s
Files scanned: %d
Directories scanned: %d

Score: %d/100
Passed: %d
Warnings: %d
Failed: %d

Findings:
`, auditReport.Path, gitRepo, auditReport.FilesScanned, auditReport.DirectoriesScanned, auditReport.Score, auditReport.PassedCount, auditReport.WarningCount, auditReport.FailedCount)

	if _, err := fmt.Fprint(w, report); err != nil {
		return err
	}

	// Print each finding
	for _, finding := range auditReport.Findings {
		symbol := getSymbol(finding.Status)
		if _, err := fmt.Fprintf(w, "%s %s - %s\n", symbol, finding.ID, finding.Message); err != nil {
			return err
		}
	}

	return nil
}

func getSymbol(status rules.Status) string {
	switch status {
	case rules.StatusPass:
		return "✓"
	case rules.StatusWarn:
		return "!"
	case rules.StatusFail:
		return "✗"
	default:
		return "?"
	}
}
