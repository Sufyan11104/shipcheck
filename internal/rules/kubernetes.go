package rules

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// CheckK8sManifestExists checks if any Kubernetes manifest exists
func CheckK8sManifestExists(path string) Finding {
	if hasK8sManifest(path) {
		return Finding{
			ID:          "k8s.manifest_exists",
			Title:       "Kubernetes manifest found",
			Category:    "k8s",
			Severity:    SeverityMedium,
			Status:      StatusPass,
			Message:     "Kubernetes manifest file detected",
			Remediation: "N/A",
		}
	}

	return Finding{
		ID:          "k8s.manifest_exists",
		Title:       "No Kubernetes manifest",
		Category:    "k8s",
		Severity:    SeverityMedium,
		Status:      StatusWarn,
		Message:     "No Kubernetes manifest files found",
		Remediation: "Create Kubernetes manifests for deployment (e.g., deployment.yaml)",
	}
}

// CheckK8sWorkloadExists checks if workload manifests exist (Deployment, StatefulSet, DaemonSet)
func CheckK8sWorkloadExists(path string) Finding {
	workloadKinds := []string{"deployment", "statefulset", "daemonset"}

	if !hasK8sManifest(path) {
		return Finding{
			ID:          "k8s.workload_exists",
			Title:       "No manifests found",
			Category:    "k8s",
			Severity:    SeverityLow,
			Status:      StatusSkip,
			Message:     "No Kubernetes manifests present; skipping workload check",
			Remediation: "N/A",
		}
	}

	if findInK8sFiles(path, workloadKinds) {
		return Finding{
			ID:          "k8s.workload_exists",
			Title:       "Workload manifest found",
			Category:    "k8s",
			Severity:    SeverityHigh,
			Status:      StatusPass,
			Message:     "Kubernetes workload manifest detected (Deployment/StatefulSet/DaemonSet)",
			Remediation: "N/A",
		}
	}

	return Finding{
		ID:          "k8s.workload_exists",
		Title:       "No workload manifest",
		Category:    "k8s",
		Severity:    SeverityHigh,
		Status:      StatusWarn,
		Message:     "No Kubernetes workload manifest found (Deployment, StatefulSet, or DaemonSet)",
		Remediation: "Create a workload manifest for your application deployment",
	}
}

// CheckK8sReadinessProbeExists checks if workloads define readiness probes
func CheckK8sReadinessProbeExists(path string) Finding {
	if !hasK8sManifest(path) {
		return Finding{
			ID:          "k8s.readiness_probe_exists",
			Title:       "No manifests found",
			Category:    "k8s",
			Severity:    SeverityLow,
			Status:      StatusSkip,
			Message:     "No Kubernetes manifests present; skipping readiness probe check",
			Remediation: "N/A",
		}
	}

	if findInK8sFiles(path, []string{"readinessprobe"}) {
		return Finding{
			ID:          "k8s.readiness_probe_exists",
			Title:       "Readiness probe found",
			Category:    "k8s",
			Severity:    SeverityMedium,
			Status:      StatusPass,
			Message:     "Readiness probe detected in container spec",
			Remediation: "N/A",
		}
	}

	return Finding{
		ID:          "k8s.readiness_probe_exists",
		Title:       "No readiness probe",
		Category:    "k8s",
		Severity:    SeverityMedium,
		Status:      StatusWarn,
		Message:     "No readiness probe detected in containers",
		Remediation: "Add a readinessProbe to help Kubernetes determine when your container is ready to accept traffic",
	}
}

// CheckK8sLivenessProbeExists checks if workloads define liveness probes
func CheckK8sLivenessProbeExists(path string) Finding {
	if !hasK8sManifest(path) {
		return Finding{
			ID:          "k8s.liveness_probe_exists",
			Title:       "No manifests found",
			Category:    "k8s",
			Severity:    SeverityLow,
			Status:      StatusSkip,
			Message:     "No Kubernetes manifests present; skipping liveness probe check",
			Remediation: "N/A",
		}
	}

	if findInK8sFiles(path, []string{"livenessprobe"}) {
		return Finding{
			ID:          "k8s.liveness_probe_exists",
			Title:       "Liveness probe found",
			Category:    "k8s",
			Severity:    SeverityMedium,
			Status:      StatusPass,
			Message:     "Liveness probe detected in container spec",
			Remediation: "N/A",
		}
	}

	return Finding{
		ID:          "k8s.liveness_probe_exists",
		Title:       "No liveness probe",
		Category:    "k8s",
		Severity:    SeverityMedium,
		Status:      StatusWarn,
		Message:     "No liveness probe detected in containers",
		Remediation: "Add a livenessProbe to automatically restart unhealthy containers",
	}
}

// CheckK8sResourceRequests checks if containers define resource requests
func CheckK8sResourceRequests(path string) Finding {
	if !hasK8sManifest(path) {
		return Finding{
			ID:          "k8s.resource_requests_exists",
			Title:       "No manifests found",
			Category:    "k8s",
			Severity:    SeverityLow,
			Status:      StatusSkip,
			Message:     "No Kubernetes manifests present; skipping resource requests check",
			Remediation: "N/A",
		}
	}

	if findInK8sFiles(path, []string{"requests:"}) {
		return Finding{
			ID:          "k8s.resource_requests_exists",
			Title:       "Resource requests found",
			Category:    "k8s",
			Severity:    SeverityMedium,
			Status:      StatusPass,
			Message:     "Container resource requests detected",
			Remediation: "N/A",
		}
	}

	return Finding{
		ID:          "k8s.resource_requests_exists",
		Title:       "No resource requests",
		Category:    "k8s",
		Severity:    SeverityMedium,
		Status:      StatusWarn,
		Message:     "No resource requests detected in containers",
		Remediation: "Define resources.requests to help Kubernetes scheduler place your pods appropriately",
	}
}

// CheckK8sResourceLimits checks if containers define resource limits
func CheckK8sResourceLimits(path string) Finding {
	if !hasK8sManifest(path) {
		return Finding{
			ID:          "k8s.resource_limits_exists",
			Title:       "No manifests found",
			Category:    "k8s",
			Severity:    SeverityLow,
			Status:      StatusSkip,
			Message:     "No Kubernetes manifests present; skipping resource limits check",
			Remediation: "N/A",
		}
	}

	if findInK8sFiles(path, []string{"limits:"}) {
		return Finding{
			ID:          "k8s.resource_limits_exists",
			Title:       "Resource limits found",
			Category:    "k8s",
			Severity:    SeverityMedium,
			Status:      StatusPass,
			Message:     "Container resource limits detected",
			Remediation: "N/A",
		}
	}

	return Finding{
		ID:          "k8s.resource_limits_exists",
		Title:       "No resource limits",
		Category:    "k8s",
		Severity:    SeverityMedium,
		Status:      StatusWarn,
		Message:     "No resource limits detected in containers",
		Remediation: "Define resources.limits to prevent runaway resource consumption",
	}
}

// CheckK8sNoLatestImageTag checks for :latest or untagged images
func CheckK8sNoLatestImageTag(path string) Finding {
	if !hasK8sManifest(path) {
		return Finding{
			ID:          "k8s.no_latest_image_tag",
			Title:       "No manifests found",
			Category:    "k8s",
			Severity:    SeverityLow,
			Status:      StatusSkip,
			Message:     "No Kubernetes manifests present; skipping image tag check",
			Remediation: "N/A",
		}
	}

	hasLatest := false
	hasNoTag := false

	filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml")) {
			file, err := os.Open(filePath)
			if err != nil {
				return nil
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				lower := strings.ToLower(line)

				// Check for :latest or image with no tag
				if strings.Contains(lower, "image:") {
					if strings.Contains(lower, ":latest") {
						hasLatest = true
					}
					// Very conservative: only warn if no colon at all after "image:"
					parts := strings.Split(line, "image:")
					if len(parts) > 1 {
						imagePart := strings.TrimSpace(parts[1])
						// Remove quotes if present
						imagePart = strings.Trim(imagePart, "\"'")
						if !strings.Contains(imagePart, ":") && !strings.Contains(imagePart, "$") {
							// Only warn if it looks like a real image name (not a variable)
							if strings.Contains(imagePart, "/") || !strings.Contains(imagePart, " ") {
								hasNoTag = true
							}
						}
					}
				}
			}
		}

		return nil
	})

	if hasLatest {
		return Finding{
			ID:          "k8s.no_latest_image_tag",
			Title:       ":latest tag detected",
			Category:    "k8s",
			Severity:    SeverityMedium,
			Status:      StatusWarn,
			Message:     "Container image uses :latest tag",
			Remediation: "Use specific image versions instead of :latest for reproducible deployments",
		}
	}

	if hasNoTag {
		return Finding{
			ID:          "k8s.no_latest_image_tag",
			Title:       "Untagged image detected",
			Category:    "k8s",
			Severity:    SeverityMedium,
			Status:      StatusWarn,
			Message:     "Container image appears to have no explicit tag",
			Remediation: "Always specify an explicit image tag for reliable deployments",
		}
	}

	return Finding{
		ID:          "k8s.no_latest_image_tag",
		Title:       "Images properly tagged",
		Category:    "k8s",
		Severity:    SeverityMedium,
		Status:      StatusPass,
		Message:     "Container images appear to use specific version tags",
		Remediation: "N/A",
	}
}

// CheckK8sReplicasConfigured checks if Deployments define proper replicas
func CheckK8sReplicasConfigured(path string) Finding {
	if !hasK8sManifest(path) {
		return Finding{
			ID:          "k8s.replicas_configured",
			Title:       "No manifests found",
			Category:    "k8s",
			Severity:    SeverityLow,
			Status:      StatusSkip,
			Message:     "No Kubernetes manifests present; skipping replicas check",
			Remediation: "N/A",
		}
	}

	// Look for replicas configuration in Deployment-like resources
	hasMultiReplicas := false
	hasDeployment := false

	filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml")) {
			file, err := os.Open(filePath)
			if err != nil {
				return nil
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				lower := strings.ToLower(line)

				if strings.Contains(lower, "kind:") && strings.Contains(lower, "deployment") {
					hasDeployment = true
				}

				if strings.Contains(lower, "replicas:") {
					// Check if replicas is > 1
					replicaPattern := regexp.MustCompile(`replicas:\s*(\d+)`)
					matches := replicaPattern.FindStringSubmatch(line)
					if len(matches) > 1 {
						if matches[1] != "1" && matches[1] != "0" {
							hasMultiReplicas = true
						}
					}
				}
			}
		}

		return nil
	})

	if hasDeployment {
		if hasMultiReplicas {
			return Finding{
				ID:          "k8s.replicas_configured",
				Title:       "Multi-replica configuration found",
				Category:    "k8s",
				Severity:    SeverityMedium,
				Status:      StatusPass,
				Message:     "Deployment configured with multiple replicas",
				Remediation: "N/A",
			}
		}

		return Finding{
			ID:          "k8s.replicas_configured",
			Title:       "No multi-replica configuration",
			Category:    "k8s",
			Severity:    SeverityMedium,
			Status:      StatusWarn,
			Message:     "No multi-replica configuration detected",
			Remediation: "Consider configuring multiple replicas for high availability",
		}
	}

	return Finding{
		ID:          "k8s.replicas_configured",
		Title:       "No Deployment found",
		Category:    "k8s",
		Severity:    SeverityLow,
		Status:      StatusSkip,
		Message:     "No Deployment resource found; replicas check not applicable",
		Remediation: "N/A",
	}
}

// hasK8sManifest checks if directory contains Kubernetes-like YAML files
func hasK8sManifest(path string) bool {
	k8sIndicators := []string{"apiversion:", "kind:", "metadata:", "spec:", "containers:"}

	var found bool
	filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil || found {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml")) {
			file, err := os.Open(filePath)
			if err != nil {
				return nil
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			lineCount := 0
			k8sLineCount := 0

			for scanner.Scan() {
				line := scanner.Text()
				lineCount++
				lower := strings.ToLower(line)

				for _, indicator := range k8sIndicators {
					if strings.Contains(lower, indicator) {
						k8sLineCount++
						break
					}
				}

				// If we find at least 2 K8s indicators, it's likely a K8s manifest
				if k8sLineCount >= 2 {
					found = true
					return nil
				}

				// Don't scan too many lines
				if lineCount > 100 {
					return nil
				}
			}
		}

		return nil
	})

	return found
}

// findInK8sFiles searches for patterns in Kubernetes manifest files
func findInK8sFiles(path string, patterns []string) bool {
	var found bool

	filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil || found {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml")) {
			file, err := os.Open(filePath)
			if err != nil {
				return nil
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				lower := strings.ToLower(line)

				for _, pattern := range patterns {
					if strings.Contains(lower, strings.ToLower(pattern)) {
						found = true
						return nil
					}
				}
			}
		}

		return nil
	})

	return found
}
