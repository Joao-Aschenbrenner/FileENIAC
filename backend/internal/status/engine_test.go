// SPDX-License-Identifier: MIT
package status_test

import (
	"testing"

	"github.com/ENIACSystems/FileENIAC/backend/internal/status"
)

func TestAllStatesDefined(t *testing.T) {
	states := []status.ImportState{
		status.StateNotImported,
		status.StateImporting,
		status.StateImported,
		status.StateClonePending,
		status.StateCloneFailed,
		status.StateRevalidationFailed,
		status.StateNeedsRefresh,
		status.StateCloned,
	}
	if len(states) != 8 {
		t.Errorf("expected 8 states, got %d", len(states))
	}
}

func TestNew(t *testing.T) {
	e := status.New()
	if e == nil {
		t.Fatal("New() returned nil")
	}
}

func TestValidTransitions_ContainsAllStates(t *testing.T) {
	e := status.New()
	trans := e.ValidTransitions()

	expectedKeys := []status.ImportState{
		status.StateNotImported,
		status.StateImporting,
		status.StateClonePending,
		status.StateCloneFailed,
		status.StateCloned,
		status.StateImported,
		status.StateRevalidationFailed,
		status.StateNeedsRefresh,
	}

	for _, k := range expectedKeys {
		if _, ok := trans[k]; !ok {
			t.Errorf("ValidTransitions missing key %q", k)
		}
	}
}

func TestValidTransitions_Specific(t *testing.T) {
	e := status.New()
	trans := e.ValidTransitions()

	tests := []struct {
		from status.ImportState
		to   []status.ImportState
	}{
		{status.StateNotImported, []status.ImportState{status.StateImporting}},
		{status.StateImporting, []status.ImportState{status.StateClonePending, status.StateCloneFailed, status.StateCloned}},
		{status.StateClonePending, []status.ImportState{status.StateCloned, status.StateCloneFailed, status.StateRevalidationFailed}},
		{status.StateCloneFailed, []status.ImportState{status.StateClonePending, status.StateCloned}},
		{status.StateCloned, []status.ImportState{status.StateImported, status.StateRevalidationFailed}},
		{status.StateImported, []status.ImportState{status.StateNeedsRefresh, status.StateRevalidationFailed}},
		{status.StateRevalidationFailed, []status.ImportState{status.StateCloned, status.StateNeedsRefresh, status.StateClonePending}},
		{status.StateNeedsRefresh, []status.ImportState{status.StateImported, status.StateRevalidationFailed}},
	}

	for _, tt := range tests {
		got := trans[tt.from]
		if len(got) != len(tt.to) {
			t.Errorf("ValidTransitions[%q] got %d transitions, want %d", tt.from, len(got), len(tt.to))
			continue
		}
		gotMap := make(map[status.ImportState]bool)
		for _, s := range got {
			gotMap[s] = true
		}
		for _, expected := range tt.to {
			if !gotMap[expected] {
				t.Errorf("ValidTransitions[%q] missing transition to %q", tt.from, expected)
			}
		}
	}
}

func TestResolve_EmptyCurrent(t *testing.T) {
	e := status.New()
	result := e.Resolve("", nil)
	if result != status.StateImporting {
		t.Errorf("expected StateImporting, got %q", result)
	}
}

func TestResolve_NilChecks(t *testing.T) {
	e := status.New()
	result := e.Resolve("", nil)
	if result != status.StateImporting {
		t.Errorf("expected StateImporting, got %q", result)
	}
}

func TestResolve_PendingWithMissingClone(t *testing.T) {
	e := status.New()
	checks := []status.CheckResult{
		{Name: "clone_exists", Passed: false},
	}
	result := e.Resolve("pending", checks)
	if result != status.StateClonePending {
		t.Errorf("expected StateClonePending, got %q", result)
	}
}

func TestResolve_PendingWithAllPassing(t *testing.T) {
	e := status.New()
	checks := []status.CheckResult{
		{Name: "clone_exists", Passed: true},
		{Name: "git_directory", Passed: true},
	}
	result := e.Resolve("pending", checks)
	if result != status.StateImporting {
		t.Errorf("expected StateImporting, got %q", result)
	}
}

func TestResolve_ClonedWithMissingGitDir(t *testing.T) {
	e := status.New()
	checks := []status.CheckResult{
		{Name: "git_directory", Passed: false},
	}
	result := e.Resolve("cloned", checks)
	if result != status.StateRevalidationFailed {
		t.Errorf("expected StateRevalidationFailed, got %q", result)
	}
}

func TestResolve_ClonedAllPassing(t *testing.T) {
	e := status.New()
	checks := []status.CheckResult{
		{Name: "clone_exists", Passed: true},
		{Name: "git_directory", Passed: true},
	}
	result := e.Resolve("cloned", checks)
	if result != status.StateCloned {
		t.Errorf("expected StateCloned, got %q", result)
	}
}

func TestResolve_ImportedWithMissing(t *testing.T) {
	e := status.New()
	checks := []status.CheckResult{
		{Name: "clone_exists", Passed: false},
	}
	result := e.Resolve("imported", checks)
	if result != status.StateClonePending {
		t.Errorf("expected StateClonePending, got %q", result)
	}
}

func TestResolve_ImportedAllPassing(t *testing.T) {
	e := status.New()
	checks := []status.CheckResult{
		{Name: "clone_exists", Passed: true},
	}
	result := e.Resolve("imported", checks)
	if result != status.StateCloned {
		t.Errorf("expected StateCloned, got %q", result)
	}
}

func TestResolve_CloneFailed(t *testing.T) {
	e := status.New()
	result := e.Resolve("clone_failed", nil)
	if result != status.StateCloneFailed {
		t.Errorf("expected StateCloneFailed, got %q", result)
	}
}

func TestResolve_UnknownStateWithFailedCheck(t *testing.T) {
	e := status.New()
	checks := []status.CheckResult{
		{Name: "some_check", Passed: false},
	}
	result := e.Resolve("unknown_state", checks)
	if result != status.ImportState("unknown_state") {
		t.Errorf("expected 'unknown_state', got %q", result)
	}
}

func TestResolve_UnknownStatePreservesCurrent(t *testing.T) {
	e := status.New()
	result := e.Resolve("weird_state", nil)
	if result != status.ImportState("weird_state") {
		t.Errorf("expected 'weird_state', got %q", result)
	}
}
