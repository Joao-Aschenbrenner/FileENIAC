package hardening

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"os"
	"time"
)

// RetryConfig controls the exponential backoff retry behavior.
type RetryConfig struct {
	MaxAttempts     int
	InitialBackoff  time.Duration
	MaxBackoff      time.Duration
	BackoffFactor   float64
}

func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:    3,
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     30 * time.Second,
		BackoffFactor:  2.0,
	}
}

// DoWithRetry executes fn up to MaxAttempts times with exponential backoff.
// Returns ErrMaxRetries if all attempts fail.
func DoWithRetry(fn func() error, cfg RetryConfig) error {
	var lastErr error

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		if err := fn(); err != nil {
			lastErr = err
			if attempt < cfg.MaxAttempts {
				backoff := time.Duration(math.Min(
					float64(cfg.InitialBackoff)*math.Pow(cfg.BackoffFactor, float64(attempt-1)),
					float64(cfg.MaxBackoff),
				))
				time.Sleep(backoff)
			}
			continue
		}
		return nil
	}

	return fmt.Errorf("%w after %d attempts: %w", ErrMaxRetries, cfg.MaxAttempts, lastErr)
}

// ChecksumFile computes SHA-256 of a file.
func ChecksumFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file for checksum: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("failed to read file for checksum: %w", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// VerifyIntegrity checks that a file has the expected SHA-256 hash.
func VerifyIntegrity(path, expectedHash string) error {
	actual, err := ChecksumFile(path)
	if err != nil {
		return err
	}
	if actual != expectedHash {
		return fmt.Errorf("%w: expected %s, got %s", ErrIntegrityMismatch, expectedHash, actual)
	}
	return nil
}
