package cmd

import (
	"fmt"

	"github.com/nick96/upngo"
	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping UpBank. Useful to test your token is correct.",
	Run: func(cmd *cobra.Command, args []string) {
		token := getToken()
		client := upngo.NewClient(token)
		if err := client.Ping(); err != nil {
			abort("UpBank ping failed: %v", err)
		}
		fmt.Printf("Successfully pinged UpBank ⚡")
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
