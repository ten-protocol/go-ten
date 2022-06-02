package smartcontract

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
)

type debugWallet struct {
	wallet.Wallet
}

func newDebugWallet(w wallet.Wallet) *debugWallet {
	return &debugWallet{w}
}

func (w *debugWallet) AwaitedSignAndSendTransaction(client ethclient.EthClient, txData types.TxData) (*types.Transaction, *types.Receipt, error) {
	signedTx, err := w.SignAndSendTransaction(client, txData)
	if err != nil {
		return nil, nil, err
	}
	receipt, err := waitTxResult(client, signedTx)
	if err != nil {
		return nil, nil, err
	}
	return signedTx, receipt, nil

}

func (w *debugWallet) SignAndSendTransaction(client ethclient.EthClient, txData types.TxData) (*types.Transaction, error) {
	signedTx, err := w.SignTransaction(txData)
	if err != nil {
		return nil, err
	}

	err = client.SendTransaction(signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func waitTxResult(client ethclient.EthClient, tx *types.Transaction) (*types.Receipt, error) {
	var receipt *types.Receipt
	var err error
	for start := time.Now(); time.Since(start) < 30*time.Second; time.Sleep(time.Second) {
		receipt, err = client.TransactionReceipt(tx.Hash())

		if err != nil {
			if err == ethereum.NotFound {
				continue
			}
			return nil, err
		}

		return receipt, nil
	}
	return nil, fmt.Errorf("transaction not minted after timeout")
}
