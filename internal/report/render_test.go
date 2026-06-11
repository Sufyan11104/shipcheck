package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Sufyan11104/shipcheck/internal/rules"
)

func TestRenderText_DefaultFormatWorks(t *testing.T) {
	var buf bytes.Buffer

	if err := Render(&buf, sampleAuditReport(), FormatText); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "ShipCheck Deployment Readiness Report") {
		t.Errorf("expected text report heading, got %s", output)
	}
	if !strings.Contains(output, "docker.dockerfile_exists") {
		t.Errorf("expected finding ID in text report, got %s", output)
	}
}

func TestRenderText_RendersSkipAsNeutralStatus(t *testing.T) {
	var buf bytes.Buffer

	if err := RenderText(&buf, sampleAuditReportWithSkip()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Skipped: 1") {
		t.Errorf("expected skipped count in text report, got %s", output)
	}
	if !strings.Contains(output, "- docker.dockerfile_non_root_user - Dockerfile not present; skipping USER check") {
		t.Errorf("expected skipped finding to use neutral symbol, got %s", output)
	}
	if strings.Contains(output, "✓ docker.dockerfile_non_root_user") {
		t.Errorf("expected skipped finding not to render as passed, got %s", output)
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
