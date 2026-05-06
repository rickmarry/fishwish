package model

type WeatherResponse struct {
	Current  CurrentWeather  `json:"current"`
	Forecast []DailyForecast `json:"forecast"`
	Cached   bool            `json:"cached"`
}

type CurrentWeather struct {
	TemperatureC     float64 `json:"temperature_c"`
	HumidityPercent  int     `json:"humidity_percent"`
	PrecipitationMM  float64 `json:"precipitation_mm"`
	WeatherCode      int     `json:"weather_code"`
	WeatherDesc      string  `json:"weather_description"`
	WindSpeedKPH     float64 `json:"wind_speed_kph"`
	WindDirectionDeg int     `json:"wind_direction_deg"`
	PressureHPA      float64 `json:"pressure_hpa"`
}

type DailyForecast struct {
	Date            string  `json:"date"`
	TempMaxC        float64 `json:"temp_max_c"`
	TempMinC        float64 `json:"temp_min_c"`
	PrecipitationMM float64 `json:"precipitation_mm"`
	WindSpeedMaxKPH float64 `json:"wind_speed_max_kph"`
	WeatherCode     int     `json:"weather_code"`
	WeatherDesc     string  `json:"weather_description"`
}

func WeatherCodeToDescription(code int) string {
	switch {
	case code == 0:
		return "Clear sky"
	case code <= 3:
		return "Partly cloudy"
	case code == 45 || code == 48:
		return "Fog"
	case code >= 51 && code <= 67:
		return "Rain"
	case code >= 71 && code <= 77:
		return "Snow"
	case code >= 80 && code <= 82:
		return "Rain showers"
	case code >= 85 && code <= 86:
		return "Snow showers"
	case code >= 95:
		return "Thunderstorm"
	default:
		return "Unknown"
	}
}
