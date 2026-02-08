# Golang API Template

> Golang API Template with a clear, scalable structure that can sustain large APIs.

## Table of Contents

- [Features](#features)
- [Directory Structure](#directory-structure)
- [Description](#description)
- [Setup](#setup)
- [Template Tour](#template-tour)
- [License](#license)

## Features

- Standardized JSON responses (success + error envelopes)
- `chi` router with clean route composition and middleware chaining
- Postgres via `sqlx` (repository pattern) + migrations via `golang-migrate`
- Optional RabbitMQ transport (publish + consume) with a simple router for handlers
- Environment-driven configuration via `godotenv` (with `-envfilename` flag support)
- Security headers via `unrolled/secure` + CORS via `go-chi/cors`
- Structured logging via `logrus`
- Hot reload via Air (`.air.toml`)
- Docker setup for Postgres and RabbitMQ + optional app container
- Linting and formatting configuration via `golangci-lint` (`.golangci.yml`)

---

## Directory Structure

```
/cmd                         # App entrypoints (binaries)
  /api                       # HTTP API server
    main.go
  /migration                 # Migration runner
    main.go
/config                      # Config schema + env loading
  config.go
/internal                    # Application code
  server.go                  # Server wiring (HTTP + Queue) and graceful shutdown
  /transport                 # Transport layer (HTTP, queue, etc.) through which external systems communicate with the API.
    /http
      http_transport.go      # Transport constructor (wires handlers)
      /middlewares           # HTTP middlewares (auth + chaining helper)
      /users                 # User HTTP handlers + routes + request mapping
    /queue
      queue.go               # Consumer startup and routing keys -> handlers
      router.go              # Router: routingKey -> handler
      queue_handlers.go      # Message handlers (example: users.created)
  /service                   # Business logic layer (use-cases)
    users.go
    types.go                 # Service wiring + request structs
  /repositories              # Data access layer (sqlx)
    user_repository.go
  /model                     # Domain + DB models (sqlx tags)
    users.go
    model_kit.go
  /errors                    # Custom error type + error codes
    errors.go
    error_codes.go
    user_email_exists.go
  /libs                      # Shared infra helpers
    /crypto                  # Password hashing helpers
    /database                # Postgres connection + migrations runner
    /queue                   # RabbitMQ implementation + types + topology
    /renderer                # Response renderer (JSON envelopes)
    /utils                   # Generic helpers (JSON body parsing, parsing utils)
  /migrations                # SQL migrations (golang-migrate format)
    1_create_tables.up.sql
/scripts                      # Helper scripts (migrations, DB init)
  run-migrations-env.sh
  run-migrations-env-prod.sh
  database-init.sh
.air.toml                     # Air hot reload config
.golangci.yml                 # golangci-lint configuration
docker-compose.yml            # Postgres + RabbitMQ for local dev
docker-compose.app.yml        # Build/run app container
Dockerfile
.env.example                  # Env template
go.mod
go.sum
LICENSE
README.md
```

## Description

**The Why**

This template aims to stay flat and easy to extend while still being explicit about boundaries. It uses a layered approach:

- **Transport layer**: HTTP handlers and queue consumers. These should be thin; parse input, call services, render output.
- **Service layer**: business logic and orchestration (including publishing events).
- **Repository layer**: data access with `sqlx` and raw SQL.

The goal is to make adding new features predictable: every new domain (users, auth, orders, products, etc.) gets its own transport + service + repository slice without forcing a deep folder tree or overly clever abstractions.

**Core decisions**

- **Router**: `github.com/go-chi/chi/v5` for fast routing + composable middleware.
- **Responses**: a single renderer (`internal/libs/renderer`) standardizes JSON envelopes:
  - \(2xx\): `{ "data": <any>, "timestamp": "<utc>" }`
  - \(4xx\)/\(5xx\): `{ "errorCode": <int>, "statusCode": <int>, "message": "<string>", "timestamp": "<utc>" }`
- **Errors**: application errors are represented by `internal/errors.HTTPError`. If you return one of those, the renderer will pick it up and respond with its `StatusCode` and `ErrorCode` consistently.
- **Config**: one config schema in `config/config.go`, loaded from env. Once loaded, config is available globally via `config.CONFIG` (simple and practical; easy to refactor into dependency injection if you prefer).
- **Queue**: RabbitMQ is optional. When enabled, the service publishes events and the consumer transport starts consumers on startup. When disabled, publishing becomes a no-op via `NoopPublisher`.

## Setup

Go version: see `go.mod` (Go 1.25+).

Install Air (hot reload):

```
go install github.com/cosmtrek/air@latest
```

Download dependencies:

```
go mod download
```

```
go mod tidy
```

Create an `.env` file in the root folder and start from `.env.example`.

**Minimum variables to run locally**

These are the ones the template actually uses to boot the API:

- `ENV` (ex: `DEV`)
- `PORT` (ex: `8080`)
- `VERSION` (ex: `0.0.1`)
- `ALLOWED_ORIGINS` (ex: `*`)
- `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`
- `RABBITMQ_ENABLED` (ex: `false` if you don't want RabbitMQ)
- `RABBITMQ_URL` (required only when `RABBITMQ_ENABLED=true`)
- `RABBITMQ_PREFETCH` (optional, defaults to `20`)

**Default networking (Docker-first)**

By default, `.env` is configured to talk to the Docker Compose services:

- `DB_HOST=go-api-template-db`
- `RABBITMQ_URL=amqp://guest:guest@go-api-template-rabbitmq:5672/`

**Run Postgres + RabbitMQ with Docker**

```
docker compose up
```

That will start:

- Postgres on `5432` (service name/hostname `go-api-template-db`)
- RabbitMQ on `5672` and the management UI on `15672` (hostname `go-api-template-rabbitmq`)

If you want to run the API on your host machine instead of in Docker, switch the `.env` hosts to `localhost` accordingly.

**Run migrations**

```
go run cmd/migration/main.go -envfilename=.env
```

Or via the helper script:

```
./scripts/run-migrations-env.sh
```

For a production-style env file:

```
./scripts/run-migrations-env-prod.sh
```

**Run the API**

- With hot reload:

```
air
```

- Without hot reload:

```
go run cmd/api/main.go (it takes the .env file as default)
```

**Dockerize the API**

Build and run the app container:

```
docker compose -f docker-compose.app.yml up --build
```

**Linting / formatting**

This repo includes a `golangci-lint` configuration (`.golangci.yml`) and enables `gofmt`/`goimports`.

- Install and run lint:

```
golangci-lint run
```

## Template Tour

### Entrypoints (`/cmd`)

- **`cmd/api/main.go`**: loads config, sets up signal handling, starts the server, and ensures shutdown cleanup runs.
- **`cmd/migration/main.go`**: loads config, connects to Postgres, and runs migrations from `internal/migrations`.

### Server wiring (`internal/server.go`)

This is where the application is composed:

- Creates the Postgres connection (`internal/libs/database`)
- Optionally connects RabbitMQ and ensures topology (exchange + queues + bindings)
- Constructs services (`internal/service`) and transports (`internal/transport/http`, `internal/transport/queue`)
- Builds the `chi` router with middleware:
  - `middleware.Logger`, `middleware.Recoverer`
  - `cors.Handler` using `ALLOWED_ORIGINS`
  - `secure.New(...).Handler` for security headers
- Registers top-level routes under `/api`
- Starts:
  - the HTTP server
  - queue consumers (only when `RABBITMQ_ENABLED=true`)
- Handles graceful shutdown and closes DB/RabbitMQ connections

### HTTP transport (`internal/transport/http`)

**Pattern**

- Each domain gets its own folder under `internal/transport/http/<domain>/`
- Each domain exposes a `RegisterRoutes(r chi.Router)` method
- `internal/transport/http/http_transport.go` wires up domain handlers and injects:
  - the relevant service(s)
  - the `ResponseRenderer`

**Example: users**

Routes are registered in `internal/transport/http/users/handler_routes.go`:

- `GET /api/users/{id}` → `GetUser`
- `POST /api/users/` → `CreateUser` (wrapped with `AuthMiddleware`)

`AuthMiddleware` (`internal/transport/http/middlewares/auth.go`) is provided as a scaffold. It currently allows all requests; you can implement JWT validation and attach claims to context.

### Response rendering (`internal/libs/renderer`)

All handlers should respond via `ResponseRenderer.JSON(w, statusCode, payload)`.

- If `payload` is an `*errors.HTTPError` (or wraps one), the renderer will respond using the embedded `StatusCode` and the standard error envelope.
- If `payload` is any other `error`, the renderer returns a \(500\) internal error envelope.
- If `statusCode` is \(2xx\), the renderer wraps the payload in the success envelope.

This lets you keep handler logic simple while still enforcing consistent responses across the API.

### Services (`internal/service`) and repositories (`internal/repositories`)

**Repository layer**

Repositories use `sqlx` directly. Example: `UserRepository` implements:

- `GetUser(ctx, id)`
- `CreateUser(ctx, user)`
- `CheckIfUserEmailExists(ctx, email)`

**Service layer**

Services implement business logic and orchestration. Example: `UserService.CreateUser`:

- checks uniqueness
- hashes the password (`internal/libs/crypto`)
- inserts the user
- publishes a `users.created` event (when RabbitMQ is enabled)

### Migrations (`internal/migrations`)

Migrations are plain SQL and run via `golang-migrate`. The included migration creates:

- `uuid-ossp` extension
- `citext` extension
- an `email` domain with a regex constraint
- a `users` table

Add new migrations by creating additional `*.up.sql` files in `internal/migrations` using the standard migrate versioned naming.

### Queue transport (`internal/transport/queue`)

When enabled, the template:

- declares an exchange (`internal/libs/queue.EventsExchangeName`)
- declares queues + bindings (`internal/libs/queue/topology.go`)
- starts consumers in goroutines (`QueueTransport.StartConsumers`)
- routes messages via `Router` (`routingKey -> handler`)

**Example message**

On user creation, the service publishes:

- exchange: `go-api-template`
- routing key: `users.created`
- body: JSON (id, email, firstName, lastName)

The consumer example handler lives in `internal/transport/queue/queue_handlers.go`.

### Adding a new domain (recommended workflow)

If you want to add a new resource like `orders`:

- **Model**: add `internal/model/orders.go`
- **Repository**: add `internal/repositories/order_repository.go`
- **Service**: add `internal/service/orders.go` and wire it in `internal/service/types.go`
- **HTTP handlers**: add `internal/transport/http/orders/` with:
  - `handler.go`, `handler_routes.go`, `types.go`, `mapper.go`
- **Transport wiring**:
  - add it to `internal/transport/http/http_transport.go`
  - register routes in `internal/server.go` under `/api`
- **(Optional) Queue**:
  - define routing keys/queue names in `internal/libs/queue/types.go`
  - add queue specs/bindings in `internal/libs/queue/topology.go`
  - register handlers in `internal/transport/queue/queue.go`

Don’t forget to rename the module (`go.mod`) and imports from `go-api-template` to your own module path.

Happy coding.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.
