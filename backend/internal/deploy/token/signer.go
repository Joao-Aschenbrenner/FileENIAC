package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

const TokenTTL = 5 * time.Minute

type Signer struct {
	secret []byte
}

func NewSigner(secret string) *Signer {
	return &Signer{
		secret: []byte(secret),
	}
}

func (s *Signer) Sign(projectName string) (string, time.Time, error) {
	timestamp := time.Now()
	message := fmt.Sprintf("%d/%s", timestamp.Unix(), projectName)

	h := hmac.New(sha256.New, s.secret)
	h.Write([]byte(message))
	signature := hex.EncodeToString(h.Sum(nil))

	return signature, timestamp, nil
}

func (s *Signer) Validate(projectName, signature string, timestamp time.Time) bool {
	if time.Since(timestamp) > TokenTTL {
		return false
	}

	message := fmt.Sprintf("%d/%s", timestamp.Unix(), projectName)
	h := hmac.New(sha256.New, s.secret)
	h.Write([]byte(message))
	expected := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expected))
}

func (s *Signer) GenerateHeaders(projectName string) (map[string]string, error) {
	sig, ts, err := s.Sign(projectName)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"X-Deploy-Token":     sig,
		"X-Deploy-Timestamp": fmt.Sprintf("%d", ts.Unix()),
		"X-Deploy-Project":   projectName,
	}, nil
}
