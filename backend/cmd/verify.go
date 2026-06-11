package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify deployment health",
	Long:  `Check if the deployed project is responding correctly.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		url, _ := cmd.Flags().GetString("url")
		timeout, _ := cmd.Flags().GetInt("timeout")

		if url == "" && projectName == "" {
			fmt.Println("Error: either --project or --url is required")
			os.Exit(1)
		}

		if url == "" {
			fmt.Printf("Error: --project flag requires configured verify_url\n")
			os.Exit(1)
		}

		client := &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}

		fmt.Printf("Checking: %s\n", url)

		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("FAILED: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			fmt.Printf("OK: HTTP %d (%.2f KB)\n", resp.StatusCode, float64(len(body))/1024)
		} else {
			fmt.Printf("WARN: HTTP %d\n", resp.StatusCode)
		}
	},
}

func init() {
	verifyCmd.Flags().StringP("project", "p", "", "project name")
	verifyCmd.Flags().StringP("url", "u", "", "URL to check")
	verifyCmd.Flags().IntP("timeout", "t", 30, "request timeout in seconds")
}