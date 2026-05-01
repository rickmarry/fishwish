module fishwish/services/weather-service

go 1.22

require (
	fishwish v0.0.0
	github.com/go-chi/chi/v5 v5.1.0
	github.com/go-chi/cors v1.2.1
	github.com/joho/godotenv v1.5.1
	github.com/redis/go-redis/v9 v9.7.0
)

replace fishwish => ../../
