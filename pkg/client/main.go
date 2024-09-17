package main

import (
	"fmt"
	"os"

	"github.com/your-username/hamilton-venus-registry/pkg/client/cmd"
	"github.com/your-username/hamilton-venus-registry/pkg/client/ui"
)

func main() {
	action := ui.RunMainTUI()

	switch action {
	case "Upload":
		name, version, file := ui.RunUploadTUI()
		if name != "" && version != "" && file != "" {
			cmd.Execute("upload", file, "--name", name, "--version", version)
		}
	case "Download":
		name, version := ui.RunDownloadTUI()
		if name != "" && version != "" {
			cmd.Execute("download", name, version)
		}
	case "Search":
		query := ui.RunSearchTUI()
		if query != "" {
			cmd.Execute("search", query)
		}
	case "Quit":
		fmt.Println("Goodbye!")
		os.Exit(0)
	default:
		fmt.Println("Invalid action")
		os.Exit(1)
	}
}
