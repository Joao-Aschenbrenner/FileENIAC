// SPDX-License-Identifier: MIT
package mirror

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ENIACSystems/FileENIAC/backend/internal/registry"
	"github.com/ENIACSystems/FileENIAC/backend/internal/transports"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
)

// mockTransport implements transports.Transport for testing mirror success paths.
type mockTransport struct {
	connectCalled    bool
	disconnectCalled bool
	listCalled       bool
	downloadCalled   bool
	entries          []transports.FileInfo
	connectError     error
	downloadError    error
	listError        error
}

func (m *mockTransport) Connect(ctx context.Context) error {
	m.connectCalled = true
	return m.connectError
}

func (m *mockTransport) Disconnect() error {
	m.disconnectCalled = true
	return nil
}

func (m *mockTransport) List(ctx context.Context, remotePath string) ([]transports.FileInfo, error) {
	m.listCalled = true
	if m.listError != nil {
		return nil, m.listError
	}
	return m.entries, nil
}

func (m *mockTransport) Stat(ctx context.Context, remotePath string) (transports.FileInfo, error) {
	return transports.FileInfo{}, nil
}

func (m *mockTransport) Upload(ctx context.Context, localPath, remotePath string) error {
	return nil
}

func (m *mockTransport) Download(ctx context.Context, remotePath, localPath string) error {
	m.downloadCalled = true
	if m.downloadError != nil {
		return m.downloadError
	}
	return os.WriteFile(localPath, []byte("mock content"), 0644)
}

func (m *mockTransport) Delete(ctx context.Context, remotePath string) error {
	return nil
}

func (m *mockTransport) Mkdir(ctx context.Context, remotePath string) error {
	return nil
}

func (m *mockTransport) Rename(ctx context.Context, from, to string) error {
	return nil
}

func TestMirrorPath(t *testing.T) {
	path := MirrorPath("/home/ws", "MyProject")
	expected := filepath.FromSlash("/home/ws/.eniac/mirror/MyProject")
	if path != expected {
		t.Errorf("expected %s, got %s", expected, path)
	}
}

func TestCreate_SuccessWithMock(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	var err error
	_, err = workspace.Init("MirrorMock", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	projID, err := registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "mock-project",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}

	_, err = registry.AddServer(workspace.Active(), &registry.Server{
		ProjectID:  projID,
		Name:       "mock-server",
		Type:       "ftps",
		Host:       "localhost",
		Port:       21,
		User:       "test",
		Password:   "test",
		TargetPath: "/remote",
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("AddServer failed: %v", err)
	}

	mock := &mockTransport{
		entries: []transports.FileInfo{
			{Name: "file1.txt", IsDir: false, Size: 11},
		},
	}

	origTransportFn := newTransportFn
	defer func() { newTransportFn = origTransportFn }()
	newTransportFn = func(cfg transports.TransportConfig) (transports.Transport, error) {
		return mock, nil
	}

	e := New()
	snapshot, err := e.Create(workspace.Active(), "mock-project")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if snapshot.Status != "completed" {
		t.Errorf("expected completed, got %s", snapshot.Status)
	}
	if snapshot.FilesCount != 1 {
		t.Errorf("expected 1 file downloaded, got %d", snapshot.FilesCount)
	}
	if snapshot.TotalSize <= 0 {
		t.Error("expected positive total size")
	}

	if !mock.connectCalled {
		t.Error("expected Connect to be called")
	}
	if !mock.disconnectCalled {
		t.Error("expected Disconnect to be called")
	}
	if !mock.listCalled {
		t.Error("expected List to be called")
	}
	if !mock.downloadCalled {
		t.Error("expected Download to be called")
	}

	mirrorPath := MirrorPath(wsPath, "mock-project")
	downloadedFile := filepath.Join(mirrorPath, "file1.txt")
	if _, statErr := os.Stat(downloadedFile); os.IsNotExist(statErr) {
		t.Error("expected file1.txt to be downloaded to mirror")
	}

	data, err := os.ReadFile(downloadedFile)
	if err != nil {
		t.Fatalf("failed to read downloaded file: %v", err)
	}
	if string(data) != "mock content" {
		t.Errorf("expected mock content, got %s", string(data))
	}
}

func TestCreate_ConnectFailure(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	_, err := workspace.Init("MirrorFail", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	projID, err := registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "fail-project",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}

	_, err = registry.AddServer(workspace.Active(), &registry.Server{
		ProjectID:  projID,
		Name:       "fail-server",
		Type:       "ftps",
		Host:       "nowhere",
		Port:       21,
		User:       "test",
		Password:   "test",
		TargetPath: "/remote",
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("AddServer failed: %v", err)
	}

	mock := &mockTransport{connectError: fmt.Errorf("mock connect failure")}
	origTransportFn := newTransportFn
	defer func() { newTransportFn = origTransportFn }()
	newTransportFn = func(cfg transports.TransportConfig) (transports.Transport, error) {
		return mock, nil
	}

	e := New()
	_, err = e.Create(workspace.Active(), "fail-project")
	if err == nil {
		t.Fatal("expected error when Connect fails")
	}

	// Mirror dir should have been created despite connection failure
	mirrorPath := MirrorPath(wsPath, "fail-project")
	if _, statErr := os.Stat(mirrorPath); os.IsNotExist(statErr) {
		t.Error("mirror directory should exist even if connect fails")
	}
}

func TestCreate_ListFailure(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	_, err := workspace.Init("MirrorListFail", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	projID, err := registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "listfail-project",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}

	_, err = registry.AddServer(workspace.Active(), &registry.Server{
		ProjectID:  projID,
		Name:       "listfail-server",
		Type:       "ftps",
		Host:       "localhost",
		Port:       21,
		User:       "test",
		Password:   "test",
		TargetPath: "/remote",
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("AddServer failed: %v", err)
	}

	mock := &mockTransport{listError: fmt.Errorf("mock list failure")}
	origTransportFn := newTransportFn
	defer func() { newTransportFn = origTransportFn }()
	newTransportFn = func(cfg transports.TransportConfig) (transports.Transport, error) {
		return mock, nil
	}

	e := New()
	_, err = e.Create(workspace.Active(), "listfail-project")
	if err == nil {
		t.Fatal("expected error when List fails")
	}
}

func TestCreate_DownloadFailure(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	_, err := workspace.Init("MirrorDlFail", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	projID, err := registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "dlfail-project",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}

	_, err = registry.AddServer(workspace.Active(), &registry.Server{
		ProjectID:  projID,
		Name:       "dlfail-server",
		Type:       "ftps",
		Host:       "localhost",
		Port:       21,
		User:       "test",
		Password:   "test",
		TargetPath: "/remote",
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("AddServer failed: %v", err)
	}

	mock := &mockTransport{
		entries: []transports.FileInfo{
			{Name: "bad.txt", IsDir: false, Size: 5},
		},
		downloadError: fmt.Errorf("mock download failure"),
	}
	origTransportFn := newTransportFn
	defer func() { newTransportFn = origTransportFn }()
	newTransportFn = func(cfg transports.TransportConfig) (transports.Transport, error) {
		return mock, nil
	}

	e := New()
	snapshot, err := e.Create(workspace.Active(), "dlfail-project")
	if err != nil {
		t.Fatalf("Create should succeed even if a download fails (file is skipped): %v", err)
	}

	// With retry exhausted, the download fails and the error is logged but not fatal
	if snapshot.Status != "completed" {
		t.Errorf("expected completed, got %s", snapshot.Status)
	}
	if snapshot.FilesCount != 0 {
		t.Errorf("expected 0 files (all downloads failed), got %d", snapshot.FilesCount)
	}
}

func TestStatus_NoSnapshot(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	_, err := workspace.Init("StatusTest", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "no-snapshot",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})

	e := New()
	_, err = e.Status(workspace.Active(), "no-snapshot")
	if err == nil {
		t.Error("expected error for project with no snapshot")
	}
}
