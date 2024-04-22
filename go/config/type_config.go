package config

// TypeConfig enum for various configurations
type TypeConfig int

const (
	enclave     = "Enclave"
	host        = "Host"
	network     = "Network"
	node        = "Node"
	testnet     = "Testnet"
	eth2Network = "Eth2Network"
	l1Deployer  = "L1Deployer"
	l2Deployer  = "L2Deployer"
	faucet      = "Faucet"
)

const (
	Enclave TypeConfig = iota
	Host
	Network
	Node
	Testnet
	Eth2Network
	L1Deployer
	L2Deployer
	Faucet
)

func (t TypeConfig) String() string {
	return [...]string{
		enclave,
		host,
		network,
		node,
		testnet,
		eth2Network,
		l1Deployer,
		l2Deployer,
		faucet,
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
	case testnet:
		return Testnet, nil
	case eth2Network:
		return Eth2Network, nil
	case l1Deployer:
		return L1Deployer, nil
	case l2Deployer:
		return L2Deployer, nil
	case faucet:
		return Faucet, nil
	default:
		panic("string " + s + " cannot be converted to TypeConfig.")
	}
}
