package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iamgp/hvr/internal/models"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for libraries in the registry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/search?q=%s", query))
		if err != nil {
			return fmt.Errorf("failed to search: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("search failed with status: %s", resp.Status)
		}

		var results []models.Library
		if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
			return fmt.Errorf("failed to decode search results: %w", err)
		}

		if len(results) == 0 {
			fmt.Println("No libraries found")
		} else {
			fmt.Printf("Found %d libraries:\n", len(results))
			for _, lib := range results {
				fmt.Printf("- %s (version %s)\n", lib.Name, lib.Version)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
