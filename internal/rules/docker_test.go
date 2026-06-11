package rules

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckDockerfileExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when Dockerfile doesn't exist
	finding := CheckDockerfileExists(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn, got %v", finding.Status)
	}

	// Test when Dockerfile exists
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	os.Create(dockerfilePath)
	finding = CheckDockerfileExists(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass, got %v", finding.Status)
	}
}

func TestCheckDockerignoreExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when .dockerignore doesn't exist
	finding := CheckDockerignoreExists(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn, got %v", finding.Status)
	}

	// Test when .dockerignore exists
	dockerignorePath := filepath.Join(tmpDir, ".dockerignore")
	os.Create(dockerignorePath)
	finding = CheckDockerignoreExists(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass, got %v", finding.Status)
	}
}

func TestCheckDockerfileNonRootUser(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when Dockerfile doesn't exist
	finding := CheckDockerfileNonRootUser(tmpDir)
	if finding.Status != StatusSkip {
		t.Errorf("expected StatusSkip when Dockerfile missing, got %v", finding.Status)
	}

	// Test when Dockerfile has no USER instruction
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	os.WriteFile(dockerfilePath, []byte("FROM alpine:latest\nRUN echo 'hello'\n"), 0644)
	finding = CheckDockerfileNonRootUser(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for no USER, got %v", finding.Status)
	}

	// Test when Dockerfile has non-root USER
	os.WriteFile(dockerfilePath, []byte("FROM alpine:latest\nRUN adduser -D myuser\nUSER myuser\n"), 0644)
	finding = CheckDockerfileNonRootUser(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for non-root USER, got %v", finding.Status)
	}

	// Test when Dockerfile has root USER
	os.WriteFile(dockerfilePath, []byte("FROM alpine:latest\nUSER root\n"), 0644)
	finding = CheckDockerfileNonRootUser(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for root USER, got %v", finding.Status)
	}
}

func TestCheckDockerfileHealthcheck(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when Dockerfile doesn't exist
	finding := CheckDockerfileHealthcheck(tmpDir)
	if finding.Status != StatusSkip {
		t.Errorf("expected StatusSkip when Dockerfile missing, got %v", finding.Status)
	}

	// Test when Dockerfile has no HEALTHCHECK
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	os.WriteFile(dockerfilePath, []byte("FROM alpine:latest\nRUN echo 'hello'\n"), 0644)
	finding = CheckDockerfileHealthcheck(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for no HEALTHCHECK, got %v", finding.Status)
	}

	// Test when Dockerfile has HEALTHCHECK
	os.WriteFile(dockerfilePath, []byte("FROM alpine:latest\nHEALTHCHECK CMD echo 'ok'\n"), 0644)
	finding = CheckDockerfileHealthcheck(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for HEALTHCHECK, got %v", finding.Status)
	}
}

func TestCheckDockerfileNoEnvCopy(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when Dockerfile doesn't exist
	finding := CheckDockerfileNoEnvCopy(tmpDir)
	if finding.Status != StatusSkip {
		t.Errorf("expected StatusSkip when Dockerfile missing, got %v", finding.Status)
	}

	// Test when Dockerfile doesn't copy .env
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	os.WriteFile(dockerfilePath, []byte("FROM alpine:latest\nCOPY . /app\n"), 0644)
	finding = CheckDockerfileNoEnvCopy(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for no .env copy, got %v", finding.Status)
	}

	// Test when Dockerfile copies .env
	os.WriteFile(dockerfilePath, []byte("FROM alpine:latest\nCOPY .env /app/.env\n"), 0644)
	finding = CheckDockerfileNoEnvCopy(tmpDir)
	if finding.Status != StatusFail {
		t.Errorf("expected StatusFail for .env copy, got %v", finding.Status)
	}
}

func TestCheckDockerfileNoSecretEnv(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when Dockerfile doesn't exist
	finding := CheckDockerfileNoSecretEnv(tmpDir)
	if finding.Status != StatusSkip {
		t.Errorf("expected StatusSkip when Dockerfile missing, got %v", finding.Status)
	}

	// Test when Dockerfile has no secret ENV
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	os.WriteFile(dockerfilePath, []byte("FROM alpine:latest\nENV APP_NAME=myapp\n"), 0644)
	finding = CheckDockerfileNoSecretEnv(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for no secret ENV, got %v", finding.Status)
	}

	// Test when Dockerfile has PASSWORD
	os.WriteFile(dockerfilePath, []byte("FROM alpine:latest\nENV DB_PASSWORD=secret123\n"), 0644)
	finding = CheckDockerfileNoSecretEnv(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for PASSWORD, got %v", finding.Status)
	}

	// Test when Dockerfile has SECRET
	os.WriteFile(dockerfilePath, []byte("FROM alpine:latest\nARG SECRET_KEY=mykey\n"), 0644)
	finding = CheckDockerfileNoSecretEnv(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for SECRET, got %v", finding.Status)
	}

	// Test when Dockerfile has TOKEN
	os.WriteFile(dockerfilePath, []byte("FROM alpine:latest\nENV API_TOKEN=abc123\n"), 0644)
	finding = CheckDockerfileNoSecretEnv(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for TOKEN, got %v", finding.Status)
	}

	// Test when Dockerfile has API_KEY
	os.WriteFile(dockerfilePath, []byte("FROM alpine:latest\nARG API_KEY=xyz789\n"), 0644)
	finding = CheckDockerfileNoSecretEnv(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for API_KEY, got %v", finding.Status)
	}

	// Test when Dockerfile has PRIVATE_KEY
	os.WriteFile(dockerfilePath, []byte("FROM alpine:latest\nENV PRIVATE_KEY=-----BEGIN\n"), 0644)
	finding = CheckDockerfileNoSecretEnv(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for PRIVATE_KEY, got %v", finding.Status)
	}
}
