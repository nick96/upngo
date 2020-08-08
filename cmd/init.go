package cmd

import (
	"fmt"
	"os"

	"github.com/nick96/upngo/keyring"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise the UpBank CLI for ease of use by adding the token to your keyring.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Enter your UpBank token below and it will be inserted into your keyring for easy of use.")
		fmt.Print("UpBank token: ")
		token, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			abort("Error: failed to read password from stdin: %v", err)
		}

		if err := keyring.SetTokenDefaultconfig(string(token)); err != nil {
			abort("Error: failed to insert UpBank token into keyring: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
