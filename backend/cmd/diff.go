package cmd

import (
	"fmt"
	"os"

	"github.com/ENIACSystems/FileENIAC/backend/internal/diff"
	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"github.com/spf13/cobra"
)

var DiffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Comparar estados do workspace",
}

var diffLocalMirrorCmd = &cobra.Command{
	Use:   "local-mirror",
	Short: "Comparar projeto local com mirror",
	Run: func(cmd *cobra.Command, args []string) {
		project, _ := cmd.Flags().GetString("project")
		if project == "" && len(args) > 0 {
			project = args[0]
		}
		if project == "" {
			fmt.Println("Uso: fileeniac diff local-mirror --project <name>")
			return
		}

		path, _ := cmd.Flags().GetString("path")
		if path == "" {
			pwd, _ := os.Getwd()
			path = pwd
		}

		_, err := workspace.Open(path)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}
		ctx := workspace.Active()
		if ctx == nil {
			fmt.Println("No active workspace context")
			return
		}

		engine := diff.New()
		report, err := engine.LocalVsMirror(ctx, project)
		if err != nil {
			log.L().Sugar().Fatalf("Diff failed: %v", err)
		}

		fmt.Printf("Diff Local vs Mirror â€” %s\n", project)
		fmt.Printf("Total: %d | Novo: %d | Modificado: %d | Removido: %d | Sincronizado: %d\n",
			report.Summary.Total, report.Summary.New, report.Summary.Modified, report.Summary.Removed, report.Summary.Synced)

		for _, f := range report.Files {
			if f.Status == "sincronizado" {
				continue
			}
			fmt.Printf("  [%s] %s\n", f.Status, f.Path)
		}
	},
}

var diffStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Exibir status de divergÃªncia do projeto",
	Run: func(cmd *cobra.Command, args []string) {
		project, _ := cmd.Flags().GetString("project")
		if project == "" && len(args) > 0 {
			project = args[0]
		}
		if project == "" {
			fmt.Println("Uso: fileeniac diff status --project <name>")
			return
		}

		path, _ := cmd.Flags().GetString("path")
		if path == "" {
			pwd, _ := os.Getwd()
			path = pwd
		}

		_, err := workspace.Open(path)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}
		ctx := workspace.Active()
		if ctx == nil {
			fmt.Println("No active workspace context")
			return
		}

		engine := diff.New()
		status, err := engine.Status(ctx, project)
		if err != nil {
			log.L().Sugar().Fatalf("Status failed: %v", err)
		}
		fmt.Printf("Projeto %s: %s\n", project, status)
	},
}

func init() {
	DiffCmd.AddCommand(diffLocalMirrorCmd)
	DiffCmd.AddCommand(diffStatusCmd)

	diffLocalMirrorCmd.Flags().StringP("project", "p", "", "Project name")
	diffLocalMirrorCmd.Flags().StringP("path", "w", "", "Workspace path (default: current dir)")
	diffStatusCmd.Flags().StringP("project", "p", "", "Project name")
	diffStatusCmd.Flags().StringP("path", "w", "", "Workspace path (default: current dir)")
}
