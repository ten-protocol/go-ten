package contractlib

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/ethadapter"
)

type NetworkConfigLib interface {
	GetContractAddr() *gethcommon.Address
	GetContractAddresses() (*common.NetworkAddresses, error)
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

func (nc *networkConfigLibImpl) GetContractAddresses() (*common.NetworkAddresses, error) {
	data, err := nc.contractABI.Pack("addresses")
	if err != nil {
		return nil, fmt.Errorf("failed to encode addresses() call: %w", err)
	}

	result, err := nc.ethClient.CallContract(context.Background(), ethereum.CallMsg{
		To:   &nc.addr,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call addresses(): %w", err)
	}

	addresses := new(common.NetworkAddresses)
	err = nc.contractABI.UnpackIntoInterface(addresses, "addresses", result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode addresses response: %w", err)
	}

	return addresses, nil
}

func (nc *networkConfigLibImpl) GetContractAddr() *gethcommon.Address {
	return &nc.addr
}
