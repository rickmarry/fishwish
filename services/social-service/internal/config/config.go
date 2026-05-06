package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DBHost       string
	DBPort       string
	DBName       string
	DBUser       string
	DBPass       string
	RedisURL     string
	MinioEndpoint string
	MinioKey     string
	MinioSecret  string
	MinioBucket  string
	MinioUseSSL  bool
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env.local")

	cfg := &Config{
		Port:         getEnv("PORT", "8085"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5433"),
		DBName:       getEnv("DB_NAME", "fishwish"),
		DBUser:       getEnv("DB_USER", "fishwish"),
		DBPass:       getEnv("DB_PASSWORD", "fishwish"),
		RedisURL:     getEnv("REDIS_URL", "redis://localhost:6380"),
		MinioEndpoint: getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioKey:     getEnv("MINIO_ACCESS_KEY", "fishwish"),
		MinioSecret:  getEnv("MINIO_SECRET_KEY", "fishwish123"),
		MinioBucket:  getEnv("MINIO_BUCKET", "fishwish-photos"),
		MinioUseSSL:  getEnv("MINIO_USE_SSL", "false") == "true",
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
