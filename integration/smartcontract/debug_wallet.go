package smartcontract

import (
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethadapter"
	"github.com/obscuronet/obscuro-playground/go/wallet"
)

var _timeout = 30 * time.Second

// debugWallet is a wrapper around the wallet that simplifies commonly used functions
type debugWallet struct {
	wallet.Wallet
}

// newDebugWallet returns a new debug wrapped wallet
func newDebugWallet(w wallet.Wallet) *debugWallet {
	return &debugWallet{w}
}

// AwaitedSignAndSendTransaction signs a tx, issues the tx and awaits the tx to be minted into a block
func (w *debugWallet) AwaitedSignAndSendTransaction(client ethadapter.EthClient, txData types.TxData) (*types.Transaction, *types.Receipt, error) {
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

// SignAndSendTransaction signs and sends a tx
func (w *debugWallet) SignAndSendTransaction(client ethadapter.EthClient, txData types.TxData) (*types.Transaction, error) {
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

// waitTxResult waits for a tx to be minted into a block
func waitTxResult(client ethadapter.EthClient, tx *types.Transaction) (*types.Receipt, error) {
	var receipt *types.Receipt
	var err error
	for start := time.Now(); time.Since(start) < _timeout; time.Sleep(time.Second) {
		receipt, err = client.TransactionReceipt(tx.Hash())

		if err != nil {
			if errors.Is(err, ethereum.NotFound) {
				continue
			}
			return nil, err
		}

		return receipt, nil
	}
	return nil, fmt.Errorf("transaction not minted after timeout")
}
