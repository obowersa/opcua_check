package commands

import (
	"fmt"
	"opcua_check/internal/config"
	"opcua_check/pkg/opcua"

	"github.com/spf13/cobra"
)

func connectionStatus(c *config.Config) {
	o := opcua.NewOpcua(c.Endpoint)

	ctx, cancel := contextWithTimeout(c)
	defer cancel()

	if err := o.CheckConnection(ctx); err != nil {
		fmt.Printf("Unable to connect to %v, got error: %v\n", c.Endpoint, err)
	} else {
		fmt.Printf("Connection to %v looks good\n", c.Endpoint)
	}
}

func buildConnectionStatusCommand(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Checks opcua connectivity",
		Long:  "Attempts to connect to the specified opcua endpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			connectionStatus(c)
			return nil
		},
	}

	return cmd
}
