package opcua

import (
	"context"
	"fmt"
	"io"

	"github.com/gopcua/opcua/ua"

	"github.com/gopcua/opcua"
)

type ProtocolError struct {
	ErrorString string
}

func (pe *ProtocolError) Error() string { return pe.ErrorString }

var (
	ErrOpcuaRead    = &ProtocolError{"unable to read variables"}
	ErrOpcuaConnect = &ProtocolError{"unable to connect"}
	ErrOpcuaStatus  = &ProtocolError{"invalid status"}
)

type Client interface {
	Connect(context.Context) error
	Read(*ua.ReadRequest) (*ua.ReadResponse, error)
	io.Closer
}

type Opcua struct {
	OpcuaClient Client
}

func NewOpcua(endpoint string) *Opcua {
	return &Opcua{opcua.NewClient(endpoint)}
}

func NewOpcuaWithClient(client Client) *Opcua {
	return &Opcua{client}
}

func (o *Opcua) CheckConnection(ctx context.Context) error {
	if err := o.OpcuaClient.Connect(ctx); err != nil {
		return fmt.Errorf("opcua err: %v base err: %w", ErrOpcuaConnect, err)
	}
	defer o.OpcuaClient.Close()

	return nil
}

func (o *Opcua) GetVariables(ctx context.Context, vars Variables) (Variables, error) {
	if err := o.OpcuaClient.Connect(ctx); err != nil {
		return nil, fmt.Errorf("opcua err: %v base err: %w", ErrOpcuaConnect, err)
	}
	defer o.OpcuaClient.Close()

	req := &ua.ReadRequest{
		MaxAge:             2000,
		NodesToRead:        vars.NodeIDs(),
		TimestampsToReturn: ua.TimestampsToReturnBoth,
	}

	// TODO: Split up if err != nil into cleaner structure
	resp, err := o.OpcuaClient.Read(req)
	if err != nil || resp == nil || resp.Results == nil || len(resp.Results) == 0 {
		return nil, fmt.Errorf("opcua err: %v base err: %w", ErrOpcuaRead, err)
	}

	if resp.Results[0].Status != ua.StatusOK {
		return nil, fmt.Errorf("opcua err: %v status not OK: %s", ErrOpcuaStatus, resp.Results[0].Status)
	}

	for i, s := range resp.Results {
		vars[i].variable = s.Value.Value()
	}

	return vars, nil
}
