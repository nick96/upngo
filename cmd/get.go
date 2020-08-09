package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/nick96/upngo"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get accounts or transactions.",
}

var getAccountCmd = &cobra.Command{
	Use:   "account [ID]",
	Short: "Get account by its ID.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		token := getToken()
		client := upngo.NewClient(token)
		account, err := client.Account(id)
		if err != nil {
			abort("Error: failed to get account by ID %s: %v", id, err)
		}

		name := account.Data.Attributes.DisplayName
		typ := account.Data.Attributes.AccountType
		amount := account.Data.Attributes.Balance.Format()

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintf(writer, "Name:\t%s\n", name)
		fmt.Fprintf(writer, "Type:\t%s\n", typ)
		fmt.Fprintf(writer, "Amount:\t%s\n", amount)
		writer.Flush()
	},
}

var getTransactionCmd = &cobra.Command{
	Use:   "transaction [ID]",
	Short: "Get transaction by its ID.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		token := getToken()
		client := upngo.NewClient(token)
		transaction, err := client.Transaction(id)
		if err != nil {
			abort("Error: failed to get transaction by ID %s: %v", id, err)
		}
		desc := transaction.Data.Attributes.Description
		msg := transaction.Data.Attributes.Message
		if msg == "" {
			msg = "N/A"
		}
		amount := transaction.Data.Attributes.Amount.Format()
		date := transaction.Data.Attributes.CreatedAt.Format(time.RFC1123)

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintf(writer, "Description:\t%s\n", desc)
		fmt.Fprintf(writer, "Message:\t%s\n", msg)
		fmt.Fprintf(writer, "Amount:\t%s\n", amount)
		fmt.Fprintf(writer, "Date:\t%s\n", date)
		writer.Flush()
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getAccountCmd, getTransactionCmd)
}
