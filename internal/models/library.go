package models

type Library struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	FilePath string `json:"file_path"`
}
