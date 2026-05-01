import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import { spotsAPI, weatherAPI, socialAPI } from "../../state/api";
import WeatherWidget from "../../components/WeatherWidget/WeatherWidget";
import ReviewSection from "../../components/ReviewSection/ReviewSection";

function SpotDetailPage() {
  const { id } = useParams();
  const [spot, setSpot] = useState(null);
  const [weather, setWeather] = useState(null);
  const [fishingScore, setFishingScore] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const load = async () => {
      try {
        const [spotRes, weatherRes] = await Promise.all([
          spotsAPI.get(id),
          weatherAPI.get({ lat: 0, lon: 0 }).catch(() => null),
        ]);
        setSpot(spotRes.data);
        if (weatherRes) {
          setWeather(weatherRes.data.conditions);
          setFishingScore(weatherRes.data.fishing_score);
        }
      } catch (e) {
        console.error(e);
      } finally {
        setLoading(false);
      }
    };
    load();
  }, [id]);

  if (loading) {
    return (
      <div className="max-w-4xl mx-auto p-6 text-center py-20">
        Loading...
      </div>
    );
  }

  if (!spot) {
    return (
      <div className="max-w-4xl mx-auto p-6 text-center py-20">
        Spot not found. <Link to="/" className="text-ocean-600">Go back to map</Link>
      </div>
    );
  }

  const typeLabels = {
    lake: "Lake",
    river: "River",
    ocean: "Ocean",
    pond: "Pond",
    stream: "Stream",
    reservoir: "Reservoir",
    estuary: "Estuary",
  };

  const difficultyColors = {
    easy: "bg-green-100 text-green-800",
    moderate: "bg-yellow-100 text-yellow-800",
    hard: "bg-red-100 text-red-800",
  };

  return (
    <div className="max-w-4xl mx-auto p-6">
      <Link to="/" className="text-ocean-600 text-sm mb-4 inline-block">
        &larr; Back to Map
      </Link>

      <div className="card">
        <div className="p-6">
          <div className="flex items-start justify-between mb-4">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">{spot.name}</h1>
              <div className="flex items-center gap-2 mt-2">
                <span className="text-sm bg-ocean-50 text-ocean-700 px-2 py-1 rounded">
                  {typeLabels[spot.type]}
                </span>
                <span className={`text-sm px-2 py-1 rounded ${difficultyColors[spot.difficulty]}`}>
                  {spot.difficulty}
                </span>
                {spot.rating > 0 && (
                  <span className="text-sm text-yellow-600">
                    {"★".repeat(Math.round(spot.rating))} {spot.rating} ({spot.review_count})
                  </span>
                )}
              </div>
            </div>
          </div>

          {spot.description && (
            <p className="text-gray-600 mb-6">{spot.description}</p>
          )}

          <div className="grid md:grid-cols-2 gap-6">
            <div>
              <h3 className="font-semibold text-gray-900 mb-2">Species</h3>
              {spot.species && spot.species.length > 0 ? (
                <div className="flex flex-wrap gap-2">
                  {spot.species.map((s) => (
                    <span key={s} className="bg-forest-50 text-forest-700 text-sm px-2 py-1 rounded">
                      {s.replace(/_/g, " ")}
                    </span>
                  ))}
                </div>
              ) : (
                <p className="text-gray-400 text-sm">No species data</p>
              )}
            </div>

            <div>
              <h3 className="font-semibold text-gray-900 mb-2">Access</h3>
              {spot.access_notes ? (
                <p className="text-gray-600 text-sm">{spot.access_notes}</p>
              ) : (
                <p className="text-gray-400 text-sm">No access info</p>
              )}
              {spot.parking && (
                <p className="text-gray-500 text-sm mt-2">Parking: {spot.parking}</p>
              )}
            </div>
          </div>
        </div>
      </div>

      {weather && (
        <div className="mt-6">
          <WeatherWidget conditions={weather} fishingScore={fishingScore} />
        </div>
      )}

      <div className="mt-6">
        <ReviewSection spotId={id} />
      </div>
    </div>
  );
}

export default SpotDetailPage;
