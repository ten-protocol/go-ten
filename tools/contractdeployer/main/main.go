package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	obscuroconfig "github.com/obscuronet/obscuro-playground/go/obscuronode/config"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
	"github.com/obscuronet/obscuro-playground/integration/erc20contract"
	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
)

func main() {
	config := parseCLIArgs()

	var contractBytes []byte
	switch config.contractType {
	case management:
		contractBytes = common.Hex2Bytes(mgmtcontractlib.MgmtContractByteCode)
	case erc20:
		contractBytes = common.Hex2Bytes(erc20contract.ContractByteCode)
	default:
		panic(fmt.Sprintf("unrecognised contract type. Expected either %s or %s", managementName, erc20Name))
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

	mgmtContractAddress := network.DeployContract(l1Client, l1Wallet, contractBytes)
	println(mgmtContractAddress.Hex())
	os.Exit(0)
}
