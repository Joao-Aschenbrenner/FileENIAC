package repair

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/registry"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"go.uber.org/zap"
)

type Report struct {
	OrphanedRepositories int      `json:"orphaned_repositories"`
	BrokenPaths          int      `json:"broken_paths"`
	Fixed                int      `json:"fixed"`
	Warnings             []string `json:"warnings"`
	Errors               []string `json:"errors,omitempty"`
}

func CheckConsistency(ctx *workspace.Context) *Report {
	report := &Report{}

	repos, err := registry.ListAllRepositories(ctx)
	if err != nil {
		report.Errors = append(report.Errors, err.Error())
		return report
	}

	for _, repo := range repos {
		if repo.ProjectID == 0 {
			report.OrphanedRepositories++
			report.Warnings = append(report.Warnings,
				fmt.Sprintf("Repository %s (ID %d) has no project association", repo.FullName, repo.GitHubID))
			continue
		}

		if repo.ClonePath != "" {
			if _, err := os.Stat(repo.ClonePath); os.IsNotExist(err) {
				report.BrokenPaths++
				report.Warnings = append(report.Warnings,
					fmt.Sprintf("Clone path missing for %s: %s", repo.FullName, repo.ClonePath))
			}
		}
	}

	// Check projects with local_path that don't exist
	projects, _ := registry.ListProjects(ctx)
	for _, p := range projects {
		if p.LocalPath != "" {
			if _, err := os.Stat(p.LocalPath); os.IsNotExist(err) {
				report.BrokenPaths++
				report.Warnings = append(report.Warnings,
					fmt.Sprintf("Project %s local path missing: %s", p.Name, p.LocalPath))
			}
		}
	}

	log.L().Info("consistency check complete",
		zap.Int("orphaned", report.OrphanedRepositories),
		zap.Int("broken_paths", report.BrokenPaths),
	)
	return report
}

func RepairOrphanedRepositories(ctx *workspace.Context) (*Report, error) {
	report := CheckConsistency(ctx)

	repos, err := registry.ListAllRepositories(ctx)
	if err != nil {
		return report, err
	}

	for _, repo := range repos {
		if repo.ProjectID == 0 {
			projectName := filepath.Base(repo.FullName)
			existing, err := registry.GetProject(ctx, projectName)
			if err != nil {
				report.Warnings = append(report.Warnings,
					fmt.Sprintf("Cannot auto-fix orphan %s: no matching project found", repo.FullName))
				continue
			}
			registry.UpdateRepositoryImport(ctx, repo.ID, existing.ID, repo.ImportStatus, repo.ClonePath)
			report.Fixed++
			report.Warnings = append(report.Warnings,
				fmt.Sprintf("Re-associated repository %s with project %s", repo.FullName, existing.Name))
		}
	}

	return report, nil
}
