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

	// Run all checks
	findings = append(findings, rules.CheckReadmeExists(e.path))
	findings = append(findings, rules.CheckGitignoreExists(e.path))
	findings = append(findings, rules.CheckEnvNotCommitted(e.path, isGitRepo))
	findings = append(findings, rules.CheckEnvExampleExists(e.path))

	// Calculate score
	score := CalculateScore(findings)

	return findings, score
}
