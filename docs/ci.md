# Using ShipCheck In CI

ShipCheck can run as a lightweight readiness gate in CI. It works best as an early warning tool alongside language tests, linters, and dedicated scanners.

## GitHub Actions Example

```yaml
name: shipcheck

on:
  pull_request:
  push:
    branches:
      - main

permissions:
  contents: read

jobs:
  audit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
      - run: go test ./...
      - run: go vet ./...
      - run: go run ./cmd/shipcheck audit . --category repo,env,ci --fail-under 70
```

## Score Thresholds

`--fail-under` exits with status code `1` when the audit score is below the threshold:

```bash
go run ./cmd/shipcheck audit . --fail-under 80
```

Use a lower threshold while a project is adopting checks, then raise it as warnings are addressed.

## Targeted Checks

`--category` limits ShipCheck to specific rule groups:

```bash
go run ./cmd/shipcheck audit . --category docker
go run ./cmd/shipcheck audit . --category repo,env,ci
```

Targeted checks are useful when a repository does not use every supported technology or when a CI job should focus on one ownership area.

## Machine-Readable Output

JSON output can be parsed by scripts:

```bash
go run ./cmd/shipcheck audit . --format json
```

Markdown output can be saved as a job artifact or posted into a pull request by separate automation:

```bash
go run ./cmd/shipcheck audit . --format markdown
```

ShipCheck does not upload reports or create pull request comments by itself.
