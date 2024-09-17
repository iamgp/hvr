package models

import (
	"testing"

	"github.com/Masterminds/semver/v3"
)

func TestLibrary(t *testing.T) {
	version, _ := semver.NewVersion("1.0.0")
	lib := Library{
		Name:        "test-lib",
		Version:     version,
		Description: "A test library",
		Author:      "Test Author",
		RepoURL:     "https://github.com/test/test-lib",
		FilePath:    "/path/to/test-lib.zip",
		Hash:        "abcdef1234567890",
		Dependencies: map[string]string{
			"dep1": "^1.0.0",
			"dep2": "~2.0.0",
		},
	}

	if lib.Name != "test-lib" {
		t.Errorf("Expected name to be 'test-lib', got '%s'", lib.Name)
	}

	if lib.Version.String() != "1.0.0" {
		t.Errorf("Expected version to be '1.0.0', got '%s'", lib.Version.String())
	}

	if len(lib.Dependencies) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(lib.Dependencies))
	}
}
