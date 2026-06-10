package engine

import (
	"strings"
	"testing"

	"github.com/Sufyan11104/shipcheck/internal/rules"
)

func TestFilterFindingsByCategory_OneCategory(t *testing.T) {
	findings := sampleCategoryFindings()

	categories, err := ParseCategoryFilter("docker")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	filtered := FilterFindingsByCategory(findings, categories)
	if len(filtered) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(filtered))
	}
	if filtered[0].Category != "docker" {
		t.Errorf("expected docker finding, got %s", filtered[0].Category)
	}
}

func TestFilterFindingsByCategory_MultipleCategories(t *testing.T) {
	findings := sampleCategoryFindings()

	categories, err := ParseCategoryFilter("docker,ci")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	filtered := FilterFindingsByCategory(findings, categories)
	if len(filtered) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(filtered))
	}
	if filtered[0].Category != "docker" || filtered[1].Category != "ci" {
		t.Errorf("expected docker and ci findings, got %+v", filtered)
	}
}

func TestParseCategoryFilter_UnknownCategory(t *testing.T) {
	_, err := ParseCategoryFilter("docker,unknown")
	if err == nil {
		t.Fatal("expected error for unknown category")
	}
	if !strings.Contains(err.Error(), "unknown category") {
		t.Errorf("expected clear unknown category error, got %v", err)
	}
}

func sampleCategoryFindings() []rules.Finding {
	return []rules.Finding{
		{ID: "docker.dockerfile_exists", Category: "docker", Severity: rules.SeverityMedium, Status: rules.StatusPass},
		{ID: "ci.workflow_file_exists", Category: "ci", Severity: rules.SeverityMedium, Status: rules.StatusWarn},
		{ID: "k8s.manifest_exists", Category: "k8s", Severity: rules.SeverityMedium, Status: rules.StatusWarn},
	}
}
