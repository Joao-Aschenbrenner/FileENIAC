// SPDX-License-Identifier: MIT
package bypass

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

type Renamer struct{}

func NewRenamer() *Renamer {
	return &Renamer{}
}

func (r *Renamer) GenerateEndpointName() (string, error) {
	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	randomSuffix := hex.EncodeToString(bytes)
	return fmt.Sprintf("_deploy_%s.php", randomSuffix), nil
}

func (r *Renamer) GetTriggerURL(baseURL, endpointName string) string {
	baseURL = strings.TrimSuffix(baseURL, "/")
	endpointName = strings.TrimPrefix(endpointName, "/")
	return baseURL + "/" + endpointName
}

func (r *Renamer) ShouldUseFallback(endpointName string) bool {
	suspiciousPatterns := []string{
		"deploy",
		"upload",
		"untar",
		"extract",
		"admin",
	}

	name := strings.ToLower(endpointName)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(name, pattern) {
			return true
		}
	}
	return false
}
