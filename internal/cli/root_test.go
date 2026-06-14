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

func TestRunWithWriter_AuditVerboseTextFormat(t *testing.T) {
	tmpDir := t.TempDir()
	writeCLIFile(t, filepath.Join(tmpDir, "README.md"), "# Test\n")

	var buf bytes.Buffer
	err := RunWithWriter([]string{"audit", tmpDir, "--verbose"}, &buf)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Full findings (--verbose):") {
		t.Errorf("expected verbose full report heading, got %s", output)
	}
	if !strings.Contains(output, "docs.readme_exists") {
		t.Errorf("expected verbose output to include passing finding, got %s", output)
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

func TestParseServeArgsDefaults(t *testing.T) {
	options, err := parseServeArgs(nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if options.path != "." {
		t.Errorf("expected default path '.', got %q", options.path)
	}
	if options.addr != "localhost:8080" {
		t.Errorf("expected default addr localhost:8080, got %q", options.addr)
	}
	if len(options.categories) != 0 {
		t.Errorf("expected no category filter, got %v", options.categories)
	}
}

func TestParseServeArgsWithFlags(t *testing.T) {
	options, err := parseServeArgs([]string{"examples/good-service", "--addr", "localhost:8081", "--category", "docker,ci"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if options.path != "examples/good-service" {
		t.Errorf("expected path examples/good-service, got %q", options.path)
	}
	if options.addr != "localhost:8081" {
		t.Errorf("expected addr localhost:8081, got %q", options.addr)
	}
	if got := strings.Join(options.categories, ","); got != "docker,ci" {
		t.Errorf("expected categories docker,ci, got %q", got)
	}
}

func TestRunWithWriter_ServeInvalidPath(t *testing.T) {
	missingPath := filepath.Join(t.TempDir(), "missing")

	var buf bytes.Buffer
	err := RunWithWriter([]string{"serve", missingPath}, &buf)
	if err == nil {
		t.Fatal("expected error for invalid serve path")
	}
	if !strings.Contains(err.Error(), "failed to access path") {
		t.Errorf("expected clear path error, got %v", err)
	}
	if strings.Contains(buf.String(), "dashboard running") {
		t.Errorf("expected no dashboard startup message, got %s", buf.String())
	}
}

func TestRunWithWriter_AuditFailUnderPrintsReportAndMessage(t *testing.T) {
	tmpDir := t.TempDir()
	writeCLIFile(t, filepath.Join(tmpDir, "README.md"), "# Test\n")

	var buf bytes.Buffer
	err := RunWithWriter([]string{"audit", tmpDir, "--fail-under", "101"}, &buf)
	if err == nil {
		t.Fatal("expected fail-under error")
	}

	output := buf.String()
	if !strings.Contains(output, "ShipCheck Deployment Readiness Report") {
		t.Errorf("expected report before fail-under message, got %s", output)
	}
	if !strings.Contains(output, "below fail-under threshold 101") {
		t.Errorf("expected fail-under message, got %s", output)
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
