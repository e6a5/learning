package middleware

import (
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/e6a5/learning/backend/08-monitoring/internal/models"
	"github.com/e6a5/learning/backend/08-monitoring/internal/repository"
)

// MonitoringMiddleware wraps HTTP handlers to collect metrics
type MonitoringMiddleware struct {
	repo *repository.MetricsRepository
}

// NewMonitoringMiddleware creates a new monitoring middleware
func NewMonitoringMiddleware(repo *repository.MetricsRepository) *MonitoringMiddleware {
	return &MonitoringMiddleware{repo: repo}
}

// responseWriter wraps http.ResponseWriter to capture response data
type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	responseSize int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.responseSize += int64(size)
	return size, err
}

// Wrap returns an HTTP handler that collects request metrics
func (m *MonitoringMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture metrics
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default status code
		}

		// Get request size
		requestSize := r.ContentLength
		if requestSize < 0 {
			requestSize = 0
		}

		// Process the request
		next.ServeHTTP(wrapped, r)

		// Calculate duration
		duration := time.Since(start)

		// Create request metrics
		metrics := models.RequestMetrics{
			Method:       r.Method,
			Path:         cleanPath(r.URL.Path),
			StatusCode:   wrapped.statusCode,
			Duration:     duration,
			RequestSize:  requestSize,
			ResponseSize: wrapped.responseSize,
			UserAgent:    r.UserAgent(),
			RemoteIP:     getRemoteIP(r),
			Timestamp:    time.Now(),
		}

		// Record metrics
		if err := m.repo.RecordRequest(metrics); err != nil {
			log.Printf("Error recording request metrics: %v", err)
		}

		// Log structured request information
		log.Printf("REQUEST: %s %s | Status: %d | Duration: %v | Size: %d bytes",
			metrics.Method, metrics.Path, metrics.StatusCode, metrics.Duration, metrics.ResponseSize)
	})
}

// cleanPath removes parameters from path for consistent metrics
func cleanPath(path string) string {
	// Remove query parameters
	if idx := strings.Index(path, "?"); idx != -1 {
		path = path[:idx]
	}

	// Remove trailing slash (except for root)
	if len(path) > 1 && strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	return path
}

// getRemoteIP extracts the real client IP address
func getRemoteIP(r *http.Request) string {
	// Check X-Forwarded-For header (proxy/load balancer)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Take the first IP if multiple are present
		if idx := strings.Index(forwarded, ","); idx != -1 {
			return strings.TrimSpace(forwarded[:idx])
		}
		return strings.TrimSpace(forwarded)
	}

	// Check X-Real-IP header (nginx proxy)
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	if strings.Contains(r.RemoteAddr, ":") {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			return host
		}
	}

	return r.RemoteAddr
}

// CorsMiddleware handles CORS headers
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs all requests in a structured format
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		log.Printf("ACCESS: %s %s %d %v %s",
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			time.Since(start),
			r.RemoteAddr,
		)
	})
}
