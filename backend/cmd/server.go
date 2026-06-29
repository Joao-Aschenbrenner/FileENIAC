// SPDX-License-Identifier: MIT
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/registry"
	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Gerenciar servidores de deploy",
}

var serverAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adicionar servidor a um projeto",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		name, _ := cmd.Flags().GetString("name")
		srvType, _ := cmd.Flags().GetString("type")
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		targetPath, _ := cmd.Flags().GetString("target-path")
		verifyURL, _ := cmd.Flags().GetString("verify-url")
		wsPath, _ := cmd.Flags().GetString("workspace-path")

		if wsPath == "" {
			pwd, _ := os.Getwd()
			wsPath = pwd
		}

		ctx, err := getWorkspaceContext(wsPath)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}

		proj, err := registry.GetProject(ctx, projectName)
		if err != nil {
			log.L().Sugar().Fatalf("Project not found: %v", err)
		}

		s := &registry.Server{
			ProjectID:  proj.ID,
			Name:       name,
			Type:       srvType,
			Host:       host,
			Port:       port,
			User:       user,
			Password:   password,
			TargetPath: targetPath,
			VerifyURL:  verifyURL,
			IsActive:   true,
		}

		id, err := registry.AddServer(ctx, s)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to add server: %v", err)
		}

		fmt.Printf("Server '%s' added to project '%s' (ID: %d)\n", name, projectName, id)
	},
}

var serverRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remover servidor",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverIDStr := args[0]
		wsPath, _ := cmd.Flags().GetString("workspace-path")

		if wsPath == "" {
			pwd, _ := os.Getwd()
			wsPath = pwd
		}

		serverID, err := strconv.ParseInt(serverIDStr, 10, 64)
		if err != nil {
			log.L().Sugar().Fatalf("Invalid server ID: %v", err)
		}

		ctx, err := getWorkspaceContext(wsPath)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}

		if err := registry.RemoveServer(ctx, serverID); err != nil {
			log.L().Sugar().Fatalf("Failed to remove server: %v", err)
		}

		fmt.Printf("Server %d removed\n", serverID)
	},
}

var serverListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar servidores",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		wsPath, _ := cmd.Flags().GetString("workspace-path")

		if wsPath == "" {
			pwd, _ := os.Getwd()
			wsPath = pwd
		}

		ctx, err := getWorkspaceContext(wsPath)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}

		var servers []*registry.Server
		if projectName != "" {
			proj, err := registry.GetProject(ctx, projectName)
			if err != nil {
				log.L().Sugar().Fatalf("Project not found: %v", err)
			}
			servers, err = registry.ListServersByProject(ctx, proj.ID)
			if err != nil {
				log.L().Sugar().Fatalf("Failed to list servers: %v", err)
			}
		} else {
			servers, err = registry.ListServers(ctx)
			if err != nil {
				log.L().Sugar().Fatalf("Failed to list servers: %v", err)
			}
		}

		if len(servers) == 0 {
			fmt.Println("No servers registered")
			return
		}

		fmt.Printf("Servers (%d):\n", len(servers))
		for _, s := range servers {
			status := "active"
			if !s.IsActive {
				status = "inactive"
			}
			fmt.Printf("  [%d] %s (%s) %s:%d [%s]\n", s.ID, s.Name, s.Type, s.Host, s.Port, status)
		}
	},
}

var serverShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Exibir detalhes do servidor",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverIDStr := args[0]
		wsPath, _ := cmd.Flags().GetString("workspace-path")

		if wsPath == "" {
			pwd, _ := os.Getwd()
			wsPath = pwd
		}

		serverID, err := strconv.ParseInt(serverIDStr, 10, 64)
		if err != nil {
			log.L().Sugar().Fatalf("Invalid server ID: %v", err)
		}

		ctx, err := getWorkspaceContext(wsPath)
		if err != nil {
			log.L().Sugar().Fatalf("Failed to open workspace: %v", err)
		}

		s, err := registry.GetServerByID(ctx, serverID)
		if err != nil {
			log.L().Sugar().Fatalf("Server not found: %v", err)
		}

		fmt.Printf("ID: %d\n", s.ID)
		fmt.Printf("Name: %s\n", s.Name)
		fmt.Printf("Project ID: %d\n", s.ProjectID)
		fmt.Printf("Type: %s\n", s.Type)
		fmt.Printf("Host: %s\n", s.Host)
		fmt.Printf("Port: %d\n", s.Port)
		fmt.Printf("User: %s\n", s.User)
		fmt.Printf("Target Path: %s\n", s.TargetPath)
		fmt.Printf("Verify URL: %s\n", s.VerifyURL)
		fmt.Printf("Active: %v\n", s.IsActive)
	},
}

func init() {
	ServerCmd.AddCommand(serverAddCmd)
	ServerCmd.AddCommand(serverRemoveCmd)
	ServerCmd.AddCommand(serverListCmd)
	ServerCmd.AddCommand(serverShowCmd)

	serverAddCmd.Flags().StringP("project", "p", "", "Project name (required)")
	serverAddCmd.Flags().StringP("name", "n", "", "Server name")
	serverAddCmd.Flags().StringP("type", "t", "ftps", "Server type (ftps, ftp)")
	serverAddCmd.Flags().StringP("host", "H", "", "Server host (required)")
	serverAddCmd.Flags().IntP("port", "P", 21, "Server port")
	serverAddCmd.Flags().StringP("user", "u", "", "Server user")
	serverAddCmd.Flags().StringP("password", "w", "", "Server password")
	serverAddCmd.Flags().StringP("target-path", "d", "/", "Remote target path")
	serverAddCmd.Flags().String("verify-url", "", "Verification URL")
	serverAddCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
	serverAddCmd.MarkFlagRequired("project")
	serverAddCmd.MarkFlagRequired("host")

	serverRemoveCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
	serverListCmd.Flags().StringP("project", "p", "", "Filter by project name")
	serverListCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
	serverShowCmd.Flags().String("workspace-path", "", "Workspace path (default: current dir)")
}
