package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download <name> <version>",
	Short: "Download a library from the registry",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		version := args[1]

		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/download?name=%s&version=%s", name, version))
		if err != nil {
			return fmt.Errorf("failed to download: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("download failed with status: %s", resp.Status)
		}

		filename := getFilenameFromHeader(resp.Header.Get("Content-Disposition"))
		if filename == "" {
			filename = fmt.Sprintf("%s-%s", name, version)
		}

		out, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return fmt.Errorf("failed to save file: %w", err)
		}

		fmt.Printf("Library downloaded successfully as %s\n", filename)
		return nil
	},
}

func getFilenameFromHeader(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.Split(header, "filename=")
	if len(parts) != 2 {
		return ""
	}
	return filepath.Base(strings.Trim(parts[1], `"`))
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
