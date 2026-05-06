# B-003 — Weather Integration (SPEC)

**Date:** 2026-05-06
**Status:** Design Complete
**Related:** `docs/backlog.md` B-003, `docs/b003-design.md`

---

## Overview

The weather-service provides real-time weather and 3-day forecasts at fishing spots via Open-Meteo API with Redis caching. Stateless service with no database.

---

## Requirements (EARS Format)

### Functional Requirements

**[PENDING] F-001: Current Weather Retrieval**
WHEN a user views a fishing spot detail page, the system shall display current weather conditions (temperature, humidity, precipitation, wind speed/direction, pressure, weather condition) for the spot's coordinates.

**[PENDING] F-002: Forecast Display**
WHEN a user views a fishing spot detail page, the system shall display a 3-day weather forecast (high/low temp, precipitation, max wind, condition) for the spot's coordinates.

**[PENDING] F-003: Cache Hit Response**
WHEN the weather for a location is cached in Redis, the system shall return the cached data with `cached: true` and skip the Open-Meteo API call.

**[PENDING] F-004: Cache Miss → API Call**
WHEN the weather for a location is not cached, the system shall call the Open-Meteo API, cache the response, and return it with `cached: false`.

**[PENDING] F-005: Coordinate Rounding for Cache Keys**
WHEN generating a Redis cache key, the system shall round latitude and longitude to 3 decimal places to avoid cache fragmentation for nearby coordinates.

**[PENDING] F-006: Input Validation**
WHEN a request is made to `GET /weather`, the system shall return 400 Bad Request if `lat` or `lng` query params are missing or out of valid range (lat: -90 to 90, lng: -180 to 180).

**[PENDING] F-007: API Failure Handling**
WHEN the Open-Meteo API returns a non-200 response or times out (>10s), the system shall return 502 or 504 respectively to the client.

**[PENDING] F-008: Redis Down Graceful Degrade**
WHEN Redis is unavailable, the system shall skip the cache layer and call the Open-Meteo API directly, logging the error but not failing the request.

**[PENDING] F-009: Weather Code to Description Mapping**
WHEN returning weather data, the system shall include human-readable weather condition descriptions mapped from WMO weather codes.

### Non-Functional Requirements

**[PENDING] NF-001: Response Time**
WHILE the weather data is cached, the system shall respond within 50ms (Redis lookup + JSON serialization).

**[PENDING] NF-002: Cache TTL**
WHERE caching is enabled, the system shall set a TTL of 1800 seconds (30 minutes) on all weather cache entries.

**[PENDING] NF-003: API Timeout**
WHILE calling the Open-Meteo API, the system shall enforce a 10-second timeout to fail fast.

---

## API Contract

### Endpoint: GET /weather

**Request:**
```
GET /weather?lat={float}&lng={float}
```

| Param | Type | Required | Validation |
|---|---|---|---|
| lat | float | yes | -90.0 ≤ lat ≤ 90.0 |
| lng | float | yes | -180.0 ≤ lng ≤ 180.0 |

**Response (200 OK):**
```json
{
  "current": {
    "temperature_c": 18.5,
    "humidity_percent": 72,
    "precipitation_mm": 0.0,
    "weather_code": 2,
    "weather_description": "Partly cloudy",
    "wind_speed_kph": 12.5,
    "wind_direction_deg": 270,
    "pressure_hpa": 1013.2
  },
  "forecast": [
    {
      "date": "2026-05-06",
      "temp_max_c": 20.1,
      "temp_min_c": 12.3,
      "precipitation_mm": 0.0,
      "wind_speed_max_kph": 15.2,
      "weather_code": 2,
      "weather_description": "Partly cloudy"
    }
  ],
  "cached": false
}
```

**Error Responses:**
- `400 Bad Request` — `{"error": "invalid lat/lng parameters", "message": "..."}`
- `502 Bad Gateway` — `{"error": "weather API unavailable", "message": "..."}`
- `504 Gateway Timeout` — `{"error": "weather API timeout", "message": "..."}`

---

## Acceptance Criteria

The feature is **Complete** when ALL of the following are verified:

1. [ ] `GET /weather?lat=39.0968&lng=-120.0324` returns 200 with valid `WeatherResponse` JSON
2. [ ] Second call to same coordinates returns 200 with `cached: true`
3. [ ] `GET /weather?lat=invalid&lng=-120.0324` returns 400
4. [ ] `GET /weather?lat=999&lng=-120.0324` returns 400 (out of range)
5. [ ] Redis down → request still succeeds (calls API directly)
6. [ ] Response includes `weather_description` mapped from WMO code
7. [ ] Frontend weather widget renders on spot detail page with current + 3-day forecast
8. [ ] Unit tests pass: `go test ./services/weather-service/... -v -count=1`

---

## Frontend Integration Spec

**Component:** `WeatherWidget.tsx` (or similar)

**Props:** `lat: number`, `lng: number`

**Behavior:**
1. On mount or lat/lng change, call `GET /api/weather-service/weather?lat={lat}&lng={lng}`
2. Show skeleton loader while fetching
3. On success, render current weather + 3-day forecast
4. On error, show "Weather unavailable" with retry button
5. Show a small "cached" badge if `cached: true` (optional transparency feature)

**Location in UI:** Spot detail page, below spot description, above reviews section.

---

## Dependencies

- **weather-service scaffold** — exists in `services/weather-service/`
- **Redis** — running in Docker Compose (port 6379)
- **Open-Meteo API** — free, no key needed: `https://api.open-meteo.com/v1/forecast`

No dependencies on other FishWish services. Weather-service is independent.
