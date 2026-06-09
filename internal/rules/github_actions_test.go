package rules

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckWorkflowsDirExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when .github/workflows doesn't exist
	finding := CheckWorkflowsDirExists(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn, got %v", finding.Status)
	}

	// Test when .github/workflows exists
	workflowsPath := filepath.Join(tmpDir, ".github", "workflows")
	os.MkdirAll(workflowsPath, 0755)
	finding = CheckWorkflowsDirExists(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass, got %v", finding.Status)
	}
}

func TestCheckWorkflowFileExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when .github/workflows doesn't exist
	finding := CheckWorkflowFileExists(tmpDir)
	if finding.Status != StatusPass { // Graceful skip
		t.Errorf("expected StatusPass (graceful skip), got %v", finding.Status)
	}

	// Test when directory exists but no workflow files
	workflowsPath := filepath.Join(tmpDir, ".github", "workflows")
	os.MkdirAll(workflowsPath, 0755)
	finding = CheckWorkflowFileExists(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn, got %v", finding.Status)
	}

	// Test when workflow file exists
	os.WriteFile(filepath.Join(workflowsPath, "test.yml"), []byte("on: push\n"), 0644)
	finding = CheckWorkflowFileExists(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass, got %v", finding.Status)
	}
}

func TestCheckTestStepExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when no workflows exist
	finding := CheckTestStepExists(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for no workflows, got %v", finding.Status)
	}

	// Test when workflow exists but no test step
	workflowsPath := filepath.Join(tmpDir, ".github", "workflows")
	os.MkdirAll(workflowsPath, 0755)
	os.WriteFile(filepath.Join(workflowsPath, "build.yml"), []byte("on: push\njobs:\n  build:\n    runs-on: ubuntu-latest\n"), 0644)
	finding = CheckTestStepExists(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for no test step, got %v", finding.Status)
	}

	// Test when workflow has test step (go test)
	os.WriteFile(filepath.Join(workflowsPath, "test.yml"), []byte("on: push\njobs:\n  test:\n    runs-on: ubuntu-latest\n    steps:\n      - run: go test ./...\n"), 0644)
	finding = CheckTestStepExists(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for test step, got %v", finding.Status)
	}
}

func TestCheckBuildStepExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when no workflows exist
	finding := CheckBuildStepExists(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for no workflows, got %v", finding.Status)
	}

	// Test when workflow has build step
	workflowsPath := filepath.Join(tmpDir, ".github", "workflows")
	os.MkdirAll(workflowsPath, 0755)
	os.WriteFile(filepath.Join(workflowsPath, "build.yml"), []byte("on: push\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n      - run: go build ./...\n"), 0644)
	finding = CheckBuildStepExists(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for build step, got %v", finding.Status)
	}
}

func TestCheckDeployAfterTests(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when no workflows exist
	finding := CheckDeployAfterTests(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass (graceful skip), got %v", finding.Status)
	}

	// Test when workflow has correct order (test before deploy)
	workflowsPath := filepath.Join(tmpDir, ".github", "workflows")
	os.MkdirAll(workflowsPath, 0755)
	os.WriteFile(filepath.Join(workflowsPath, "ci.yml"), []byte("on: push\njobs:\n  test:\n    runs-on: ubuntu-latest\n    steps:\n      - run: go test\n      - run: deploy\n"), 0644)
	finding = CheckDeployAfterTests(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for correct order, got %v", finding.Status)
	}

	// Test when workflow has wrong order (deploy before test)
	os.WriteFile(filepath.Join(workflowsPath, "ci.yml"), []byte("on: push\njobs:\n  deploy:\n    runs-on: ubuntu-latest\n    steps:\n      - run: deploy\n      - run: go test\n"), 0644)
	finding = CheckDeployAfterTests(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for deploy before test, got %v", finding.Status)
	}
}

func TestCheckActionsPinned(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when no workflows exist
	finding := CheckActionsPinned(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass (no actions), got %v", finding.Status)
	}

	// Test when workflow has pinned actions
	workflowsPath := filepath.Join(tmpDir, ".github", "workflows")
	os.MkdirAll(workflowsPath, 0755)
	os.WriteFile(filepath.Join(workflowsPath, "ci.yml"), []byte("on: push\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v4\n"), 0644)
	finding = CheckActionsPinned(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for pinned actions, got %v", finding.Status)
	}

	// Test when workflow has unpinned actions
	os.WriteFile(filepath.Join(workflowsPath, "ci.yml"), []byte("on: push\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout\n"), 0644)
	finding = CheckActionsPinned(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for unpinned actions, got %v", finding.Status)
	}
}

func TestCheckNoSecretEcho(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when no workflows exist
	finding := CheckNoSecretEcho(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass (no secrets), got %v", finding.Status)
	}

	// Test when workflow doesn't echo secrets
	workflowsPath := filepath.Join(tmpDir, ".github", "workflows")
	os.MkdirAll(workflowsPath, 0755)
	os.WriteFile(filepath.Join(workflowsPath, "ci.yml"), []byte("on: push\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo \"Hello\"\n"), 0644)
	finding = CheckNoSecretEcho(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for no secret echo, got %v", finding.Status)
	}

	// Test when workflow echoes secrets
	os.WriteFile(filepath.Join(workflowsPath, "ci.yml"), []byte("on: push\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo ${{ secrets.TOKEN }}\n"), 0644)
	finding = CheckNoSecretEcho(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for secret echo, got %v", finding.Status)
	}
}

func TestCheckPermissionsDeclared(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when no workflows exist
	finding := CheckPermissionsDeclared(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass (no workflows), got %v", finding.Status)
	}

	// Test when workflow lacks permissions block
	workflowsPath := filepath.Join(tmpDir, ".github", "workflows")
	os.MkdirAll(workflowsPath, 0755)
	os.WriteFile(filepath.Join(workflowsPath, "ci.yml"), []byte("on: push\njobs:\n  build:\n    runs-on: ubuntu-latest\n"), 0644)
	finding = CheckPermissionsDeclared(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for no permissions block, got %v", finding.Status)
	}

	// Test when workflow has permissions block
	os.WriteFile(filepath.Join(workflowsPath, "ci.yml"), []byte("on: push\npermissions:\n  contents: read\njobs:\n  build:\n    runs-on: ubuntu-latest\n"), 0644)
	finding = CheckPermissionsDeclared(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for permissions block, got %v", finding.Status)
	}
}
