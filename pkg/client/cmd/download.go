package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download <name> <version>",
	Short: "Download a library from the registry",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		version := args[1]

		url := fmt.Sprintf("http://localhost:8080/download?name=%s&version=%s", name, version)
		fmt.Printf("Downloading from: %s\n", url)

		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to download: %w", err)
		}
		defer resp.Body.Close()

		fmt.Printf("Response status: %s\n", resp.Status)

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("download failed with status: %s, body: %s", resp.Status, string(body))
		}

		filename := getFilenameFromHeader(resp.Header.Get("Content-Disposition"))
		if filename == "" {
			filename = fmt.Sprintf("%s-%s.zip", name, version)
		}
		fmt.Printf("Saving file as: %s\n", filename)

		out, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer out.Close()

		n, err := io.Copy(out, resp.Body)
		if err != nil {
			return fmt.Errorf("failed to save file: %w", err)
		}
		fmt.Printf("Wrote %d bytes to file\n", n)

		modTimeStr := resp.Header.Get("X-File-ModTime")
		if modTimeStr != "" {
			modTime, err := strconv.ParseInt(modTimeStr, 10, 64)
			if err == nil {
				err = os.Chtimes(filename, time.Now(), time.Unix(modTime, 0))
				if err != nil {
					fmt.Printf("Warning: Failed to set modification time: %v\n", err)
				} else {
					fmt.Printf("Set modification time to: %s\n", time.Unix(modTime, 0))
				}
			}
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
