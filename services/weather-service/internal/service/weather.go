package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"fishwish/pkg/weather"
)

type WeatherService struct {
	client *weather.Client
}

func NewWeatherService(apiKey string) *WeatherService {
	return &WeatherService{
		client: weather.NewClient(apiKey),
	}
}

func (s *WeatherService) GetConditions(ctx context.Context, lat, lon float64) (*weather.Conditions, error) {
	return s.client.GetConditions(lat, lon)
}

func (s *WeatherService) GetForecast(ctx context.Context, lat, lon float64) ([]weather.Conditions, error) {
	return nil, nil
}

func (s *WeatherService) GetTides(ctx context.Context, lat, lon float64) ([]TidePrediction, error) {
	return []TidePrediction{}, nil
}

type TidePrediction struct {
	Time   string  `json:"time"`
	Height float64 `json:"height_ft"`
	Type   string  `json:"type"`
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

func (s *WeatherService) GetCache(ctx context.Context, key string) (*weather.Conditions, error) {
	return nil, nil
}

func (s *WeatherService) SetCache(ctx context.Context, key string, data *weather.Conditions) error {
	return nil
}

func cacheKey(lat, lon float64) string {
	return fmt.Sprintf("weather:%.2f:%.2f", lat, lon)
}
