# ShipCheck Rules

ShipCheck groups findings by the part of the deployment workflow they touch. The checks are designed to produce practical early warnings, not complete validation.

## Repository, Docs, And Env

What ShipCheck checks:

- README file presence
- `.gitignore` file presence
- `.env` file presence when auditing a Git repository path
- `.env.example` file presence

Why it matters:

These files set basic expectations for maintainability and local setup. A committed or root-level `.env` file is a signal that secrets or local-only values may accidentally enter source control.

Possible findings:

- `README file found`
- `No .gitignore file found`
- `.env file found - ensure it is in .gitignore`
- `No .env.example file found`

Limitations:

ShipCheck checks file presence and simple repository signals. It does not inspect README quality, parse `.gitignore`, scan Git history, or prove that a `.env` file was committed.

## Docker

What ShipCheck checks:

- Dockerfile presence
- `.dockerignore` presence
- Non-root `USER` instruction
- `HEALTHCHECK` instruction
- Direct `.env` copy into the image
- Secret-like `ARG` or `ENV` names

Why it matters:

These are common container hardening and operability signals. They can catch simple issues such as running as root, missing health checks, or accidentally including local environment files in an image.

Possible findings:

- `Dockerfile found`
- `No non-root USER instruction detected in Dockerfile`
- `Dockerfile may be copying .env file directly`
- `Potential secret-like Docker ARG or ENV name detected`

Limitations:

The Docker checks are line-oriented heuristics for a root-level Dockerfile. ShipCheck does not build images, evaluate all Dockerfile semantics, understand every multi-stage edge case, or replace Hadolint.

## GitHub Actions

What ShipCheck checks:

- `.github/workflows` directory and workflow YAML files
- Test and build command signals
- Deploy-like terms appearing before test commands
- Version-pinned action references
- Simple secret echo patterns
- Explicit `permissions` block

Why it matters:

CI workflows often become the deployment gate. Missing tests, unpinned actions, broad permissions, or secret logging patterns can increase release risk.

Possible findings:

- `Workflow files contain test step`
- `Workflow file may have deployment before testing`
- `Workflow actions are not pinned to a version`
- `No explicit permissions block detected in workflows`

Limitations:

ShipCheck scans workflow text and simple YAML-like patterns. It does not fully parse GitHub Actions semantics, evaluate reusable workflows, understand all shell behavior, or replace actionlint.

## Kubernetes

What ShipCheck checks:

- Kubernetes-like manifest presence
- Workload kinds such as Deployment, StatefulSet, or DaemonSet
- Readiness and liveness probes
- Resource requests and limits
- `:latest` or untagged image signals
- Multi-replica Deployment configuration

Why it matters:

These checks highlight basic runtime reliability signals. Probes, resource settings, stable image tags, and multiple replicas help make services easier to operate.

Possible findings:

- `Kubernetes workload manifest detected`
- `No readiness probe detected in containers`
- `Container image uses :latest tag`
- `No multi-replica configuration detected`

Limitations:

ShipCheck does not validate Kubernetes schemas, apply manifests, inspect clusters, or evaluate policies. It uses text heuristics and is not a replacement for kube-score, kubeconform, or admission policy tooling.

## Terraform

What ShipCheck checks:

- Terraform or tfvars file presence
- `terraform fmt` and `terraform validate` automation hints
- `required_providers` block
- Provider version constraints
- Backend block
- Secret-like variable defaults
- `.terraform.lock.hcl` presence

Why it matters:

Terraform configuration benefits from formatting, validation, provider pinning, remote state, and careful secret handling. These checks flag common gaps before infrastructure code reaches CI or review.

Possible findings:

- `Terraform files detected`
- `No required_providers block detected in Terraform files`
- `Terraform provider declarations do not appear to include version constraints`
- `Terraform variable defaults may contain secret-like names or values`

Limitations:

ShipCheck reads Terraform files as text. It does not run `terraform init`, `terraform validate`, providers, plans, cloud APIs, or a full HCL parser. It is not a replacement for Terraform, Checkov, Trivy, or policy-as-code tools.
