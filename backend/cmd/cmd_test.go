// SPDX-License-Identifier: MIT
package cmd

import (
	"os"
	"strings"
	"testing"
)

func executeWithArgs(args ...string) (string, error) {
	tmpOut, _ := os.CreateTemp("", "cmdtest-*.out")
	tmpErr, _ := os.CreateTemp("", "cmdtest-*.err")
	defer os.Remove(tmpOut.Name())
	defer os.Remove(tmpErr.Name())

	oldOut := os.Stdout
	oldErr := os.Stderr
	os.Stdout = tmpOut
	os.Stderr = tmpErr

	rootCmd.SetOut(tmpOut)
	rootCmd.SetErr(tmpErr)
	rootCmd.SetArgs(args)

	_, err := rootCmd.ExecuteC()

	tmpOut.Close()
	tmpErr.Close()
	os.Stdout = oldOut
	os.Stderr = oldErr

	outData, _ := os.ReadFile(tmpOut.Name())
	errData, _ := os.ReadFile(tmpErr.Name())
	return string(outData) + string(errData), err
}

func TestRootHelp(t *testing.T) {
	output, err := executeWithArgs("--help")
	if err != nil {
		t.Fatalf("help failed: %v", err)
	}
	if !strings.Contains(output, "fileeniac") {
		t.Error("help should contain 'fileeniac'")
	}
	if !strings.Contains(output, "Usage:") {
		t.Error("help should contain 'Usage:'")
	}
}

func TestRootHelpShortFlag(t *testing.T) {
	output, err := executeWithArgs("-h")
	if err != nil {
		t.Fatalf("-h failed: %v", err)
	}
	if !strings.Contains(output, "fileeniac") {
		t.Error("help should contain 'fileeniac'")
	}
}

func TestInvalidCommand(t *testing.T) {
	_, err := executeWithArgs("nonexistent")
	if err == nil {
		t.Fatal("expected error for invalid command")
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Errorf("expected 'unknown command', got: %v", err)
	}
}

func TestInvalidFlag(t *testing.T) {
	_, err := executeWithArgs("--invalid-flag-xyz")
	if err == nil {
		t.Fatal("expected error for invalid flag")
	}
}

func TestWorkspaceHelp(t *testing.T) {
	output, err := executeWithArgs("workspace", "--help")
	if err != nil {
		t.Fatalf("workspace help failed: %v", err)
	}
	if !strings.Contains(output, "init") {
		t.Error("workspace help should contain 'init' subcommand")
	}
	if !strings.Contains(output, "open") {
		t.Error("workspace help should contain 'open' subcommand")
	}
	if !strings.Contains(output, "status") {
		t.Error("workspace help should contain 'status' subcommand")
	}
	if !strings.Contains(output, "scan") {
		t.Error("workspace help should contain 'scan' subcommand")
	}
}

func TestDeployHelp(t *testing.T) {
	output, err := executeWithArgs("deploy", "--help")
	if err != nil {
		t.Fatalf("deploy help failed: %v", err)
	}
	if !strings.Contains(output, "run") {
		t.Error("deploy help should contain 'run' subcommand")
	}
	if !strings.Contains(output, "verify") {
		t.Error("deploy help should contain 'verify' subcommand")
	}
	if !strings.Contains(output, "rollback") {
		t.Error("deploy help should contain 'rollback' subcommand")
	}
	if !strings.Contains(output, "history") {
		t.Error("deploy help should contain 'history' subcommand")
	}
}

func TestProjectHelp(t *testing.T) {
	output, err := executeWithArgs("project", "--help")
	if err != nil {
		t.Fatalf("project help failed: %v", err)
	}
	if !strings.Contains(output, "add") {
		t.Error("project help should contain 'add' subcommand")
	}
	if !strings.Contains(output, "remove") {
		t.Error("project help should contain 'remove' subcommand")
	}
	if !strings.Contains(output, "list") {
		t.Error("project help should contain 'list' subcommand")
	}
	if !strings.Contains(output, "show") {
		t.Error("project help should contain 'show' subcommand")
	}
}

func TestMirrorHelp(t *testing.T) {
	output, err := executeWithArgs("mirror", "--help")
	if err != nil {
		t.Fatalf("mirror help failed: %v", err)
	}
	if !strings.Contains(output, "create") {
		t.Error("mirror help should contain 'create' subcommand")
	}
	if !strings.Contains(output, "status") {
		t.Error("mirror help should contain 'status' subcommand")
	}
}

func TestSyncHelp(t *testing.T) {
	output, err := executeWithArgs("sync", "--help")
	if err != nil {
		t.Fatalf("sync help failed: %v", err)
	}
	if !strings.Contains(output, "plan") {
		t.Error("sync help should contain 'plan' subcommand")
	}
	if !strings.Contains(output, "apply") {
		t.Error("sync help should contain 'apply' subcommand")
	}
	if !strings.Contains(output, "manifest") {
		t.Error("sync help should contain 'manifest' subcommand")
	}
	if !strings.Contains(output, "reconcile") {
		t.Error("sync help should contain 'reconcile' subcommand")
	}
}

func TestVersionCmd_Help(t *testing.T) {
	output, err := executeWithArgs("version", "--help")
	if err != nil {
		t.Fatalf("version --help failed: %v", err)
	}
	if !strings.Contains(output, "Exibir") {
		t.Errorf("version help should contain description, got: %q", output)
	}
}

func TestRootCmd_HasSubcommands(t *testing.T) {
	expected := []string{
		"version", "workspace", "project", "deploy",
		"diff", "sync", "mirror", "config",
		"server", "repo", "auth", "serve",
		"desktop", "native",
	}

	for _, name := range expected {
		sub, _, err := rootCmd.Find([]string{name})
		if err != nil || sub == nil {
			t.Errorf("expected subcommand %q to exist", name)
		}
	}
}

func TestInvalidSubcommandFormat(t *testing.T) {
	output, err := executeWithArgs("nonexistent")
	if err == nil {
		t.Fatal("expected error for invalid top-level command")
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Errorf("expected 'unknown command', got: %v", err)
	}
	_ = output
}

func TestInvalidSubSubcommand(t *testing.T) {
	output, err := executeWithArgs("workspace", "nonexistent")
	if err == nil {
		// Cobra parent commands without Run print help on invalid subcommand
		_ = output
		return
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Errorf("expected 'unknown command', got: %v", err)
	}
}

func TestWorkspaceInit_MissingRequiredFlag(t *testing.T) {
	_, err := executeWithArgs("workspace", "init")
	if err == nil {
		t.Fatal("expected error when --name is missing (required)")
	}
	if !strings.Contains(err.Error(), "required flag") &&
		!strings.Contains(err.Error(), "--name") {
		t.Errorf("expected 'required flag' or '--name', got: %v", err)
	}
}

func TestRootHelp_ShowsSubcommands(t *testing.T) {
	output, err := executeWithArgs("--help")
	if err != nil {
		t.Fatalf("help failed: %v", err)
	}

	subcommands := []string{"version", "workspace", "project", "deploy", "serve", "native"}
	for _, sc := range subcommands {
		if !strings.Contains(output, sc) {
			t.Errorf("help should mention %q subcommand", sc)
		}
	}
}

func TestShortHelpFlag(t *testing.T) {
	for _, cmd := range []string{"workspace", "project", "deploy", "diff", "sync", "mirror"} {
		output, err := executeWithArgs(cmd, "-h")
		if err != nil {
			t.Fatalf("%s -h failed: %v", cmd, err)
		}
		if !strings.Contains(output, cmd) {
			t.Errorf("help for %q should contain its name", cmd)
		}
	}
}

func TestConfigCmd_HasFlags(t *testing.T) {
	_, _, err := rootCmd.Find([]string{"config", "get"})
	if err != nil {
		t.Fatal("expected 'config get' subcommand")
	}
	_, _, err = rootCmd.Find([]string{"config", "set"})
	if err != nil {
		t.Fatal("expected 'config set' subcommand")
	}
	_, _, err = rootCmd.Find([]string{"config", "list"})
	if err != nil {
		t.Fatal("expected 'config list' subcommand")
	}
}

func TestServerCmd_HasSubcommands(t *testing.T) {
	for _, sc := range []string{"add", "remove", "list", "show"} {
		_, _, err := rootCmd.Find([]string{"server", sc})
		if err != nil {
			t.Errorf("expected 'server %s' subcommand", sc)
		}
	}
}

func TestAuthCmd_HasSubcommands(t *testing.T) {
	for _, sc := range []string{"login", "status", "logout"} {
		_, _, err := rootCmd.Find([]string{"auth", sc})
		if err != nil {
			t.Errorf("expected 'auth %s' subcommand", sc)
		}
	}
}

func TestNoArgs_RootShowsHelp(t *testing.T) {
	output, err := executeWithArgs()
	if err != nil {
		t.Fatalf("root with no args should not error: %v", err)
	}
	if !strings.Contains(output, "Usage:") {
		t.Errorf("root with no args should show help, got: %q", output[:len(output)])
	}
	if !strings.Contains(output, "Available Commands:") {
		t.Error("root help should list available commands")
	}
}

func TestDoubleDash(t *testing.T) {
	output, err := executeWithArgs("--", "version")
	if err != nil {
		t.Fatalf("-- should be ignored: %v", err)
	}
	_ = output
}

func TestNonexistentFlag(t *testing.T) {
	_, err := executeWithArgs("version", "--nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent flag")
	}
}
