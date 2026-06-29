// SPDX-License-Identifier: MIT
package transports

import (
	"context"
	"time"
)

type TransportConfig struct {
	Protocol string
	Host     string
	Port     int
	User     string
	Pass     string
	Timeout  time.Duration
}

type FileInfo struct {
	Name    string
	Path    string
	Size    int64
	IsDir   bool
	ModTime time.Time
}

type Transport interface {
	Connect(ctx context.Context) error
	Disconnect() error
	Upload(ctx context.Context, localPath, remotePath string) error
	Download(ctx context.Context, remotePath, localPath string) error
	Delete(ctx context.Context, remotePath string) error
	List(ctx context.Context, remotePath string) ([]FileInfo, error)
	Stat(ctx context.Context, remotePath string) (FileInfo, error)
}
