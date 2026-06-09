package rules

import (
	"os"
	"path/filepath"
)

// CheckReadmeExists checks if README exists in the directory
func CheckReadmeExists(path string) Finding {
	readmeFiles := []string{"README.md", "README.txt", "README"}

	for _, name := range readmeFiles {
		filePath := filepath.Join(path, name)
		if _, err := os.Stat(filePath); err == nil {
			return Finding{
				ID:          "docs.readme_exists",
				Title:       "README found",
				Category:    "docs",
				Severity:    SeverityHigh,
				Status:      StatusPass,
				Message:     "README file found",
				Remediation: "N/A",
			}
		}
	}

	return Finding{
		ID:          "docs.readme_exists",
		Title:       "README missing",
		Category:    "docs",
		Severity:    SeverityHigh,
		Status:      StatusFail,
		Message:     "No README file found",
		Remediation: "Create a README.md file to document the project",
	}
}

// CheckGitignoreExists checks if .gitignore exists in the directory
func CheckGitignoreExists(path string) Finding {
	filePath := filepath.Join(path, ".gitignore")
	if _, err := os.Stat(filePath); err == nil {
		return Finding{
			ID:          "repo.gitignore_exists",
			Title:       ".gitignore found",
			Category:    "repo",
			Severity:    SeverityMedium,
			Status:      StatusPass,
			Message:     ".gitignore file found",
			Remediation: "N/A",
		}
	}

	return Finding{
		ID:          "repo.gitignore_exists",
		Title:       ".gitignore missing",
		Category:    "repo",
		Severity:    SeverityMedium,
		Status:      StatusFail,
		Message:     "No .gitignore file found",
		Remediation: "Create a .gitignore file to exclude sensitive files from version control",
	}
}

// CheckEnvNotCommitted checks if .env file is in the repository
func CheckEnvNotCommitted(path string, isGitRepo bool) Finding {
	if !isGitRepo {
		// Can't check if not a git repo, treat as warning
		return Finding{
			ID:          "env.env_not_committed",
			Title:       "Git repo required",
			Category:    "env",
			Severity:    SeverityLow,
			Status:      StatusWarn,
			Message:     "Not a Git repository; cannot verify .env status",
			Remediation: "Initialize a Git repository for proper version control",
		}
	}

	// Check if .env file exists in working directory (not committed, just present)
	envPath := filepath.Join(path, ".env")
	if _, err := os.Stat(envPath); err == nil {
		return Finding{
			ID:          "env.env_not_committed",
			Title:       ".env file detected",
			Category:    "env",
			Severity:    SeverityHigh,
			Status:      StatusWarn,
			Message:     ".env file found - ensure it is in .gitignore to prevent accidental commits",
			Remediation: "Add .env to .gitignore and use .env.example for template",
		}
	}

	// .env doesn't exist, which is good
	return Finding{
		ID:          "env.env_not_committed",
		Title:       "No committed .env file",
		Category:    "env",
		Severity:    SeverityHigh,
		Status:      StatusPass,
		Message:     "No committed .env file detected",
		Remediation: "N/A",
	}
}

// CheckEnvExampleExists checks if .env.example exists
func CheckEnvExampleExists(path string) Finding {
	filePath := filepath.Join(path, ".env.example")
	if _, err := os.Stat(filePath); err == nil {
		return Finding{
			ID:          "env.env_example_exists",
			Title:       ".env.example found",
			Category:    "env",
			Severity:    SeverityMedium,
			Status:      StatusPass,
			Message:     ".env.example file found",
			Remediation: "N/A",
		}
	}

	return Finding{
		ID:          "env.env_example_exists",
		Title:       ".env.example missing",
		Category:    "env",
		Severity:    SeverityLow,
		Status:      StatusWarn,
		Message:     "No .env.example file found",
		Remediation: "Create a .env.example file as a template for environment variables",
	}
}
