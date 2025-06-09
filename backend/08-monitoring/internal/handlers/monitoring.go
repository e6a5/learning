package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/e6a5/learning/backend/08-monitoring/internal/models"
	"github.com/e6a5/learning/backend/08-monitoring/internal/repository"
	"github.com/e6a5/learning/backend/08-monitoring/internal/utils"
)

// MonitoringHandler handles monitoring-related HTTP requests
type MonitoringHandler struct {
	repo           *repository.MetricsRepository
	healthCheckers []repository.HealthChecker
	promRegistry   *prometheus.Registry
}

// NewMonitoringHandler creates a new monitoring handler
func NewMonitoringHandler(repo *repository.MetricsRepository, checkers []repository.HealthChecker) *MonitoringHandler {
	return &MonitoringHandler{
		repo:           repo,
		healthCheckers: checkers,
		promRegistry:   prometheus.NewRegistry(),
	}
}

// HealthCheck handles GET /health - comprehensive health check
func (h *MonitoringHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	response := h.repo.PerformHealthChecks(ctx, h.healthCheckers)

	statusCode := http.StatusOK
	if response.HasCriticalFailures() {
		statusCode = http.StatusServiceUnavailable
	} else if !response.IsHealthy() {
		statusCode = http.StatusOK // 200 for degraded but still serving
	}

	utils.RespondJSON(w, statusCode, response)
}

// LivenessCheck handles GET /health/live - simple liveness probe
func (h *MonitoringHandler) LivenessCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "alive",
		"timestamp": time.Now(),
		"uptime":    time.Since(h.repo.GetSystemMetrics().Timestamp),
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

// ReadinessCheck handles GET /health/ready - readiness probe
func (h *MonitoringHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	response := h.repo.PerformHealthChecks(ctx, h.healthCheckers)

	statusCode := http.StatusOK
	if response.HasCriticalFailures() {
		statusCode = http.StatusServiceUnavailable
	}

	readinessResponse := map[string]interface{}{
		"ready":     response.IsHealthy(),
		"status":    response.Status,
		"timestamp": time.Now(),
		"checks":    len(response.Checks),
	}

	utils.RespondJSON(w, statusCode, readinessResponse)
}

// GetMetrics handles GET /metrics - Prometheus-style metrics
func (h *MonitoringHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	promhttp.HandlerFor(h.promRegistry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}

// GetCustomMetrics handles GET /api/metrics - custom JSON metrics
func (h *MonitoringHandler) GetCustomMetrics(w http.ResponseWriter, r *http.Request) {
	requestMetrics := h.repo.GetRequestMetrics()
	errorMetrics := h.repo.GetErrorMetrics()
	customMetrics := h.repo.GetCustomMetrics()
	systemMetrics := h.repo.GetSystemMetrics()

	response := map[string]interface{}{
		"request_metrics": requestMetrics,
		"error_metrics":   errorMetrics,
		"custom_metrics":  customMetrics,
		"system_metrics":  systemMetrics,
		"timestamp":       time.Now(),
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

// PostCustomMetric handles POST /api/metrics - submit custom metric
func (h *MonitoringHandler) PostCustomMetric(w http.ResponseWriter, r *http.Request) {
	var metric models.CustomMetric

	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON format",
		})
		return
	}

	metric.Timestamp = time.Now()

	if err := h.repo.RecordCustomMetric(metric); err != nil {
		log.Printf("Error recording custom metric: %v", err)
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Metric recorded successfully",
		"metric":  metric,
	})
}

// GetSystemInfo handles GET /api/system - system information
func (h *MonitoringHandler) GetSystemInfo(w http.ResponseWriter, r *http.Request) {
	systemMetrics := h.repo.GetSystemMetrics()

	response := map[string]interface{}{
		"system_metrics": systemMetrics,
		"process_info": map[string]interface{}{
			"goroutines": systemMetrics.GoroutineCount,
			"memory": map[string]interface{}{
				"heap_alloc":   fmt.Sprintf("%.2f MB", float64(systemMetrics.HeapAlloc)/1024/1024),
				"heap_inuse":   fmt.Sprintf("%.2f MB", float64(systemMetrics.HeapInUse)/1024/1024),
				"memory_usage": fmt.Sprintf("%.2f MB", float64(systemMetrics.MemoryUsage)/1024/1024),
				"memory_total": fmt.Sprintf("%.2f MB", float64(systemMetrics.MemoryTotal)/1024/1024),
			},
		},
		"timestamp": time.Now(),
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

// DemoEndpoint handles GET /api/demo - endpoint to generate metrics
func (h *MonitoringHandler) DemoEndpoint(w http.ResponseWriter, r *http.Request) {
	// Get optional error parameter to simulate errors
	errorParam := r.URL.Query().Get("error")

	// Get optional delay parameter to simulate slow responses
	delayParam := r.URL.Query().Get("delay")

	if delayParam != "" {
		if delay, err := strconv.Atoi(delayParam); err == nil && delay > 0 && delay < 5000 {
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}

	// Simulate different types of responses
	switch errorParam {
	case "400":
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Simulated bad request error",
		})
		return
	case "500":
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Simulated internal server error",
		})
		return
	case "timeout":
		time.Sleep(6 * time.Second) // Longer than typical timeout
		utils.RespondJSON(w, http.StatusRequestTimeout, map[string]string{
			"error": "Simulated timeout",
		})
		return
	}

	// Record custom demo metric
	metric := models.CustomMetric{
		Name:  "demo_requests_total",
		Type:  "counter",
		Value: 1,
		Labels: map[string]string{
			"endpoint": "demo",
			"method":   r.Method,
		},
		Timestamp: time.Now(),
	}

	if err := h.repo.RecordCustomMetric(metric); err != nil {
		log.Printf("Error recording demo metric: %v", err)
	}

	response := map[string]interface{}{
		"message":   "Demo endpoint successful",
		"timestamp": time.Now(),
		"method":    r.Method,
		"path":      r.URL.Path,
		"tips": []string{
			"Try adding ?error=400 to simulate bad request",
			"Try adding ?error=500 to simulate server error",
			"Try adding ?delay=1000 to simulate slow response",
			"Check /metrics for Prometheus metrics",
			"Check /api/metrics for JSON metrics",
		},
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

// GetStatus handles GET /api/status - application status overview
func (h *MonitoringHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	healthResponse := h.repo.PerformHealthChecks(ctx, h.healthCheckers)
	systemMetrics := h.repo.GetSystemMetrics()
	requestMetrics := h.repo.GetRequestMetrics()

	// Calculate total requests
	totalRequests := int64(0)
	for _, count := range requestMetrics {
		totalRequests += count
	}

	response := map[string]interface{}{
		"application": map[string]interface{}{
			"status":      healthResponse.Status,
			"version":     healthResponse.Version,
			"environment": healthResponse.Environment,
			"uptime":      healthResponse.Uptime.Seconds(),
		},
		"traffic": map[string]interface{}{
			"total_requests": totalRequests,
			"endpoints":      len(requestMetrics),
		},
		"performance": map[string]interface{}{
			"goroutines": systemMetrics.GoroutineCount,
			"memory_mb":  float64(systemMetrics.MemoryUsage) / 1024 / 1024,
			"heap_mb":    float64(systemMetrics.HeapAlloc) / 1024 / 1024,
		},
		"health_checks": map[string]interface{}{
			"total":    len(healthResponse.Checks),
			"healthy":  countHealthyChecks(healthResponse.Checks),
			"degraded": countDegradedChecks(healthResponse.Checks),
			"failed":   countFailedChecks(healthResponse.Checks),
		},
		"timestamp": time.Now(),
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

// Helper functions for health check counting
func countHealthyChecks(checks []models.HealthCheck) int {
	count := 0
	for _, check := range checks {
		if check.Status == models.HealthStatusHealthy {
			count++
		}
	}
	return count
}

func countDegradedChecks(checks []models.HealthCheck) int {
	count := 0
	for _, check := range checks {
		if check.Status == models.HealthStatusDegraded {
			count++
		}
	}
	return count
}

func countFailedChecks(checks []models.HealthCheck) int {
	count := 0
	for _, check := range checks {
		if check.Status == models.HealthStatusUnhealthy {
			count++
		}
	}
	return count
}
