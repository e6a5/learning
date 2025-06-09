package repository

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/e6a5/learning/backend/08-monitoring/internal/models"
)

// MetricsRepository handles metrics storage and retrieval
type MetricsRepository struct {
	mu            sync.RWMutex
	requestCount  map[string]int64
	errorCount    map[string]int64
	customMetrics map[string]models.CustomMetric
	startTime     time.Time
	version       string
	environment   string
}

// HealthChecker defines interface for health checks
type HealthChecker interface {
	Check(ctx context.Context) models.HealthCheck
}

// DatabaseHealthChecker checks database connectivity
type DatabaseHealthChecker struct {
	name string
	url  string
}

// ExternalServiceHealthChecker checks external service health
type ExternalServiceHealthChecker struct {
	name string
	url  string
}

// NewMetricsRepository creates a new metrics repository
func NewMetricsRepository(version, environment string) *MetricsRepository {
	return &MetricsRepository{
		requestCount:  make(map[string]int64),
		errorCount:    make(map[string]int64),
		customMetrics: make(map[string]models.CustomMetric),
		startTime:     time.Now(),
		version:       version,
		environment:   environment,
	}
}

// RecordRequest records HTTP request metrics
func (r *MetricsRepository) RecordRequest(metrics models.RequestMetrics) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("%s:%s", metrics.Method, metrics.Path)
	r.requestCount[key]++

	if metrics.StatusCode >= 400 {
		errorKey := fmt.Sprintf("%s:%d", key, metrics.StatusCode)
		r.errorCount[errorKey]++
	}

	return nil
}

// RecordCustomMetric stores a custom metric
func (r *MetricsRepository) RecordCustomMetric(metric models.CustomMetric) error {
	if err := metric.Validate(); err != nil {
		return fmt.Errorf("invalid metric: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	key := r.buildMetricKey(metric.Name, metric.Labels)
	r.customMetrics[key] = metric

	return nil
}

// GetRequestMetrics returns request count metrics
func (r *MetricsRepository) GetRequestMetrics() map[string]int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]int64)
	for k, v := range r.requestCount {
		result[k] = v
	}
	return result
}

// GetErrorMetrics returns error count metrics
func (r *MetricsRepository) GetErrorMetrics() map[string]int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]int64)
	for k, v := range r.errorCount {
		result[k] = v
	}
	return result
}

// GetCustomMetrics returns all custom metrics
func (r *MetricsRepository) GetCustomMetrics() []models.CustomMetric {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []models.CustomMetric
	for _, metric := range r.customMetrics {
		result = append(result, metric)
	}
	return result
}

// GetSystemMetrics returns current system metrics
func (r *MetricsRepository) GetSystemMetrics() models.SystemMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return models.SystemMetrics{
		CPUUsage:       0, // Simplified - would need CPU monitoring
		MemoryUsage:    int64(m.Alloc),
		MemoryTotal:    int64(m.Sys),
		GoroutineCount: runtime.NumGoroutine(),
		HeapInUse:      int64(m.HeapInuse),
		HeapAlloc:      int64(m.HeapAlloc),
		Timestamp:      time.Now(),
	}
}

// PerformHealthChecks runs all health checks and returns results
func (r *MetricsRepository) PerformHealthChecks(ctx context.Context, checkers []HealthChecker) models.HealthResponse {
	var checks []models.HealthCheck
	overallStatus := models.HealthStatusHealthy

	// Run all health checks
	for _, checker := range checkers {
		check := checker.Check(ctx)
		checks = append(checks, check)

		// Determine overall status
		if check.Status == models.HealthStatusUnhealthy {
			overallStatus = models.HealthStatusUnhealthy
		} else if check.Status == models.HealthStatusDegraded && overallStatus == models.HealthStatusHealthy {
			overallStatus = models.HealthStatusDegraded
		}
	}

	return models.HealthResponse{
		Status:      overallStatus,
		Version:     r.version,
		Uptime:      time.Since(r.startTime),
		Timestamp:   time.Now(),
		Checks:      checks,
		Environment: r.environment,
	}
}

// buildMetricKey creates a unique key for metrics with labels
func (r *MetricsRepository) buildMetricKey(name string, labels map[string]string) string {
	key := name
	for k, v := range labels {
		key += fmt.Sprintf(",%s=%s", k, v)
	}
	return key
}

// NewDatabaseHealthChecker creates a database health checker
func NewDatabaseHealthChecker(name, url string) *DatabaseHealthChecker {
	return &DatabaseHealthChecker{name: name, url: url}
}

// Check performs database health check
func (d *DatabaseHealthChecker) Check(ctx context.Context) models.HealthCheck {
	start := time.Now()

	// Simplified database check - in reality, you'd ping the actual database
	status := models.HealthStatusHealthy
	message := "Database connection successful"

	// Simulate database check with timeout
	select {
	case <-time.After(100 * time.Millisecond): // Simulate DB response time
		// Check passed
	case <-ctx.Done():
		status = models.HealthStatusUnhealthy
		message = "Database check timeout"
	}

	duration := time.Since(start)

	check, _ := models.NewHealthCheck(d.name, message, status, duration)
	check.Details = map[string]interface{}{
		"connection_url": d.url,
		"type":           "database",
	}

	return *check
}

// NewExternalServiceHealthChecker creates an external service health checker
func NewExternalServiceHealthChecker(name, url string) *ExternalServiceHealthChecker {
	return &ExternalServiceHealthChecker{name: name, url: url}
}

// Check performs external service health check
func (e *ExternalServiceHealthChecker) Check(ctx context.Context) models.HealthCheck {
	start := time.Now()

	status := models.HealthStatusHealthy
	message := "External service responding"

	// Create HTTP client with timeout
	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequestWithContext(ctx, "GET", e.url, nil)
	if err != nil {
		check, _ := models.NewHealthCheck(e.name, fmt.Sprintf("Failed to create request: %v", err),
			models.HealthStatusUnhealthy, time.Since(start))
		return *check
	}

	resp, err := client.Do(req)
	if err != nil {
		status = models.HealthStatusUnhealthy
		message = fmt.Sprintf("Request failed: %v", err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			status = models.HealthStatusDegraded
			message = fmt.Sprintf("Service returned status %d", resp.StatusCode)
		}
	}

	duration := time.Since(start)
	check, _ := models.NewHealthCheck(e.name, message, status, duration)
	check.Details = map[string]interface{}{
		"service_url": e.url,
		"type":        "external_service",
	}

	return *check
}
