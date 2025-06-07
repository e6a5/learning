package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/e6a5/learning/backend/07-error-handling/internal/circuit"
	"github.com/e6a5/learning/backend/07-error-handling/internal/handlers"
	"github.com/e6a5/learning/backend/07-error-handling/internal/middleware"
	"github.com/e6a5/learning/backend/07-error-handling/internal/models"
	"github.com/e6a5/learning/backend/07-error-handling/internal/retry"
)

// App holds application dependencies - small, focused
type App struct {
	db             *sql.DB
	redis          *redis.Client
	dbCircuit      *circuit.Breaker
	redisCircuit   *circuit.Breaker
	userCache      map[int]*models.User
	cacheMutex     sync.RWMutex
	requestCounter int64
	counterMutex   sync.Mutex
}

func main() {
	// Load environment and configure logging
	setupLogging()

	// Initialize application with dependencies
	app := &App{
		userCache:    make(map[int]*models.User),
		dbCircuit:    circuit.New("database", 5, 30*time.Second),
		redisCircuit: circuit.New("redis", 3, 15*time.Second),
	}

	// Initialize databases with retry logic
	if err := app.initializeDependencies(); err != nil {
		logrus.WithError(err).Warn("Failed to initialize some dependencies, continuing with degraded functionality")
	}

	// Setup HTTP server
	router := app.setupRoutes()
	port := getEnv("PORT", "8080")

	logrus.WithFields(logrus.Fields{
		"port":    port,
		"version": "1.0.0",
	}).Info("🔥 Error Handling Server starting")

	// Start server
	if err := http.ListenAndServe(":"+port, router); err != nil {
		logrus.WithError(err).Fatal("Server failed to start")
	}
}

func setupLogging() {
	if err := godotenv.Load(); err != nil {
		logrus.Info("No .env file found, using defaults")
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
}

func (app *App) setupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Apply middleware chain
	router.Use(middleware.PanicRecovery(app.sendErrorResponse))
	router.Use(middleware.RequestID(&app.requestCounter, &app.counterMutex))
	router.Use(middleware.Logging())
	router.Use(middleware.RateLimit())

	// Initialize handlers
	userHandler := handlers.NewUserHandler(
		app.sendJSONResponse,
		app.sendErrorResponse,
		app.sendErrorResponseWithFallback,
	)

	// API routes
	router.HandleFunc("/", app.homeHandler).Methods("GET")
	router.HandleFunc("/health", app.healthHandler).Methods("GET")

	// User routes with dependency injection
	router.HandleFunc("/users", userHandler.GetUsers(app.dbCircuit.Call, app.userCache)).Methods("GET")
	router.HandleFunc("/users", userHandler.CreateUser(app.dbCircuit.Call, app.userCache)).Methods("POST")
	router.HandleFunc("/users/{id:[0-9]+}", userHandler.GetUser(app.dbCircuit.Call, app.userCache)).Methods("GET")

	// Error simulation routes
	router.HandleFunc("/simulate/panic", app.simulatePanicHandler).Methods("GET")
	router.HandleFunc("/simulate/db-error", app.simulateDBErrorHandler).Methods("GET")
	router.HandleFunc("/simulate/validation-error", app.simulateValidationErrorHandler).Methods("POST")

	// Circuit breaker management
	router.HandleFunc("/circuit-breaker/status", app.circuitBreakerStatusHandler).Methods("GET")
	router.HandleFunc("/circuit-breaker/reset", app.resetCircuitBreakersHandler).Methods("POST")

	return router
}

func (app *App) initializeDependencies() error {
	var errors []error

	// Initialize MySQL with retry
	if err := app.initializeMySQL(); err != nil {
		errors = append(errors, err)
	}

	// Initialize Redis with retry
	if err := app.initializeRedis(); err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return errors[0] // Return first error for simplicity
	}
	return nil
}

func (app *App) initializeMySQL() error {
	config := models.RetryConfig{
		MaxAttempts:   5,
		BaseDelay:     1 * time.Second,
		MaxDelay:      30 * time.Second,
		BackoffFactor: 2.0,
		Jitter:        true,
	}

	return retry.WithRetry("mysql-connection", config, func() error {
		dsn := getEnv("DB_DSN", "user:password@tcp(localhost:3306)/testdb")
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return err
		}

		if err := db.Ping(); err != nil {
			db.Close()
			return err
		}

		app.db = db
		logrus.Info("MySQL connection established")
		return nil
	})
}

func (app *App) initializeRedis() error {
	config := models.RetryConfig{
		MaxAttempts:   3,
		BaseDelay:     500 * time.Millisecond,
		MaxDelay:      10 * time.Second,
		BackoffFactor: 2.0,
		Jitter:        true,
	}

	return retry.WithRetry("redis-connection", config, func() error {
		addr := getEnv("REDIS_ADDR", "localhost:6379")
		client := redis.NewClient(&redis.Options{Addr: addr})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Ping(ctx).Err(); err != nil {
			client.Close()
			return err
		}

		app.redis = client
		logrus.Info("Redis connection established")
		return nil
	})
}

// Simple handlers that focus on HTTP concerns only
func (app *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	response := models.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message":     "Welcome to Error Handling Learning Lab!",
			"server_time": time.Now(),
			"endpoints": []string{
				"GET /", "GET /health", "GET /users", "POST /users", "GET /users/{id}",
				"GET /simulate/panic", "GET /simulate/db-error", "POST /simulate/validation-error",
				"GET /circuit-breaker/status", "POST /circuit-breaker/reset",
			},
		},
	}
	app.sendJSONResponse(w, http.StatusOK, response)
}

func (app *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	health := app.buildHealthResponse()
	response := models.APIResponse{Success: true, Data: health}
	app.sendJSONResponse(w, http.StatusOK, response)
}

func (app *App) buildHealthResponse() map[string]interface{} {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"services":  make(map[string]interface{}),
	}

	// Check database health
	if app.db != nil && app.db.Ping() == nil {
		health["services"].(map[string]interface{})["database"] = map[string]interface{}{"status": "healthy"}
	} else {
		health["services"].(map[string]interface{})["database"] = map[string]interface{}{"status": "unavailable"}
	}

	// Check Redis health
	if app.redis != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if app.redis.Ping(ctx).Err() == nil {
			health["services"].(map[string]interface{})["redis"] = map[string]interface{}{"status": "healthy"}
		} else {
			health["services"].(map[string]interface{})["redis"] = map[string]interface{}{"status": "unhealthy"}
		}
	}

	// Add circuit breaker status
	health["circuit_breakers"] = map[string]interface{}{
		"database": map[string]interface{}{
			"state":    app.dbCircuit.GetState(),
			"failures": app.dbCircuit.GetFailures(),
		},
		"redis": map[string]interface{}{
			"state":    app.redisCircuit.GetState(),
			"failures": app.redisCircuit.GetFailures(),
		},
	}

	return health
}

// Error simulation handlers - focused on single responsibility
func (app *App) simulatePanicHandler(w http.ResponseWriter, r *http.Request) {
	logrus.WithField("request_id", r.Header.Get("X-Request-ID")).Info("Simulating panic")
	panic("This is a simulated panic for testing recovery")
}

func (app *App) simulateDBErrorHandler(w http.ResponseWriter, r *http.Request) {
	// Force circuit breaker to open
	for i := 0; i < 6; i++ {
		app.dbCircuit.Call(func() error { return nil }) // Simulate failures
	}

	response := models.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message":       "Database circuit breaker should now be open",
			"circuit_state": app.dbCircuit.GetState(),
		},
	}
	app.sendJSONResponse(w, http.StatusOK, response)
}

func (app *App) simulateValidationErrorHandler(w http.ResponseWriter, r *http.Request) {
	app.sendErrorResponse(w, models.APIError{
		Type:      models.ValidationError,
		Code:      "SIMULATED_VALIDATION_ERROR",
		Message:   "This is a simulated validation error",
		Details:   map[string]interface{}{"field": "test_field", "value": "invalid_value"},
		RequestID: r.Header.Get("X-Request-ID"),
		Timestamp: time.Now(),
		Retryable: false,
	}, http.StatusBadRequest)
}

func (app *App) circuitBreakerStatusHandler(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"database": map[string]interface{}{
			"state":         app.dbCircuit.GetState(),
			"failures":      app.dbCircuit.GetFailures(),
			"last_failure":  app.dbCircuit.GetLastFailTime(),
			"success_count": app.dbCircuit.GetSuccessCount(),
		},
		"redis": map[string]interface{}{
			"state":         app.redisCircuit.GetState(),
			"failures":      app.redisCircuit.GetFailures(),
			"last_failure":  app.redisCircuit.GetLastFailTime(),
			"success_count": app.redisCircuit.GetSuccessCount(),
		},
	}

	response := models.APIResponse{Success: true, Data: status}
	app.sendJSONResponse(w, http.StatusOK, response)
}

func (app *App) resetCircuitBreakersHandler(w http.ResponseWriter, r *http.Request) {
	app.dbCircuit.Reset()
	app.redisCircuit.Reset()
	logrus.Info("Circuit breakers reset")

	response := models.APIResponse{
		Success: true,
		Data:    map[string]interface{}{"message": "All circuit breakers have been reset"},
	}
	app.sendJSONResponse(w, http.StatusOK, response)
}

// HTTP utility functions - focused on HTTP concerns
func (app *App) sendJSONResponse(w http.ResponseWriter, statusCode int, data models.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logrus.WithError(err).Error("Failed to encode JSON response")
	}
}

func (app *App) sendErrorResponse(w http.ResponseWriter, apiError models.APIError, statusCode int) {
	response := models.APIResponse{Success: false, Error: &apiError}
	app.sendJSONResponse(w, statusCode, response)
}

func (app *App) sendErrorResponseWithFallback(w http.ResponseWriter, apiError models.APIError, fallbackData interface{}, statusCode int) {
	response := models.APIResponse{Success: false, Error: &apiError, FallbackData: fallbackData}
	app.sendJSONResponse(w, statusCode, response)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
