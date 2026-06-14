package bypass

import (
	"strings"
	"testing"
)

func TestRenamer_GenerateEndpointName(t *testing.T) {
	renamer := NewRenamer()

	name1, err := renamer.GenerateEndpointName()
	if err != nil {
		t.Fatalf("GenerateEndpointName failed: %v", err)
	}

	if !strings.HasPrefix(name1, "_deploy_") {
		t.Errorf("expected prefix '_deploy_', got '%s'", name1)
	}

	if !strings.HasSuffix(name1, ".php") {
		t.Errorf("expected suffix '.php', got '%s'", name1)
	}

	name2, _ := renamer.GenerateEndpointName()
	if name1 == name2 {
		t.Error("two generated names should be different")
	}
}

func TestRenamer_GetTriggerURL(t *testing.T) {
	renamer := NewRenamer()

	tests := []struct {
		baseURL      string
		endpointName string
		expected     string
	}{
		{"https://example.com/projects/test", "_deploy_abc123.php", "https://example.com/projects/test/_deploy_abc123.php"},
		{"https://example.com/", "_deploy_abc.php", "https://example.com/_deploy_abc.php"},
		{"https://example.com", "/_deploy_abc.php", "https://example.com/_deploy_abc.php"},
	}

	for _, tt := range tests {
		got := renamer.GetTriggerURL(tt.baseURL, tt.endpointName)
		if got != tt.expected {
			t.Errorf("GetTriggerURL(%q, %q) = %q, want %q", tt.baseURL, tt.endpointName, got, tt.expected)
		}
	}
}

func TestRenamer_ShouldUseFallback(t *testing.T) {
	renamer := NewRenamer()

	if !renamer.ShouldUseFallback("_deploy_abc.php") {
		t.Error("'deploy' in name should trigger fallback")
	}

	if !renamer.ShouldUseFallback("admin_unt ar.php") {
		t.Error("'untar' in name should trigger fallback")
	}

	if renamer.ShouldUseFallback("_x7k9m.php") {
		t.Error("random name without suspicious patterns should not trigger fallback")
	}
}