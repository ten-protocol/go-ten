package network

import (
	"math/big"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/log"
)

func createEthClientConnection(id int64, port uint) ethclient.EthClient {
	hostConfig := config.HostConfig{
		ID:                  common.BigToAddress(big.NewInt(id)),
		L1NodeHost:          Localhost,
		L1NodeWebsocketPort: port,
		L1ConnectionTimeout: DefaultL1ConnectionTimeout,
	}
	ethnode, err := ethclient.NewEthClient(hostConfig)
	if err != nil {
		panic(err)
	}
	return ethnode
}

func deployContract(workerClient ethclient.EthClient, w wallet.Wallet, contractBytes []byte) *common.Address {
	deployContractTx := types.LegacyTx{
		Nonce:    w.GetNonceAndIncrement(),
		GasPrice: big.NewInt(2000000000),
		Gas:      1025_000_000,
		Data:     contractBytes,
	}

	signedTx, err := w.SignTransaction(&deployContractTx)
	if err != nil {
		panic(err)
	}

	err = workerClient.SendTransaction(signedTx)
	if err != nil {
		panic(err)
	}

	var receipt *types.Receipt
	for start := time.Now(); time.Since(start) < 80*time.Second; time.Sleep(2 * time.Second) {
		receipt, err = workerClient.TransactionReceipt(signedTx.Hash())
		if err == nil && receipt != nil {
			if receipt.Status != types.ReceiptStatusSuccessful {
				panic("unable to deploy contract")
			}
			break
		}

		log.Info("Contract deploy tx has not been mined into a block after %s...", time.Since(start))
	}

	log.Info("Contract successfully deployed to %s", receipt.ContractAddress)
	return &receipt.ContractAddress
}
