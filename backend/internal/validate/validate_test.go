// SPDX-License-Identifier: MIT
package validate_test

import (
	"os/exec"
	"testing"

	"github.com/ENIACSystems/FileENIAC/backend/internal/validate"
)

func hasGit(t *testing.T) bool {
	t.Helper()
	if err := exec.Command("git", "--version").Run(); err != nil {
		t.Skip("git not available")
		return false
	}
	return true
}

func TestValidateClone_ExistingDir(t *testing.T) {
	dir := t.TempDir()
	result := validate.ValidateClone(dir, "main")
	if result.Valid {
		t.Errorf("expected invalid for existing dir without .git, got checks: %+v", result.Checks)
	}
}

func TestValidateClone_NonExistingDir(t *testing.T) {
	result := validate.ValidateClone("/nonexistent/path/12345", "main")
	if result.Valid {
		t.Error("expected invalid for non-existing dir")
	}
	var hasDirCheck, hasGitCheck bool
	for _, c := range result.Checks {
		if c.Name == "directory_exists" && !c.Passed {
			hasDirCheck = true
		}
		if c.Name == "git_directory" && !c.Passed {
			hasGitCheck = true
		}
	}
	if !hasDirCheck {
		t.Error("expected directory_exists check with Passed=false")
	}
	if !hasGitCheck {
		t.Error("expected git_directory check with Passed=false")
	}
}

func TestValidateClone_ExistingDirWithGit(t *testing.T) {
	if !hasGit(t) {
		return
	}
	dir := t.TempDir()
	cmd := exec.Command("git", "init", dir)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init failed: %s: %v", out, err)
	}

	result := validate.ValidateClone(dir, "main")
	if !result.Valid {
		t.Errorf("expected valid, got checks: %+v", result.Checks)
	}
	var hasGitCheck bool
	for _, c := range result.Checks {
		if c.Name == "git_directory" && c.Passed {
			hasGitCheck = true
		}
	}
	if !hasGitCheck {
		t.Error("expected git_directory check with Passed=true")
	}
}

func TestCloneIntegrity_NonExistingDir(t *testing.T) {
	result := validate.CloneIntegrity("/nonexistent/path/67890", "", "")
	if result.Valid {
		t.Error("expected invalid for non-existing dir")
	}
}

func TestCloneIntegrity_WithGitRepo(t *testing.T) {
	if !hasGit(t) {
		return
	}
	dir := t.TempDir()
	initGitRepo(t, dir)

	result := validate.CloneIntegrity(dir, "", "")
	if !result.Valid {
		t.Errorf("expected valid, got checks: %+v", result.Checks)
	}
}

func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	cmd := exec.Command("git", "init", dir)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init failed: %s: %v", out, err)
	}
	gitCmd := exec.Command("git", "-C", dir, "config", "user.email", "test@test.com")
	if out, err := gitCmd.CombinedOutput(); err != nil {
		t.Fatalf("git config email failed: %s: %v", out, err)
	}
	gitCmd = exec.Command("git", "-C", dir, "config", "user.name", "Test")
	if out, err := gitCmd.CombinedOutput(); err != nil {
		t.Fatalf("git config user failed: %s: %v", out, err)
	}
	gitCmd = exec.Command("git", "-C", dir, "checkout", "-b", "main")
	if out, err := gitCmd.CombinedOutput(); err != nil {
		t.Fatalf("git checkout -b main failed: %s: %v", out, err)
	}
	cmt := exec.Command("git", "-C", dir, "commit", "--allow-empty", "-m", "initial")
	if out, err := cmt.CombinedOutput(); err != nil {
		t.Fatalf("git commit failed: %s: %v", out, err)
	}
	gitCmd = exec.Command("git", "-C", dir, "remote", "add", "origin", "https://github.com/test/test.git")
	if out, err := gitCmd.CombinedOutput(); err != nil {
		t.Fatalf("git remote add failed: %s: %v", out, err)
	}
}

func TestCloneIntegrity_WithBranchCheck(t *testing.T) {
	if !hasGit(t) {
		return
	}
	dir := t.TempDir()
	initGitRepo(t, dir)

	result := validate.CloneIntegrity(dir, "main", "")
	if !result.Valid {
		t.Errorf("expected valid with correct branch 'main', got: %+v", result.Checks)
	}
}

func TestCloneIntegrity_WithWrongBranchCheck(t *testing.T) {
	if !hasGit(t) {
		return
	}
	dir := t.TempDir()
	initGitRepo(t, dir)

	result := validate.CloneIntegrity(dir, "develop", "")
	if result.Valid {
		t.Error("expected invalid with wrong branch 'develop'")
	}
}

func TestValidateImport_EmptyProjectName(t *testing.T) {
	result := validate.ValidateImport("", "/tmp")
	if result.Valid {
		t.Error("expected invalid for empty project name")
	}
	var hasNameCheck bool
	for _, c := range result.Checks {
		if c.Name == "project_name" && !c.Passed {
			hasNameCheck = true
		}
	}
	if !hasNameCheck {
		t.Error("expected project_name check with Passed=false")
	}
}

func TestValidateImport_EmptyLocalPath(t *testing.T) {
	result := validate.ValidateImport("myproject", "")
	if result.Valid {
		t.Error("expected invalid for empty local path")
	}
	var hasPathCheck bool
	for _, c := range result.Checks {
		if c.Name == "local_path" && !c.Passed {
			hasPathCheck = true
		}
	}
	if !hasPathCheck {
		t.Error("expected local_path check with Passed=false")
	}
}

func TestValidateImport_NonExistingPath(t *testing.T) {
	result := validate.ValidateImport("myproject", "/nonexistent/path/99999")
	if result.Valid {
		t.Error("expected invalid for non-existing path")
	}
}

func TestValidateImport_Valid(t *testing.T) {
	dir := t.TempDir()
	result := validate.ValidateImport("myproject", dir)
	if !result.Valid {
		t.Errorf("expected valid, got: %+v", result.Checks)
	}
}

func TestValidateAssociation_ZeroServers(t *testing.T) {
	result := validate.ValidateAssociation(0)
	if result.Valid {
		t.Error("expected invalid for 0 servers")
	}
	var hasCheck bool
	for _, c := range result.Checks {
		if c.Name == "server_association" && !c.Passed {
			hasCheck = true
		}
	}
	if !hasCheck {
		t.Error("expected server_association check with Passed=false")
	}
}

func TestValidateAssociation_WithServers(t *testing.T) {
	result := validate.ValidateAssociation(2)
	if !result.Valid {
		t.Errorf("expected valid for >0 servers, got: %+v", result.Checks)
	}
	var hasCheck bool
	for _, c := range result.Checks {
		if c.Name == "server_association" && c.Passed {
			hasCheck = true
		}
	}
	if !hasCheck {
		t.Error("expected server_association check with Passed=true")
	}
}
