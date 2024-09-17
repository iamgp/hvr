package main

import (
	"fmt"
	"os"

	"github.com/your-username/hamilton-venus-registry/pkg/client/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
