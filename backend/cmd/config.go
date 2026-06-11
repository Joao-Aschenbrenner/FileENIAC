package cmd

import (
	"fmt"
	"os"

	"github.com/eniacsystems/eniac-deploy/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show configuration",
	Long:  `Display the current configuration without sensitive data.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfgPath := getConfigPath()
		loader := config.NewLoader()

		cfg, err := loader.LoadConfig(cfgPath)
		if err != nil {
			fmt.Printf("No configuration found at %s\n", cfgPath)
			fmt.Println("\nRun 'eniac-deploy init' to create a new configuration.")
			os.Exit(0)
		}

		fmt.Printf("Configuration: %s\n\n", cfgPath)
		fmt.Printf("Global:\n")
		fmt.Printf("  Secret Env: %s\n", cfg.Global.SecretEnv)
		fmt.Printf("  History DB: %s\n", cfg.Global.HistoryDB)
		fmt.Printf("  Fallback: %v\n\n", cfg.Global.Fallback)

		fmt.Printf("Projects (%d):\n", len(cfg.Projects))
		for name, p := range cfg.Projects {
			fmt.Printf("\n  [%s]\n", name)
			fmt.Printf("    Working Dir: %s\n", p.WorkingDir)
			fmt.Printf("    Server: %s:%d\n", p.FTPS.Host, p.FTPS.Port)
			fmt.Printf("    User: %s\n", p.FTPS.User)
			fmt.Printf("    Target Path: %s\n", p.Deploy.TargetPath)
			fmt.Printf("    Migrations: %v\n", p.RunMigrations)
		}
	},
}

func init() {
}