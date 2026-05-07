# B-011 — Map: Zoom to Spot on Marker Click (SPEC)

**Date:** 2026-05-07
**Status:** Design Complete
**Related:** `docs/backlog.md` B-011, `docs/b011-design.md`

---

## Overview

Clicking a spot marker on the map flies the map to that spot's location at zoom level 12, so the user can see the precise position and surrounding area.

---

## Requirements (EARS Format)

### Functional Requirements

**[PENDING] F-001: Zoom to Spot on Marker Click**
WHEN a user clicks a spot marker on the map, the system shall fly the map to the spot's coordinates at zoom level 12.

**[PENDING] F-002: Zoom on Sidebar Spot Card Click**
WHEN a user clicks a spot card in the sidebar, the system shall fly the map to the spot's coordinates at zoom level 12 (same behavior as marker click).

**[PENDING] F-003: Marker Highlight Preserved**
WHERE a spot marker is the currently selected spot, the system shall continue to display it with the highlight style (scale increase and color change).

**[PENDING] F-004: Rapid Click Handling**
WHEN a user clicks multiple spot markers in quick succession, the system shall abort the in-flight animation and fly to the most recently clicked spot.

### Non-Functional Requirements

**[PENDING] NF-001: Animation Duration**
WHERE flyTo animation is triggered, the system shall complete the animation within approximately 1.2 seconds using MapLibre's default easing.

**[PENDING] NF-002: No Map Jump on Initial Load**
WHERE no spot is selected, the system shall not trigger any flyTo animation on initial map render.

---

## Acceptance Criteria

The feature is **Complete** when all of the following are verified:

1. [ ] Click a spot marker → map flies to it at zoom 12
2. [ ] Click a sidebar spot card → map flies to that spot at zoom 12
3. [ ] Selected spot marker is highlighted (bigger, different color) after zoom
4. [ ] Rapid clicking multiple spots → animation aborts cleanly, last clicked spot wins
5. [ ] Map does not jump on initial load (no selectedSpot → no flyTo)
