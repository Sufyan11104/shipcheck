package engine

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Sufyan11104/shipcheck/internal/rules"
)

// Engine runs security and hygiene checks
type Engine struct {
	path string
}

// NewEngine creates a new audit engine
func NewEngine(path string) *Engine {
	return &Engine{path: path}
}

// RunChecks executes all checks and returns findings
func (e *Engine) RunChecks(isGitRepo bool) ([]rules.Finding, int) {
	return e.RunChecksWithCategories(isGitRepo, nil)
}

// RunChecksWithCategories executes checks with awareness of explicitly requested categories.
func (e *Engine) RunChecksWithCategories(isGitRepo bool, categories []string) ([]rules.Finding, int) {
	var findings []rules.Finding

	// Run generic checks
	findings = append(findings, rules.CheckReadmeExists(e.path))
	findings = append(findings, rules.CheckGitignoreExists(e.path))
	findings = append(findings, rules.CheckEnvNotCommitted(e.path, isGitRepo))
	findings = append(findings, rules.CheckEnvExampleExists(e.path))

	// Run Docker checks
	findings = append(findings, rules.CheckDockerfileExists(e.path))
	findings = append(findings, rules.CheckDockerignoreExists(e.path))
	findings = append(findings, rules.CheckDockerfileNonRootUser(e.path))
	findings = append(findings, rules.CheckDockerfileHealthcheck(e.path))
	findings = append(findings, rules.CheckDockerfileNoEnvCopy(e.path))
	findings = append(findings, rules.CheckDockerfileNoSecretEnv(e.path))

	// Run GitHub Actions checks
	findings = append(findings, rules.CheckWorkflowsDirExists(e.path))
	findings = append(findings, rules.CheckWorkflowFileExists(e.path))
	findings = append(findings, rules.CheckTestStepExists(e.path))
	findings = append(findings, rules.CheckBuildStepExists(e.path))
	findings = append(findings, rules.CheckDeployAfterTests(e.path))
	findings = append(findings, rules.CheckActionsPinned(e.path))
	findings = append(findings, rules.CheckNoSecretEcho(e.path))
	findings = append(findings, rules.CheckPermissionsDeclared(e.path))

	// Run Kubernetes checks
	findings = append(findings, rules.CheckK8sManifestExists(e.path))
	findings = append(findings, rules.CheckK8sWorkloadExists(e.path))
	findings = append(findings, rules.CheckK8sReadinessProbeExists(e.path))
	findings = append(findings, rules.CheckK8sLivenessProbeExists(e.path))
	findings = append(findings, rules.CheckK8sResourceRequests(e.path))
	findings = append(findings, rules.CheckK8sResourceLimits(e.path))
	findings = append(findings, rules.CheckK8sNoLatestImageTag(e.path))
	findings = append(findings, rules.CheckK8sReplicasConfigured(e.path))

	// Run Terraform/IaC checks
	findings = append(findings, rules.CheckTerraformFilesExist(e.path))
	findings = append(findings, rules.CheckTerraformFmtRecommended(e.path))
	findings = append(findings, rules.CheckTerraformValidateRecommended(e.path))
	findings = append(findings, rules.CheckTerraformRequiredProvidersExists(e.path))
	findings = append(findings, rules.CheckTerraformProviderVersionsConstrained(e.path))
	findings = append(findings, rules.CheckTerraformBackendConfigured(e.path))
	findings = append(findings, rules.CheckTerraformNoSuspiciousVariableDefaults(e.path))
	findings = append(findings, rules.CheckTerraformLockfilePresent(e.path))

	findings = e.applyOptionalCategoryContext(findings, categories)

	// Calculate score
	score := CalculateScore(findings)

	return findings, score
}

func (e *Engine) applyOptionalCategoryContext(findings []rules.Finding, categories []string) []rules.Finding {
	explicit := explicitCategorySet(categories)
	active := e.detectActiveOptionalCategories()

	for i := range findings {
		category := strings.ToLower(findings[i].Category)
		if !optionalCategory(category) || explicit[category] || active[category] {
			continue
		}

		findings[i].Status = rules.StatusSkip
		findings[i].Severity = rules.SeverityLow
		findings[i].Remediation = "N/A"
	}

	return findings
}

func explicitCategorySet(categories []string) map[string]bool {
	explicit := make(map[string]bool, len(categories))
	for _, category := range categories {
		explicit[strings.ToLower(category)] = true
	}
	return explicit
}

func optionalCategory(category string) bool {
	switch category {
	case "docker", "ci", "k8s", "terraform":
		return true
	default:
		return false
	}
}

func (e *Engine) detectActiveOptionalCategories() map[string]bool {
	return map[string]bool{
		"docker":    fileExists(filepath.Join(e.path, "Dockerfile")) || fileExists(filepath.Join(e.path, ".dockerignore")),
		"ci":        dirExists(filepath.Join(e.path, ".github", "workflows")),
		"k8s":       rules.HasK8sManifest(e.path),
		"terraform": rules.HasTerraformFiles(e.path),
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
