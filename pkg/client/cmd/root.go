package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hvr",
	Short: "Hamilton Venus Registry CLI",
	Long:  `A command-line interface for interacting with the Hamilton Venus Registry.`,
}

func Execute(args ...string) {
	rootCmd.SetArgs(args)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add any global flags here
}
