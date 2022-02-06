package opcua

import (
	"context"
	"fmt"
	"testing"

	"github.com/gopcua/opcua/ua"
)

type MockOpcuaClient struct {
}

func (m *MockOpcuaClient) Connect(ctx context.Context) error {
	return nil
}

func (m *MockOpcuaClient) Read(request *ua.ReadRequest) (*ua.ReadResponse, error) {
	return nil, nil
}

func (m *MockOpcuaClient) Close() error {
	return nil
}

type MockOpcuaClientConnErr struct {
	MockOpcuaClient
}

func (m *MockOpcuaClientConnErr) Connect(ctx context.Context) error {
	return fmt.Errorf("connection error")
}

type MockOpcuaClientStatusErr struct {
	MockOpcuaClient
}

func (m *MockOpcuaClientStatusErr) Read(request *ua.ReadRequest) (*ua.ReadResponse, error) {
	r := ua.ReadResponse{
		Results: []*ua.DataValue{{Status: ua.StatusBad}},
	}

	return &r, nil
}

type MockOpcuaClientStatusOK struct {
	MockOpcuaClient
}

func (m *MockOpcuaClientStatusOK) Read(request *ua.ReadRequest) (*ua.ReadResponse, error) {
	var d []*ua.DataValue

	v, _ := ua.NewVariant(5.0)

	for range request.NodesToRead {
		d = append(d, &ua.DataValue{Value: v, Status: ua.StatusOK})
	}

	r := ua.ReadResponse{
		Results: d,
	}

	return &r, nil
}

func TestNewOpcua(t *testing.T) {
	type args struct {
		endpoint string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"NilEndpoint", args{""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOpcua(tt.args.endpoint); (got.OpcuaClient == Client(nil)) == tt.want {
				t.Errorf("NewOpcua() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOpcuaWithClient(t *testing.T) {
	type args struct {
		client Client
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"ValidMockOpcuaClient", args{&MockOpcuaClient{}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOpcuaWithClient(tt.args.client); (got.OpcuaClient == Client(nil)) == tt.want {
				t.Errorf("NewOpcuaWithClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpcua_GetVariables(t *testing.T) {
	type fields struct {
		OpcuaClient Client
	}

	type args struct {
		ctx  context.Context
		vars Variables
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Variables
		wantErr bool
	}{
		{"Nil", fields{&MockOpcuaClient{}}, args{context.Background(), Variables{}}, Variables{}, true},
		{"InvalidConnError", fields{&MockOpcuaClientConnErr{}}, args{context.Background(), Variables{}}, Variables{}, true},
		{"InvalidReqStatusBad", fields{&MockOpcuaClientStatusErr{}}, args{context.Background(), Variables{}}, Variables{}, true},
		{"InvalidReqStatusOKNilVar", fields{&MockOpcuaClientStatusOK{}}, args{context.Background(), Variables{}}, Variables{}, true},
		{"ValidReqStatusOK", fields{&MockOpcuaClientStatusOK{}}, args{context.Background(), Variables{Variable{"test", nil, nil}}}, Variables{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Opcua{
				OpcuaClient: tt.fields.OpcuaClient,
			}
			_, err := o.GetVariables(tt.args.ctx, tt.args.vars)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVariables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestOpcua_CheckConnection(t *testing.T) {
	type fields struct {
		OpcuaClient Client
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"InvalidConnError", fields{&MockOpcuaClientConnErr{}}, args{context.Background()}, true},
		{"ValidConn", fields{&MockOpcuaClientStatusOK{}}, args{context.Background()}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Opcua{
				OpcuaClient: tt.fields.OpcuaClient,
			}
			if err := o.CheckConnection(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("CheckConnection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
