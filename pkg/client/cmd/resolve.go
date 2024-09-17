package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var resolveCmd = &cobra.Command{
	Use:   "resolve <name> <version>",
	Short: "Resolve dependencies for a library",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		version := args[1]

		dependencies, err := resolveDependencies(name, version)
		if err != nil {
			return err
		}

		fmt.Printf("Dependencies for %s version %s:\n", name, version)
		for _, dep := range dependencies {
			fmt.Printf("- %s (%s)\n", dep.Name, dep.Version)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(resolveCmd)
}
