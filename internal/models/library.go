package models

type Library struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      string `json:"author"`
	RepoURL     string `json:"repo_url"`
	FilePath    string `json:"file_path"`
	Hash        string `json:"hash"`
}
