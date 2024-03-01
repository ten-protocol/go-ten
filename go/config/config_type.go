package config

import (
	"fmt"
)

const (
	host     = "host"
	enclave  = "enclave"
	edgeless = "edgeless"
	unknown  = "unknown"
)

// ConfType represents the type of the configuration.
type ConfType int

const (
	Host ConfType = iota
	Enclave
	Edgeless
	Unknown
)

func (c ConfType) String() string {
	switch c {
	case Host:
		return host
	case Enclave:
		return enclave
	case Edgeless:
		return edgeless
	case Unknown:
		return unknown
	default:
		return unknown
	}
}

func ToConfType(s string) (ConfType, error) {
	switch s {
	case host:
		return Host, nil
	case enclave:
		return Enclave, nil
	case edgeless:
		return Edgeless, nil
	default:
		return Unknown, fmt.Errorf("string '%s' cannot be converted to a node type", s)
	}
}

//func (c ConfType) configFromFlags(flags map[string]*flag.TenFlag) (interface{}, error) {
//	switch c {
//	case Host:
//		nodeType, err := common.ToNodeType(flags[NodeTypeFlag].String())
//		if err != nil {
//			return nil, fmt.Errorf("unrecognised node type '%s'", flags[NodeTypeFlag].String())
//		}
//
//		cfg := &HostConfig{}
//		cfg.IsGenesis = flags[isGenesisFlag].Bool()
//		cfg.NodeType = nodeType
//		cfg.
//
//		return "host"
//	case Enclave:
//		cfg := &EnclaveConfig{}
//		return "enclave"
//	case Edgeless:
//		//return "edgeless"
//	case Unknown:
//		//return "unknown"
//	default:
//		return "unknown"
//	}
//}
