# ShipCheck Demo

Run these commands from the repository root.

## Current Repository

```bash
go run ./cmd/shipcheck audit .
```

## Example Services

```bash
go run ./cmd/shipcheck audit examples/good-service
go run ./cmd/shipcheck audit examples/risky-service
```

The good service includes container, CI, Kubernetes, and Terraform basics. The risky service intentionally includes common issues such as a weak Dockerfile, unpinned workflow action, missing probes, and incomplete Terraform settings.

These examples are nested fixtures, so Git repository detection reports `no` when run in place.

## Formats

```bash
go run ./cmd/shipcheck audit examples/good-service --format json
go run ./cmd/shipcheck audit examples/good-service --format markdown
```

## Category Filtering

```bash
go run ./cmd/shipcheck audit examples/good-service --category docker
go run ./cmd/shipcheck audit examples/risky-service --category docker
```

## Fail Under

```bash
go run ./cmd/shipcheck audit examples/good-service --fail-under 80
go run ./cmd/shipcheck audit examples/risky-service --fail-under 80
```

The good service should score higher because its checks can validate healthy source files. The risky service should score lower because it has warnings and failures across Docker, CI, Kubernetes, Terraform, and environment handling.
