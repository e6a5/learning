package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/e6a5/learning/backend/03-redis-intro/internal/models"
)

// RespondJSON sends a JSON response with the given status code and data
func RespondJSON(w http.ResponseWriter, statusCode int, data models.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

// GetEnv gets an environment variable with a default value
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
