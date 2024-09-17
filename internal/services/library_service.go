package services

import (
	"io"
	"time"

	"github.com/your-username/hamilton-venus-registry/internal/models"
	"github.com/your-username/hamilton-venus-registry/internal/storage"
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

func (s *LibraryService) Upload(name, version string, data io.Reader, modTime time.Time) error {
	filePath, err := s.fileStore.Save(name, version, data, modTime)
	if err != nil {
		return err
	}

	library := models.Library{
		Name:     name,
		Version:  version,
		FilePath: filePath,
	}

	return s.db.Save(library)
}

func (s *LibraryService) Download(name, version string) ([]byte, time.Time, error) {
	library, err := s.db.Get(name, version)
	if err != nil {
		return nil, time.Time{}, err
	}

	fileContent, modTime, err := s.fileStore.Get(library.FilePath)
	if err != nil {
		return nil, time.Time{}, err
	}

	return fileContent, modTime, nil
}

func (s *LibraryService) Search(query string) ([]models.Library, error) {
	return s.db.Search(query)
}

// Implement methods for library management
