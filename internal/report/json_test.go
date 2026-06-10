package report

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestRenderJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer

	if err := RenderJSON(&buf, sampleAuditReport()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !json.Valid(buf.Bytes()) {
		t.Fatalf("expected valid JSON, got %s", buf.String())
	}
}

func TestRenderJSON_IncludesScoreAndFindings(t *testing.T) {
	var buf bytes.Buffer

	if err := RenderJSON(&buf, sampleAuditReport()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var decoded AuditReport
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if decoded.Score != 88 {
		t.Errorf("expected score 88, got %d", decoded.Score)
	}
	if len(decoded.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(decoded.Findings))
	}
	if decoded.Findings[0].ID != "docker.dockerfile_exists" {
		t.Errorf("expected docker finding, got %s", decoded.Findings[0].ID)
	}
}
