package main

import (
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"

	"github.com/e6a5/learning/backend/03-redis-intro/internal/handlers"
	"github.com/e6a5/learning/backend/03-redis-intro/internal/repository"
	"github.com/e6a5/learning/backend/03-redis-intro/internal/utils"
)

func main() {
	// Initialize Redis connection
	redisClient, err := initializeRedis()
	if err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}

	// Initialize dependencies
	cacheRepo := repository.NewCacheRepository(redisClient)
	cacheHandler := handlers.NewCacheHandler(cacheRepo)

	// Setup HTTP server
	router := setupRoutes(cacheHandler)
	port := utils.GetEnv("PORT", "8080")

	log.Println("🚀 Redis Server running at http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func initializeRedis() (*redis.Client, error) {
	addr := utils.GetEnv("REDIS_ADDR", "redis:6379")
	password := utils.GetEnv("REDIS_PASSWORD", "")
	db := 0 // Default database

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test connection
	if err := repository.NewCacheRepository(client).Ping(); err != nil {
		return nil, err
	}

	return client, nil
}

func setupRoutes(cacheHandler *handlers.CacheHandler) *mux.Router {
	router := mux.NewRouter()

	// Key-Value operations
	router.HandleFunc("/cache/{key}", cacheHandler.GetValue).Methods("GET")
	router.HandleFunc("/cache", cacheHandler.SetValue).Methods("POST")
	router.HandleFunc("/cache/{key}", cacheHandler.DeleteValue).Methods("DELETE")
	router.HandleFunc("/cache", cacheHandler.GetAllKeys).Methods("GET")

	// Cache operations
	router.HandleFunc("/cache/{key}/ttl", cacheHandler.GetTTL).Methods("GET")
	router.HandleFunc("/cache/{key}/expire", cacheHandler.SetExpire).Methods("POST")

	// Health check
	router.HandleFunc("/health", cacheHandler.HealthCheck).Methods("GET")

	return router
}
