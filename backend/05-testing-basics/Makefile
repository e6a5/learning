# 🧪 Makefile for Testing Validation Functions

.PHONY: help test coverage bench clean fmt vet deps

# Default target
help:
	@echo "🧪 Testing Basics - Available Commands:"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  test            - Run all unit tests"
	@echo "  test-verbose    - Run tests with verbose output"
	@echo ""
	@echo "📊 Coverage:"
	@echo "  coverage        - Run tests with coverage report"
	@echo "  coverage-html   - Generate HTML coverage report"
	@echo ""
	@echo "🚀 Performance:"
	@echo "  bench           - Run benchmark tests"
	@echo ""
	@echo "🔧 Development:"
	@echo "  fmt             - Format code"
	@echo "  vet             - Run go vet"
	@echo "  deps            - Download dependencies"
	@echo "  clean           - Clean test artifacts"

# Basic testing
test:
	@echo "🧪 Running unit tests..."
	go test ./models -v

test-verbose:
	@echo "🧪 Running tests with verbose output..."
	go test ./models -v -count=1

# Coverage testing
coverage:
	@echo "📊 Running tests with coverage..."
	go test ./models -coverprofile=coverage.out
	go tool cover -func=coverage.out

coverage-html: coverage
	@echo "📊 Generating HTML coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"
	@echo "💡 Open with: open coverage.html"

# Benchmark testing
bench:
	@echo "🚀 Running benchmark tests..."
	go test ./models -bench=. -benchmem

# Development commands
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

vet:
	@echo "🔍 Running go vet..."
	go vet ./...

deps:
	@echo "📦 Downloading dependencies..."
	go mod download
	go mod tidy

# Clean up
clean:
	@echo "🧹 Cleaning up test artifacts..."
	rm -f coverage.out coverage.html
	go clean -testcache
	@echo "✅ Cleanup complete"

# Quick development workflow
dev: fmt vet test
	@echo "✅ Development checks passed"

# Show what we actually test
show-tests:
	@echo "📋 Test Functions:"
	@grep "^func Test" models/user_test.go
	@echo ""
	@echo "📋 Benchmark Functions:"
	@grep "^func Benchmark" models/user_test.go 