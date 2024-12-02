package common

import "fmt"

const (
	activeSequencer = "sequencer"
	validator       = "validator"
	backupSequencer = "backup_sequencer"
	unknown         = "unknown"
)

// NodeType represents the type of the node.
type NodeType int

const (
	ActiveSequencer NodeType = iota
	Validator
	BackupSequencer
	Unknown
)

func (n NodeType) String() string {
	switch n {
	case ActiveSequencer:
		return activeSequencer
	case Validator:
		return validator
	case BackupSequencer:
		return backupSequencer
	case Unknown:
		return unknown
	default:
		return unknown
	}
}

func (n *NodeType) UnmarshalText(text []byte) error {
	nodeType, err := ToNodeType(string(text))
	if err != nil {
		return err
	}
	*n = nodeType
	return nil
}

func ToNodeType(s string) (NodeType, error) {
	switch s {
	case activeSequencer:
		return ActiveSequencer, nil
	case validator:
		return Validator, nil
	case backupSequencer:
		return BackupSequencer, nil
	default:
		return Unknown, fmt.Errorf("string '%s' cannot be converted to a node type", s)
	}
}
