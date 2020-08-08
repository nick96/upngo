package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/nick96/upngo"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List accounts or transactions",
}

var listAccountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "List accounts",
	Run: func(cmd *cobra.Command, args []string) {
		token := getToken()
		client := upngo.NewClient(token)
		accounts, err := client.Accounts()
		if err != nil {
			abort("Failed to get upbank accounts: %v", err)
		}
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		for _, account := range accounts.Data {
			name := account.Attributes.DisplayName
			typ := account.Attributes.AccountType
			amount := account.Attributes.Balance.Format()
			fmt.Fprintf(writer, "%s\t%s\t%s\n", name, typ, amount)
		}
		writer.Flush()
	},
}

var listTransactionsCmd = &cobra.Command{
	Use:   "transactions",
	Short: "List transactions",
	Run: func(cmd *cobra.Command, args []string) {
		token := getToken()
		client := upngo.NewClient(token)
		transactions, err := client.Transactions()
		if err != nil {
			abort("Failed to get upbank transactions: %v", err)
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		for _, transaction := range transactions.Data {
			desc := transaction.Attributes.Description
			msg := transaction.Attributes.Message
			amount := transaction.Attributes.Amount.Format()
			date := transaction.Attributes.CreatedAt.Format(time.RFC1123)
			fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", desc, msg, amount, date)
		}
		writer.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listAccountsCmd, listTransactionsCmd)
}
