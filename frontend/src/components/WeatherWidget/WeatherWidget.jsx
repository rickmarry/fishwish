function WeatherWidget({ conditions, fishingScore }) {
  if (!conditions) return null;

  const scoreColor =
    fishingScore >= 75
      ? "text-green-600"
      : fishingScore >= 50
        ? "text-yellow-600"
        : "text-red-600";

  const scoreLabel =
    fishingScore >= 75
      ? "Excellent"
      : fishingScore >= 50
        ? "Good"
        : fishingScore >= 25
          ? "Fair"
          : "Poor";

  return (
    <div className="card p-6">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold text-gray-900">Conditions</h2>
        <div className="text-center">
          <div className={`text-3xl font-bold ${scoreColor}`}>{fishingScore}</div>
          <div className={`text-sm ${scoreColor}`}>{scoreLabel}</div>
        </div>
      </div>

      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
        <Stat label="Air Temp" value={`${Math.round(conditions.temperature_f)}°F`} />
        <Stat label="Water Temp" value={conditions.water_temp_f ? `${Math.round(conditions.water_temp_f)}°F` : "N/A"} />
        <Stat label="Wind" value={`${Math.round(conditions.wind_speed_mph)} mph ${conditions.wind_direction}`} />
        <Stat label="Pressure" value={`${conditions.barometric_pressure_in.toFixed(2)} in`} />
        <Stat label="Humidity" value={`${Math.round(conditions.humidity)}%`} />
        <Stat label="Cloud Cover" value={`${Math.round(conditions.cloud_cover)}%`} />
        <Stat label="Rain Chance" value={`${Math.round(conditions.precip_chance)}%`} />
        <Stat label="Tide" value={conditions.tide_state !== "unknown" ? `${conditions.tide_height_ft} ft ${conditions.tide_state}` : "N/A"} />
      </div>
    </div>
  );
}

function Stat({ label, value }) {
  return (
    <div className="text-center">
      <div className="text-lg font-semibold text-gray-900">{value}</div>
      <div className="text-xs text-gray-500">{label}</div>
    </div>
  );
}

export default WeatherWidget;
