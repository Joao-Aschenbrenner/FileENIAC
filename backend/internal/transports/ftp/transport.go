// SPDX-License-Identifier: MIT
package ftp

import (
	"context"
	"fmt"
	"path"

	ftplib "github.com/jlaffaye/ftp"

	deployftp "github.com/ENIACSystems/FileENIAC/backend/internal/deploy/ftp"
	"github.com/ENIACSystems/FileENIAC/backend/internal/transports"
)

type Transport struct {
	cfg transports.TransportConfig
	cli *deployftp.Client
}

func init() {
	transports.Register("ftp", func(cfg transports.TransportConfig) (transports.Transport, error) {
		return &Transport{cfg: cfg}, nil
	})
}

func (t *Transport) Connect(ctx context.Context) error {
	ftpCfg := deployftp.Config{
		Host:    t.cfg.Host,
		Port:    t.cfg.Port,
		User:    t.cfg.User,
		Pass:    t.cfg.Pass,
		Timeout: t.cfg.Timeout,
	}
	t.cli = deployftp.NewClient(ftpCfg)
	return t.cli.Connect()
}

func (t *Transport) Disconnect() error {
	if t.cli != nil {
		return t.cli.Disconnect()
	}
	return nil
}

func (t *Transport) Upload(ctx context.Context, localPath, remotePath string) error {
	if t.cli == nil {
		return fmt.Errorf("not connected")
	}
	return t.cli.Upload(localPath, remotePath)
}

func (t *Transport) Download(ctx context.Context, remotePath, localPath string) error {
	if t.cli == nil {
		return fmt.Errorf("not connected")
	}
	return t.cli.Download(remotePath, localPath)
}

func (t *Transport) Delete(ctx context.Context, remotePath string) error {
	if t.cli == nil {
		return fmt.Errorf("not connected")
	}
	return t.cli.Delete(remotePath)
}

func (t *Transport) List(ctx context.Context, remotePath string) ([]transports.FileInfo, error) {
	if t.cli == nil {
		return nil, fmt.Errorf("not connected")
	}
	entries, err := t.cli.List(remotePath)
	if err != nil {
		return nil, err
	}
	result := make([]transports.FileInfo, 0, len(entries))
	for _, e := range entries {
		result = append(result, toFileInfo(e, remotePath))
	}
	return result, nil
}

func (t *Transport) Stat(ctx context.Context, remotePath string) (transports.FileInfo, error) {
	if t.cli == nil {
		return transports.FileInfo{}, fmt.Errorf("not connected")
	}
	parentDir := path.Dir(remotePath)
	entries, err := t.cli.List(parentDir)
	if err != nil {
		return transports.FileInfo{}, err
	}
	targetName := path.Base(remotePath)
	for _, e := range entries {
		if e.Name == targetName {
			return toFileInfo(e, parentDir), nil
		}
	}
	return transports.FileInfo{}, fmt.Errorf("file not found: %s", remotePath)
}

func toFileInfo(e *ftplib.Entry, parentDir string) transports.FileInfo {
	return transports.FileInfo{
		Name:    e.Name,
		Path:    path.Join(parentDir, e.Name),
		Size:    int64(e.Size),
		IsDir:   e.Type == ftplib.EntryTypeFolder,
		ModTime: e.Time,
	}
}
