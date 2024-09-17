package storage

import (
	"io"
	"os"
	"path/filepath"
)

type FileStore interface {
	Save(name, version string, data io.Reader) (string, error)
	Get(path string) (io.ReadCloser, error)
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

func (fs *LocalFileStore) Save(name, version string, data io.Reader) (string, error) {
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
	return filename, nil
}

func (fs *LocalFileStore) Get(path string) (io.ReadCloser, error) {
	return os.Open(path)
}
