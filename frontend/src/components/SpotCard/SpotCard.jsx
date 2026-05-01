import { Link } from "react-router-dom";

function SpotCard({ spot, onClick, isSelected }) {
  const typeIcons = {
    lake: "&#127754;",
    river: "&#127754;",
    ocean: "&#127754;",
    pond: "&#127754;",
    stream: "&#127754;",
    reservoir: "&#127754;",
    estuary: "&#127754;",
  };

  const difficultyColors = {
    easy: "text-green-600",
    moderate: "text-yellow-600",
    hard: "text-red-600",
  };

  return (
    <div
      onClick={onClick}
      className={`card p-4 cursor-pointer transition-all ${
        isSelected ? "ring-2 ring-ocean-500 bg-ocean-50" : "hover:shadow-md"
      }`}
    >
      <div className="flex items-start gap-3">
        <span className="text-2xl" dangerouslySetInnerHTML={{ __html: typeIcons[spot.type] || "&#128031;" }} />
        <div className="flex-1 min-w-0">
          <h3 className="font-semibold text-gray-900 truncate">{spot.name}</h3>
          <div className="flex items-center gap-2 mt-1">
            <span className="text-xs text-gray-500 capitalize">{spot.type}</span>
            <span className={`text-xs ${difficultyColors[spot.difficulty]}`}>
              {spot.difficulty}
            </span>
          </div>

          {spot.species && spot.species.length > 0 && (
            <div className="flex flex-wrap gap-1 mt-2">
              {spot.species.slice(0, 3).map((s) => (
                <span key={s} className="text-xs bg-gray-100 text-gray-600 px-1.5 py-0.5 rounded">
                  {s.replace(/_/g, " ")}
                </span>
              ))}
              {spot.species.length > 3 && (
                <span className="text-xs text-gray-400">+{spot.species.length - 3}</span>
              )}
            </div>
          )}

          {spot.rating > 0 && (
            <div className="flex items-center gap-1 mt-2 text-sm text-yellow-600">
              {"★".repeat(Math.round(spot.rating))}
              <span className="text-gray-400 text-xs">({spot.review_count})</span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default SpotCard;
