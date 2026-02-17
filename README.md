# Subscriptions Service

Test service for managing user subscriptions.

## Stack

- Go 1.25
- PostgreSQL 16
- pgxpool
- chi router
- Goose migrations
- Docker Compose

## Architecture

Three-layer structure:

- `internal/httpapi` - transport/handlers (HTTP + chi routing)
- `internal/service` - business layer (use cases)
- `internal/repository` - data access layer (PostgreSQL)

## Run locally

1. Install dependencies:

```bash
go mod tidy
```

2. Start PostgreSQL:

```bash
docker compose up -d db
```

3. Apply migrations:

```bash
make migrate-up
```

4. Run API:

```bash
make run
```

## Run full stack with Docker

```bash
docker compose up --build
```

## Endpoints

- `HEAD /health`
- `POST /api/v1/subscriptions`
- `GET /api/v1/subscriptions`
- `GET /api/v1/subscriptions/{id}`
- `PUT /api/v1/subscriptions/{id}`
- `DELETE /api/v1/subscriptions/{id}`
- `GET /api/v1/subscriptions/total?from=MM-YYYY&to=MM-YYYY`

## Swagger

Swagger UI is available at `GET /swagger/index.html` after starting the API.

## DB connection retries

Database connection uses fixed retry policy in code (`pkg/postgres/postgres.go`):

- retries: `5`
- pause between attempts: `3s`

## Migration naming

Migration files use standard goose format: `YYYYMMDDHHMMSS_name.sql`.
