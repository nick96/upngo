package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var verbose *bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "upngo",
	Short: "Talk to your bank from the CLI!",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	verbose = rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose logging")
	cobra.OnInitialize(func() {
		if !*verbose {
			log.SetOutput(ioutil.Discard)
		}
	})
}
