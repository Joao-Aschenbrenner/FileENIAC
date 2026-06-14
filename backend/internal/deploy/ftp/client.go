package ftp

import (
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

type Client struct {
	conn      *ftp.ServerConn
	host      string
	user      string
	pass      string
	port      int
	connected bool
}

type Config struct {
	Host    string
	Port    int
	User    string
	Pass    string
	Timeout time.Duration
}

func NewClient(cfg Config) *Client {
	if cfg.Timeout == 0 {
		cfg.Timeout = 120 * time.Second
	}
	return &Client{
		host: cfg.Host,
		port: cfg.Port,
		user: cfg.User,
		pass: cfg.Pass,
	}
}

func (c *Client) Connect() error {
	addr := fmt.Sprintf("%s:%d", c.host, c.port)

	conn, err := ftp.Dial(addr,
		ftp.DialWithTimeout(30*time.Second),
		ftp.DialWithExplicitTLS(&tls.Config{InsecureSkipVerify: false, ServerName: c.host}),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to FTP server: %w", err)
	}

	if err := conn.Login(c.user, c.pass); err != nil {
		conn.Quit()
		return fmt.Errorf("failed to login: %w", err)
	}

	c.conn = conn
	c.connected = true
	return nil
}

func (c *Client) Disconnect() error {
	if c.conn != nil {
		err := c.conn.Quit()
		c.connected = false
		return err
	}
	return nil
}

func (c *Client) IsConnected() bool {
	return c.connected && c.conn != nil
}

func (c *Client) EnsureDir(dirPath string) error {
	parts := strings.Split(dirPath, "/")
	current := ""
	for _, part := range parts {
		if part == "" {
			continue
		}
		current += "/" + part
		err := c.conn.MakeDir(current)
		if err != nil {
			errStr := err.Error()
			if !strings.Contains(errStr, "file exists") && !strings.Contains(errStr, "already exists") {
				return fmt.Errorf("failed to create directory %s: %w", current, err)
			}
		}
	}
	return nil
}

func (c *Client) Upload(localPath, remotePath string) error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected")
	}

	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %w", err)
	}
	defer file.Close()

	remoteDir := path.Dir(remotePath)
	if err := c.EnsureDir(remoteDir); err != nil {
		return err
	}

	err = c.conn.Stor(remotePath, file)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

func (c *Client) Download(remotePath, localPath string) error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected")
	}

	resp, err := c.conn.Retr(remotePath)
	if err != nil {
		return fmt.Errorf("failed to retrieve file: %w", err)
	}
	defer resp.Close()

	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, resp)
	if err != nil {
		return fmt.Errorf("failed to write local file: %w", err)
	}

	return nil
}

func (c *Client) Delete(remotePath string) error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected")
	}

	err := c.conn.Delete(remotePath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (c *Client) Rename(oldPath, newPath string) error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected")
	}

	err := c.conn.Rename(oldPath, newPath)
	if err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}

	return nil
}

func (c *Client) List(path string) ([]*ftp.Entry, error) {
	if !c.IsConnected() {
		return nil, fmt.Errorf("not connected")
	}

	entries, err := c.conn.List(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list directory: %w", err)
	}

	return entries, nil
}

func (c *Client) GetCurrentDir() (string, error) {
	if !c.IsConnected() {
		return "", fmt.Errorf("not connected")
	}

	pwd, err := c.conn.CurrentDir()
	if err != nil {
		return "", fmt.Errorf("failed to get current dir: %w", err)
	}

	return pwd, nil
}

func (c *Client) ChangeDir(path string) error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected")
	}

	err := c.conn.ChangeDir(path)
	if err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}

	return nil
}

func (c *Client) NoOp() error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected")
	}

	return c.conn.NoOp()
}
