# FishFinder — Active Backlog

Active items only — `Backlog` and `Design Complete` status. Completed items are in `docs/backlog-complete.md`.

**Status vocabulary:** `Backlog` → `Design Complete` (LLD + SPEC written) → `Implementation Complete` (moved to `backlog-complete.md`)

---

## Build Order

**Any Design Complete item is ready to build — Rick picks whatever he wants.** The sequence below is advisory: a suggested order when there's no strong preference, based on dependencies. It is not a constraint.

### Recommended sequence
1. **B-001 — User Registration & Auth** — Backlog. Foundation for all user-specific features.
2. **B-002 — Spot Search & Discovery** — Backlog. Core feature: search spots by location, species, amenities.
3. **B-004 — Spot Reviews & Ratings** — Backlog. Depends on social-service (Implementation Partial).

### Dependency map
- B-002 (Spot Search): depends on spot-service (Implementation Complete)
- B-004 (Reviews): depends on social-service (Implementation Partial)

### Parking lot (genuinely unsequenced)
- B-005 — Catch Log & Personal Best Tracking
- B-006 — Species Guide with Seasonal Patterns
- B-007 — Spot Photos & User Generated Content
- B-008 — Fishing Buddies & Trip Planning
- B-009 — Tide Charts & Solunar Tables
- B-010 — Mobile App (React Native)
- B-012 — Depth Charts / Bathymetry Overlay — Design Complete

---

## Sequenced Backlog — Needs Design

### B-001 — User Registration & Auth
**Status:** Backlog

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

**User story:** As a user, I want to search for fishing spots near me (or near a destination) so I can find new places to fish.

**Core loop:**
1. User enters search criteria: location (lat/lng or place name), radius, species filter, amenities
2. Frontend calls `search-service` with criteria
3. Search service queries PostGIS (`ST_DWithin` + `ST_Distance`) against spots table
4. Results returned with distance, species match info, and available amenities
5. User clicks a spot → detail view

**Design decisions to make:**
- Geocoding for place name → lat/lng (use Open-Meteo geocoding API, free)
- Pagination strategy for results
- How to surface species match (spot_species join)

**Suggested architecture:**
- `search-service`: `POST /search` with `SearchRequest{lat, lng, radius_miles, species[], amenities[]}` → `SearchResponse{spots[], total}`
- PostGIS query using `ST_DWithin(geom, ST_MakePoint($1,$2)::geometry, $3)` with GIST index
- Frontend: search form, results list with distance, map view (future)

---

### B-004 — Spot Reviews & Ratings
**Status:** Backlog
**Depends on:** social-service (Implementation Partial)

**User story:** As a user, I want to rate and review fishing spots so I can share my experience and help others decide if a spot is worth visiting.

**Core loop:**
1. User views a spot → clicks "Write Review"
2. Modal: rating (1-5 stars), text review, optional photo upload
3. Frontend calls `social-service`: `POST /spots/{id}/reviews`
4. Review stored in `reviews` table
5. Reviews displayed on spot detail page

**Design decisions to make:**
- Photo upload flow: frontend → social-service → MinIO, or presigned URL direct to MinIO?
- Review editing/deletion (only by author)
- Aggregate rating calculation (stored vs. computed on read)

**Suggested architecture:**
- `social-service`: review CRUD endpoints, photo upload to MinIO
- `reviews` table: `id, spot_id, user_id, rating, content, photo_urls (JSONB), created_at`
- Frontend: review form, star rating component, photo upload, reviews list on spot detail

---

## Unsequenced Active Backlog

### B-005 — Catch Log & Personal Best Tracking
**Status:** Backlog

**User story:** As an angler, I want to log my catches so I can track my personal bests and keep a fishing journal.

**Core loop:**
1. User clicks "Log a Catch"
2. Form: species, weight, length, photo, date, spot
3. Submit → `social-service`: `POST /catches`
4. "Personal Bests" view shows top catch per species

**Suggested architecture:**
- `catch_logs` table: `id, user_id, spot_id, species_id, weight, length, photo_url, caught_at`
- `social-service`: catch CRUD, personal bests query (`SELECT DISTINCT ON (species_id) ... ORDER BY weight DESC`)
- Frontend: catch log form, history list, personal bests page

---

### B-006 — Species Guide with Seasonal Patterns
**Status:** Backlog

**User story:** As a user, I want to browse species and learn what's biting and where so I can plan my fishing trips better.

**Suggested architecture:**
- Extend `species` table: `description, preferred_bait, seasonal_pattern (JSONB with month→activity map)`
- Frontend: species grid, detail page, seasonal activity chart (recharts)

---

### B-007 — Spot Photos & User Generated Content
**Status:** Backlog

**User story:** As a user, I want to upload and view photos of fishing spots so I can see what the spot looks like before visiting.

**Suggested architecture:**
- `spot_photos` table: `id, spot_id, user_id, url, caption, taken_at`
- Extend `social-service` photo handling (shared with reviews)
- Frontend: photo gallery component, upload modal, lightbox view

---

### B-008 — Fishing Buddies & Trip Planning
**Status:** Backlog (stub — needs full entry)

Connect with other anglers, plan trips together, share spots privately.

---

### B-009 — Tide Charts & Solunar Tables
**Status:** Backlog (stub — needs full entry)

Integrate tide data and solunar forecasts (moon phase, sunrise/sunset) for saltwater and coastal spots.

---

### B-010 — Mobile App (React Native)
**Status:** Backlog (stub — needs full entry)

Wrap the existing React app in React Native or build a native mobile app for on-the-water access.

---

### B-012 — Depth Charts / Bathymetry Overlay
**Status:** Design Complete

**User story:** As an angler, I want to see depth contours and underwater terrain on the map so I can find drop-offs, holes, and structure where fish are likely to hold.

**Core loop:**
1. User opens the map
2. Depth chart layer renders as an overlay on the existing MapLibre map
3. Depth contours, soundings, and shaded bathymetry are visible beneath spot markers
4. User can toggle the depth layer on/off

**Suggested architecture:**
- Add a vector tile source + layer to the existing MapLibre map style in `MapView.jsx`
- Toggle control in the UI (button or layer switcher)
- No backend changes needed if using a tile service

**Options researched:**

| Provider | Pricing | Coverage | Resolution | Integration |
|---|---|---|---|---|
| VectorCharts.com | 25k req/mo **free**, then $1/1k | US + international | Good (NOAA ENC vector tiles) | MapLibre vector tile source |
| GEBCO WMS | Free | Global | 450m grid (coarse) | WMS tile layer in MapLibre |
| MarineCharts.io | $49/mo starter | US coastal + inland | Good (NOAA ENC) | MapLibre vector tile source |
| Navionics Web API | Free (Standard) or Paid (Enhanced) | Global | HD SonarChart | JS library (OpenLayers-based) |
| Mapbox Bathymetry v2 | Free tier with Mapbox account | Global (GEBCO data) | GEBCO-derived | Mapbox vector tiles |

**Recommendation for V1:** VectorCharts.com has a generous free tier, proper US NOAA ENC data with depth contours, and drops in as a MapLibre vector tile source — minimal code change.

**Caveats:**
- Inland lake coverage varies significantly between providers (GEBCO is poor, Navionics is best)
- Free tiers often require attribution
- If user base grows, 25k requests/month on VectorCharts may be exceeded quickly
- Navionics Standard (free) does NOT allow overlaying your own content — defeats the purpose
- Consider making this a toggleable layer (performance impact of rendering both OSM raster tiles and vector depth tiles)


