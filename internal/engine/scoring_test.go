package engine

import (
	"testing"

	"github.com/Sufyan11104/shipcheck/internal/rules"
)

func TestCalculateScore_AllPass(t *testing.T) {
	findings := []rules.Finding{
		{Status: rules.StatusPass, Severity: rules.SeverityHigh},
		{Status: rules.StatusPass, Severity: rules.SeverityHigh},
	}

	score := CalculateScore(findings)
	if score != 100 {
		t.Errorf("expected 100 for all passing, got %d", score)
	}
}

func TestCalculateScore_HighFail(t *testing.T) {
	findings := []rules.Finding{
		{Status: rules.StatusFail, Severity: rules.SeverityHigh},
	}

	score := CalculateScore(findings)
	if score != 75 {
		t.Errorf("expected 75 for one high fail, got %d", score)
	}
}

func TestCalculateScore_MixedFindings(t *testing.T) {
	findings := []rules.Finding{
		{Status: rules.StatusPass, Severity: rules.SeverityHigh},
		{Status: rules.StatusSkip, Severity: rules.SeverityHigh},
		{Status: rules.StatusWarn, Severity: rules.SeverityHigh},
		{Status: rules.StatusFail, Severity: rules.SeverityMedium},
	}

	score := CalculateScore(findings)
	// 100 - 15 (warn) - 15 (fail medium) = 70
	expected := 70
	if score != expected {
		t.Errorf("expected %d, got %d", expected, score)
	}
}

func TestCalculateScore_SkipDoesNotReduceScore(t *testing.T) {
	findings := []rules.Finding{
		{Status: rules.StatusSkip, Severity: rules.SeverityHigh},
		{Status: rules.StatusSkip, Severity: rules.SeverityMedium},
	}

	score := CalculateScore(findings)
	if score != 100 {
		t.Errorf("expected 100 for skipped findings, got %d", score)
	}
}

func TestCalculateScore_ClampsToZero(t *testing.T) {
	findings := []rules.Finding{
		{Status: rules.StatusFail, Severity: rules.SeverityHigh},
		{Status: rules.StatusFail, Severity: rules.SeverityHigh},
		{Status: rules.StatusFail, Severity: rules.SeverityHigh},
		{Status: rules.StatusFail, Severity: rules.SeverityHigh},
	}

	score := CalculateScore(findings)
	if score != 0 {
		t.Errorf("expected 0 (clamped), got %d", score)
	}
}

func TestSummarizeFindings(t *testing.T) {
	findings := []rules.Finding{
		{Status: rules.StatusPass},
		{Status: rules.StatusPass},
		{Status: rules.StatusWarn},
		{Status: rules.StatusFail},
		{Status: rules.StatusFail},
		{Status: rules.StatusFail},
		{Status: rules.StatusSkip},
		{Status: rules.StatusSkip},
	}

	passed, warned, failed := SummarizeFindings(findings)

	if passed != 2 {
		t.Errorf("expected 2 passed, got %d", passed)
	}
	if warned != 1 {
		t.Errorf("expected 1 warned, got %d", warned)
	}
	if failed != 3 {
		t.Errorf("expected 3 failed, got %d", failed)
	}
}

func TestSummarizeFindingsWithSkipped(t *testing.T) {
	findings := []rules.Finding{
		{Status: rules.StatusPass},
		{Status: rules.StatusWarn},
		{Status: rules.StatusFail},
		{Status: rules.StatusSkip},
		{Status: rules.StatusSkip},
	}

	passed, warned, failed, skipped := SummarizeFindingsWithSkipped(findings)

	if passed != 1 {
		t.Errorf("expected 1 passed, got %d", passed)
	}
	if warned != 1 {
		t.Errorf("expected 1 warned, got %d", warned)
	}
	if failed != 1 {
		t.Errorf("expected 1 failed, got %d", failed)
	}
	if skipped != 2 {
		t.Errorf("expected 2 skipped, got %d", skipped)
	}
}
