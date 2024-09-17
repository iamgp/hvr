package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload <file>",
	Short: "Upload a library to the registry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		name, _ := cmd.Flags().GetString("name")
		version, _ := cmd.Flags().GetString("version")

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", filepath.Base(filePath))
		if err != nil {
			return fmt.Errorf("failed to create form file: %w", err)
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return fmt.Errorf("failed to copy file content: %w", err)
		}

		writer.WriteField("name", name)
		writer.WriteField("version", version)

		// Get file info
		fileInfo, err := file.Stat()
		if err != nil {
			return fmt.Errorf("failed to get file info: %w", err)
		}

		// Add modification time to the form data
		modTime := fileInfo.ModTime().Unix()
		writer.WriteField("modTime", fmt.Sprintf("%d", modTime))

		err = writer.Close()
		if err != nil {
			return fmt.Errorf("failed to close multipart writer: %w", err)
		}

		req, err := http.NewRequest("POST", "http://localhost:8080/upload", body)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send request: %w", err)
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		if resp.StatusCode != http.StatusCreated {
			if resp.StatusCode == http.StatusConflict {
				return fmt.Errorf("upload failed: %s", string(responseBody))
			}
			return fmt.Errorf("upload failed with status: %s, body: %s", resp.Status, string(responseBody))
		}

		fmt.Printf("Library %s version %s uploaded successfully\n", name, version)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().String("name", "", "Name of the library")
	uploadCmd.Flags().String("version", "", "Version of the library")
	uploadCmd.MarkFlagRequired("name")
	uploadCmd.MarkFlagRequired("version")
}
