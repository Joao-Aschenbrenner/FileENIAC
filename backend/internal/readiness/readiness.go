package readiness

import (
	"github.com/ENIACSystems/FileENIAC/backend/internal/clone"
	"github.com/ENIACSystems/FileENIAC/backend/internal/registry"
	"github.com/ENIACSystems/FileENIAC/backend/internal/validate"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
)

type Check struct {
	Name    string `json:"name"`
	Passed  bool   `json:"passed"`
	Message string `json:"message,omitempty"`
}

type Result struct {
	Ready  bool    `json:"ready"`
	Checks []Check `json:"checks"`
}

func CheckDeploy(ctx *workspace.Context, projectName string) *Result {
	result := &Result{Ready: true}

	if ctx == nil {
		result.Checks = append(result.Checks, Check{"workspace_loaded", false, "No active workspace"})
		result.Ready = false
		return result
	}
	result.Checks = append(result.Checks, Check{"workspace_loaded", true, ""})

	if projectName == "" {
		result.Checks = append(result.Checks, Check{"project_selected", false, "No project selected"})
		result.Ready = false
		return result
	}

	proj, err := registry.GetProject(ctx, projectName)
	if err != nil {
		result.Checks = append(result.Checks, Check{"project_exists", false, "Project not found: " + err.Error()})
		result.Ready = false
		return result
	}
	result.Checks = append(result.Checks, Check{"project_exists", true, ""})

	if proj.LocalPath != "" {
		v := validate.ValidateClone(proj.LocalPath, proj.Branch)
		if !v.Valid {
			result.Checks = append(result.Checks, Check{"clone_integrity", false, "Clone has issues"})
			result.Ready = false
		} else {
			result.Checks = append(result.Checks, Check{"clone_integrity", true, ""})
		}
	}

	// Check if clone is intact
	if proj.LocalPath != "" {
		if !clone.IsCloned(proj.LocalPath) {
			result.Checks = append(result.Checks, Check{"clone_exists", false, "Local clone not found"})
			result.Ready = false
		} else {
			result.Checks = append(result.Checks, Check{"clone_exists", true, ""})
		}
	}

	servers, _ := registry.ListServersByProject(ctx, proj.ID)
	if len(servers) == 0 {
		result.Checks = append(result.Checks, Check{"server_active", false, "No active server for this project"})
		result.Ready = false
	} else {
		result.Checks = append(result.Checks, Check{"server_active", true, ""})
	}

	return result
}

func CheckSync(ctx *workspace.Context, projectName string) *Result {
	result := CheckDeploy(ctx, projectName)
	// Additional sync-specific checks would go here
	// For now, sync readiness = deploy readiness - server check
	// Remove the server requirement for sync
	checks := make([]Check, 0, len(result.Checks))
	for _, c := range result.Checks {
		if c.Name != "server_active" {
			checks = append(checks, c)
		}
	}
	result.Checks = checks
	// Re-evaluate: ready if all checks pass
	result.Ready = true
	for _, c := range result.Checks {
		if !c.Passed {
			result.Ready = false
			break
		}
	}
	return result
}
