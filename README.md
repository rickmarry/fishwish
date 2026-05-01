# FishFinder

A fishing spot discovery platform built with Go microservices, React, and AWS.

## Quick Start

```bash
# Install dependencies and set up local environment
make setup

# Start all services (infra + hot-reload dev)
make dev

# Or run everything in Docker
make dev-docker
```

## Architecture

- **Backend**: Go microservices (user, spot, search, weather, social)
- **Frontend**: React + Vite + MapLibre
- **Local Infra**: PostgreSQL + PostGIS, Redis, LocalStack, MinIO
- **Production Infra**: AWS (ECS, RDS, DynamoDB, S3, Cognito, OpenSearch)

## Services

| Service | Port | Description |
|---------|------|-------------|
| user-service | 8081 | Authentication, profiles, preferences |
| spot-service | 8082 | Fishing spots CRUD, geo queries |
| search-service | 8083 | Full-text search, filtering, ranking |
| weather-service | 8084 | Weather, tides, conditions |
| social-service | 8085 | Reviews, ratings, catch logs |

## Local Dev

```bash
make up        # Start infrastructure only
make migrate   # Run database migrations
make seed      # Load sample fishing spots
make test      # Run all tests
make down      # Stop everything
```

## Tech Stack

- Go 1.22+, chi router, golang-migrate
- React 18, Vite, Zustand, MapLibre GL
- Docker Compose, LocalStack, MinIO
- PostgreSQL + PostGIS, Redis
