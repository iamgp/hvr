package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var jsonOutput bool

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for libraries",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("search query is required")
		}
		query := args[0]

		// TODO: Implement actual search logic here
		// For now, just return a dummy result that uses the query
		results := []map[string]string{
			{"name": "test-lib", "version": "1.0.0"},
			{"name": "another-lib", "version": "2.0.0"},
		}

		// Filter results based on the query
		filteredResults := make([]map[string]string, 0)
		for _, result := range results {
			if strings.Contains(result["name"], query) {
				filteredResults = append(filteredResults, result)
			}
		}

		if jsonOutput {
			jsonData, err := json.Marshal(filteredResults)
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(jsonData))
		} else {
			if len(filteredResults) == 0 {
				fmt.Println("No libraries found")
			} else {
				for _, result := range filteredResults {
					fmt.Printf("%s (%s)\n", result["name"], result["version"])
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results in JSON format")
}
