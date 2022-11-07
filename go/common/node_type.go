package common

import "fmt"

const (
	aggregator = "aggregator"
	validator  = "validator"
	unknown    = "unknown"
)

// NodeType represents the type of the node.
type NodeType int

const (
	Aggregator NodeType = iota
	Validator
	Unknown
)

func (n NodeType) String() string {
	switch n {
	case Aggregator:
		return aggregator
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
	case aggregator:
		return Aggregator, nil
	case validator:
		return Validator, nil
	default:
		return Unknown, fmt.Errorf("string '%s' cannot be converted to a node type", s)
	}
}
