package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/iamgp/hvr/internal/dependency"
	"github.com/iamgp/hvr/internal/models"
	"github.com/iamgp/hvr/internal/storage"
)

type LibraryService struct {
	db        *storage.SQLiteDatabase
	fileStore storage.FileStore
	resolver  *dependency.Resolver
}

func NewLibraryService(db *storage.SQLiteDatabase, fs storage.FileStore) *LibraryService {
	return &LibraryService{
		db:        db,
		fileStore: fs,
		resolver:  dependency.NewResolver(db),
	}
}

func (s *LibraryService) Upload(name, versionStr, description, author, repoURL string, dependencies map[string]string, data io.Reader, modTime time.Time) error {
	version, err := semver.NewVersion(versionStr)
	if err != nil {
		return fmt.Errorf("invalid version: %w", err)
	}

	// Check if the library version already exists
	_, err = s.db.Get(name, version.String())
	if err == nil {
		return fmt.Errorf("library version already exists: %s %s", name, version)
	}

	// Calculate hash
	hasher := sha256.New()
	teeReader := io.TeeReader(data, hasher)

	filePath, err := s.fileStore.Save(name, version.String(), teeReader, modTime)
	if err != nil {
		return err
	}

	hash := hex.EncodeToString(hasher.Sum(nil))

	library := models.Library{
		Name:         name,
		Version:      version,
		Description:  description,
		Author:       author,
		RepoURL:      repoURL,
		FilePath:     filePath,
		Hash:         hash,
		Dependencies: dependencies,
	}

	return s.db.Save(library)
}

func (s *LibraryService) Download(name, versionStr string) ([]byte, time.Time, string, error) {
	var library models.Library
	var err error

	if versionStr == "latest" {
		library, err = s.db.GetLatest(name)
	} else {
		version, err := semver.NewVersion(versionStr)
		if err != nil {
			return nil, time.Time{}, "", fmt.Errorf("invalid version: %w", err)
		}
		library, err = s.db.Get(name, version.String())
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

func (s *LibraryService) ResolveLibraryDependencies(name, version string) ([]models.Library, error) {
	library, err := s.db.Get(name, version)
	if err != nil {
		return nil, fmt.Errorf("failed to get library %s version %s: %w", name, version, err)
	}

	return s.resolver.ResolveDependencies(library)
}
