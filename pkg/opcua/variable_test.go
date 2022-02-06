package opcua

import (
	"testing"

	"github.com/gopcua/opcua/ua"
)

func variableEqual(v1, v2 Variable) bool {
	switch {
	case v1.name != v2.name:
		return false
	case v1.variable != v2.variable:
		return false
	case (v1.nodeID == nil) && (v2.nodeID == nil):
		return true
	case v1.nodeID.NodeID.String() != v2.nodeID.NodeID.String():
		return false
	default:
		return true
	}
}

func readID(node string) *ua.ReadValueID {
	id, _ := ua.ParseNodeID(node)
	return &ua.ReadValueID{NodeID: id}
}

func TestNewVariables(t *testing.T) {
	type args struct {
		nodes []string
	}

	tests := []struct {
		name  string
		args  args
		want  *Variables
		equal bool
	}{
		{"Nil", args{[]string{}}, &Variables{}, true},
		{"Empty", args{[]string{""}}, &Variables{}, true},
		{"SingleLeadingWhitespace", args{[]string{" ns=4;s=OPC_Dummy_Var_1"}}, &Variables{}, true},
		{"Single", args{[]string{"ns=4;s=OPC_Dummy_Var_1"}}, &Variables{Variable{"ns=4;s=OPC_Dummy_Var_1", nil, readID("ns=4;s=OPC_Dummy_Var_1")}}, true},
		{"Double", args{[]string{"ns=4;s=OPC_Dummy_Var_1", "ns=4;s=OPC_Dummy_Var_2"}}, &Variables{Variable{"ns=4;s=OPC_Dummy_Var_1", nil, readID("ns=4;s=OPC_Dummy_Var_1")}, Variable{"ns=4;s=OPC_Dummy_Var_2", nil, readID("ns=4;s=OPC_Dummy_Var_2")}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewVariables(tt.args.nodes)
			if len(*got) != len(*tt.want) {
				t.Errorf("NewVariables() = %v, want %v", got, tt.want)
			}
			for i, n := range *got {
				if !variableEqual(n, (*tt.want)[i]) {
					t.Errorf("NewVariables() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestNewVariable(t *testing.T) {
	type args struct {
		node string
	}

	tests := []struct {
		name    string
		args    args
		want    *Variable
		wantErr bool
	}{
		{"Nil", args{""}, nil, true},
		{"InvalidSimple", args{"test"}, nil, true},
		{"InvalidNoNode", args{"ns=;s=OPC_Dummy_Var_2"}, nil, true},
		{"InvalidNoVar", args{"ns=4;s="}, nil, true},
		{"InvalidLeadingWhitespace", args{" ns=4;s=OPC_Dummy_Var_2"}, nil, true},
		{"ValidNodeSmall", args{"ns=4;s=OPC_Dummy_Var_2"}, &Variable{"ns=4;s=OPC_Dummy_Var_2", nil, readID("ns=4;s=OPC_Dummy_Var_2")}, false},
		{"ValidNodeLarge", args{"ns=456;s=OPC_Dummy_Var_2"}, &Variable{"ns=456;s=OPC_Dummy_Var_2", nil, readID("ns=456;s=OPC_Dummy_Var_2")}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewVariable(tt.args.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVariable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil {
				if variableEqual(*got, *tt.want) != true {
					t.Errorf("NewVariable() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestVariables_NodeIDs(t *testing.T) {
	tests := []struct {
		name string
		vs   Variables
	}{
		{"Nil", Variables{}},
		{"Len1", Variables{Variable{"test", nil, nil}}},
		{"Len2", Variables{Variable{"test", nil, nil}, Variable{"test2", nil, nil}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vs.NodeIDs(); len(tt.vs) != len(got) {
				t.Errorf("NodeIDs() = %v, want %v", len(tt.vs), len(got))
			}
		})
	}
}

func TestVariables_String(t *testing.T) {
	tests := []struct {
		name string
		vs   Variables
		want string
	}{
		{"Nil", Variables{}, ""},
		{"Len1", Variables{Variable{name: "test", variable: 0}}, "test: 0\n"},
		{"DuplicateName", Variables{Variable{name: "test", variable: 0}, Variable{name: "test", variable: 5}}, "test: 0\ntest: 5\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vs.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
