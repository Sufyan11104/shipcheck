package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Sufyan11104/shipcheck/internal/rules"
)

func TestRenderText_DefaultFormatWorks(t *testing.T) {
	var buf bytes.Buffer

	if err := Render(&buf, sampleTextAuditReport(), FormatText); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "ShipCheck Deployment Readiness Report") {
		t.Errorf("expected text report heading, got %s", output)
	}
	if strings.Contains(output, "repo.gitignore_exists") {
		t.Errorf("expected pass finding to be hidden by default, got %s", output)
	}
	if !strings.Contains(output, "docker.dockerfile_non_root_user") {
		t.Errorf("expected warning finding in text report, got %s", output)
	}
	if !strings.Contains(output, "docker.dockerfile_no_env_copy") {
		t.Errorf("expected failure finding in text report, got %s", output)
	}
}

func TestRenderText_VerboseRendersPassAndSkipStatuses(t *testing.T) {
	var buf bytes.Buffer

	if err := RenderTextWithOptions(&buf, sampleTextAuditReport(), TextOptions{Verbose: true}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Full findings (--verbose):") {
		t.Errorf("expected verbose full report heading, got %s", output)
	}
	if !strings.Contains(output, "✓ repo.gitignore_exists") {
		t.Errorf("expected pass finding in verbose report, got %s", output)
	}
	if !strings.Contains(output, "- k8s.manifest_exists") {
		t.Errorf("expected skipped finding to use neutral symbol in verbose report, got %s", output)
	}
}

func TestRenderText_ScoreLabels(t *testing.T) {
	tests := []struct {
		score int
		label string
	}{
		{score: 95, label: "Excellent"},
		{score: 80, label: "Good"},
		{score: 66, label: "Needs attention"},
		{score: 42, label: "High risk"},
	}

	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			report := sampleTextAuditReport()
			report.Score = tt.score
			var buf bytes.Buffer

			if err := RenderText(&buf, report); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !strings.Contains(buf.String(), tt.label) {
				t.Errorf("expected score label %q in output, got %s", tt.label, buf.String())
			}
		})
	}
}

func TestRenderText_RendersCategoryHeadingsAndRemediation(t *testing.T) {
	var buf bytes.Buffer

	if err := RenderText(&buf, sampleTextAuditReport()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "\nDocker\n") {
		t.Errorf("expected Docker category heading, got %s", output)
	}
	if !strings.Contains(output, "Fix: Add a USER instruction with a non-root user for better security") {
		t.Errorf("expected remediation for warning, got %s", output)
	}
	if !strings.Contains(output, "Fix: Avoid copying .env into Docker images") {
		t.Errorf("expected remediation for failure, got %s", output)
	}
	if !strings.Contains(output, "Passed and skipped checks are hidden by default") {
		t.Errorf("expected concise output hint, got %s", output)
	}
}

func TestRender_InvalidFormat(t *testing.T) {
	var buf bytes.Buffer

	err := Render(&buf, sampleAuditReport(), "xml")
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
	if !strings.Contains(err.Error(), "unknown report format") {
		t.Errorf("expected clear invalid format error, got %v", err)
	}
}

func sampleAuditReport() AuditReport {
	return AuditReport{
		Path:               ".",
		GitRepository:      true,
		FilesScanned:       10,
		DirectoriesScanned: 2,
		Score:              88,
		PassedCount:        1,
		WarningCount:       0,
		FailedCount:        0,
		SkippedCount:       0,
		Findings: []ReportFinding{
			{
				ID:          "docker.dockerfile_exists",
				Title:       "Dockerfile found",
				Category:    "docker",
				Severity:    rules.SeverityMedium,
				Status:      rules.StatusPass,
				Message:     "Dockerfile found",
				Remediation: "N/A",
			},
		},
	}
}

func sampleAuditReportWithSkip() AuditReport {
	report := sampleAuditReport()
	report.SkippedCount = 1
	report.Findings = append(report.Findings, ReportFinding{
		ID:          "docker.dockerfile_non_root_user",
		Title:       "Dockerfile not found",
		Category:    "docker",
		Severity:    rules.SeverityLow,
		Status:      rules.StatusSkip,
		Message:     "Dockerfile not present; skipping USER check",
		Remediation: "N/A",
	})
	return report
}

func sampleTextAuditReport() AuditReport {
	return AuditReport{
		Path:               ".",
		GitRepository:      true,
		FilesScanned:       10,
		DirectoriesScanned: 2,
		Score:              66,
		PassedCount:        1,
		WarningCount:       1,
		FailedCount:        1,
		SkippedCount:       1,
		Findings: []ReportFinding{
			{
				ID:          "repo.gitignore_exists",
				Title:       ".gitignore found",
				Category:    "repo",
				Severity:    rules.SeverityMedium,
				Status:      rules.StatusPass,
				Message:     ".gitignore file found",
				Remediation: "N/A",
			},
			{
				ID:          "docker.dockerfile_non_root_user",
				Title:       "No non-root USER instruction",
				Category:    "docker",
				Severity:    rules.SeverityHigh,
				Status:      rules.StatusWarn,
				Message:     "No non-root USER instruction detected in Dockerfile",
				Remediation: "Add a USER instruction with a non-root user for better security",
			},
			{
				ID:          "docker.dockerfile_no_env_copy",
				Title:       ".env copy detected",
				Category:    "docker",
				Severity:    rules.SeverityHigh,
				Status:      rules.StatusFail,
				Message:     "Dockerfile may be copying .env file directly",
				Remediation: "Avoid copying .env into Docker images; use environment variables or .env.example instead",
			},
			{
				ID:          "k8s.manifest_exists",
				Title:       "No Kubernetes manifest",
				Category:    "k8s",
				Severity:    rules.SeverityLow,
				Status:      rules.StatusSkip,
				Message:     "No Kubernetes manifest files found",
				Remediation: "N/A",
			},
		},
	}
}
