# FishFinder

## What This Is

Fishing spot discovery platform. Built by Rick Marry — Principal/Staff Software Engineer with 20+ years in distributed systems, search infrastructure, and cloud-native platforms.

Full architecture: `docs/architecture.md`

---

## Architecture at a Glance

5 Go microservices + React frontend. Local-first dev with Docker Compose.

```
React Frontend (Vite + Tailwind + Zustand)
  │  HTTP /api/*
Go API Services:
  user-service:8081     Postgres + Redis
  spot-service:8082     Postgres + Redis (owns migrations)
  search-service:8083   Postgres (PostGIS)
  weather-service:8084  Redis only (calls Open-Meteo API)
  social-service:8085   Postgres + MinIO
```

**Communication patterns:**
- Frontend → services: HTTP REST via Vite proxy
- Services are independent — no inter-service calls yet (by design)

---

## Build Order

Services are built in this sequence — each depends on the previous. **Current status lives in `memory/project_current_work.md` (repo root), not here.**

| # | Service | Notes |
|---|---|---|
| 1 | `spot-service` | Core data model, migrations, seed data |
| 2 | `user-service` | Auth, user profiles, Redis caching |
| 3 | `search-service` | PostGIS spatial queries |
| 4 | `weather-service` | Open-Meteo API, Redis caching |
| 5 | `social-service` | Reviews, catch logs, MinIO photos |
| 6 | `frontend` | React SPA |

---

## Design Discipline — Active Mandate

We follow a spec-driven development stack to prevent intent drift. **Do not skip any layer.**

```
HLD → LLD → Specs/Contracts → Implementation → Tests enforce contracts
```

**Status as of 2026-05-06:**
- HLD — in progress (`docs/architecture.md`)
- ADRs — establish as needed (`docs/adr/`), write for every non-obvious system-level decision
- Infrastructure — Postgres, Redis, MinIO running locally via Docker Compose
- Migrations + seed data — applied (15 species, 10 spots)
- Services — spot-service and user-service partially implemented
- LLD/SPEC.md — not yet written; discipline applies for any new services or significant changes

---

## Definition of Done

Before declaring any implementation task complete, verify every applicable item:

**Design artifacts**
- [ ] ADR written if a non-obvious system-level decision was made (`docs/adr/`) — included in the same PR, never deferred
- [ ] `DESIGN.md` written for new services or significant features
- [ ] `SPEC.md` (EARS format) written for new services or significant features

**Contracts**
- [ ] API request/response structs defined in `internal/model/`
- [ ] OpenAPI spec updated if exposing public endpoints

**Implementation**
- [ ] Migration written if schema changed (lives in spot-service)
- [ ] `docker-compose.yml` updated if new service added
- [ ] Lint and typecheck pass (`make test`)

**Documentation**
- [ ] `services/<name>/docs/service.md` created or updated

**Git workflow**
- [ ] GitHub issue created
- [ ] Feature branch created from main
- [ ] Conventional commit message, no Co-Authored-By line
- [ ] PR created referencing the issue with `Closes #N`
- [ ] PR merged — `gh pr merge <number> --squash --delete-branch`
- [ ] `git checkout main && git pull`
- [ ] **Update `docs/backlog.md` Build Order** — mark the completed item, remove it from "Next up", promote the next item. Commit alongside `memory/project_current_work.md` in the same turn. These two files must always be in sync.

---

## Shared AI Assistant Rules

**IMPORTANT:** Before starting any work, read and follow the rules in `docs/AI_RULES.md`.

That file contains the complete set of workflow rules, conventions, and processes that apply to **all AI assistants** (Claude, Gemini, Junie, etc.) working on this project. Key rules include:
- Git workflow (issue → branch → commit → push → PR → merge → checkout main)
- Commit message conventions (conventional commits, no co-author lines)
- Design discipline (HLD → LLD → Specs → Implementation → Tests)
- ADR process for architecture decisions

This section exists to provide Claude-specific context below, but the shared rules in `docs/AI_RULES.md` take precedence for workflow and conventions.

---

## Conventions

### Every service gets an architecture doc
`services/<name>/docs/service.md` — purpose, design decisions, API, config, local dev instructions.
See `services/spot-service/docs/service.md` as the reference example (create it first).

### Service structure (standardized)
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

### Migrations live in spot-service
Even though they create tables for all services (users, spots, species, reviews, catch_logs). Run from `services/spot-service`.

### Each service has its own go.mod
Do NOT try to run `go build` from the root for a service. Always `cd services/<name>` first or use `make test` which handles paths.

### Auth middleware is stubbed
JWT validation works but context key injection (`WithUserID`) returns `interface{}` instead of `context.Context`. Fix before production.

---

## Key Design Decisions (top-level)

**Why a monorepo?**
Single repo for all services — easier local development, shared `pkg/`, unified Makefile. Each service still has its own go.mod using `replace fishwish => ../../`.

**Why Go for all services?**
Rick has deep Go experience. Go's concurrency model and standard library are well-suited for REST APIs with pgxpool and Redis.

**Why PostGIS for search?**
Spatial queries (find spots within X miles of a point) are first-class in PostGIS. No need for a separate search engine yet.

**Why weather-service has no DB?**
It calls Open-Meteo API (free, no key needed) and caches in Redis. Stateless by design.

**Why MinIO for social-service?**
S3-compatible storage for user photos. MinIO gives local parity with production AWS S3.

---

## Local Development

```bash
make setup     # first run: install air, migrate, npm deps
make up        # start infra only (postgres, redis, localstack, minio)
make dev       # infra + all services with hot reload
make migrate   # run migrations up
make seed      # load 10 sample spots + 15 species
make test      # go test ./services/... -v -count=1
```

---

## Session Commands

### rs (resume session)
Do ALL of the following before giving any summary:

1. **Read** `memory/project_current_work.md` (repo root) — current status, next steps
2. **Read** `docs/backlog.md` — full inventory of all backlog items and their status
3. Give a concise summary: what's in progress, and the list of all **Design Complete** items (ready to build)

### ss (save session)
After a PR merges or at end of session:

1. Update `memory/project_current_work.md` — mark completed items, update next steps
2. Update `docs/backlog.md` Build Order — remove completed item, promote next. These two files must always be in sync.
3. Commit both: `git add memory/project_current_work.md docs/backlog.md && git commit -m "docs: update project status - <what changed>"`

---

## Context for Claude

- Rick is a Principal/Staff Software Engineer — don't over-explain Go, distributed systems, or cloud patterns
- This project runs in parallel with an active job search — keep code quality high
- Every service needs a `docs/service.md` — create it as part of building the service, not after
- Prefer explicit, readable Go over clever Go — interview reviewability matters
- **The SDLC is Rick's own — frame it that way in interviews.** He didn't read about HLD → LLD → EARS specs → enforced contracts → implementation. He ran into intent drift building a multi-service system with agentic tooling, diagnosed the failure mode himself, and built a process to prevent it. The discipline predates any blog post he's read. When this comes up in interviews, the framing is: "this is how I work, not a methodology I adopted."
- **Enforce the design discipline.** If Rick moves toward implementation without a DESIGN.md and SPEC.md for the service, flag it.
