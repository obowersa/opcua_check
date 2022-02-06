package commands

import (
	"opcua_check/internal/config"

	"github.com/spf13/cobra"
)

func buildGetCommand(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Commands which get information",
		Long:  "Groups our get commands together. These commands provide no modification options",
	}
	cmd.AddCommand(buildVariableCommands(c))
	cmd.AddCommand(buildConnectionCommands(c))

	return cmd
}

func buildVariableCommands(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "variables <optional_variable_group>",
		Short: "Commands for interacting with opcua variables",
		Long:  "Queries an OPCUA endpoint for specified variables. Acces a combination of variable flags and config file variables",
		Example: `
Usage examples for opcua_check get variables
--------------------------------------------
Query using variable group defined in configfile
>opcua_check get dummy_plc_1 --endpoint=opc.tcp://localhost:5002

Query using variables flag
>opcua_check get --variables="ns=4;s=PLC_PLC_One_Dummy_Var_1,ns=4;s=OPC_PLC_Two_Dummy_Var_2" --endpoint=opc.tcp://localhost:5002

Query using combination of both
opcua_check get dummmy_plc_1 --variables="ns=4;s=PLC_PLC_One_Dummy_Var_1" --endpoint opc.tcp://localhost:5001
`,

		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 || cmd.Flags().Changed("variables") {
				getVariables(c, args)
			} else {
				_ = cmd.Help()
			}
		},
	}

	addEndpointFlags(cmd, c)
	cmd.Flags().StringSliceVarP(&c.Variables, "variables", "v", []string{}, "variables to query")

	return cmd
}

func buildConnectionCommands(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connection",
		Short: "Commands for checking opcua connection information",
	}

	addEndpointFlags(cmd, c)
	cmd.AddCommand(buildConnectionStatusCommand(c))

	return cmd
}
