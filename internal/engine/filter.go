package engine

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Sufyan11104/shipcheck/internal/rules"
)

var allowedCategories = map[string]bool{
	"ci":        true,
	"docker":    true,
	"docs":      true,
	"env":       true,
	"k8s":       true,
	"repo":      true,
	"terraform": true,
}

// ParseCategoryFilter parses a comma-separated category list.
func ParseCategoryFilter(value string) ([]string, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}

	seen := make(map[string]bool)
	var categories []string
	for _, part := range strings.Split(value, ",") {
		category := strings.ToLower(strings.TrimSpace(part))
		if category == "" {
			continue
		}

		if !allowedCategories[category] {
			return nil, fmt.Errorf("unknown category %q (valid: %s)", category, strings.Join(AllowedCategories(), ", "))
		}

		if !seen[category] {
			categories = append(categories, category)
			seen[category] = true
		}
	}

	return categories, nil
}

// FilterFindingsByCategory returns findings matching selected categories.
func FilterFindingsByCategory(findings []rules.Finding, categories []string) []rules.Finding {
	if len(categories) == 0 {
		return findings
	}

	selected := make(map[string]bool, len(categories))
	for _, category := range categories {
		selected[category] = true
	}

	filtered := make([]rules.Finding, 0, len(findings))
	for _, finding := range findings {
		if selected[strings.ToLower(finding.Category)] {
			filtered = append(filtered, finding)
		}
	}

	return filtered
}

// AllowedCategories returns the supported audit category names.
func AllowedCategories() []string {
	categories := make([]string, 0, len(allowedCategories))
	for category := range allowedCategories {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	return categories
}
