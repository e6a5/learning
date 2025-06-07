# 🔐 Authentication & Security

SERVICE_NAME := app
DB_SERVICE := db

# Build services
build:
	docker compose build

# Build and start services
up:
	docker compose up --detach

# View service logs
logs:
	docker compose logs -f $(SERVICE_NAME)

# Stop services
down:
	docker compose down

# Restart with rebuild
restart:
	docker compose down
	docker compose up --detach --build

# Database logs
db-logs:
	docker logs auth-lab-db

# Database CLI access
db-cli:
	docker exec -it auth-lab-db mysql -uuser -ppass authlab

# Show running services
ps:
	docker compose ps

# Development: run Go app locally (requires db service)
run:
	go mod tidy
	go run main.go

# Test endpoints
test-status:
	curl http://localhost:8081/

test-auth:
	@echo "🔐 Testing authentication flow..."
	@echo "1. Registering a new user..."
	@curl -s -X POST http://localhost:8081/auth/register \
		-H "Content-Type: application/json" \
		-d '{"username":"testuser","password":"securepass123","email":"test@example.com"}' | jq .
	@echo "\n2. Logging in..."
	@curl -s -X POST http://localhost:8081/auth/login \
		-H "Content-Type: application/json" \
		-d '{"username":"testuser","password":"securepass123"}' | jq .

test-admin:
	@echo "👤 Testing with default admin account..."
	@curl -s -X POST http://localhost:8081/auth/login \
		-H "Content-Type: application/json" \
		-d '{"username":"admin","password":"admin123"}' | jq .

test-rate-limit:
	@echo "⚡ Testing rate limiting..."
	@for i in {1..10}; do \
		curl -s -X POST http://localhost:8081/auth/login \
			-H "Content-Type: application/json" \
			-d '{"username":"test","password":"test"}' \
			-w "Request $$i: %{http_code}\n" -o /dev/null; \
	done

# Clean up everything
clean:
	docker compose down -v --remove-orphans

# Full rebuild
rebuild:
	make clean
	make up 