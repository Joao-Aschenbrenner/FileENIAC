package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew_SetsToken(t *testing.T) {
	svc := New("test-token-123")
	if svc == nil {
		t.Fatal("New returned nil")
	}
	if svc.token != "test-token-123" {
		t.Errorf("expected token 'test-token-123', got %q", svc.token)
	}
}

func TestNew_EmptyToken(t *testing.T) {
	svc := New("")
	if svc == nil {
		t.Fatal("New returned nil")
	}
}

func TestValidateToken_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "token valid-token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if r.Header.Get("Accept") != "application/vnd.github.v3+json" {
			http.Error(w, "bad accept", http.StatusNotAcceptable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"login": "testuser"})
	}))
	defer server.Close()

	svc := New("valid-token")
	svc.client = server.Client()
	svc.baseURL = server.URL

	user, err := svc.ValidateToken()
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if user != "testuser" {
		t.Errorf("expected 'testuser', got %q", user)
	}
}

func TestValidateToken_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	svc := New("bad-token")
	svc.client = server.Client()
	svc.baseURL = server.URL

	_, err := svc.ValidateToken()
	if err == nil {
		t.Fatal("expected error for unauthorized token")
	}
}

func TestValidateToken_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	svc := New("ok-token")
	svc.client = server.Client()
	svc.baseURL = server.URL

	_, err := svc.ValidateToken()
	if err == nil {
		t.Fatal("expected error for server error")
	}
}

func TestListOrganizations_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "token org-token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"login": "org1", "id": 1, "url": "https://api.github.com/orgs/org1", "avatar_url": ""},
		})
	}))
	defer server.Close()

	svc := New("org-token")
	svc.client = server.Client()
	svc.baseURL = server.URL

	orgs, err := svc.ListOrganizations()
	if err != nil {
		t.Fatalf("ListOrganizations failed: %v", err)
	}
	if len(orgs) != 1 {
		t.Fatalf("expected 1 org, got %d", len(orgs))
	}
	if orgs[0].Login != "org1" {
		t.Errorf("expected 'org1', got %q", orgs[0].Login)
	}
}

func TestListOrganizations_Empty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]map[string]interface{}{})
	}))
	defer server.Close()

	svc := New("empty-org-token")
	svc.client = server.Client()
	svc.baseURL = server.URL

	orgs, err := svc.ListOrganizations()
	if err != nil {
		t.Fatalf("ListOrganizations failed: %v", err)
	}
	if len(orgs) != 0 {
		t.Errorf("expected 0 orgs, got %d", len(orgs))
	}
}

func TestListOrganizations_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	svc := New("ok-token")
	svc.client = server.Client()
	svc.baseURL = server.URL

	_, err := svc.ListOrganizations()
	if err == nil {
		t.Fatal("expected error for server error")
	}
}

func TestRepository_Owner(t *testing.T) {
	r := &Repository{FullName: "owner/repo"}
	if r.Owner() != "owner" {
		t.Errorf("expected 'owner', got %q", r.Owner())
	}
}

func TestRepository_Owner_Empty(t *testing.T) {
	r := &Repository{FullName: ""}
	if r.Owner() != "" {
		t.Errorf("expected '', got %q", r.Owner())
	}
}

func TestRepository_RepoName(t *testing.T) {
	r := &Repository{FullName: "owner/repo"}
	if r.RepoName() != "repo" {
		t.Errorf("expected 'repo', got %q", r.RepoName())
	}
}

func TestRepository_RepoName_NoFullName(t *testing.T) {
	r := &Repository{Name: "myrepo", FullName: ""}
	if r.RepoName() != "myrepo" {
		t.Errorf("expected 'myrepo', got %q", r.RepoName())
	}
}
