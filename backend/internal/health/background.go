// SPDX-License-Identifier: MIT
package health

import (
	"sync"
	"time"

	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/registry"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"go.uber.org/zap"
)

type BackgroundRunner struct {
	interval time.Duration
	stopCh   chan struct{}
	wg       sync.WaitGroup
	running  bool
	mu       sync.Mutex
}

type Snapshot struct {
	Timestamp      time.Time `json:"timestamp"`
	TokenValid     bool      `json:"token_valid"`
	GitHubUser     string    `json:"github_user,omitempty"`
	ProjectsCount  int       `json:"projects_count"`
	ServersCount   int       `json:"servers_count"`
	ClonesValid    int       `json:"clones_valid"`
	ClonesBroken   int       `json:"clones_broken"`
	DivergentCount int       `json:"divergent_count"`
	Status         string    `json:"status"`
}

var (
	currentSnapshot Snapshot
	snapshotMu      sync.RWMutex
)

func NewBackgroundRunner(interval time.Duration) *BackgroundRunner {
	return &BackgroundRunner{
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

func (b *BackgroundRunner) Start(ctx *workspace.Context) {
	b.mu.Lock()
	if b.running {
		b.mu.Unlock()
		return
	}
	b.running = true
	b.mu.Unlock()

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		ticker := time.NewTicker(b.interval)
		defer ticker.Stop()

		b.runCheck(ctx)
		for {
			select {
			case <-ticker.C:
				b.runCheck(ctx)
			case <-b.stopCh:
				log.L().Info("background health runner stopped")
				return
			}
		}
	}()
	log.L().Info("background health runner started", zap.Duration("interval", b.interval))
}

func (b *BackgroundRunner) Stop() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.running {
		return
	}
	close(b.stopCh)
	b.running = false
	b.wg.Wait()
}

func (b *BackgroundRunner) runCheck(ctx *workspace.Context) {
	snap := Snapshot{Timestamp: time.Now()}

	token, err := ctx.DB.GetSetting("github_token")
	if err == nil && token != "" {
		snap.TokenValid = true
		snap.GitHubUser, _ = ctx.DB.GetSetting("github_user")
	}

	projects, _ := registry.ListProjects(ctx)
	snap.ProjectsCount = len(projects)

	for _, p := range projects {
		if p.DivergenceStatus == "divergente" || p.DivergenceStatus == "unknown" {
			snap.DivergentCount++
		}
	}

	servers, _ := registry.ListServers(ctx)
	snap.ServersCount = len(servers)

	if snap.ProjectsCount == 0 || !snap.TokenValid {
		snap.Status = "degraded"
	} else {
		snap.Status = "healthy"
	}

	snapshotMu.Lock()
	currentSnapshot = snap
	snapshotMu.Unlock()
}

func GetSnapshot() Snapshot {
	snapshotMu.RLock()
	defer snapshotMu.RUnlock()
	return currentSnapshot
}
