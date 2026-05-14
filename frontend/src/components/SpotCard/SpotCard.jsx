import { useState } from "react";

function SpotCard({ spot, onClick, isSelected, onDelete, onRename, onSave }) {
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState(spot.name);
  const isUserPin = spot.type === "pin";

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

  const saveName = () => {
    const trimmed = draft.trim();
    if (trimmed && trimmed !== spot.name) {
      onRename?.(spot.id, trimmed);
    }
    setEditing(false);
  };

  const handleCardClick = (e) => {
    if (editing) return;
    onClick?.();
  };

  const handleDelete = (e) => {
    e.stopPropagation();
    onDelete?.(spot.id);
  };

  const handleKeyDown = (e) => {
    if (e.key === "Enter") saveName();
    if (e.key === "Escape") {
      setDraft(spot.name);
      setEditing(false);
    }
  };

  return (
    <div
      onClick={handleCardClick}
      className={`card p-4 cursor-pointer transition-all relative ${
        isSelected ? "ring-2 ring-ocean-500 bg-ocean-50" : "hover:shadow-md"
      }`}
    >
      {isUserPin && onDelete && (
        <button
          onClick={handleDelete}
          className="absolute top-2 right-2 w-5 h-5 flex items-center justify-center text-xs text-gray-400 hover:text-red-500 hover:bg-red-50 rounded"
          title="Delete spot"
        >
          &#10005;
        </button>
      )}

      <div className="flex items-start gap-3">
        <span className="text-2xl" dangerouslySetInnerHTML={{ __html: typeIcons[spot.type] || "&#128031;" }} />
        <div className="flex-1 min-w-0">
          {editing ? (
            <input
              autoFocus
              value={draft}
              onChange={(e) => setDraft(e.target.value)}
              onBlur={saveName}
              onKeyDown={handleKeyDown}
              onClick={(e) => e.stopPropagation()}
              className="w-full font-semibold text-gray-900 border-b border-ocean-400 outline-none bg-transparent"
            />
          ) : (
            <h3
              className={`font-semibold text-gray-900 truncate ${isUserPin && onRename ? "cursor-text hover:text-ocean-600" : ""}`}
              onDoubleClick={(e) => {
                if (!isUserPin || !onRename) return;
                e.stopPropagation();
                setEditing(true);
              }}
              title={isUserPin && onRename ? "Double-click to rename" : undefined}
            >
              {spot.name}
            </h3>
          )}
          <div className="flex items-center gap-2 mt-1">
            <span className="text-xs text-gray-500 capitalize">{isUserPin ? "pin" : spot.type}</span>
            {!isUserPin && (
              <span className={`text-xs ${difficultyColors[spot.difficulty]}`}>
                {spot.difficulty}
              </span>
            )}
            {isUserPin && (
              <span className="text-xs text-gray-400">
                {spot.lat.toFixed(4)}, {spot.lon.toFixed(4)}
              </span>
            )}
          </div>

          {isUserPin && onSave && (
            <button
              onClick={(e) => { e.stopPropagation(); onSave(spot.id); }}
              className="mt-2 text-xs font-medium px-2 py-0.5 rounded bg-ocean-600 text-white hover:bg-ocean-700"
            >
              Save
            </button>
          )}

          {spot.species && spot.species.length > 0 && !isUserPin && (
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

          {spot.rating > 0 && !isUserPin && (
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
