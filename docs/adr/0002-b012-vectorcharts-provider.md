# ADR 0002 — Choose VectorCharts.com for Bathymetry Data

**Date:** 2026-05-07
**Status:** Accepted
**Related:** `docs/b012-design.md`, `docs/b012-spec.md`

---

## Context

FishFinder needs depth chart / bathymetry data to show anglers underwater terrain. Options range from free global datasets (GEBCO) to premium marine chart APIs (Navionics, MarineCharts). The chosen provider must balance data quality, cost, and integration effort for a small-scale fishing spot app.

## Options Considered

- **VectorCharts.com** — NOAA ENC data delivered as MapLibre-compatible vector tiles. Free tier: 25k requests/month, then $1/1k. US + international coverage. Simple API key auth.

- **GEBCO WMS** — Free global bathymetry at 450m grid resolution. Lower detail, not suitable for identifying fishing structure (drop-offs, holes). WMS raster tiles, harder to style.

- **MarineCharts.io** — $49/month starter plan. Same NOAA ENC data as VectorCharts. Better for production scale but overkill for V1 budget.

- **Navionics Web API (Standard)** — Free but restricted: cannot overlay own content on the chart. Defeats the purpose for FishFinder (need to show spot markers on depth data).

- **Navionics Web API (Enhanced)** — Full-featured but requires contacting Garmin for custom pricing. Likely thousands/year. Unclear timeline to get access.

- **Mapbox Bathymetry v2** — Uses GEBCO data through Mapbox. Requires Mapbox account and GL JS. Adds another platform dependency.

## Decision

Use **VectorCharts.com** for V1.

## Consequences

**Easier:**
- Free to start (25k req/mo covers ~500 user sessions)
- Drops in as a MapLibre vector tile source — no backend changes
- NOAA ENC data is authoritative for US waters
- Upgrade path is usage-based (no contract renegotiation)

**Harder:**
- API key exposed in frontend (acceptable for client-side tile services)
- US coverage is good; inland lake coverage may vary
- Will need to revisit if app scales past 25k requests/month
- Future migration to a different provider means updating the tile source URL and possible layer config changes
