package commands

import (
	"fmt"
	"opcua_check/internal/config"

	"github.com/spf13/cobra"
)

func BuildRootCommand(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use: "opcua_check",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := c.InitializeConfig(cmd)
			if err != nil {
				fmt.Println(err)
				return
			}
		},
	}

	cmd.AddCommand(buildGetCommand(c))
	return cmd
}
