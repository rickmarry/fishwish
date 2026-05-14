import { useState } from "react";

const STORAGE_KEY = "fishwish_units";

function getInitialUnits() {
  return localStorage.getItem(STORAGE_KEY) === "imperial" ? "imperial" : "metric";
}

function cToF(c) { return c * 9 / 5 + 32; }
function kphToMph(k) { return k * 0.621371; }
function hPaToInHg(h) { return h * 0.02953; }
function mmToIn(mm) { return mm * 0.0393701; }

function WeatherWidget({ current, forecast }) {
  const [units, setUnits] = useState(getInitialUnits);
  const isMetric = units === "metric";

  if (!current) return null;

  const toggle = () => {
    const next = isMetric ? "imperial" : "metric";
    setUnits(next);
    localStorage.setItem(STORAGE_KEY, next);
  };

  return (
    <div className="card p-6">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold text-gray-900">Weather</h2>
        <button
          onClick={toggle}
          className="text-xs font-medium px-2 py-1 rounded bg-gray-100 text-gray-600 hover:bg-gray-200"
        >
          {isMetric ? "\u00b0F / mph" : "\u00b0C / kph"}
        </button>
      </div>

      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
        <Stat
          label="Temp"
          value={`${Math.round(isMetric ? current.temperature_c : cToF(current.temperature_c))}\u00b0${isMetric ? "C" : "F"}`}
        />
        <Stat label="Humidity" value={`${current.humidity_percent}%`} />
        <Stat
          label="Wind"
          value={`${Math.round(isMetric ? current.wind_speed_kph : kphToMph(current.wind_speed_kph))} ${isMetric ? "kph" : "mph"} ${degToDir(current.wind_direction_deg)}`}
        />
        <Stat
          label="Pressure"
          value={`${isMetric ? Math.round(current.pressure_hpa) : hPaToInHg(current.pressure_hpa).toFixed(2)} ${isMetric ? "hPa" : "inHg"}`}
        />
        <Stat
          label="Precip"
          value={`${isMetric ? current.precipitation_mm.toFixed(1) : mmToIn(current.precipitation_mm).toFixed(2)} ${isMetric ? "mm" : "in"}`}
        />
        <Stat label="Conditions" value={current.weather_description} />
      </div>

      {forecast && forecast.length > 0 && (
        <div className="mt-6">
          <h3 className="font-semibold text-gray-900 mb-2">3-Day Forecast</h3>
          <div className="grid grid-cols-3 gap-4">
            {forecast.map((day) => (
              <div key={day.date} className="text-center bg-gray-50 rounded p-3">
                <div className="text-sm text-gray-500">{new Date(day.date).toLocaleDateString("en-US", { weekday: "short", timeZone: "UTC" })}</div>
                <div className="text-sm font-semibold">{day.weather_description}</div>
                <div className="text-sm">
                  {`${Math.round(isMetric ? day.temp_min_c : cToF(day.temp_min_c))}\u00b0 / ${Math.round(isMetric ? day.temp_max_c : cToF(day.temp_max_c))}\u00b0`}
                </div>
                <div className="text-xs text-gray-500">
                  {(isMetric ? day.precipitation_mm : mmToIn(day.precipitation_mm)).toFixed(1)}{isMetric ? "mm" : "in"} rain
                </div>
                <div className="text-xs text-gray-500 mt-0.5">
                  {Math.round(isMetric ? day.wind_speed_max_kph : kphToMph(day.wind_speed_max_kph))} {isMetric ? "kph" : "mph"} {degToDir(day.wind_direction_deg)}
                </div>
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
