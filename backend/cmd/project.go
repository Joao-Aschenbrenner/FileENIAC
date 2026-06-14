package cmd

import (
	"fmt"
	"os"

	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/registry"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"github.com/spf13/cobra"
)

var ProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Gerenciar projetos",
}

var projectAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adicionar projeto ao workspace",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		path, _ := cmd.Flags().GetString("path")
		remotePath, _ := cmd.Flags().GetString("remote-path")
		wsPath, _ := cmd.Flags().GetString("workspace-path")

		if wsPath == "" {
			pwd, _ := os.Getwd()
			wsPath = pwd
		}

		ctx, err := getWorkspaceContext(wsPath)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}

		p := &registry.Project{
			Name:       name,
			LocalPath:  path,
			RemotePath: remotePath,
			Branch:     "main",
			IsActive:   true,
		}

		id, err := registry.AddProject(ctx, p)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to add project: %v", err)
		}

		fmt.Printf("Project '%s' added (ID: %d)\n", name, id)
	},
}

var projectRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remover projeto do workspace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		wsPath, _ := cmd.Flags().GetString("workspace-path")

		if wsPath == "" {
			pwd, _ := os.Getwd()
			wsPath = pwd
		}

		ctx, err := getWorkspaceContext(wsPath)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}

		project, err := registry.GetProject(ctx, name)
		if err != nil {
			log.L().Sugar().Fatalf("Project not found: %v", err)
		}

		if err := registry.RemoveProject(ctx, project.ID); err != nil {
			log.L().Sugar().Fatalf("Failed to remove project: %v", err)
		}

		fmt.Printf("Project '%s' removed\n", name)
	},
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar projetos do workspace",
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

		if len(projects) == 0 {
			fmt.Println("No projects registered")
			return
		}

		fmt.Printf("Projects (%d):\n", len(projects))
		for _, p := range projects {
			status := "active"
			if !p.IsActive {
				status = "inactive"
			}
			fmt.Printf("  %s (%s) - %s [%s]\n", p.Name, p.Environment, p.LocalPath, status)
		}
	},
}

var projectShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Exibir detalhes do projeto",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		wsPath, _ := cmd.Flags().GetString("workspace-path")

		if wsPath == "" {
			pwd, _ := os.Getwd()
			wsPath = pwd
		}

		ctx, err := getWorkspaceContext(wsPath)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}

		project, err := registry.GetProject(ctx, name)
		if err != nil {
			log.L().Sugar().Fatalf("Project not found: %v", err)
		}

		fmt.Printf("ID: %d\n", project.ID)
		fmt.Printf("Name: %s\n", project.Name)
		fmt.Printf("Local Path: %s\n", project.LocalPath)
		fmt.Printf("Remote Path: %s\n", project.RemotePath)
		fmt.Printf("Branch: %s\n", project.Branch)
		fmt.Printf("Git URL: %s\n", project.GitURL)
		fmt.Printf("Environment: %s\n", project.Environment)
		fmt.Printf("Active: %v\n", project.IsActive)
	},
}

func init() {
	ProjectCmd.AddCommand(projectAddCmd)
	ProjectCmd.AddCommand(projectRemoveCmd)
	ProjectCmd.AddCommand(projectListCmd)
	ProjectCmd.AddCommand(projectShowCmd)

	projectAddCmd.Flags().StringP("name", "n", "", "Project name")
	projectAddCmd.Flags().StringP("path", "p", "", "Local project path")
	projectAddCmd.Flags().StringP("remote-path", "r", "", "Remote server path")
	projectAddCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
	projectAddCmd.MarkFlagRequired("name")
	projectAddCmd.MarkFlagRequired("path")

	projectRemoveCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
	projectListCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
	projectShowCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
}

func getWorkspaceContext(path string) (*workspace.Context, error) {
	_, err := workspace.Open(path)
	if err != nil {
		return nil, err
	}
	ctx := workspace.Active()
	if ctx == nil {
		return nil, fmt.Errorf("failed to get active workspace context after opening")
	}
	return ctx, nil
}
