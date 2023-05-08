package node

type Node interface {
	Start() error
	Stop() error
	Upgrade(networkCfg *NetworkConfig) error
}
