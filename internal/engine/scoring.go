package engine

import (
	"github.com/Sufyan11104/shipcheck/internal/rules"
)

// CalculateScore converts findings into a 0-100 readiness score
func CalculateScore(findings []rules.Finding) int {
	score := 100

	for _, finding := range findings {
		switch finding.Status {
		case rules.StatusFail:
			// Failed findings subtract more
			switch finding.Severity {
			case rules.SeverityHigh:
				score -= 25
			case rules.SeverityMedium:
				score -= 15
			case rules.SeverityLow:
				score -= 5
			}
		case rules.StatusWarn:
			// Warnings subtract less
			switch finding.Severity {
			case rules.SeverityHigh:
				score -= 15
			case rules.SeverityMedium:
				score -= 8
			case rules.SeverityLow:
				score -= 3
			}
		case rules.StatusPass:
			// Passing checks do not reduce score
		}
	}

	// Clamp score between 0 and 100
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// SummarizeFinding counts findings by status
func SummarizeFindings(findings []rules.Finding) (passed, warned, failed int) {
	for _, f := range findings {
		switch f.Status {
		case rules.StatusPass:
			passed++
		case rules.StatusWarn:
			warned++
		case rules.StatusFail:
			failed++
		}
	}
	return
}
