package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/eniacsystems/eniac-deploy/internal/config"
	"github.com/eniacsystems/eniac-deploy/internal/deployer"
	"github.com/eniacsystems/eniac-deploy/internal/history"
	"github.com/eniacsystems/eniac-deploy/internal/logger"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Deploy project to server",
	Long:  `Package and deploy a project to the configured FTPS server.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		useFallback, _ := cmd.Flags().GetBool("fallback")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if projectName == "" {
			fmt.Println("Error: --project flag is required")
			os.Exit(1)
		}

		cfgPath := getConfigPath()
		loader := config.NewLoader()
		cfg, err := loader.LoadConfig(cfgPath)
		if err != nil {
			exitWithError("failed to load config", err)
		}

		project, err := loader.GetProject(cfg, projectName)
		if err != nil {
			exitWithError("project not found", err)
		}

		secret := os.Getenv(cfg.Global.SecretEnv)
		if secret == "" {
			secret = os.Getenv("ENIAC_DEPLOY_SECRET")
		}
		if secret == "" {
			fmt.Println("Error: ENIAC_DEPLOY_SECRET environment variable is not set")
			os.Exit(1)
		}

		fmt.Printf("Deploying project: %s\n", projectName)
		fmt.Printf("Working dir: %s\n", project.WorkingDir)
		fmt.Printf("Server: %s:%d\n", project.FTPS.Host, project.FTPS.Port)
		fmt.Printf("Fallback mode: %v\n", useFallback)

		if dryRun {
			fmt.Println("\n[DRY RUN] Would deploy with following configuration:")
			fmt.Printf("  Project: %s\n", project.Name)
			fmt.Printf("  Target: %s\n", project.Deploy.TargetPath)
			fmt.Printf("  Migrations: %v\n", project.RunMigrations)
			return
		}

		start := time.Now()

		orch := deployer.NewOrchestrator(project, secret, deployer.Options{
			UseFallback: useFallback,
		})

		if dbPath := getHistoryDBPath(cfg); dbPath != "" {
			db, err := history.NewDB(dbPath)
			if err != nil {
				logger.Warn("Failed to initialize history DB: %v", err)
			} else {
				orch.SetHistoryDB(db)
				defer db.Close()
			}
		}

		if err := orch.Push(context.Background()); err != nil {
			fmt.Fprintf(os.Stderr, "Deploy failed: %v\n", err)
			orch.Cleanup()
			os.Exit(1)
		}

		orch.Disconnect()
		orch.Cleanup()

		fmt.Printf("\nDeploy completed in %v\n", time.Since(start))
	},
}

func init() {
	pushCmd.Flags().StringP("project", "p", "", "project name (required)")
	pushCmd.Flags().Bool("fallback", false, "use FTPS mirror fallback instead of HTTP trigger")
	pushCmd.Flags().Bool("dry-run", false, "show what would be deployed without deploying")
	pushCmd.MarkFlagRequired("project")
}

func getHistoryDBPath(cfg *config.Config) string {
	if cfg.Global.HistoryDB == "" {
		return ""
	}

	path := cfg.Global.HistoryDB
	if path[0] == '~' {
		if home, err := os.UserHomeDir(); err == nil {
			path = filepath.Join(home, path[1:])
		}
	}

	return path
}