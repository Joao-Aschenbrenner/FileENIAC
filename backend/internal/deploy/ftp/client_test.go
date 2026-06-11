package ftp

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	cfg := Config{
		Host:     "ftp.example.com",
		Port:     21,
		User:     "testuser",
		Pass:     "testpass",
		Timeout:  30 * time.Second,
	}

	client := NewClient(cfg)
	
	if client.host != "ftp.example.com" {
		t.Errorf("expected host 'ftp.example.com', got '%s'", client.host)
	}

	if client.port != 21 {
		t.Errorf("expected port 21, got %d", client.port)
	}

	if client.IsConnected() {
		t.Error("new client should not be connected")
	}
}

func TestClient_IsConnected(t *testing.T) {
	client := NewClient(Config{Host: "test"})
	
	if client.IsConnected() {
		t.Error("uninitialized client should not be connected")
	}
}

func TestClient_Config_Defaults(t *testing.T) {
	cfg := Config{
		Host: "ftp.example.com",
		Port: 21,
		User: "test",
		Pass: "test",
	}

	client := NewClient(cfg)
	
	if client == nil {
		t.Error("NewClient should not return nil")
	}

	if client.host != "ftp.example.com" {
		t.Errorf("expected host 'ftp.example.com', got '%s'", client.host)
	}
}