package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/e6a5/learning/backend/08-monitoring/internal/handlers"
	"github.com/e6a5/learning/backend/08-monitoring/internal/middleware"
	"github.com/e6a5/learning/backend/08-monitoring/internal/repository"
)

func main() {
	// Configuration from environment
	port := getEnv("PORT", "8080")
	version := getEnv("VERSION", "1.0.0")
	environment := getEnv("ENVIRONMENT", "development")

	log.Printf("Starting monitoring service version %s in %s environment", version, environment)

	// Initialize dependencies
	metricsRepo := repository.NewMetricsRepository(version, environment)

	// Set up health checkers
	healthCheckers := []repository.HealthChecker{
		repository.NewDatabaseHealthChecker("database", "mysql://localhost:3306"),
		repository.NewExternalServiceHealthChecker("api", "https://httpbin.org/status/200"),
	}

	// Initialize handlers
	monitoringHandler := handlers.NewMonitoringHandler(metricsRepo, healthCheckers)

	// Initialize middleware
	monitoringMiddleware := middleware.NewMonitoringMiddleware(metricsRepo)

	// Setup routes
	router := setupRoutes(monitoringHandler, monitoringMiddleware)

	// Start server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func setupRoutes(handler *handlers.MonitoringHandler, monitoringMW *middleware.MonitoringMiddleware) *mux.Router {
	router := mux.NewRouter()

	// Apply global middleware
	router.Use(middleware.CorsMiddleware)
	router.Use(middleware.LoggingMiddleware)
	router.Use(monitoringMW.Wrap)

	// Health check endpoints (no monitoring to avoid recursive metrics)
	healthRouter := router.PathPrefix("/health").Subrouter()
	healthRouter.HandleFunc("", handler.HealthCheck).Methods("GET")
	healthRouter.HandleFunc("/live", handler.LivenessCheck).Methods("GET")
	healthRouter.HandleFunc("/ready", handler.ReadinessCheck).Methods("GET")

	// Metrics endpoints
	router.HandleFunc("/metrics", handler.GetMetrics).Methods("GET")

	// API endpoints
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/metrics", handler.GetCustomMetrics).Methods("GET")
	apiRouter.HandleFunc("/metrics", handler.PostCustomMetric).Methods("POST")
	apiRouter.HandleFunc("/system", handler.GetSystemInfo).Methods("GET")
	apiRouter.HandleFunc("/status", handler.GetStatus).Methods("GET")
	apiRouter.HandleFunc("/demo", handler.DemoEndpoint).Methods("GET")

	return router
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
