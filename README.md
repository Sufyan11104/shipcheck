# ShipCheck

ShipCheck is a Go CLI static analyser that scans repositories for deployment-readiness signals across common DevOps files.

It is a portfolio and learning project for turning practical deployment hygiene checks into a small, CI-friendly command-line tool. ShipCheck does not try to prove that an application is safe to deploy. It gives early warnings about missing or risky repository, Docker, GitHub Actions, Kubernetes, and Terraform patterns before they become operational problems.

## Why It Exists

Small services often accumulate deployment risk in plain text files: a Dockerfile without a non-root user, a workflow with no permissions block, a Kubernetes deployment without probes, or Terraform without provider constraints. ShipCheck collects those checks into one lightweight audit so a developer can quickly see where a repository needs attention.

## What It Checks

- Repository hygiene: README, `.gitignore`, `.env`, and `.env.example` signals
- Docker readiness: Dockerfile presence, `.dockerignore`, non-root `USER`, `HEALTHCHECK`, `.env` copy detection, and secret-like `ARG`/`ENV` names
- GitHub Actions: workflow presence, test/build steps, deploy-before-test ordering, action pinning, secret echo patterns, and explicit permissions
- Kubernetes: manifest/workload detection, readiness and liveness probes, resource requests/limits, image tags, and replica configuration
- Terraform/IaC: Terraform files, fmt/validate automation hints, required providers, provider version constraints, backend block, suspicious variable defaults, and lockfile presence

See [docs/rules.md](docs/rules.md) for rule details and limitations.

## Quick Start

```bash
go run ./cmd/shipcheck version
go run ./cmd/shipcheck audit .
```

Build a local binary:

```bash
make build
./bin/shipcheck audit .
```

## Example Commands

Audit the current repository:

```bash
go run ./cmd/shipcheck audit .
```

Audit one example service:

```bash
go run ./cmd/shipcheck audit examples/good-service
go run ./cmd/shipcheck audit examples/risky-service
```

Focus on Docker checks:

```bash
go run ./cmd/shipcheck audit examples/risky-service --category docker
```

Use ShipCheck as a CI threshold:

```bash
go run ./cmd/shipcheck audit . --category repo,env,ci --fail-under 70
```

## Example Output

```text
ShipCheck Deployment Readiness Report
Path: examples/risky-service
Score: 0/100
Passed: 11
Warnings: 22
Failed: 1
Skipped: 0

! docker.dockerfile_non_root_user - No non-root USER instruction detected in Dockerfile
x docker.dockerfile_no_env_copy - Dockerfile may be copying .env file directly
! k8s.readiness_probe_exists - No readiness probe detected in containers
```

## Report Formats

Text output is the default:

```bash
go run ./cmd/shipcheck audit .
```

JSON and Markdown are available for automation, reports, or pull request comments:

```bash
go run ./cmd/shipcheck audit . --format json
go run ./cmd/shipcheck audit . --format markdown
```

## Category Filtering

Use `--category` to limit the audit to one or more rule groups:

```bash
go run ./cmd/shipcheck audit . --category docker
go run ./cmd/shipcheck audit . --category repo,env,ci
```

Supported categories are `ci`, `docker`, `docs`, `env`, `k8s`, `repo`, and `terraform`.

## CI Thresholds

`--fail-under` exits with a non-zero status when the calculated score is below the requested threshold:

```bash
go run ./cmd/shipcheck audit . --fail-under 80
```

ShipCheck's own GitHub Actions workflow runs tests, `go vet`, and targeted ShipCheck audits on pushes and pull requests. See [docs/ci.md](docs/ci.md) for CI usage notes.

## Example Repositories

- `examples/good-service`: a small service fixture with healthy Docker, CI, Kubernetes, and Terraform patterns.
- `examples/risky-service`: a small service fixture with intentional readiness issues using fake placeholder values only.

Demo commands are documented in [docs/demo.md](docs/demo.md).

## Limitations

- ShipCheck is static analysis only.
- Checks are heuristic and intentionally lightweight.
- ShipCheck does not run Docker, Kubernetes, Terraform, or cloud APIs.
- It does not guarantee deployment safety.
- It is not a replacement for dedicated tools such as Hadolint, actionlint, kube-score, Checkov, Trivy, or Terraform's own validation commands.

## Roadmap

- SARIF output for code scanning integrations
- Configuration file for rule selection and thresholds
- Rule suppression for accepted risks
- Richer YAML and HCL parsing

## Development

```bash
go test ./...
go vet ./...
make build
```

Useful project areas:

- `cmd/shipcheck`: CLI entrypoint
- `internal/cli`: command parsing and flag handling
- `internal/scanner`: file scanning and repository metadata
- `internal/rules`: rule implementations
- `internal/engine`: rule orchestration, scoring, and filtering
- `internal/report`: text, JSON, and Markdown output
- `examples/`: demo fixtures

## CV / Project Summary

- Built ShipCheck, a Go CLI DevOps static-analysis tool that audits repository, Docker, CI, Kubernetes, and Terraform readiness signals, supports JSON/Markdown output and CI thresholds, and includes realistic demo fixtures plus GitHub Actions CI.
