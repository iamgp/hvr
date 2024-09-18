package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestClientServerIntegration(t *testing.T) {
	// Print current working directory and its contents
	pwd, _ := os.Getwd()
	t.Logf("Current working directory: %s", pwd)
	files, _ := ioutil.ReadDir(".")
	for _, f := range files {
		t.Logf("File: %s", f.Name())
	}

	// Reset the database
	resetCmd := exec.Command("make", "reset-db")
	if output, err := resetCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to reset database: %v\nOutput: %s", err, output)
	}

	// Build the server binary
	buildCmd := exec.Command("go", "build", "-o", "hvr-server", "./cmd/server")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build server: %v\nOutput: %s", err, output)
	}

	// Build the client binary
	buildCmd = exec.Command("go", "build", "-o", "hvr", "./pkg/client")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build client: %v\nOutput: %s", err, output)
	}

	// Start the server
	serverCmd := exec.Command("./hvr-server")
	if err := serverCmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer serverCmd.Process.Kill()

	// Wait for the server to start
	time.Sleep(5 * time.Second)

	// Run client commands
	tempDir, err := os.MkdirTemp("", "hvr-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name           string
		args           []string
		wantErr        bool
		expectedOutput string
		additionalCheck func(*testing.T, string)
	}{
		{"Upload", []string{"upload", "testdata/test-lib.zip", "--name", "test-lib", "--version", "1.0.0"}, false, "Library test-lib version 1.0.0 uploaded successfully", nil},
		{"Upload Duplicate", []string{"upload", "testdata/test-lib.zip", "--name", "test-lib", "--version", "1.0.0"}, true, "Error: library version already exists: test-lib 1.0.0", nil},
		{"Download", []string{"download", "test-lib", "1.0.0"}, false, "Library downloaded and verified successfully", func(t *testing.T, output string) {
			if _, err := os.Stat("test-lib-1.0.0.zip"); os.IsNotExist(err) {
				t.Errorf("Downloaded file not found: test-lib-1.0.0.zip")
			}
		}},
		{"Download Non-existent", []string{"download", "non-existent-lib", "1.0.0"}, true, "failed to download library", nil},
		{"Search", []string{"search", "test"}, false, "test-lib", nil},
		{"Search No Results", []string{"search", "nonexistent"}, false, "No libraries found", nil},
			// Remove the "List Versions" test if the command doesn't exist
		// {"Upload New Version", []string{"upload", "testdata/test-lib-v2.zip", "--name", "test-lib", "--version", "2.0.0"}, false, "Successfully uploaded test-lib version 2.0.0", nil},
		{"Download to Specific Directory", []string{"download", "test-lib", "1.0.0", "-o", tempDir}, false, "Library downloaded and verified successfully", func(t *testing.T, output string) {
			expectedFile := filepath.Join(tempDir, "test-lib-1.0.0.zip")
			if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
				t.Errorf("Downloaded file not found: %s", expectedFile)
			}
		}},
		{"Install", []string{"install", "test-lib", "1.0.0"}, false, "Library test-lib version 1.0.0 installed successfully", nil},
		// Remove the "Uninstall" test if the command doesn't exist
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./hvr", tt.args...)
			output, err := cmd.CombinedOutput()

			if (err != nil) != tt.wantErr {
				t.Errorf("Command %s failed: %v\nOutput: %s", tt.name, err, output)
			}

			if !strings.Contains(string(output), tt.expectedOutput) {
				t.Errorf("Command %s output does not contain expected string.\nExpected to contain: %s\nGot: %s", tt.name, tt.expectedOutput, output)
			}

			// Check if help text is printed unexpectedly
			if strings.Contains(string(output), "Usage:") && !tt.wantErr {
				t.Errorf("Command %s printed help text, which suggests it didn't run correctly. Output: %s", tt.name, output)
			}

			if tt.additionalCheck != nil {
				tt.additionalCheck(t, string(output))
			}
		})
	}

	// Update JSON Output test
	t.Run("JSON Output", func(t *testing.T) {
		cmd := exec.Command("./hvr", "search", "test", "--json")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("JSON output command failed: %v\nOutput: %s", err, output)
		}

		var result []map[string]string
		if err := json.Unmarshal(output, &result); err != nil {
			t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, output)
		}

		if len(result) == 0 {
			t.Errorf("Expected non-empty JSON result")
		}

		if result[0]["name"] != "test-lib" {
			t.Errorf("Expected library name 'test-lib', got %v", result[0]["name"])
		}
	})

	// Remove the Version Comparison test if the command doesn't exist
}

func TestServerStartupFailure(t *testing.T) {
	cmd := exec.Command("./hvr-server", "invalid")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("Expected server to fail with invalid port")
	}
	if !strings.Contains(string(output), "Invalid port number") {
		t.Errorf("Expected error message about invalid port, got: %s", string(output))
	}
}
