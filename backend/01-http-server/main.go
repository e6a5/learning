package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/e6a5/learning/backend/01-http-server/internal/handlers"
	"github.com/e6a5/learning/backend/01-http-server/internal/middleware"
	"github.com/e6a5/learning/backend/01-http-server/internal/repository"
	"github.com/e6a5/learning/backend/01-http-server/internal/utils"
)

func main() {
	// Initialize application
	setupLogging()

	// Initialize dependencies
	userRepo := repository.NewUserRepository()
	userHandler := handlers.NewUserHandler(userRepo)
	learnHandler := handlers.NewLearnHandler()

	// Setup HTTP server
	router := setupRoutes(userHandler, learnHandler)
	port := utils.GetEnv("PORT", "8080")

	logrus.WithFields(logrus.Fields{
		"port":    port,
		"version": "1.0.0",
	}).Info("🚀 HTTP Server starting")

	// Start the server
	logrus.Fatal(http.ListenAndServe(":"+port, router))
}

func setupLogging() {
	if err := godotenv.Load(); err != nil {
		logrus.Info("No .env file found, using defaults")
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
}

func setupRoutes(userHandler *handlers.UserHandler, learnHandler *handlers.LearnHandler) *mux.Router {
	router := mux.NewRouter()

	// Apply middleware
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.CORSMiddleware)

	// System routes
	router.HandleFunc("/", learnHandler.Home).Methods("GET")
	router.HandleFunc("/health", learnHandler.Health).Methods("GET")

	// User routes
	router.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	router.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id:[0-9]+}", userHandler.GetUser).Methods("GET")

	// Learning routes
	router.HandleFunc("/learn/basics", learnHandler.Basics).Methods("GET")
	router.HandleFunc("/learn/packages", learnHandler.Packages).Methods("GET")
	router.HandleFunc("/learn/modules", learnHandler.Modules).Methods("GET")
	router.HandleFunc("/learn/examples", learnHandler.Examples).Methods("GET")

	return router
}
