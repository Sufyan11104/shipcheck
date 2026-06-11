# ShipCheck

**Status: Early WIP - Stage 7**

A Go CLI static-analysis tool that audits repositories for deployment readiness.

Example repositories are available under `examples/` and demo commands are documented in `docs/demo.md`.

## Current Capabilities (Stage 7)

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
- Text, JSON, and Markdown report output
- CI-friendly audit flags for score thresholds and category filtering

## Planned Features

Future stages will add:
- SARIF output
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

# Select report format
./bin/shipcheck audit <path> --format text
./bin/shipcheck audit <path> --format json
./bin/shipcheck audit <path> --format markdown

# Filter by category
./bin/shipcheck audit <path> --category docker,ci

# Fail when the score is below a threshold
./bin/shipcheck audit <path> --fail-under 80
```

## Testing

```bash
make test
```

CI runs tests and vet on pushes and pull requests. The workflow also runs ShipCheck against this repository and the `examples/good-service` fixture.

## Development

- **CLI entrypoint**: `cmd/shipcheck/main.go`
- **CLI routing**: `internal/cli/root.go`
- **Scanner logic**: `internal/scanner/scanner.go`
- **Rule definitions**: `internal/rules/` (types.go, checks.go, docker.go, github_actions.go, kubernetes.go, terraform.go)
- **Audit engine**: `internal/engine/engine.go`
- **Scoring logic**: `internal/engine/scoring.go`
- **Report formatting**: `internal/report/` (text.go, json.go, markdown.go)
- **Version**: `internal/version/version.go`
