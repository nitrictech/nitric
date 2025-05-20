package main

import (
	"os"

	"github.com/nitrictech/nitric/cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
