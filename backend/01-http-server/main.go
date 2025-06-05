package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	// Third-party packages from GitHub
	"github.com/e6a5/learning/backend/01-http-server/basics"
	"github.com/gorilla/mux"     // Popular HTTP router
	"github.com/joho/godotenv"   // Environment variable loader
	"github.com/sirupsen/logrus" // Structured logging
	// Local package demonstrating Go fundamentals
)

// User represents a user in our system
// This demonstrates Go struct basics
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	JoinedAt string `json:"joined_at"`
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// In-memory storage (slice of pointers to User structs)
var users []*User
var nextID int = 1

// Initialize sample data
func init() {
	// This function runs before main()
	users = append(users, &User{
		ID:       1,
		Name:     "Alice Johnson",
		Email:    "alice@example.com",
		JoinedAt: time.Now().Format("2006-01-02 15:04:05"),
	})
	nextID = 2
}

func main() {
	// Load environment variables from .env file
	// This shows how to use third-party packages
	if err := godotenv.Load(); err != nil {
		logrus.Info("No .env file found, using defaults")
	}

	// Configure structured logging
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// Get port from environment variable or use default
	port := getEnv("PORT", "8080")

	// Create a new router using gorilla/mux (third-party package)
	router := mux.NewRouter()

	// Add middleware for logging requests
	router.Use(loggingMiddleware)

	// Define routes with different HTTP methods
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/health", healthHandler).Methods("GET")
	router.HandleFunc("/users", getUsersHandler).Methods("GET")
	router.HandleFunc("/users", createUserHandler).Methods("POST")
	router.HandleFunc("/users/{id:[0-9]+}", getUserHandler).Methods("GET")

	// Static route for learning about routing
	router.HandleFunc("/learn/basics", basicsHandler).Methods("GET")
	router.HandleFunc("/learn/packages", packagesHandler).Methods("GET")
	router.HandleFunc("/learn/modules", modulesHandler).Methods("GET")
	router.HandleFunc("/learn/examples", runExamplesHandler).Methods("GET")

	logrus.WithFields(logrus.Fields{
		"port":    port,
		"version": "1.0.0",
	}).Info("🚀 HTTP Server starting")

	// Start the server
	logrus.Fatal(http.ListenAndServe(":"+port, router))
}

// Middleware function to log all requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the request
		logrus.WithFields(logrus.Fields{
			"method":   r.Method,
			"path":     r.URL.Path,
			"duration": time.Since(start),
			"ip":       r.RemoteAddr,
		}).Info("Request processed")
	})
}

// Home handler - demonstrates basic HTTP response
func homeHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
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

	sendJSONResponse(w, http.StatusOK, response)
}

// Health check handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Success: true,
		Message: "Server is healthy",
		Data: map[string]interface{}{
			"status":    "UP",
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
			"uptime":    "running",
		},
	}

	sendJSONResponse(w, http.StatusOK, response)
}

// Get all users - demonstrates slice operations
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Success: true,
		Message: fmt.Sprintf("Found %d users", len(users)),
		Data:    users,
	}

	sendJSONResponse(w, http.StatusOK, response)
}

// Create a new user - demonstrates JSON parsing and slice append
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User

	// Parse JSON from request body
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		response := Response{
			Success: false,
			Message: "Invalid JSON format",
		}
		sendJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Validate required fields
	if newUser.Name == "" || newUser.Email == "" {
		response := Response{
			Success: false,
			Message: "Name and email are required",
		}
		sendJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Set ID and timestamp
	newUser.ID = nextID
	nextID++
	newUser.JoinedAt = time.Now().Format("2006-01-02 15:04:05")

	// Add to slice
	users = append(users, &newUser)

	logrus.WithFields(logrus.Fields{
		"user_id": newUser.ID,
		"name":    newUser.Name,
		"email":   newUser.Email,
	}).Info("New user created")

	response := Response{
		Success: true,
		Message: "User created successfully",
		Data:    newUser,
	}

	sendJSONResponse(w, http.StatusCreated, response)
}

// Get user by ID - demonstrates URL parameters and error handling
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL using gorilla/mux
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert string to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid user ID",
		}
		sendJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Find user in slice
	for _, user := range users {
		if user.ID == id {
			response := Response{
				Success: true,
				Message: "User found",
				Data:    user,
			}
			sendJSONResponse(w, http.StatusOK, response)
			return
		}
	}

	// User not found
	response := Response{
		Success: false,
		Message: "User not found",
	}
	sendJSONResponse(w, http.StatusNotFound, response)
}

// Learning endpoints - teach Go concepts

func basicsHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
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

	sendJSONResponse(w, http.StatusOK, response)
}

func packagesHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
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

	sendJSONResponse(w, http.StatusOK, response)
}

func modulesHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
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

	sendJSONResponse(w, http.StatusOK, response)
}

// runExamplesHandler demonstrates using a local package
func runExamplesHandler(w http.ResponseWriter, r *http.Request) {
	// Capture the output from basics.RunAllExamples()
	response := Response{
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

	sendJSONResponse(w, http.StatusOK, response)
}

// Helper function to send JSON responses
func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logrus.WithError(err).Error("Failed to encode JSON response")
	}
}

// Helper function to get environment variable with default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
