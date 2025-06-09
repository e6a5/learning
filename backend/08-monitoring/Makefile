# 📊 08-monitoring Makefile
# Observability and Health Monitoring Development Tools

# Configuration
APP_NAME := monitoring-service
PORT := 8080
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "1.0.0")
ENVIRONMENT := development

# Go configuration
GO_CMD := go
GO_BUILD := $(GO_CMD) build
GO_TEST := $(GO_CMD) test
GO_RUN := $(GO_CMD) run
GO_MOD := $(GO_CMD) mod
GO_VET := $(GO_CMD) vet
GO_FMT := gofmt

# Docker configuration
DOCKER_CMD := docker
COMPOSE_CMD := docker compose

# Colors for output
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

.PHONY: help build run test test-unit test-integration test-coverage clean lint format deps docker-build docker-run docker-stop docker-clean compose-up compose-down compose-logs health-check load-test demo all

## Default target
all: deps format lint test build

## 📖 Show help information
help:
	@echo "$(GREEN)📊 08-monitoring: System Observability & Health Monitoring$(NC)"
	@echo ""
	@echo "$(YELLOW)🚀 Development Commands:$(NC)"
	@echo "  make run           - Run the monitoring service locally"
	@echo "  make build         - Build the monitoring service binary"
	@echo "  make deps          - Download and tidy dependencies"
	@echo "  make clean         - Clean build artifacts and logs"
	@echo ""
	@echo "$(YELLOW)🧪 Testing Commands:$(NC)"
	@echo "  make test          - Run all tests with coverage"
	@echo "  make test-unit     - Run unit tests only"
	@echo "  make test-integration - Run integration tests"
	@echo "  make test-coverage - Generate detailed test coverage report"
	@echo "  make health-check  - Test health endpoints"
	@echo "  make load-test     - Generate load for testing metrics"
	@echo "  make demo          - Run demo scenarios"
	@echo ""
	@echo "$(YELLOW)🔍 Code Quality:$(NC)"
	@echo "  make lint          - Run linting checks"
	@echo "  make format        - Format code with gofmt"
	@echo "  make vet           - Run go vet static analysis"
	@echo ""
	@echo "$(YELLOW)🐳 Docker Commands:$(NC)"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run in Docker container"
	@echo "  make docker-stop   - Stop Docker container"
	@echo "  make docker-clean  - Clean Docker resources"
	@echo ""
	@echo "$(YELLOW)🏗️ Compose Commands:$(NC)"
	@echo "  make compose-up    - Start full monitoring stack"
	@echo "  make compose-down  - Stop monitoring stack"
	@echo "  make compose-logs  - Show stack logs"
	@echo ""
	@echo "$(YELLOW)📊 Monitoring URLs:$(NC)"
	@echo "  • Application:     http://localhost:$(PORT)"
	@echo "  • Health Check:    http://localhost:$(PORT)/health"
	@echo "  • Metrics:         http://localhost:$(PORT)/metrics"
	@echo "  • Prometheus:      http://localhost:9090"
	@echo "  • Grafana:         http://localhost:3000 (admin/admin)"

## 🚀 Development Commands

## Run the monitoring service locally
run:
	@echo "$(GREEN)🚀 Starting monitoring service...$(NC)"
	@export PORT=$(PORT) VERSION=$(VERSION) ENVIRONMENT=$(ENVIRONMENT) && $(GO_RUN) main.go

## Build the monitoring service binary
build:
	@echo "$(GREEN)🔨 Building monitoring service...$(NC)"
	@$(GO_BUILD) -ldflags "-X main.version=$(VERSION)" -o bin/$(APP_NAME) main.go
	@echo "$(GREEN)✅ Build complete: bin/$(APP_NAME)$(NC)"

## Download and tidy dependencies
deps:
	@echo "$(GREEN)📦 Downloading dependencies...$(NC)"
	@$(GO_MOD) download
	@$(GO_MOD) tidy
	@$(GO_MOD) verify

## Clean build artifacts and logs
clean:
	@echo "$(GREEN)🧹 Cleaning build artifacts...$(NC)"
	@rm -rf bin/
	@rm -f *.log
	@rm -f coverage.out coverage.html
	@$(GO_CMD) clean -cache -testcache

## 🧪 Testing Commands

## Run all tests with coverage
test: test-unit test-integration

## Run unit tests only
test-unit:
	@echo "$(GREEN)🧪 Running unit tests...$(NC)"
	@$(GO_TEST) -v -race -short ./internal/... -coverprofile=unit-coverage.out
	@$(GO_CMD) tool cover -html=unit-coverage.out -o unit-coverage.html
	@echo "$(GREEN)✅ Unit tests complete. Coverage: unit-coverage.html$(NC)"

## Run integration tests
test-integration:
	@echo "$(GREEN)🔗 Running integration tests...$(NC)"
	@$(GO_TEST) -v -race -run Integration ./...
	@echo "$(GREEN)✅ Integration tests complete$(NC)"

## Generate detailed test coverage report
test-coverage:
	@echo "$(GREEN)📊 Generating test coverage report...$(NC)"
	@$(GO_TEST) -v -race ./... -coverprofile=coverage.out
	@$(GO_CMD) tool cover -html=coverage.out -o coverage.html
	@$(GO_CMD) tool cover -func=coverage.out
	@echo "$(GREEN)✅ Coverage report generated: coverage.html$(NC)"

## Test health endpoints
health-check:
	@echo "$(GREEN)🩺 Testing health endpoints...$(NC)"
	@echo "Testing liveness probe..."
	@curl -s http://localhost:$(PORT)/health/live | jq . || echo "Service not running"
	@echo ""
	@echo "Testing readiness probe..."
	@curl -s http://localhost:$(PORT)/health/ready | jq . || echo "Service not running"
	@echo ""
	@echo "Testing comprehensive health check..."
	@curl -s http://localhost:$(PORT)/health | jq . || echo "Service not running"

## Generate load for testing metrics
load-test:
	@echo "$(GREEN)🔥 Generating load for testing...$(NC)"
	@echo "Sending 100 requests to demo endpoint..."
	@for i in {1..100}; do \
		curl -s http://localhost:$(PORT)/api/demo > /dev/null & \
	done; wait
	@echo "$(GREEN)✅ Load test complete. Check metrics at http://localhost:$(PORT)/api/metrics$(NC)"

## Run demo scenarios
demo:
	@echo "$(GREEN)🎭 Running demo scenarios...$(NC)"
	@echo ""
	@echo "$(YELLOW)1. Normal request:$(NC)"
	@curl -s http://localhost:$(PORT)/api/demo | jq .
	@echo ""
	@echo "$(YELLOW)2. Bad request (400):$(NC)"
	@curl -s "http://localhost:$(PORT)/api/demo?error=400" | jq .
	@echo ""
	@echo "$(YELLOW)3. Server error (500):$(NC)"
	@curl -s "http://localhost:$(PORT)/api/demo?error=500" | jq .
	@echo ""
	@echo "$(YELLOW)4. Slow request (2s delay):$(NC)"
	@curl -s "http://localhost:$(PORT)/api/demo?delay=2000" | jq .
	@echo ""
	@echo "$(YELLOW)5. Check metrics:$(NC)"
	@curl -s http://localhost:$(PORT)/api/metrics | jq '.request_metrics'

## 🔍 Code Quality

## Run linting checks
lint:
	@echo "$(GREEN)🔍 Running linting checks...$(NC)"
	@$(GO_VET) ./...
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)⚠️  golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
	fi

## Format code with gofmt
format:
	@echo "$(GREEN)🎨 Formatting code...$(NC)"
	@$(GO_FMT) -s -w .
	@if command -v goimports > /dev/null; then \
		goimports -w .; \
	else \
		echo "$(YELLOW)⚠️  goimports not installed. Install with: go install golang.org/x/tools/cmd/goimports@latest$(NC)"; \
	fi

## Run go vet static analysis
vet:
	@echo "$(GREEN)🔎 Running static analysis...$(NC)"
	@$(GO_VET) ./...

## 🐳 Docker Commands

## Build Docker image
docker-build:
	@echo "$(GREEN)🐳 Building Docker image...$(NC)"
	@$(DOCKER_CMD) build -t $(APP_NAME):$(VERSION) .
	@$(DOCKER_CMD) tag $(APP_NAME):$(VERSION) $(APP_NAME):latest

## Run in Docker container
docker-run: docker-build
	@echo "$(GREEN)🐳 Running Docker container...$(NC)"
	@$(DOCKER_CMD) run --rm --name $(APP_NAME) \
		-p $(PORT):$(PORT) \
		-e VERSION=$(VERSION) \
		-e ENVIRONMENT=$(ENVIRONMENT) \
		$(APP_NAME):latest

## Stop Docker container
docker-stop:
	@echo "$(GREEN)🛑 Stopping Docker container...$(NC)"
	@$(DOCKER_CMD) stop $(APP_NAME) || true

## Clean Docker resources
docker-clean:
	@echo "$(GREEN)🧹 Cleaning Docker resources...$(NC)"
	@$(DOCKER_CMD) rmi $(APP_NAME):$(VERSION) $(APP_NAME):latest || true
	@$(DOCKER_CMD) system prune -f

## 🏗️ Compose Commands

## Start full monitoring stack
compose-up:
	@echo "$(GREEN)🏗️ Starting monitoring stack...$(NC)"
	@$(COMPOSE_CMD) up -d
	@echo "$(GREEN)✅ Stack started. Access points:$(NC)"
	@echo "  • Application:     http://localhost:$(PORT)"
	@echo "  • Prometheus:      http://localhost:9090"
	@echo "  • Grafana:         http://localhost:3000 (admin/admin)"

## Stop monitoring stack
compose-down:
	@echo "$(GREEN)🛑 Stopping monitoring stack...$(NC)"
	@$(COMPOSE_CMD) down

## Show stack logs
compose-logs:
	@echo "$(GREEN)📋 Showing stack logs...$(NC)"
	@$(COMPOSE_CMD) logs -f

## 📊 Monitoring Utilities

## Submit custom metric example
custom-metric:
	@echo "$(GREEN)📊 Submitting custom metric...$(NC)"
	@curl -X POST http://localhost:$(PORT)/api/metrics \
		-H "Content-Type: application/json" \
		-d '{"name":"test_metric","type":"counter","value":1,"labels":{"source":"makefile"}}'

## Check system information
system-info:
	@echo "$(GREEN)💻 Getting system information...$(NC)"
	@curl -s http://localhost:$(PORT)/api/system | jq .

## Check application status
status:
	@echo "$(GREEN)📈 Getting application status...$(NC)"
	@curl -s http://localhost:$(PORT)/api/status | jq .

## 🔧 Development Utilities

## Watch for file changes and restart (requires entr)
watch:
	@echo "$(GREEN)👀 Watching for changes... (requires 'entr')$(NC)"
	@if command -v entr > /dev/null; then \
		find . -name "*.go" | entr -r make run; \
	else \
		echo "$(RED)❌ 'entr' not installed. Install with: brew install entr$(NC)"; \
	fi

## Initialize development environment
init:
	@echo "$(GREEN)🎬 Initializing development environment...$(NC)"
	@make deps
	@make format
	@make lint
	@make test-unit
	@echo "$(GREEN)✅ Development environment ready!$(NC)"

## Benchmark testing
bench:
	@echo "$(GREEN)⚡ Running benchmarks...$(NC)"
	@$(GO_TEST) -bench=. -benchmem ./...

## Security check (requires gosec)
security:
	@echo "$(GREEN)🔒 Running security checks...$(NC)"
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "$(YELLOW)⚠️  gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest$(NC)"; \
	fi

## 📋 Information Commands

## Show project information
info:
	@echo "$(GREEN)📋 Project Information:$(NC)"
	@echo "  Name:        $(APP_NAME)"
	@echo "  Version:     $(VERSION)"
	@echo "  Environment: $(ENVIRONMENT)"
	@echo "  Port:        $(PORT)"
	@echo "  Go Version:  $(shell $(GO_CMD) version)"
	@echo "  Git Branch:  $(shell git branch --show-current 2>/dev/null || echo 'unknown')"
	@echo "  Git Commit:  $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"

## Show dependency versions
deps-info:
	@echo "$(GREEN)📦 Dependency Information:$(NC)"
	@$(GO_MOD) list -m all 