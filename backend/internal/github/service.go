// SPDX-License-Identifier: MIT
package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Service struct {
	token   string
	client  *http.Client
	baseURL string
}

type Organization struct {
	Login  string `json:"login"`
	ID     int64  `json:"id"`
	URL    string `json:"url"`
	Avatar string `json:"avatar_url"`
}

type Repository struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Description   string `json:"description"`
	HTMLURL       string `json:"html_url"`
	CloneURL      string `json:"clone_url"`
	DefaultBranch string `json:"default_branch"`
	Language      string `json:"language"`
	Private       bool   `json:"private"`
	Fork          bool   `json:"fork"`
	Organization  string `json:"organization,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
}

type Branch struct {
	Name      string `json:"name"`
	CommitSHA string `json:"commit_sha,omitempty"`
}

func New(token string) *Service {
	return &Service{
		token:   token,
		client:  &http.Client{Timeout: 15 * time.Second},
		baseURL: "https://api.github.com",
	}
}

func (s *Service) ghURL(path string) string {
	return s.baseURL + path
}

func (s *Service) ValidateToken() (string, error) {
	req, err := http.NewRequest("GET", s.ghURL("/user"), nil)
	if err != nil {
		return "", fmt.Errorf("github request: %w", err)
	}
	req.Header.Set("Authorization", "token "+s.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("github connection: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github auth failed: HTTP %d", resp.StatusCode)
	}

	var user struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", fmt.Errorf("github decode: %w", err)
	}

	return user.Login, nil
}

func (s *Service) ListOrganizations() ([]Organization, error) {
	req, _ := http.NewRequest("GET", s.ghURL("/user/orgs"), nil)
	req.Header.Set("Authorization", "token "+s.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github orgs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github orgs: HTTP %d", resp.StatusCode)
	}

	var orgs []Organization
	if err := json.NewDecoder(resp.Body).Decode(&orgs); err != nil {
		return nil, fmt.Errorf("github decode: %w", err)
	}

	return orgs, nil
}

func (s *Service) ListRepositories(org string) ([]Repository, error) {
	url := s.ghURL(fmt.Sprintf("/orgs/%s/repos?per_page=100&type=all", org))
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "token "+s.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github repos: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github repos: HTTP %d", resp.StatusCode)
	}

	var repos []Repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("github decode: %w", err)
	}

	for i := range repos {
		repos[i].Organization = org
		if repos[i].DefaultBranch == "" {
			repos[i].DefaultBranch = "main"
		}
	}

	return repos, nil
}

func (s *Service) ListUserRepositories() ([]Repository, error) {
	req, _ := http.NewRequest("GET", s.ghURL("/user/repos?per_page=100&type=owner"), nil)
	req.Header.Set("Authorization", "token "+s.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github user repos: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github user repos: HTTP %d", resp.StatusCode)
	}

	var repos []Repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("github decode: %w", err)
	}

	for i := range repos {
		repos[i].Organization = repos[i].Owner()
		if repos[i].DefaultBranch == "" {
			repos[i].DefaultBranch = "main"
		}
	}

	return repos, nil
}

func (r *Repository) Owner() string {
	if r.FullName == "" {
		return ""
	}
	parts := splitPath(r.FullName)
	if len(parts) >= 1 {
		return parts[0]
	}
	return ""
}

func (r *Repository) RepoName() string {
	if r.FullName == "" {
		return r.Name
	}
	parts := splitPath(r.FullName)
	if len(parts) >= 2 {
		return parts[1]
	}
	return r.Name
}

func splitPath(s string) []string {
	var parts []string
	current := ""
	for _, c := range s {
		if c == '/' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
