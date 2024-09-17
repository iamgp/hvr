package handlers

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/your-username/hamilton-venus-registry/internal/services"
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

		// Now upload the zipped file
		err = s.Upload(name, version, zipBuffer)
		if err != nil {
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
		if name == "" || version == "" {
			http.Error(w, "Missing name or version parameter", http.StatusBadRequest)
			return
		}

		file, err := s.Download(name, version)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s-%s.zip", name, version))
		_, err = io.Copy(w, file)
		if err != nil {
			log.Printf("Error sending file: %v", err)
			http.Error(w, "Error sending file", http.StatusInternalServerError)
			return
		}
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
