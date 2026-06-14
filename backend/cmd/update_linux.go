//go:build !windows

package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var fromPath string

var UpdateCmd = &cobra.Command{
	Use:   "update-from [path]",
	Short: "Atualizar eniac a partir de um novo binÃ¡rio",
	Long: `Substitui o binÃ¡rio eniac atual por um novo binÃ¡rio.
Uso: eniac update-from /caminho/para/novo/eniac

O comando copia o novo binÃ¡rio sobre o atual no diretÃ³rio de instalaÃ§Ã£o.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fromPath = args[0]

		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("nÃ£o foi possÃ­vel determinar o caminho do executÃ¡vel: %w", err)
		}

		from := fromPath
		if from == "" {
			return fmt.Errorf("informe o caminho do novo binÃ¡rio")
		}

		src, err := os.Stat(from)
		if err != nil {
			return fmt.Errorf("arquivo nÃ£o encontrado: %s", from)
		}
		if src.IsDir() {
			return fmt.Errorf("caminho Ã© um diretÃ³rio, informe um arquivo: %s", from)
		}
		if src.Size() == 0 {
			return fmt.Errorf("arquivo vazio: %s", from)
		}

		log.L().Info("atualizando eniac",
			zap.String("de", from),
			zap.String("para", exe),
		)

		backup := exe + ".bak"
		if err := copyFile(exe, backup); err != nil {
			return fmt.Errorf("erro ao criar backup: %w", err)
		}

		if err := copyFile(from, exe); err != nil {
			restoreErr := copyFile(backup, exe)
			if restoreErr != nil {
				log.L().Error("falha ao restaurar backup", zap.Error(restoreErr))
			}
			return fmt.Errorf("erro ao copiar novo binÃ¡rio: %w", err)
		}

		log.L().Info("eniac atualizado (backup salvo em " + backup + ")")

		fmt.Println("eniac atualizado com sucesso!")
		fmt.Println("Backup salvo em:", backup)

		return nil
	},
}

func killProcess(path string) {
	exec.Command("pkill", "-f", filepath.Base(path)).Run()
}

func launchNative(eniacPath, nativeApp string) {
	cmd := exec.Command(eniacPath, "native", "--app", nativeApp)
	cmd.Start()
}

func copyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	if _, err := io.Copy(d, s); err != nil {
		return err
	}
	return d.Sync()
}
