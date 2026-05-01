import { useState, useEffect } from "react";
import MapView from "../../components/MapView/MapView";
import SpotCard from "../../components/SpotCard/SpotCard";
import { useSpots } from "../../hooks/useSpots";
import { useGeoLocation } from "../../hooks/useGeoLocation";

function MapPage() {
  const { spots, loading, fetchNearby } = useSpots();
  const { location: userLocation } = useGeoLocation();
  const [selectedSpot, setSelectedSpot] = useState(null);
  const [showSidebar, setShowSidebar] = useState(true);

  useEffect(() => {
    if (userLocation) {
      fetchNearby(userLocation.lat, userLocation.lon);
    }
  }, [userLocation]);

  const handleSpotClick = (spot) => {
    setSelectedSpot(spot);
  };

  return (
    <div className="relative h-[calc(100vh-4rem)]">
      <MapView
        spots={spots}
        userLocation={userLocation}
        selectedSpot={selectedSpot}
        onSpotClick={handleSpotClick}
      />

      {showSidebar && (
        <aside className="absolute top-4 left-4 bottom-4 w-80 bg-white rounded-xl shadow-lg overflow-hidden flex flex-col">
          <div className="p-4 border-b border-gray-200 bg-white">
            <h2 className="font-semibold text-gray-900">Nearby Spots</h2>
            <p className="text-sm text-gray-500">
              {loading ? "Loading..." : `${spots.length} spots found`}
            </p>
          </div>

          <div className="flex-1 overflow-y-auto p-4 space-y-3">
            {loading && (
              <div className="text-center py-8 text-gray-400">Loading spots...</div>
            )}
            {!loading && spots.length === 0 && (
              <div className="text-center py-8 text-gray-400">
                No spots found nearby.
              </div>
            )}
            {spots.map((spot) => (
              <SpotCard
                key={spot.id}
                spot={spot}
                onClick={() => handleSpotClick(spot)}
                isSelected={selectedSpot?.id === spot.id}
              />
            ))}
          </div>
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
