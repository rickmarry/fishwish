# B-011 — Map: Zoom to Spot on Marker Click (LLD)

**Date:** 2026-05-07
**Status:** Design Complete
**Related:** `docs/backlog.md` B-011

---

## What This Is

When a user clicks a spot marker on the map, the map flies to that spot's coordinates at zoom level 12 so the user can see the precise location and surrounding area without manually zooming.

---

## Why This Matters

Without zoom-on-click, clicking a marker highlights it but the map stays at whatever zoom it was at (default zoom 4 shows the entire US). The user has to manually zoom in to understand where the spot actually is. This is a basic map interaction expectation.

---

## Interaction Flow

```
User clicks spot marker on map
  ↓
onSpotClick fires with spot object
  ↓
MapPage sets selectedSpot = spot
  ↓
MapView useEffect detects selectedSpot change
  ↓
mapRef.current.flyTo({ center: [spot.lon, spot.lat], zoom: 12 })
  ↓
Map animates to new center and zoom level
```

---

## Zoom Level

**Zoom 12** — city/neighborhood level. Shows ~5km across. Enough context to see nearby roads, water bodies, and terrain without being too zoomed in to lose orientation.

Comparison:
- Zoom 4 (default) — entire continental US
- Zoom 8 — regional (e.g., Lake Tahoe area)
- Zoom 10 — county level
- **Zoom 12** — town/neighborhood
- Zoom 14 — street level (too zoomed for a fishing spot)

---

## Behavior Details

**Animation:** MapLibre `flyTo` with default parameters (~1.2s duration, easing curve). If a spot is already selected and the user clicks it again, the animation re-runs (idempotent).

**Marker highlight preserved:** Selected spot marker remains highlighted (scale+color change independent of the zoom).

**Caveats handled:**
1. **Rapid clicks** — MapLibre aborts in-flight flyTo animations gracefully, so clicking spots quickly just jumps to the latest one
2. **No deselection** — clicking empty map space currently doesn't clear `selectedSpot`. The map won't jump back to default because the useEffect only fires on non-null `selectedSpot` changes
3. **Sidebar clicks** — clicking a spot card in the sidebar already calls the same `handleSpotClick`, so zoom works from both the map and sidebar

---

## Files Changed

```
frontend/src/components/MapView/MapView.jsx   — add useEffect for flyTo
```

No new dependencies. No backend changes. No migrations.

---

## Testing

**Manual test:**
1. Load map at localhost:3006
2. Click a spot marker → map flies to it at zoom 12
3. Click a different spot marker → map flies to that one
4. Click sidebar spot card → same behavior
5. Click marker rapidly → animation aborts cleanly

No unit test needed (pure UI interaction, MapLibre is a third-party library).
