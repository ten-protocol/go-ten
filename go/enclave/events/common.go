package events

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/contracts/generated/NetworkConfig"
)

func ConvertLogsToNetworkUpgrades(logs []types.Log, eventName string, networkConfigABI abi.ABI) ([]NetworkConfig.NetworkConfigUpgraded, error) {
	networkUpgrades := make([]NetworkConfig.NetworkConfigUpgraded, 0)
	for _, log := range logs {
		var event NetworkConfig.NetworkConfigUpgraded
		err := networkConfigABI.UnpackIntoInterface(&event, eventName, log.Data)
		if err != nil {
			return nil, err
		}

		networkUpgrades = append(networkUpgrades, event)
	}

	return networkUpgrades, nil
}
