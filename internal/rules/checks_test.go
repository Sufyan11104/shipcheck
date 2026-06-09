package rules

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckReadmeExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when README doesn't exist
	finding := CheckReadmeExists(tmpDir)
	if finding.Status != StatusFail {
		t.Errorf("expected StatusFail, got %v", finding.Status)
	}

	// Test when README.md exists
	readmePath := filepath.Join(tmpDir, "README.md")
	os.Create(readmePath)
	finding = CheckReadmeExists(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass, got %v", finding.Status)
	}
}

func TestCheckGitignoreExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when .gitignore doesn't exist
	finding := CheckGitignoreExists(tmpDir)
	if finding.Status != StatusFail {
		t.Errorf("expected StatusFail, got %v", finding.Status)
	}

	// Test when .gitignore exists
	gitignorePath := filepath.Join(tmpDir, ".gitignore")
	os.Create(gitignorePath)
	finding = CheckGitignoreExists(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass, got %v", finding.Status)
	}
}

func TestCheckEnvNotCommitted(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when not a git repo
	finding := CheckEnvNotCommitted(tmpDir, false)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for non-git repo, got %v", finding.Status)
	}

	// Test when .env doesn't exist
	finding = CheckEnvNotCommitted(tmpDir, true)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass when .env missing, got %v", finding.Status)
	}

	// Test when .env exists
	envPath := filepath.Join(tmpDir, ".env")
	os.Create(envPath)
	finding = CheckEnvNotCommitted(tmpDir, true)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn when .env exists, got %v", finding.Status)
	}
}

func TestCheckEnvExampleExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when .env.example doesn't exist
	finding := CheckEnvExampleExists(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn, got %v", finding.Status)
	}

	// Test when .env.example exists
	envExamplePath := filepath.Join(tmpDir, ".env.example")
	os.Create(envExamplePath)
	finding = CheckEnvExampleExists(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass, got %v", finding.Status)
	}
}
