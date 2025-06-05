# 📦 Makefile for Redis Intro (fully containerized version)

SERVICE_NAME := app
REDIS_SERVICE := redis

build:
	docker compose build

up:
	docker compose up --detach

logs:
	docker compose logs -f $(SERVICE_NAME)

down:
	docker compose down

restart:
	docker compose down
	docker compose up --detach --build

redis-logs:
	docker compose logs -f $(REDIS_SERVICE)

redis-cli:
	docker compose exec $(REDIS_SERVICE) redis-cli

ps:
	docker compose ps

# Test endpoints
test-health:
	curl http://localhost:8080/health

test-set:
	curl -X POST http://localhost:8080/cache \
		-H "Content-Type: application/json" \
		-d '{"key":"hello","value":"world","ttl":300}'

test-get:
	curl http://localhost:8080/cache/hello

test-keys:
	curl http://localhost:8080/cache

test-delete:
	curl -X DELETE http://localhost:8080/cache/hello

test-ttl:
	curl http://localhost:8080/cache/hello/ttl

clean:
	docker compose down -v --remove-orphans

rebuild:
	make clean
	make up

help:
	@echo "Available commands:"
	@echo "  build      - Build the Docker images"
	@echo "  up         - Start services in background"
	@echo "  down       - Stop services"
	@echo "  logs       - Show app logs"
	@echo "  redis-logs - Show Redis logs"
	@echo "  redis-cli  - Open Redis CLI"
	@echo "  test-*     - Test various endpoints"
	@echo "  clean      - Remove all containers and volumes"
	@echo "  rebuild    - Clean and rebuild everything" 