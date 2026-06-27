package cmd

import (
	"fmt"
	"os"

	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/mirror"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"github.com/spf13/cobra"
)

var MirrorCmd = &cobra.Command{
	Use:   "mirror",
	Short: "Gerenciar espelhos locais de servidores",
}

var mirrorCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Criar espelho local do servidor",
	Long:  "Faz download do conteÃºdo remoto via FTPS para .eniac/mirror/{project}/. OperaÃ§Ã£o somente leitura no servidor.",
	Run: func(cmd *cobra.Command, args []string) {
		project, _ := cmd.Flags().GetString("project")
		if project == "" && len(args) > 0 {
			project = args[0]
		}
		if project == "" {
			fmt.Println("Uso: fileeniac mirror create --project <name>")
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

		engine := mirror.New()
		snapshot, err := engine.Create(ctx, project)
		if err != nil {
			log.L().Sugar().Fatalf("Mirror failed: %v", err)
		}

		fmt.Printf("Mirror criado para %s\n", project)
		fmt.Printf("  Snapshot: %s\n", snapshot.ID)
		fmt.Printf("  Arquivos: %d\n", snapshot.FilesCount)
		fmt.Printf("  Tamanho: %d bytes\n", snapshot.TotalSize)
		fmt.Printf("  Status: %s\n", snapshot.Status)
	},
}

var mirrorStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status do Ãºltimo espelho do projeto",
	Run: func(cmd *cobra.Command, args []string) {
		project, _ := cmd.Flags().GetString("project")
		if project == "" && len(args) > 0 {
			project = args[0]
		}
		if project == "" {
			fmt.Println("Uso: fileeniac mirror status --project <name>")
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

		engine := mirror.New()
		snapshot, err := engine.Status(ctx, project)
		if err != nil {
			log.L().Sugar().Fatalf("Mirror status failed: %v", err)
		}

		fmt.Printf("Mirror â€” %s\n", project)
		fmt.Printf("  Snapshot: %s\n", snapshot.ID)
		fmt.Printf("  Arquivos: %d\n", snapshot.FilesCount)
		fmt.Printf("  Tamanho: %d bytes\n", snapshot.TotalSize)
		fmt.Printf("  Status: %s\n", snapshot.Status)
		fmt.Printf("  Iniciado: %s\n", snapshot.StartedAt)
		if snapshot.CompletedAt != "" {
			fmt.Printf("  ConcluÃ­do: %s\n", snapshot.CompletedAt)
		}
	},
}

func init() {
	MirrorCmd.AddCommand(mirrorCreateCmd)
	MirrorCmd.AddCommand(mirrorStatusCmd)

	mirrorCreateCmd.Flags().StringP("project", "p", "", "Project name")
	mirrorCreateCmd.Flags().StringP("path", "w", "", "Workspace path (default: current dir)")
	mirrorStatusCmd.Flags().StringP("project", "p", "", "Project name")
	mirrorStatusCmd.Flags().StringP("path", "w", "", "Workspace path (default: current dir)")
}
