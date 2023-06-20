package node

import (
	"encoding/json"
	"os"
	"path"
)

// This is the location where the metadata will be stored
const _networkCfgFile = ".obscuro-network.json"

// NetworkConfig is key network information required to start a node connecting to that network.
// We persist it as a json file on our testnet hosts so that they can read it off when restart/upgrading
type NetworkConfig struct {
	ManagementContractAddress string
	MessageBusAddress         string
	L1StartHash               string // L1 block hash from which to process for L2 data (mgmt contract deploy block)
}

func WriteNetworkConfigToDisk(cfg *Config) error {
	n := NetworkConfig{
		ManagementContractAddress: cfg.managementContractAddr,
		MessageBusAddress:         cfg.messageBusContractAddress,
		L1StartHash:               cfg.l1Start,
	}
	jsonStr, err := json.Marshal(n)
	if err != nil {
		return err
	}

	// store in the user home dir
	dirname, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// create the file as read-only, expect it to be immutable data for the lifetime of the obscuro network for the node
	err = os.WriteFile(path.Join(dirname, _networkCfgFile), jsonStr, 0o644) //nolint:gosec
	if err != nil {
		return err
	}
	return nil
}

func ReadNetworkConfigFromDisk() (*NetworkConfig, error) {
	bytes, err := os.ReadFile(_networkCfgFile)
	if err != nil {
		return nil, err
	}
	var cfg NetworkConfig
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
