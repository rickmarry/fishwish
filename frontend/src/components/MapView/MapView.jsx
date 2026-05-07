import { useEffect, useRef, useState } from "react";

const DEFAULT_CENTER = [-98.5795, 39.8283];
const DEFAULT_ZOOM = 4;

const VECTORCHARTS_KEY = import.meta.env.VITE_VECTORCHARTS_KEY;

function MapView({ spots = [], userLocation, selectedSpot, onSpotClick }) {
  const containerRef = useRef(null);
  const mapRef = useRef(null);
  const [showDepth, setShowDepth] = useState(false);

  useEffect(() => {
    if (!containerRef.current || mapRef.current) return;

    const loadMap = async () => {
      const maplibregl = (await import("maplibre-gl")).default;

      const style = {
        version: 8,
        name: "FishFinder Style",
        sources: {
          osm: {
            type: "raster",
            tiles: ["https://a.tile.openstreetmap.org/{z}/{x}/{y}.png"],
            tileSize: 256,
            attribution: "© OpenStreetMap Contributors",
            maxzoom: 19,
          },
        },
        layers: [
          {
            id: "osm",
            type: "raster",
            source: "osm",
          },
        ],
      };

      if (VECTORCHARTS_KEY && VECTORCHARTS_KEY !== "your_vectorcharts_api_key_here") {
        style.sources.depth = {
          type: "vector",
          tiles: [
            `https://api.vectorcharts.com/v1/tiles/{z}/{x}/{y}.pbf?key=${VECTORCHARTS_KEY}`,
          ],
          maxzoom: 14,
        };
        style.layers.push(
          {
            id: "depth-areas",
            type: "fill",
            source: "depth",
            "source-layer": "deptharea",
            paint: {
              "fill-color": [
                "interpolate",
                ["linear"],
                ["get", "depth"],
                -8000, "#0a1628",
                -200, "#1a3a5c",
                -50, "#2a5a8c",
                -10, "#4a8aba",
                -2, "#7abada",
                0, "#c8e8f0",
              ],
              "fill-opacity": 0.4,
            },
            layout: { visibility: "none" },
          },
          {
            id: "depth-contours",
            type: "line",
            source: "depth",
            "source-layer": "depthcontour",
            paint: {
              "line-color": "#4a6a8a",
              "line-opacity": 0.6,
              "line-width": 1,
            },
            layout: { visibility: "none" },
          }
        );
      }

      mapRef.current = new maplibregl.Map({
        container: containerRef.current,
        style,
        center: DEFAULT_CENTER,
        zoom: DEFAULT_ZOOM,
      });

      mapRef.current.addControl(new maplibregl.NavigationControl(), "top-right");

      if (userLocation) {
        mapRef.current.flyTo({
          center: [userLocation.lon, userLocation.lat],
          zoom: 10,
        });
      }
    };

    loadMap();

    return () => {
      if (mapRef.current) {
        mapRef.current.remove();
        mapRef.current = null;
      }
    };
  }, []);

  useEffect(() => {
    if (!mapRef.current || !VECTORCHARTS_KEY) return;
    ["depth-areas", "depth-contours"].forEach((id) => {
      const layer = mapRef.current.getLayer(id);
      if (layer) {
        mapRef.current.setLayoutProperty(
          id,
          "visibility",
          showDepth ? "visible" : "none"
        );
      }
    });
  }, [showDepth]);

  useEffect(() => {
    if (!mapRef.current || spots.length === 0) return;

    spots.forEach((spot) => {
      const el = document.createElement("div");
      el.className = `w-4 h-4 rounded-full border-2 border-white shadow-md cursor-pointer transition-transform ${
        selectedSpot?.id === spot.id
          ? "bg-ocean-600 scale-150"
          : "bg-forest-500 hover:scale-125"
      }`;
      el.title = spot.name;
      el.addEventListener("click", () => onSpotClick?.(spot));

      new maplibregl.Marker(el)
        .setLngLat([spot.lon, spot.lat])
        .addTo(mapRef.current);
    });
  }, [spots, selectedSpot]);

  useEffect(() => {
    if (!mapRef.current || !selectedSpot) return;
    mapRef.current.flyTo({ center: [selectedSpot.lon, selectedSpot.lat], zoom: 12 });
  }, [selectedSpot]);

  useEffect(() => {
    if (!mapRef.current || !userLocation) return;

    const el = document.createElement("div");
    el.className = "w-4 h-4 rounded-full bg-blue-500 border-2 border-white shadow-md";

    new maplibregl.Marker(el)
      .setLngLat([userLocation.lon, userLocation.lat])
      .addTo(mapRef.current);
  }, [userLocation]);

  return (
    <div className="relative w-full h-full">
      <div ref={containerRef} className="w-full h-full bg-gray-100" />
      {VECTORCHARTS_KEY && VECTORCHARTS_KEY !== "your_vectorcharts_api_key_here" && (
        <button
          onClick={() => setShowDepth(!showDepth)}
          className={`absolute bottom-4 right-4 z-10 px-3 py-1.5 rounded-lg text-xs font-medium shadow-md transition-colors ${
            showDepth
              ? "bg-ocean-600 text-white"
              : "bg-white text-gray-700 hover:bg-gray-50"
          }`}
        >
          {showDepth ? "Depth: ON" : "Depth: OFF"}
        </button>
      )}
    </div>
  );
}

export default MapView;
