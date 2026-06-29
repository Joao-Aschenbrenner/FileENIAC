// SPDX-License-Identifier: MIT
package refresh

import (
	"fmt"

	gh "github.com/ENIACSystems/FileENIAC/backend/internal/github"
	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/registry"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"go.uber.org/zap"
)

type RefreshResult struct {
	Organizations int `json:"organizations"`
	Repositories  int `json:"repositories"`
	ChangesFound  int `json:"changes_found"`
	Errors        int `json:"errors"`
}

func RefreshGitHub(ctx *workspace.Context) (*RefreshResult, error) {
	result := &RefreshResult{}

	token, err := getToken(ctx)
	if err != nil {
		return result, fmt.Errorf("github not configured: %w", err)
	}

	svc := gh.New(token)

	orgs, err := svc.ListOrganizations()
	if err != nil {
		return result, fmt.Errorf("list orgs: %w", err)
	}
	result.Organizations = len(orgs)

	for _, org := range orgs {
		repos, err := svc.ListRepositories(org.Login)
		if err != nil {
			log.L().Warn("refresh repos", zap.String("org", org.Login), zap.Error(err))
			result.Errors++
			continue
		}

		imported, _ := registry.ListRepositoriesByOrg(ctx, org.Login)
		importedMap := make(map[int64]*registry.Repository)
		for _, ir := range imported {
			importedMap[ir.GitHubID] = ir
		}

		for _, repo := range repos {
			existing, found := importedMap[repo.ID]
			if !found {
				_, err := registry.AddRepository(ctx, &registry.Repository{
					GitHubID:      repo.ID,
					Name:          repo.Name,
					FullName:      repo.FullName,
					Description:   repo.Description,
					HTMLURL:       repo.HTMLURL,
					CloneURL:      repo.CloneURL,
					DefaultBranch: repo.DefaultBranch,
					Language:      repo.Language,
					Private:       repo.Private,
					Organization:  org.Login,
				})
				if err == nil {
					result.ChangesFound++
				}
			} else {
				if existing.Name != repo.Name || existing.DefaultBranch != repo.DefaultBranch || existing.CloneURL != repo.CloneURL {
					registry.UpdateRepositoryFromGitHub(ctx, existing.ID, repo.Name, repo.FullName, repo.CloneURL, repo.DefaultBranch, repo.Description)
					result.ChangesFound++
				}
			}
		}
		result.Repositories += len(repos)
	}

	log.L().Info("github refresh complete",
		zap.Int("orgs", result.Organizations),
		zap.Int("repos", result.Repositories),
		zap.Int("changes", result.ChangesFound),
	)
	return result, nil
}

func getToken(ctx *workspace.Context) (string, error) {
	enc, err := ctx.DB.GetSetting("github_token")
	if err != nil || enc == "" {
		return "", fmt.Errorf("no github token")
	}
	v, err := registry.VaultFromCtx(ctx)
	if err != nil {
		return "", err
	}
	return v.Decrypt(enc)
}
