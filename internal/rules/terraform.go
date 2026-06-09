package rules

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type terraformBlock struct {
	name    string
	content string
}

var terraformSecretKeywords = []string{
	"password",
	"secret",
	"token",
	"api_key",
	"access_key",
	"private_key",
	"credential",
}

// CheckTerraformFilesExist checks if Terraform or tfvars files exist.
func CheckTerraformFilesExist(path string) Finding {
	tfFiles, tfVarsFiles := findTerraformFiles(path)
	if len(tfFiles)+len(tfVarsFiles) > 0 {
		return Finding{
			ID:          "terraform.files_exist",
			Title:       "Terraform files found",
			Category:    "terraform",
			Severity:    SeverityMedium,
			Status:      StatusPass,
			Message:     "Terraform files detected",
			Remediation: "N/A",
		}
	}

	return Finding{
		ID:          "terraform.files_exist",
		Title:       "No Terraform files",
		Category:    "terraform",
		Severity:    SeverityMedium,
		Status:      StatusWarn,
		Message:     "No Terraform files found",
		Remediation: "Add Terraform configuration if infrastructure-as-code is used for this project",
	}
}

// CheckTerraformFmtRecommended checks whether terraform fmt appears in CI or local automation.
func CheckTerraformFmtRecommended(path string) Finding {
	if !hasTerraformConfigFiles(path) {
		return Finding{
			ID:          "terraform.fmt_recommended",
			Title:       "No Terraform configuration",
			Category:    "terraform",
			Severity:    SeverityLow,
			Status:      StatusPass,
			Message:     "No .tf files present; skipping terraform fmt recommendation",
			Remediation: "N/A",
		}
	}

	if terraformAutomationCommandFound(path, "terraform fmt") {
		return Finding{
			ID:          "terraform.fmt_recommended",
			Title:       "terraform fmt found",
			Category:    "terraform",
			Severity:    SeverityInfo,
			Status:      StatusPass,
			Message:     "terraform fmt appears in CI or local automation",
			Remediation: "N/A",
		}
	}

	return Finding{
		ID:          "terraform.fmt_recommended",
		Title:       "terraform fmt recommended",
		Category:    "terraform",
		Severity:    SeverityInfo,
		Status:      StatusWarn,
		Message:     "Terraform files detected; run terraform fmt in CI or local workflow",
		Remediation: "Add terraform fmt -check -recursive to CI or local automation",
	}
}

// CheckTerraformValidateRecommended checks whether terraform validate appears in CI or local automation.
func CheckTerraformValidateRecommended(path string) Finding {
	if !hasTerraformConfigFiles(path) {
		return Finding{
			ID:          "terraform.validate_recommended",
			Title:       "No Terraform configuration",
			Category:    "terraform",
			Severity:    SeverityLow,
			Status:      StatusPass,
			Message:     "No .tf files present; skipping terraform validate recommendation",
			Remediation: "N/A",
		}
	}

	if terraformAutomationCommandFound(path, "terraform validate") {
		return Finding{
			ID:          "terraform.validate_recommended",
			Title:       "terraform validate found",
			Category:    "terraform",
			Severity:    SeverityInfo,
			Status:      StatusPass,
			Message:     "terraform validate appears in CI or local automation",
			Remediation: "N/A",
		}
	}

	return Finding{
		ID:          "terraform.validate_recommended",
		Title:       "terraform validate recommended",
		Category:    "terraform",
		Severity:    SeverityInfo,
		Status:      StatusWarn,
		Message:     "Terraform files detected; run terraform validate in CI or local workflow",
		Remediation: "Add terraform validate to CI or local automation after terraform init",
	}
}

// CheckTerraformRequiredProvidersExists checks if Terraform declares required providers.
func CheckTerraformRequiredProvidersExists(path string) Finding {
	content := readTerraformConfigFiles(path)
	if content == "" {
		return Finding{
			ID:          "terraform.required_providers_exists",
			Title:       "No Terraform configuration",
			Category:    "terraform",
			Severity:    SeverityLow,
			Status:      StatusPass,
			Message:     "No .tf files present; skipping required_providers check",
			Remediation: "N/A",
		}
	}

	if len(extractTerraformBlocks(content, "required_providers")) > 0 {
		return Finding{
			ID:          "terraform.required_providers_exists",
			Title:       "required_providers found",
			Category:    "terraform",
			Severity:    SeverityMedium,
			Status:      StatusPass,
			Message:     "Terraform required_providers block detected",
			Remediation: "N/A",
		}
	}

	return Finding{
		ID:          "terraform.required_providers_exists",
		Title:       "No required_providers block",
		Category:    "terraform",
		Severity:    SeverityMedium,
		Status:      StatusWarn,
		Message:     "No required_providers block detected in Terraform files",
		Remediation: "Declare providers in a terraform required_providers block with source and version constraints",
	}
}

// CheckTerraformProviderVersionsConstrained checks if required provider declarations include versions.
func CheckTerraformProviderVersionsConstrained(path string) Finding {
	content := readTerraformConfigFiles(path)
	if content == "" {
		return Finding{
			ID:          "terraform.provider_versions_constrained",
			Title:       "No Terraform configuration",
			Category:    "terraform",
			Severity:    SeverityLow,
			Status:      StatusPass,
			Message:     "No .tf files present; skipping provider version check",
			Remediation: "N/A",
		}
	}

	blocks := extractTerraformBlocks(content, "required_providers")
	if len(blocks) == 0 {
		return Finding{
			ID:          "terraform.provider_versions_constrained",
			Title:       "No provider constraints",
			Category:    "terraform",
			Severity:    SeverityMedium,
			Status:      StatusWarn,
			Message:     "No required_providers block found to verify provider version constraints",
			Remediation: "Add provider source and version constraints in required_providers",
		}
	}

	if requiredProviderBlocksHaveVersionConstraints(blocks) {
		return Finding{
			ID:          "terraform.provider_versions_constrained",
			Title:       "Provider versions constrained",
			Category:    "terraform",
			Severity:    SeverityMedium,
			Status:      StatusPass,
			Message:     "Terraform provider declarations appear to include version constraints",
			Remediation: "N/A",
		}
	}

	return Finding{
		ID:          "terraform.provider_versions_constrained",
		Title:       "Provider versions not constrained",
		Category:    "terraform",
		Severity:    SeverityMedium,
		Status:      StatusWarn,
		Message:     "Terraform provider declarations do not appear to include version constraints",
		Remediation: "Add version constraints to each provider in required_providers",
	}
}

// CheckTerraformBackendConfigured checks if Terraform config declares a backend block.
func CheckTerraformBackendConfigured(path string) Finding {
	content := readTerraformConfigFiles(path)
	if content == "" {
		return Finding{
			ID:          "terraform.backend_configured",
			Title:       "No Terraform configuration",
			Category:    "terraform",
			Severity:    SeverityLow,
			Status:      StatusPass,
			Message:     "No .tf files present; skipping backend check",
			Remediation: "N/A",
		}
	}

	backendPattern := regexp.MustCompile(`(?i)\bbackend\s+"[^"]+"\s*{`)
	if backendPattern.MatchString(content) {
		return Finding{
			ID:          "terraform.backend_configured",
			Title:       "Backend block found",
			Category:    "terraform",
			Severity:    SeverityLow,
			Status:      StatusPass,
			Message:     "Terraform backend block detected",
			Remediation: "N/A",
		}
	}

	return Finding{
		ID:          "terraform.backend_configured",
		Title:       "No backend block",
		Category:    "terraform",
		Severity:    SeverityLow,
		Status:      StatusWarn,
		Message:     "No backend block detected; remote state may be needed for team or production use.",
		Remediation: "Consider configuring a remote backend for shared or production Terraform state",
	}
}

// CheckTerraformNoSuspiciousVariableDefaults checks for secret-like variable defaults.
func CheckTerraformNoSuspiciousVariableDefaults(path string) Finding {
	content := readTerraformConfigFiles(path)
	if content == "" {
		return Finding{
			ID:          "terraform.no_suspicious_variable_defaults",
			Title:       "No Terraform configuration",
			Category:    "terraform",
			Severity:    SeverityLow,
			Status:      StatusPass,
			Message:     "No .tf files present; skipping variable default check",
			Remediation: "N/A",
		}
	}

	if hasSuspiciousTerraformVariableDefault(content) {
		return Finding{
			ID:          "terraform.no_suspicious_variable_defaults",
			Title:       "Suspicious variable default",
			Category:    "terraform",
			Severity:    SeverityHigh,
			Status:      StatusWarn,
			Message:     "Terraform variable defaults may contain secret-like names or values",
			Remediation: "Avoid committing secret defaults; use secret managers, environment variables, or CI/CD secret storage",
		}
	}

	return Finding{
		ID:          "terraform.no_suspicious_variable_defaults",
		Title:       "No suspicious variable defaults",
		Category:    "terraform",
		Severity:    SeverityHigh,
		Status:      StatusPass,
		Message:     "No obvious secret-like Terraform variable defaults detected",
		Remediation: "N/A",
	}
}

// CheckTerraformLockfilePresent checks if the Terraform dependency lock file is present.
func CheckTerraformLockfilePresent(path string) Finding {
	if !hasTerraformConfigFiles(path) {
		return Finding{
			ID:          "terraform.lockfile_present",
			Title:       "No Terraform configuration",
			Category:    "terraform",
			Severity:    SeverityLow,
			Status:      StatusPass,
			Message:     "No .tf files present; skipping Terraform lockfile check",
			Remediation: "N/A",
		}
	}

	if terraformLockfileExists(path) {
		return Finding{
			ID:          "terraform.lockfile_present",
			Title:       "Terraform lockfile found",
			Category:    "terraform",
			Severity:    SeverityLow,
			Status:      StatusPass,
			Message:     ".terraform.lock.hcl file detected",
			Remediation: "N/A",
		}
	}

	return Finding{
		ID:          "terraform.lockfile_present",
		Title:       "Terraform lockfile missing",
		Category:    "terraform",
		Severity:    SeverityLow,
		Status:      StatusWarn,
		Message:     ".terraform.lock.hcl not found; it may not exist before terraform init",
		Remediation: "Commit .terraform.lock.hcl after terraform init for reproducible provider selections",
	}
}

func findTerraformFiles(path string) (tfFiles []string, tfVarsFiles []string) {
	filepath.WalkDir(path, func(filePath string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if entry.IsDir() {
			switch entry.Name() {
			case ".git", ".terraform":
				return filepath.SkipDir
			}
			return nil
		}

		name := strings.ToLower(entry.Name())
		switch {
		case strings.HasSuffix(name, ".tf"):
			tfFiles = append(tfFiles, filePath)
		case strings.HasSuffix(name, ".tfvars"):
			tfVarsFiles = append(tfVarsFiles, filePath)
		}

		return nil
	})

	return tfFiles, tfVarsFiles
}

func hasTerraformConfigFiles(path string) bool {
	tfFiles, _ := findTerraformFiles(path)
	return len(tfFiles) > 0
}

func readTerraformConfigFiles(path string) string {
	tfFiles, _ := findTerraformFiles(path)
	if len(tfFiles) == 0 {
		return ""
	}

	var content strings.Builder
	for _, filePath := range tfFiles {
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}
		content.WriteString("\n")
		content.WriteString(stripTerraformComments(string(data)))
	}

	return content.String()
}

func stripTerraformComments(content string) string {
	blockCommentPattern := regexp.MustCompile(`(?s)/\*.*?\*/`)
	content = blockCommentPattern.ReplaceAllString(content, "")

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
			lines[i] = ""
		}
	}

	return strings.Join(lines, "\n")
}

func terraformAutomationCommandFound(path, command string) bool {
	candidates := []string{
		filepath.Join(path, ".github", "workflows"),
		filepath.Join(path, "Makefile"),
		filepath.Join(path, "makefile"),
		filepath.Join(path, "Taskfile.yml"),
		filepath.Join(path, "Taskfile.yaml"),
		filepath.Join(path, "justfile"),
	}

	for _, candidate := range candidates {
		info, err := os.Stat(candidate)
		if err != nil {
			continue
		}

		if info.IsDir() {
			if terraformCommandFoundInDir(candidate, command) {
				return true
			}
			continue
		}

		if terraformCommandFoundInFile(candidate, command) {
			return true
		}
	}

	return false
}

func terraformCommandFoundInDir(path, command string) bool {
	found := false
	filepath.WalkDir(path, func(filePath string, entry os.DirEntry, err error) error {
		if err != nil || found {
			return nil
		}

		if entry.IsDir() {
			return nil
		}

		name := strings.ToLower(entry.Name())
		if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
			found = terraformCommandFoundInFile(filePath, command)
		}

		return nil
	})

	return found
}

func terraformCommandFoundInFile(path, command string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	return strings.Contains(strings.ToLower(string(data)), strings.ToLower(command))
}

func extractTerraformBlocks(content, keyword string) []terraformBlock {
	pattern := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(keyword) + `\b\s*(?:"([^"]+)")?\s*(?:=\s*)?{`)
	matches := pattern.FindAllStringSubmatchIndex(content, -1)
	blocks := make([]terraformBlock, 0, len(matches))

	for _, match := range matches {
		name := ""
		if len(match) >= 4 && match[2] >= 0 {
			name = content[match[2]:match[3]]
		}

		openBrace := strings.LastIndex(content[match[0]:match[1]], "{")
		if openBrace < 0 {
			continue
		}
		openBrace += match[0]

		closeBrace := findMatchingBrace(content, openBrace)
		if closeBrace <= openBrace {
			continue
		}

		blocks = append(blocks, terraformBlock{
			name:    name,
			content: content[openBrace+1 : closeBrace],
		})
	}

	return blocks
}

func findMatchingBrace(content string, openBrace int) int {
	depth := 0
	for i := openBrace; i < len(content); i++ {
		switch content[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}

	return -1
}

func requiredProviderBlocksHaveVersionConstraints(blocks []terraformBlock) bool {
	sourcePattern := regexp.MustCompile(`(?i)\bsource\s*=`)
	versionPattern := regexp.MustCompile(`(?i)\bversion\s*=\s*["'][^"']*[0-9~><=!][^"']*["']`)

	for _, block := range blocks {
		sourceCount := len(sourcePattern.FindAllStringIndex(block.content, -1))
		versionCount := len(versionPattern.FindAllStringIndex(block.content, -1))
		if versionCount == 0 {
			return false
		}
		if sourceCount > 0 && versionCount < sourceCount {
			return false
		}
	}

	return true
}

func hasSuspiciousTerraformVariableDefault(content string) bool {
	for _, block := range extractTerraformBlocks(content, "variable") {
		defaultExpression, hasDefault := terraformVariableDefaultExpression(block.content)
		if !hasDefault || isEmptyTerraformDefault(defaultExpression) {
			continue
		}

		if containsTerraformSecretKeyword(block.name) || containsTerraformSecretKeyword(defaultExpression) {
			return true
		}
	}

	return false
}

func terraformVariableDefaultExpression(blockContent string) (string, bool) {
	defaultPattern := regexp.MustCompile(`(?i)\bdefault\s*=\s*(.*)`)
	lines := strings.Split(blockContent, "\n")

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
			continue
		}

		matches := defaultPattern.FindStringSubmatch(line)
		if len(matches) == 0 {
			continue
		}

		expression := strings.TrimSpace(matches[1])
		if expression == "{" || strings.HasPrefix(expression, "{") || expression == "[" || strings.HasPrefix(expression, "[") {
			return strings.Join(lines[i:], "\n"), true
		}

		return expression, true
	}

	return "", false
}

func isEmptyTerraformDefault(expression string) bool {
	normalized := strings.TrimSpace(expression)
	normalized = strings.Trim(normalized, `"'`)

	switch strings.ToLower(normalized) {
	case "", "null", "[]", "{}":
		return true
	default:
		return false
	}
}

func containsTerraformSecretKeyword(value string) bool {
	lower := strings.ToLower(value)
	for _, keyword := range terraformSecretKeywords {
		if strings.Contains(lower, keyword) {
			return true
		}
	}

	return false
}

func terraformLockfileExists(path string) bool {
	found := false
	filepath.WalkDir(path, func(filePath string, entry os.DirEntry, err error) error {
		if err != nil || found {
			return nil
		}

		if entry.IsDir() {
			if entry.Name() == ".git" || entry.Name() == ".terraform" {
				return filepath.SkipDir
			}
			return nil
		}

		if entry.Name() == ".terraform.lock.hcl" {
			found = true
		}

		return nil
	})

	return found
}
