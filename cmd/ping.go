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
		fmt.Printf("Successfully pinged UpBank ⚡\n")
	},
}

var pingWebhookCmd = &cobra.Command{
	Use:   "webhook [ID]",
	Short: "Ping UpBank webhook by ID.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token := getToken()
		client := upngo.NewClient(token)
		id := args[0]
		if _, err := client.PingWebhook(id); err != nil {
			abort("Webhook ping failed: %v", err)
		}
		fmt.Printf("Successfully pinged webhook ⚡\n")
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.AddCommand(pingWebhookCmd)
}
