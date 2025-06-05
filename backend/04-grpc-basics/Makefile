# 📦 Makefile for gRPC Basics (fully containerized version with buf)

# Generate protobuf code locally (requires buf installed)
gen:
	buf generate

# Lint protobuf files
lint:
	buf lint

# Format protobuf files
format:
	buf format -w

# Check for breaking changes
breaking:
	buf breaking --against '.git#branch=main'

# Update dependencies
deps:
	buf mod update

# Docker operations
build:
	docker compose build

build-server:
	docker compose build server

build-client:
	docker compose build client

# Start server only
up:
	docker compose up server -d

# Start server with client demo
demo:
	docker compose --profile client up --build

# View logs
logs:
	docker compose logs -f server

logs-client:
	docker compose --profile client logs client

# Development
down:
	docker compose down

restart:
	docker compose down
	docker compose up server -d --build

ps:
	docker compose ps

# Testing with grpcurl (requires grpcurl installed)
test-create:
	grpcurl -plaintext -d '{"name":"John Doe","email":"john@example.com"}' \
		localhost:50051 user.UserService/CreateUser

test-get:
	grpcurl -plaintext -d '{"id":1}' \
		localhost:50051 user.UserService/GetUser

test-list:
	grpcurl -plaintext -d '{"page":1,"limit":10}' \
		localhost:50051 user.UserService/ListUsers

test-watch:
	grpcurl -plaintext -d '{}' \
		localhost:50051 user.UserService/WatchUsers

# List available services
list-services:
	grpcurl -plaintext localhost:50051 list

describe-service:
	grpcurl -plaintext localhost:50051 describe user.UserService

# Development helpers
clean:
	docker compose down -v --remove-orphans
	docker system prune -f

rebuild:
	make clean
	make build
	make up

# Quick demo run
quick-demo:
	@echo "🚀 Starting gRPC demo..."
	@make up
	@echo "⏳ Waiting for server to start..."
	@sleep 5
	@make demo

# Buf development workflow
dev-setup:
	@echo "📦 Setting up buf development environment..."
	@command -v buf >/dev/null 2>&1 || { echo "Installing buf..."; \
		go install github.com/bufbuild/buf/cmd/buf@latest; }
	@echo "✅ Buf is ready!"

proto-workflow:
	@echo "🔧 Running full protobuf workflow..."
	make format
	make lint
	make gen
	@echo "✅ Protobuf workflow complete!"

help:
	@echo "Available commands:"
	@echo "  Proto & Buf:"
	@echo "    gen           - Generate protobuf code with buf"
	@echo "    lint          - Lint protobuf files"
	@echo "    format        - Format protobuf files"
	@echo "    breaking      - Check for breaking changes"
	@echo "    deps          - Update buf dependencies"
	@echo "    proto-workflow- Run full protobuf workflow"
	@echo ""
	@echo "  Docker & Services:"
	@echo "    build         - Build all Docker images"
	@echo "    up            - Start gRPC server"
	@echo "    demo          - Run server + client demo"
	@echo "    logs          - Show server logs"
	@echo "    down          - Stop services"
	@echo ""
	@echo "  Testing:"
	@echo "    test-*        - Test gRPC methods with grpcurl"
	@echo "    quick-demo    - Full demo (server + client)"
	@echo ""
	@echo "  Development:"
	@echo "    dev-setup     - Install buf and setup environment"
	@echo "    clean         - Remove all containers and images" 