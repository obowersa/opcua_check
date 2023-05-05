package main

import (
	"os"

	"github.com/obowersa/opcuacheck/internal/commands"
	"github.com/obowersa/opcuacheck/internal/config"
)

func main() {
	cmd := commands.BuildRootCommand(config.NewConfig())

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
