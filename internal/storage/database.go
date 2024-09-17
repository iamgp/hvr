package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/your-username/hamilton-venus-registry/internal/models"
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
			file_path TEXT,
			PRIMARY KEY (name, version)
		)
	`)
	return err
}

func (db *SQLiteDatabase) Save(library models.Library) error {
	_, err := db.db.Exec("INSERT OR REPLACE INTO libraries (name, version, file_path) VALUES (?, ?, ?)",
		library.Name, library.Version, library.FilePath)
	return err
}

func (db *SQLiteDatabase) Get(name, version string) (models.Library, error) {
	var library models.Library
	err := db.db.QueryRow("SELECT name, version, file_path FROM libraries WHERE name = ? AND version = ?",
		name, version).Scan(&library.Name, &library.Version, &library.FilePath)
	if err == sql.ErrNoRows {
		return models.Library{}, fmt.Errorf("library %s version %s not found", name, version)
	}
	return library, err
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
	var library models.Library
	err := db.db.QueryRow("SELECT name, version, file_path FROM libraries WHERE name = ? ORDER BY version DESC LIMIT 1", name).Scan(&library.Name, &library.Version, &library.FilePath)
	if err != nil {
		return models.Library{}, err
	}
	return library, nil
}
