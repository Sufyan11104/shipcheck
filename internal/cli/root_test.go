package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunWithWriter_AuditDefaultTextFormat(t *testing.T) {
	tmpDir := t.TempDir()
	writeCLIFile(t, filepath.Join(tmpDir, "README.md"), "# Test\n")

	var buf bytes.Buffer
	err := RunWithWriter([]string{"audit", tmpDir}, &buf)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "ShipCheck Deployment Readiness Report") {
		t.Errorf("expected text report output, got %s", output)
	}
}

func TestRunWithWriter_AuditInvalidFormat(t *testing.T) {
	tmpDir := t.TempDir()

	var buf bytes.Buffer
	err := RunWithWriter([]string{"audit", tmpDir, "--format", "xml"}, &buf)
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
	if !strings.Contains(err.Error(), "unknown report format") {
		t.Errorf("expected clear invalid format error, got %v", err)
	}
}

func TestRunWithWriter_AuditUnknownCategory(t *testing.T) {
	tmpDir := t.TempDir()

	var buf bytes.Buffer
	err := RunWithWriter([]string{"audit", tmpDir, "--category", "docker,unknown"}, &buf)
	if err == nil {
		t.Fatal("expected error for unknown category")
	}
	if !strings.Contains(err.Error(), "unknown category") {
		t.Errorf("expected clear unknown category error, got %v", err)
	}
}

func TestEvaluateFailUnder(t *testing.T) {
	if err := EvaluateFailUnder(80, 80); err != nil {
		t.Fatalf("expected no error at threshold, got %v", err)
	}

	err := EvaluateFailUnder(79, 80)
	if err == nil {
		t.Fatal("expected fail-under error")
	}

	exitErr, ok := err.(interface{ ExitCode() int })
	if !ok {
		t.Fatalf("expected error with ExitCode, got %T", err)
	}
	if exitErr.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d", exitErr.ExitCode())
	}
	if !strings.Contains(err.Error(), "Score 79 is below fail-under threshold 80") {
		t.Errorf("expected threshold message, got %v", err)
	}
}

func writeCLIFile(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("failed to create parent directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
}
