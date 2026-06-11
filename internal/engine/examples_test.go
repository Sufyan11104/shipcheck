package engine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Sufyan11104/shipcheck/internal/rules"
)

func TestExampleServicesAuditBehavior(t *testing.T) {
	goodFindings, goodScore := auditExampleService(t, "good-service")
	riskyFindings, riskyScore := auditExampleService(t, "risky-service")

	if goodScore <= riskyScore {
		t.Fatalf("expected good-service score %d to be higher than risky-service score %d", goodScore, riskyScore)
	}

	goodWarnings := countFindingsByStatus(goodFindings, rules.StatusWarn)
	riskyWarnings := countFindingsByStatus(riskyFindings, rules.StatusWarn)
	if goodWarnings >= riskyWarnings {
		t.Fatalf("expected good-service to have fewer warnings than risky-service; got %d and %d", goodWarnings, riskyWarnings)
	}

	for _, category := range []string{"docker", "ci", "k8s", "terraform", "env"} {
		if !hasWarningInCategory(riskyFindings, category) {
			t.Errorf("expected risky-service to trigger a %s warning", category)
		}
	}
}

func auditExampleService(t *testing.T, name string) ([]rules.Finding, int) {
	t.Helper()

	path := filepath.Join(repoRoot(t), "examples", name)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected example service %q to exist: %v", name, err)
	}

	return NewEngine(path).RunChecks(true)
}

func repoRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("failed to find repository root")
		}
		dir = parent
	}
}

func countFindingsByStatus(findings []rules.Finding, status rules.Status) int {
	count := 0
	for _, finding := range findings {
		if finding.Status == status {
			count++
		}
	}
	return count
}

func hasWarningInCategory(findings []rules.Finding, category string) bool {
	for _, finding := range findings {
		if finding.Category == category && finding.Status == rules.StatusWarn {
			return true
		}
	}
	return false
}
