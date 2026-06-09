# ShipCheck

**Status: Early WIP - Stage 2**

A Go CLI static-analysis tool that audits repositories for deployment readiness.

## Current Capabilities (Stage 2)

- CLI with `version` and `audit` commands
- Directory scanning (file and directory counting)
- Git repository detection
- **Rule engine with 4 initial repository hygiene checks:**
  - README detection
  - .gitignore detection
  - .env file hygiene (warns if committed)
  - .env.example template detection
- **Deployment readiness scoring (0-100)**
  - Automatic score calculation based on findings
  - Different weight penalties for high/medium/low severity findings
  - Pass/warn/fail status categorization

## Planned Features

Future stages will add:
- More comprehensive checks (Docker, Kubernetes, Terraform, GitHub Actions)
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
- **Rule definitions**: `internal/rules/types.go`, `internal/rules/checks.go`
- **Audit engine**: `internal/engine/engine.go`
- **Scoring logic**: `internal/engine/scoring.go`
- **Report formatting**: `internal/report/text.go`
- **Version**: `internal/version/version.go`
