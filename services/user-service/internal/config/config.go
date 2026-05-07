package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	DBHost   string
	DBPort   string
	DBName   string
	DBUser   string
	DBPass   string
	RedisURL string
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env.local")

	cfg := &Config{
		Port:     getEnv("PORT", "8086"),
		DBHost:   getEnv("DB_HOST", "localhost"),
		DBPort:   getEnv("DB_PORT", "5433"),
		DBName:   getEnv("DB_NAME", "fishwish"),
		DBUser:   getEnv("DB_USER", "fishwish"),
		DBPass:   getEnv("DB_PASSWORD", "fishwish"),
		RedisURL: getEnv("REDIS_URL", "redis://localhost:6380"),
	}

	return cfg, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
