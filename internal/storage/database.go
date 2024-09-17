package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/iamgp/hvr/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDatabase struct {
	db *sql.DB
}

func NewSQLiteDatabase(dbPath string) (*SQLiteDatabase, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := createTable(db); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &SQLiteDatabase{db: db}, nil
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS libraries (
			name TEXT,
			version TEXT,
			description TEXT,
			author TEXT,
			repo_url TEXT,
			file_path TEXT,
			hash TEXT,
			dependencies TEXT,
			PRIMARY KEY (name, version)
		)
	`)
	return err
}

func (db *SQLiteDatabase) Save(library models.Library) error {
	dependenciesJSON, err := json.Marshal(library.Dependencies)
	if err != nil {
		return fmt.Errorf("failed to marshal dependencies: %w", err)
	}

	_, err = db.db.Exec("INSERT OR REPLACE INTO libraries (name, version, description, author, repo_url, file_path, hash, dependencies) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		library.Name, library.Version.String(), library.Description, library.Author, library.RepoURL, library.FilePath, library.Hash, string(dependenciesJSON))
	return err
}

func (db *SQLiteDatabase) Get(name, version string) (models.Library, error) {
	var library models.Library
	var versionStr, dependenciesJSON string
	err := db.db.QueryRow("SELECT name, version, description, author, repo_url, file_path, hash, dependencies FROM libraries WHERE name = ? AND version = ?",
		name, version).Scan(&library.Name, &versionStr, &library.Description, &library.Author, &library.RepoURL, &library.FilePath, &library.Hash, &dependenciesJSON)
	if err == sql.ErrNoRows {
		return models.Library{}, fmt.Errorf("library %s version %s not found", name, version)
	}
	if err != nil {
		return models.Library{}, err
	}

	library.Version, err = semver.NewVersion(versionStr)
	if err != nil {
		return models.Library{}, fmt.Errorf("invalid version: %w", err)
	}

	err = json.Unmarshal([]byte(dependenciesJSON), &library.Dependencies)
	if err != nil {
		return models.Library{}, fmt.Errorf("failed to unmarshal dependencies: %w", err)
	}

	return library, nil
}

func (db *SQLiteDatabase) Search(query string) ([]models.Library, error) {
	rows, err := db.db.Query("SELECT name, version FROM libraries WHERE name LIKE ?", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var libraries []models.Library
	for rows.Next() {
		var lib models.Library
		if err := rows.Scan(&lib.Name, &lib.Version); err != nil {
			return nil, err
		}
		libraries = append(libraries, lib)
	}
	return libraries, rows.Err()
}

func (db *SQLiteDatabase) Close() error {
	return db.db.Close()
}

func (db *SQLiteDatabase) GetLatest(name string) (models.Library, error) {
	rows, err := db.db.Query("SELECT name, version, description, author, repo_url, file_path, hash, dependencies FROM libraries WHERE name = ? ORDER BY version DESC", name)
	if err != nil {
		return models.Library{}, err
	}
	defer rows.Close()

	var latestLibrary models.Library
	var latestVersion *semver.Version

	for rows.Next() {
		var library models.Library
		var versionStr, dependenciesJSON string
		err := rows.Scan(&library.Name, &versionStr, &library.Description, &library.Author, &library.RepoURL, &library.FilePath, &library.Hash, &dependenciesJSON)
		if err != nil {
			return models.Library{}, err
		}

		version, err := semver.NewVersion(versionStr)
		if err != nil {
			continue // Skip invalid versions
		}

		if latestVersion == nil || version.GreaterThan(latestVersion) {
			latestVersion = version
			latestLibrary = library
			latestLibrary.Version = version
			json.Unmarshal([]byte(dependenciesJSON), &latestLibrary.Dependencies)
		}
	}

	if latestVersion == nil {
		return models.Library{}, fmt.Errorf("no valid versions found for library %s", name)
	}

	return latestLibrary, nil
}

// Add this new method to the SQLiteDatabase struct
func (db *SQLiteDatabase) GetAllVersions(name string) ([]*semver.Version, error) {
	rows, err := db.db.Query("SELECT version FROM libraries WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []*semver.Version
	for rows.Next() {
		var versionStr string
		err := rows.Scan(&versionStr)
		if err != nil {
			return nil, err
		}
		version, err := semver.NewVersion(versionStr)
		if err != nil {
			continue // Skip invalid versions
		}
		versions = append(versions, version)
	}

	return versions, nil
}
