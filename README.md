# ShipCheck

**Status: Early WIP - Stage 6**

A Go CLI static-analysis tool that audits repositories for deployment readiness.

## Current Capabilities (Stage 6)

- CLI with `version` and `audit` commands
- Directory scanning (file and directory counting)
- Git repository detection
- **Rule engine with 34 checks across 6 categories:**
  - **Repository hygiene (4):** README, .gitignore, .env, .env.example
  - **Docker readiness (6):** Dockerfile presence, .dockerignore, non-root USER, HEALTHCHECK, .env copy detection, secret-like ENV/ARG detection
  - **GitHub Actions CI/CD (8):** Workflows directory, workflow files, test steps, build steps, deploy order, action version pinning, secret echo detection, permissions block
  - **Kubernetes manifests (8):** Manifest detection, workload types, readiness probes, liveness probes, resource requests, resource limits, image tag validation, replica configuration
  - **Terraform/IaC readiness (8):** Terraform file detection, fmt and validate workflow recommendations, required providers, provider version constraints, backend block detection, suspicious variable defaults, dependency lockfile detection
- **Deployment readiness scoring (0-100)**
  - Automatic score calculation based on findings
  - Different weight penalties for high/medium/low severity findings
  - Pass/warn/fail status categorization

## Planned Features

Future stages will add:
- Multiple report formats (JSON, Markdown, SARIF)
- Rule engine with suppressions and custom rules
- Configuration file support
- CI/CD integration with configurable thresholds

## Building

```bash
make build
```

## Running

```bash
# Show version
make version

# Audit current directory
make run

# Or directly
./bin/shipcheck audit <path>
```

## Testing

```bash
make test
```

## Development

- **CLI entrypoint**: `cmd/shipcheck/main.go`
- **CLI routing**: `internal/cli/root.go`
- **Scanner logic**: `internal/scanner/scanner.go`
- **Rule definitions**: `internal/rules/` (types.go, checks.go, docker.go, github_actions.go, kubernetes.go, terraform.go)
- **Audit engine**: `internal/engine/engine.go`
- **Scoring logic**: `internal/engine/scoring.go`
- **Report formatting**: `internal/report/text.go`
- **Version**: `internal/version/version.go`
