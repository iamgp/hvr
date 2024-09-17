package storage

import (
	"os"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/iamgp/hvr/internal/models"
)

func TestSQLiteDatabase(t *testing.T) {
	dbPath := "test.db"
	db, err := NewSQLiteDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer os.Remove(dbPath)
	defer db.Close()

	version, _ := semver.NewVersion("1.0.0")
	lib := models.Library{
		Name:        "test-lib",
		Version:     version,
		Description: "A test library",
		Author:      "Test Author",
		RepoURL:     "https://github.com/test/test-lib",
		FilePath:    "/path/to/test-lib.zip",
		Hash:        "abcdef1234567890",
		Dependencies: map[string]string{
			"dep1": "^1.0.0",
		},
	}

	err = db.Save(lib)
	if err != nil {
		t.Fatalf("Failed to save library: %v", err)
	}

	retrieved, err := db.Get("test-lib", "1.0.0")
	if err != nil {
		t.Fatalf("Failed to retrieve library: %v", err)
	}

	if retrieved.Name != lib.Name {
		t.Errorf("Expected name to be '%s', got '%s'", lib.Name, retrieved.Name)
	}

	if retrieved.Version.String() != lib.Version.String() {
		t.Errorf("Expected version to be '%s', got '%s'", lib.Version.String(), retrieved.Version.String())
	}
}
