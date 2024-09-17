package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"os"

	"github.com/iamgp/hvr/pkg/client/metadata"
	"github.com/spf13/cobra"
)

var uploadMetaCmd = &cobra.Command{
	Use:   "uploadmeta <metadata-file>",
	Short: "Upload a library to the registry using a metadata file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		metadataFile := args[0]

		// Parse the metadata file
		meta, err := metadata.ParseMetadataFile(metadataFile)
		if err != nil {
			return fmt.Errorf("failed to parse metadata file: %w", err)
		}

		// Create a temporary zip file
		tempFile, err := os.CreateTemp("", "library-*.zip")
		if err != nil {
			return fmt.Errorf("failed to create temporary file: %w", err)
		}
		defer os.Remove(tempFile.Name())

		// Create a zip writer
		zipWriter := zip.NewWriter(tempFile)

		// Add files to the zip
		for _, file := range meta.Files {
			err = addFileToZip(zipWriter, file)
			if err != nil {
				return fmt.Errorf("failed to add file to zip: %w", err)
			}
		}

		// Close the zip writer
		err = zipWriter.Close()
		if err != nil {
			return fmt.Errorf("failed to close zip writer: %w", err)
		}

		// Reopen the zip file for reading
		tempFile.Seek(0, 0)

		// Use the existing upload logic
		return uploadLibrary(tempFile.Name(), meta.Name, meta.Version, meta.Description, meta.Author, meta.RepoURL)
	},
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filename
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}

func init() {
	rootCmd.AddCommand(uploadMetaCmd)
}
