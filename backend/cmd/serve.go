// SPDX-License-Identifier: MIT
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
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

		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		addr := fmt.Sprintf("%s:%d", host, port)

		srv := api.New(addr)
		srv.ServeFrontend()

		ctx := workspace.Active()
		if ctx != nil {
			bg := bghealth.NewBackgroundRunner(30 * time.Second)
			bg.Start(ctx)
			srv.SetBackgroundRunner(bg)
			log.L().Info("background health runner started", zap.Duration("interval", 30*time.Second))
		}

		actualAddr, err := srv.ListenDynamic()
		if err != nil {
			log.L().Sugar().Fatalf("Falha ao iniciar servidor: %v", err)
		}

		_, portStr, _ := strings.Cut(actualAddr, ":")
		fmt.Printf("FILEENIAC_READY port=%s token=%s\n", portStr, srv.Token())

		log.L().Sugar().Infof("API server listening on %s", actualAddr)
		fmt.Printf("\n  FileENIAC API rodando em: \x1b[36mhttp://%s\x1b[0m\n\n", actualAddr)

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig
		log.L().Info("shutting down")
		srv.Close()
	},
}

func init() {
	ServeCmd.Flags().String("host", "127.0.0.1", "Server listen host")
	ServeCmd.Flags().IntP("port", "p", 0, "Server listen port (0 = random)")
}
