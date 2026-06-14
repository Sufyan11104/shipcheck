# ShipCheck Deployment Readiness Report

## Summary

| Metric | Value |
| --- | --- |
| Path | . |
| Git repository | yes |
| Files scanned | 306 |
| Directories scanned | 185 |
| Score | 69/100 |
| Passed | 23 |
| Warnings | 5 |
| Failed | 0 |
| Skipped | 6 |

## Findings

| Status | ID | Category | Severity | Message | Remediation |
| --- | --- | --- | --- | --- | --- |
| pass | docs.readme_exists | docs | high | README file found | N/A |
| pass | repo.gitignore_exists | repo | medium | .gitignore file found | N/A |
| pass | env.env_not_committed | env | high | No committed .env file detected | N/A |
| pass | env.env_example_exists | env | medium | .env.example file found | N/A |
| skip | docker.dockerfile_exists | docker | low | No Dockerfile found | N/A |
| skip | docker.dockerignore_exists | docker | low | No .dockerignore file found | N/A |
| skip | docker.dockerfile_non_root_user | docker | low | Dockerfile not present; skipping USER check | N/A |
| skip | docker.dockerfile_healthcheck | docker | low | Dockerfile not present; skipping HEALTHCHECK check | N/A |
| skip | docker.dockerfile_no_env_copy | docker | low | Dockerfile not present; skipping .env COPY check | N/A |
| skip | docker.dockerfile_no_secret_env | docker | low | Dockerfile not present; skipping secret ENV check | N/A |
| pass | ci.workflows_dir_exists | ci | medium | .github/workflows directory found | N/A |
| pass | ci.workflow_file_exists | ci | medium | At least one workflow YAML file found | N/A |
| pass | ci.test_step_exists | ci | high | Workflow files contain test step | N/A |
| pass | ci.build_step_exists | ci | high | Workflow files contain build step | N/A |
| pass | ci.deploy_after_tests | ci | high | No obvious deploy-before-test pattern detected | N/A |
| pass | ci.actions_pinned | ci | medium | Workflow actions are pinned to versions | N/A |
| pass | ci.no_secret_echo | ci | high | Workflows do not appear to echo secrets | N/A |
| pass | ci.permissions_declared | ci | medium | Workflow includes explicit permissions block | N/A |
| pass | k8s.manifest_exists | k8s | medium | Kubernetes manifest file detected | N/A |
| pass | k8s.workload_exists | k8s | high | Kubernetes workload manifest detected (Deployment/StatefulSet/DaemonSet) | N/A |
| pass | k8s.readiness_probe_exists | k8s | medium | Readiness probe detected in container spec | N/A |
| pass | k8s.liveness_probe_exists | k8s | medium | Liveness probe detected in container spec | N/A |
| pass | k8s.resource_requests_exists | k8s | medium | Container resource requests detected | N/A |
| pass | k8s.resource_limits_exists | k8s | medium | Container resource limits detected | N/A |
| warn | k8s.no_latest_image_tag | k8s | medium | Container image uses :latest tag | Use specific image versions instead of :latest for reproducible deployments |
| pass | k8s.replicas_configured | k8s | medium | Deployment configured with multiple replicas | N/A |
| pass | terraform.files_exist | terraform | medium | Terraform files detected | N/A |
| warn | terraform.fmt_recommended | terraform | info | Terraform files detected; run terraform fmt in CI or local workflow | Add terraform fmt -check -recursive to CI or local automation |
| warn | terraform.validate_recommended | terraform | info | Terraform files detected; run terraform validate in CI or local workflow | Add terraform validate to CI or local automation after terraform init |
| pass | terraform.required_providers_exists | terraform | medium | Terraform required_providers block detected | N/A |
| warn | terraform.provider_versions_constrained | terraform | medium | Terraform provider declarations do not appear to include version constraints | Add version constraints to each provider in required_providers |
| pass | terraform.backend_configured | terraform | low | Terraform backend block detected | N/A |
| warn | terraform.no_suspicious_variable_defaults | terraform | high | Terraform variable defaults may contain secret-like names or values | Avoid committing secret defaults; use secret managers, environment variables, or CI/CD secret storage |
| pass | terraform.lockfile_present | terraform | low | .terraform.lock.hcl file detected | N/A |
