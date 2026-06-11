package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/eniacsystems/eniac-deploy/internal/config"
	"github.com/eniacsystems/eniac-deploy/internal/deployer"
	"github.com/eniacsystems/eniac-deploy/internal/ftp"
	"github.com/eniacsystems/eniac-deploy/internal/history"
	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback to previous deployment",
	Long:  `Restore the .env.bak file from the last successful deployment.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
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

		dbPath := getHistoryDBPath(cfg)
		if dbPath == "" {
			fmt.Println("Error: history DB not configured")
			os.Exit(1)
		}

		db, err := history.NewDB(dbPath)
		if err != nil {
			exitWithError("failed to open history DB", err)
		}
		defer db.Close()

		crud := history.NewCRUD(db)

		lastSuccess, err := crud.GetLastSuccessful(projectName)
		if err != nil {
			exitWithError("failed to get last successful deployment", err)
		}

		if lastSuccess == nil {
			fmt.Println("Error: no successful deployment found for rollback")
			os.Exit(1)
		}

		fmt.Printf("Last successful deployment: #%d at %s\n", lastSuccess.ID, lastSuccess.Timestamp)

		if dryRun {
			fmt.Println("\n[DRY RUN] Would rollback to deployment #" + fmt.Sprintf("%d", lastSuccess.ID))
			return
		}

		ftpClient := ftp.NewClient(ftp.Config{
			Host: project.FTPS.Host,
			Port: project.FTPS.Port,
			User: project.FTPS.User,
			Pass: project.FTPS.Pass,
		})

		if err := ftpClient.Connect(); err != nil {
			exitWithError("failed to connect to FTPS", err)
		}
		defer ftpClient.Disconnect()

		rollbacker := deployer.NewRollbacker(project, ftpClient, crud)
		if err := rollbacker.Execute(context.Background()); err != nil {
			exitWithError("rollback failed", err)
		}

		fmt.Println("Rollback completed successfully")
	},
}

func init() {
	rollbackCmd.Flags().StringP("project", "p", "", "project name (required)")
	rollbackCmd.Flags().Bool("dry-run", false, "show what would be rolled back without making changes")
	rollbackCmd.MarkFlagRequired("project")
}