package config

// TypeConfig enum for various configurations
type TypeConfig int

const (
	enclave     = "Enclave"
	host        = "Host"
	network     = "Network"
	node        = "Node"
	eth2Network = "Eth2Network"
	l1Deployer  = "L1Deployer"
	l2Deployer  = "L2Deployer"
)

const (
	Enclave TypeConfig = iota
	Host
	Network
	Node
	Eth2Network
	L1Deployer
	L2Deployer
)

func (t TypeConfig) String() string {
	return [...]string{
		enclave,
		host,
		network,
		node,
		eth2Network,
		l1Deployer,
		l2Deployer,
	}[t]
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
	case eth2Network:
		return Eth2Network, nil
	case l1Deployer:
		return L1Deployer, nil
	case l2Deployer:
		return L2Deployer, nil
	default:
		panic("string " + s + " cannot be converted to TypeConfig.")
	}
}
