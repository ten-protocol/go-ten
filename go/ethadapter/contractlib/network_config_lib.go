package contractlib

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ten-protocol/go-ten/contracts/generated/NetworkConfig"
	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/ethadapter"
)

type NetworkConfigLib interface {
	GetContractAddr() *gethcommon.Address
	GetContractAddresses() (*common.NetworkConfigAddresses, error)
	IsMock() bool
}

type networkConfigLibImpl struct {
	addr        gethcommon.Address
	ethClient   ethclient.Client
	contractABI abi.ABI
}

func NewNetworkConfigLib(address gethcommon.Address, ethClient ethclient.Client) (NetworkConfigLib, error) {
	return &networkConfigLibImpl{
		addr:        address,
		ethClient:   ethClient,
		contractABI: ethadapter.NetworkConfigABI,
	}, nil
}

func (nc *networkConfigLibImpl) GetContractAddresses() (*common.NetworkConfigAddresses, error) {
	networkConfigContract, err := NetworkConfig.NewNetworkConfigCaller(nc.addr, &nc.ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create NetworkConfig caller: %w", err)
	}

	// Call the Addresses method using the generated binding
	addresses, err := networkConfigContract.Addresses(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to call addresses(): %w", err)
	}

	// Convert to your NetworkAddresses type
	return &common.NetworkConfigAddresses{
		CrossChain:             addresses.CrossChain,
		MessageBus:             addresses.MessageBus,
		NetworkEnclaveRegistry: addresses.NetworkEnclaveRegistry,
		RollupContract:         addresses.RollupContract,
	}, nil
}

func (nc *networkConfigLibImpl) GetContractAddr() *gethcommon.Address {
	return &nc.addr
}

func (nc *networkConfigLibImpl) IsMock() bool {
	return false
}
