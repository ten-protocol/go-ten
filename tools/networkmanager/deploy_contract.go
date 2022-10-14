package networkmanager

import (
	"os"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/contracts/managementcontract"
	"github.com/obscuronet/go-obscuro/integration/erc20contract"

	"github.com/ethereum/go-ethereum/common"
	obscuroconfig "github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"
)

// DeployContract deploys a management contract or ERC20 contract to the L1 network, and prints its address.
func DeployContract(config Config, logger gethlog.Logger) {
	var contractBytes []byte
	switch config.Command { //nolint:exhaustive
	case DeployMgmtContract:
		bytecode, err := managementcontract.Bytecode()
		if err != nil {
			panic(err)
		}
		contractBytes = bytecode
	case DeployERC20Contract:
		contractBytes = erc20contract.L1BytecodeWithDefaultSupply(config.erc20Token)
	default:
		panic("unrecognised command type")
	}

	hostConfig := obscuroconfig.HostConfig{
		PrivateKeyString: config.privateKeys[0], // We deploy the contract using the first private key.
		L1ChainID:        config.l1ChainID,
	}

	l1Client, err := ethadapter.NewEthClient(config.l1NodeHost, config.l1NodeWebsocketPort, config.l1RPCTimeout, common.HexToAddress("0x0"), logger)
	if err != nil {
		panic(err)
	}

	l1Wallet := wallet.NewInMemoryWalletFromConfig(hostConfig, logger)
	nonce, err := l1Client.Nonce(l1Wallet.Address())
	if err != nil {
		panic(err)
	}
	l1Wallet.SetNonce(nonce)

	var contractAddress *common.Address
	contractAddress, err = network.DeployContract(l1Client, l1Wallet, contractBytes)
	if err != nil {
		panic(err)
	}

	println(contractAddress.Hex())
	os.Exit(0)
}
