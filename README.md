# ShipCheck

**Status: Early WIP - Stage 3**

A Go CLI static-analysis tool that audits repositories for deployment readiness.

## Current Capabilities (Stage 3)

- CLI with `version` and `audit` commands
- Directory scanning (file and directory counting)
- Git repository detection
- **Rule engine with 10 checks across 3 categories:**
  - **Repository hygiene (4):** README, .gitignore, .env, .env.example
  - **Docker readiness (6):** Dockerfile presence, .dockerignore, non-root USER, HEALTHCHECK, .env copy detection, secret-like ENV/ARG detection
- **Deployment readiness scoring (0-100)**
  - Automatic score calculation based on findings
  - Different weight penalties for high/medium/low severity findings
  - Pass/warn/fail status categorization

## Planned Features

Future stages will add:
- GitHub Actions workflow validation
- Kubernetes manifest audits
- Terraform code analysis
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
- **Rule definitions**: `internal/rules/` (types.go, checks.go, docker.go)
- **Audit engine**: `internal/engine/engine.go`
- **Scoring logic**: `internal/engine/scoring.go`
- **Report formatting**: `internal/report/text.go`
- **Version**: `internal/version/version.go`
