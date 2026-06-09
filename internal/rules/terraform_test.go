package rules

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckTerraformFilesExist_NoTerraformFiles(t *testing.T) {
	tmpDir := t.TempDir()

	finding := CheckTerraformFilesExist(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn, got %v", finding.Status)
	}
}

func TestCheckTerraformFilesExist_TfFileDetected(t *testing.T) {
	tmpDir := t.TempDir()
	writeTestFile(t, filepath.Join(tmpDir, "main.tf"), `resource "null_resource" "example" {}`)

	finding := CheckTerraformFilesExist(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass, got %v", finding.Status)
	}
}

func TestCheckTerraformFilesExist_TfvarsFileDetected(t *testing.T) {
	tmpDir := t.TempDir()
	writeTestFile(t, filepath.Join(tmpDir, "terraform.tfvars"), `environment = "dev"`)

	finding := CheckTerraformFilesExist(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass, got %v", finding.Status)
	}

	deeperFinding := CheckTerraformRequiredProvidersExists(tmpDir)
	if deeperFinding.Status != StatusPass {
		t.Errorf("expected StatusPass skip for tfvars-only directory, got %v", deeperFinding.Status)
	}
}

func TestCheckTerraformFmtAndValidateRecommended(t *testing.T) {
	tmpDir := t.TempDir()
	writeTestFile(t, filepath.Join(tmpDir, "main.tf"), `resource "null_resource" "example" {}`)

	fmtFinding := CheckTerraformFmtRecommended(tmpDir)
	if fmtFinding.Status != StatusWarn || fmtFinding.Severity != SeverityInfo {
		t.Errorf("expected info warning for missing terraform fmt automation, got %v/%v", fmtFinding.Status, fmtFinding.Severity)
	}

	validateFinding := CheckTerraformValidateRecommended(tmpDir)
	if validateFinding.Status != StatusWarn || validateFinding.Severity != SeverityInfo {
		t.Errorf("expected info warning for missing terraform validate automation, got %v/%v", validateFinding.Status, validateFinding.Severity)
	}

	workflowPath := filepath.Join(tmpDir, ".github", "workflows", "terraform.yml")
	writeTestFile(t, workflowPath, `on: push
jobs:
  terraform:
    runs-on: ubuntu-latest
    steps:
      - run: terraform fmt -check -recursive
      - run: terraform validate
`)

	fmtFinding = CheckTerraformFmtRecommended(tmpDir)
	if fmtFinding.Status != StatusPass {
		t.Errorf("expected StatusPass when terraform fmt is automated, got %v", fmtFinding.Status)
	}

	validateFinding = CheckTerraformValidateRecommended(tmpDir)
	if validateFinding.Status != StatusPass {
		t.Errorf("expected StatusPass when terraform validate is automated, got %v", validateFinding.Status)
	}
}

func TestCheckTerraformRequiredProvidersExists(t *testing.T) {
	tmpDir := t.TempDir()
	writeTestFile(t, filepath.Join(tmpDir, "main.tf"), `resource "null_resource" "example" {}`)

	finding := CheckTerraformRequiredProvidersExists(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for missing required_providers, got %v", finding.Status)
	}

	writeTestFile(t, filepath.Join(tmpDir, "versions.tf"), `terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}`)

	finding = CheckTerraformRequiredProvidersExists(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for required_providers block, got %v", finding.Status)
	}
}

func TestCheckTerraformProviderVersionsConstrained(t *testing.T) {
	tmpDir := t.TempDir()
	writeTestFile(t, filepath.Join(tmpDir, "versions.tf"), `terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
}`)

	finding := CheckTerraformProviderVersionsConstrained(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for missing provider version, got %v", finding.Status)
	}

	writeTestFile(t, filepath.Join(tmpDir, "versions.tf"), `terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}`)

	finding = CheckTerraformProviderVersionsConstrained(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for constrained provider version, got %v", finding.Status)
	}
}

func TestCheckTerraformBackendConfigured(t *testing.T) {
	tmpDir := t.TempDir()
	writeTestFile(t, filepath.Join(tmpDir, "main.tf"), `resource "null_resource" "example" {}`)

	finding := CheckTerraformBackendConfigured(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for missing backend block, got %v", finding.Status)
	}

	writeTestFile(t, filepath.Join(tmpDir, "backend.tf"), `terraform {
  backend "s3" {
    bucket = "example-state"
    key    = "terraform.tfstate"
    region = "eu-west-2"
  }
}`)

	finding = CheckTerraformBackendConfigured(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for backend block, got %v", finding.Status)
	}
}

func TestCheckTerraformNoSuspiciousVariableDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	writeTestFile(t, filepath.Join(tmpDir, "variables.tf"), `variable "app_name" {
  default = "super-secret-value"
}`)

	finding := CheckTerraformNoSuspiciousVariableDefaults(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for suspicious default value, got %v", finding.Status)
	}

	writeTestFile(t, filepath.Join(tmpDir, "variables.tf"), `variable "db_password" {
  default = "change-me"
}`)

	finding = CheckTerraformNoSuspiciousVariableDefaults(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for suspicious variable name, got %v", finding.Status)
	}

	writeTestFile(t, filepath.Join(tmpDir, "variables.tf"), `variable "environment" {
  default = "prod"
}`)

	finding = CheckTerraformNoSuspiciousVariableDefaults(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for safe variable default, got %v", finding.Status)
	}
}

func TestCheckTerraformLockfilePresent(t *testing.T) {
	tmpDir := t.TempDir()
	writeTestFile(t, filepath.Join(tmpDir, "main.tf"), `resource "null_resource" "example" {}`)

	finding := CheckTerraformLockfilePresent(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for missing lockfile, got %v", finding.Status)
	}

	writeTestFile(t, filepath.Join(tmpDir, ".terraform.lock.hcl"), `provider "registry.terraform.io/hashicorp/aws" {}`)

	finding = CheckTerraformLockfilePresent(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for lockfile, got %v", finding.Status)
	}
}

func TestTerraformNoFilesSkipDeeperChecksCleanly(t *testing.T) {
	tmpDir := t.TempDir()

	checks := []Finding{
		CheckTerraformFmtRecommended(tmpDir),
		CheckTerraformValidateRecommended(tmpDir),
		CheckTerraformRequiredProvidersExists(tmpDir),
		CheckTerraformProviderVersionsConstrained(tmpDir),
		CheckTerraformBackendConfigured(tmpDir),
		CheckTerraformNoSuspiciousVariableDefaults(tmpDir),
		CheckTerraformLockfilePresent(tmpDir),
	}

	for _, finding := range checks {
		if finding.Status != StatusPass {
			t.Errorf("expected %s to skip with StatusPass, got %v", finding.ID, finding.Status)
		}
	}
}

func writeTestFile(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("failed to create parent directory: %v", err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
}
