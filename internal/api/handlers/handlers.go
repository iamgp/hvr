package handlers

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/iamgp/hvr/internal/services"
)

func UploadHandler(s *services.LibraryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			log.Printf("Error parsing multipart form: %v", err)
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		version := r.FormValue("version")
		file, header, err := r.FormFile("file")
		if err != nil {
			log.Printf("Error retrieving file: %v", err)
			http.Error(w, "Error retrieving file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		modTimeStr := r.FormValue("modTime")
		modTime := time.Now() // Default to current time
		if modTimeStr != "" {
			modTimeUnix, err := strconv.ParseInt(modTimeStr, 10, 64)
			if err == nil {
				modTime = time.Unix(modTimeUnix, 0)
			}
		}

		log.Printf("Uploading file: %s, name: %s, version: %s", header.Filename, name, version)

		// Create a buffer to store our zipped file
		zipBuffer := new(bytes.Buffer)
		zipWriter := zip.NewWriter(zipBuffer)

		// Create a new file inside the zip archive
		zipFile, err := zipWriter.Create(header.Filename)
		if err != nil {
			log.Printf("Error creating zip file: %v", err)
			http.Error(w, "Error creating zip file", http.StatusInternalServerError)
			return
		}

		// Copy the uploaded file data to the zip file
		_, err = io.Copy(zipFile, file)
		if err != nil {
			log.Printf("Error copying file to zip: %v", err)
			http.Error(w, "Error copying file to zip", http.StatusInternalServerError)
			return
		}

		// Close the zip writer
		zipWriter.Close()

		// Now upload the zipped file with the modification time
		description := r.FormValue("description")
		author := r.FormValue("author")
		repoURL := r.FormValue("repoURL")

		err = s.Upload(name, version, description, author, repoURL, zipBuffer, modTime)
		if err != nil {
			if strings.Contains(err.Error(), "library version already exists") {
				log.Printf("Attempt to overwrite existing version: %v", err)
				http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusConflict)
				return
			}
			log.Printf("Error uploading file: %v", err)
			http.Error(w, fmt.Sprintf("Error uploading file: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Library uploaded successfully"})
	}
}

func DownloadHandler(s *services.LibraryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := r.URL.Query().Get("name")
		version := r.URL.Query().Get("version")

		if name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}

		if version == "" {
			version = "latest"
		}

		fileContent, modTime, hash, err := s.Download(name, version)
		if err != nil {
			log.Printf("Error downloading file: %v", err)
			http.Error(w, fmt.Sprintf("Error downloading file: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s-%s.zip", name, version))
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("X-File-ModTime", fmt.Sprintf("%d", modTime.Unix()))
		w.Header().Set("X-File-Hash", hash)

		_, err = w.Write(fileContent)
		if err != nil {
			log.Printf("Error writing file content to response: %v", err)
			http.Error(w, "Error sending file", http.StatusInternalServerError)
			return
		}

		log.Printf("File %s-%s.zip downloaded successfully", name, version)
	}
}

func SearchHandler(s *services.LibraryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		query := r.URL.Query().Get("q")
		if query == "" {
			http.Error(w, "Missing q parameter", http.StatusBadRequest)
			return
		}

		results, err := s.Search(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(results)
	}
}
