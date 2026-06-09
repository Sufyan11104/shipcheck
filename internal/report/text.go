package report

import (
	"fmt"

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

	report := fmt.Sprintf(`ShipCheck Deployment Readiness Report
Path: %s
Git repository: %s
Files scanned: %d
Directories scanned: %d
Score: not calculated yet
Checks: coming in Stage 2
`, result.Path, gitRepo, result.FileCount, result.DirectoryCount)

	fmt.Print(report)
	return nil
}
