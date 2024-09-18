package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestInstallCommand(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := ioutil.TempDir("", "hvr-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to the temp directory
	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)

	// Run the install command
	cmd := installCmd
	cmd.SetArgs([]string{"test-lib", "1.0.0", "--dir", "vendor"})

	// Capture output
	output := new(bytes.Buffer)
	cmd.SetOut(output)
	cmd.SetErr(output)

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("Install command failed: %v", err)
	}

	// Check if the vendor directory was created
	vendorDir := filepath.Join(tempDir, "vendor")
	if _, err := os.Stat(vendorDir); os.IsNotExist(err) {
		t.Errorf("Vendor directory was not created")
	}

	// Check if the dummy file was created
	dummyFile := filepath.Join(vendorDir, "test-lib-1.0.0.txt")
	if _, err := os.Stat(dummyFile); os.IsNotExist(err) {
		t.Errorf("Dummy file was not created")
	}

	// Check the output
	expectedOutput := "Library test-lib version 1.0.0 installed successfully in vendor\n"
	if output.String() != expectedOutput {
		t.Errorf("Unexpected output.\nExpected: %q\nGot: %q", expectedOutput, output.String())
	}
}

func TestInstallCmd(t *testing.T) {
	// Setup temporary directory
	tempDir, err := os.MkdirTemp("", "hvr-install-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		errMsg   string
		checkFn  func(*testing.T, string)
	}{
		{
			name:    "Install Success",
			args:    []string{"test-lib", "1.0.0", "--dir", tempDir},
			wantErr: false,
			checkFn: func(t *testing.T, dir string) {
				if _, err := os.Stat(filepath.Join(dir, "test-lib-1.0.0.txt")); os.IsNotExist(err) {
					t.Errorf("Installed library not found")
				}
			},
		},
		{
			name:    "Install Failure - Invalid Version",
			args:    []string{"test-lib", "", "--dir", tempDir},
			wantErr: true,
			errMsg:  "invalid version format: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new root command for each test
			rootCmd := &cobra.Command{Use: "hvr"}
			rootCmd.AddCommand(installCmd)

			// Set the arguments
			rootCmd.SetArgs(append([]string{"install"}, tt.args...))

			// Capture output
			output := new(bytes.Buffer)
			rootCmd.SetOut(output)
			rootCmd.SetErr(output)

			// Execute the command
			err := rootCmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("Install command error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Expected error message to contain %q, got %q", tt.errMsg, err.Error())
			}

			if tt.checkFn != nil {
				tt.checkFn(t, tempDir)
			}

			// Check output for help text only if we don't expect an error
			if !tt.wantErr && strings.Contains(output.String(), "Usage:") {
				t.Errorf("Command printed help text, which suggests it didn't run correctly. Output: %s", output.String())
			}
		})
	}
}
