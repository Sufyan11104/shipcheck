package rules

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckK8sManifestExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when no manifest exists
	finding := CheckK8sManifestExists(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn, got %v", finding.Status)
	}

	// Test when manifest exists
	manifestContent := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  replicas: 2
  template:
    spec:
      containers:
      - name: app
        image: myapp:v1.0
`
	os.WriteFile(filepath.Join(tmpDir, "deployment.yaml"), []byte(manifestContent), 0644)
	finding = CheckK8sManifestExists(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass, got %v", finding.Status)
	}
}

func TestCheckK8sWorkloadExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when no manifest exists
	finding := CheckK8sWorkloadExists(tmpDir)
	if finding.Status != StatusSkip {
		t.Errorf("expected StatusSkip (graceful skip), got %v", finding.Status)
	}

	// Test when manifest exists but no workload
	nonWorkloadContent := `apiVersion: v1
kind: ConfigMap
metadata:
  name: config
data:
  key: value
`
	os.WriteFile(filepath.Join(tmpDir, "config.yaml"), []byte(nonWorkloadContent), 0644)
	finding = CheckK8sWorkloadExists(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for no workload, got %v", finding.Status)
	}

	// Test when Deployment exists
	deploymentContent := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
`
	os.WriteFile(filepath.Join(tmpDir, "deployment.yaml"), []byte(deploymentContent), 0644)
	finding = CheckK8sWorkloadExists(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for Deployment, got %v", finding.Status)
	}
}

func TestCheckK8sReadinessProbeExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when no manifest exists
	finding := CheckK8sReadinessProbeExists(tmpDir)
	if finding.Status != StatusSkip {
		t.Errorf("expected StatusSkip (graceful skip), got %v", finding.Status)
	}

	// Test when manifest exists without readiness probe
	deploymentContent := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: app
        image: myapp:v1.0
`
	os.WriteFile(filepath.Join(tmpDir, "deployment.yaml"), []byte(deploymentContent), 0644)
	finding = CheckK8sReadinessProbeExists(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for no readiness probe, got %v", finding.Status)
	}

	// Test when manifest has readiness probe
	probeContent := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: app
        image: myapp:v1.0
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
`
	os.WriteFile(filepath.Join(tmpDir, "deployment.yaml"), []byte(probeContent), 0644)
	finding = CheckK8sReadinessProbeExists(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for readiness probe, got %v", finding.Status)
	}
}

func TestCheckK8sLivenessProbeExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when no manifest exists
	finding := CheckK8sLivenessProbeExists(tmpDir)
	if finding.Status != StatusSkip {
		t.Errorf("expected StatusSkip (graceful skip), got %v", finding.Status)
	}

	// Test when manifest exists without liveness probe
	deploymentContent := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: app
        image: myapp:v1.0
`
	os.WriteFile(filepath.Join(tmpDir, "deployment.yaml"), []byte(deploymentContent), 0644)
	finding = CheckK8sLivenessProbeExists(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for no liveness probe, got %v", finding.Status)
	}

	// Test when manifest has liveness probe
	probeContent := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: app
        image: myapp:v1.0
        livenessProbe:
          exec:
            command:
            - /bin/sh
            - -c
            - ps aux | grep app
`
	os.WriteFile(filepath.Join(tmpDir, "deployment.yaml"), []byte(probeContent), 0644)
	finding = CheckK8sLivenessProbeExists(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for liveness probe, got %v", finding.Status)
	}
}

func TestCheckK8sResourceRequests(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when no manifest exists
	finding := CheckK8sResourceRequests(tmpDir)
	if finding.Status != StatusSkip {
		t.Errorf("expected StatusSkip (graceful skip), got %v", finding.Status)
	}

	// Test when manifest exists without resource requests
	deploymentContent := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: app
        image: myapp:v1.0
`
	os.WriteFile(filepath.Join(tmpDir, "deployment.yaml"), []byte(deploymentContent), 0644)
	finding = CheckK8sResourceRequests(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for no resource requests, got %v", finding.Status)
	}

	// Test when manifest has resource requests
	resourceContent := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: app
        image: myapp:v1.0
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
`
	os.WriteFile(filepath.Join(tmpDir, "deployment.yaml"), []byte(resourceContent), 0644)
	finding = CheckK8sResourceRequests(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for resource requests, got %v", finding.Status)
	}
}

func TestCheckK8sResourceLimits(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when no manifest exists
	finding := CheckK8sResourceLimits(tmpDir)
	if finding.Status != StatusSkip {
		t.Errorf("expected StatusSkip (graceful skip), got %v", finding.Status)
	}

	// Test when manifest exists without resource limits
	deploymentContent := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: app
        image: myapp:v1.0
`
	os.WriteFile(filepath.Join(tmpDir, "deployment.yaml"), []byte(deploymentContent), 0644)
	finding = CheckK8sResourceLimits(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for no resource limits, got %v", finding.Status)
	}

	// Test when manifest has resource limits
	resourceContent := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: app
        image: myapp:v1.0
        resources:
          limits:
            memory: "128Mi"
            cpu: "500m"
`
	os.WriteFile(filepath.Join(tmpDir, "deployment.yaml"), []byte(resourceContent), 0644)
	finding = CheckK8sResourceLimits(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for resource limits, got %v", finding.Status)
	}
}

func TestCheckK8sNoLatestImageTag(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when no manifest exists
	finding := CheckK8sNoLatestImageTag(tmpDir)
	if finding.Status != StatusSkip {
		t.Errorf("expected StatusSkip (graceful skip), got %v", finding.Status)
	}

	// Test when manifest has specific tag
	deploymentContent := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: app
        image: myapp:v1.0
`
	os.WriteFile(filepath.Join(tmpDir, "deployment.yaml"), []byte(deploymentContent), 0644)
	finding = CheckK8sNoLatestImageTag(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for specific tag, got %v", finding.Status)
	}

	// Test when manifest uses :latest
	latestContent := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: app
        image: myapp:latest
`
	os.WriteFile(filepath.Join(tmpDir, "deployment.yaml"), []byte(latestContent), 0644)
	finding = CheckK8sNoLatestImageTag(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for :latest tag, got %v", finding.Status)
	}
}

func TestCheckK8sReplicasConfigured(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when no manifest exists
	finding := CheckK8sReplicasConfigured(tmpDir)
	if finding.Status != StatusSkip {
		t.Errorf("expected StatusSkip (graceful skip), got %v", finding.Status)
	}

	// Test when Deployment has no replicas field
	deploymentContent := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: app
        image: myapp:v1.0
`
	os.WriteFile(filepath.Join(tmpDir, "deployment.yaml"), []byte(deploymentContent), 0644)
	finding = CheckK8sReplicasConfigured(tmpDir)
	if finding.Status != StatusWarn {
		t.Errorf("expected StatusWarn for no replicas, got %v", finding.Status)
	}

	// Test when Deployment has multiple replicas
	multiReplicaContent := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: app
        image: myapp:v1.0
`
	os.WriteFile(filepath.Join(tmpDir, "deployment.yaml"), []byte(multiReplicaContent), 0644)
	finding = CheckK8sReplicasConfigured(tmpDir)
	if finding.Status != StatusPass {
		t.Errorf("expected StatusPass for multiple replicas, got %v", finding.Status)
	}
}
