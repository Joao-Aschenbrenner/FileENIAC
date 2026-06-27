package validate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Result struct {
	Valid  bool    `json:"valid"`
	Checks []Check `json:"checks"`
}

type Check struct {
	Name   string `json:"name"`
	Passed bool   `json:"passed"`
	Error  string `json:"error,omitempty"`
}

func ValidateClone(repoDir, branch string) *Result {
	result := &Result{Valid: true}

	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		result.Checks = append(result.Checks, Check{
			Name: "directory_exists", Passed: false, Error: "directory not found",
		})
		result.Valid = false
	} else {
		result.Checks = append(result.Checks, Check{
			Name: "directory_exists", Passed: true,
		})
	}

	gitDir := filepath.Join(repoDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		result.Checks = append(result.Checks, Check{
			Name: "git_directory", Passed: false, Error: ".git not found",
		})
		result.Valid = false
	} else {
		result.Checks = append(result.Checks, Check{
			Name: "git_directory", Passed: true,
		})
	}

	return result
}

func CloneIntegrity(repoDir, expectedBranch, expectedRemote string) *Result {
	result := &Result{Valid: true}

	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		result.Checks = append(result.Checks, Check{
			Name: "directory_exists", Passed: false, Error: "directory not found",
		})
		result.Valid = false
		return result
	}
	result.Checks = append(result.Checks, Check{Name: "directory_exists", Passed: true})

	gitDir := filepath.Join(repoDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		result.Checks = append(result.Checks, Check{
			Name: "git_directory", Passed: false, Error: ".git not found",
		})
		result.Valid = false
		return result
	}
	result.Checks = append(result.Checks, Check{Name: "git_directory", Passed: true})

	currentBranch, err := execOutput(repoDir, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		result.Checks = append(result.Checks, Check{
			Name: "branch_readable", Passed: false, Error: err.Error(),
		})
		result.Valid = false
	} else {
		branchOk := expectedBranch == "" || strings.TrimSpace(currentBranch) == expectedBranch
		result.Checks = append(result.Checks, Check{
			Name: "branch_correct", Passed: branchOk, Error: fmt.Sprintf("expected %s, got %s", expectedBranch, strings.TrimSpace(currentBranch)),
		})
		if !branchOk {
			result.Valid = false
		}
	}

	remoteURL, err := execOutput(repoDir, "git", "config", "--get", "remote.origin.url")
	if err != nil {
		result.Checks = append(result.Checks, Check{
			Name: "remote_origin", Passed: false, Error: "remote.origin not configured",
		})
		result.Valid = false
	} else {
		remoteOk := expectedRemote == "" || strings.TrimSpace(remoteURL) == expectedRemote
		result.Checks = append(result.Checks, Check{
			Name: "remote_correct", Passed: remoteOk, Error: fmt.Sprintf("expected %s", expectedRemote),
		})
		if !remoteOk {
			result.Valid = false
		}
	}

	commitSHA, err := execOutput(repoDir, "git", "rev-parse", "HEAD")
	if err != nil {
		result.Checks = append(result.Checks, Check{
			Name: "commit_readable", Passed: false, Error: err.Error(),
		})
		result.Valid = false
	} else {
		result.Checks = append(result.Checks, Check{
			Name: "commit_readable", Passed: true,
		})
		_ = commitSHA
	}

	gitFsck, err := execOutput(repoDir, "git", "fsck", "--no-progress")
	if err != nil {
		result.Checks = append(result.Checks, Check{
			Name: "git_fsck", Passed: false, Error: "repository integrity check failed",
		})
		result.Valid = false
	} else {
		result.Checks = append(result.Checks, Check{
			Name: "git_fsck", Passed: true,
		})
		_ = gitFsck
	}

	return result
}

func execOutput(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s: %w", name, err)
	}
	return string(out), nil
}

func ValidateImport(projectName, localPath string) *Result {
	result := &Result{Valid: true}

	if projectName == "" {
		result.Checks = append(result.Checks, Check{
			Name: "project_name", Passed: false, Error: "project name is required",
		})
		result.Valid = false
	} else {
		result.Checks = append(result.Checks, Check{
			Name: "project_name", Passed: true,
		})
	}

	if localPath == "" {
		result.Checks = append(result.Checks, Check{
			Name: "local_path", Passed: false, Error: "local path is required",
		})
		result.Valid = false
	} else if _, err := os.Stat(localPath); os.IsNotExist(err) {
		result.Checks = append(result.Checks, Check{
			Name: "local_path", Passed: false, Error: fmt.Sprintf("path does not exist: %s", localPath),
		})
		result.Valid = false
	} else {
		result.Checks = append(result.Checks, Check{
			Name: "local_path", Passed: true,
		})
	}

	return result
}

func ValidateAssociation(serverCount int) *Result {
	result := &Result{Valid: true}
	if serverCount == 0 {
		result.Checks = append(result.Checks, Check{
			Name: "server_association", Passed: false, Error: "no servers associated. Deploy will not work.",
		})
		result.Valid = false
	} else {
		result.Checks = append(result.Checks, Check{
			Name: "server_association", Passed: true,
		})
	}
	return result
}
