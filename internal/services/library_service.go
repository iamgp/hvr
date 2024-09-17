package services

import (
	"io"

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

func (s *LibraryService) Upload(name, version string, data io.Reader) error {
	filePath, err := s.fileStore.Save(name, version, data)
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

func (s *LibraryService) Download(name, version string) (io.ReadCloser, error) {
	library, err := s.db.Get(name, version)
	if err != nil {
		return nil, err
	}

	return s.fileStore.Get(library.FilePath)
}

func (s *LibraryService) Search(query string) ([]models.Library, error) {
	return s.db.Search(query)
}

// Implement methods for library management
