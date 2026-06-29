// SPDX-License-Identifier: MIT
package hardening

import "errors"

var (
	ErrCircuitOpen       = errors.New("circuit breaker is open")
	ErrMaxRetries        = errors.New("max retries exceeded")
	ErrIntegrityMismatch = errors.New("integrity check failed")
)
