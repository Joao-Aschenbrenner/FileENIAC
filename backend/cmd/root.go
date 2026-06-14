package cmd

import (
	"fmt"
	"os"

	"github.com/ENIACSystems/FileENIAC/backend/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fileeniac",
	Short: "FileENIAC - Gerenciamento de projetos e deploys",
	Long: `FileENIAC é uma plataforma para gerenciamento de workspace local,
deploys FTPS, histórico e monitoramento.

Uso:
  fileeniac native             Iniciar aplicativo desktop nativo (recomendado)
  fileeniac desktop            Iniciar no navegador (fallback)
  fileeniac init               Inicializar workspace
  fileeniac project add        Adicionar projeto
  fileeniac deploy run         Executar deploy
  fileeniac deploy verify      Verificar deploy
  fileeniac deploy rollback    Reverter deploy
  fileeniac deploy history     Histórico de deploys`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("FileENIAC v%s\n", version.Version)
		fmt.Printf("Build: %s\n", version.BuildDate)
		cmd.Help()
	},
}

var versionFlag bool

func init() {
	rootCmd.AddCommand(VersionCmd)
	rootCmd.AddCommand(WorkspaceCmd)
	rootCmd.AddCommand(ProjectCmd)
	rootCmd.AddCommand(DeployCmd)
	rootCmd.AddCommand(DiffCmd)
	rootCmd.AddCommand(SyncCmd)
	rootCmd.AddCommand(MirrorCmd)
	rootCmd.AddCommand(ConfigCmd)
	rootCmd.AddCommand(ServerCmd)
	rootCmd.AddCommand(RepoCmd)
	rootCmd.AddCommand(AuthCmd)
	rootCmd.AddCommand(ServeCmd)
	rootCmd.AddCommand(DesktopCmd)
	rootCmd.AddCommand(NativeCmd)
	rootCmd.AddCommand(UpdateCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(WorkspaceCmd)
	rootCmd.AddCommand(ProjectCmd)
	rootCmd.AddCommand(DeployCmd)
	rootCmd.AddCommand(DiffCmd)
	rootCmd.AddCommand(SyncCmd)
	rootCmd.AddCommand(MirrorCmd)
	rootCmd.AddCommand(ConfigCmd)
	rootCmd.AddCommand(ServerCmd)
	rootCmd.AddCommand(RepoCmd)
	rootCmd.AddCommand(AuthCmd)
	rootCmd.AddCommand(ServeCmd)
	rootCmd.AddCommand(DesktopCmd)
	rootCmd.AddCommand(NativeCmd)
}
