# Go API Template (Golang)

Opinionated starter for building a scalable Go HTTP API with Postgres (and optional RabbitMQ), with a clean transport/service/repository split and consistent JSON responses.

## The Idea

This template aims to stay **flat, explicit, and easy to extend** while keeping clear boundaries:

- **Transport**: thin edge (HTTP handlers / queue consumers). Parse input, call services, render output.
- **Service**: business logic + orchestration (and optional event publishing).
- **Repository**: data access via `sqlx` + SQL (easy to debug and optimize).

Adding a new domain should feel repeatable: create a small slice in transport/service/repository, wire it, ship it.

## Core decisions

- **Router**: `chi` (`github.com/go-chi/chi/v5`) for composable routes and middleware.
- **Responses**: a single renderer (`internal/libs/renderer`) standardizes JSON envelopes.
- **Errors**: return `internal/errors.HTTPError` to get consistent `statusCode` + `errorCode`.
- **Config**: one schema in `config/config.go`, loaded from env; accessible via `config.CONFIG`.
- **Queue**: RabbitMQ is optional; when disabled, publishing becomes a no-op (`NoopPublisher`).

## Features

- **HTTP routing**: `chi` with composable middleware
- **Postgres**: `sqlx` repositories + SQL migrations via `golang-migrate`
- **Queue (optional)**: RabbitMQ publish/consume with routing-key → handler router
- **Config**: environment-driven (supports `-envfilename`)
- **API ergonomics**: standardized JSON success/error envelopes
- **Ops/dev**: Docker Compose (DB + RabbitMQ), Air hot-reload, `golangci-lint`

## Quickstart (local dev)

Prereqs: Go (see `go.mod`), Docker (recommended for DB/queue).

1. Create env file:

```bash
cp .env.example .env
```

2. Start dependencies:

```bash
docker compose up
```

3. Run migrations:

```bash
go run cmd/migration/main.go -envfilename=.env
```

4. Run the API:

```bash
go run cmd/api/main.go
```

Optional hot reload:

```bash
go install github.com/cosmtrek/air@latest
air
```

## Configuration

Create `.env` from `.env.example`. Minimum keys to boot:

- **App**: `ENV`, `PORT`, `VERSION`, `ALLOWED_ORIGINS`
- **Database**: `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`
- **RabbitMQ (optional)**:
  - `RABBITMQ_ENABLED` (`true|false`)
  - `RABBITMQ_URL` (required when enabled)
  - `RABBITMQ_PREFETCH` (optional, default `20`)

Docker-first defaults (used by `docker-compose.yml`):

- `DB_HOST=go-api-template-db`
- `RABBITMQ_URL=amqp://guest:guest@go-api-template-rabbitmq:5672/`

If you run services outside Docker, switch hosts to `localhost`.

## Project layout (high level)

```
cmd/                    # Entrypoints (api server, migration runner)
config/                 # Env/config schema + loader
internal/
  server.go             # Wiring + graceful shutdown
  transport/            # HTTP + queue entrypoints (thin)
  service/              # Use-cases / orchestration
  repositories/         # sqlx data access
  model/                # Domain + DB models
  errors/               # Typed HTTP errors + error codes
  libs/                 # Shared infra (db, queue, renderer, utils, crypto)
  migrations/           # SQL migrations
scripts/                # Helper scripts
```

## Architecture (quick tour)

If you’re looking for “where to change what”, start here:

- **Entrypoints**: `cmd/api/main.go`, `cmd/migration/main.go`
- **Composition + lifecycle**: `internal/server.go` (router, middleware, DB, optional queue, graceful shutdown)
- **HTTP wiring**: `internal/transport/http/http_transport.go` (+ per-domain folders under `internal/transport/http/`)
- **Queue wiring**: `internal/transport/queue/queue.go` + `internal/transport/queue/router.go`
- **Migrations**: `internal/migrations/` (plain SQL, `golang-migrate` format)

## Conventions

- **Layering**:
  - **Transport** parses input, calls service, renders output
  - **Service** contains business logic (and publishes events)
  - **Repository** talks to Postgres (`sqlx` + SQL)
- **Response envelopes** (via `internal/libs/renderer`):
  - \(2xx\): `{ "data": <any>, "timestamp": "<utc>" }`
  - \(4xx/5xx\): `{ "errorCode": <int>, "statusCode": <int>, "message": "<string>", "timestamp": "<utc>" }`
- **Errors**: return `internal/errors.HTTPError` to control `statusCode` + `errorCode` consistently.

## Examples included

- **Users HTTP routes**: see `internal/transport/http/users/handler_routes.go` (e.g. `GET /api/users/{id}`, `POST /api/users/`)
- **Auth middleware scaffold**: `internal/transport/http/middlewares/auth.go` (replace with JWT/session validation)
- **Queue handler example**: see `internal/transport/queue/` for a sample routing-key handler

## Adding a new domain (checklist)

For a resource like `orders`:

- **Model**: `internal/model/orders.go`
- **Repo**: `internal/repositories/order_repository.go`
- **Service**: `internal/service/orders.go` (wire in `internal/service/types.go`)
- **HTTP**: `internal/transport/http/orders/` (handler + routes + types/mapper)
- **Wire**: add to `internal/transport/http/http_transport.go` and register under `/api`
- **Queue (optional)**: add routing keys/topology/handlers under `internal/libs/queue` + `internal/transport/queue`

Also: rename the module in `go.mod` and update imports from `go-api-template` to your module path.

## Lint

```bash
golangci-lint run
```

## License

MIT — see [LICENSE](./LICENSE).
