package heartbeat

import (
	"os"
	"sync"
	"time"

	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
)

var (
	mu      sync.Mutex
	t       *time.Timer
	timeout time.Duration
	stopped bool
)

func Start(timeoutDuration time.Duration) {
	mu.Lock()
	defer mu.Unlock()

	if t != nil {
		t.Stop()
	}
	timeout = timeoutDuration
	stopped = false
	t = time.NewTimer(timeout)

	go func() {
		<-t.C
		mu.Lock()
		if stopped {
			mu.Unlock()
			return
		}
		mu.Unlock()
		log.L().Warn("heartbeat timeout - no heartbeat received, shutting down")
		os.Exit(0)
	}()
}

func Reset() {
	mu.Lock()
	defer mu.Unlock()
	if t != nil && !stopped {
		t.Reset(timeout)
	}
}

func Stop() {
	mu.Lock()
	defer mu.Unlock()
	stopped = true
	if t != nil {
		t.Stop()
	}
}
