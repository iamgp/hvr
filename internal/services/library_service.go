package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/iamgp/hvr/internal/models"
	"github.com/iamgp/hvr/internal/storage"
)

type LibraryService struct {
	db        *storage.SQLiteDatabase
	fileStore storage.FileStore
}

func NewLibraryService(db *storage.SQLiteDatabase, fs storage.FileStore) *LibraryService {
	return &LibraryService{
		db:        db,
		fileStore: fs,
	}
}

func (s *LibraryService) Upload(name, version, description, author, repoURL string, data io.Reader, modTime time.Time) error {
	// Check if the library version already exists
	_, err := s.db.Get(name, version)
	if err == nil {
		// Library version already exists
		return fmt.Errorf("library version already exists: %s %s", name, version)
	}

	// Calculate hash
	hasher := sha256.New()
	teeReader := io.TeeReader(data, hasher)

	filePath, err := s.fileStore.Save(name, version, teeReader, modTime)
	if err != nil {
		return err
	}

	hash := hex.EncodeToString(hasher.Sum(nil))

	library := models.Library{
		Name:        name,
		Version:     version,
		Description: description,
		Author:      author,
		RepoURL:     repoURL,
		FilePath:    filePath,
		Hash:        hash,
	}

	return s.db.Save(library)
}

func (s *LibraryService) Download(name, version string) ([]byte, time.Time, string, error) {
	var library models.Library
	var err error

	if version == "latest" {
		library, err = s.db.GetLatest(name)
	} else {
		library, err = s.db.Get(name, version)
	}

	if err != nil {
		return nil, time.Time{}, "", err
	}

	fileContent, modTime, err := s.fileStore.Get(library.FilePath)
	if err != nil {
		return nil, time.Time{}, "", err
	}

	return fileContent, modTime, library.Hash, nil
}

func (s *LibraryService) Search(query string) ([]models.Library, error) {
	return s.db.Search(query)
}

// Implement methods for library management
