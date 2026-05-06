function WeatherWidget({ current, forecast }) {
  if (!current) return null;

  return (
    <div className="card p-6">
      <h2 className="text-lg font-semibold text-gray-900 mb-4">Weather</h2>

      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
        <Stat label="Temp" value={`${Math.round(current.temperature_c)}°C`} />
        <Stat label="Humidity" value={`${current.humidity_percent}%`} />
        <Stat label="Wind" value={`${Math.round(current.wind_speed_kph)} kph ${degToDir(current.wind_direction_deg)}`} />
        <Stat label="Pressure" value={`${current.pressure_hpa} hPa`} />
        <Stat label="Precip" value={`${current.precipitation_mm} mm`} />
        <Stat label="Conditions" value={current.weather_description} />
      </div>

      {forecast && forecast.length > 0 && (
        <div className="mt-6">
          <h3 className="font-semibold text-gray-900 mb-2">3-Day Forecast</h3>
          <div className="grid grid-cols-3 gap-4">
            {forecast.map((day) => (
              <div key={day.date} className="text-center bg-gray-50 rounded p-3">
                <div className="text-sm text-gray-500">{new Date(day.date).toLocaleDateString('en-US', { weekday: 'short' })}</div>
                <div className="text-sm font-semibold">{day.weather_description}</div>
                <div className="text-sm">{Math.round(day.temp_min_c)}° / {Math.round(day.temp_max_c)}°</div>
                <div className="text-xs text-gray-500">{day.precipitation_mm}mm rain</div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

function degToDir(deg) {
  const dirs = ['N', 'NE', 'E', 'SE', 'S', 'SW', 'W', 'NW'];
  return dirs[Math.round(deg / 45) % 8];
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
