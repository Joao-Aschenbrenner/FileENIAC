package cmd

import (
	"time"

	"github.com/ENIACSystems/FileENIAC/backend/internal/api"
	bghealth "github.com/ENIACSystems/FileENIAC/backend/internal/health"
	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/update"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Iniciar API HTTP para desktop",
	Long:  `Inicia um servidor HTTP REST para ser consumido pelo Desktop App FileENIAC.`,
	Run: func(cmd *cobra.Command, args []string) {
		if update.CheckAndApply() {
			return
		}

		addr, _ := cmd.Flags().GetString("addr")

		srv := api.New(addr)

		// Start background health runner
		ctx := workspace.Active()
		if ctx != nil {
			bg := bghealth.NewBackgroundRunner(30 * time.Second)
			bg.Start(ctx)
			srv.SetBackgroundRunner(bg)
			log.L().Info("background health runner started", zap.Duration("interval", 30*time.Second))
		}

		log.L().Sugar().Infof("API server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil {
			log.L().Sugar().Fatalf("Server error: %v", err)
		}
	},
}

func init() {
	ServeCmd.Flags().StringP("addr", "a", ":8080", "Server listen address")
}
