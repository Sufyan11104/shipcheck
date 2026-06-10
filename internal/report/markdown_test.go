package report

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderMarkdown_IncludesScoreAndFindings(t *testing.T) {
	var buf bytes.Buffer

	if err := RenderMarkdown(&buf, sampleAuditReport()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "# ShipCheck Deployment Readiness Report") {
		t.Errorf("expected Markdown heading, got %s", output)
	}
	if !strings.Contains(output, "| Score | 88/100 |") {
		t.Errorf("expected score in summary table, got %s", output)
	}
	if !strings.Contains(output, "docker.dockerfile_exists") {
		t.Errorf("expected finding ID in findings table, got %s", output)
	}
}
