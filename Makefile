.PHONY: up down dev migrate seed test clean build proto

up:
	docker compose up -d postgres redis localstack minio
	@echo "Waiting for services to be healthy..."
	@until docker compose exec postgres pg_isready -U fishwish 2>/dev/null; do sleep 1; done
	@echo "Infrastructure ready."

down:
	docker compose down

dev: up
	@echo "Starting Go services with hot reload..."
	@cd services/user-service && air -c .air.toml &
	@cd services/spot-service && air -c .air.toml &
	@cd services/search-service && air -c .air.toml &
	@cd services/weather-service && air -c .air.toml &
	@cd services/social-service && air -c .air.toml &
	@cd frontend && npm run dev &
	@echo "All services started. Ctrl+C to stop."
	@wait

dev-docker:
	docker compose up --build

migrate:
	@cd services/spot-service && go run cmd/migrate/main.go up

migrate-down:
	@cd services/spot-service && go run cmd/migrate/main.go down

seed:
	@go run ./scripts/seed/main.go

test:
	@go test ./services/... -v -count=1

test-coverage:
	@go test ./services/... -v -count=1 -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html

build:
	@docker compose build

clean: down
	docker compose down -v --rmi local
	rm -rf coverage.html coverage.out

setup:
	@echo "Installing Go tools..."
	@go install github.com/air-verse/air@latest
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "Installing frontend dependencies..."
	@cd frontend && npm install
	@echo "Setup complete. Run 'make dev' to start."
