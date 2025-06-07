package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/e6a5/learning/backend/07-error-handling/internal/models"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	sendJSONResponse              func(http.ResponseWriter, int, models.APIResponse)
	sendErrorResponse             func(http.ResponseWriter, models.APIError, int)
	sendErrorResponseWithFallback func(http.ResponseWriter, models.APIError, interface{}, int)
}

// NewUserHandler creates a new user handler
func NewUserHandler(
	sendJSONResponse func(http.ResponseWriter, int, models.APIResponse),
	sendErrorResponse func(http.ResponseWriter, models.APIError, int),
	sendErrorResponseWithFallback func(http.ResponseWriter, models.APIError, interface{}, int),
) *UserHandler {
	return &UserHandler{
		sendJSONResponse:              sendJSONResponse,
		sendErrorResponse:             sendErrorResponse,
		sendErrorResponseWithFallback: sendErrorResponseWithFallback,
	}
}

// GetUsers handles GET /users requests with circuit breaker and fallback
func (h *UserHandler) GetUsers(dbCall func(func() error) error, userCache map[int]*models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users []models.User

		// Try to get users from database with circuit breaker
		err := dbCall(func() error {
			// This would be injected database connection in real implementation
			// For now, simulating the database call structure
			return fmt.Errorf("database connection not available in handler")
		})

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":      err.Error(),
				"request_id": r.Header.Get("X-Request-ID"),
			}).Warn("Failed to fetch users from database, using fallback")

			// Use cached data as fallback
			var cachedUsers []models.User
			for _, user := range userCache {
				cachedUsers = append(cachedUsers, *user)
			}

			fallbackData := map[string]interface{}{
				"users":      cachedUsers,
				"cache_info": "Data from local cache due to database unavailability",
				"cache_age":  "unknown",
			}

			h.sendErrorResponseWithFallback(w, models.APIError{
				Type:      models.ServiceUnavailable,
				Code:      "DATABASE_UNAVAILABLE",
				Message:   "Unable to fetch latest users, showing cached data",
				RequestID: r.Header.Get("X-Request-ID"),
				Timestamp: time.Now(),
				Retryable: true,
			}, fallbackData, http.StatusPartialContent)
			return
		}

		response := models.APIResponse{
			Success: true,
			Data: map[string]interface{}{
				"users": users,
				"count": len(users),
			},
		}

		h.sendJSONResponse(w, http.StatusOK, response)
	}
}

// CreateUser handles POST /users requests with validation
func (h *UserHandler) CreateUser(dbCall func(func() error) error, userCache map[int]*models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		// Parse and validate input
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			h.sendErrorResponse(w, models.APIError{
				Type:      models.ValidationError,
				Code:      "INVALID_JSON",
				Message:   "Request body contains invalid JSON",
				Details:   map[string]interface{}{"error": err.Error()},
				RequestID: r.Header.Get("X-Request-ID"),
				Timestamp: time.Now(),
				Retryable: false,
			}, http.StatusBadRequest)
			return
		}

		// Validate required fields
		if err := validateUser(&user); err != nil {
			err.RequestID = r.Header.Get("X-Request-ID")
			h.sendErrorResponse(w, *err, http.StatusBadRequest)
			return
		}

		// Try to create user in database
		err := dbCall(func() error {
			user.JoinedAt = time.Now()
			user.ID = 1 // This would be set by database
			return nil  // Simulated success
		})

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":      err.Error(),
				"user_name":  user.Name,
				"user_email": user.Email,
				"request_id": r.Header.Get("X-Request-ID"),
			}).Error("Failed to create user in database")

			h.sendErrorResponse(w, models.APIError{
				Type:      models.DatabaseError,
				Code:      "USER_CREATION_FAILED",
				Message:   "Unable to create user at this time",
				Details:   map[string]interface{}{"retryable": true},
				RequestID: r.Header.Get("X-Request-ID"),
				Timestamp: time.Now(),
				Retryable: true,
			}, http.StatusServiceUnavailable)
			return
		}

		// Cache the user locally
		userCache[user.ID] = &user

		response := models.APIResponse{
			Success: true,
			Data:    user,
			Metadata: map[string]interface{}{
				"created_at": time.Now(),
			},
		}

		h.sendJSONResponse(w, http.StatusCreated, response)
	}
}

// GetUser handles GET /users/{id} requests with cache fallback
func (h *UserHandler) GetUser(dbCall func(func() error) error, userCache map[int]*models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.sendErrorResponse(w, models.APIError{
				Type:      models.ValidationError,
				Code:      "INVALID_USER_ID",
				Message:   "User ID must be a valid number",
				Details:   map[string]interface{}{"provided_id": idStr},
				RequestID: r.Header.Get("X-Request-ID"),
				Timestamp: time.Now(),
				Retryable: false,
			}, http.StatusBadRequest)
			return
		}

		var user models.User

		// Try to get user from database
		err = dbCall(func() error {
			// Simulated database call
			if id == 1 {
				user = models.User{ID: 1, Name: "Alice", Email: "alice@example.com", JoinedAt: time.Now()}
				return nil
			}
			return sql.ErrNoRows
		})

		if err != nil {
			// Try cache as fallback
			if cachedUser, exists := userCache[id]; exists {
				response := models.APIResponse{
					Success:      true,
					Data:         *cachedUser,
					FallbackData: map[string]interface{}{"source": "cache"},
					Metadata:     map[string]interface{}{"cached": true},
				}
				h.sendJSONResponse(w, http.StatusOK, response)
				return
			}

			if err == sql.ErrNoRows {
				h.sendErrorResponse(w, models.APIError{
					Type:      models.ValidationError,
					Code:      "USER_NOT_FOUND",
					Message:   fmt.Sprintf("User with ID %d not found", id),
					RequestID: r.Header.Get("X-Request-ID"),
					Timestamp: time.Now(),
					Retryable: false,
				}, http.StatusNotFound)
			} else {
				h.sendErrorResponse(w, models.APIError{
					Type:      models.DatabaseError,
					Code:      "USER_FETCH_FAILED",
					Message:   "Unable to fetch user at this time",
					RequestID: r.Header.Get("X-Request-ID"),
					Timestamp: time.Now(),
					Retryable: true,
				}, http.StatusServiceUnavailable)
			}
			return
		}

		response := models.APIResponse{
			Success: true,
			Data:    user,
		}

		h.sendJSONResponse(w, http.StatusOK, response)
	}
}

func validateUser(user *models.User) *models.APIError {
	var errors []map[string]interface{}

	if user.Name == "" {
		errors = append(errors, map[string]interface{}{
			"field":   "name",
			"message": "Name is required",
		})
	}

	if user.Email == "" {
		errors = append(errors, map[string]interface{}{
			"field":   "email",
			"message": "Email is required",
		})
	} else if !isValidEmail(user.Email) {
		errors = append(errors, map[string]interface{}{
			"field":   "email",
			"message": "Email format is invalid",
			"value":   user.Email,
		})
	}

	if len(errors) > 0 {
		return &models.APIError{
			Type:      models.ValidationError,
			Code:      "VALIDATION_FAILED",
			Message:   "User validation failed",
			Details:   map[string]interface{}{"field_errors": errors},
			Timestamp: time.Now(),
			Retryable: false,
		}
	}

	return nil
}

func isValidEmail(email string) bool {
	// Simple email validation (production would use proper regex)
	return len(email) > 0 &&
		len(email) <= 254 &&
		contains(email, "@") &&
		contains(email, ".")
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
