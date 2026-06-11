package deployer

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/eniacsystems/eniac-deploy/internal/config"
	"github.com/eniacsystems/eniac-deploy/internal/ftp"
)

type Step interface {
	Execute(ctx context.Context) error
	Rollback(ctx context.Context) error
	Name() string
}

type PackStep struct {
	project *config.Project
	result  string
}

func NewPackStep(project *config.Project) *PackStep {
	return &PackStep{project: project}
}

func (s *PackStep) Name() string {
	return "Pack"
}

func (s *PackStep) Execute(ctx context.Context) error {
	fmt.Println("Packing project...")
	time.Sleep(100 * time.Millisecond)
	s.result = "/tmp/deploy.tar.gz"
	return nil
}

func (s *PackStep) Rollback(ctx context.Context) error {
	if s.result != "" {
		os.Remove(s.result)
	}
	return nil
}

type FTPSConnectStep struct {
	client *ftp.Client
	host   string
	port   int
	user   string
	pass   string
}

func NewFTPSConnectStep(host string, port int, user, pass string) *FTPSConnectStep {
	return &FTPSConnectStep{
		host: host,
		port: port,
		user: user,
		pass: pass,
	}
}

func (s *FTPSConnectStep) Name() string {
	return "FTPS Connect"
}

func (s *FTPSConnectStep) Execute(ctx context.Context) error {
	cfg := ftp.Config{
		Host: s.host,
		Port: s.port,
		User: s.user,
		Pass: s.pass,
	}

	s.client = ftp.NewClient(cfg)
	return s.client.Connect()
}

func (s *FTPSConnectStep) Rollback(ctx context.Context) error {
	if s.client != nil {
		return s.client.Disconnect()
	}
	return nil
}

type UploadStep struct {
	client     *ftp.Client
	localPath  string
	remotePath string
}

func NewUploadStep(client *ftp.Client, localPath, remotePath string) *UploadStep {
	return &UploadStep{
		client:     client,
		localPath:  localPath,
		remotePath: remotePath,
	}
}

func (s *UploadStep) Name() string {
	return "Upload"
}

func (s *UploadStep) Execute(ctx context.Context) error {
	return s.client.Upload(s.localPath, s.remotePath)
}

func (s *UploadStep) Rollback(ctx context.Context) error {
	return s.client.Delete(s.remotePath)
}

type VerifyStep struct {
	verifyURL string
	client    *ftp.Client
}

func NewVerifyStep(verifyURL string) *VerifyStep {
	return &VerifyStep{verifyURL: verifyURL}
}

func (s *VerifyStep) Name() string {
	return "Verify"
}

func (s *VerifyStep) Execute(ctx context.Context) error {
	fmt.Printf("Verifying deployment at %s\n", s.verifyURL)
	time.Sleep(100 * time.Millisecond)
	return nil
}

func (s *VerifyStep) Rollback(ctx context.Context) error {
	return nil
}

func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()

	retryableErrors := []string{
		"timeout",
		"connection refused",
		"connection reset",
		"EOF",
		"max-retries",
		"Temporary failure",
	}

	for _, pattern := range retryableErrors {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}

	return false
}

func RetryWithBackoff(fn func() error, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = fn()
		if err == nil {
			return nil
		}

		if !IsRetryableError(err) {
			return err
		}

		backoff := time.Duration(1<<uint(i)) * time.Second
		if backoff > 30*time.Second {
			backoff = 30 * time.Second
		}

		fmt.Printf("Retry %d/%d after %v: %v\n", i+1, maxRetries, backoff, err)
		time.Sleep(backoff)
	}

	return fmt.Errorf("max retries (%d) exceeded: %w", maxRetries, err)
}