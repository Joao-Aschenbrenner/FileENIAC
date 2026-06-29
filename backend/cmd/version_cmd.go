// SPDX-License-Identifier: MIT
package cmd

import (
	"fmt"

	"github.com/ENIACSystems/FileENIAC/backend/internal/version"
	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Exibir versão do FileENIAC",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("FileENIAC v%s\n", version.Version)
		fmt.Printf("Build: %s\n", version.BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(VersionCmd)
}
