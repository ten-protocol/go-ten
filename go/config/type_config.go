package config

// TypeConfig enum for various configurations
type TypeConfig int

const (
	enclave = "Enclave"
	host    = "Host"
	network = "Network"
	node    = "Node"
	shared  = "Shared"
)

const (
	Enclave TypeConfig = iota
	Host
	Network
	Node
	Shared
)

func (t TypeConfig) String() string {
	return [...]string{enclave, host, network, node, shared}[t]
}

func ToTypeConfig(s string) (TypeConfig, error) {
	switch s {
	case enclave:
		return Enclave, nil
	case host:
		return Host, nil
	case network:
		return Network, nil
	case node:
		return Node, nil
	case shared:
		return Shared, nil
	default:
		panic("string " + s + " cannot be converted to TypeConfig.")
	}
}
