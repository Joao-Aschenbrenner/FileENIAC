//go:build !windows

package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/ENIACSystems/FileENIAC/backend/internal/api"
	bghealth "github.com/ENIACSystems/FileENIAC/backend/internal/health"
	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var DesktopCmd = &cobra.Command{
	Use:   "desktop",
	Short: "Iniciar FileENIAC (API + Frontend + Navegador)",
	Long:  `Inicia o servidor completo com API, frontend embutido e abre o navegador automaticamente.`,
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("addr")
		noBrowser, _ := cmd.Flags().GetBool("no-browser")

		srv := api.New(addr)
		srv.ServeFrontend()

		ctx := workspace.Active()
		if ctx != nil {
			bg := bghealth.NewBackgroundRunner(30 * time.Second)
			bg.Start(ctx)
			srv.SetBackgroundRunner(bg)
			log.L().Info("background health runner started", zap.Duration("interval", 30*time.Second))
		}

		if !noBrowser {
			time.AfterFunc(500*time.Millisecond, func() {
				url := fmt.Sprintf("http://localhost%s", addr)
				log.L().Info("abrindo navegador", zap.String("url", url))
				openBrowser(url)
			})
		}

		log.L().Sugar().Infof("FileENIAC rodando em http://localhost%s", addr)
		fmt.Printf("\n  FileENIAC rodando em: \x1b[36mhttp://localhost%s\x1b[0m\n\n", addr)

		if err := srv.ListenAndServe(); err != nil {
			log.L().Sugar().Fatalf("Server error: %v", err)
		}
	},
}

func openBrowser(url string) {
	switch runtime.GOOS {
	case "darwin":
		exec.Command("open", url).Start()
	default:
		exec.Command("xdg-open", url).Start()
	}
}

func init() {
	DesktopCmd.Flags().StringP("addr", "a", ":8080", "Server listen address")
	DesktopCmd.Flags().Bool("no-browser", false, "NÃ£o abrir navegador automaticamente")
}
