package networkmanager

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	obscuroconfig "github.com/obscuronet/obscuro-playground/go/obscuronode/config"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
	"github.com/obscuronet/obscuro-playground/integration/erc20contract"
	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
)

var (
	mgmtContractBytes  = common.Hex2Bytes(mgmtcontractlib.MgmtContractByteCode)
	erc20ContractBytes = common.Hex2Bytes(erc20contract.ContractByteCode)
)

// DeployContract deploys a management contract or ERC20 contract to the L1 network, and prints its address.
func DeployContract(config Config) {
	var contractBytes []byte
	switch config.Command {
	case DeployMgmtContract:
		contractBytes = mgmtContractBytes
	case DeployERC20Contract:
		contractBytes = erc20ContractBytes
	default:
		panic("unrecognised command type")
	}

	hostConfig := obscuroconfig.HostConfig{
		L1NodeHost:          config.l1NodeHost,
		L1NodeWebsocketPort: config.l1NodeWebsocketPort,
		L1ConnectionTimeout: config.l1ConnectionTimeout,
		PrivateKeyString:    config.privateKeyString,
		ChainID:             config.chainID,
	}
	l1Client, err := ethclient.NewEthClient(hostConfig)
	if err != nil {
		panic(err)
	}
	l1Wallet := wallet.NewInMemoryWalletFromString(hostConfig)

	var contractAddress *common.Address
	nonce := l1Wallet.GetNonceAndIncrement()
	contractAddress, err = network.DeployContract(l1Client, l1Wallet, contractBytes, nonce)
	for err != nil {
		// If the error isn't a nonce-too-low error, we report it as a legitimate error.
		if err.Error() != core.ErrNonceTooLow.Error() {
			panic(fmt.Errorf("contract deployment failed. Cause: %w", err))
		}
		// TODO - Smarter approach to finding correct nonce.
		// We loop until we have reached the required nonce.
		nonce++
		contractAddress, err = network.DeployContract(l1Client, l1Wallet, contractBytes, nonce)
	}
	println(contractAddress.Hex())
	os.Exit(0)
}
