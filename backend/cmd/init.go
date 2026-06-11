package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/eniacsystems/eniac-deploy/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project configuration",
	Long:  `Create a new eniac-deploy.toml configuration file for a project.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("name")
		workingDir, _ := cmd.Flags().GetString("dir")
		outputPath, _ := cmd.Flags().GetString("output")

		if projectName == "" {
			projectName = filepath.Base(workingDir)
		}

		if outputPath == "" {
			outputPath = "eniac-deploy.toml"
		}

		cfg := config.DefaultConfig()
		loader := config.NewLoader()

		project := loader.InitProject(projectName, workingDir)
		cfg.Projects[projectName] = project

		if err := loader.SaveConfig(cfg, outputPath); err != nil {
			exitWithError("failed to save config", err)
		}

		fmt.Printf("Created %s\n", outputPath)
		fmt.Println("\nNext steps:")
		fmt.Printf("  1. Edit %s and configure your FTPS credentials\n", outputPath)
		fmt.Printf("  2. Set ENIAC_DEPLOY_SECRET environment variable\n")
		fmt.Printf("  3. Run: eniac-deploy push --project %s\n", projectName)
	},
}

func init() {
	initCmd.Flags().StringP("name", "n", "", "project name")
	initCmd.Flags().StringP("dir", "d", ".", "working directory")
	initCmd.Flags().StringP("output", "o", "eniac-deploy.toml", "output config file")
}