package commands

import (
	"context"
	"time"

	"github.com/obowersa/opcuacheck/internal/config"

	"github.com/spf13/cobra"
)

func contextWithTimeout(c *config.Config) (context.Context, context.CancelFunc) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.Timeout)*time.Second)

	return ctx, cancel
}

func addEndpointFlags(cmd *cobra.Command, c *config.Config) {
	cmd.PersistentFlags().StringVarP(&c.ConfigFile, "config", "c", "./configs/config.yaml", "path to the config file, defaults to ./configs/config.yaml")
	cmd.PersistentFlags().StringVarP(&c.Output, "output", "o", "", "specific output format, example: --output json")
	cmd.PersistentFlags().StringVarP(&c.Endpoint, "endpoint", "e", "", "opcua endpoint with opc.tcp://<addr>:<port> format")
	cmd.PersistentFlags().IntVarP(&c.Timeout, "timeout", "t", 3, "number of seconds to wait before timing out an opcua connection")
	_ = cmd.MarkPersistentFlagRequired("endpoint")
}
