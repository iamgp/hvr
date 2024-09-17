package main

import (
	"os/exec"
	"testing"
	"time"
)

func TestClientServerIntegration(t *testing.T) {
	// Start the server
	serverCmd := exec.Command("./hvr-server")
	err := serverCmd.Start()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer serverCmd.Process.Kill()

	// Wait for the server to start
	time.Sleep(2 * time.Second)

	// Run client commands
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"Upload", []string{"upload", "testdata/test-lib.zip", "--name", "test-lib", "--version", "1.0.0"}, false},
		{"Download", []string{"download", "test-lib", "1.0.0"}, false},
		{"Search", []string{"search", "test"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./hvr", tt.args...)
			output, err := cmd.CombinedOutput()
			if (err != nil) != tt.wantErr {
				t.Errorf("Command %s failed: %v\nOutput: %s", tt.name, err, output)
			}
			// TODO: Add more specific checks for each command's output
		})
	}
}
