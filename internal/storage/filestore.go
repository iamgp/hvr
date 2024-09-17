package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FileStore interface {
	Save(name, version string, data io.Reader, modTime time.Time) (string, error)
	Get(path string) ([]byte, time.Time, error)
}

type LocalFileStore struct {
	baseDir string
}

func NewLocalFileStore(baseDir string) (*LocalFileStore, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}
	return &LocalFileStore{baseDir: baseDir}, nil
}

func (fs *LocalFileStore) Save(name, version string, data io.Reader, modTime time.Time) (string, error) {
	filename := filepath.Join(fs.baseDir, name, version+".zip")
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return "", err
	}
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = io.Copy(file, data)
	if err != nil {
		return "", err
	}

	// Set the modification time of the file
	err = os.Chtimes(filename, time.Now(), modTime)
	if err != nil {
		return "", fmt.Errorf("failed to set modification time: %w", err)
	}

	return filename, nil
}

func (fs *LocalFileStore) Get(path string) ([]byte, time.Time, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, time.Time{}, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, time.Time{}, err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, time.Time{}, err
	}

	return content, fileInfo.ModTime(), nil
}
