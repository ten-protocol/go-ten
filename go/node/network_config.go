package node

import (
	"encoding/json"
	"github.com/ten-protocol/go-ten/go/config"
	"os"
	"path"
)

// This is the location where the metadata will be stored
const _networkCfgFile = ".ten-network.json"

func WriteNetworkConfigToDisk(cfg config.NetworkConfig) error {
	n := cfg.GetNetwork()
	jsonStr, err := json.Marshal(n)
	if err != nil {
		return err
	}

	// store in the user home dir
	filePath, err := tenFilePath()
	if err != nil {
		return err
	}

	// create the file as read-only, expect it to be immutable data for the lifetime of the obscuro network for the node
	err = os.WriteFile(filePath, jsonStr, 0o644) //nolint:gosec
	if err != nil {
		return err
	}
	return nil
}

func ReadNetworkConfigFromDisk() (*config.NetworkInputConfig, error) {
	// store in the user home dir
	filePath, err := tenFilePath()
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var cfg config.NetworkInputConfig
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func tenFilePath() (string, error) {
	// store in the user home dir
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(dirname, _networkCfgFile), nil
}
