package cmd

import (
	"fmt"
	"os"

	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/registry"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"github.com/spf13/cobra"
)

var RepoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Gerenciar repositÃ³rios Git dos projetos",
}

var repoAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Associar repositÃ³rio Git a um projeto",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		url, _ := cmd.Flags().GetString("url")
		branch, _ := cmd.Flags().GetString("branch")
		token, _ := cmd.Flags().GetString("token")
		wsPath, _ := cmd.Flags().GetString("workspace-path")

		if wsPath == "" {
			pwd, _ := os.Getwd()
			wsPath = pwd
		}

		ctx, err := getWorkspaceContext(wsPath)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}

		proj, err := registry.GetProject(ctx, projectName)
		if err != nil {
			log.L().Sugar().Fatalf("Project not found: %v", err)
		}

		if token != "" {
			if err := storeGitToken(ctx, url, token); err != nil {
				log.L().Sugar().Fatalf("Failed to store token: %v", err)
			}
		}

		if err := registry.UpdateRepo(ctx, proj.ID, url, branch); err != nil {
			log.L().Sugar().Fatalf("Failed to update repo: %v", err)
		}

		fmt.Printf("Repository '%s' (branch: %s) associated with project '%s'\n", url, branch, projectName)
	},
}

var repoRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remover associaÃ§Ã£o de repositÃ³rio Git",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		wsPath, _ := cmd.Flags().GetString("workspace-path")

		if wsPath == "" {
			pwd, _ := os.Getwd()
			wsPath = pwd
		}

		ctx, err := getWorkspaceContext(wsPath)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}

		proj, err := registry.GetProject(ctx, projectName)
		if err != nil {
			log.L().Sugar().Fatalf("Project not found: %v", err)
		}

		if err := registry.UpdateRepo(ctx, proj.ID, "", "main"); err != nil {
			log.L().Sugar().Fatalf("Failed to remove repo: %v", err)
		}

		fmt.Printf("Repository association removed from project '%s'\n", projectName)
	},
}

var repoListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar repositÃ³rios Git associados",
	Run: func(cmd *cobra.Command, args []string) {
		wsPath, _ := cmd.Flags().GetString("workspace-path")

		if wsPath == "" {
			pwd, _ := os.Getwd()
			wsPath = pwd
		}

		ctx, err := getWorkspaceContext(wsPath)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}

		projects, err := registry.ListProjects(ctx)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to list projects: %v", err)
		}

		count := 0
		for _, p := range projects {
			if p.GitURL != "" {
				count++
				fmt.Printf("  %s -> %s [%s]\n", p.Name, p.GitURL, p.Branch)
			}
		}

		if count == 0 {
			fmt.Println("No repositories registered")
		} else {
			fmt.Printf("\nRepositories (%d):\n", count)
		}
	},
}

var repoShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Exibir detalhes do repositÃ³rio",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		wsPath, _ := cmd.Flags().GetString("workspace-path")

		if wsPath == "" {
			pwd, _ := os.Getwd()
			wsPath = pwd
		}

		ctx, err := getWorkspaceContext(wsPath)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}

		proj, err := registry.GetProject(ctx, projectName)
		if err != nil {
			log.L().Sugar().Fatalf("Project not found: %v", err)
		}

		if proj.GitURL == "" {
			fmt.Printf("Project '%s' has no repository associated\n", projectName)
			return
		}

		fmt.Printf("Project: %s\n", proj.Name)
		fmt.Printf("Git URL: %s\n", proj.GitURL)
		fmt.Printf("Branch: %s\n", proj.Branch)
		fmt.Printf("Local Path: %s\n", proj.LocalPath)
	},
}

func storeGitToken(ctx *workspace.Context, gitURL, token string) error {
	key := "git_token:" + gitURL
	v, err := registry.VaultFromCtx(ctx)
	if err != nil {
		return fmt.Errorf("vault init: %w", err)
	}
	enc, err := v.Encrypt(token)
	if err != nil {
		return fmt.Errorf("encrypt token: %w", err)
	}
	return ctx.DB.SetSetting(key, enc)
}

func init() {
	RepoCmd.AddCommand(repoAddCmd)
	RepoCmd.AddCommand(repoRemoveCmd)
	RepoCmd.AddCommand(repoListCmd)
	RepoCmd.AddCommand(repoShowCmd)

	repoAddCmd.Flags().StringP("project", "p", "", "Project name (required)")
	repoAddCmd.Flags().StringP("url", "u", "", "Git repository URL (required)")
	repoAddCmd.Flags().StringP("branch", "b", "main", "Git branch")
	repoAddCmd.Flags().String("token", "", "GitHub token for authentication")
	repoAddCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
	repoAddCmd.MarkFlagRequired("project")
	repoAddCmd.MarkFlagRequired("url")

	repoRemoveCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
	repoListCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
	repoShowCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
}
