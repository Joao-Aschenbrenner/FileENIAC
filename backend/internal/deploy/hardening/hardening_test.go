// SPDX-License-Identifier: MIT
package hardening

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCircuitBreaker_ClosedInitial(t *testing.T) {
	cb := NewCircuitBreaker(DefaultCircuitBreakerConfig())
	if cb.State() != StateClosed {
		t.Errorf("expected CLOSED, got %s", cb.State())
	}
}

func TestCircuitBreaker_OpensAfterThreshold(t *testing.T) {
	cfg := DefaultCircuitBreakerConfig()
	cfg.Threshold = 2
	cfg.ResetTimeout = 999 * 24 * time.Hour // effectively never resets
	cb := NewCircuitBreaker(cfg)

	errFail := errors.New("fail")
	err := cb.Execute(func() error { return errFail })
	if !errors.Is(err, errFail) {
		t.Fatal("expected original error")
	}

	err = cb.Execute(func() error { return errFail })
	if !errors.Is(err, errFail) {
		t.Fatal("expected original error")
	}

	err = cb.Execute(func() error { return nil })
	if !errors.Is(err, ErrCircuitOpen) {
		t.Fatalf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreaker_ResetsOnSuccess(t *testing.T) {
	cfg := DefaultCircuitBreakerConfig()
	cfg.Threshold = 2
	cfg.ResetTimeout = 999 * 24 * time.Hour
	cb := NewCircuitBreaker(cfg)

	errFail := errors.New("fail")
	cb.Execute(func() error { return errFail })
	cb.Execute(func() error { return errFail })
	// Circuit open
	cb.Execute(func() error { return errFail })
	if cb.State() != StateOpen {
		t.Fatalf("expected OPEN, got %s", cb.State())
	}

	cb.Reset()
	if cb.State() != StateClosed {
		t.Errorf("expected CLOSED after reset, got %s", cb.State())
	}

	err := cb.Execute(func() error { return nil })
	if err != nil {
		t.Fatalf("expected success after reset, got %v", err)
	}
}

func TestRetry_SuccessFirstAttempt(t *testing.T) {
	attempts := 0
	err := DoWithRetry(func() error {
		attempts++
		return nil
	}, DefaultRetryConfig())
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", attempts)
	}
}

func TestRetry_FailsAfterMaxAttempts(t *testing.T) {
	attempts := 0
	err := DoWithRetry(func() error {
		attempts++
		return errors.New("fail")
	}, RetryConfig{
		MaxAttempts:    3,
		InitialBackoff: 1,
		MaxBackoff:     10,
		BackoffFactor:  1.0,
	})
	if !errors.Is(err, ErrMaxRetries) {
		t.Fatalf("expected ErrMaxRetries, got %v", err)
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestRetry_SucceedsOnRetry(t *testing.T) {
	attempts := 0
	err := DoWithRetry(func() error {
		attempts++
		if attempts < 2 {
			return errors.New("transient")
		}
		return nil
	}, RetryConfig{
		MaxAttempts:    3,
		InitialBackoff: 1,
		MaxBackoff:     10,
		BackoffFactor:  1.0,
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}

func TestChecksumFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(path, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	sum, err := ChecksumFile(path)
	if err != nil {
		t.Fatalf("ChecksumFile failed: %v", err)
	}

	// SHA-256 of "hello"
	if !strings.HasPrefix(sum, "2cf24dba5fb0a30e26e83b2ac5b9e29e") {
		t.Errorf("unexpected checksum: %s", sum)
	}
}

func TestVerifyIntegrity(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(path, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	sum, _ := ChecksumFile(path)

	if err := VerifyIntegrity(path, sum); err != nil {
		t.Errorf("expected pass, got %v", err)
	}

	if err := VerifyIntegrity(path, "badhash"); err == nil {
		t.Error("expected error for bad hash")
	} else if !errors.Is(err, ErrIntegrityMismatch) {
		t.Errorf("expected ErrIntegrityMismatch, got %v", err)
	}
}
