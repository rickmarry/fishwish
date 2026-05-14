# Local Development & Debugging

## Killing Services by Port

When you want to debug a service in IntelliJ, kill just that service's port first:

```bash
lsof -ti :8081 | xargs kill  # user-service
lsof -ti :8082 | xargs kill  # spot-service
lsof -ti :8083 | xargs kill  # search-service
lsof -ti :8084 | xargs kill  # weather-service
lsof -ti :8085 | xargs kill  # social-service
```

The `-t` flag returns just the PID, safe to pipe to `xargs kill`.

## Kill All Services

```bash
for port in 8081 8082 8083 8084 8085; do lsof -ti :$port | xargs kill 2>/dev/null; done
```

## IntelliJ Run Configurations

Run configs are in `.run/` directory. Each config:
- Sets working directory to `services/<name>/`
- Loads env vars from `.env.local`
- Runs `cmd/server/main.go`

To debug: kill the service port above, then start that run config in IntelliJ in debug mode.

## Service Ports

| Service | Port | Has DB | Has Redis | Has MinIO |
|---------|------|--------|-----------|-----------|
| user-service | 8081 | ✓ | ✓ | |
| spot-service | 8082 | ✓ | ✓ | |
| search-service | 8083 | ✓ | | |
| weather-service | 8084 | | ✓ | |
| social-service | 8085 | ✓ | ✓ | ✓ |

## Frontend

```bash
cd frontend && npm run dev
```

Vite dev server on port 3006. Proxies `/api/*` to service ports.

### Depth Chart / Bathymetry Setup

The map has a toggleable depth chart overlay. Two providers are available:

**Default — GEBCO WMS** (free, no key needed):
- Global bathymetry via GEBCO_2024 grid (15 arc-second resolution)
- Works out of the box — just restart the frontend

**Alternative — VectorCharts.com** (requires API key):
- Higher-resolution NOAA ENC vector tiles with styled contours
- To use: set these in `frontend/.env.local`:
  ```
  VITE_DEPTH_PROVIDER=vectorcharts
  VITE_VECTORCHARTS_KEY=pk_your_key_here
  ```
- Get a key at https://vectorcharts.com

Without a valid VectorCharts key, the provider automatically falls back to GEBCO.

## Infrastructure

```bash
make up          # start postgres, redis, localstack, minio
make dev         # infra + all services with hot reload (air)
make dev-docker  # everything in containers
```
