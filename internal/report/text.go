package report

import (
	"fmt"

	"github.com/Sufyan11104/shipcheck/internal/engine"
	"github.com/Sufyan11104/shipcheck/internal/rules"
	"github.com/Sufyan11104/shipcheck/internal/scanner"
)

// PrintTextReport prints a text-based deployment readiness report
func PrintTextReport(result *scanner.ScanResult) error {
	if result.Error != nil {
		return result.Error
	}

	gitRepo := "no"
	if result.IsGitRepository {
		gitRepo = "yes"
	}

	// Run the audit engine
	eng := engine.NewEngine(result.Path)
	findings, score := eng.RunChecks(result.IsGitRepository)

	// Summarize findings
	passed, warned, failed := engine.SummarizeFindings(findings)

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
`, result.Path, gitRepo, result.FileCount, result.DirectoryCount, score, passed, warned, failed)

	fmt.Print(report)

	// Print each finding
	for _, finding := range findings {
		symbol := getSymbol(finding.Status)
		fmt.Printf("%s %s - %s\n", symbol, finding.ID, finding.Message)
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
