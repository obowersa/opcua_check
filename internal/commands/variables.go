package commands

import (
	"encoding/json"
	"fmt"

	"github.com/obowersa/opcuacheck/internal/config"
	"github.com/obowersa/opcuacheck/pkg/opcua"
)

func getVariables(c *config.Config, args []string) {
	opc := opcua.NewOpcua(c.Endpoint)

	ctx, cancel := contextWithTimeout(c)
	defer cancel()

	if len(args) == 1 {
		c.Variables = append(c.Variables, c.BaseVariables[args[0]]...)
	}

	vs, err := opc.GetVariables(ctx, *opcua.NewVariables(c.Variables))
	if err != nil {
		fmt.Println(err)
		return
	}

	if c.Output == "json" {
		p, _ := json.Marshal(vs)
		fmt.Println(string(p))
		return
	}
	fmt.Println(&vs)
}
