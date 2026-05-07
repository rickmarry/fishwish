# FishFinder AGENTS.md

## Project

Fishing spot discovery platform. Go microservices + React frontend. Local-first dev with Docker Compose.

## Dev Commands

```
make setup     # first run: install air, migrate, npm deps
make up        # start infra only (postgres, redis, localstack, minio)
make dev       # infra + all services with hot reload (requires `air` installed)
make dev-docker # everything in containers (no local Go needed)
make migrate   # run migrations up
make migrate-down
make seed      # load 10 sample spots + 15 species
make test      # go test ./services/... -v -count=1
```

**Depth charts:** Set `VITE_VECTORCHARTS_KEY` in `frontend/.env.local` (get a free key at vectorcharts.com) to enable the bathymetry overlay toggle on the map.

## Architecture

**5 Go microservices**, each with its own `go.mod` using `replace fishwish => ../../` to import shared code from `pkg/`.

| Service | Port | Depends on |
|---------|------|------------|
| user-service | 8081 | postgres, redis |
| spot-service | 8082 | postgres, redis |
| search-service | 8083 | postgres |
| weather-service | 8084 | redis only (no DB, calls Open-Meteo API) |
| social-service | 8085 | postgres, minio |

**Frontend**: React 18 + Vite 5 + Tailwind + Zustand. Dev server proxies `/api/*` to service ports (docker hostnames in compose, `localhost` for local dev).

## Service Structure (standardized)

```
services/<name>/
  cmd/server/main.go        # entrypoint, chi router setup
  cmd/migrate/main.go       # spot-service only — migration runner
  internal/
    config/config.go        # loads .env.local via godotenv
    model/<name>.go         # request/response structs
    repository/db.go        # pgxpool DB struct + repo methods
    service/<name>.go       # business logic
    handler/<name>.go       # HTTP handlers, writeJSON helper
  .air.toml                 # hot reload config
  .env.local                # env vars for local dev
  Dockerfile                # multi-stage alpine build
  go.mod                    # replace fishwish => ../../
```

## Key Gotchas

- **Weather service has no PostgreSQL** — only Redis for caching. It calls Open-Meteo API (free, no key needed).
- **Migrations live in spot-service** even though they create tables for all services (users, spots, species, reviews, catch_logs). Run from `services/spot-service`.
- **Each service has its own go.mod** — do NOT try to run `go build` from the root for a service. Always `cd services/<name>` first or use `make test` which handles paths.
- **Root go.mod exists but is for shared pkg only** — services use `replace` directive.
- **Auth middleware in pkg/auth/ is stubbed** — JWT validation works but context key injection (`WithUserID`) returns `interface{}` instead of `context.Context`. Fix before production.
- **Vite proxy uses docker service names** in compose (`user-service:8081`), but `VITE_*` env vars point to `localhost`. Docker dev goes through the proxy; local `make dev` uses direct ports.

## Database

- PostgreSQL 16 + PostGIS 3.4. Spots use `GEOMETRY(Point, 4326)` with GIST index.
- Migrations use golang-migrate, file-based in `spot-service/internal/repository/migrations/`.
- LocalStack provides DynamoDB + S3 (for production parity; not heavily used yet).
- MinIO provides S3-compatible storage for user photos.

## Seeded Data

`scripts/seed/main.go` loads 10 real fishing spots (Tahoe, Everglades, Kenai River, etc.) and 15 species with spot_species join table records.
