package models

import "github.com/Masterminds/semver/v3"

type Library struct {
	Name         string            `json:"name"`
	Version      *semver.Version   `json:"version"`
	Description  string            `json:"description"`
	Author       string            `json:"author"`
	RepoURL      string            `json:"repo_url"`
	FilePath     string            `json:"file_path"`
	Hash         string            `json:"hash"`
	Dependencies map[string]string `json:"dependencies"`
}
