package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

var rdb *redis.Client
var ctx = context.Background()

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	TTL   int    `json:"ttl,omitempty"` // Time to live in seconds
}

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func main() {
	// Connect to Redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password
		DB:       0,  // default DB
	})

	// Test connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	router := mux.NewRouter()

	// Key-Value operations
	router.HandleFunc("/cache/{key}", getValue).Methods("GET")
	router.HandleFunc("/cache", setValue).Methods("POST")
	router.HandleFunc("/cache/{key}", deleteValue).Methods("DELETE")
	router.HandleFunc("/cache", getAllKeys).Methods("GET")

	// Cache operations
	router.HandleFunc("/cache/{key}/ttl", getTTL).Methods("GET")
	router.HandleFunc("/cache/{key}/expire", setExpire).Methods("POST")

	// Health check
	router.HandleFunc("/health", healthCheck).Methods("GET")

	log.Println("🚀 Redis Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getValue(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		respondJSON(w, http.StatusNotFound, Response{Error: "Key not found"})
		return
	} else if err != nil {
		respondJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Data: KeyValue{Key: key, Value: val},
	})
}

func setValue(w http.ResponseWriter, r *http.Request) {
	var kv KeyValue
	if err := json.NewDecoder(r.Body).Decode(&kv); err != nil {
		respondJSON(w, http.StatusBadRequest, Response{Error: "Invalid JSON"})
		return
	}

	if kv.Key == "" || kv.Value == "" {
		respondJSON(w, http.StatusBadRequest, Response{Error: "Key and value are required"})
		return
	}

	var err error
	if kv.TTL > 0 {
		err = rdb.Set(ctx, kv.Key, kv.Value, time.Duration(kv.TTL)*time.Second).Err()
	} else {
		err = rdb.Set(ctx, kv.Key, kv.Value, 0).Err()
	}

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	respondJSON(w, http.StatusCreated, Response{
		Message: "Key set successfully",
		Data:    kv,
	})
}

func deleteValue(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	deleted, err := rdb.Del(ctx, key).Result()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	if deleted == 0 {
		respondJSON(w, http.StatusNotFound, Response{Error: "Key not found"})
		return
	}

	respondJSON(w, http.StatusOK, Response{Message: "Key deleted successfully"})
}

func getAllKeys(w http.ResponseWriter, r *http.Request) {
	keys, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Data: map[string]interface{}{"keys": keys, "count": len(keys)},
	})
}

func getTTL(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	ttl, err := rdb.TTL(ctx, key).Result()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Data: map[string]interface{}{
			"key": key,
			"ttl": ttl.Seconds(),
		},
	})
}

func setExpire(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	var req struct {
		TTL int `json:"ttl"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, Response{Error: "Invalid JSON"})
		return
	}

	if req.TTL <= 0 {
		respondJSON(w, http.StatusBadRequest, Response{Error: "TTL must be positive"})
		return
	}

	success, err := rdb.Expire(ctx, key, time.Duration(req.TTL)*time.Second).Result()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	if !success {
		respondJSON(w, http.StatusNotFound, Response{Error: "Key not found"})
		return
	}

	respondJSON(w, http.StatusOK, Response{Message: "Expiration set successfully"})
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		respondJSON(w, http.StatusServiceUnavailable, Response{Error: "Redis unavailable"})
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Message: "Service healthy",
		Data:    map[string]string{"redis": "connected"},
	})
}

func respondJSON(w http.ResponseWriter, statusCode int, data Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
