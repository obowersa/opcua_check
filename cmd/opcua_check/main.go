package main

import (
	"opcua_check/internal/commands"
	"opcua_check/internal/config"
	"os"
)

func main() {
	cmd := commands.BuildRootCommand(config.NewConfig())

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
