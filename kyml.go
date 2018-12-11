package main

import (
	"os"

	"github.com/frigus02/kyml/pkg/commands"
)

func main() {
	if err := commands.NewDefaultCommand().Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
