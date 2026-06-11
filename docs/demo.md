# ShipCheck Demo

Run these commands from the repository root.

## Current Repository

Audit ShipCheck itself:

```bash
go run ./cmd/shipcheck audit .
```

For the checks used by the project CI:

```bash
go run ./cmd/shipcheck audit . --category repo,env,ci --fail-under 70
```

## Example Services

The good fixture includes healthy Docker, CI, Kubernetes, and Terraform signals:

```bash
go run ./cmd/shipcheck audit examples/good-service
```

The risky fixture intentionally includes common readiness issues with fake placeholder values only:

```bash
go run ./cmd/shipcheck audit examples/risky-service
```

These examples are nested fixtures, so Git repository detection reports `no` when run in place.

## JSON Output

```bash
go run ./cmd/shipcheck audit examples/good-service --format json
```

JSON is useful when another script needs to inspect the score, counts, or finding statuses.

## Markdown Output

```bash
go run ./cmd/shipcheck audit examples/good-service --format markdown
```

Markdown is useful for saving a readable report as an artifact or copying a summary into a review.

## Docker-Only Audit

```bash
go run ./cmd/shipcheck audit examples/good-service --category docker
go run ./cmd/shipcheck audit examples/risky-service --category docker
```

The good fixture should pass the Docker checks. The risky fixture should show warnings and a failure around missing hardening and `.env` copy behavior.

## Fail-Under Threshold

```bash
go run ./cmd/shipcheck audit examples/good-service --fail-under 90
go run ./cmd/shipcheck audit examples/risky-service --fail-under 80
```

The good service should score much higher because most checks can validate healthy source files. The risky service should score lower because it intentionally has warnings or failures across Docker, CI, Kubernetes, Terraform, and environment handling.
