package contractlib

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/enclave/crosschain"
	"github.com/ten-protocol/go-ten/go/ethadapter"
)

type NetworkConfigLib struct {
	addr        gethcommon.Address
	ethClient   ethadapter.EthClient
	contractABI abi.ABI
}

type NetworkAddresses struct {
	CrossChain             gethcommon.Address
	MessageBus             gethcommon.Address
	NetworkEnclaveRegistry gethcommon.Address
	RollupContract         gethcommon.Address
}

func NewNetworkConfigLib(address gethcommon.Address, ethClient ethadapter.EthClient) (*NetworkConfigLib, error) {
	return &NetworkConfigLib{
		addr:        address,
		ethClient:   ethClient,
		contractABI: crosschain.NetworkConfigABI,
	}, nil
}

func (nc *NetworkConfigLib) GetContractAddresses() (*NetworkAddresses, error) {
	data, err := nc.contractABI.Pack("addresses")
	if err != nil {
		return nil, fmt.Errorf("failed to encode addresses() call: %w", err)
	}

	result, err := nc.ethClient.CallContract(ethereum.CallMsg{
		To:   &nc.addr,
		Data: data,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call addresses(): %w", err)
	}

	addresses := new(NetworkAddresses)
	err = nc.contractABI.UnpackIntoInterface(addresses, "addresses", result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode addresses response: %w", err)
	}

	return addresses, nil
}

func (nc *NetworkConfigLib) GetContractAddr() *gethcommon.Address {
	return &nc.addr
}
