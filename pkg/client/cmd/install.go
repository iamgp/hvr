package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// serverURL is the URL of the server that the client communicates with
var serverURL string

var installDir string

var installCmd = &cobra.Command{
	Use:   "install [library] [version]",
	Short: "Install a library",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("library name is required")
		}
		name := args[0]
		version := "latest"
		if len(args) > 1 {
			version = args[1]
		}

		// Validate version (simple check for now)
		if version != "latest" && !isValidVersion(version) {
			return fmt.Errorf("invalid version format: %s", version)
		}

		// Create the installation directory
		if err := os.MkdirAll(installDir, 0755); err != nil {
			return fmt.Errorf("failed to create installation directory: %w", err)
		}

		// TODO: Implement actual installation logic here
		// For now, just create an empty file to simulate installation
		dummyFile := filepath.Join(installDir, fmt.Sprintf("%s-%s.txt", name, version))
		if _, err := os.Create(dummyFile); err != nil {
			return fmt.Errorf("failed to create dummy file: %w", err)
		}

		fmt.Printf("Library %s version %s installed successfully in %s\n", name, version, installDir)
		return nil
	},
}

func isValidVersion(version string) bool {
	// Add more sophisticated version validation if needed
	return len(version) > 0
}

func init() {
	// Set a default value for serverURL
	serverURL = "http://localhost:8080"  // Adjust this to your default server URL

	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringVarP(&installDir, "dir", "d", "vendor", "Installation directory")
}
