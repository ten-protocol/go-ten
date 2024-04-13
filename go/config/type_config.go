package config

// TypeConfig enum for various configurations
type TypeConfig int

const (
	enclave = "Enclave"
	host    = "Host"
	network = "Network"
	node    = "Node"
)

const (
	Enclave TypeConfig = iota
	Host
	Network
	Node
)

func (t TypeConfig) String() string {
	return [...]string{"Enclave", "Host", "Network"}[t]
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
	default:
		panic("string " + s + " cannot be converted to TypeConfig.")
	}
}
