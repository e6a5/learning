# 📦 Makefile for HTTP Server Basics (Go Foundations & Package Management)

# Go Module Commands
init:
	go mod init github.com/e6a5/learning/backend/01-http-server

tidy:
	go mod tidy

# Add popular packages (examples for learning)
add-gin:
	go get github.com/gin-gonic/gin

add-testify:
	go get github.com/stretchr/testify

add-cobra:
	go get github.com/spf13/cobra

# Remove a package
remove-package:
	@echo "Example: go mod edit -droprequire=github.com/some/package"
	@echo "Then run: make tidy"

# Development Commands
run:
	go run main.go

run-examples:
	@echo "🎓 Running Go Fundamentals Examples..."
	@echo "The examples are now integrated as a package in the HTTP server"
	@echo "Start the server with 'make run' and visit:"
	@echo "  http://localhost:8080/learn/examples"
	@echo ""
	@echo "Or run individual functions:"
	@echo "  curl http://localhost:8080/learn/examples"

build:
	go build -o server main.go

# Environment setup
env:
	@if [ ! -f .env ]; then cp env.example .env; echo "Created .env file from example"; fi

# Testing and Quality
test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

# View dependencies
deps:
	go list -m all

deps-json:
	go list -json -m all

deps-graph:
	go mod graph

# Module information
mod-info:
	@echo "📦 Module Information:"
	@echo "Module: $$(go list -m)"
	@echo "Go Version: $$(go version)"
	@echo ""
	@echo "📋 Direct Dependencies:"
	@go list -m -f '{{ .Path }} {{ .Version }}' all | grep -v "$$(go list -m)" | head -10
	@echo ""
	@echo "🔍 All Dependencies:"
	@go list -m all | wc -l | awk '{print "Total: " $$1 " modules"}'

# Interactive testing
test-api:
	@echo "🧪 Testing API endpoints:"
	@echo ""
	@echo "1️⃣ Home page:"
	@curl -s http://localhost:8080/ | jq '.' 2>/dev/null || curl -s http://localhost:8080/
	@echo ""
	@echo ""
	@echo "2️⃣ Health check:"
	@curl -s http://localhost:8080/health | jq '.' 2>/dev/null || curl -s http://localhost:8080/health
	@echo ""

test-users:
	@echo "👥 Testing User endpoints:"
	@echo ""
	@echo "1️⃣ Get all users:"
	@curl -s http://localhost:8080/users | jq '.' 2>/dev/null || curl -s http://localhost:8080/users
	@echo ""
	@echo ""
	@echo "2️⃣ Create a user:"
	@curl -s -X POST http://localhost:8080/users \
		-H "Content-Type: application/json" \
		-d '{"name":"John Doe","email":"john@example.com"}' | jq '.' 2>/dev/null || curl -s -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name":"John Doe","email":"john@example.com"}'
	@echo ""
	@echo ""
	@echo "3️⃣ Get user by ID:"
	@curl -s http://localhost:8080/users/1 | jq '.' 2>/dev/null || curl -s http://localhost:8080/users/1
	@echo ""

test-learning:
	@echo "🎓 Testing Learning endpoints:"
	@echo ""
	@echo "1️⃣ Go Basics:"
	@curl -s http://localhost:8080/learn/basics | jq '.data.variables' 2>/dev/null || echo "Install jq for formatted output"
	@echo ""
	@echo "2️⃣ Packages:"
	@curl -s http://localhost:8080/learn/packages | jq '.data.third_party_used' 2>/dev/null || echo "Install jq for formatted output"
	@echo ""

# Package management workflow
package-workflow:
	@echo "📦 Go Package Management Workflow Demo:"
	@echo ""
	@echo "1️⃣ Initialize module (if not exists):"
	@echo "   go mod init github.com/yourusername/project"
	@echo ""
	@echo "2️⃣ Add dependencies:"
	@echo "   go get github.com/gorilla/mux"
	@echo "   go get github.com/sirupsen/logrus@latest"
	@echo ""
	@echo "3️⃣ Clean up dependencies:"
	@echo "   go mod tidy"
	@echo ""
	@echo "4️⃣ View dependencies:"
	@echo "   go list -m all"
	@echo ""
	@echo "5️⃣ Update dependencies:"
	@echo "   go get -u ./..."
	@echo ""
	@echo "Current module: $$(go list -m)"

# Learning examples
learn-go:
	@echo "🎓 Learn Go Fundamentals:"
	@echo ""
	@echo "Run interactive Go examples:"
	@echo "   make run-examples"
	@echo ""
	@echo "Available learning paths:"
	@echo "   1. Go Basics: Variables, types, control structures"
	@echo "   2. Functions: Multiple returns, error handling"
	@echo "   3. Structs: Data structures and methods"
	@echo "   4. Collections: Slices and maps"
	@echo "   5. Pointers: Memory management"
	@echo "   6. Error Handling: Go's approach to errors"

# Full demo
demo:
	@echo "🚀 Starting HTTP Server Demo..."
	@make env
	@echo "⏳ Starting server in background..."
	@make run &
	@sleep 3
	@echo ""
	@make test-api
	@echo ""
	@make test-users
	@echo ""
	@make test-learning
	@echo ""
	@echo "🎉 Demo completed! Check the server logs above."
	@echo "💡 Press Ctrl+C to stop the server"

# Development helpers
clean:
	go clean
	rm -f server

# Quick setup for new learners
setup:
	@echo "🎓 Setting up Go HTTP Server Learning Environment..."
	@make env
	@make tidy
	@echo ""
	@echo "✅ Setup complete! Now you can:"
	@echo "   1. Run the server: make run"
	@echo "   2. Test the API: make test-api"
	@echo "   3. Learn about modules: make mod-info"
	@echo "   4. See package workflow: make package-workflow"
	@echo "   5. Learn Go basics: make run-examples"

help:
	@echo "Available commands:"
	@echo ""
	@echo "📦 Module Management:"
	@echo "  init          - Initialize Go module"
	@echo "  tidy          - Clean up dependencies"
	@echo "  deps          - View all dependencies"
	@echo "  mod-info      - Show module information"
	@echo ""
	@echo "🔧 Development:"
	@echo "  run           - Run the HTTP server"
	@echo "  run-examples  - Run Go fundamentals examples"
	@echo "  build         - Build binary"
	@echo "  env           - Create .env file from example"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  test-api      - Test basic API endpoints"
	@echo "  test-users    - Test user management"
	@echo "  test-learning - Test learning endpoints"
	@echo "  demo          - Full interactive demo"
	@echo ""
	@echo "📚 Learning:"
	@echo "  learn-go      - Show Go learning paths"
	@echo "  package-workflow - Show package management workflow"
	@echo "  setup         - Complete setup for new learners" 