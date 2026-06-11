# Architecture

ShipCheck is intentionally small. The code is split by CLI handling, scanning, rule execution, scoring, and reporting.

## Entry Point

`cmd/shipcheck/main.go` is the executable entrypoint. It delegates command handling to the internal CLI package.

## CLI

`internal/cli` parses commands and flags such as `audit`, `--format`, `--category`, and `--fail-under`. It coordinates scanner, engine, and report packages, then returns errors instead of exiting directly from deep logic.

## Scanner

`internal/scanner` validates the target path, counts files and directories, and records whether the scanned path itself contains a `.git` directory.

## Rules

`internal/rules` contains individual checks for repository hygiene, Docker, GitHub Actions, Kubernetes, and Terraform. Rules return structured findings with category, severity, status, message, and remediation text.

## Engine

`internal/engine` runs the rule set, calculates the readiness score, summarizes finding counts, and filters findings by category.

## Reports

`internal/report` renders the shared audit model as text, JSON, or Markdown.

## Examples

`examples/` contains small fake repositories used for demos and regression tests:

- `examples/good-service` shows mostly healthy deployment-readiness signals.
- `examples/risky-service` intentionally triggers common warnings and failures with fake placeholder values.
