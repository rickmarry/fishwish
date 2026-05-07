import { useEffect, useRef } from "react";

const DEFAULT_CENTER = [39.8283, -98.5795];
const DEFAULT_ZOOM = 4;

function MapView({ spots = [], userLocation, selectedSpot, onSpotClick }) {
  const containerRef = useRef(null);
  const mapRef = useRef(null);

  useEffect(() => {
    if (!containerRef.current || mapRef.current) return;

    const loadMap = async () => {
      const maplibregl = (await import("maplibre-gl")).default;

      mapRef.current = new maplibregl.Map({
        container: containerRef.current,
        style: {
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
        },
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
    <div ref={containerRef} className="w-full h-full bg-gray-100" />
  );
}

export default MapView;
