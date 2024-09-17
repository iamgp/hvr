package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hvr",
	Short: "Hamilton Venus Registry CLI",
	Long:  `A command-line interface for interacting with the Hamilton Venus Registry.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(uploadMetaCmd)
	rootCmd.AddCommand(resolveCmd)
	rootCmd.AddCommand(installCmd)
}
