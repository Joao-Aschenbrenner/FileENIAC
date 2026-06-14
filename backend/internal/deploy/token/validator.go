package token

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Validator struct {
	secret string
}

func NewValidator(secret string) *Validator {
	return &Validator{
		secret: secret,
	}
}

func (v *Validator) ValidateRequest(r *http.Request) error {
	sig := r.Header.Get("X-Deploy-Token")
	tsStr := r.Header.Get("X-Deploy-Timestamp")
	project := r.Header.Get("X-Deploy-Project")

	if sig == "" || tsStr == "" || project == "" {
		return fmt.Errorf("missing required headers")
	}

	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid timestamp: %w", err)
	}
	timestamp := time.Unix(ts, 0)

	signer := NewSigner(v.secret)
	if !signer.Validate(project, sig, timestamp) {
		return fmt.Errorf("invalid token")
	}

	return nil
}