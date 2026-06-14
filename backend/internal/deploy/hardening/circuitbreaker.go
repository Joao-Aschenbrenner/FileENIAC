package hardening

import (
	"fmt"
	"sync"
	"time"
)

// CircuitBreaker implements a simple circuit breaker pattern.
// States: CLOSED (normal), OPEN (failing), HALF_OPEN (probing).
type CircuitBreaker struct {
	mu              sync.Mutex
	failures        int
	lastFailure     time.Time
	threshold       int
	resetTimeout    time.Duration
	halfOpenTimeout time.Duration
	state           string
}

const (
	StateClosed   = "CLOSED"
	StateOpen     = "OPEN"
	StateHalfOpen = "HALF_OPEN"
)

type CircuitBreakerConfig struct {
	Threshold       int
	ResetTimeout    time.Duration
	HalfOpenTimeout time.Duration
}

func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		Threshold:       3,
		ResetTimeout:    30 * time.Second,
		HalfOpenTimeout: 5 * time.Second,
	}
}

func NewCircuitBreaker(cfg CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		threshold:       cfg.Threshold,
		resetTimeout:    cfg.ResetTimeout,
		halfOpenTimeout: cfg.HalfOpenTimeout,
		state:           StateClosed,
	}
}

func (cb *CircuitBreaker) State() string {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// Execute runs fn if the circuit is closed or half-open.
// Returns ErrCircuitOpen if the circuit is open and not yet ready to probe.
func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mu.Lock()

	if cb.state == StateOpen {
		if time.Since(cb.lastFailure) > cb.resetTimeout {
			cb.state = StateHalfOpen
		} else {
			cb.mu.Unlock()
			return fmt.Errorf("%w: circuit open for %s", ErrCircuitOpen, time.Since(cb.lastFailure).Round(time.Second))
		}
	}
	cb.mu.Unlock()

	err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		cb.lastFailure = time.Now()

		if cb.failures >= cb.threshold {
			cb.state = StateOpen
		} else if cb.state == StateHalfOpen {
			cb.state = StateOpen
		}
		return err
	}

	// Success resets
	cb.failures = 0
	if cb.state == StateHalfOpen {
		cb.state = StateClosed
	}
	return nil
}

func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures = 0
	cb.state = StateClosed
}
