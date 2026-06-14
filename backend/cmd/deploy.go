package cmd

import (
	"fmt"
	"os"

	"github.com/ENIACSystems/FileENIAC/backend/internal/deploy"
	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"github.com/spf13/cobra"
)

var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Gerenciar deploys",
}

var deployRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Executar deploy de um projeto",
	Run: func(cmd *cobra.Command, args []string) {
		project, _ := cmd.Flags().GetString("project")
		fallback, _ := cmd.Flags().GetBool("fallback")
		wsPath, _ := cmd.Flags().GetString("workspace-path")

		if wsPath == "" {
			pwd, _ := os.Getwd()
			wsPath = pwd
		}

		ctx, err := getDeployContext(wsPath)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}

		svc := deploy.NewService(ctx.DB)
		result, err := svc.Deploy(ctx, project, fallback)
		if err != nil {
			log.L().Sugar().Fatalf("Deploy failed: %v", err)
		}

		fmt.Printf("Deploy %s: %s\n", result.DeployID, result.Message)
	},
}

var deployVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verificar deploy de um projeto",
	Run: func(cmd *cobra.Command, args []string) {
		project, _ := cmd.Flags().GetString("project")
		wsPath, _ := cmd.Flags().GetString("workspace-path")

		if wsPath == "" {
			pwd, _ := os.Getwd()
			wsPath = pwd
		}

		ctx, err := getDeployContext(wsPath)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}

		svc := deploy.NewService(ctx.DB)
		result, err := svc.Verify(ctx, project)
		if err != nil {
			log.L().Sugar().Fatalf("Verify failed: %v", err)
		}

		fmt.Printf("Status: %s\n", result.Status)
		fmt.Printf("Message: %s\n", result.Message)
	},
}

var deployRollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Reverter deploy de um projeto",
	Run: func(cmd *cobra.Command, args []string) {
		project, _ := cmd.Flags().GetString("project")
		wsPath, _ := cmd.Flags().GetString("workspace-path")

		if wsPath == "" {
			pwd, _ := os.Getwd()
			wsPath = pwd
		}

		ctx, err := getDeployContext(wsPath)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}

		svc := deploy.NewService(ctx.DB)
		result, err := svc.Rollback(ctx, project)
		if err != nil {
			log.L().Sugar().Fatalf("Rollback failed: %v", err)
		}

		fmt.Printf("Rollback: %s\n", result.Message)
	},
}

var deployHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "HistÃ³rico de deploys de um projeto",
	Run: func(cmd *cobra.Command, args []string) {
		project, _ := cmd.Flags().GetString("project")
		limit, _ := cmd.Flags().GetInt("limit")
		wsPath, _ := cmd.Flags().GetString("workspace-path")

		if wsPath == "" {
			pwd, _ := os.Getwd()
			wsPath = pwd
		}

		ctx, err := getDeployContext(wsPath)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}

		svc := deploy.NewService(ctx.DB)
		logs, err := svc.GetHistory(ctx, project, limit)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to get history: %v", err)
		}

		if len(logs) == 0 {
			fmt.Println("No deploys found")
			return
		}

		fmt.Printf("Deploy history for '%s' (%d):\n", project, len(logs))
		for _, l := range logs {
			status := l.Status
			if l.Status == "success" {
				status = "OK"
			} else if l.Status == "failed" {
				status = "FAIL"
			}
			fmt.Printf("  [%s] %s - %s\n", status, l.DeployID, l.CompletedAt)
			if l.CommitMessage != "" {
				fmt.Printf("        %s\n", l.CommitMessage)
			}
		}
	},
}

func init() {
	DeployCmd.AddCommand(deployRunCmd)
	DeployCmd.AddCommand(deployVerifyCmd)
	DeployCmd.AddCommand(deployRollbackCmd)
	DeployCmd.AddCommand(deployHistoryCmd)

	deployRunCmd.Flags().StringP("project", "p", "", "Project name")
	deployRunCmd.Flags().Bool("fallback", false, "Use FTPS mirror fallback")
	deployRunCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
	deployRunCmd.MarkFlagRequired("project")

	deployVerifyCmd.Flags().StringP("project", "p", "", "Project name")
	deployVerifyCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
	deployVerifyCmd.MarkFlagRequired("project")

	deployRollbackCmd.Flags().StringP("project", "p", "", "Project name")
	deployRollbackCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
	deployRollbackCmd.MarkFlagRequired("project")

	deployHistoryCmd.Flags().StringP("project", "p", "", "Project name")
	deployHistoryCmd.Flags().Int("limit", 20, "Number of records")
	deployHistoryCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
	deployHistoryCmd.MarkFlagRequired("project")
}

func getDeployContext(path string) (*workspace.Context, error) {
	_, err := workspace.Open(path)
	if err != nil {
		return nil, err
	}
	_, err = workspace.Open(path)
	if err != nil {
		return nil, err
	}
	return workspace.Active(), nil
}
