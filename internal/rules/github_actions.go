package rules

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// CheckWorkflowsDirExists checks if .github/workflows directory exists
func CheckWorkflowsDirExists(path string) Finding {
	workflowsPath := filepath.Join(path, ".github", "workflows")
	if _, err := os.Stat(workflowsPath); err == nil {
		return Finding{
			ID:          "ci.workflows_dir_exists",
			Title:       ".github/workflows found",
			Category:    "ci",
			Severity:    SeverityMedium,
			Status:      StatusPass,
			Message:     ".github/workflows directory found",
			Remediation: "N/A",
		}
	}

	return Finding{
		ID:          "ci.workflows_dir_exists",
		Title:       ".github/workflows missing",
		Category:    "ci",
		Severity:    SeverityMedium,
		Status:      StatusWarn,
		Message:     "No .github/workflows directory found",
		Remediation: "Create a .github/workflows directory and add CI/CD workflow files",
	}
}

// CheckWorkflowFileExists checks if at least one workflow YAML file exists
func CheckWorkflowFileExists(path string) Finding {
	workflowsPath := filepath.Join(path, ".github", "workflows")
	entries, err := os.ReadDir(workflowsPath)
	if err != nil {
		// Directory doesn't exist, skip gracefully
		return Finding{
			ID:          "ci.workflow_file_exists",
			Title:       ".github/workflows not found",
			Category:    "ci",
			Severity:    SeverityLow,
			Status:      StatusPass,
			Message:     ".github/workflows directory not present; skipping workflow file check",
			Remediation: "N/A",
		}
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			name := strings.ToLower(entry.Name())
			if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
				return Finding{
					ID:          "ci.workflow_file_exists",
					Title:       "Workflow file found",
					Category:    "ci",
					Severity:    SeverityMedium,
					Status:      StatusPass,
					Message:     "At least one workflow YAML file found",
					Remediation: "N/A",
				}
			}
		}
	}

	return Finding{
		ID:          "ci.workflow_file_exists",
		Title:       "No workflow files found",
		Category:    "ci",
		Severity:    SeverityMedium,
		Status:      StatusWarn,
		Message:     "No .yml or .yaml workflow files found in .github/workflows",
		Remediation: "Add at least one workflow YAML file for CI/CD",
	}
}

// CheckTestStepExists checks if workflows contain test commands
func CheckTestStepExists(path string) Finding {
	testPatterns := []string{
		"go test",
		"npm test",
		"pnpm test",
		"yarn test",
		"pytest",
		"cargo test",
		"mvn test",
		"gradle test",
	}

	if foundInWorkflows(path, testPatterns) {
		return Finding{
			ID:          "ci.test_step_exists",
			Title:       "Test step found",
			Category:    "ci",
			Severity:    SeverityHigh,
			Status:      StatusPass,
			Message:     "Workflow files contain test step",
			Remediation: "N/A",
		}
	}

	return Finding{
		ID:          "ci.test_step_exists",
		Title:       "No test step detected",
		Category:    "ci",
		Severity:    SeverityHigh,
		Status:      StatusWarn,
		Message:     "No test commands detected in workflow files",
		Remediation: "Add a test step to your CI workflow (e.g., go test, npm test)",
	}
}

// CheckBuildStepExists checks if workflows contain build commands
func CheckBuildStepExists(path string) Finding {
	buildPatterns := []string{
		"go build",
		"npm run build",
		"pnpm build",
		"yarn build",
		"docker build",
		"cargo build",
		"mvn package",
		"gradle build",
	}

	if foundInWorkflows(path, buildPatterns) {
		return Finding{
			ID:          "ci.build_step_exists",
			Title:       "Build step found",
			Category:    "ci",
			Severity:    SeverityHigh,
			Status:      StatusPass,
			Message:     "Workflow files contain build step",
			Remediation: "N/A",
		}
	}

	return Finding{
		ID:          "ci.build_step_exists",
		Title:       "No build step detected",
		Category:    "ci",
		Severity:    SeverityMedium,
		Status:      StatusWarn,
		Message:     "No build commands detected in workflow files",
		Remediation: "Add a build step to your CI workflow (e.g., go build, npm run build)",
	}
}

// CheckDeployAfterTests checks if deploy commands appear before test commands
func CheckDeployAfterTests(path string) Finding {
	deployPatterns := []string{"deploy", "release", "publish"}
	testPatterns := []string{"go test", "npm test", "pnpm test", "yarn test", "pytest", "cargo test", "mvn test", "gradle test"}

	workflowsPath := filepath.Join(path, ".github", "workflows")
	entries, err := os.ReadDir(workflowsPath)
	if err != nil {
		return Finding{
			ID:          "ci.deploy_after_tests",
			Title:       "Workflows not found",
			Category:    "ci",
			Severity:    SeverityLow,
			Status:      StatusPass,
			Message:     ".github/workflows directory not present; skipping deploy order check",
			Remediation: "N/A",
		}
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			name := strings.ToLower(entry.Name())
			if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
				filePath := filepath.Join(workflowsPath, entry.Name())
				file, err := os.Open(filePath)
				if err != nil {
					continue
				}
				defer file.Close()

				content := ""
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					content += strings.ToLower(scanner.Text()) + "\n"
				}

				// Find first occurrence of deploy and test patterns
				firstDeploy := len(content)
				firstTest := len(content)

				for _, pattern := range deployPatterns {
					idx := strings.Index(content, pattern)
					if idx != -1 && idx < firstDeploy {
						firstDeploy = idx
					}
				}

				for _, pattern := range testPatterns {
					idx := strings.Index(content, pattern)
					if idx != -1 && idx < firstTest {
						firstTest = idx
					}
				}

				// Warn if deploy appears before test
				if firstDeploy < firstTest && firstDeploy < len(content) {
					return Finding{
						ID:          "ci.deploy_after_tests",
						Title:       "Deploy may run before tests",
						Category:    "ci",
						Severity:    SeverityHigh,
						Status:      StatusWarn,
						Message:     "Workflow file may have deployment before testing",
						Remediation: "Ensure tests run before deployment steps",
					}
				}
			}
		}
	}

	return Finding{
		ID:          "ci.deploy_after_tests",
		Title:       "Deploy after tests verified",
		Category:    "ci",
		Severity:    SeverityHigh,
		Status:      StatusPass,
		Message:     "No obvious deploy-before-test pattern detected",
		Remediation: "N/A",
	}
}

// CheckActionsPinned checks that GitHub Actions use version pinning
func CheckActionsPinned(path string) Finding {
	workflowsPath := filepath.Join(path, ".github", "workflows")
	entries, err := os.ReadDir(workflowsPath)
	if err != nil {
		return Finding{
			ID:          "ci.actions_pinned",
			Title:       "Workflows not found",
			Category:    "ci",
			Severity:    SeverityLow,
			Status:      StatusPass,
			Message:     ".github/workflows directory not present; skipping action pinning check",
			Remediation: "N/A",
		}
	}

	usesRegex := regexp.MustCompile(`uses:\s*([a-zA-Z0-9\-._/]+)(?:@|#|$)`)
	unpinnedRegex := regexp.MustCompile(`uses:\s*[a-zA-Z0-9\-._/]+\s*$`)

	foundUnpinned := false
	foundPinned := false

	for _, entry := range entries {
		if !entry.IsDir() {
			name := strings.ToLower(entry.Name())
			if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
				filePath := filepath.Join(workflowsPath, entry.Name())
				file, err := os.Open(filePath)
				if err != nil {
					continue
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := scanner.Text()
					if strings.Contains(line, "uses:") {
						if usesRegex.MatchString(line) && strings.Contains(line, "@") {
							foundPinned = true
						} else if unpinnedRegex.MatchString(line) {
							foundUnpinned = true
						}
					}
				}
			}
		}
	}

	if foundUnpinned {
		return Finding{
			ID:          "ci.actions_pinned",
			Title:       "Unpinned action detected",
			Category:    "ci",
			Severity:    SeverityMedium,
			Status:      StatusWarn,
			Message:     "Workflow actions are not pinned to a version",
			Remediation: "Pin action versions using @v3, @v4, etc. Use full SHA pinning for higher security",
		}
	}

	if foundPinned {
		return Finding{
			ID:          "ci.actions_pinned",
			Title:       "Actions pinned to versions",
			Category:    "ci",
			Severity:    SeverityMedium,
			Status:      StatusPass,
			Message:     "Workflow actions are pinned to versions",
			Remediation: "N/A",
		}
	}

	return Finding{
		ID:          "ci.actions_pinned",
		Title:       "No actions found",
		Category:    "ci",
		Severity:    SeverityLow,
		Status:      StatusPass,
		Message:     "No GitHub Actions used in workflows or workflows not found",
		Remediation: "N/A",
	}
}

// CheckNoSecretEcho checks for commands that echo secrets
func CheckNoSecretEcho(path string) Finding {
	secretEchoPatterns := []string{
		"echo ${{ secrets.",
		"echo \"${{ secrets.",
		"printf ${{ secrets.",
		"echo '${{ secrets.",
	}

	workflowsPath := filepath.Join(path, ".github", "workflows")
	entries, err := os.ReadDir(workflowsPath)
	if err != nil {
		return Finding{
			ID:          "ci.no_secret_echo",
			Title:       "Workflows not found",
			Category:    "ci",
			Severity:    SeverityLow,
			Status:      StatusPass,
			Message:     ".github/workflows directory not present; skipping secret echo check",
			Remediation: "N/A",
		}
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			name := strings.ToLower(entry.Name())
			if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
				filePath := filepath.Join(workflowsPath, entry.Name())
				file, err := os.Open(filePath)
				if err != nil {
					continue
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := strings.ToLower(scanner.Text())
					for _, pattern := range secretEchoPatterns {
						if strings.Contains(line, strings.ToLower(pattern)) {
							return Finding{
								ID:          "ci.no_secret_echo",
								Title:       "Secret echo detected",
								Category:    "ci",
								Severity:    SeverityHigh,
								Status:      StatusWarn,
								Message:     "Workflow may be echoing secrets to logs",
								Remediation: "Avoid echoing secrets; use masked output or environment variables",
							}
						}
					}
				}
			}
		}
	}

	return Finding{
		ID:          "ci.no_secret_echo",
		Title:       "No secret echo detected",
		Category:    "ci",
		Severity:    SeverityHigh,
		Status:      StatusPass,
		Message:     "Workflows do not appear to echo secrets",
		Remediation: "N/A",
	}
}

// CheckPermissionsDeclared checks for explicit permissions blocks
func CheckPermissionsDeclared(path string) Finding {
	workflowsPath := filepath.Join(path, ".github", "workflows")
	entries, err := os.ReadDir(workflowsPath)
	if err != nil {
		return Finding{
			ID:          "ci.permissions_declared",
			Title:       "Workflows not found",
			Category:    "ci",
			Severity:    SeverityLow,
			Status:      StatusPass,
			Message:     ".github/workflows directory not present; skipping permissions check",
			Remediation: "N/A",
		}
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			name := strings.ToLower(entry.Name())
			if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
				filePath := filepath.Join(workflowsPath, entry.Name())
				file, err := os.Open(filePath)
				if err != nil {
					continue
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := scanner.Text()
					trimmed := strings.TrimSpace(line)
					if trimmed == "permissions:" || strings.HasPrefix(trimmed, "permissions:") {
						return Finding{
							ID:          "ci.permissions_declared",
							Title:       "Permissions declared",
							Category:    "ci",
							Severity:    SeverityMedium,
							Status:      StatusPass,
							Message:     "Workflow includes explicit permissions block",
							Remediation: "N/A",
						}
					}
				}
			}
		}
	}

	return Finding{
		ID:          "ci.permissions_declared",
		Title:       "No permissions block",
		Category:    "ci",
		Severity:    SeverityMedium,
		Status:      StatusWarn,
		Message:     "No explicit permissions block detected in workflows",
		Remediation: "Add a top-level permissions block to limit workflow access",
	}
}

// foundInWorkflows is a helper function to search for patterns in workflow files
func foundInWorkflows(path string, patterns []string) bool {
	workflowsPath := filepath.Join(path, ".github", "workflows")
	entries, err := os.ReadDir(workflowsPath)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			name := strings.ToLower(entry.Name())
			if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
				filePath := filepath.Join(workflowsPath, entry.Name())
				file, err := os.Open(filePath)
				if err != nil {
					continue
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := strings.ToLower(scanner.Text())
					for _, pattern := range patterns {
						if strings.Contains(line, strings.ToLower(pattern)) {
							return true
						}
					}
				}
			}
		}
	}

	return false
}
