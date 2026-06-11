package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "eniac-deploy",
	Short: "ENIAC-DEPLOY v1 - FTP/FTPS deployment tool",
	Long: `ENIAC-DEPLOY is a CLI tool for deploying projects to FTP/FTPS servers.
It supports atomic .env swaps, HMAC token authentication, and rollback functionality.

Usage:
  eniac-deploy init --project myproject
  eniac-deploy push --project myproject
  eniac-deploy rollback --project myproject
  eniac-deploy history --project myproject
`,
	Version: "1.0.0",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "eniac-deploy.toml", "config file path")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(verifyCmd)
	rootCmd.AddCommand(rollbackCmd)
	rootCmd.AddCommand(historyCmd)
	rootCmd.AddCommand(configCmd)
}

func getConfigPath() string {
	if cfgFile != "" {
		return cfgFile
	}
	return "eniac-deploy.toml"
}

func exitWithError(msg string, err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Details: %v\n", err)
	}
	os.Exit(1)
}