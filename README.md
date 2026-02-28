# Backend Go Template

A production-ready Go backend template built with [Echo v5](https://echo.labstack.com/), GORM, Redis, and PostgreSQL. Follows a clean architecture pattern (transport → application → domain → infrastructure) to keep business logic decoupled from framework and infrastructure concerns.

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Local Development](#local-development)
  - [Docker](#docker)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Development Commands](#development-commands)
- [CI](#ci)

---

## Features

- JWT-based authentication (login, register, refresh, logout)
- Google OAuth login
- PostgreSQL via GORM ORM
- Redis cache layer (standalone / cluster / sentinel)
- Prometheus metrics (`/metrics`)
- Health-check endpoints (`/liveness`, `/readiness`)
- Structured JSON logging with `log/slog`
- Request validation with `go-playground/validator`
- Mock generation with `mockery`
- Multi-stage Docker build producing a minimal static binary

## Tech Stack

| Concern | Library |
|---|---|
| HTTP Framework | [Echo v5](https://github.com/labstack/echo) |
| ORM | [GORM](https://gorm.io) + `pgx` driver |
| Cache | [go-redis v9](https://github.com/redis/go-redis) |
| Config | [Viper](https://github.com/spf13/viper) |
| JWT | [golang-jwt/jwt v5](https://github.com/golang-jwt/jwt) |
| Validation | [go-playground/validator v10](https://github.com/go-playground/validator) |
| Metrics | [echoprometheus](https://github.com/labstack/echo-contrib) |
| Health checks | [hellofresh/health-go v5](https://github.com/hellofresh/health-go) |
| HTTP client | [resty v3](https://resty.dev) |
| Testing | [testify](https://github.com/stretchr/testify) |

## Project Structure

```
.
├── cmd/
│   └── main.go              # Entry point – initialises and starts the server
├── config/
│   ├── config.go            # Config structs & Viper loader
│   └── config.yml           # Default configuration file
├── internal/
│   ├── server.go            # Wire-up: connects DB/cache, builds the Echo instance
│   ├── application/
│   │   └── auth/            # Auth use-cases (commands & queries)
│   ├── domain/
│   │   ├── dtos/            # Request / response data transfer objects
│   │   └── models/          # Domain models
│   ├── infras/
│   │   ├── httpclient/      # Outbound HTTP client adapters
│   │   ├── repository/      # GORM repository implementations
│   │   └── security/        # JWT token manager
│   └── transport/
│       └── http/
│           ├── echo.go      # Echo setup (middleware, routes, error handler)
│           ├── handler/     # HTTP handlers
│           ├── helper/      # Response writers & error helpers
│           └── router/      # Route registration
├── pkgs/
│   ├── cache/               # Redis cache abstraction
│   ├── db/
│   │   ├── orm/             # GORM connection factory
│   │   └── rdclient/        # Redis client factory
│   ├── decorator/           # Command / query decorators (logging, metrics, …)
│   ├── hltcheck/            # Health-check service factory
│   ├── logger/              # JSON slog handler
│   ├── svcerr/              # Domain error types
│   └── utils/               # Shared utilities (validation, etc.)
├── mocks/                   # Auto-generated mocks (mockery)
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── .golangci.yml
```

## Getting Started

### Prerequisites

- [Go 1.26+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/)
- [`golangci-lint`](https://golangci-lint.run/usage/install/) (for linting)

### Local Development

1. **Clone the repository**

   ```bash
   git clone https://github.com/tdatIT/backend-go-template.git
   cd backend-go-template
   ```

2. **Start infrastructure services** (PostgreSQL, Redis, RabbitMQ)

   ```bash
   make docker-up
   ```

3. **Configure the application**

   The default configuration file is `config/config.yml`. Review and adjust the values (database credentials, JWT secret, etc.) as needed for your environment. You can also override settings via environment variables — any config key maps to an env var with dots replaced by underscores (e.g. `DATABASE_HOST`).

4. **Run the application**

   ```bash
   make run
   ```

   The HTTP server starts on `:5000` by default.

### Docker

Build and run the application in a container:

```bash
# Build the image
make docker-build

# Run with docker compose (starts the app + dependencies)
docker compose up
```

The `Dockerfile` uses a multi-stage build: a `golang:1.26-alpine` builder stage compiles a static binary, which is then copied into a minimal `alpine` runtime image.

## Configuration

Configuration is loaded from `config/config.yml` (or the path set by the `CONFIG_PATH` env var) and can be overridden by environment variables.

| Section | Key | Default | Description |
|---|---|---|---|
| `server` | `name` | `backend-go` | Service name (used in logs & metrics) |
| `server` | `httpPort` | `:5000` | HTTP listen address |
| `server` | `debugMode` | `false` | Enable debug mode |
| `database` | `host` | `localhost` | PostgreSQL host |
| `database` | `port` | `5432` | PostgreSQL port |
| `database` | `userName` | `postgres` | DB username |
| `database` | `password` | — | DB password |
| `database` | `database` | `backend_go` | Database name |
| `redis` | `mode` | `standalone` | `standalone`, `cluster`, or `sentinel` |
| `redis` | `address` | `127.0.0.1:6379` | Redis address(es) |
| `auth` | `jwtSecret` | `change_me` | **Change in production!** JWT signing secret |
| `auth` | `accessTokenTTL` | `15m` | Access token lifetime |
| `auth` | `refreshTokenTTL` | `720h` | Refresh token lifetime (30 days) |
| `logger` | `level` | `info` | Log level: `debug`, `info`, `warn`, `error` |

## API Endpoints

### Auth

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/v1/auth/login` | Login with username and password |
| `POST` | `/api/v1/auth/via-google` | Login with a Google OAuth token |
| `POST` | `/api/v1/auth/register` | Register a new account |
| `POST` | `/api/v1/auth/refresh` | Refresh an access token (Bearer refresh token) |
| `POST` | `/api/v1/auth/logout` | Logout and invalidate the session (Bearer access token) |

### Observability

| Method | Path | Description |
|---|---|---|
| `GET` | `/metrics` | Prometheus metrics |
| `GET` | `/liveness` | Liveness probe (returns `200 ok`) |
| `GET` | `/readiness` | Readiness probe (checks DB & Redis) |

## Development Commands

```bash
make build        # Compile binary to bin/app
make run          # Run the application
make test         # Run all tests
make lint         # Run golangci-lint
make fmt          # Format source files with gofmt
make tidy         # Run go mod tidy
make docker-build # Build Docker image (app:latest)
make docker-up    # Start infrastructure services via Docker Compose
make docker-down  # Stop infrastructure services
```

## CI

GitHub Actions (`.github/workflows/ci.yml`) runs on every push and pull request:

1. `go test ./...`
2. `golangci-lint`
3. Docker image build
