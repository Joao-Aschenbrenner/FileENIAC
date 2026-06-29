// SPDX-License-Identifier: MIT
package update

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/version"
	"go.uber.org/zap"
)

const UpdateDirName = "update"
const BackupDirPrefix = "backup-"

func CheckAndApply() bool {
	exe, err := os.Executable()
	if err != nil {
		return false
	}

	installDir := filepath.Dir(exe)
	updateDir := filepath.Join(installDir, UpdateDirName)

	if _, err := os.Stat(updateDir); os.IsNotExist(err) {
		return false
	}

	updateExe := filepath.Join(updateDir, "fileeniac.exe")
	if _, err := os.Stat(updateExe); os.IsNotExist(err) {
		log.L().Debug("nenhuma atualizacao encontrada em", zap.String("dir", updateDir))
		return false
	}

	log.L().Info("atualizacao detectada", zap.String("dir", updateDir))

	updateTauri := filepath.Join(updateDir, "FileENIAC.exe")
	updateDll := filepath.Join(updateDir, "WebView2Loader.dll")

	if _, err := os.Stat(updateTauri); os.IsNotExist(err) {
		log.L().Warn("FileENIAC.exe nao encontrado no diretorio de update")
		return false
	}

	backupDir := filepath.Join(installDir, BackupDirPrefix+time.Now().Format("20060102-150405"))
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		log.L().Error("falha ao criar diretorio de backup", zap.Error(err))
		return false
	}

	files := []string{"fileeniac.exe", "FileENIAC.exe", "WebView2Loader.dll"}
	for _, f := range files {
		src := filepath.Join(installDir, f)
		if _, err := os.Stat(src); err == nil {
			dst := filepath.Join(backupDir, f)
			if err := copyFile(src, dst); err != nil {
				log.L().Error("falha ao fazer backup", zap.String("file", f), zap.Error(err))
			}
		}
	}

	oldExe := filepath.Join(installDir, "fileeniac.exe.old")
	os.Remove(oldExe)

	if err := os.Rename(exe, oldExe); err != nil {
		log.L().Error("falha ao renomear fileeniac.exe atual", zap.Error(err))
		return false
	}
	log.L().Info("fileeniac.exe renomeado para fileeniac.exe.old")

	updateFiles := map[string]string{
		"FileENIAC.exe":      updateTauri,
		"WebView2Loader.dll": updateDll,
	}

	for name, src := range updateFiles {
		dst := filepath.Join(installDir, name)
		if _, err := os.Stat(src); err == nil {
			if err := copyFile(src, dst); err != nil {
				log.L().Error("falha ao copiar atualizacao", zap.String("file", name), zap.Error(err))
				return false
			}
			log.L().Info("atualizado", zap.String("file", name))
		}
	}

	if err := copyFile(updateExe, exe); err != nil {
		log.L().Error("falha ao copiar novo fileeniac.exe", zap.Error(err))
		return false
	}
	log.L().Info("novo fileeniac.exe copiado")

	log.L().Info("atualizacao concluida. backup salvo em " + backupDir)

	fmt.Printf("\n  \x1b[32mâœ“ FileENIAC atualizado para v%s\x1b[0m\n", version.Version)
	fmt.Printf("  Backup salvo em: %s\n\n", backupDir)

	os.RemoveAll(updateDir)

	startNewInstance(exe, installDir)

	return true
}

func startNewInstance(exePath, installDir string) {
	nativeExe := filepath.Join(installDir, "FileENIAC.exe")
	cmd := exec.Command(exePath, "native")
	if _, err := os.Stat(nativeExe); err == nil {
		cmd = exec.Command(exePath, "native", "--app", nativeExe)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.L().Error("falha ao reiniciar app apos update", zap.Error(err))
		return
	}
	log.L().Info("app reiniciado com nova versao", zap.Int("pid", cmd.Process.Pid))

	os.Exit(0)
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
