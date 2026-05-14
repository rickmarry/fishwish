import { useState, useEffect } from "react";
import MapView from "../../components/MapView/MapView";
import SpotCard from "../../components/SpotCard/SpotCard";
import WeatherWidget from "../../components/WeatherWidget/WeatherWidget";
import { useSpots } from "../../hooks/useSpots";
import { useGeoLocation } from "../../hooks/useGeoLocation";
import { spotsAPI, weatherAPI } from "../../state/api";

function MapPage() {
  const { spots, loading, fetchNearby } = useSpots();
  const { location: userLocation } = useGeoLocation();
  const [selectedSpot, setSelectedSpot] = useState(null);
  const [selectedSpotWeather, setSelectedSpotWeather] = useState(null);
  const [userSpots, setUserSpots] = useState([]);
  const [showSidebar, setShowSidebar] = useState(true);

  const allSpots = [...spots, ...userSpots];

  useEffect(() => {
    if (userLocation) {
      fetchNearby(userLocation.lat, userLocation.lon);
    }
  }, [userLocation]);

  useEffect(() => {
    if (!selectedSpot) {
      setSelectedSpotWeather(null);
      return;
    }
    weatherAPI
      .get({ lat: selectedSpot.lat, lng: selectedSpot.lon })
      .then((res) => setSelectedSpotWeather(res.data))
      .catch(() => setSelectedSpotWeather(null));
  }, [selectedSpot]);

  const handleSpotClick = (spot) => {
    setSelectedSpot(spot);
  };

  const handleMapClick = (lat, lon) => {
    const spot = {
      id: `user-${Date.now()}`,
      name: "New Spot",
      lat,
      lon,
      type: "pin",
      difficulty: "easy",
      rating: 0,
      review_count: 0,
      species: [],
      best_seasons: [],
      created_at: new Date().toISOString(),
    };
    setUserSpots((prev) => [spot, ...prev]);
    setSelectedSpot(spot);
  };

  const handleDeleteSpot = (id) => {
    setUserSpots((prev) => prev.filter((s) => s.id !== id));
    setSelectedSpot((prev) => (prev?.id === id ? null : prev));
  };

  const handleRenameSpot = (id, name) => {
    setUserSpots((prev) =>
      prev.map((s) => (s.id === id ? { ...s, name } : s))
    );
    setSelectedSpot((prev) => (prev?.id === id ? { ...prev, name } : prev));
  };

  const handleSaveSpot = async (id) => {
    const spot = userSpots.find((s) => s.id === id);
    if (!spot) return;

    try {
      const res = await spotsAPI.create({
        name: spot.name,
        lat: spot.lat,
        lon: spot.lon,
        type: "lake",
        difficulty: spot.difficulty,
      });
      setUserSpots((prev) =>
        prev.map((s) =>
          s.id === id ? { ...res.data, species: [], best_seasons: [] } : s
        )
      );
      setSelectedSpot((prev) =>
        prev?.id === id ? { ...res.data, species: [], best_seasons: [] } : prev
      );
    } catch {
      // save failed silently
    }
  };

  return (
    <div className="relative h-[calc(100vh-4rem)]">
      <MapView
        spots={allSpots}
        userLocation={userLocation}
        selectedSpot={selectedSpot}
        onSpotClick={handleSpotClick}
        onMapClick={handleMapClick}
      />

      {showSidebar && (
        <aside className="absolute top-4 left-4 bottom-4 w-80 bg-white rounded-xl shadow-lg overflow-hidden flex flex-col">
          <div className="p-4 border-b border-gray-200 bg-white">
            <h2 className="font-semibold text-gray-900">Spots</h2>
            <p className="text-sm text-gray-500">
              {loading ? "Loading..." : `${allSpots.length} ${allSpots.length === 1 ? "spot" : "spots"}`}
            </p>
          </div>

          <div className="flex-1 overflow-y-auto p-4 space-y-3">
            {loading && (
              <div className="text-center py-8 text-gray-400">Loading spots...</div>
            )}
            {!loading && allSpots.length === 0 && (
              <div className="text-center py-8 text-gray-400">
                Click the map to add a spot.
              </div>
            )}
            {allSpots.map((spot) => (
              <SpotCard
                key={spot.id}
                spot={spot}
                onClick={() => handleSpotClick(spot)}
                isSelected={selectedSpot?.id === spot.id}
                onDelete={spot.type === "pin" ? handleDeleteSpot : undefined}
                onRename={spot.type === "pin" ? handleRenameSpot : undefined}
                onSave={spot.type === "pin" ? handleSaveSpot : undefined}
              />
            ))}
          </div>

          {selectedSpot && selectedSpotWeather && (
            <div className="border-t border-gray-200 overflow-y-auto">
              <WeatherWidget
                current={selectedSpotWeather.current}
                forecast={selectedSpotWeather.forecast}
              />
            </div>
          )}
        </aside>
      )}

      <button
        onClick={() => setShowSidebar(!showSidebar)}
        className="absolute top-4 right-4 z-10 bg-white p-2 rounded-lg shadow-md hover:bg-gray-50"
      >
        {showSidebar ? "Hide" : "Show"} List
      </button>
    </div>
  );
}

export default MapPage;
