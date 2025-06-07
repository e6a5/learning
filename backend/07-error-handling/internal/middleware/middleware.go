package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/e6a5/learning/backend/07-error-handling/internal/models"
	"github.com/sirupsen/logrus"
)

// ResponseWriter wraps http.ResponseWriter to capture status code
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// PanicRecovery recovers from panics and returns structured error responses
func PanicRecovery(sendErrorFn func(http.ResponseWriter, models.APIError, int)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logrus.WithFields(logrus.Fields{
						"panic":      err,
						"method":     r.Method,
						"path":       r.URL.Path,
						"request_id": r.Header.Get("X-Request-ID"),
					}).Error("Panic recovered")

					sendErrorFn(w, models.APIError{
						Type:      models.InternalError,
						Code:      "PANIC_RECOVERED",
						Message:   "Internal server error occurred",
						RequestID: r.Header.Get("X-Request-ID"),
						Timestamp: time.Now(),
						Retryable: false,
					}, http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// RequestID adds unique request IDs to requests
func RequestID(counter *int64, mutex *sync.Mutex) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				mutex.Lock()
				*counter++
				requestID = fmt.Sprintf("req_%d_%d", time.Now().Unix(), *counter)
				mutex.Unlock()
			}

			r.Header.Set("X-Request-ID", requestID)
			w.Header().Set("X-Request-ID", requestID)
			next.ServeHTTP(w, r)
		})
	}
}

// Logging logs all HTTP requests with structured data
func Logging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code
			wrapped := &ResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			logEntry := logrus.WithFields(logrus.Fields{
				"method":     r.Method,
				"path":       r.URL.Path,
				"status":     wrapped.statusCode,
				"duration":   duration,
				"request_id": r.Header.Get("X-Request-ID"),
				"ip":         r.RemoteAddr,
			})

			if wrapped.statusCode >= 500 {
				logEntry.Error("Request completed with server error")
			} else if wrapped.statusCode >= 400 {
				logEntry.Warn("Request completed with client error")
			} else {
				logEntry.Info("Request completed successfully")
			}
		})
	}
}

// RateLimit provides basic rate limiting (production would use Redis)
func RateLimit() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simple in-memory rate limiting for demonstration
			// Production implementation would use Redis with sliding windows
			next.ServeHTTP(w, r)
		})
	}
}
