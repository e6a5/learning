# 🔥 Error Handling Lab - Makefile

# Service names
MYSQL_SERVICE := mysql
REDIS_SERVICE := redis
APP_NAME := error-handling-server

# Development commands
.PHONY: help setup clean run logs test

help: ## Show this help message
	@echo "🔥 Error Handling Learning Lab"
	@echo ""
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*##"; printf "\033[36m\033[0m"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Setup Commands

setup: ## Start MySQL and Redis dependencies
	@echo "🚀 Starting dependencies..."
	docker compose up -d $(MYSQL_SERVICE) $(REDIS_SERVICE)
	@echo "⏳ Waiting for services to be ready..."
	@sleep 10
	@echo "✅ Dependencies ready!"

setup-with-chaos: ## Start dependencies with chaos engineering tools  
	@echo "🚀 Starting all services including chaos tools..."
	docker compose --profile chaos up -d
	@sleep 15
	@echo "✅ All services ready!"

##@ Development Commands

run: ## Run the error handling server locally
	@echo "🔥 Starting Error Handling Server locally..."
	go mod tidy
	go run main.go

run-local: run ## Alias for run

dev: setup run ## Setup dependencies and run server locally

build: ## Build the Go application
	@echo "🔨 Building application..."
	go build -o bin/error-handling-server main.go
	@echo "✅ Build complete!"

##@ Docker Commands

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker compose build app
	@echo "✅ Docker image built!"

docker-run: ## Run the full application stack with Docker
	@echo "🚀 Starting full application stack..."
	docker compose up -d
	@echo "⏳ Waiting for all services to be ready..."
	@sleep 30
	@echo "✅ Application stack ready!"
	@echo "🌐 Application: http://localhost:8080"
	@echo "📊 Health: http://localhost:8080/health"

docker-dev: docker-build docker-run ## Build and run with Docker

docker-logs: ## View application logs in Docker
	docker compose logs -f app

docker-stop: ## Stop the Docker application stack
	@echo "🛑 Stopping application stack..."
	docker compose down

docker-restart: ## Restart the Docker application stack
	@echo "🔄 Restarting application stack..."
	docker compose down
	docker compose up -d --build
	@sleep 30

##@ Testing Commands

test: ## Run basic functionality tests
	@echo "🧪 Testing basic endpoints..."
	curl -s http://localhost:8080/ | jq .
	curl -s http://localhost:8080/health | jq .

test-error-scenarios: ## Test various error scenarios
	@echo "🔥 Testing error scenarios..."
	@echo "1. Testing validation error:"
	curl -s -X POST http://localhost:8080/simulate/validation-error | jq .
	@echo "\n2. Testing panic recovery:"
	curl -s http://localhost:8080/simulate/panic | jq .
	@echo "\n3. Testing database circuit breaker:"
	curl -s http://localhost:8080/simulate/db-error | jq .
	@echo "\n4. Checking circuit breaker status:"
	curl -s http://localhost:8080/circuit-breaker/status | jq .

test-database-down: ## Test behavior when database is unavailable
	@echo "💥 Stopping database to test error handling..."
	docker compose stop $(MYSQL_SERVICE)
	@echo "🧪 Testing endpoints with database down..."
	curl -s http://localhost:8080/users | jq .
	@echo "\n🔄 Restarting database..."
	docker compose start $(MYSQL_SERVICE)
	@sleep 10

test-redis-down: ## Test behavior when Redis is unavailable  
	@echo "💥 Stopping Redis to test error handling..."
	docker compose stop $(REDIS_SERVICE)
	@echo "🧪 Testing with Redis down..."
	curl -s http://localhost:8080/health | jq .
	@echo "\n🔄 Restarting Redis..."
	docker compose start $(REDIS_SERVICE)

test-slow-query: ## Test slow query handling
	@echo "🐌 Testing slow query simulation..."
	time curl -s http://localhost:8080/simulate/slow-query | jq .

##@ Chaos Engineering

chaos-db-kill: ## Kill database container randomly
	@echo "💥 Chaos: Killing database..."
	docker compose kill $(MYSQL_SERVICE)
	@sleep 2
	docker compose start $(MYSQL_SERVICE)

chaos-db-slow: ## Simulate slow database  
	@echo "🐌 Chaos: Slowing down database..."
	docker exec error-handling-mysql mysql -uroot -prootpassword -e "SET GLOBAL innodb_flush_log_at_trx_commit=0;"

chaos-network-partition: ## Simulate network issues
	@echo "🌐 Chaos: Simulating network partition..."
	docker network disconnect error-handling-network error-handling-mysql || true
	@sleep 10
	docker network connect error-handling-network error-handling-mysql

##@ Circuit Breaker Management

circuit-status: ## Check circuit breaker status
	@echo "⚡ Circuit Breaker Status:"
	curl -s http://localhost:8080/circuit-breaker/status | jq .

circuit-reset: ## Reset all circuit breakers
	@echo "🔄 Resetting circuit breakers..."
	curl -s -X POST http://localhost:8080/circuit-breaker/reset | jq .

##@ Monitoring Commands

logs: ## View application logs (if running in container)
	docker compose logs -f || echo "Application not running in container"

logs-mysql: ## View MySQL logs
	docker compose logs $(MYSQL_SERVICE)

logs-redis: ## View Redis logs  
	docker compose logs $(REDIS_SERVICE)

monitor: ## Monitor system during error scenarios
	@echo "📊 Starting monitoring session..."
	@echo "Press Ctrl+C to stop"
	@while true; do \
		echo "=== $$(date) ==="; \
		curl -s http://localhost:8080/health | jq '.data.services, .data.circuit_breakers'; \
		sleep 5; \
	done

##@ Load Testing

load-test: ## Run basic load test
	@echo "⚡ Running load test..."
	@if command -v ab >/dev/null 2>&1; then \
		ab -n 100 -c 10 http://localhost:8080/users; \
	else \
		echo "❌ Apache Bench (ab) not installed. Install with: brew install apache-bench"; \
	fi

load-test-with-failures: ## Load test while simulating failures
	@echo "💥 Load testing with chaos..."
	@make load-test &
	@sleep 2
	@make chaos-db-kill
	@wait

##@ Cleanup Commands

clean: ## Stop and remove all containers and volumes
	@echo "🧹 Cleaning up..."
	docker compose down -v --remove-orphans
	docker volume prune -f

restart: ## Restart all services
	@echo "🔄 Restarting services..."
	docker compose down
	docker compose up -d
	@sleep 10
	@echo "✅ Services restarted!"

##@ Database Commands

db-cli: ## Connect to MySQL database
	docker exec -it error-handling-mysql mysql -uapp_user -papp_password error_handling_db

db-logs: logs-mysql ## Alias for MySQL logs

db-reset: ## Reset database data
	@echo "🗄️ Resetting database..."
	docker compose down $(MYSQL_SERVICE)
	docker volume rm $$(docker volume ls -q | grep error.*mysql) 2>/dev/null || true
	docker compose up -d $(MYSQL_SERVICE)
	@sleep 15

##@ Examples

example-usage: ## Show example API usage
	@echo "🎯 Example API Usage:"
	@echo ""
	@echo "1. Basic health check:"
	@echo "   curl http://localhost:8080/health"
	@echo ""
	@echo "2. Create user:"
	@echo "   curl -X POST http://localhost:8080/users \\"
	@echo "     -H 'Content-Type: application/json' \\"
	@echo "     -d '{\"name\":\"John Doe\",\"email\":\"john@example.com\"}'"
	@echo ""
	@echo "3. Get users (with fallback on DB failure):"
	@echo "   curl http://localhost:8080/users"
	@echo ""
	@echo "4. Test error scenarios:"
	@echo "   curl http://localhost:8080/simulate/panic"
	@echo "   curl -X POST http://localhost:8080/simulate/validation-error"
	@echo ""
	@echo "🐳 Quick Docker Start:"
	@echo "   make docker-dev     # Build and run everything"
	@echo "   make docker-logs    # View logs"
	@echo "   make docker-stop    # Stop all services"

# Default target
.DEFAULT_GOAL := help 