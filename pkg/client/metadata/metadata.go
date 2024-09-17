package metadata

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type Metadata struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	RepoURL     string   `json:"repo_url"`
	Files       []string `json:"files"`
}

func ParseMetadataFile(filename string) (*Metadata, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var meta Metadata
	err = json.Unmarshal(data, &meta)
	if err != nil {
		return nil, err
	}

	// Resolve glob patterns in Files
	var resolvedFiles []string
	for _, pattern := range meta.Files {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		resolvedFiles = append(resolvedFiles, matches...)
	}
	meta.Files = resolvedFiles

	return &meta, nil
}
