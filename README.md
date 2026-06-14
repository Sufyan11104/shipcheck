# ShipCheck

ShipCheck is a Go CLI that audits repositories for deployment-readiness signals across Docker, GitHub Actions, Kubernetes, Terraform, and environment-file hygiene.

## Why ShipCheck

Repositories often contain deployment risks in plain text configuration files before anything reaches production. ShipCheck gives a quick static overview of whether a project has the basic signals expected from a deployable service: tests in CI, pinned actions, container health checks, Kubernetes probes, Terraform provider constraints, and safe environment-file conventions.

## Features

- Static repository audit
- Dockerfile and `.dockerignore` checks
- GitHub Actions CI/CD checks
- Kubernetes manifest readiness checks
- Terraform/IaC checks
- Environment-file hygiene checks
- Text, JSON, and Markdown reports
- Category filtering
- CI-friendly score thresholds

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

## Local Dashboard

Run a local browser dashboard for the current repository:

```bash
go run ./cmd/shipcheck serve .
open http://localhost:8080
```

Use a different address when needed:

```bash
go run ./cmd/shipcheck serve examples/good-service --addr localhost:8081
open http://localhost:8081
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
! k8s.readiness_probe_exists - No readiness probe detected in containers
! terraform.backend_configured - No backend block detected; remote state may be needed for team or production use.
```

## Commands And Flags

Audit a repository:

```bash
go run ./cmd/shipcheck audit .
go run ./cmd/shipcheck audit examples/good-service
```

Choose an output format:

```bash
go run ./cmd/shipcheck audit . --format text
go run ./cmd/shipcheck audit . --format json
go run ./cmd/shipcheck audit . --format markdown
```

Run selected categories:

```bash
go run ./cmd/shipcheck audit . --category docker,ci,k8s
go run ./cmd/shipcheck audit . --category terraform,env,docs,repo
```

Fail when the score is below a threshold:

```bash
go run ./cmd/shipcheck audit . --fail-under 80
```

Supported flags:

- `--format text|json|markdown`
- `--category docker,ci,k8s,terraform,env,docs,repo`
- `--fail-under 80`

## Rule Categories

| Category | Checks |
| --- | --- |
| `repo` / `docs` / `env` | README, `.gitignore`, `.env`, and `.env.example` signals |
| Docker | Dockerfile presence, `.dockerignore`, non-root `USER`, `HEALTHCHECK`, `.env` copy detection, secret-like `ARG`/`ENV` names |
| GitHub Actions | Workflow files, test/build steps, deploy-before-test ordering, action pinning, secret echo patterns, explicit permissions |
| Kubernetes | Manifest and workload detection, readiness/liveness probes, resource requests/limits, image tags, replica configuration |
| Terraform | Terraform files, fmt/validate automation hints, required providers, provider version constraints, backend block, variable defaults, lockfile |

See [docs/rules.md](docs/rules.md) for more detail on rule behavior and limitations.

## Example Repositories

ShipCheck includes two small fixtures for comparing healthy and risky deployment patterns:

- `examples/good-service`: healthy Docker, GitHub Actions, Kubernetes, and Terraform signals.
- `examples/risky-service`: intentional readiness issues using fake placeholder values only.

Try them with:

```bash
go run ./cmd/shipcheck audit examples/good-service
go run ./cmd/shipcheck audit examples/risky-service
```

More demo commands are in [docs/demo.md](docs/demo.md).

## Using In CI

ShipCheck can be used as a lightweight CI gate with a score threshold:

```bash
go run ./cmd/shipcheck audit . --category repo,env,ci --fail-under 70
```

Minimal GitHub Actions example:

```yaml
- uses: actions/checkout@v4
- uses: actions/setup-go@v5
  with:
    go-version: "1.22"
- run: go test ./...
- run: go vet ./...
- run: go run ./cmd/shipcheck audit . --category repo,env,ci --fail-under 70
```

See [docs/ci.md](docs/ci.md) for CI usage notes.

## Output Formats

- `text`: human-readable terminal output
- `json`: structured output for automation
- `markdown`: readable reports for artifacts, issues, or comments

## Limitations

- Static analysis only
- Heuristic checks
- Does not run Docker, `kubectl`, Terraform, or cloud APIs
- Does not replace specialised tools such as Hadolint, actionlint, kube-score, Checkov, Trivy, or `terraform validate`
- Does not guarantee deployment safety

## Development

```bash
go test ./...
go vet ./...
make build
```

## Project Layout

- `cmd/shipcheck`: CLI entrypoint
- `internal/cli`: command parsing and flag handling
- `internal/scanner`: file scanning and repository metadata
- `internal/rules`: rule implementations
- `internal/engine`: rule orchestration, scoring, and filtering
- `internal/report`: text, JSON, and Markdown output
- `examples`: demo repositories
- `docs`: rule, CI, demo, and architecture notes

## Roadmap

- SARIF output
- `.shipcheck.yaml` configuration
- Rule suppressions
- Richer YAML/HCL parsing
- More rule packs
