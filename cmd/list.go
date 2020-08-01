package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"golang.org/x/text/currency"

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
			unit := currency.MustParseISO(account.Attributes.Balance.CurrencyCode)
			balance, _ := strconv.ParseFloat(account.Attributes.Balance.Value, 64)
			amount := unit.Amount(balance)
			fmtdAmount := currency.NarrowSymbol(amount)


			fmt.Fprintf(writer, "%s\t%s\t%s\n", name, typ, fmtdAmount)
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
		for _, transaction := range transactions {
			fmt.Printf("%v\n", transaction)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listAccountsCmd, listTransactionsCmd)
}
