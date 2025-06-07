package circuit

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// State represents the circuit breaker state
type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

// Breaker implements the circuit breaker pattern
type Breaker struct {
	name         string
	maxFailures  int
	resetTimeout time.Duration
	state        State
	failures     int
	lastFailTime time.Time
	successCount int
	mutex        sync.RWMutex
}

// New creates a new circuit breaker
func New(name string, maxFailures int, resetTimeout time.Duration) *Breaker {
	return &Breaker{
		name:         name,
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        Closed,
	}
}

// Call executes the given function with circuit breaker protection
func (cb *Breaker) Call(fn func() error) error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	switch cb.state {
	case Open:
		if time.Since(cb.lastFailTime) > cb.resetTimeout {
			cb.state = HalfOpen
			cb.successCount = 0
			logrus.WithField("circuit", cb.name).Info("Circuit breaker moved to half-open state")
		} else {
			return fmt.Errorf("circuit breaker is open for %s", cb.name)
		}
	}

	err := fn()

	if err != nil {
		cb.failures++
		cb.lastFailTime = time.Now()

		if cb.state == HalfOpen || cb.failures >= cb.maxFailures {
			cb.state = Open
			logrus.WithFields(logrus.Fields{
				"circuit":  cb.name,
				"failures": cb.failures,
			}).Warn("Circuit breaker opened")
		}
		return err
	}

	// Success
	if cb.state == HalfOpen {
		cb.successCount++
		if cb.successCount >= 3 { // Require 3 successes to close
			cb.state = Closed
			cb.failures = 0
			logrus.WithField("circuit", cb.name).Info("Circuit breaker closed")
		}
	} else {
		cb.failures = 0
	}

	return nil
}

// GetState returns the current state of the circuit breaker
func (cb *Breaker) GetState() string {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	switch cb.state {
	case Closed:
		return "closed"
	case Open:
		return "open"
	case HalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// GetFailures returns the current failure count
func (cb *Breaker) GetFailures() int {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.failures
}

// GetLastFailTime returns the last failure time
func (cb *Breaker) GetLastFailTime() time.Time {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.lastFailTime
}

// GetSuccessCount returns the current success count in half-open state
func (cb *Breaker) GetSuccessCount() int {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.successCount
}

// Reset resets the circuit breaker to closed state
func (cb *Breaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.state = Closed
	cb.failures = 0
	cb.successCount = 0
}
