package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

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
	Use:   "transaction",
	Short: "Get transaction by its ID.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getAccountCmd, getTransactionCmd)
}
