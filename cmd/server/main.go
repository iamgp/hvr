package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/your-username/hamilton-venus-registry/internal/api/handlers"
	"github.com/your-username/hamilton-venus-registry/internal/services"
	"github.com/your-username/hamilton-venus-registry/internal/storage"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Hamilton Venus Registry!")
}

func main() {
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

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
