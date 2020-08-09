package cmd

import (
	"fmt"

	"github.com/nick96/upngo"
	"github.com/spf13/cobra"
)

var (
	webhookDescription string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
}

var addWebhookCmd = &cobra.Command{
	Use:   "webhook [URL]",
	Short: "Register a webhook at URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		token := getToken()
		client := upngo.NewClient(token)
		webhook, err := client.RegisterWebhook(url, upngo.WithDescription(webhookDescription))
		if err != nil {
			abort("Failed to register webhook at %s: %v", url, err)
		}

		fmt.Printf(`
Successfully registered webhook at %s 💸

Here's the secret key:
\t%s
Use it to verify requests send to the webhook URL.
`, url, webhook.Data.Attributes.SecretKey)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.AddCommand(addWebhookCmd)

	addWebhookCmd.Flags().StringVarP(
		&webhookDescription,
		"description",
		"d",
		"",
		"Webhook description (optional)",
	)
}
