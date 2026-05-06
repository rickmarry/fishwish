# B-003 — Weather Integration (LLD)

**Date:** 2026-05-06
**Status:** Design Complete
**Related:** `docs/backlog.md` B-003, `docs/adr/0001-b003-weather-service-design.md`

---

## What This Is

Real-time weather and forecasts at fishing spots via Open-Meteo API (free, no registration, no key). Stateless service — no database. Redis for caching responses (30-minute TTL).

---

## Why This Matters

Fishermen check weather before heading out — wind, precipitation, pressure, and temperature all affect bite activity. A spot detail page without weather is incomplete. This is a core feature, not a nice-to-have.

---

## Service Architecture

```
Frontend → weather-service:8084 → Redis (cache)
                     ↓ (cache miss)
                  Open-Meteo API
```

**weather-service** is stateless — no Postgres, no persistent storage. Redis is cache-only.

---

## Data Model

### Open-Meteo API

**Endpoint:** `https://api.open-meteo.com/v1/forecast`

**Request params:**
```
latitude, longitude
current=temperature_2m,relative_humidity_2m,precipitation,weather_code,wind_speed_10m,wind_direction_10m,pressure_msl
daily=temperature_2m_max,temperature_2m_min,precipitation_sum,wind_speed_10m_max,weather_code
timezone=auto
forecast_days=3
```

**Response shape (simplified):**
```json
{
  "current": {
    "temperature_2m": 18.5,
    "relative_humidity_2m": 72,
    "precipitation": 0.0,
    "weather_code": 2,
    "wind_speed_10m": 12.5,
    "wind_direction_10m": 270,
    "pressure_msl": 1013.2
  },
  "daily": {
    "time": ["2026-05-06", "2026-05-07", "2026-05-08"],
    "temperature_2m_max": [20.1, 18.5, 19.8],
    "temperature_2m_min": [12.3, 11.8, 13.1],
    "precipitation_sum": [0.0, 2.5, 0.5],
    "wind_speed_10m_max": [15.2, 18.7, 12.3],
    "weather_code": [2, 61, 2]
  }
}
```

### Internal Model (`internal/model/weather.go`)

```go
type WeatherResponse struct {
    Current     CurrentWeather  `json:"current"`
    Forecast    []DailyForecast `json:"forecast"`
    Cached      bool            `json:"cached"`
}

type CurrentWeather struct {
    TemperatureC     float64 `json:"temperature_c"`
    HumidityPercent  int     `json:"humidity_percent"`
    PrecipitationMM  float64 `json:"precipitation_mm"`
    WeatherCode      int     `json:"weather_code"`
    WindSpeedKPH     float64 `json:"wind_speed_kph"`
    WindDirectionDeg int     `json:"wind_direction_deg"`
    PressureHPA      float64 `json:"pressure_hpa"`
}

type DailyForecast struct {
    Date             string  `json:"date"`
    TempMaxC         float64 `json:"temp_max_c"`
    TempMinC         float64 `json:"temp_min_c"`
    PrecipitationMM  float64 `json:"precipitation_mm"`
    WindSpeedMaxKPH  float64 `json:"wind_speed_max_kph"`
    WeatherCode      int     `json:"weather_code"`
}
```

### Weather Codes (WMO standard, simplified)
| Code Range | Condition |
|---|---|
| 0 | Clear sky |
| 1-3 | Mainly clear, partly cloudy, overcast |
| 45, 48 | Fog |
| 51-67 | Drizzle, rain, freezing rain |
| 71-77 | Snow, snow grains |
| 80-82 | Rain showers |
| 85-86 | Snow showers |
| 95-99 | Thunderstorm |

Helper function: `WeatherCodeToDescription(code int) string`

---

## API Endpoint

### GET /weather?lat={float}&lng={float}

**Request:** Query params `lat` (float, -90 to 90) and `lng` (float, -180 to 180)

**Response:** `WeatherResponse` JSON

**Errors:**
- `400 Bad Request` — missing or invalid lat/lng
- `502 Bad Gateway` — Open-Meteo API unreachable
- `504 Gateway Timeout` — Open-Meteo API timeout (>10s)

**Cache strategy:**
- Redis key: `weather:{lat}:{lng}` (rounded to 3 decimal places to avoid cache fragmentation)
- TTL: 1800 seconds (30 minutes)
- Cache miss → call API → store in Redis → return
- Cache hit → return cached data with `cached: true`
- Stale cache: not used. If Redis is down, skip cache and call API directly (degrade gracefully)

---

## Handler Flow

```
GET /weather?lat=39.0968&lng=-120.0324
  ↓
Validate lat/lng params
  ↓
Round lat/lng to 3 decimal places (precision: ~110m)
  ↓
Redis GET weather:{rounded_lat}:{rounded_lng}
  ↓
Cache hit? → Return cached JSON with cached=true
  ↓
Cache miss?
  ↓
HTTP GET to Open-Meteo API (10s timeout)
  ↓
Parse response → map to WeatherResponse
  ↓
Redis SETEX weather:{key} 1800 {json} (fire-and-forget, don't fail request if Redis write fails)
  ↓
Return WeatherResponse with cached=false
```

---

## Redis Cache

**Client:** `github.com/redis/go-redis/v9` (already in go.mod per AGENTS.md)

**Key format:** `weather:{lat}:{lng}` where lat/lng are rounded to 3 decimal places

**Why round:** Open-Meteo API returns identical data for coordinates within ~110m. Without rounding, `39.0968` and `39.0969` would be separate cache entries for the same spot.

**TTL:** 1800 seconds (30 minutes). Weather changes, but not instantly. 30 minutes balances freshness vs. API call frequency.

**Failure mode:** If Redis is unavailable, skip cache entirely and call API directly. Log the error but don't fail the request.

---

## Open-Meteo API Client

**Library:** Standard `net/http` (no SDK needed — simple REST API)

**Timeout:** 10 seconds (weather data is not critical path, fail fast)

**Error handling:**
- Non-200 response → return 502 to client
- Timeout → return 504 to client
- Network error → return 502 to client

**No API key needed** — Open-Meteo is free for non-commercial use, no registration required.

---

## Frontend Integration

**Call pattern:** When user views spot detail page, frontend calls:
```
GET /api/weather-service/weather?lat={spot.lat}&lng={spot.lng}
```
(Vite proxy routes `/api/weather-service` to `localhost:8084`)

**Weather widget:** Display on spot detail page:
- Current: temperature, condition icon, wind, pressure
- 3-day forecast: day, high/low, precipitation, condition icon
- "Cached" indicator if data is from cache (optional, transparency)

**Loading state:** Show skeleton while fetching
**Error state:** Show "Weather unavailable" with retry button

---

## Files Changed

```
services/weather-service/
  internal/
    model/weather.go        # WeatherResponse, CurrentWeather, DailyForecast, helper functions
    handler/weather.go      # GET /weather handler
    service/weather.go      # Open-Meteo API client, Redis cache logic
    config/config.go        # Add OPEN_METEO_BASE_URL, WEATHER_CACHE_TTL env vars (optional)
  cmd/server/main.go        # Register /weather route
```

No migrations needed (no database).
No proto files (REST API, not gRPC).
No changes to other services.

---

## Testing

**Unit tests:**
- `weather_test.go` — handler tests with httptest, mock Redis
- `service_test.go` — API client tests with httptest mock server
- `model_test.go` — WeatherCodeToDescription function tests

**Integration test (future):**
- Start service + Redis, call endpoint, verify response shape

**Manual test:**
```bash
curl "http://localhost:8084/weather?lat=39.0968&lng=-120.0324"
```

---

## Open Questions

None. Design is complete and ready for implementation.
