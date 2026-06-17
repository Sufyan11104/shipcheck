package report

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/Sufyan11104/shipcheck/internal/rules"
)

func TestRenderSARIF_Golden(t *testing.T) {
	var buf bytes.Buffer

	if err := RenderSARIF(&buf, sampleSARIFAuditReport()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	golden, err := os.ReadFile("testdata/sample.sarif.golden")
	if err != nil {
		t.Fatalf("failed to read golden SARIF fixture: %v", err)
	}

	if buf.String() != string(golden) {
		t.Fatalf("SARIF output differed from golden fixture\nwant:\n%s\ngot:\n%s", string(golden), buf.String())
	}
}

func TestRenderSARIF_MetadataAndRuleDescriptors(t *testing.T) {
	log := renderSARIFForTest(t, sampleSARIFAuditReport())

	if log.Version != "2.1.0" {
		t.Errorf("expected SARIF version 2.1.0, got %q", log.Version)
	}
	if log.Schema != sarifSchemaURL {
		t.Errorf("expected schema %q, got %q", sarifSchemaURL, log.Schema)
	}
	if len(log.Runs) != 1 {
		t.Fatalf("expected one SARIF run, got %d", len(log.Runs))
	}

	driver := log.Runs[0].Tool.Driver
	if driver.Name != "ShipCheck" {
		t.Errorf("expected tool name ShipCheck, got %q", driver.Name)
	}
	if driver.InformationURI != shipCheckInfoURI {
		t.Errorf("expected tool info URI %q, got %q", shipCheckInfoURI, driver.InformationURI)
	}
	if driver.Version == "" {
		t.Error("expected tool version")
	}
	if len(driver.Rules) != 2 {
		t.Fatalf("expected descriptors for warning/failure rules only, got %d", len(driver.Rules))
	}
	if driver.Rules[0].ID != "docker.dockerfile_non_root_user" || driver.Rules[1].ID != "env.env_not_committed" {
		t.Errorf("expected stable rule descriptor order, got %s then %s", driver.Rules[0].ID, driver.Rules[1].ID)
	}
	if driver.Rules[0].Properties.Category != "docker" {
		t.Errorf("expected descriptor category docker, got %q", driver.Rules[0].Properties.Category)
	}
}

func TestRenderSARIF_MapsFindingsToResults(t *testing.T) {
	log := renderSARIFForTest(t, sampleSARIFAuditReport())
	results := log.Runs[0].Results

	if len(results) != 2 {
		t.Fatalf("expected pass and skip findings to be omitted, got %d results", len(results))
	}

	warning := results[0]
	if warning.RuleID != "docker.dockerfile_non_root_user" {
		t.Errorf("expected warning result rule ID, got %q", warning.RuleID)
	}
	if warning.Level != "warning" {
		t.Errorf("expected medium severity to map to SARIF warning, got %q", warning.Level)
	}
	if got := warning.Properties["shipcheckStatus"]; got != "warn" {
		t.Errorf("expected original status warn, got %v", got)
	}
	if len(warning.Locations) != 1 || warning.Locations[0].PhysicalLocation.ArtifactLocation.URI != "services/api/Dockerfile" {
		t.Errorf("expected normalized relative location, got %+v", warning.Locations)
	}

	failure := results[1]
	if failure.RuleID != "env.env_not_committed" {
		t.Errorf("expected failure result rule ID, got %q", failure.RuleID)
	}
	if failure.Level != "error" {
		t.Errorf("expected high severity to map to SARIF error, got %q", failure.Level)
	}
	if got := failure.Properties["shipcheckSeverity"]; got != "high" {
		t.Errorf("expected original severity high, got %v", got)
	}
	if evidence, ok := failure.Properties["evidence"].(string); !ok || evidence != "SECRET_KEY=[redacted] PUBLIC=ok" {
		t.Errorf("expected sanitized evidence, got %#v", failure.Properties["evidence"])
	}
}

func TestRenderSARIF_SeverityToLevelMapping(t *testing.T) {
	tests := []struct {
		severity rules.Severity
		want     string
	}{
		{severity: rules.SeverityInfo, want: "note"},
		{severity: rules.SeverityLow, want: "note"},
		{severity: rules.SeverityMedium, want: "warning"},
		{severity: rules.SeverityHigh, want: "error"},
		{severity: rules.Severity("critical"), want: "error"},
	}

	for _, tt := range tests {
		t.Run(string(tt.severity), func(t *testing.T) {
			if got := sarifLevel(tt.severity); got != tt.want {
				t.Errorf("expected %s to map to %s, got %s", tt.severity, tt.want, got)
			}
		})
	}
}

func TestRenderSARIF_SafeLocationHandling(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantURI  string
		wantNone bool
	}{
		{name: "repository wide", path: "", wantNone: true},
		{name: "absolute unix", path: "/Users/sufyan/project/Dockerfile", wantNone: true},
		{name: "absolute windows", path: `C:\Users\sufyan\project\Dockerfile`, wantNone: true},
		{name: "traversal", path: "../Dockerfile", wantNone: true},
		{name: "windows separator", path: `services\api\Dockerfile`, wantURI: "services/api/Dockerfile"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := AuditReport{
				Findings: []ReportFinding{
					{
						ID:          "test.rule",
						Title:       "Test rule",
						Category:    "repo",
						Severity:    rules.SeverityMedium,
						Status:      rules.StatusWarn,
						Message:     "Test warning",
						Remediation: "Fix it",
						Path:        tt.path,
					},
				},
			}
			log := renderSARIFForTest(t, report)
			locations := log.Runs[0].Results[0].Locations
			if tt.wantNone {
				if len(locations) != 0 {
					t.Fatalf("expected no SARIF location, got %+v", locations)
				}
				return
			}
			if len(locations) != 1 {
				t.Fatalf("expected one SARIF location, got %d", len(locations))
			}
			if got := locations[0].PhysicalLocation.ArtifactLocation.URI; got != tt.wantURI {
				t.Errorf("expected URI %q, got %q", tt.wantURI, got)
			}
		})
	}
}

func TestRenderSARIF_DeterministicOutput(t *testing.T) {
	var first bytes.Buffer
	var second bytes.Buffer

	if err := RenderSARIF(&first, sampleSARIFAuditReport()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := RenderSARIF(&second, sampleSARIFAuditReport()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if first.String() != second.String() {
		t.Fatalf("expected deterministic SARIF output")
	}
}

func TestRenderSARIF_OmitsAbsolutePaths(t *testing.T) {
	var buf bytes.Buffer
	report := AuditReport{
		Findings: []ReportFinding{
			{
				ID:       "env.env_not_committed",
				Title:    ".env file detected",
				Category: "env",
				Severity: rules.SeverityHigh,
				Status:   rules.StatusFail,
				Message:  ".env file found",
				Path:     "/Users/sufyanelmansuri/project/.env",
			},
		},
	}

	if err := RenderSARIF(&buf, report); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := buf.String()
	if strings.Contains(output, "/Users/") || strings.Contains(output, "sufyanelmansuri") {
		t.Fatalf("expected absolute local path to be omitted, got %s", output)
	}
}

func renderSARIFForTest(t *testing.T, auditReport AuditReport) sarifLog {
	t.Helper()

	var buf bytes.Buffer
	if err := RenderSARIF(&buf, auditReport); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var log sarifLog
	if err := json.Unmarshal(buf.Bytes(), &log); err != nil {
		t.Fatalf("expected valid SARIF JSON, got %v\n%s", err, buf.String())
	}
	return log
}

func sampleSARIFAuditReport() AuditReport {
	return AuditReport{
		Path:               ".",
		GitRepository:      true,
		FilesScanned:       4,
		DirectoriesScanned: 1,
		Score:              62,
		PassedCount:        1,
		WarningCount:       1,
		FailedCount:        1,
		SkippedCount:       1,
		Findings: []ReportFinding{
			{
				ID:       "repo.gitignore_exists",
				Title:    ".gitignore found",
				Category: "repo",
				Severity: rules.SeverityMedium,
				Status:   rules.StatusPass,
				Message:  ".gitignore file found",
			},
			{
				ID:          "docker.dockerfile_non_root_user",
				Title:       "No non-root USER instruction",
				Category:    "docker",
				Severity:    rules.SeverityMedium,
				Status:      rules.StatusWarn,
				Message:     "No non-root USER instruction detected in Dockerfile",
				Remediation: "Add a non-root user and switch to it with USER before the final command.",
				Path:        `services\api\Dockerfile`,
				Evidence:    "USER root",
			},
			{
				ID:          "env.env_not_committed",
				Title:       ".env file detected",
				Category:    "env",
				Severity:    rules.SeverityHigh,
				Status:      rules.StatusFail,
				Message:     ".env file found - ensure it is in .gitignore to prevent accidental commits",
				Remediation: "Add .env to .gitignore and use .env.example for template",
				Path:        ".env",
				Evidence:    "SECRET_KEY=abc123 PUBLIC=ok",
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
