package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/e6a5/learning/backend/03-redis-intro/internal/models"
	"github.com/e6a5/learning/backend/03-redis-intro/internal/repository"
	"github.com/e6a5/learning/backend/03-redis-intro/internal/utils"
)

// CacheHandler handles cache-related HTTP requests
type CacheHandler struct {
	repo *repository.CacheRepository
}

// NewCacheHandler creates a new cache handler
func NewCacheHandler(repo *repository.CacheRepository) *CacheHandler {
	return &CacheHandler{repo: repo}
}

// GetValue handles GET /cache/{key} - retrieves a cached value
func (h *CacheHandler) GetValue(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	kv, err := h.repo.Get(key)
	if err != nil {
		log.Printf("Error getting key %s: %v", key, err)
		if err.Error() == "key not found: "+key {
			utils.RespondJSON(w, http.StatusNotFound, models.APIResponse{Error: "Key not found"})
		} else {
			utils.RespondJSON(w, http.StatusInternalServerError, models.APIResponse{Error: "Internal server error"})
		}
		return
	}

	utils.RespondJSON(w, http.StatusOK, models.APIResponse{Data: kv})
}

// SetValue handles POST /cache - sets a cached value
func (h *CacheHandler) SetValue(w http.ResponseWriter, r *http.Request) {
	var req models.SetCacheRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, models.APIResponse{Error: "Invalid JSON"})
		return
	}

	if err := req.Validate(); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, models.APIResponse{Error: err.Error()})
		return
	}

	if err := h.repo.Set(req.Key, req.Value, req.TTL); err != nil {
		log.Printf("Error setting key %s: %v", req.Key, err)
		utils.RespondJSON(w, http.StatusInternalServerError, models.APIResponse{Error: "Internal server error"})
		return
	}

	kv := models.NewKeyValue(req.Key, req.Value, req.TTL)
	utils.RespondJSON(w, http.StatusCreated, models.APIResponse{
		Message: "Key set successfully",
		Data:    kv,
	})
}

// DeleteValue handles DELETE /cache/{key} - deletes a cached value
func (h *CacheHandler) DeleteValue(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	if err := h.repo.Delete(key); err != nil {
		log.Printf("Error deleting key %s: %v", key, err)
		if err.Error() == "key not found: "+key {
			utils.RespondJSON(w, http.StatusNotFound, models.APIResponse{Error: "Key not found"})
		} else {
			utils.RespondJSON(w, http.StatusInternalServerError, models.APIResponse{Error: "Internal server error"})
		}
		return
	}

	utils.RespondJSON(w, http.StatusOK, models.APIResponse{Message: "Key deleted successfully"})
}

// GetAllKeys handles GET /cache - retrieves all keys
func (h *CacheHandler) GetAllKeys(w http.ResponseWriter, r *http.Request) {
	pattern := r.URL.Query().Get("pattern")

	keys, err := h.repo.GetAllKeys(pattern)
	if err != nil {
		log.Printf("Error getting all keys: %v", err)
		utils.RespondJSON(w, http.StatusInternalServerError, models.APIResponse{Error: "Internal server error"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, models.APIResponse{
		Data: map[string]interface{}{
			"keys":  keys,
			"count": len(keys),
		},
	})
}

// GetTTL handles GET /cache/{key}/ttl - gets TTL for a key
func (h *CacheHandler) GetTTL(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	ttl, err := h.repo.GetTTL(key)
	if err != nil {
		log.Printf("Error getting TTL for key %s: %v", key, err)
		utils.RespondJSON(w, http.StatusInternalServerError, models.APIResponse{Error: "Internal server error"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, models.APIResponse{
		Data: map[string]interface{}{
			"key": key,
			"ttl": ttl.Seconds(),
		},
	})
}

// SetExpire handles POST /cache/{key}/expire - sets TTL for a key
func (h *CacheHandler) SetExpire(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	var req models.SetExpireRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, models.APIResponse{Error: "Invalid JSON"})
		return
	}

	if err := req.Validate(); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, models.APIResponse{Error: err.Error()})
		return
	}

	if err := h.repo.SetExpire(key, req.TTL); err != nil {
		log.Printf("Error setting expire for key %s: %v", key, err)
		if err.Error() == "key not found: "+key {
			utils.RespondJSON(w, http.StatusNotFound, models.APIResponse{Error: "Key not found"})
		} else {
			utils.RespondJSON(w, http.StatusInternalServerError, models.APIResponse{Error: "Internal server error"})
		}
		return
	}

	utils.RespondJSON(w, http.StatusOK, models.APIResponse{Message: "Expiration set successfully"})
}

// HealthCheck handles GET /health - checks Redis connectivity
func (h *CacheHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.Ping(); err != nil {
		log.Printf("Health check failed: %v", err)
		utils.RespondJSON(w, http.StatusServiceUnavailable, models.APIResponse{Error: "Redis unavailable"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, models.APIResponse{
		Message: "Service healthy",
		Data:    map[string]string{"redis": "connected"},
	})
}
