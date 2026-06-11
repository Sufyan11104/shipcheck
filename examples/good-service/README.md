# Good Service

Small Go HTTP service used as a ShipCheck demo fixture.

## Setup

```bash
go test ./...
go run ./cmd/good-service
```

Copy `.env.example` to `.env` for local development only. Do not commit `.env`.

## Deployment Notes

- Build the container image from the root of this service.
- Run CI before promoting an image.
- Deploy with `k8s/deployment.yaml`.
- Review Terraform changes with `terraform fmt` and `terraform validate` before applying.
