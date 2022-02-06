package opcua

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/gopcua/opcua/ua"
)

var (
	ErrInvalidNode = fmt.Errorf("invalid node format")
)

type Variable struct {
	name     string
	variable interface{}
	nodeID   *ua.ReadValueID
}

type Variables []Variable

func NewVariable(node string) (*Variable, error) {
	if err := nodeValidate(node); err != nil {
		return nil, err
	}

	id, _ := ua.ParseNodeID(node)

	return &Variable{node, nil, &ua.ReadValueID{NodeID: id}}, nil
}

func (v Variable) String() string {
	return fmt.Sprintf("%s: %v", v.name, v.variable)
}

func (v Variable) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(&struct {
		Key   string
		Value interface{}
	}{
		Key:   v.name,
		Value: v.variable,
	})

	if err != nil {
		return nil, fmt.Errorf("error unmarshalling json err: %w", err)
	}
	return j, nil
}

func NewVariables(nodes []string) *Variables {
	vs := Variables{}

	for _, n := range nodes {
		v, err := NewVariable(n)
		if err != nil {
			fmt.Printf("error: %v var: %v\n", err, n)
			continue
		}

		vs = append(vs, *v)
	}

	return &vs
}

func (vs *Variables) NodeIDs() []*ua.ReadValueID {
	var rv []*ua.ReadValueID
	for _, v := range *vs {
		rv = append(rv, v.nodeID)
	}

	return rv
}

func (vs *Variables) String() string {
	var s string
	for _, v := range *vs {
		s += fmt.Sprintf("%v\n", v)
	}

	return s
}

func nodeValidate(node string) error {
	res, _ := regexp.MatchString("^ns=[0-9]+;s=(.)+", node)

	if !res {
		return ErrInvalidNode
	}

	return nil
}
