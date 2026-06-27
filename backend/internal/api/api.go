package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/ENIACSystems/FileENIAC/backend/internal/heartbeat"

	"github.com/ENIACSystems/FileENIAC/backend/internal/clone"
	"github.com/ENIACSystems/FileENIAC/backend/internal/deploy"
	"github.com/ENIACSystems/FileENIAC/backend/internal/diff"
	gh "github.com/ENIACSystems/FileENIAC/backend/internal/github"
	bghealth "github.com/ENIACSystems/FileENIAC/backend/internal/health"
	"github.com/ENIACSystems/FileENIAC/backend/internal/history"
	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/mirror"
	"github.com/ENIACSystems/FileENIAC/backend/internal/readiness"
	"github.com/ENIACSystems/FileENIAC/backend/internal/refresh"
	"github.com/ENIACSystems/FileENIAC/backend/internal/registry"
	"github.com/ENIACSystems/FileENIAC/backend/internal/repair"
	"github.com/ENIACSystems/FileENIAC/backend/internal/status"
	syncpkg "github.com/ENIACSystems/FileENIAC/backend/internal/sync"
	"github.com/ENIACSystems/FileENIAC/backend/internal/validate"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"go.uber.org/zap"
)

type Server struct {
	addr       string
	mux        *http.ServeMux
	mu         sync.RWMutex
	srv        *http.Server
	background *bghealth.BackgroundRunner
}

func New(addr string) *Server {
	s := &Server{
		addr: addr,
		mux:  http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) SetBackgroundRunner(bg *bghealth.BackgroundRunner) {
	s.background = bg
}

func (s *Server) routes() {
	s.mux.HandleFunc("/api/health", s.handleHealth())
	s.mux.HandleFunc("/api/workspace", s.requireWorkspace(s.handleWorkspace()))
	s.mux.HandleFunc("/api/projects", s.requireWorkspace(s.handleProjects()))
	s.mux.HandleFunc("/api/projects/", s.requireWorkspace(s.handleProjectByID()))
	s.mux.HandleFunc("/api/servers", s.requireWorkspace(s.handleServers()))
	s.mux.HandleFunc("/api/servers/", s.requireWorkspace(s.handleServerByID()))
	s.mux.HandleFunc("/api/settings", s.requireWorkspace(s.handleSettings()))
	s.mux.HandleFunc("/api/history", s.requireWorkspace(s.handleHistory()))
	s.mux.HandleFunc("/api/events", s.requireWorkspace(s.handleEvents()))
	s.mux.HandleFunc("/api/deploys", s.requireWorkspace(s.handleDeploys()))
	s.mux.HandleFunc("/api/deploy", s.requireWorkspace(s.handleDeployExec()))
	s.mux.HandleFunc("/api/rollback", s.requireWorkspace(s.handleRollback()))
	s.mux.HandleFunc("/api/verify", s.requireWorkspace(s.handleVerify()))
	s.mux.HandleFunc("/api/diff", s.requireWorkspace(s.handleDiff()))
	s.mux.HandleFunc("/api/syncs", s.requireWorkspace(s.handleSyncs()))
	s.mux.HandleFunc("/api/sync", s.requireWorkspace(s.handleSyncExec()))
	s.mux.HandleFunc("/api/mirror", s.requireWorkspace(s.handleMirror()))
	s.mux.HandleFunc("/api/health/check", s.requireWorkspace(s.handleHealthCheck()))
	s.mux.HandleFunc("/api/health/background", s.requireWorkspace(s.handleBackgroundHealth()))
	s.mux.HandleFunc("/api/github/status", s.requireWorkspace(s.handleGitHubStatus()))
	s.mux.HandleFunc("/api/github/login", s.requireWorkspace(s.handleGitHubLogin()))
	s.mux.HandleFunc("/api/github/logout", s.requireWorkspace(s.handleGitHubLogout()))
	s.mux.HandleFunc("/api/github/organizations", s.requireWorkspace(s.handleGitHubOrganizations()))
	s.mux.HandleFunc("/api/github/repositories", s.requireWorkspace(s.handleGitHubRepositories()))
	s.mux.HandleFunc("/api/github/import", s.requireWorkspace(s.handleGitHubImport()))
	s.mux.HandleFunc("/api/github/clone", s.requireWorkspace(s.handleGitHubClone()))
	s.mux.HandleFunc("/api/repositories", s.requireWorkspace(s.handleRepositories()))
	s.mux.HandleFunc("/api/repositories/", s.requireWorkspace(s.handleRepositoryByID()))
	s.mux.HandleFunc("/api/refresh/github", s.requireWorkspace(s.handleRefreshGitHub()))
	s.mux.HandleFunc("/api/revalidate", s.requireWorkspace(s.handleRevalidate()))
	s.mux.HandleFunc("/api/readiness/deploy", s.requireWorkspace(s.handleReadinessDeploy()))
	s.mux.HandleFunc("/api/readiness/sync", s.requireWorkspace(s.handleReadinessSync()))
	s.mux.HandleFunc("/api/repair/check", s.requireWorkspace(s.handleRepairCheck()))
	s.mux.HandleFunc("/api/repair/fix", s.requireWorkspace(s.handleRepairFix()))
	s.mux.HandleFunc("/api/heartbeat", s.handleHeartbeat())
}

func (s *Server) ListenAndServe() error {
	srv := &http.Server{Addr: s.addr, Handler: s.corsMiddleware(s.mux)}
	s.mu.Lock()
	s.srv = srv
	s.mu.Unlock()
	return srv.ListenAndServe()
}

func (s *Server) ListenDynamic() (string, error) {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return "", err
	}
	actualAddr := listener.Addr().String()
	log.L().Info("api server listening (dynamic)", zap.String("addr", actualAddr))
	srv := &http.Server{Handler: s.corsMiddleware(s.mux)}
	s.mu.Lock()
	s.srv = srv
	s.mu.Unlock()
	go srv.Serve(listener)
	return actualAddr, nil
}

func (s *Server) handleHeartbeat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondError(w, http.StatusMethodNotAllowed, "use POST")
			return
		}
		heartbeat.Reset()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) Close() error {
	s.mu.RLock()
	srv := s.srv
	s.mu.RUnlock()
	if srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(ctx)
	}
	return nil
}

func (s *Server) Addr() string {
	s.mu.RLock()
	srv := s.srv
	s.mu.RUnlock()
	if srv != nil {
		return srv.Addr
	}
	return s.addr
}

func respond(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respond(w, status, map[string]string{"error": msg})
}

func (s *Server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func (s *Server) requireWorkspace(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := workspace.Active()
		if ctx == nil {
			wsPath := r.URL.Query().Get("workspace")
			if wsPath == "" {
				pwd, _ := os.Getwd()
				wsPath = pwd
			}
			_, err := workspace.Open(wsPath)
			if err != nil {
				respondError(w, http.StatusNotFound, fmt.Sprintf("workspace not found: %v", err))
				return
			}
			ctx = workspace.Active()
		}
		_ = ctx
		next(w, r)
	}
}

// GET /workspace
func (s *Server) handleWorkspace() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := workspace.Active()
		if ctx == nil {
			respondError(w, http.StatusNotFound, "no active workspace")
			return
		}
		respond(w, http.StatusOK, ctx.Workspace.Status())
	}
}

// GET/POST /projects
func (s *Server) handleProjects() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := workspace.Active()
		switch r.Method {
		case http.MethodGet:
			projects, err := registry.ListProjects(ctx)
			if err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respond(w, http.StatusOK, projects)
		case http.MethodPost:
			var p registry.Project
			if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
				respondError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
				return
			}
			id, err := registry.AddProject(ctx, &p)
			if err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			p.ID = id
			respond(w, http.StatusCreated, p)
		default:
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

// GET/DELETE /projects/:name
func (s *Server) handleProjectByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := workspace.Active()
		name := filepath.Base(r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			p, err := registry.GetProject(ctx, name)
			if err != nil {
				respondError(w, http.StatusNotFound, err.Error())
				return
			}
			respond(w, http.StatusOK, p)
		case http.MethodDelete:
			p, err := registry.GetProject(ctx, name)
			if err != nil {
				respondError(w, http.StatusNotFound, err.Error())
				return
			}
			if err := registry.RemoveProject(ctx, p.ID); err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respond(w, http.StatusOK, map[string]string{"status": "removed"})
		default:
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

// GET/POST /servers?project=
func (s *Server) handleServers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := workspace.Active()
		projectName := r.URL.Query().Get("project")
		switch r.Method {
		case http.MethodGet:
			var servers []*registry.Server
			var err error
			if projectName != "" {
				p, pErr := registry.GetProject(ctx, projectName)
				if pErr != nil {
					respondError(w, http.StatusNotFound, pErr.Error())
					return
				}
				servers, err = registry.ListServersByProject(ctx, p.ID)
			} else {
				servers, err = registry.ListServers(ctx)
			}
			if err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			for _, srv := range servers {
				srv.Password = ""
			}
			respond(w, http.StatusOK, servers)
		case http.MethodPost:
			var srv registry.Server
			if err := json.NewDecoder(r.Body).Decode(&srv); err != nil {
				respondError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
				return
			}
			id, err := registry.AddServer(ctx, &srv)
			if err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			srv.ID = id
			srv.Password = ""
			respond(w, http.StatusCreated, srv)
		default:
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

// GET/DELETE /servers/:id
func (s *Server) handleServerByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := workspace.Active()
		idStr := filepath.Base(r.URL.Path)
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid server ID")
			return
		}
		switch r.Method {
		case http.MethodGet:
			srv, err := registry.GetServerByID(ctx, id)
			if err != nil {
				respondError(w, http.StatusNotFound, err.Error())
				return
			}
			srv.Password = ""
			respond(w, http.StatusOK, srv)
		case http.MethodDelete:
			if err := registry.RemoveServer(ctx, id); err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respond(w, http.StatusOK, map[string]string{"status": "removed"})
		default:
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

// GET/POST /settings
func (s *Server) handleSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := workspace.Active()
		switch r.Method {
		case http.MethodGet:
			settings, err := ctx.DB.ListSettings()
			if err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respond(w, http.StatusOK, settings)
		case http.MethodPost:
			var body map[string]string
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				respondError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
				return
			}
			for k, v := range body {
				if err := ctx.DB.SetSetting(k, v); err != nil {
					respondError(w, http.StatusInternalServerError, err.Error())
					return
				}
			}
			respond(w, http.StatusOK, map[string]string{"status": "saved"})
		default:
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

// GET /history?project=&type=&limit=&offset=
func (s *Server) handleHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := workspace.Active()
		hs := history.NewService(ctx.DB)

		projectName := r.URL.Query().Get("project")
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit <= 0 || limit > 200 {
			limit = 50
		}
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		if offset < 0 {
			offset = 0
		}
		eventType := r.URL.Query().Get("type")

		if projectName != "" {
			proj, err := registry.GetProject(ctx, projectName)
			if err != nil {
				respondError(w, http.StatusNotFound, err.Error())
				return
			}
			logs, err := hs.GetDeployHistory(proj.ID, limit)
			if err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respond(w, http.StatusOK, logs)
			return
		}

		events, err := hs.GetEventList(eventType, limit, offset)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respond(w, http.StatusOK, events)
	}
}

// GET /events?type=&limit=&offset=
func (s *Server) handleEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := workspace.Active()
		hs := history.NewService(ctx.DB)

		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit <= 0 || limit > 200 {
			limit = 50
		}
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		if offset < 0 {
			offset = 0
		}
		eventType := r.URL.Query().Get("type")

		events, err := hs.GetEventList(eventType, limit, offset)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respond(w, http.StatusOK, events)
	}
}

// GET /deploys?project=&limit=
func (s *Server) handleDeploys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := workspace.Active()
		hs := history.NewService(ctx.DB)

		projectName := r.URL.Query().Get("project")
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit <= 0 || limit > 200 {
			limit = 20
		}

		if projectName == "" {
			respondError(w, http.StatusBadRequest, "project query param required")
			return
		}

		proj, err := registry.GetProject(ctx, projectName)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}

		logs, err := hs.GetDeployHistory(proj.ID, limit)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respond(w, http.StatusOK, logs)
	}
}

// POST /deploy {"project":"name"}
func (s *Server) handleDeployExec() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		ctx := workspace.Active()
		var body struct {
			Project     string `json:"project"`
			UseFallback bool   `json:"use_fallback"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			respondError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}
		if body.Project == "" {
			respondError(w, http.StatusBadRequest, "project required")
			return
		}

		svc := deploy.NewService(ctx.DB)
		result, err := svc.Deploy(ctx, body.Project, body.UseFallback)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respond(w, http.StatusOK, result)
	}
}

// POST /rollback {"project":"name"}
func (s *Server) handleRollback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		ctx := workspace.Active()
		var body struct {
			Project string `json:"project"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			respondError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}
		if body.Project == "" {
			respondError(w, http.StatusBadRequest, "project required")
			return
		}

		svc := deploy.NewService(ctx.DB)
		result, err := svc.Rollback(ctx, body.Project)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respond(w, http.StatusOK, result)
	}
}

// POST /verify {"project":"name"}
func (s *Server) handleVerify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		ctx := workspace.Active()
		var body struct {
			Project string `json:"project"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			respondError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}
		if body.Project == "" {
			respondError(w, http.StatusBadRequest, "project required")
			return
		}

		svc := deploy.NewService(ctx.DB)
		result, err := svc.Verify(ctx, body.Project)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respond(w, http.StatusOK, result)
	}
}

// GET /diff?project=
func (s *Server) handleDiff() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := workspace.Active()
		projectName := r.URL.Query().Get("project")
		if projectName == "" {
			respondError(w, http.StatusBadRequest, "project query param required")
			return
		}

		e := diff.New()
		report, err := e.LocalVsMirror(ctx, projectName)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respond(w, http.StatusOK, report)
	}
}

// GET /syncs?project=&limit=
func (s *Server) handleSyncs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := workspace.Active()
		projectName := r.URL.Query().Get("project")
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit <= 0 || limit > 200 {
			limit = 20
		}

		var manifests []map[string]interface{}
		var err error
		if projectName != "" {
			manifests, err = ctx.DB.QueryManifestsByProject(projectName, limit)
		} else {
			manifests, err = ctx.DB.QueryManifests(limit)
		}
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respond(w, http.StatusOK, manifests)
	}
}

// POST /sync {"project":"name","action":"mirror_update","direction":"local_to_mirror","confirm":false}
func (s *Server) handleSyncExec() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		ctx := workspace.Active()
		var body struct {
			Project   string `json:"project"`
			Action    string `json:"action"`
			Direction string `json:"direction"`
			Confirm   bool   `json:"confirm"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			respondError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}
		if body.Project == "" {
			respondError(w, http.StatusBadRequest, "project required")
			return
		}
		if body.Direction == "" {
			body.Direction = "local_to_mirror"
		}

		e := diff.New()
		report, err := e.LocalVsMirror(ctx, body.Project)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		se := syncpkg.New()
		suggestion := se.Plan(report)

		if err := se.Apply(ctx, body.Project, report, body.Direction, body.Confirm); err != nil {
			se.GenerateManifest(ctx, body.Project, body.Action, report, "failed")
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		manifest, err := se.GenerateManifest(ctx, body.Project, body.Action, report, "completed")
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		se.Reconcile(ctx, body.Project, "sincronizado")

		respond(w, http.StatusOK, map[string]interface{}{
			"suggestion": suggestion,
			"manifest":   manifest,
			"diff":       report,
		})
	}
}

// POST /mirror {"project":"name"}
func (s *Server) handleMirror() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		ctx := workspace.Active()
		var body struct {
			Project string `json:"project"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			respondError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}
		if body.Project == "" {
			respondError(w, http.StatusBadRequest, "project required")
			return
		}

		me := mirror.New()
		snapshot, err := me.Create(ctx, body.Project)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respond(w, http.StatusOK, snapshot)
	}
}

// GET /health/check
func (s *Server) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := workspace.Active()
		hs := history.NewService(ctx.DB)

		projects, _ := registry.ListProjects(ctx)
		servers, _ := registry.ListServers(ctx)
		events, _ := hs.GetRecentEvents(5)

		status := "healthy"
		if len(projects) == 0 {
			status = "degraded"
		}

		divergentCount := 0
		for _, p := range projects {
			if p.DivergenceStatus == "divergente" || p.DivergenceStatus == "unknown" {
				divergentCount++
			}
		}

		respond(w, http.StatusOK, map[string]interface{}{
			"status":      status,
			"projects":    len(projects),
			"servers":     len(servers),
			"divergent":   divergentCount,
			"last_events": events,
		})
	}
}

func getGitHubToken(ctx *workspace.Context) (string, error) {
	token, err := ctx.DB.GetSetting("github_token")
	if err != nil || token == "" {
		return "", fmt.Errorf("not authenticated with GitHub")
	}
	v, err := registry.VaultFromCtx(ctx)
	if err != nil {
		return "", fmt.Errorf("vault: %w", err)
	}
	dec, err := v.Decrypt(token)
	if err != nil {
		return "", fmt.Errorf("token decrypt: %w", err)
	}
	return dec, nil
}

func githubService(ctx *workspace.Context) (*gh.Service, error) {
	token, err := getGitHubToken(ctx)
	if err != nil {
		return nil, err
	}
	return gh.New(token), nil
}

// GET /github/status
func (s *Server) handleGitHubStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := workspace.Active()
		token, err := ctx.DB.GetSetting("github_token")
		if err != nil || token == "" {
			respond(w, http.StatusOK, map[string]interface{}{
				"authenticated": false,
				"user":          "",
			})
			return
		}
		user, _ := ctx.DB.GetSetting("github_user")
		respond(w, http.StatusOK, map[string]interface{}{
			"authenticated": true,
			"user":          user,
		})
	}
}

// POST /github/login {"token":"ghp_..."}
func (s *Server) handleGitHubLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		ctx := workspace.Active()
		var body struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			respondError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
		if body.Token == "" {
			respondError(w, http.StatusBadRequest, "token required")
			return
		}
		svc := gh.New(body.Token)
		user, err := svc.ValidateToken()
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}
		v, err := registry.VaultFromCtx(ctx)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		enc, err := v.Encrypt(body.Token)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "encryption failed")
			return
		}
		if err := ctx.DB.SetSetting("github_token", enc); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err := ctx.DB.SetSetting("github_user", user); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		ctx.DB.RecordEvent("github_login", fmt.Sprintf("GitHub authenticated as %s", user), nil)
		log.L().Info("github login", zap.String("user", user))
		respond(w, http.StatusOK, map[string]string{"user": user})
	}
}

// POST /github/logout
func (s *Server) handleGitHubLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		ctx := workspace.Active()
		ctx.DB.SetSetting("github_token", "")
		ctx.DB.SetSetting("github_user", "")
		ctx.DB.RecordEvent("github_logout", "GitHub disconnected", nil)
		respond(w, http.StatusOK, map[string]string{"status": "logged_out"})
	}
}

// GET /github/organizations
func (s *Server) handleGitHubOrganizations() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := workspace.Active()
		svc, err := githubService(ctx)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}
		orgs, err := svc.ListOrganizations()
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respond(w, http.StatusOK, orgs)
	}
}

// GET /github/repositories?org=<org>
func (s *Server) handleGitHubRepositories() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := workspace.Active()
		svc, err := githubService(ctx)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}
		org := r.URL.Query().Get("org")
		var repos []gh.Repository
		if org != "" {
			repos, err = svc.ListRepositories(org)
		} else {
			repos, err = svc.ListUserRepositories()
		}
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Mark already imported repos
		imported, _ := registry.ListAllRepositories(ctx)
		importedMap := make(map[int64]bool)
		for _, ir := range imported {
			importedMap[ir.GitHubID] = true
		}
		type repoStatus struct {
			gh.Repository
			Imported bool `json:"imported"`
		}
		result := make([]repoStatus, len(repos))
		for i, repo := range repos {
			result[i] = repoStatus{Repository: repo, Imported: importedMap[repo.ID]}
		}

		respond(w, http.StatusOK, result)
	}
}

// POST /github/import {"repos":[{"github_id":...,"name":"...","clone_url":"...","default_branch":"...","full_name":"...","html_url":"...","organization":"..."}], "clone_dir":"C:/projects"}
func (s *Server) handleGitHubImport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		ctx := workspace.Active()
		var body struct {
			Repos    []gh.Repository `json:"repos"`
			CloneDir string          `json:"clone_dir"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			respondError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
		if len(body.Repos) == 0 {
			respondError(w, http.StatusBadRequest, "repos required")
			return
		}

		type importResult struct {
			Repo         gh.Repository    `json:"repo"`
			ProjectID    int64            `json:"project_id"`
			RepositoryID int64            `json:"repository_id"`
			CloneResult  *clone.Result    `json:"clone_result,omitempty"`
			Validation   *validate.Result `json:"validation,omitempty"`
			Error        string           `json:"error,omitempty"`
		}

		wsPath := ctx.Workspace.Path
		results := make([]importResult, 0, len(body.Repos))

		for _, repo := range body.Repos {
			r := importResult{Repo: repo}

			cloneDir := filepath.Join(wsPath, "projects", repo.Name)

			project, err := registry.AddProject(ctx, &registry.Project{
				Name:        repo.Name,
				DisplayName: repo.Name,
				LocalPath:   cloneDir,
				RemotePath:  "/",
				Branch:      repo.DefaultBranch,
				GitURL:      repo.CloneURL,
				Environment: "production",
			})
			if err != nil {
				r.Error = fmt.Sprintf("project registration: %s", err)
				results = append(results, r)
				continue
			}
			r.ProjectID = project

			repoID, err := registry.AddRepository(ctx, &registry.Repository{
				GitHubID:      repo.ID,
				Name:          repo.Name,
				FullName:      repo.FullName,
				Description:   repo.Description,
				HTMLURL:       repo.HTMLURL,
				CloneURL:      repo.CloneURL,
				DefaultBranch: repo.DefaultBranch,
				Language:      repo.Language,
				Private:       repo.Private,
				Organization:  repo.Organization,
			})
			if err != nil {
				r.Error = fmt.Sprintf("repository registration: %s", err)
				results = append(results, r)
				continue
			}
			r.RepositoryID = repoID

			cloneResult, err := clone.Clone(repo.CloneURL, cloneDir, repo.DefaultBranch)
			if err != nil {
				registry.UpdateRepositoryImport(ctx, repoID, project, "clone_failed", "")
				r.Error = fmt.Sprintf("clone: %s", err)
				results = append(results, r)
				continue
			}
			r.CloneResult = cloneResult

			registry.UpdateRepositoryImport(ctx, repoID, project, "cloned", cloneDir)

			vResult := validate.ValidateClone(cloneDir, repo.DefaultBranch)
			r.Validation = vResult

			importStatus := "imported"
			if !vResult.Valid {
				importStatus = "validation_failed"
			}
			ctx.DB.Exec(`UPDATE projects SET import_status=?, clone_path=?, organization=?, repo_name=?, provider='github', github_id=?, last_sync_commit=? WHERE id=?`,
				importStatus, cloneDir, repo.Organization, repo.Name, repo.ID, cloneResult.CommitSHA, project)

			ctx.DB.RecordEvent("github_import", fmt.Sprintf("Imported %s/%s", repo.Organization, repo.Name), map[string]interface{}{"project_id": project, "github_id": repo.ID})

			results = append(results, r)
		}

		respond(w, http.StatusOK, results)
	}
}

// POST /github/clone {"repo_id":...,"project_id":...,"clone_url":"...","branch":"...","clone_dir":"..."}
func (s *Server) handleGitHubClone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		ctx := workspace.Active()
		var body struct {
			RepoID    int64  `json:"repo_id"`
			ProjectID int64  `json:"project_id"`
			CloneURL  string `json:"clone_url"`
			Branch    string `json:"branch"`
			CloneDir  string `json:"clone_dir"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			respondError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
		if body.CloneURL == "" || body.CloneDir == "" {
			respondError(w, http.StatusBadRequest, "clone_url and clone_dir required")
			return
		}
		if body.Branch == "" {
			body.Branch = "main"
		}

		cloneResult, err := clone.Clone(body.CloneURL, body.CloneDir, body.Branch)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if body.RepoID > 0 {
			registry.UpdateRepositoryImport(ctx, body.RepoID, body.ProjectID, "cloned", body.CloneDir)
		}
		if body.ProjectID > 0 {
			ctx.DB.Exec(`UPDATE projects SET clone_path=?, import_status='cloned', last_sync_commit=? WHERE id=?`,
				body.CloneDir, cloneResult.CommitSHA, body.ProjectID)
		}

		respond(w, http.StatusOK, cloneResult)
	}
}

// GET /repositories?org=
func (s *Server) handleRepositories() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := workspace.Active()
		org := r.URL.Query().Get("org")
		var repos []*registry.Repository
		var err error
		if org != "" {
			repos, err = registry.ListRepositoriesByOrg(ctx, org)
		} else {
			repos, err = registry.ListAllRepositories(ctx)
		}
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respond(w, http.StatusOK, repos)
	}
}

// GET /repositories/:id
func (s *Server) handleRepositoryByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := workspace.Active()
		idStr := filepath.Base(r.URL.Path)
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid repo ID")
			return
		}
		repo, err := registry.GetRepositoryByGitHubID(ctx, id)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respond(w, http.StatusOK, repo)
	}
}

// POST /refresh/github
func (s *Server) handleRefreshGitHub() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		ctx := workspace.Active()
		result, err := refresh.RefreshGitHub(ctx)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		ctx.DB.RecordEvent("github_refresh", fmt.Sprintf("Refreshed %d orgs, %d repos, %d changes", result.Organizations, result.Repositories, result.ChangesFound), nil)
		respond(w, http.StatusOK, result)
	}
}

// POST /revalidate
func (s *Server) handleRevalidate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		ctx := workspace.Active()
		se := status.New()
		projects, _ := registry.ListProjects(ctx)
		repos, _ := registry.ListAllRepositories(ctx)

		type revalidationItem struct {
			Project    string             `json:"project"`
			Status     status.ImportState `json:"status"`
			CloneCheck *validate.Result   `json:"clone_check,omitempty"`
			Integrity  *validate.Result   `json:"integrity,omitempty"`
		}

		items := make([]revalidationItem, 0)
		for _, p := range projects {
			item := revalidationItem{Project: p.Name}
			item.CloneCheck = validate.ValidateClone(p.LocalPath, p.Branch)
			item.Integrity = validate.CloneIntegrity(p.LocalPath, p.Branch, p.GitURL)
			checks := []status.CheckResult{
				{Name: "clone_exists", Passed: item.CloneCheck.Valid},
			}
			importStatus := "imported"
			item.Status = se.Resolve(importStatus, checks)
			items = append(items, item)
		}

		respond(w, http.StatusOK, map[string]interface{}{
			"projects_checked": len(projects),
			"repos_checked":    len(repos),
			"items":            items,
		})
	}
}

// GET /readiness/deploy?project=
func (s *Server) handleReadinessDeploy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := workspace.Active()
		projectName := r.URL.Query().Get("project")
		result := readiness.CheckDeploy(ctx, projectName)
		respond(w, http.StatusOK, result)
	}
}

// GET /readiness/sync?project=
func (s *Server) handleReadinessSync() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := workspace.Active()
		projectName := r.URL.Query().Get("project")
		result := readiness.CheckSync(ctx, projectName)
		respond(w, http.StatusOK, result)
	}
}

// GET /repair/check
func (s *Server) handleRepairCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := workspace.Active()
		report := repair.CheckConsistency(ctx)
		respond(w, http.StatusOK, report)
	}
}

// POST /repair/fix
func (s *Server) handleRepairFix() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		ctx := workspace.Active()
		report, err := repair.RepairOrphanedRepositories(ctx)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		ctx.DB.RecordEvent("repair_fix", fmt.Sprintf("Repaired %d issues, %d warnings", report.Fixed, len(report.Warnings)), nil)
		respond(w, http.StatusOK, report)
	}
}

// GET /health/background
func (s *Server) handleBackgroundHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snap := bghealth.GetSnapshot()
		if snap.Timestamp.IsZero() {
			snap.Timestamp = time.Now()
			snap.Status = "unknown"
		}
		respond(w, http.StatusOK, snap)
	}
}
