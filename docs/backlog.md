# FishFinder — Product Backlog

This is the single source of truth for FishFinder feature ideas and future work. Items are roughly grouped by theme, not prioritized within sections.

**Status vocabulary:** `Backlog` → `Design Complete` (LLD + SPEC written) → `Implementation Partial` (built but has pending requirements in spec) → `Implementation Complete` (all spec requirements satisfied, merged)

---

## Build Order

**Any Design Complete item is ready to build — Rick picks whatever he wants.** The sequence below is advisory: a suggested order when there's no strong preference, based on dependencies. It is not a constraint.

Updated whenever new items are added or dependencies change.

### Next up (in order)
1. **B-001 — User Registration & Auth** — Backlog. Foundation for all user-specific features.
2. **B-002 — Spot Search & Discovery** — Backlog. Core feature: search spots by location, species, amenities.
3. **B-003 — Weather Integration** — Backlog. Real-time weather at fishing spots via Open-Meteo API.

### Dependency map
- B-002 (Spot Search): depends on spot-service (Implementation Complete)
- B-003 (Weather): depends on weather-service (Implementation Complete)
- B-004 (Reviews): depends on social-service (Implementation Partial)

### Retired IDs (assigned but no backlog entry)
| ID | Notes |
|---|---|
| B-000 | Initial scaffold + migrations + seed data — Implementation Complete (commit fe64465) |

### Parking lot (no sequence yet)
- B-005 — Catch Log & Personal Best Tracking
- B-006 — Species Guide with Seasonal Patterns
- B-007 — Spot Photos & User Generated Content

---

## Core Infrastructure (Pre-Backlog Era — Retroactively Documented)

These items were built before the backlog format was established. Added retroactively so the backlog is a true inventory of the system.

### B-000 — Initial Scaffold + Infrastructure
**Status:** Implementation Complete (scaffold only)
**Docs:** `AGENTS.md`, `CLAUDE.md`

**What it is:** Initial project scaffold with 5 Go microservices (user, spot, search, weather, social), React frontend (Vite + Tailwind + Zustand), Docker Compose infrastructure (Postgres 16 + PostGIS, Redis, MinIO, LocalStack disabled), and database migrations with seed data (15 species, 10 fishing spots).

**What's actually done:**
- `docker-compose.yml` with Postgres, Redis, MinIO (LocalStack disabled)
- Migrations applied, seed data loaded (15 species, 10 spots)
- Service scaffolds created with standard structure (`cmd/server/main.go`, `internal/` dirs)
- Frontend scaffold with Vite dev server running

**What's NOT done (needs verification):**
- Whether each service's HTTP handlers are actually implemented
- Whether services start and respond to requests
- Whether frontend can fetch data from backend

**Key design decisions:**
- PostGIS for spatial queries (not a separate search engine)
- Weather service is stateless (no DB) — calls Open-Meteo API, caches in Redis
- Migrations live in spot-service even though they create tables for all services
- MinIO for local S3 parity with production
- LocalStack disabled (was causing auth issues, not heavily used yet)

---

## User Features

### B-001 — User Registration & Auth
**Status:** Backlog

User registration with JWT-based auth. Users can create accounts, log in, and have their identity propagated through the API via the auth middleware.

**User story:** As a fishing enthusiast, I want to create an account and log in so that I can save favorite spots, write reviews, and track my catches.

**Core loop:**
1. User registers with email + password
2. User logs in → receives JWT
3. JWT propagated via `Authorization: Bearer` header
4. Protected endpoints extract user ID from token context

**Design decisions to make:**
- Password hashing strategy (bcrypt recommended)
- JWT secret management (env var, rotate-able)
- Auth middleware context key fix (currently returns `interface{}` instead of `context.Context`)

**Suggested architecture:**
- `user-service`: registration + login endpoints, bcrypt hashing, JWT signing
- Auth middleware in `pkg/auth/` — fix `WithUserID` to return `context.Context`
- Frontend: register/login forms, token storage in Zustand store + `Authorization` header on API calls

---

### B-002 — Spot Search & Discovery
**Status:** Backlog

Search fishing spots by geographic location (within X miles of a point), species available, and amenities. Uses PostGIS spatial queries.

**User story:** As a user, I want to search for fishing spots near me (or near a destination) so I can find new places to fish.

**Core loop:**
1. User enters search criteria: location (lat/lng or place name), radius, species filter, amenities
2. Frontend calls `search-service` with criteria
3. Search service queries PostGIS (`ST_DWithin` + `ST_Distance`) against spots table
4. Results returned with distance, species match info, and available amenities
5. User clicks a spot → detail view with full info

**Design decisions to make:**
- Geocoding for place name → lat/lng (use Open-Meteo geocoding API, free)
- Pagination strategy for results
- How to surface species match (spot_species join)

**Suggested architecture:**
- `search-service`: `POST /search` with `SearchRequest{lat, lng, radius_miles, species[], amenities[]}` → `SearchResponse{spots[], total}`
- PostGIS query using `ST_DWithin(geom, ST_MakePoint($1,$2)::geometry, $3)` with GIST index
- Frontend: search form, results list with distance, map view (future)

---

### B-003 — Weather Integration
**Status:** Implementation Complete (2026-05-06)
**Docs:** `docs/b003-design.md`, `docs/b003-spec.md`, `docs/adr/0001-b003-weather-service-design.md`

**What it is:** Real-time weather and 3-day forecasts at fishing spots via Open-Meteo API (free, no key). Stateless service with Redis caching (30-minute TTL). Coordinates rounded to 3 decimal places for cache efficiency.

**What's done:**
- `GET /weather?lat=X&lng=Y` returns current conditions + 3-day forecast
- Redis caching with key `weather:{lat}:{lng}`, TTL 1800s
- WMO weather codes mapped to human-readable descriptions
- Input validation (lat: -90 to 90, lng: -180 to 180)
- Error handling: 400 for bad params, 502 for API failure, 504 for timeout
- Graceful Redis degradation (skip cache if Redis down, call API directly)
- Open-Meteo API client in `pkg/weather/client.go` (reused)

**Verification:**
- `curl http://localhost:8084/weather?lat=39.0968&lng=-120.0324` → 200 with current + forecast
- Cache hit on 2nd call → `cached: true`
- Invalid lat → 400 Bad Request

**Key design decisions:**
- Stateless service (no DB) — weather is ephemeral data
- Open-Meteo over OpenWeatherMap — no API key, no registration
- Redis cache with 30min TTL — balances freshness vs. API calls
- Coordinate rounding to 3 decimals — ~110m precision, avoids cache fragmentation

---

## Social Features

### B-004 — Spot Reviews & Ratings
**Status:** Backlog

Users can rate and review fishing spots. Reviews include rating (1-5), text, and photos (stored in MinIO via social-service).

**User story:** As a user, I want to rate and review fishing spots so I can share my experience and help others decide if a spot is worth visiting.

**Core loop:**
1. User views a spot → clicks "Write Review"
2. Modal: rating (1-5 stars), text review, optional photo upload
3. Frontend calls `social-service`: `POST /spots/{id}/reviews`
4. Review stored in `reviews` table (spot_id, user_id, rating, content, photo_urls[])
5. Reviews displayed on spot detail page with reviewer info and date

**Design decisions to make:**
- Photo upload flow: frontend → social-service → MinIO, or presigned URL direct to MinIO?
- Review editing/deletion (only by author)
- Aggregate rating calculation (stored vs. computed on read)

**Suggested architecture:**
- `social-service`: review CRUD endpoints, photo upload to MinIO
- `reviews` table: `id, spot_id, user_id, rating, content, photo_urls (JSONB), created_at`
- Frontend: review form, star rating component, photo upload, reviews list on spot detail

---

### B-005 — Catch Log & Personal Best Tracking
**Status:** Backlog

Users log their catches: species, weight, length, photo, date, and linked spot. Tracks personal bests per species.

**User story:** As an angler, I want to log my catches so I can track my personal bests and keep a fishing journal.

**Core loop:**
1. User clicks "Log a Catch"
2. Form: species (from seeded list), weight, length, photo, date, spot (search or select)
3. Submit → `social-service`: `POST /catches`
4. Catch stored in `catch_logs` table
5. User can view their catch history and filter by species
6. "Personal Bests" view shows top catch per species

**Suggested architecture:**
- `catch_logs` table: `id, user_id, spot_id, species_id, weight, length, photo_url, caught_at`
- `social-service`: catch CRUD, personal bests query (`SELECT DISTINCT ON (species_id) ... ORDER BY weight DESC`)
- Frontend: catch log form, history list, personal bests page

---

## Content & Discovery

### B-006 — Species Guide with Seasonal Patterns
**Status:** Backlog

Interactive species guide using the 15 seeded species. Shows seasonal patterns (when they're active), preferred bait/lures, and which spots have that species.

**User story:** As a user, I want to browse species and learn what's biting and where so I can plan my fishing trips better.

**Core loop:**
1. User navigates to "Species" page
2. Grid of species cards with image, name, seasonal activity chart
3. Click species → detail page: description, seasonal pattern, preferred bait, spots where it's found
4. Link to search with that species pre-filtered

**Suggested architecture:**
- Extend `species` table with: `description, preferred_bait, seasonal_pattern (JSONB with month→activity map)`
- `search-service` or `spot-service`: endpoint to get spots by species
- Frontend: species grid, detail page, seasonal activity chart (recharts)

---

### B-007 — Spot Photos & User Generated Content
**Status:** Backlog

Users can upload photos of spots (not just review photos). Photo gallery per spot with captions, taken date, and map location within the spot.

**User story:** As a user, I want to upload and view photos of fishing spots so I can see what the spot looks like before visiting.

**Core loop:**
1. User views spot → "Photos" tab
2. Photo gallery grid (thumbnail view)
3. Click photo → full-size view with caption, date, uploader
4. "Add Photo" button → upload modal (MinIO via social-service)
5. Photos displayed on spot detail and in user's profile

**Suggested architecture:**
- Extend `social-service` photo handling (shared with reviews)
- `spot_photos` table: `id, spot_id, user_id, url, caption, taken_at, location (optional geom)`
- Frontend: photo gallery component, upload modal, lightbox view

---

## Future / V2

### B-008 — Fishing Buddies & Trip Planning
**Status:** Backlog

Connect with other anglers, plan trips together, share spots privately.

### B-009 — Tide Charts & Solunar Tables
**Status:** Backlog

Integrate tide data and solunar forecasts (moon phase, sunrise/sunset) for saltwater and coastal spots.

### B-010 — Mobile App (React Native)
**Status:** Backlog

Wrap the existing React app in React Native or build native mobile app for on-the-water access.
