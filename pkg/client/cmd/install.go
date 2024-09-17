package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/iamgp/hvr/pkg/client/metadata"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install dependencies specified in hvr.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read hvr.json file
		data, err := ioutil.ReadFile("hvr.json")
		if err != nil {
			return fmt.Errorf("failed to read hvr.json: %w", err)
		}

		var meta metadata.Metadata
		err = json.Unmarshal(data, &meta)
		if err != nil {
			return fmt.Errorf("failed to parse hvr.json: %w", err)
		}

		// Get the installation directory
		installDir, _ := cmd.Flags().GetString("dir")
		if installDir == "" {
			installDir = "./vendor"
		}

		// Create the installation directory
		err = os.MkdirAll(installDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create installation directory: %w", err)
		}

		// Install each dependency
		for depName, depVersion := range meta.Dependencies {
			fmt.Printf("Installing %s (%s)\n", depName, depVersion)

			// Resolve the dependency
			resolvedDeps, err := resolveDependencies(depName, depVersion)
			if err != nil {
				return fmt.Errorf("failed to resolve dependencies for %s: %w", depName, err)
			}

			// Download and install the resolved dependencies
			for _, dep := range resolvedDeps {
				fmt.Printf("  Installing %s (%s)\n", dep.Name, dep.Version)
				err = downloadLibrary(dep.Name, dep.Version.String(), filepath.Join(installDir, dep.Name))
				if err != nil {
					return fmt.Errorf("failed to install %s: %w", dep.Name, err)
				}
			}
		}

		fmt.Printf("Successfully installed all dependencies in %s\n", installDir)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringP("dir", "d", "", "Directory to install libraries (default is ./vendor)")
}
