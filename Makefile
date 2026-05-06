.PHONY: up down dev migrate seed test clean build proto

up:
	docker compose up -d postgres redis minio
	@echo "Waiting for services to be healthy..."
	@until docker compose exec postgres pg_isready -U fishwish 2>/dev/null; do sleep 1; done
	@echo "Infrastructure ready."

down:
	docker compose down

dev: up
	@echo "Stopping any running services..."
	@-lsof -i :8081 -i :8082 -i :8083 -i :8084 -i :8085 2>/dev/null | grep LISTEN | awk '{print $$2}' | sort -u | xargs kill 2>/dev/null || true
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
	@for svc in user-service spot-service search-service weather-service social-service; do \
		echo "=== $$svc ==="; \
		cd $(CURDIR)/services/$$svc && go test ./... -v -count=1; \
	done

test-coverage:
	@for svc in user-service spot-service search-service weather-service social-service; do \
		echo "=== $$svc ==="; \
		cd $(CURDIR)/services/$$svc && go test ./... -v -count=1 -coverprofile=coverage.out; \
	done

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
