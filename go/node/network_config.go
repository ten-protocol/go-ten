package node

import (
	"encoding/json"
	"io/ioutil"
)

// This is the location where the metadata will be stored on all testnet VMs
const _networkCfgFilePath = "/home/obscuro/network.json"

// networkConfig is key network information required to start a node connecting to that network.
// We persist it as a json file on our testnet hosts so that they can read it off when restart/upgrading
type networkConfig struct {
	managementContractAddress string
	messageBusAddress         string
}

func WriteNetworkConfigToDisk(cfg networkConfig) error {
	jsonStr, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_networkCfgFilePath, jsonStr, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ReadNetworkConfigFromDisk() (*networkConfig, error) {
	bytes, err := ioutil.ReadFile(_networkCfgFilePath)
	if err != nil {
		return nil, err
	}
	var cfg networkConfig
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
