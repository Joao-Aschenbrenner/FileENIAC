package deploy

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/ENIACSystems/FileENIAC/backend/internal/deploy/packer"
	"github.com/ENIACSystems/FileENIAC/backend/internal/history"
	"github.com/ENIACSystems/FileENIAC/backend/internal/registry"
	"github.com/ENIACSystems/FileENIAC/backend/internal/transports"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
)

// --- Mocks ---

type mockFTPClient struct {
	connectErr  bool
	uploadErr   bool
	downloadErr bool
	deleteErr   bool

	connectCalled    bool
	disconnectCalled bool
	uploadCalled     bool
	uploadLocal      string
	uploadRemote     string
	downloadCalled   bool
	downloadRemote   string
	downloadLocal    string
	deletePaths      []string
	storedContent    []byte
}

func (m *mockFTPClient) Connect(_ context.Context) error {
	m.connectCalled = true
	if m.connectErr {
		return errMock("connect failed")
	}
	return nil
}

func (m *mockFTPClient) Disconnect() error {
	m.disconnectCalled = true
	return nil
}

func (m *mockFTPClient) Upload(_ context.Context, local, remote string) error {
	m.uploadCalled = true
	m.uploadLocal = local
	m.uploadRemote = remote
	if m.uploadErr {
		return errMock("upload failed")
	}
	data, err := os.ReadFile(local)
	if err != nil {
		return err
	}
	m.storedContent = data
	return nil
}

func (m *mockFTPClient) Download(_ context.Context, remote, local string) error {
	m.downloadCalled = true
	m.downloadRemote = remote
	m.downloadLocal = local
	if m.downloadErr {
		return errMock("download failed")
	}
	return os.WriteFile(local, m.storedContent, 0644)
}

func (m *mockFTPClient) Delete(_ context.Context, remote string) error {
	m.deletePaths = append(m.deletePaths, remote)
	if m.deleteErr {
		return errMock("delete failed")
	}
	return nil
}

func (m *mockFTPClient) List(_ context.Context, _ string) ([]transports.FileInfo, error) {
	return nil, nil
}

func (m *mockFTPClient) Stat(_ context.Context, _ string) (transports.FileInfo, error) {
	return transports.FileInfo{}, errMock("not implemented")
}

type mockPacker struct {
	packErr   bool
	fileCount int
}

func (m *mockPacker) Pack(sourceDir, outputPath string) (*packer.Result, error) {
	if m.packErr {
		return nil, errMock("pack failed")
	}
	content := []byte("mock-artifact-content-for-testing")
	if err := os.WriteFile(outputPath, content, 0644); err != nil {
		return nil, err
	}
	return &packer.Result{
		ArchivePath: outputPath,
		FileCount:   m.fileCount,
		SizeBytes:   int64(len(content)),
	}, nil
}

func (m *mockPacker) SetExcludes(excludes []string) {}

func errMock(msg string) error {
	return &mockError{msg: msg}
}

type mockError struct{ msg string }

func (e *mockError) Error() string { return "mock: " + e.msg }

// --- Test helpers ---

func setupTestWorkspace(t *testing.T) (*workspace.Context, *Service) {
	t.Helper()
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projPath := filepath.Join(tmpDir, "project")

	if err := os.MkdirAll(projPath, 0755); err != nil {
		t.Fatalf("mkdir project: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projPath, "index.html"), []byte("<h1>Hello</h1>"), 0644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	if _, err := workspace.Init("TestWS", wsPath, ""); err != nil {
		t.Fatalf("Init: %v", err)
	}
	if _, err := workspace.Open(wsPath); err != nil {
		t.Fatalf("Open: %v", err)
	}
	ctx := workspace.Active()
	if ctx == nil {
		t.Fatal("context is nil")
	}
	t.Cleanup(func() { _ = ctx.DB.Close() })

	projID, err := registry.AddProject(ctx, &registry.Project{
		Name:       "test-project",
		LocalPath:  projPath,
		RemotePath: "/remote",
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("AddProject: %v", err)
	}

	_, err = registry.AddServer(ctx, &registry.Server{
		ProjectID:  projID,
		Name:       "production",
		Type:       "ftps",
		Host:       "ftp.example.com",
		Port:       21,
		User:       "user",
		Password:   "pass",
		TargetPath: "/var/www",
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("AddServer: %v", err)
	}

	s := NewService(ctx.DB)
	return ctx, s
}

func setupProjectOnly(t *testing.T) (*workspace.Context, *Service) {
	t.Helper()
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projPath := filepath.Join(tmpDir, "project")

	if err := os.MkdirAll(projPath, 0755); err != nil {
		t.Fatalf("mkdir project: %v", err)
	}

	if _, err := workspace.Init("TestWS", wsPath, ""); err != nil {
		t.Fatalf("Init: %v", err)
	}
	if _, err := workspace.Open(wsPath); err != nil {
		t.Fatalf("Open: %v", err)
	}
	ctx := workspace.Active()
	if ctx == nil {
		t.Fatal("context is nil")
	}
	t.Cleanup(func() { _ = ctx.DB.Close() })

	_, err := registry.AddProject(ctx, &registry.Project{
		Name:       "no-server-project",
		LocalPath:  projPath,
		RemotePath: "/remote",
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("AddProject: %v", err)
	}

	s := NewService(ctx.DB)
	return ctx, s
}

// --- Validate tests ---

func TestValidate_ExistingProject(t *testing.T) {
	ctx, s := setupTestWorkspace(t)

	err := s.Validate(ctx, "test-project")
	if err != nil {
		t.Errorf("Validate should succeed for existing project with server, got: %v", err)
	}
}

func TestValidate_NonExistingProject(t *testing.T) {
	ctx, s := setupTestWorkspace(t)

	err := s.Validate(ctx, "nonexistent")
	if err == nil {
		t.Fatal("Validate should fail for non-existing project")
	}
}

func TestValidate_NoServer(t *testing.T) {
	ctx, s := setupProjectOnly(t)

	err := s.Validate(ctx, "no-server-project")
	if err == nil {
		t.Fatal("Validate should fail when project has no server configured")
	}
}

// --- Deploy tests ---

func TestDeploy_ProjectNotFound(t *testing.T) {
	ctx, s := setupTestWorkspace(t)

	_, err := s.Deploy(ctx, "nonexistent", false)
	if err == nil {
		t.Fatal("Deploy should fail for non-existing project")
	}
}

func TestDeploy_NoServer(t *testing.T) {
	ctx, s := setupProjectOnly(t)

	_, err := s.Deploy(ctx, "no-server-project", false)
	if err == nil {
		t.Fatal("Deploy should fail when project has no server")
	}
}

func TestDeploy_PackFailure(t *testing.T) {
	ctx, s := setupTestWorkspace(t)

	orig := newPackerFn
	defer func() { newPackerFn = orig }()

	newPackerFn = func(excludes []string) artifactPacker {
		return &mockPacker{packErr: true}
	}

	_, err := s.Deploy(ctx, "test-project", false)
	if err == nil {
		t.Fatal("Deploy should fail when packing fails")
	}
}

func TestDeploy_FTPConnectFailure(t *testing.T) {
	ctx, s := setupTestWorkspace(t)

	origPacker := newPackerFn
	origFTP := newTransportFn
	defer func() {
		newPackerFn = origPacker
		newTransportFn = origFTP
	}()

	newPackerFn = func(excludes []string) artifactPacker {
		return &mockPacker{fileCount: 5}
	}
	newTransportFn = func(cfg transports.TransportConfig) (transports.Transport, error) {
		return &mockFTPClient{connectErr: true}, nil
	}

	_, err := s.Deploy(ctx, "test-project", false)
	if err == nil {
		t.Fatal("Deploy should fail when FTP connection fails")
	}
}

func TestDeploy_FullSuccess(t *testing.T) {
	ctx, s := setupTestWorkspace(t)

	origPacker := newPackerFn
	origFTP := newTransportFn
	defer func() {
		newPackerFn = origPacker
		newTransportFn = origFTP
	}()

	fileCount := 7
	newPackerFn = func(excludes []string) artifactPacker {
		return &mockPacker{fileCount: fileCount}
	}
	newTransportFn = func(cfg transports.TransportConfig) (transports.Transport, error) {
		return &mockFTPClient{}, nil
	}

	result, err := s.Deploy(ctx, "test-project", false)
	if err != nil {
		t.Fatalf("Deploy should succeed: %v", err)
	}
	if result.Status != "success" {
		t.Errorf("expected status 'success', got %q", result.Status)
	}
	if result.DeployID == "" {
		t.Error("expected non-empty deploy ID")
	}
	if result.FilesCount != fileCount {
		t.Errorf("expected FilesCount %d, got %d", fileCount, result.FilesCount)
	}

	// Verify artifact hash is stored correctly (the ArtifactHash bug fix)
	lastDeploy, err := s.history.GetLastDeploy(1)
	if err != nil {
		t.Fatalf("GetLastDeploy: %v", err)
	}
	if lastDeploy.ArtifactHash == "" {
		t.Error("ArtifactHash should not be empty in stored deploy log")
	}
	if lastDeploy.Status != "success" {
		t.Errorf("expected status 'success' in deploy log, got %q", lastDeploy.Status)
	}
}

// --- Rollback tests ---

func TestRollback_ProjectNotFound(t *testing.T) {
	ctx, s := setupTestWorkspace(t)

	_, err := s.Rollback(ctx, "nonexistent")
	if err == nil {
		t.Fatal("Rollback should fail for non-existing project")
	}
}

func TestRollback_NoDeployHistory(t *testing.T) {
	ctx, s := setupTestWorkspace(t)

	_, err := s.Rollback(ctx, "test-project")
	if err == nil {
		t.Fatal("Rollback should fail when no deploy history exists")
	}
}

func TestRollback_FTPServerConnectFailure(t *testing.T) {
	ctx, s := setupTestWorkspace(t)

	origPacker := newPackerFn
	origFTP := newTransportFn
	defer func() {
		newPackerFn = origPacker
		newTransportFn = origFTP
	}()

	// Need a deploy record for rollback to find
	s.history.RecordDeploy(&history.DeployLog{
		ProjectID: 1,
		DeployID:  "dep-rollback-test",
		Status:    "success",
	})

	// FTP connect fails - rollback should still proceed (log-only)
	newTransportFn = func(cfg transports.TransportConfig) (transports.Transport, error) {
		return &mockFTPClient{connectErr: true}, nil
	}

	result, err := s.Rollback(ctx, "test-project")
	if err != nil {
		t.Fatalf("Rollback should proceed even if FTP connect fails: %v", err)
	}
	if result.Status != "rolled_back" {
		t.Errorf("expected status 'rolled_back', got %q", result.Status)
	}
}

func TestRollback_DeletesFromServer(t *testing.T) {
	ctx, s := setupTestWorkspace(t)

	origPacker := newPackerFn
	origFTP := newTransportFn
	defer func() {
		newPackerFn = origPacker
		newTransportFn = origFTP
	}()

	s.history.RecordDeploy(&history.DeployLog{
		ProjectID: 1,
		DeployID:  "dep-rollback-del",
		Status:    "success",
	})

	mockFTP := &mockFTPClient{}
	newTransportFn = func(cfg transports.TransportConfig) (transports.Transport, error) {
		return mockFTP, nil
	}

	result, err := s.Rollback(ctx, "test-project")
	if err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}
	if result.Status != "rolled_back" {
		t.Errorf("expected status 'rolled_back', got %q", result.Status)
	}
	if !mockFTP.connectCalled {
		t.Error("expected FTP Connect to be called")
	}
	if len(mockFTP.deletePaths) < 2 {
		t.Errorf("expected at least 2 Delete calls (artifact + manifest), got %d", len(mockFTP.deletePaths))
	}
	if mockFTP.disconnectCalled != true {
		t.Error("expected FTP Disconnect to be called")
	}
}

func TestRollback_DeleteFailure(t *testing.T) {
	ctx, s := setupTestWorkspace(t)

	origFTP := newTransportFn
	defer func() { newTransportFn = origFTP }()

	s.history.RecordDeploy(&history.DeployLog{
		ProjectID: 1,
		DeployID:  "dep-rollback-del-fail",
		Status:    "success",
	})

	mockFTP := &mockFTPClient{deleteErr: true}
	newTransportFn = func(cfg transports.TransportConfig) (transports.Transport, error) {
		return mockFTP, nil
	}

	result, err := s.Rollback(ctx, "test-project")
	if err != nil {
		t.Fatalf("Rollback should proceed even if FTP delete fails: %v", err)
	}
	if result.Status != "rolled_back" {
		t.Errorf("expected status 'rolled_back', got %q", result.Status)
	}
}

// --- Verify tests ---

func TestVerify_ProjectNotFound(t *testing.T) {
	ctx, s := setupTestWorkspace(t)

	_, err := s.Verify(ctx, "nonexistent")
	if err == nil {
		t.Fatal("Verify should fail for non-existing project")
	}
}

func TestVerify_NoDeploys(t *testing.T) {
	ctx, s := setupTestWorkspace(t)

	result, err := s.Verify(ctx, "test-project")
	if err != nil {
		t.Fatalf("Verify should not error when no deploys exist: %v", err)
	}
	if result.Status != "unknown" {
		t.Errorf("expected status 'unknown', got %q", result.Status)
	}
}

func TestVerify_LastDeployStatus(t *testing.T) {
	ctx, s := setupTestWorkspace(t)

	s.history.RecordDeploy(&history.DeployLog{
		ProjectID: 1,
		DeployID:  "dep-verify-001",
		Status:    "success",
	})

	result, err := s.Verify(ctx, "test-project")
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	if result.DeployID != "dep-verify-001" {
		t.Errorf("expected deploy ID 'dep-verify-001', got %q", result.DeployID)
	}
	if result.Status != "success" {
		t.Errorf("expected status 'success', got %q", result.Status)
	}
}

func TestVerify_FailedDeployStatus(t *testing.T) {
	ctx, s := setupTestWorkspace(t)

	s.history.RecordDeploy(&history.DeployLog{
		ProjectID: 1,
		DeployID:  "dep-verify-fail",
		Status:    "failed",
	})

	result, err := s.Verify(ctx, "test-project")
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	if result.Status != "failed" {
		t.Errorf("expected status 'failed', got %q", result.Status)
	}
}
