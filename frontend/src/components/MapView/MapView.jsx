import { useEffect, useRef, useState } from "react";
import { activeProvider, providerKey } from "../../config/depth-providers";

const DEFAULT_CENTER = [-98.5795, 39.8283];
const DEFAULT_ZOOM = 4;

let maplibregl = null;

function MapView({ spots = [], userLocation, selectedSpot, onSpotClick, onMapClick }) {
  const containerRef = useRef(null);
  const mapRef = useRef(null);
  const [showDepth, setShowDepth] = useState(false);

  useEffect(() => {
    if (!containerRef.current || mapRef.current) return;

    const loadMap = async () => {
      maplibregl = (await import("maplibre-gl")).default;

      const style = {
        version: 8,
        name: "FishFinder Style",
        sources: {
          osm: {
            type: "raster",
            tiles: ["https://a.tile.openstreetmap.org/{z}/{x}/{y}.png"],
            tileSize: 256,
            attribution: "\u00a9 OpenStreetMap Contributors",
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

      style.sources.depth = activeProvider.getSourceConfig(providerKey);
      activeProvider.addLayers(style);

      mapRef.current = new maplibregl.Map({
        container: containerRef.current,
        style,
        center: DEFAULT_CENTER,
        zoom: DEFAULT_ZOOM,
      });

      mapRef.current.addControl(new maplibregl.NavigationControl(), "top-right");

      mapRef.current.on("click", (e) => {
        onMapClick?.(e.lngLat.lat, e.lngLat.lng);
      });

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
    if (!mapRef.current) return;
    const layerIds = activeProvider.getLayerIds();
    layerIds.forEach((id) => {
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

  const markersRef = useRef([]);

  useEffect(() => {
    markersRef.current.forEach((m) => m.remove());
    markersRef.current = [];

    if (!mapRef.current) return;

    spots.forEach((spot) => {
      const el = document.createElement("div");
      const isUserPin = spot.type === "pin";
      const isSelected = selectedSpot?.id === spot.id;
      el.className = [
        "w-4 h-4 rounded-full border-2 border-white shadow-md cursor-pointer transition-transform",
        isSelected
          ? "bg-ocean-600 scale-150"
          : isUserPin
            ? "bg-amber-500 hover:scale-125"
            : "bg-forest-500 hover:scale-125",
      ].join(" ");
      el.title = spot.name;
      el.addEventListener("click", () => onSpotClick?.(spot));

      const marker = new maplibregl.Marker(el)
        .setLngLat([spot.lon, spot.lat])
        .addTo(mapRef.current);
      markersRef.current.push(marker);
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
    </div>
  );
}

export default MapView;
