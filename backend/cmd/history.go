package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/eniacsystems/eniac-deploy/internal/config"
	"github.com/eniacsystems/eniac-deploy/internal/history"
	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show deployment history",
	Long:  `Display the deployment history for a project.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		limit, _ := cmd.Flags().GetInt("limit")

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
		records, err := crud.GetByProject(projectName, limit)
		if err != nil {
			exitWithError("failed to get history", err)
		}

		if len(records) == 0 {
			fmt.Printf("No deployment history for project '%s'\n", projectName)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
		fmt.Fprintf(w, "ID\tProject\tTimestamp\tStatus\tArtifact\t\n")
		fmt.Fprintf(w, "--\t--------\t-----------------\t------\t--------\t\n")

		for _, rec := range records {
			timestamp := rec.Timestamp.Format("2006-01-02 15:04:05")
			artifact := rec.ArtifactHash
			if len(artifact) > 12 {
				artifact = artifact[:12]
			}
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
				rec.ID, rec.ProjectID, timestamp, rec.Status, artifact)
		}
		w.Flush()
	},
}

func init() {
	historyCmd.Flags().StringP("project", "p", "", "project name (required)")
	historyCmd.Flags().IntP("limit", "l", 20, "maximum number of records to show")
	historyCmd.MarkFlagRequired("project")
}