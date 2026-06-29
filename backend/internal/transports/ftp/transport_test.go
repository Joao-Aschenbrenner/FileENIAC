// SPDX-License-Identifier: MIT
package ftp

import (
	"context"
	"testing"

	"github.com/ENIACSystems/FileENIAC/backend/internal/transports"
)

func TestFTPRegistered(t *testing.T) {
	protocols := transports.Registered()
	found := false
	for _, p := range protocols {
		if p == "ftp" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected 'ftp' to be registered via init()")
	}
}

func TestFTPNew(t *testing.T) {
	tr, err := transports.New(transports.TransportConfig{
		Protocol: "ftp",
		Host:     "127.0.0.1",
		Port:     21,
	})
	if err != nil {
		t.Fatalf("New(ftp) failed: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil transport")
	}
}

func TestFTPConnectFails(t *testing.T) {
	tr, err := transports.New(transports.TransportConfig{
		Protocol: "ftp",
		Host:     "127.0.0.1",
		Port:     1,
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	ctx := context.Background()
	err = tr.Connect(ctx)
	if err == nil {
		t.Fatal("expected Connect to fail (no FTP server)")
	}
}

func TestFTPMethodsRequireConnection(t *testing.T) {
	tr, err := transports.New(transports.TransportConfig{
		Protocol: "ftp",
		Host:     "127.0.0.1",
		Port:     21,
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	ctx := context.Background()

	err = tr.Upload(ctx, "/local/file", "/remote/file")
	if err == nil {
		t.Error("expected Upload to fail before Connect")
	}

	err = tr.Download(ctx, "/remote/file", "/local/file")
	if err == nil {
		t.Error("expected Download to fail before Connect")
	}

	err = tr.Delete(ctx, "/remote/file")
	if err == nil {
		t.Error("expected Delete to fail before Connect")
	}

	_, err = tr.List(ctx, "/remote/dir")
	if err == nil {
		t.Error("expected List to fail before Connect")
	}

	_, err = tr.Stat(ctx, "/remote/file")
	if err == nil {
		t.Error("expected Stat to fail before Connect")
	}
}

func TestFTPDisconnectWithoutConnect(t *testing.T) {
	tr := &Transport{}

	err := tr.Disconnect()
	if err != nil {
		t.Errorf("Disconnect without Connect should return nil, got: %v", err)
	}
}

func TestFTPTransportImplementsInterface(t *testing.T) {
	tr := &Transport{}
	var iface transports.Transport = tr
	_ = iface
}
