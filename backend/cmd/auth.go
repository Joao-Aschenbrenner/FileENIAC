package cmd

import (
	"fmt"
	"os"

	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/registry"
	"github.com/spf13/cobra"
)

const gitHubTokenKey = "github_token"

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "AutenticaÃ§Ã£o com GitHub",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Armazenar token de acesso do GitHub",
	Run: func(cmd *cobra.Command, args []string) {
		token, _ := cmd.Flags().GetString("token")
		wsPath, _ := cmd.Flags().GetString("workspace-path")

		if wsPath == "" {
			pwd, _ := os.Getwd()
			wsPath = pwd
		}

		ctx, err := getWorkspaceContext(wsPath)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}

		v, err := registry.VaultFromCtx(ctx)
		if err != nil {
			log.L().Sugar().Fatalf("Vault init: %v", err)
		}

		enc, err := v.Encrypt(token)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to encrypt token: %v", err)
		}

		if err := ctx.DB.SetSetting(gitHubTokenKey, enc); err != nil {
			log.L().Sugar().Fatalf("Failed to store token: %v", err)
		}

		fmt.Println("GitHub token stored successfully")
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Verificar status da autenticaÃ§Ã£o GitHub",
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

		encToken, err := ctx.DB.GetSetting(gitHubTokenKey)
		if err != nil {
			fmt.Println("Not authenticated with GitHub")
			fmt.Println("Run 'fileeniac auth login --token <token>' to authenticate")
			return
		}

		v, err := registry.VaultFromCtx(ctx)
		if err != nil {
			log.L().Sugar().Fatalf("Vault init: %v", err)
		}

		token, err := v.Decrypt(encToken)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to decrypt token: %v", err)
		}

		display := token
		if len(display) > 8 {
			display = display[:4] + "..." + display[len(display)-4:]
		}

		fmt.Printf("Authenticated with GitHub (token: %s)\n", display)
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remover token de acesso do GitHub",
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

		if err := ctx.DB.SetSetting(gitHubTokenKey, ""); err != nil {
			log.L().Sugar().Fatalf("Failed to remove token: %v", err)
		}

		fmt.Println("GitHub token removed")
	},
}

func init() {
	AuthCmd.AddCommand(authLoginCmd)
	AuthCmd.AddCommand(authStatusCmd)
	AuthCmd.AddCommand(authLogoutCmd)

	authLoginCmd.Flags().StringP("token", "t", "", "GitHub personal access token (required)")
	authLoginCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
	authLoginCmd.MarkFlagRequired("token")

	authStatusCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
	authLogoutCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
}
