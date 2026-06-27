package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ENIACSystems/FileENIAC/backend/internal/api"
	bghealth "github.com/ENIACSystems/FileENIAC/backend/internal/health"
	"github.com/ENIACSystems/FileENIAC/backend/internal/heartbeat"
	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/update"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var NativeCmd = &cobra.Command{
	Use:   "native",
	Short: "Iniciar FileENIAC como aplicativo desktop nativo",
	Long: `Inicia o servidor backend e abre o aplicativo desktop nativo (WebView2).
Nao abre o navegador - usa uma janela de aplicativo real.`,
	Run: func(cmd *cobra.Command, args []string) {
		if update.CheckAndApply() {
			return
		}

		addr, _ := cmd.Flags().GetString("addr")
		tauriPath, _ := cmd.Flags().GetString("app")

		if tauriPath == "" {
			exe, err := os.Executable()
			if err != nil {
				log.L().Sugar().Fatalf("Nao foi possivel determinar caminho do executavel: %v", err)
			}
			dir := filepath.Dir(exe)
			tauriPath = filepath.Join(dir, "FileENIAC.exe")
			log.L().Info("caminho resolvido", zap.String("exe", exe), zap.String("dir", dir), zap.String("tauriPath", tauriPath))
		} else {
			abs, err := filepath.Abs(tauriPath)
			if err == nil {
				tauriPath = abs
			}
		}

		log.L().Info("verificando aplicativo desktop", zap.String("path", tauriPath))
		if _, err := os.Stat(tauriPath); os.IsNotExist(err) {
			log.L().Sugar().Warnf("Desktop app nao encontrado em %s, tentando busca alternativa", tauriPath)
			alternatives := []string{
				filepath.Join("C:\\Program Files\\FileENIAC", "FileENIAC.exe"),
				filepath.Join("C:\\Program Files (x86)\\FileENIAC", "FileENIAC.exe"),
			}
			tauriPath = ""
			for _, alt := range alternatives {
				if _, err := os.Stat(alt); err == nil {
					tauriPath = alt
					log.L().Info("encontrado em caminho alternativo", zap.String("path", alt))
					break
				}
			}
			if tauriPath == "" {
				log.L().Sugar().Fatalf("Desktop app nao encontrado. Use --app para especificar o caminho do FileENIAC.exe")
			}
		}

		srv := api.New(addr)
		srv.ServeFrontend()

		actualAddr, err := srv.ListenDynamic()
		if err != nil {
			log.L().Sugar().Fatalf("Falha ao iniciar servidor: %v", err)
		}

		_, portStr, _ := strings.Cut(actualAddr, ":")
		os.Setenv("FILEENIAC_API_PORT", portStr)

		ctx := workspace.Active()
		if ctx != nil {
			bg := bghealth.NewBackgroundRunner(30 * time.Second)
			bg.Start(ctx)
			srv.SetBackgroundRunner(bg)
		}

		heartbeat.Start(30 * time.Second)

		time.AfterFunc(800*time.Millisecond, func() {
			log.L().Info("abrindo aplicativo desktop nativo", zap.String("path", tauriPath))
			c := exec.Command(tauriPath)
			c.Env = append(os.Environ(), fmt.Sprintf("FILEENIAC_API_PORT=%s", portStr))
			if err := c.Start(); err != nil {
				log.L().Sugar().Errorf("Falha ao abrir aplicativo desktop: %v", err)
			} else {
				log.L().Info("aplicativo desktop iniciado", zap.Int("pid", c.Process.Pid))
				go func() {
					if err := c.Wait(); err != nil {
						log.L().Sugar().Infof("Aplicativo desktop encerrado: %v", err)
					}
				}()
			}
		})

		portInt, _ := strconv.Atoi(portStr)
		log.L().Sugar().Infof("FileENIAC rodando na porta %d (janela nativa)", portInt)
		fmt.Printf("\n  FileENIAC rodando na porta: \x1b[36m%d\x1b[0m\n\n", portInt)

		select {}
	},
}

func init() {
	NativeCmd.Flags().StringP("addr", "a", ":0", "Server listen address (use :0 for random port)")
	NativeCmd.Flags().String("app", "", "Caminho para FileENIAC.exe (default: mesmo diretorio do fileeniac.exe)")
}
