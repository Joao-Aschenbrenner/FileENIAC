// SPDX-License-Identifier: MIT
package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
)

func TestHealth(t *testing.T) {
	srv := New(":0")
	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	srv.mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	if body["status"] != "ok" {
		t.Errorf("expected ok, got %s", body["status"])
	}
}

func TestCORSPermitsTauriPreflightWithWorkspaceHeader(t *testing.T) {
	srv := New(":0")
	handler := srv.corsMiddleware(srv.authMiddleware(srv.rateLimitMiddleware(srv.mux)))

	req := httptest.NewRequest("OPTIONS", "/api/projects", nil)
	req.Header.Set("Origin", "http://tauri.localhost")
	req.Header.Set("Access-Control-Request-Headers", "authorization,x-workspace")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected preflight 200, got %d", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://tauri.localhost" {
		t.Fatalf("expected tauri origin echo, got %q", got)
	}
	if got := w.Header().Get("Access-Control-Allow-Headers"); !strings.Contains(got, "X-Workspace") {
		t.Fatalf("expected X-Workspace in allow headers, got %q", got)
	}
}

func TestWorkspace_NoWorkspace(t *testing.T) {
	srv := New(":0")
	req := httptest.NewRequest("GET", "/api/workspace", nil)
	w := httptest.NewRecorder()
	srv.mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestWorkspace_WithContext(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	_, err := workspace.Init("APITest", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	workspace.Active().DB.Close()

	srv := New(":0")
	req := httptest.NewRequest("GET", "/api/workspace?workspace="+wsPath, nil)
	w := httptest.NewRecorder()
	srv.mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var body map[string]interface{}
	json.NewDecoder(w.Body).Decode(&body)
	if body["name"] != "APITest" {
		t.Errorf("expected APITest, got %v", body["name"])
	}
}

func TestProjects_CRUD(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	_, err := workspace.Init("ProjAPITest", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	srv := New(":0")
	baseURL := "?workspace=" + wsPath

	projectPath := filepath.Join(tmpDir, "myproject")
	os.MkdirAll(projectPath, 0700)

	// Create project via POST
	payload := map[string]string{
		"name":       "api-project",
		"local_path": projectPath,
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/projects"+baseURL, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	// List projects
	listReq := httptest.NewRequest("GET", "/api/projects"+baseURL, nil)
	lw := httptest.NewRecorder()
	srv.mux.ServeHTTP(lw, listReq)

	if lw.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", lw.Code)
	}

	var projects []map[string]interface{}
	json.NewDecoder(lw.Body).Decode(&projects)
	if len(projects) == 0 {
		t.Error("expected at least 1 project")
	}
}

func TestSettings(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	_, err := workspace.Init("SettingsAPITest", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	srv := New(":0")
	baseURL := "?workspace=" + wsPath

	req := httptest.NewRequest("GET", "/api/settings"+baseURL, nil)
	w := httptest.NewRecorder()
	srv.mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
