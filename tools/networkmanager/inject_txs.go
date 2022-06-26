package networkmanager

import (
	"fmt"
	"time"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
	"github.com/obscuronet/obscuro-playground/integration/simulation"
)

func InjectTransactions(nmConfig Config) {
	hostConfig := config.HostConfig{
		ID:                  common.HexToAddress(""),
		L1NodeHost:          nmConfig.l1NodeHost,
		L1NodeWebsocketPort: nmConfig.l1NodeWebsocketPort,
		L1ConnectionTimeout: nmConfig.l1ConnectionTimeout,
		PrivateKeyString:    nmConfig.privateKeyString,
		ChainID:             nmConfig.chainID,
	}

	// TODO - Consider extending this command to support multiple L1 clients and L2 clients.
	l1Client, err := ethclient.NewEthClientFromConfig(hostConfig)
	if err != nil {
		panic(fmt.Sprintf("could not create L1 client. Cause: %s", err))
	}
	l2Client := obscuroclient.NewClient(nmConfig.obscuroClientAddress)

	l1Wallet := wallet.NewInMemoryWalletFromConfig(hostConfig)
	nonce, err := l1Client.Nonce(l1Wallet.Address())
	if err != nil {
		panic(err)
	}
	l1Wallet.SetNonce(nonce)

	txInjector := simulation.NewTransactionInjector(
		1*time.Second,
		stats.NewStats(1),
		[]ethclient.EthClient{l1Client},
		&params.SimWallets{
			// todo
		},
		&nmConfig.mgmtContractAddress,
		[]obscuroclient.Client{l2Client},
		mgmtcontractlib.NewMgmtContractLib(&nmConfig.mgmtContractAddress),
		erc20contractlib.NewERC20ContractLib(&nmConfig.mgmtContractAddress, &nmConfig.erc20ContractAddress),
	)

	println("Injecting transactions into network...")
	txInjector.Start()
}
