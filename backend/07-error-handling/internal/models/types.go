package models

import "time"

// ErrorType represents different categories of errors
type ErrorType string

const (
	ValidationError     ErrorType = "validation_error"
	DatabaseError       ErrorType = "database_error"
	NetworkError        ErrorType = "network_error"
	AuthenticationError ErrorType = "authentication_error"
	RateLimitError      ErrorType = "rate_limit_error"
	InternalError       ErrorType = "internal_error"
	ServiceUnavailable  ErrorType = "service_unavailable"
)

// APIError represents a structured error response
type APIError struct {
	Type      ErrorType   `json:"type"`
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	RequestID string      `json:"request_id"`
	Timestamp time.Time   `json:"timestamp"`
	Retryable bool        `json:"retryable"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success      bool        `json:"success"`
	Data         interface{} `json:"data,omitempty"`
	Error        *APIError   `json:"error,omitempty"`
	FallbackData interface{} `json:"fallback_data,omitempty"`
	Metadata     interface{} `json:"metadata,omitempty"`
}

// User represents a user in the system
type User struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	JoinedAt time.Time `json:"joined_at"`
}

// RetryConfig holds configuration for retry logic
type RetryConfig struct {
	MaxAttempts   int
	BaseDelay     time.Duration
	MaxDelay      time.Duration
	BackoffFactor float64
	Jitter        bool
}
