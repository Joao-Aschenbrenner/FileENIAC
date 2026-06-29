// SPDX-License-Identifier: MIT
package token

import (
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestSigner_Sign(t *testing.T) {
	signer := NewSigner("test-secret")

	sig, ts, err := signer.Sign("test-project")
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	if sig == "" {
		t.Error("signature should not be empty")
	}

	if ts.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestSigner_Validate(t *testing.T) {
	secret := "test-secret"
	signer := NewSigner(secret)

	sig, ts, _ := signer.Sign("test-project")

	if !signer.Validate("test-project", sig, ts) {
		t.Error("valid signature should pass validation")
	}

	if signer.Validate("test-project", sig, ts.Add(-10*time.Minute)) {
		t.Error("expired token should fail validation")
	}

	if signer.Validate("wrong-project", sig, ts) {
		t.Error("wrong project should fail validation")
	}

	if signer.Validate("test-project", "wrong-signature", ts) {
		t.Error("wrong signature should fail validation")
	}
}

func TestSigner_GenerateHeaders(t *testing.T) {
	signer := NewSigner("test-secret")

	headers, err := signer.GenerateHeaders("test-project")
	if err != nil {
		t.Fatalf("GenerateHeaders failed: %v", err)
	}

	if headers["X-Deploy-Token"] == "" {
		t.Error("X-Deploy-Token should not be empty")
	}

	if headers["X-Deploy-Timestamp"] == "" {
		t.Error("X-Deploy-Timestamp should not be empty")
	}

	if headers["X-Deploy-Project"] != "test-project" {
		t.Errorf("expected project 'test-project', got '%s'", headers["X-Deploy-Project"])
	}
}

func TestValidator_ValidateRequest(t *testing.T) {
	secret := "test-secret"
	validator := NewValidator(secret)
	signer := NewSigner(secret)

	sig, ts, _ := signer.Sign("test-project")

	req, _ := http.NewRequest("POST", "/deploy", nil)
	req.Header.Set("X-Deploy-Token", sig)
	req.Header.Set("X-Deploy-Timestamp", strconv.FormatInt(ts.Unix(), 10))
	req.Header.Set("X-Deploy-Project", "test-project")

	if err := validator.ValidateRequest(req); err != nil {
		t.Errorf("valid request should pass: %v", err)
	}

	req.Header.Del("X-Deploy-Token")
	if err := validator.ValidateRequest(req); err == nil {
		t.Error("missing token should fail")
	}

	req.Header.Set("X-Deploy-Token", "invalid")
	if err := validator.ValidateRequest(req); err == nil {
		t.Error("invalid token should fail")
	}
}
