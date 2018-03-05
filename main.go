package main

import (
	"os"

	"github.com/jaipradeesh/iptablesbeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
