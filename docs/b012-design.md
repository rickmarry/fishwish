# B-012 — Depth Charts / Bathymetry Overlay (LLD)

**Date:** 2026-05-07
**Status:** Design Complete
**Related:** `docs/backlog.md` B-012, `docs/adr/0002-b012-vectorcharts-provider.md`

---

## What This Is

A toggleable depth chart overlay on the fishing spot map. Shows underwater terrain — depth contours, soundings, and shaded bathymetry — so anglers can find drop-offs, holes, and structure where fish hold.

---

## Why This Matters

Depth is the single most important factor in finding fish. A map that shows where spots are but not the underwater terrain forces users to cross-reference external tools. This closes the loop.

---

## Integration Architecture

```
Existing MapLibre Map (OSM raster tiles)
  └── add vector tile source → VectorCharts API
       └── add layer → depth contours + soundings
       └── toggle button → show/hide the layer
```

No backend changes. No new service. The frontend adds a tile source directly to the existing MapLibre instance.

---

## Provider: VectorCharts.com

**Chosen for V1** because:
- 25k requests/month free (sufficient for single-user/small-scale)
- NOAA ENC data with proper depth contours, soundings, depth areas
- MapLibre-compatible vector tiles (PBF format)
- US + international coverage
- Simple API key authentication

**API Details:**
- Endpoint: `https://api.vectorcharts.com/v1/tiles/{z}/{x}/{y}.pbf?key=API_KEY`
- Returns Mapbox Vector Tiles (protocol buffer)
- Zoom range: 0-14 (depth contours visible from ~zoom 10+)
- Attribution required per their terms

**Pricing escalation:**
- Free: 0-25k requests/month
- $1/1k: 25k-100k requests/month
- Custom: 100k+

---

## MapLibre Integration

### Source Configuration
Add to the MapLibre map style sources:

```javascript
mapRef.current.addSource('depth', {
  type: 'vector',
  tiles: ['https://api.vectorcharts.com/v1/tiles/{z}/{x}/{y}.pbf?key=API_KEY'],
  maxzoom: 14,
});
```

### Layer Configuration
```javascript
mapRef.current.addLayer({
  id: 'depth-contours',
  type: 'line',
  source: 'depth',
  'source-layer': 'depthcontour',
  paint: {
    'line-color': '#4a6a8a',
    'line-opacity': 0.6,
    'line-width': 1,
  },
});

mapRef.current.addLayer({
  id: 'depth-areas',
  type: 'fill',
  source: 'depth',
  'source-layer': 'deptharea',
  paint: {
    'fill-color': [
      'interpolate', ['linear'], ['get', 'depth'],
      -8000, '#0a1628',
      -200, '#1a3a5c',
      -50, '#2a5a8c',
      -10, '#4a8aba',
      -2, '#7abada',
      0, '#c8e8f0',
    ],
    'fill-opacity': 0.4,
  },
});
```

### Layer Order
Depth layers go BELOW the spot markers so markers stay visible on top:
```
OSM raster tiles (base)
  └── Depth areas (fill)
  └── Depth contours (lines)
       └── Spot markers
       └── User location marker
```

---

## Toggle Control

Add a button to toggle the depth layer visibility:

```javascript
const [showDepth, setShowDepth] = useState(false);

// Toggle depth layers
useEffect(() => {
  if (!mapRef.current) return;
  ['depth-areas', 'depth-contours'].forEach((id) => {
    const vis = mapRef.current.getLayer(id)
      ? mapRef.current.getLayoutProperty(id, 'visibility')
      : 'none';
    mapRef.current.setLayoutProperty(
      id,
      'visibility',
      showDepth ? 'visible' : 'none'
    );
  });
}, [showDepth]);
```

UI position: bottom-right of the map, near or grouped with the existing NavigationControl.

---

## API Key Handling

The VectorCharts API key needs to be in the frontend (tile URLs are exposed in the browser). Store as an environment variable:

```
VITE_VECTORCHARTS_KEY=pk_...
```

In `MapView.jsx`:
```javascript
const apiKey = import.meta.env.VITE_VECTORCHARTS_KEY;
```

---

## Files Changed

```
frontend/
  .env.local                              — add VITE_VECTORCHARTS_KEY
  src/components/MapView/MapView.jsx      — add depth source, layers, toggle
```

No backend changes, no migrations, no new services.

---

## Testing

**Manual:**
1. Open map → depth layer is off by default
2. Click toggle → depth contours and shaded areas appear beneath markers
3. Toggle off → depth disappears
4. Pan/zoom → tiles load at appropriate zoom levels

**Free tier budget:**
- Estimate ~50 tile requests per map session (panning a few spots)
- 25k free limit → ~500 sessions/month
- Monitor via VectorCharts dashboard
