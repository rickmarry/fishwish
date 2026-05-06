# ADR 0001 — Weather Service Design

**Date:** 2026-05-06
**Status:** Accepted
**Related:** `docs/backlog.md` B-003, `docs/b003-design.md`, `docs/b003-spec.md`

---

## Context

The weather-service needs to provide real-time weather and forecasts at fishing spots. We need to decide:
1. Which weather API to use
2. Whether to use a database or go stateless
3. How to cache responses

## Options Considered

### Weather API Options

- **Open-Meteo** — free, no registration, no key, WMO standard weather codes, 15 min update frequency
- **OpenWeatherMap** — free tier available, requires API key registration, 60 calls/minute limit
- **WeatherAPI.com** — free tier (1M calls/month), requires API key, more detailed data

**Decision:** Open-Meteo

**Why:** No registration friction, no API key to manage, no rate limit concerns for local dev. WMO standard codes are well-documented. The data quality is sufficient for fishing conditions (temp, wind, precipitation, pressure).

### Database vs Stateless

- **Postgres** — persistent storage, historical weather data possible, but adds complexity and cost
- **Stateless + Redis only** — no DB, cache weather responses, simpler architecture

**Decision:** Stateless with Redis caching only

**Why:** Weather is ephemeral data — storing it in Postgres adds no value since we always want current conditions. Other services (spot, user, social) already use Postgres; keeping weather-service without a DB is a deliberate architectural choice that demonstrates stateless service patterns.

### Cache Strategy

- **No cache** — every request hits Open-Meteo API, slow (network latency), wasteful for repeated spot views
- **In-memory cache** — fast, but lost on service restart, can't share across service instances
- **Redis cache** — shared across instances, TTL support, already used by other services

**Decision:** Redis with 30-minute TTL

**Why:** Weather doesn't change instantly. 30 minutes balances freshness (fishermen need current conditions) with API efficiency. Redis is already in the stack for user-service and spot-service.

---

## Decision

1. Use **Open-Meteo API** for weather data (`https://api.open-meteo.com/v1/forecast`)
2. Build weather-service as a **stateless service** with no database
3. Use **Redis** for caching weather responses with **1800s TTL** (30 minutes)
4. Cache key format: `weather:{lat}:{lng}` with coordinates rounded to 3 decimal places (~110m precision)
5. Round coordinates to avoid cache fragmentation for nearby spots

---

## Consequences

**What becomes easier:**
- No database migrations or schema management for weather-service
- Horizontal scaling is trivial (no state, shared Redis cache)
- No API key management or registration overhead
- Open-Meteo requires no auth, reducing failure points

**What becomes harder:**
- No historical weather data (would need DB to store it)
- Redis dependency — if Redis is down, we degrade to direct API calls (handled gracefully)
- Cache warming not possible — first request for a spot always hits the API

**Future implications:**
- If we need historical weather (e.g., "what were conditions when I caught this fish?"), we'll need to add Postgres or a time-series DB
- If Open-Meteo rate limits become an issue, we can add a backup weather API with fallback logic
- Multi-region deployment will need Redis replication or separate Redis per region
