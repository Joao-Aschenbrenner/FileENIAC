// SPDX-License-Identifier: MIT
package cmd

import (
	"fmt"
	"os"

	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace/discovery"
	"github.com/spf13/cobra"
)

var WorkspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Gerenciar workspace",
}

var workspaceInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Inicializar workspace",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		path, _ := cmd.Flags().GetString("path")
		desc, _ := cmd.Flags().GetString("description")

		if path == "" {
			pwd, _ := os.Getwd()
			path = pwd
		}

		ws, err := workspace.Init(name, path, desc)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to init workspace: %v", err)
		}

		fmt.Printf("Workspace '%s' initialized at %s\n", ws.Name, ws.Path)
	},
}

var workspaceOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Abrir workspace existente",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		if path == "" {
			pwd, _ := os.Getwd()
			path = pwd
		}

		ws, err := workspace.Open(path)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}

		fmt.Printf("Workspace '%s' opened from %s\n", ws.Name, ws.Path)
	},
}

var workspaceStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Exibir status do workspace",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		if path == "" {
			pwd, _ := os.Getwd()
			path = pwd
		}

		ws, err := workspace.Open(path)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}

		status := ws.Status()
		if status == nil {
			fmt.Println("No active workspace context")
			return
		}

		fmt.Printf("Workspace: %s\n", status["name"])
		if desc := status["description"]; desc != "" {
			fmt.Printf("Description: %s\n", desc)
		}
		fmt.Printf("Path: %s\n", status["path"])
		fmt.Printf("Projects: %d\n", status["projects"])
		fmt.Printf("Servers: %d\n", status["servers"])
		fmt.Printf("Deploys: %d\n", status["deploys"])
		fmt.Printf("Events: %d\n", status["events"])
	},
}

var workspaceScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Escaneamento de workspaces",
	Long:  "Varre diretÃ³rios em busca de workspaces FileENIAC (.eniac/). A descoberta Ã© somente leitura.",
	Run: func(cmd *cobra.Command, args []string) {
		root, _ := cmd.Flags().GetString("path")
		if root == "" {
			pwd, _ := os.Getwd()
			root = pwd
		}
		depth, _ := cmd.Flags().GetInt("depth")
		if depth <= 0 {
			depth = 3
		}

		scanner := discovery.New()
		results, err := scanner.ScanDeep(root, depth)
		if err != nil {
			log.L().Sugar().Fatalf("Scan failed: %v", err)
		}

		if len(results) == 0 {
			fmt.Println("Nenhum workspace FileENIAC encontrado.")
			return
		}

		fmt.Printf("Workspaces encontrados: %d\n", len(results))
		for _, r := range results {
			if !r.Found {
				continue
			}
			fmt.Printf("  %s (%s)\n", r.Workspace, r.Path)
			fmt.Printf("    Projetos: %d, Servidores: %d\n", len(r.Projects), len(r.Servers))
			for _, p := range r.Projects {
				fmt.Printf("      - %s (env: %s, branch: %s)\n", p.Name, p.Environment, p.Branch)
			}
		}
	},
}

func init() {
	WorkspaceCmd.AddCommand(workspaceInitCmd)
	WorkspaceCmd.AddCommand(workspaceOpenCmd)
	WorkspaceCmd.AddCommand(workspaceStatusCmd)
	WorkspaceCmd.AddCommand(workspaceScanCmd)

	workspaceInitCmd.Flags().StringP("name", "n", "", "Workspace name")
	workspaceInitCmd.Flags().StringP("path", "p", "", "Workspace path (default: current dir)")
	workspaceInitCmd.Flags().StringP("description", "d", "", "Workspace description")
	workspaceInitCmd.MarkFlagRequired("name")

	workspaceOpenCmd.Flags().StringP("path", "p", "", "Workspace path (default: current dir)")
	workspaceStatusCmd.Flags().StringP("path", "p", "", "Workspace path (default: current dir)")
	workspaceScanCmd.Flags().StringP("path", "p", "", "Root directory to scan (default: current dir)")
	workspaceScanCmd.Flags().Int("depth", 3, "Max scan depth")
}
