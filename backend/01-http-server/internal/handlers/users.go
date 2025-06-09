package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/e6a5/learning/backend/01-http-server/internal/models"
	"github.com/e6a5/learning/backend/01-http-server/internal/repository"
	"github.com/e6a5/learning/backend/01-http-server/internal/utils"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	repo *repository.UserRepository
}

// NewUserHandler creates a new user handler
func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// GetUsers handles GET /users - returns all users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := h.repo.GetAll()

	response := models.Response{
		Success: true,
		Message: "Found " + strconv.Itoa(len(users)) + " users",
		Data:    users,
	}

	utils.SendJSONResponse(w, http.StatusOK, response)
}

// CreateUser handles POST /users - creates a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest

	// Parse JSON from request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := models.Response{
			Success: false,
			Message: "Invalid JSON format",
		}
		utils.SendJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		response := models.Response{
			Success: false,
			Message: err.Error(),
		}
		utils.SendJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Create user
	user := h.repo.Create(req.Name, req.Email)

	logrus.WithFields(logrus.Fields{
		"user_id": user.ID,
		"name":    user.Name,
		"email":   user.Email,
	}).Info("New user created")

	response := models.Response{
		Success: true,
		Message: "User created successfully",
		Data:    user,
	}

	utils.SendJSONResponse(w, http.StatusCreated, response)
}

// GetUser handles GET /users/{id} - returns a specific user
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert string to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := models.Response{
			Success: false,
			Message: "Invalid user ID",
		}
		utils.SendJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Get user from repository
	user, err := h.repo.GetByID(id)
	if err != nil {
		response := models.Response{
			Success: false,
			Message: "User not found",
		}
		utils.SendJSONResponse(w, http.StatusNotFound, response)
		return
	}

	response := models.Response{
		Success: true,
		Message: "User found",
		Data:    user,
	}

	utils.SendJSONResponse(w, http.StatusOK, response)
}
