package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/iamgp/hvr/internal/api/handlers"
	"github.com/iamgp/hvr/internal/services"
	"github.com/iamgp/hvr/internal/storage"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Hamilton Venus Registry!")
}

func main() {
	port := "8080" // default port
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	// Validate the port
	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("Invalid port number: %s", port)
	}

	db, err := storage.NewSQLiteDatabase("./hvpm.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	fileStore, err := storage.NewLocalFileStore("./library_files")
	if err != nil {
		log.Fatalf("Failed to initialize file store: %v", err)
	}

	libraryService := services.NewLibraryService(db, fileStore)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/upload", handlers.UploadHandler(libraryService))
	http.HandleFunc("/download", handlers.DownloadHandler(libraryService))
	http.HandleFunc("/search", handlers.SearchHandler(libraryService))
	http.HandleFunc("/resolve", handlers.ResolveDependenciesHandler(libraryService))

	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
