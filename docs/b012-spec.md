# B-012 — Depth Charts / Bathymetry Overlay (SPEC)

**Date:** 2026-05-07
**Status:** Design Complete
**Related:** `docs/backlog.md` B-012, `docs/b012-design.md`, `docs/adr/0002-b012-vectorcharts-provider.md`

---

## Overview

A toggleable bathymetry overlay on the fishing spot map, showing depth contours and shaded depth areas using NOAA ENC data via VectorCharts.com vector tiles.

---

## Requirements (EARS Format)

### Functional Requirements

**[PENDING] F-001: Depth Layer Toggle**
WHEN a user is viewing the map, the system shall provide a control to toggle the depth chart layer on and off.

**[PENDING] F-002: Depth Contours Display**
WHEN the depth layer is enabled, the system shall display depth contour lines at the appropriate zoom levels (visible from approximately zoom 10+).

**[PENDING] F-003: Depth Area Shading**
WHEN the depth layer is enabled, the system shall display depth areas with color shading (darker for deeper, lighter for shallower).

**[PENDING] F-004: Layer Ordering**
WHERE the depth layer is enabled, the system shall render depth areas and contour lines below spot markers and user location markers.

**[PENDING] F-005: Default State**
WHEN the map first loads, the system shall have the depth layer disabled by default.

**[PENDING] F-006: Attribution**
WHERE the depth layer uses VectorCharts.com data, the system shall display attribution text as required by the provider's terms.

### Non-Functional Requirements

**[PENDING] NF-001: Tile Requests**
WHILE a user is actively using the depth layer, the system shall limit tile requests to stay within the 25k/month free tier for typical usage (~500 sessions/month).

**[PENDING] NF-002: No Performance Impact When Disabled**
WHERE the depth layer is disabled, the system shall not load any depth tile data or impact map rendering performance.

---

## Acceptance Criteria

The feature is **Complete** when all of the following are verified:

1. [ ] Depth toggle button is visible on the map (bottom-right)
2. [ ] Click toggle → depth contours and shaded areas appear on the map
3. [ ] Click toggle again → depth layer disappears
4. [ ] Depth layer renders below spot markers (markers remain visible and clickable)
5. [ ] Depth layer is off on initial page load
6. [ ] Toggle works correctly across page navigation (map re-mounts)
7. [ ] Attribution visible when depth layer is enabled
