package engine

import (
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

	// Calculate score
	score := CalculateScore(findings)

	return findings, score
}
