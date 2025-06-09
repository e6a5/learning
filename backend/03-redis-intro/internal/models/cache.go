package models

import "fmt"

// KeyValue represents a Redis key-value pair
type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	TTL   int    `json:"ttl,omitempty"` // Time to live in seconds
}

// SetCacheRequest represents the request to set a cache value
type SetCacheRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	TTL   int    `json:"ttl,omitempty"`
}

// SetExpireRequest represents the request to set TTL for a key
type SetExpireRequest struct {
	TTL int `json:"ttl"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Validate validates the set cache request
func (r SetCacheRequest) Validate() error {
	if r.Key == "" {
		return &ValidationError{Field: "key", Message: "Key is required"}
	}
	if r.Value == "" {
		return &ValidationError{Field: "value", Message: "Value is required"}
	}
	if r.TTL < 0 {
		return &ValidationError{Field: "ttl", Message: "TTL must be non-negative"}
	}
	return nil
}

// Validate validates the set expire request
func (r SetExpireRequest) Validate() error {
	if r.TTL <= 0 {
		return &ValidationError{Field: "ttl", Message: "TTL must be positive"}
	}
	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// NewKeyValue creates a new KeyValue instance
func NewKeyValue(key, value string, ttl int) *KeyValue {
	return &KeyValue{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}
}
