package transports

import (
	"context"
	"testing"
)

type mockTransport struct {
	cfg TransportConfig
}

func (m *mockTransport) Connect(ctx context.Context) error { return nil }
func (m *mockTransport) Disconnect() error                  { return nil }
func (m *mockTransport) Upload(ctx context.Context, local, remote string) error {
	return nil
}
func (m *mockTransport) Download(ctx context.Context, remote, local string) error {
	return nil
}
func (m *mockTransport) Delete(ctx context.Context, remote string) error { return nil }
func (m *mockTransport) List(ctx context.Context, remote string) ([]FileInfo, error) {
	return nil, nil
}
func (m *mockTransport) Stat(ctx context.Context, remote string) (FileInfo, error) {
	return FileInfo{}, nil
}

func init() {
	Register("mock", func(cfg TransportConfig) (Transport, error) {
		return &mockTransport{cfg: cfg}, nil
	})
}

func TestRegisterAndLookup(t *testing.T) {
	Register("test_proto", func(cfg TransportConfig) (Transport, error) {
		return &mockTransport{cfg: cfg}, nil
	})

	got, ok := lookup("test_proto")
	if !ok {
		t.Fatal("expected protocol to be registered")
	}
	if got == nil {
		t.Fatal("expected non-nil constructor")
	}
}

func TestRegisterDuplicate(t *testing.T) {
	Register("dup_proto", func(cfg TransportConfig) (Transport, error) {
		return &mockTransport{cfg: cfg}, nil
	})
	Register("dup_proto", func(cfg TransportConfig) (Transport, error) {
		return &mockTransport{cfg: cfg}, nil
	})

	got, ok := lookup("dup_proto")
	if !ok {
		t.Fatal("expected duplicate registration to succeed (overwrite)")
	}
	if got == nil {
		t.Fatal("expected non-nil constructor after duplicate")
	}
}

func TestLookupUnregistered(t *testing.T) {
	_, ok := lookup("nonexistent_protocol")
	if ok {
		t.Error("expected unregistered protocol to return false")
	}
}

func TestNew_ValidProtocol(t *testing.T) {
	tr, err := New(TransportConfig{Protocol: "mock"})
	if err != nil {
		t.Fatalf("New(mock) failed: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil transport")
	}
}

func TestNew_InvalidProtocol(t *testing.T) {
	tr, err := New(TransportConfig{Protocol: "invalid_protocol_xyz"})
	if err == nil {
		t.Fatal("expected error for invalid protocol")
	}
	if tr != nil {
		t.Fatal("expected nil transport on error")
	}
}

func TestRegistered_IncludesMock(t *testing.T) {
	protocols := Registered()
	found := false
	for _, p := range protocols {
		if p == "mock" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'mock' in registered protocols")
	}
}

func TestNew_WithConfig(t *testing.T) {
	Register("config_test", func(cfg TransportConfig) (Transport, error) {
		if cfg.Host != "example.com" || cfg.Port != 2222 {
			t.Errorf("unexpected config: %+v", cfg)
		}
		return &mockTransport{cfg: cfg}, nil
	})

	tr, err := New(TransportConfig{
		Protocol: "config_test",
		Host:     "example.com",
		Port:     2222,
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil transport")
	}
}

func TestTransportImplementsInterface(t *testing.T) {
	var tr Transport = &mockTransport{}
	_ = tr
}
