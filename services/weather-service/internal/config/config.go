package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	RedisAddr  string
	WeatherAPI string
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env.local")

	cfg := &Config{
		Port:      getEnv("PORT", "8084"),
		RedisAddr: redisAddrFromURL(getEnv("REDIS_URL", "redis://localhost:6380")),
		WeatherAPI: getEnv("WEATHER_API_KEY", ""),
	}

	return cfg, nil
}

func redisAddrFromURL(url string) string {
	// Convert redis://host:port to host:port
	if strings.HasPrefix(url, "redis://") {
		return strings.TrimPrefix(url, "redis://")
	}
	return url
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
