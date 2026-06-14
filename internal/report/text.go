package report

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Sufyan11104/shipcheck/internal/engine"
	"github.com/Sufyan11104/shipcheck/internal/rules"
	"github.com/Sufyan11104/shipcheck/internal/scanner"
)

// TextOptions controls human-readable text report rendering.
type TextOptions struct {
	Verbose bool
}

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
	return RenderTextWithOptions(w, auditReport, TextOptions{})
}

// RenderTextWithOptions writes a text-based deployment readiness report.
func RenderTextWithOptions(w io.Writer, auditReport AuditReport, options TextOptions) error {
	gitRepo := "no"
	if auditReport.GitRepository {
		gitRepo = "yes"
	}

	report := fmt.Sprintf(`ShipCheck Deployment Readiness Report
Path: %s
Git repository: %s
Files scanned: %d
Directories scanned: %d

Score: %d/100 — %s
Passed: %d
Warnings: %d
Failed: %d
Skipped: %d

`, auditReport.Path, gitRepo, auditReport.FilesScanned, auditReport.DirectoriesScanned, auditReport.Score, scoreLabel(auditReport.Score), auditReport.PassedCount, auditReport.WarningCount, auditReport.FailedCount, auditReport.SkippedCount)

	if _, err := fmt.Fprint(w, report); err != nil {
		return err
	}

	if options.Verbose {
		if _, err := fmt.Fprintln(w, "Full findings (--verbose):"); err != nil {
			return err
		}
		return renderGroupedFindings(w, auditReport.Findings)
	}

	if _, err := fmt.Fprintln(w, "Findings requiring attention:"); err != nil {
		return err
	}
	visible := visibleTextFindings(auditReport.Findings, false)
	if len(visible) == 0 {
		if _, err := fmt.Fprintln(w, "  No warnings or failures."); err != nil {
			return err
		}
	} else if err := renderGroupedFindings(w, visible); err != nil {
		return err
	}

	_, err := fmt.Fprintln(w, "\nPassed and skipped checks are hidden by default. Use --verbose to show all findings.")
	return err
}

func renderGroupedFindings(w io.Writer, findings []ReportFinding) error {
	grouped := groupFindingsByCategory(findings)
	for _, category := range categoryOrder() {
		items := grouped[category.key]
		if len(items) == 0 {
			continue
		}

		if _, err := fmt.Fprintf(w, "\n%s\n", category.title); err != nil {
			return err
		}
		for _, finding := range items {
			if err := renderTextFinding(w, finding); err != nil {
				return err
			}
		}
	}

	return nil
}

func renderTextFinding(w io.Writer, finding ReportFinding) error {
	symbol := getSymbol(finding.Status)
	if _, err := fmt.Fprintf(w, "  %s %s\n", symbol, finding.ID); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "    %s\n", finding.Message); err != nil {
		return err
	}
	if shouldShowRemediation(finding) {
		if _, err := fmt.Fprintf(w, "    Fix: %s\n", finding.Remediation); err != nil {
			return err
		}
	}
	return nil
}

func visibleTextFindings(findings []ReportFinding, verbose bool) []ReportFinding {
	if verbose {
		return findings
	}

	visible := make([]ReportFinding, 0, len(findings))
	for _, finding := range findings {
		if finding.Status == rules.StatusWarn || finding.Status == rules.StatusFail {
			visible = append(visible, finding)
		}
	}
	return visible
}

func groupFindingsByCategory(findings []ReportFinding) map[string][]ReportFinding {
	grouped := make(map[string][]ReportFinding)
	for _, finding := range findings {
		key := strings.ToLower(finding.Category)
		grouped[key] = append(grouped[key], finding)
	}
	return grouped
}

type categoryHeading struct {
	key   string
	title string
}

func categoryOrder() []categoryHeading {
	return []categoryHeading{
		{key: "repo", title: "Repository"},
		{key: "env", title: "Environment"},
		{key: "docker", title: "Docker"},
		{key: "ci", title: "GitHub Actions"},
		{key: "k8s", title: "Kubernetes"},
		{key: "terraform", title: "Terraform"},
		{key: "docs", title: "Documentation"},
	}
}

func shouldShowRemediation(finding ReportFinding) bool {
	if finding.Status != rules.StatusWarn && finding.Status != rules.StatusFail {
		return false
	}
	return finding.Remediation != "" && finding.Remediation != "N/A"
}

func scoreLabel(score int) string {
	switch {
	case score >= 90:
		return "Excellent"
	case score >= 70:
		return "Good"
	case score >= 50:
		return "Needs attention"
	default:
		return "High risk"
	}
}

func getSymbol(status rules.Status) string {
	switch status {
	case rules.StatusPass:
		return "✓"
	case rules.StatusWarn:
		return "!"
	case rules.StatusFail:
		return "✗"
	case rules.StatusSkip:
		return "-"
	default:
		return "?"
	}
}
