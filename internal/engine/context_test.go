package engine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Sufyan11104/shipcheck/internal/rules"
)

func TestDefaultAudit_DockerOnlyRepoSkipsUnusedOptionalCategories(t *testing.T) {
	tmpDir := t.TempDir()
	writeContextTestFile(t, tmpDir, "README.md", "# Docker service\n")
	writeContextTestFile(t, tmpDir, ".gitignore", ".env\n")
	writeContextTestFile(t, tmpDir, ".env.example", "APP_ENV=local\n")
	writeContextTestFile(t, tmpDir, ".dockerignore", ".env\n")
	writeContextTestFile(t, tmpDir, "Dockerfile", `FROM alpine:3.20
RUN adduser -D app
HEALTHCHECK CMD echo ok
USER app
CMD ["sh"]
`)

	findings, score := NewEngine(tmpDir).RunChecks(true)

	if score != 100 {
		t.Fatalf("expected Docker-only repo score 100 without unused category penalties, got %d", score)
	}
	assertFindingStatus(t, findings, "docker.dockerfile_exists", rules.StatusPass)
	assertFindingStatus(t, findings, "ci.workflows_dir_exists", rules.StatusSkip)
	assertFindingStatus(t, findings, "k8s.manifest_exists", rules.StatusSkip)
	assertFindingStatus(t, findings, "terraform.files_exist", rules.StatusSkip)
}

func TestDefaultAudit_TerraformOnlyRepoSkipsUnusedOptionalCategories(t *testing.T) {
	tmpDir := t.TempDir()
	writeContextTestFile(t, tmpDir, "README.md", "# Terraform module\n")
	writeContextTestFile(t, tmpDir, ".gitignore", ".env\n.terraform/\n")
	writeContextTestFile(t, tmpDir, ".env.example", "APP_ENV=local\n")
	writeContextTestFile(t, tmpDir, "Makefile", "fmt:\n\tterraform fmt\nvalidate:\n\tterraform validate\n")
	writeContextTestFile(t, tmpDir, "versions.tf", `terraform {
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "~> 2.5"
    }
  }
}
`)
	writeContextTestFile(t, tmpDir, "main.tf", `terraform {
  backend "local" {
    path = "terraform.tfstate"
  }
}

provider "local" {}

variable "environment" {
  type    = string
  default = "demo"
}
`)
	writeContextTestFile(t, tmpDir, ".terraform.lock.hcl", `provider "registry.terraform.io/hashicorp/local" {
  version     = "2.5.1"
  constraints = "~> 2.5"
}
`)

	findings, score := NewEngine(tmpDir).RunChecks(true)

	if score != 100 {
		t.Fatalf("expected Terraform-only repo score 100 without unused category penalties, got %d", score)
	}
	assertFindingStatus(t, findings, "terraform.files_exist", rules.StatusPass)
	assertFindingStatus(t, findings, "docker.dockerfile_exists", rules.StatusSkip)
	assertFindingStatus(t, findings, "ci.workflows_dir_exists", rules.StatusSkip)
	assertFindingStatus(t, findings, "k8s.manifest_exists", rules.StatusSkip)
}

func TestExplicitCategoryAudit_MissingCategoryEvidenceStillWarns(t *testing.T) {
	tests := []struct {
		name       string
		category   string
		findingID  string
		wantStatus rules.Status
	}{
		{name: "docker", category: "docker", findingID: "docker.dockerfile_exists", wantStatus: rules.StatusWarn},
		{name: "ci", category: "ci", findingID: "ci.workflows_dir_exists", wantStatus: rules.StatusWarn},
		{name: "k8s", category: "k8s", findingID: "k8s.manifest_exists", wantStatus: rules.StatusWarn},
		{name: "terraform", category: "terraform", findingID: "terraform.files_exist", wantStatus: rules.StatusWarn},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings, _ := NewEngine(t.TempDir()).RunChecksWithCategories(true, []string{tt.category})
			findings = FilterFindingsByCategory(findings, []string{tt.category})

			assertFindingStatus(t, findings, tt.findingID, tt.wantStatus)
		})
	}
}

func writeContextTestFile(t *testing.T, root, name, content string) {
	t.Helper()

	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file %s: %v", name, err)
	}
}

func assertFindingStatus(t *testing.T, findings []rules.Finding, id string, status rules.Status) {
	t.Helper()

	for _, finding := range findings {
		if finding.ID == id {
			if finding.Status != status {
				t.Fatalf("expected %s status %s, got %s", id, status, finding.Status)
			}
			return
		}
	}

	t.Fatalf("expected finding %s", id)
}
