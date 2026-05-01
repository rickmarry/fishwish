package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	RedisURL    string
	WeatherAPI  string
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env.local")

	cfg := &Config{
		Port:     getEnv("PORT", "8084"),
		RedisURL: getEnv("REDIS_URL", "redis://localhost:6379"),
		WeatherAPI: getEnv("WEATHER_API_KEY", ""),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
