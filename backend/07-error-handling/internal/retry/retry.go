package retry

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/e6a5/learning/backend/07-error-handling/internal/models"
	"github.com/sirupsen/logrus"
)

// WithRetry executes the given function with retry logic
func WithRetry(operation string, config models.RetryConfig, fn func() error) error {
	var lastErr error

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		lastErr = fn()
		if lastErr == nil {
			if attempt > 1 {
				logrus.WithFields(logrus.Fields{
					"operation": operation,
					"attempt":   attempt,
				}).Info("Operation succeeded after retry")
			}
			return nil
		}

		if attempt == config.MaxAttempts {
			break
		}

		delay := calculateBackoffDelay(config, attempt)
		logrus.WithFields(logrus.Fields{
			"operation": operation,
			"attempt":   attempt,
			"error":     lastErr.Error(),
			"delay":     delay,
		}).Warn("Operation failed, retrying")

		time.Sleep(delay)
	}

	return fmt.Errorf("operation %s failed after %d attempts: %w", operation, config.MaxAttempts, lastErr)
}

func calculateBackoffDelay(config models.RetryConfig, attempt int) time.Duration {
	delay := float64(config.BaseDelay) * math.Pow(config.BackoffFactor, float64(attempt-1))

	if delay > float64(config.MaxDelay) {
		delay = float64(config.MaxDelay)
	}

	if config.Jitter {
		jitterRange := delay * 0.1
		jitter := (rand.Float64() - 0.5) * 2 * jitterRange
		delay += jitter
	}

	return time.Duration(delay)
}
