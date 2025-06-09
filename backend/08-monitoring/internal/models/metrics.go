package models

import (
	"fmt"
	"time"
)

// HealthStatus represents the health state of a service
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
)

// HealthCheck represents a single health check result
type HealthCheck struct {
	Name      string                 `json:"name"`
	Status    HealthStatus           `json:"status"`
	Message   string                 `json:"message,omitempty"`
	Duration  time.Duration          `json:"duration_ms"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// HealthResponse represents the overall health response
type HealthResponse struct {
	Status      HealthStatus  `json:"status"`
	Version     string        `json:"version"`
	Uptime      time.Duration `json:"uptime_seconds"`
	Timestamp   time.Time     `json:"timestamp"`
	Checks      []HealthCheck `json:"checks"`
	Environment string        `json:"environment"`
}

// CustomMetric represents a custom application metric
type CustomMetric struct {
	Name      string            `json:"name"`
	Type      string            `json:"type"` // counter, gauge, histogram
	Value     float64           `json:"value"`
	Labels    map[string]string `json:"labels,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

// RequestMetrics represents HTTP request metrics
type RequestMetrics struct {
	Method       string        `json:"method"`
	Path         string        `json:"path"`
	StatusCode   int           `json:"status_code"`
	Duration     time.Duration `json:"duration_ms"`
	RequestSize  int64         `json:"request_size_bytes"`
	ResponseSize int64         `json:"response_size_bytes"`
	UserAgent    string        `json:"user_agent,omitempty"`
	RemoteIP     string        `json:"remote_ip,omitempty"`
	Timestamp    time.Time     `json:"timestamp"`
}

// SystemMetrics represents system-level metrics
type SystemMetrics struct {
	CPUUsage       float64   `json:"cpu_usage_percent"`
	MemoryUsage    int64     `json:"memory_usage_bytes"`
	MemoryTotal    int64     `json:"memory_total_bytes"`
	GoroutineCount int       `json:"goroutine_count"`
	HeapInUse      int64     `json:"heap_inuse_bytes"`
	HeapAlloc      int64     `json:"heap_alloc_bytes"`
	Timestamp      time.Time `json:"timestamp"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Validate validates a custom metric
func (m CustomMetric) Validate() error {
	if m.Name == "" {
		return &ValidationError{Field: "name", Message: "Metric name is required"}
	}
	if len(m.Name) > 100 {
		return &ValidationError{Field: "name", Message: "Metric name must be less than 100 characters"}
	}
	if m.Type == "" {
		return &ValidationError{Field: "type", Message: "Metric type is required"}
	}
	if m.Type != "counter" && m.Type != "gauge" && m.Type != "histogram" {
		return &ValidationError{Field: "type", Message: "Metric type must be counter, gauge, or histogram"}
	}
	return nil
}

// NewHealthCheck creates a new health check with validation
func NewHealthCheck(name, message string, status HealthStatus, duration time.Duration) (*HealthCheck, error) {
	if name == "" {
		return nil, &ValidationError{Field: "name", Message: "Health check name is required"}
	}
	if len(name) > 50 {
		return nil, &ValidationError{Field: "name", Message: "Health check name must be less than 50 characters"}
	}

	return &HealthCheck{
		Name:      name,
		Status:    status,
		Message:   message,
		Duration:  duration,
		Timestamp: time.Now(),
	}, nil
}

// IsHealthy returns true if the overall status is healthy
func (h HealthResponse) IsHealthy() bool {
	return h.Status == HealthStatusHealthy
}

// HasCriticalFailures returns true if any checks are unhealthy
func (h HealthResponse) HasCriticalFailures() bool {
	for _, check := range h.Checks {
		if check.Status == HealthStatusUnhealthy {
			return true
		}
	}
	return false
}
