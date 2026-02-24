# AGENTS

## Project overview
- Go module: `github.com/tdatIT/backend-go-template` (`go.mod`).
- Single entrypoint: `cmd/main.go` (currently prints "Hello world").

## Local development
- Build: `make build` (outputs `bin/app` from `./cmd/main.go`).
- Run: `make run` (runs `./cmd/main.go`).
- Tests: `make test` (runs `go test ./...`).
- Format: `make fmt` (runs `gofmt -w .`).
- Tidy: `make tidy` (runs `go mod tidy`).

## Linting
- `make lint` uses `golangci-lint` with `govet`, `staticcheck`, `errcheck`, `revive`, `gofmt` (`.golangci.yml`).

## Docker
- Multi-stage build in `Dockerfile` produces a static binary from `./cmd/main.go` and exposes port 5000.
- Build image: `make docker-build` (tags `app:latest`).

## Local dependencies via Compose
- `docker-compose.yml` provisions:
  - Postgres 17 (`app`/`app`, DB `app`) on `5432`.
  - Redis 8 on `6379`.
  - RabbitMQ 3.13 management on `5672` and `15672` (`app`/`app`).
- Start/stop: `make docker-up` / `make docker-down`.

## CI
- GitHub Actions in `.github/workflows/ci.yml` runs `go test ./...`, `golangci-lint`, and a Docker build.

