// SPDX-License-Identifier: MIT
package ftp

import "fmt"

func (c *Client) Test() error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected to server")
	}

	if err := c.conn.NoOp(); err != nil {
		return fmt.Errorf("NoOp failed: %w", err)
	}

	return nil
}

func (c *Client) CheckDir(path string) (bool, error) {
	if !c.IsConnected() {
		return false, fmt.Errorf("not connected")
	}

	entries, err := c.conn.List(path)
	if err != nil {
		return false, nil
	}

	return len(entries) >= 0, nil
}

func (c *Client) CheckFile(remotePath string) (bool, error) {
	if !c.IsConnected() {
		return false, fmt.Errorf("not connected")
	}

	entries, err := c.conn.List(remotePath)
	if err != nil {
		return false, nil
	}

	return len(entries) > 0, nil
}

func (c *Client) RawCommand(cmd string) (string, error) {
	if !c.IsConnected() {
		return "", fmt.Errorf("not connected")
	}

	if err := c.conn.NoOp(); err != nil {
		return "", fmt.Errorf("connection not healthy: %w", err)
	}

	return "raw command not supported in FTPS mode", nil
}
