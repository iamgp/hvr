package cmd

import (
	"crypto/sha256"
	"encoding/hex"
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

var outputDir string

func downloadLibrary(name, version, destPath string) error {
	url := fmt.Sprintf("http://localhost:8080/download?name=%s&version=%s", name, version)
	fmt.Printf("Downloading from: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("download failed with status: %s, body: %s", resp.Status, string(body))
	}

	filename := getFilenameFromHeader(resp.Header.Get("Content-Disposition"))
	if filename == "" {
		filename = fmt.Sprintf("%s-%s.zip", name, version)
	}
	filePath := filepath.Join(destPath, filename)

	err = os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	expectedHash := resp.Header.Get("X-File-Hash")
	hasher := sha256.New()
	teeReader := io.TeeReader(resp.Body, hasher)

	n, err := io.Copy(out, teeReader)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	fmt.Printf("Wrote %d bytes to file\n", n)

	actualHash := hex.EncodeToString(hasher.Sum(nil))
	if actualHash != expectedHash {
		os.Remove(filePath) // Delete the file if hash doesn't match
		return fmt.Errorf("hash mismatch: expected %s, got %s", expectedHash, actualHash)
	}

	modTimeStr := resp.Header.Get("X-File-ModTime")
	if modTimeStr != "" {
		modTime, err := strconv.ParseInt(modTimeStr, 10, 64)
		if err == nil {
			err = os.Chtimes(filePath, time.Now(), time.Unix(modTime, 0))
			if err != nil {
				fmt.Printf("Warning: Failed to set modification time: %v\n", err)
			} else {
				fmt.Printf("Set modification time to: %s\n", time.Unix(modTime, 0))
			}
		}
	}

	fmt.Printf("Library downloaded and verified successfully as %s\n", filePath)
	return nil
}

var downloadCmd = &cobra.Command{
	Use:   "download <name> [version]",
	Short: "Download a library",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("library name is required")
		}
		name := args[0]
		version := "latest"
		if len(args) > 1 {
			version = args[1]
		}

		// Use outputDir if specified, otherwise use current directory
		downloadPath := "."
		if outputDir != "" {
			downloadPath = outputDir
		}

		// Implement actual download logic here
		err := downloadLibrary(name, version, downloadPath)
		if err != nil {
			return fmt.Errorf("failed to download library: %w", err)
		}

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
	downloadCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for downloaded files")
}
