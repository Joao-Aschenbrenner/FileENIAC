// SPDX-License-Identifier: MIT
package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	workspace.Active().DB.Close()
}

func TestWorkspace_PrepareCreatesWorkspaceInSelectedFolder(t *testing.T) {
	tmpDir := t.TempDir()
	selected := filepath.Join(tmpDir, "ENIAC_SYSTEMS")

	srv := New(":0")
	payload := map[string]string{"path": selected}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/workspace", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if _, err := os.Stat(filepath.Join(selected, ".eniac", "config.toml")); err != nil {
		t.Fatalf("expected workspace config to be created: %v", err)
	}

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["name"] != "ENIAC_SYSTEMS" {
		t.Fatalf("expected default workspace name ENIAC_SYSTEMS, got %v", resp["name"])
	}
	if resp["path"] != selected {
		t.Fatalf("expected workspace path %s, got %v", selected, resp["path"])
	}
	workspace.Active().DB.Close()
}

func TestWorkspaces_ListRootWithoutCreatingWorkspaceAtRoot(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "ENIAC_SYSTEMS")
	workspacePath := filepath.Join(root, "ClienteA")
	if _, err := workspace.Init("ClienteA", workspacePath, ""); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	workspace.Active().DB.Close()

	srv := New(":0")
	req := httptest.NewRequest("GET", "/api/workspaces?root="+url.QueryEscape(root), nil)
	w := httptest.NewRecorder()

	srv.mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if _, err := os.Stat(filepath.Join(root, ".eniac")); !os.IsNotExist(err) {
		t.Fatalf("expected root to stay as allocation folder without .eniac, got err=%v", err)
	}

	var resp []workspaceSummary
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 workspace, got %d", len(resp))
	}
	if resp[0].Name != "ClienteA" || resp[0].Path != workspacePath {
		t.Fatalf("unexpected workspace summary: %+v", resp[0])
	}
}

func TestWorkspaces_ListCreatesBaseFolderOnly(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "new-base")

	srv := New(":0")
	req := httptest.NewRequest("GET", "/api/workspaces?root="+url.QueryEscape(root), nil)
	w := httptest.NewRecorder()

	srv.mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if _, err := os.Stat(root); err != nil {
		t.Fatalf("expected base folder to be created: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".eniac")); !os.IsNotExist(err) {
		t.Fatalf("expected no workspace at base folder, got err=%v", err)
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
