package engine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Sufyan11104/shipcheck/internal/rules"
)

func TestEngine_RunChecks_AllPass(t *testing.T) {
	tmpDir := t.TempDir()

	// Create all required files
	os.Create(filepath.Join(tmpDir, "README.md"))
	os.Create(filepath.Join(tmpDir, ".gitignore"))
	os.Create(filepath.Join(tmpDir, ".env.example"))

	// Create a proper Dockerfile with all good practices
	dockerfileContent := `FROM alpine:latest
RUN adduser -D appuser
HEALTHCHECK CMD echo 'ok'
COPY app /app
ENV APP_NAME=myapp
USER appuser
`
	os.WriteFile(filepath.Join(tmpDir, "Dockerfile"), []byte(dockerfileContent), 0644)
	os.Create(filepath.Join(tmpDir, ".dockerignore"))

	// Create GitHub Actions workflow
	workflowsPath := filepath.Join(tmpDir, ".github", "workflows")
	os.MkdirAll(workflowsPath, 0755)
	workflowContent := `on: push
permissions:
  contents: read
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: go test ./...
      - run: go build ./...
`
	os.WriteFile(filepath.Join(workflowsPath, "ci.yml"), []byte(workflowContent), 0644)

	engine := NewEngine(tmpDir)
	findings, score := engine.RunChecks(true)

	if len(findings) != 18 {
		t.Errorf("expected 18 findings, got %d", len(findings))
	}

	// Check that we have the right findings
	expectedIDs := map[string]bool{
		"docs.readme_exists":              false,
		"repo.gitignore_exists":           false,
		"env.env_not_committed":           false,
		"env.env_example_exists":          false,
		"docker.dockerfile_exists":        false,
		"docker.dockerignore_exists":      false,
		"docker.dockerfile_non_root_user": false,
		"docker.dockerfile_healthcheck":   false,
		"docker.dockerfile_no_env_copy":   false,
		"docker.dockerfile_no_secret_env": false,
		"ci.workflows_dir_exists":         false,
		"ci.workflow_file_exists":         false,
		"ci.test_step_exists":             false,
		"ci.build_step_exists":            false,
		"ci.deploy_after_tests":           false,
		"ci.actions_pinned":               false,
		"ci.no_secret_echo":               false,
		"ci.permissions_declared":         false,
	}

	for _, f := range findings {
		if _, exists := expectedIDs[f.ID]; exists {
			expectedIDs[f.ID] = true
		}
	}

	for id, found := range expectedIDs {
		if !found {
			t.Errorf("expected finding %s not found", id)
		}
	}

	// Score should be 100 if all pass
	if score != 100 {
		t.Errorf("expected score 100, got %d", score)
	}
}

func TestEngine_RunChecks_WithWarnings(t *testing.T) {
	tmpDir := t.TempDir()

	// Create some files but not .env.example or Dockerfile
	os.Create(filepath.Join(tmpDir, "README.md"))
	os.Create(filepath.Join(tmpDir, ".gitignore"))

	engine := NewEngine(tmpDir)
	findings, score := engine.RunChecks(true)

	// Find the .env.example finding
	var envExampleFinding *rules.Finding
	for i, f := range findings {
		if f.ID == "env.env_example_exists" {
			envExampleFinding = &findings[i]
			break
		}
	}

	if envExampleFinding == nil {
		t.Fatal("expected .env.example finding")
	}

	if envExampleFinding.Status != rules.StatusWarn {
		t.Errorf("expected warn status for missing .env.example, got %v", envExampleFinding.Status)
	}

	// Score should be less than 100
	if score >= 100 {
		t.Errorf("expected score < 100 with warnings, got %d", score)
	}
}

func TestEngine_RunChecks_NotGitRepo(t *testing.T) {
	tmpDir := t.TempDir()

	os.Create(filepath.Join(tmpDir, "README.md"))
	os.Create(filepath.Join(tmpDir, ".gitignore"))
	os.Create(filepath.Join(tmpDir, ".env.example"))

	engine := NewEngine(tmpDir)
	findings, _ := engine.RunChecks(false)

	// Find the env check
	var envFinding *rules.Finding
	for i, f := range findings {
		if f.ID == "env.env_not_committed" {
			envFinding = &findings[i]
			break
		}
	}

	if envFinding == nil {
		t.Fatal("expected env finding")
	}

	if envFinding.Status != rules.StatusWarn {
		t.Errorf("expected warn status for non-git repo, got %v", envFinding.Status)
	}
}
