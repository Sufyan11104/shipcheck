# ShipCheck

**Status: Early WIP - Stage 1**

A Go CLI static-analysis tool that audits repositories for deployment readiness.

## Current Capabilities (Stage 1)

- Basic CLI structure with `version` and `audit` commands
- Directory scanning (file and directory counting)
- Git repository detection
- Placeholder deployment readiness report

## Planned Features

Future stages will add:
- Comprehensive deployment readiness scoring
- Docker configuration analysis
- GitHub Actions workflow validation
- Kubernetes manifest audits
- Terraform code analysis
- Environment file hygiene checks
- Secrets detection
- Multiple report formats (JSON, Markdown, SARIF)
- Rule engine with suppressions
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
- **Report formatting**: `internal/report/text.go`
- **Version**: `internal/version/version.go`
