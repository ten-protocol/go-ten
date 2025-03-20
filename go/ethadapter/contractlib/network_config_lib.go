package contractlib

import (
	"fmt"

	"github.com/ethereum/go-ethereum"

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
	AddAdditionalAddress(name string, address gethcommon.Address) (ethereum.CallMsg, error)
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

	addresses, err := networkConfigContract.Addresses(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to call addresses(): %w", err)
	}

	return &common.NetworkConfigAddresses{
		CrossChain:            addresses.CrossChain,
		EnclaveRegistry:       addresses.NetworkEnclaveRegistry,
		RollupContract:        addresses.RollupContract,
		L1MessageBus:          addresses.MessageBus,
		L1Bridge:              addresses.L1Bridge,
		L2Bridge:              addresses.L2Bridge,
		L1CrossChainMessenger: addresses.L1crossChainMessenger,
		L2CrossChainMessenger: addresses.L2crossChainMessenger,
	}, nil
}

func (nc *networkConfigLibImpl) GetContractAddr() *gethcommon.Address {
	return &nc.addr
}

func (nc *networkConfigLibImpl) AddAdditionalAddress(name string, address gethcommon.Address) (ethereum.CallMsg, error) {
	data, err := nc.contractABI.Pack(ethadapter.AddAdditionalAddressMethod, name, address)
	if err != nil {
		return ethereum.CallMsg{}, fmt.Errorf("could not pack the call data. Cause: %w", err)
	}
	return ethereum.CallMsg{To: &nc.addr, Data: data}, nil
}

func (nc *networkConfigLibImpl) IsMock() bool {
	return false
}
