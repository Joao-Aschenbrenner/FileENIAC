package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Gerenciar configuraÃ§Ãµes do workspace",
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Obter valor de configuraÃ§Ã£o",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		wsPath, _ := cmd.Flags().GetString("workspace-path")
		if wsPath == "" {
			pwd, _ := os.Getwd()
			wsPath = pwd
		}

		ctx, err := getWorkspaceContext(wsPath)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}

		if key == "name" {
			fmt.Println(ctx.Config.Workspace.Name)
			return
		}
		if key == "description" {
			fmt.Println(ctx.Config.Workspace.Description)
			return
		}
		if strings.HasPrefix(key, "workspace.") {
			log.L().Sugar().Fatalf("Use 'fileeniac config set' to modify workspace.%s", key)
		}

		value, err := ctx.DB.GetSetting(key)
		if err != nil {
			log.L().Sugar().Fatalf("Setting '%s' not found: %v", key, err)
		}
		fmt.Println(value)
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Definir valor de configuraÃ§Ã£o",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key, value := args[0], args[1]
		wsPath, _ := cmd.Flags().GetString("workspace-path")
		if wsPath == "" {
			pwd, _ := os.Getwd()
			wsPath = pwd
		}

		ctx, err := getWorkspaceContext(wsPath)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}

		if key == "name" || key == "description" {
			if key == "name" {
				ctx.Config.Workspace.Name = value
			} else {
				ctx.Config.Workspace.Description = value
			}
			wsDir := wsPath + "/.eniac"
			if err := saveConfigToDir(ctx.Config, wsDir); err != nil {
				log.L().Sugar().Fatalf("Failed to save config: %v", err)
			}
			fmt.Printf("%s updated to '%s'\n", key, value)
			return
		}

		if err := ctx.DB.SetSetting(key, value); err != nil {
			log.L().Sugar().Fatalf("Failed to set '%s': %v", key, err)
		}
		fmt.Printf("%s set to '%s'\n", key, value)
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar configuraÃ§Ãµes do workspace",
	Run: func(cmd *cobra.Command, args []string) {
		wsPath, _ := cmd.Flags().GetString("workspace-path")
		if wsPath == "" {
			pwd, _ := os.Getwd()
			wsPath = pwd
		}

		ctx, err := getWorkspaceContext(wsPath)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}

		fmt.Printf("workspace.name = %s\n", ctx.Config.Workspace.Name)
		fmt.Printf("workspace.description = %s\n", ctx.Config.Workspace.Description)

		settings, err := ctx.DB.ListSettings()
		if err != nil {
			log.L().Sugar().Fatalf("Failed to list settings: %v", err)
		}
		for k, v := range settings {
			if k == "schema_version" {
				continue
			}
			fmt.Printf("%s = %s\n", k, v)
		}
	},
}

func saveConfigToDir(cfg *workspace.Config, wsDir string) error {
	return workspace.SaveConfig(cfg, wsDir)
}

func init() {
	ConfigCmd.AddCommand(configGetCmd)
	ConfigCmd.AddCommand(configSetCmd)
	ConfigCmd.AddCommand(configListCmd)

	configGetCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
	configSetCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
	configListCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
}
