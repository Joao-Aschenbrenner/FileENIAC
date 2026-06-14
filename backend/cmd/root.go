package cmd

import (
	"fmt"
	"os"

	"github.com/ENIACSystems/FileENIAC/backend/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "eniac",
	Short: "FileENIAC - Gerenciamento de projetos e deploys",
	Long: `FileENIAC é uma plataforma para gerenciamento de workspace local,
deploys FTPS, histórico e monitoramento.

Uso:
  eniac native             Iniciar aplicativo desktop nativo (recomendado)
  eniac desktop            Iniciar no navegador (fallback)
  FileENIAC init     Inicializar workspace
  eniac project add        Adicionar projeto
  eniac deploy run         Executar deploy
  eniac deploy verify      Verificar deploy
  eniac deploy rollback    Reverter deploy
  eniac deploy history     Histórico de deploys`,
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
