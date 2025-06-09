package handlers

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/e6a5/learning/backend/01-http-server/basics"
	"github.com/e6a5/learning/backend/01-http-server/internal/models"
	"github.com/e6a5/learning/backend/01-http-server/internal/utils"
)

// LearnHandler handles learning-related HTTP requests
type LearnHandler struct{}

// NewLearnHandler creates a new learn handler
func NewLearnHandler() *LearnHandler {
	return &LearnHandler{}
}

// Home handles GET / - welcome page
func (h *LearnHandler) Home(w http.ResponseWriter, r *http.Request) {
	response := models.Response{
		Success: true,
		Message: "Welcome to Go HTTP Server Learning Lab!",
		Data: map[string]interface{}{
			"server_time": time.Now().Format("2006-01-02 15:04:05"),
			"go_version":  "1.23.4",
			"endpoints": []string{
				"GET /",
				"GET /health",
				"GET /users",
				"POST /users",
				"GET /users/{id}",
				"GET /learn/basics",
				"GET /learn/packages",
				"GET /learn/modules",
				"GET /learn/examples",
			},
		},
	}

	utils.SendJSONResponse(w, http.StatusOK, response)
}

// Health handles GET /health - health check
func (h *LearnHandler) Health(w http.ResponseWriter, r *http.Request) {
	response := models.Response{
		Success: true,
		Message: "Server is healthy",
		Data: map[string]interface{}{
			"status":    "UP",
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
			"uptime":    "running",
		},
	}

	utils.SendJSONResponse(w, http.StatusOK, response)
}

// Basics handles GET /learn/basics - Go basics tutorial
func (h *LearnHandler) Basics(w http.ResponseWriter, r *http.Request) {
	response := models.Response{
		Success: true,
		Message: "Go Basics Tutorial",
		Data: map[string]interface{}{
			"variables": map[string]string{
				"var name string":    "Explicit type declaration",
				"name := \"John\"":   "Short variable declaration",
				"const Pi = 3.14159": "Constant declaration",
			},
			"types": []string{
				"string", "int", "int64", "float64", "bool",
				"[]string (slice)", "map[string]int", "*User (pointer)",
			},
			"control_structures": map[string]string{
				"if/else": "Conditional execution",
				"for":     "Loops (only loop in Go)",
				"switch":  "Multi-way branching",
				"range":   "Iterate over slices, maps, channels",
			},
			"functions": map[string]string{
				"Multiple return values": "func divide(a, b int) (int, error)",
				"Named return values":    "func calc() (sum, diff int)",
				"Variadic functions":     "func printf(format string, args ...interface{})",
			},
		},
	}

	utils.SendJSONResponse(w, http.StatusOK, response)
}

// Packages handles GET /learn/packages - Go packages tutorial
func (h *LearnHandler) Packages(w http.ResponseWriter, r *http.Request) {
	response := models.Response{
		Success: true,
		Message: "Go Packages & Third-party Libraries",
		Data: map[string]interface{}{
			"standard_library": map[string]string{
				"fmt":           "Formatted I/O (print, sprintf)",
				"net/http":      "HTTP client and server",
				"encoding/json": "JSON encoding/decoding",
				"time":          "Time and duration",
				"os":            "Operating system interface",
				"strconv":       "String conversions",
			},
			"third_party_used": map[string]string{
				"github.com/gorilla/mux":     "Powerful HTTP router with URL variables",
				"github.com/sirupsen/logrus": "Structured logging with levels and fields",
				"github.com/joho/godotenv":   "Load environment variables from .env file",
			},
			"popular_packages": map[string]string{
				"github.com/gin-gonic/gin":    "High-performance HTTP web framework",
				"github.com/lib/pq":           "PostgreSQL driver",
				"github.com/go-redis/redis":   "Redis client",
				"github.com/stretchr/testify": "Testing toolkit with assertions",
				"github.com/spf13/cobra":      "CLI application framework",
			},
		},
	}

	utils.SendJSONResponse(w, http.StatusOK, response)
}

// Modules handles GET /learn/modules - Go modules tutorial
func (h *LearnHandler) Modules(w http.ResponseWriter, r *http.Request) {
	response := models.Response{
		Success: true,
		Message: "Go Modules Tutorial",
		Data: map[string]interface{}{
			"what_are_modules": "Go modules are collections of Go packages stored in a file tree with a go.mod file at its root",
			"commands": map[string]string{
				"go mod init <module-name>": "Initialize a new module",
				"go mod tidy":               "Add missing and remove unused modules",
				"go get <package>":          "Add or update a dependency",
				"go get <package>@version":  "Get a specific version",
				"go mod download":           "Download modules to local cache",
				"go mod verify":             "Verify dependencies have expected content",
				"go list -m all":            "View all dependencies",
			},
			"go_mod_file_structure": map[string]string{
				"module":  "Declares the module path",
				"go":      "Sets the expected Go language version",
				"require": "Lists required dependencies with versions",
				"exclude": "Excludes specific versions of modules",
				"replace": "Replaces a module with another",
			},
			"version_examples": map[string]string{
				"v1.2.3":                   "Exact version",
				"latest":                   "Latest version",
				"v1.2.0-beta.1":            "Pre-release version",
				"v0.0.0-commitdate-commit": "Pseudo-version",
			},
			"best_practices": []string{
				"Use semantic versioning",
				"Keep dependencies minimal",
				"Regularly update dependencies",
				"Use go mod tidy to clean up",
				"Commit go.mod and go.sum files",
			},
		},
	}

	utils.SendJSONResponse(w, http.StatusOK, response)
}

// Examples handles GET /learn/examples - runs Go examples
func (h *LearnHandler) Examples(w http.ResponseWriter, r *http.Request) {
	response := models.Response{
		Success: true,
		Message: "Running Go Fundamentals Examples from local 'basics' package",
		Data: map[string]interface{}{
			"note": "This endpoint demonstrates importing and using a local Go package",
			"package_structure": map[string]string{
				"package":  "basics",
				"location": "./basics/go-basics.go",
				"import":   "github.com/e6a5/learning/backend/01-http-server/basics",
			},
			"available_functions": []string{
				"basics.RunAllExamples() - Run all demonstrations",
				"basics.DemonstrateVariables() - Variables and types",
				"basics.DemonstrateControlStructures() - If/else, loops, switch",
				"basics.DemonstrateFunctions() - Multiple returns, errors",
				"basics.DemonstrateStructs() - Struct definition and usage",
				"basics.DemonstrateCollections() - Slices and maps",
				"basics.DemonstratePointers() - Memory management",
				"basics.DemonstrateErrorHandling() - Error patterns",
			},
			"example_usage": map[string]string{
				"import":     "import \"github.com/e6a5/learning/backend/01-http-server/basics\"",
				"call":       "basics.RunAllExamples()",
				"individual": "basics.DemonstrateVariables()",
			},
			"tip": "Check the terminal/logs to see the actual examples output when this endpoint is called",
		},
	}

	// Actually run the examples (output will go to terminal/logs)
	logrus.Info("Running Go fundamentals examples from basics package...")
	basics.RunAllExamples()
	logrus.Info("Go fundamentals examples completed")

	utils.SendJSONResponse(w, http.StatusOK, response)
}
