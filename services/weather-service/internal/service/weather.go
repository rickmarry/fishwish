package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"fishwish/pkg/weather"
	"fishwish/services/weather-service/internal/model"

	"github.com/redis/go-redis/v9"
)

type WeatherService struct {
	client  *weather.Client
	redis   *redis.Client
	baseURL string
}

func NewWeatherService(redisURL, baseURL string) *WeatherService {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	if baseURL == "" {
		baseURL = "https://api.open-meteo.com/v1/forecast"
	}
	return &WeatherService{
		client:  weather.NewClient(""),
		redis:   rdb,
		baseURL: baseURL,
	}
}

func (s *WeatherService) GetWeather(ctx context.Context, lat, lng float64) (*model.WeatherResponse, error) {
	// Round coordinates to 3 decimal places for cache key
	rLat := roundTo3(lat)
	rLng := roundTo3(lng)
	today := time.Now().Format("2006-01-02")
	key := fmt.Sprintf("weather:%.3f:%.3f:%s", rLat, rLng, today)

	// Try cache first
	cached, err := s.redis.Get(ctx, key).Result()
	if err == nil {
		var resp model.WeatherResponse
		if json.Unmarshal([]byte(cached), &resp) == nil {
			resp.Cached = true
			return &resp, nil
		}
	}

	// Cache miss — call Open-Meteo API with explicit dates
	end := time.Now().Add(2 * 24 * time.Hour).Format("2006-01-02")
	url := fmt.Sprintf(
		"%s?latitude=%.4f&longitude=%.4f"+
			"&current=temperature_2m,relative_humidity_2m,precipitation,weather_code,wind_speed_10m,wind_direction_10m,pressure_msl"+
			"&daily=temperature_2m_max,temperature_2m_min,precipitation_sum,wind_speed_10m_max,weather_code"+
			"&timezone=auto&start_date=%s&end_date=%s",
		s.baseURL, rLat, rLng, today, end,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("weather API unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	var data struct {
		Current struct {
			Temperature     float64 `json:"temperature_2m"`
			Humidity        float64 `json:"relative_humidity_2m"`
			Precipitation   float64 `json:"precipitation"`
			WeatherCode     int     `json:"weather_code"`
			WindSpeed       float64 `json:"wind_speed_10m"`
			WindDirection   int     `json:"wind_direction_10m"`
			Pressure        float64 `json:"pressure_msl"`
		} `json:"current"`
		Daily struct {
			Time           []string  `json:"time"`
			TempMax         []float64 `json:"temperature_2m_max"`
			TempMin         []float64 `json:"temperature_2m_min"`
			Precipitation   []float64 `json:"precipitation_sum"`
			WindSpeedMax   []float64 `json:"wind_speed_10m_max"`
			WeatherCode    []int     `json:"weather_code"`
		} `json:"daily"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	// Build forecast
	forecast := make([]model.DailyForecast, 0, len(data.Daily.Time))
	for i := range data.Daily.Time {
		forecast = append(forecast, model.DailyForecast{
			Date:           data.Daily.Time[i],
			TempMaxC:       data.Daily.TempMax[i],
			TempMinC:       data.Daily.TempMin[i],
			PrecipitationMM: data.Daily.Precipitation[i],
			WindSpeedMaxKPH: data.Daily.WindSpeedMax[i],
			WeatherCode:    data.Daily.WeatherCode[i],
			WeatherDesc:    model.WeatherCodeToDescription(data.Daily.WeatherCode[i]),
		})
	}

	result := &model.WeatherResponse{
		Current: model.CurrentWeather{
			TemperatureC:     data.Current.Temperature,
			HumidityPercent:  int(data.Current.Humidity),
			PrecipitationMM:  data.Current.Precipitation,
			WeatherCode:      data.Current.WeatherCode,
			WeatherDesc:      model.WeatherCodeToDescription(data.Current.WeatherCode),
			WindSpeedKPH:     data.Current.WindSpeed,
			WindDirectionDeg: data.Current.WindDirection,
			PressureHPA:      data.Current.Pressure,
		},
		Forecast: forecast,
		Cached:   false,
	}

	// Cache the result (fire-and-forget, don't fail request if Redis write fails)
	if jsonBytes, err := json.Marshal(result); err == nil {
		_ = s.redis.Set(ctx, key, string(jsonBytes), 30*time.Minute).Err()
	}

	return result, nil
}

func (s *WeatherService) GetForecast(ctx context.Context, lat, lon float64) ([]model.DailyForecast, error) {
	resp, err := s.GetWeather(ctx, lat, lon)
	if err != nil {
		return nil, err
	}
	return resp.Forecast, nil
}

func (s *WeatherService) GetTides(ctx context.Context, lat, lon float64) ([]TidePrediction, error) {
	return []TidePrediction{}, nil
}

func (s *WeatherService) GetFishingScore(ctx context.Context, conditions *weather.Conditions) int {
	score := 50
	if conditions.TemperatureF >= 50 && conditions.TemperatureF <= 85 {
		score += 10
	}
	if conditions.BarometricPress >= 29.8 && conditions.BarometricPress <= 30.2 {
		score += 15
	}
	if conditions.WindSpeedMph < 10 {
		score += 10
	} else if conditions.WindSpeedMph > 20 {
		score -= 15
	}
	if conditions.PrecipChance < 30 {
		score += 5
	}
	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}
	return score
}

type TidePrediction struct {
	Time   string  `json:"time"`
	Height float64 `json:"height_ft"`
	Type   string  `json:"type"`
}

func roundTo3(v float64) float64 {
	r, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", v), 64)
	return r
}
