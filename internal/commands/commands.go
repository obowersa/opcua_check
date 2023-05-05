package commands

import (
	"fmt"

	"github.com/obowersa/opcuacheck/internal/config"
	"github.com/spf13/cobra"
)

func BuildRootCommand(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use: "opcuacheck",
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
