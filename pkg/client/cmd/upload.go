package cmd

import (
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload <file>",
	Short: "Upload a library to the registry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		name, _ := cmd.Flags().GetString("name")
		version, _ := cmd.Flags().GetString("version")
		description, _ := cmd.Flags().GetString("description")
		author, _ := cmd.Flags().GetString("author")
		repoURL, _ := cmd.Flags().GetString("repo-url")

		return uploadLibrary(filePath, name, version, description, author, repoURL)
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().String("name", "", "Name of the library")
	uploadCmd.Flags().String("version", "", "Version of the library")
	uploadCmd.Flags().String("description", "", "Description of the library")
	uploadCmd.Flags().String("author", "", "Author of the library")
	uploadCmd.Flags().String("repo-url", "", "Repository URL of the library")
	uploadCmd.MarkFlagRequired("name")
	uploadCmd.MarkFlagRequired("version")
}
