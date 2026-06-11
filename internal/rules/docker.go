package rules

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// CheckDockerfileExists checks if Dockerfile exists in the directory
func CheckDockerfileExists(path string) Finding {
	dockerfilePath := filepath.Join(path, "Dockerfile")
	if _, err := os.Stat(dockerfilePath); err == nil {
		return Finding{
			ID:          "docker.dockerfile_exists",
			Title:       "Dockerfile found",
			Category:    "docker",
			Severity:    SeverityMedium,
			Status:      StatusPass,
			Message:     "Dockerfile found",
			Remediation: "N/A",
		}
	}

	return Finding{
		ID:          "docker.dockerfile_exists",
		Title:       "Dockerfile missing",
		Category:    "docker",
		Severity:    SeverityMedium,
		Status:      StatusWarn,
		Message:     "No Dockerfile found",
		Remediation: "Create a Dockerfile to containerize your application",
	}
}

// CheckDockerignoreExists checks if .dockerignore exists
func CheckDockerignoreExists(path string) Finding {
	dockerignorePath := filepath.Join(path, ".dockerignore")
	if _, err := os.Stat(dockerignorePath); err == nil {
		return Finding{
			ID:          "docker.dockerignore_exists",
			Title:       ".dockerignore found",
			Category:    "docker",
			Severity:    SeverityLow,
			Status:      StatusPass,
			Message:     ".dockerignore file found",
			Remediation: "N/A",
		}
	}

	return Finding{
		ID:          "docker.dockerignore_exists",
		Title:       ".dockerignore missing",
		Category:    "docker",
		Severity:    SeverityLow,
		Status:      StatusWarn,
		Message:     "No .dockerignore file found",
		Remediation: "Create a .dockerignore file to exclude unnecessary files from Docker builds",
	}
}

// CheckDockerfileNonRootUser checks if Dockerfile uses a non-root USER instruction
func CheckDockerfileNonRootUser(path string) Finding {
	dockerfilePath := filepath.Join(path, "Dockerfile")
	file, err := os.Open(dockerfilePath)
	if err != nil {
		// Dockerfile doesn't exist, skip this check with info status
		return Finding{
			ID:          "docker.dockerfile_non_root_user",
			Title:       "Dockerfile not found",
			Category:    "docker",
			Severity:    SeverityLow,
			Status:      StatusSkip,
			Message:     "Dockerfile not present; skipping USER check",
			Remediation: "N/A",
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToUpper(trimmed), "USER") {
			parts := strings.Fields(trimmed)
			if len(parts) > 1 {
				user := parts[1]
				// Check if user is not root and not a numeric UID like 0
				if user != "root" && user != "0" {
					return Finding{
						ID:          "docker.dockerfile_non_root_user",
						Title:       "Non-root USER found",
						Category:    "docker",
						Severity:    SeverityHigh,
						Status:      StatusPass,
						Message:     "Dockerfile uses a non-root USER instruction",
						Remediation: "N/A",
					}
				}
			}
		}
	}

	return Finding{
		ID:          "docker.dockerfile_non_root_user",
		Title:       "No non-root USER instruction",
		Category:    "docker",
		Severity:    SeverityHigh,
		Status:      StatusWarn,
		Message:     "No non-root USER instruction detected in Dockerfile",
		Remediation: "Add a USER instruction with a non-root user for better security",
	}
}

// CheckDockerfileHealthcheck checks if Dockerfile has a HEALTHCHECK instruction
func CheckDockerfileHealthcheck(path string) Finding {
	dockerfilePath := filepath.Join(path, "Dockerfile")
	file, err := os.Open(dockerfilePath)
	if err != nil {
		return Finding{
			ID:          "docker.dockerfile_healthcheck",
			Title:       "Dockerfile not found",
			Category:    "docker",
			Severity:    SeverityLow,
			Status:      StatusSkip,
			Message:     "Dockerfile not present; skipping HEALTHCHECK check",
			Remediation: "N/A",
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToUpper(trimmed), "HEALTHCHECK") {
			return Finding{
				ID:          "docker.dockerfile_healthcheck",
				Title:       "HEALTHCHECK found",
				Category:    "docker",
				Severity:    SeverityMedium,
				Status:      StatusPass,
				Message:     "Dockerfile includes a HEALTHCHECK instruction",
				Remediation: "N/A",
			}
		}
	}

	return Finding{
		ID:          "docker.dockerfile_healthcheck",
		Title:       "No HEALTHCHECK instruction",
		Category:    "docker",
		Severity:    SeverityMedium,
		Status:      StatusWarn,
		Message:     "No HEALTHCHECK instruction detected in Dockerfile",
		Remediation: "Consider adding a HEALTHCHECK instruction for container monitoring",
	}
}

// CheckDockerfileNoEnvCopy checks if Dockerfile avoids copying .env directly
func CheckDockerfileNoEnvCopy(path string) Finding {
	dockerfilePath := filepath.Join(path, "Dockerfile")
	file, err := os.Open(dockerfilePath)
	if err != nil {
		return Finding{
			ID:          "docker.dockerfile_no_env_copy",
			Title:       "Dockerfile not found",
			Category:    "docker",
			Severity:    SeverityLow,
			Status:      StatusSkip,
			Message:     "Dockerfile not present; skipping .env COPY check",
			Remediation: "N/A",
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToUpper(trimmed), "COPY") {
			if strings.Contains(strings.ToUpper(trimmed), ".ENV") {
				return Finding{
					ID:          "docker.dockerfile_no_env_copy",
					Title:       ".env copy detected",
					Category:    "docker",
					Severity:    SeverityHigh,
					Status:      StatusFail,
					Message:     "Dockerfile may be copying .env file directly",
					Remediation: "Avoid copying .env into Docker images; use environment variables or .env.example instead",
				}
			}
		}
	}

	return Finding{
		ID:          "docker.dockerfile_no_env_copy",
		Title:       "No .env copy detected",
		Category:    "docker",
		Severity:    SeverityHigh,
		Status:      StatusPass,
		Message:     "Dockerfile does not copy .env file",
		Remediation: "N/A",
	}
}

// CheckDockerfileNoSecretEnv checks for secret-like ARG or ENV names
func CheckDockerfileNoSecretEnv(path string) Finding {
	dockerfilePath := filepath.Join(path, "Dockerfile")
	file, err := os.Open(dockerfilePath)
	if err != nil {
		return Finding{
			ID:          "docker.dockerfile_no_secret_env",
			Title:       "Dockerfile not found",
			Category:    "docker",
			Severity:    SeverityLow,
			Status:      StatusSkip,
			Message:     "Dockerfile not present; skipping secret ENV check",
			Remediation: "N/A",
		}
	}
	defer file.Close()

	secretKeywords := []string{"PASSWORD", "SECRET", "TOKEN", "API_KEY", "PRIVATE_KEY"}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		upperTrimmed := strings.ToUpper(trimmed)

		if strings.HasPrefix(upperTrimmed, "ARG ") || strings.HasPrefix(upperTrimmed, "ENV ") {
			for _, keyword := range secretKeywords {
				if strings.Contains(upperTrimmed, keyword) {
					return Finding{
						ID:          "docker.dockerfile_no_secret_env",
						Title:       "Potential secret-like variable",
						Category:    "docker",
						Severity:    SeverityHigh,
						Status:      StatusWarn,
						Message:     "Potential secret-like Docker ARG or ENV name detected (password, secret, token, api_key, private_key)",
						Remediation: "Avoid storing secrets in Dockerfile; use Docker secrets, environment variables at runtime, or secret management tools",
					}
				}
			}
		}
	}

	return Finding{
		ID:          "docker.dockerfile_no_secret_env",
		Title:       "No obvious secret variables",
		Category:    "docker",
		Severity:    SeverityHigh,
		Status:      StatusPass,
		Message:     "No obvious secret-like ARG or ENV variables detected",
		Remediation: "N/A",
	}
}
