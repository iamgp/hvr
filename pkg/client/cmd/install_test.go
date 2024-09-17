package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestInstallCommand(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := ioutil.TempDir("", "hvr-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock hvr.json file
	hvrJSON := `{
		"name": "test-project",
		"version": "1.0.0",
		"dependencies": {
			"test-lib": "^1.0.0"
		}
	}`
	err = ioutil.WriteFile(filepath.Join(tempDir, "hvr.json"), []byte(hvrJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create hvr.json: %v", err)
	}

	// Change to the temp directory
	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)

	// Run the install command
	cmd := installCmd
	cmd.SetArgs([]string{"--dir", "vendor"})
	err = cmd.Execute()
	if err != nil {
		t.Fatalf("Install command failed: %v", err)
	}

	// Check if the vendor directory was created
	if _, err := os.Stat("vendor"); os.IsNotExist(err) {
		t.Errorf("Vendor directory was not created")
	}

	// TODO: Add more checks for installed dependencies
}
