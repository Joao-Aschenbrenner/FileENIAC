package cmd

import (
	"fmt"
	"os"

	"github.com/ENIACSystems/FileENIAC/backend/internal/diff"
	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/sync"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"github.com/spf13/cobra"
)

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sincronizar e reconciliar estados",
}

var syncPlanCmd = &cobra.Command{
	Use:   "plan",
	Short: "Analisar divergÃªncias e sugerir sync",
	Run: func(cmd *cobra.Command, args []string) {
		project, _ := cmd.Flags().GetString("project")
		if project == "" && len(args) > 0 {
			project = args[0]
		}
		if project == "" {
			fmt.Println("Uso: eniac sync plan --project <name>")
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

		diffEngine := diff.New()
		report, err := diffEngine.LocalVsMirror(ctx, project)
		if err != nil {
			log.L().Sugar().Fatalf("Diff failed: %v", err)
		}

		syncEngine := sync.New()
		suggestion := syncEngine.Plan(report)

		fmt.Printf("Sync Plan â€” %s\n", project)
		fmt.Printf("AÃ§Ã£o: %s\n", suggestion.Action)
		fmt.Printf("Arquivos: %d\n", suggestion.FileCount)
		fmt.Printf("DescriÃ§Ã£o: %s\n", suggestion.Description)
	},
}

var syncApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Executar sync e gerar manifesto",
	Run: func(cmd *cobra.Command, args []string) {
		project, _ := cmd.Flags().GetString("project")
		if project == "" && len(args) > 0 {
			project = args[0]
		}
		if project == "" {
			fmt.Println("Uso: eniac sync apply --project <name>")
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

		diffEngine := diff.New()
		report, err := diffEngine.LocalVsMirror(ctx, project)
		if err != nil {
			log.L().Sugar().Fatalf("Diff failed: %v", err)
		}

		syncEngine := sync.New()
		suggestion := syncEngine.Plan(report)

		if suggestion.FileCount == 0 {
			fmt.Println("Nada a sincronizar. Estados jÃ¡ estÃ£o consistentes.")
			return
		}

		fmt.Printf("DivergÃªncias: %d arquivos\n", suggestion.FileCount)
		fmt.Print("Deseja continuar? (yes/no): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Sync cancelado pelo usuÃ¡rio.")
			return
		}

		manifest, err := syncEngine.GenerateManifest(ctx, project, "sync_apply", report, "completed")
		if err != nil {
			log.L().Sugar().Fatalf("Manifest failed: %v", err)
		}

		// Mark as synced after successful operation
		syncEngine.Reconcile(ctx, project, "sincronizado")

		fmt.Printf("Sync concluÃ­do. Manifesto: %s\n", manifest.ManifestID)
	},
}

var syncReconcileCmd = &cobra.Command{
	Use:   "reconcile",
	Short: "Marcar projeto como reconciliado",
	Run: func(cmd *cobra.Command, args []string) {
		project, _ := cmd.Flags().GetString("project")
		status, _ := cmd.Flags().GetString("status")
		if project == "" && len(args) > 0 {
			project = args[0]
		}
		if project == "" {
			fmt.Println("Uso: eniac sync reconcile --project <name> --status <status>")
			return
		}
		if status == "" {
			status = "sincronizado"
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

		syncEngine := sync.New()
		if err := syncEngine.Reconcile(ctx, project, status); err != nil {
			log.L().Sugar().Fatalf("Reconcile failed: %v", err)
		}
		fmt.Printf("Projeto %s marcado como: %s\n", project, status)
	},
}

func init() {
	SyncCmd.AddCommand(syncPlanCmd)
	SyncCmd.AddCommand(syncApplyCmd)
	SyncCmd.AddCommand(syncReconcileCmd)

	syncPlanCmd.Flags().StringP("project", "p", "", "Project name")
	syncPlanCmd.Flags().StringP("path", "w", "", "Workspace path (default: current dir)")
	syncApplyCmd.Flags().StringP("project", "p", "", "Project name")
	syncApplyCmd.Flags().StringP("path", "w", "", "Workspace path (default: current dir)")
	syncReconcileCmd.Flags().StringP("project", "p", "", "Project name")
	syncReconcileCmd.Flags().StringP("status", "s", "", "Divergence status")
	syncReconcileCmd.Flags().StringP("path", "w", "", "Workspace path (default: current dir)")
}
