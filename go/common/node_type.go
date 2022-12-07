package common

import "fmt"

const (
	sequencer = "sequencer"
	validator = "validator"
	unknown   = "unknown"
)

// NodeType represents the type of the node.
type NodeType int

const (
	Sequencer NodeType = iota
	Validator
	Unknown
)

func (n NodeType) String() string {
	switch n {
	case Sequencer:
		return sequencer
	case Validator:
		return validator
	case Unknown:
		return unknown
	default:
		return unknown
	}
}

func ToNodeType(s string) (NodeType, error) {
	switch s {
	case sequencer:
		return Sequencer, nil
	case validator:
		return Validator, nil
	default:
		return Unknown, fmt.Errorf("string '%s' cannot be converted to a node type", s)
	}
}
