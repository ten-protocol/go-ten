package smartcontract

import (
	"context"
	"fmt"
	"time"

	"github.com/ten-protocol/go-ten/integration/common/testlog"

	"github.com/ethereum/go-ethereum"
	"github.com/ten-protocol/go-ten/go/common/retry"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/wallet"
)

// debugWallet is a wrapper around the wallet that simplifies commonly used functions
type debugWallet struct {
	wallet.Wallet
	receiptTimeout time.Duration
}

func newDebugWallet(w wallet.Wallet, timeout time.Duration) *debugWallet {
	return &debugWallet{
		Wallet:         w,
		receiptTimeout: timeout,
	}
}

// AwaitedSignAndSendTransaction signs a tx, issues the tx and awaits the tx to be minted into a block
func (w *debugWallet) AwaitedSignAndSendTransaction(client ethadapter.EthClient, txData types.TxData) (*types.Transaction, *types.Receipt, error) {
	txData, err := ethadapter.SetTxGasPrice(context.Background(), client, txData, w.Address(), w.GetNonceAndIncrement(), 0, testlog.Logger())
	if err != nil {
		w.SetNonce(w.GetNonce() - 1)
		return nil, nil, err
	}

	signedTx, err := w.SignTransaction(txData)
	if err != nil {
		return nil, nil, err
	}

	err = client.SendTransaction(signedTx)
	if err != nil {
		return nil, nil, err
	}

	var receipt *types.Receipt
	err = retry.Do(func() error {
		receipt, err = client.TransactionReceipt(signedTx.Hash())
		if err != nil {
			return err
		}
		if receipt == nil {
			return fmt.Errorf("no receipt yet")
		}
		return nil
	}, retry.NewTimeoutStrategy(w.receiptTimeout, time.Second))

	return signedTx, receipt, err
}

func (w *debugWallet) debugTransaction(client ethadapter.EthClient, tx *types.Transaction) ([]byte, error) {
	return client.EthClient().CallContract(
		context.Background(),
		ethereum.CallMsg{
			From:       w.Address(),
			To:         tx.To(),
			Gas:        tx.Gas(),
			GasPrice:   tx.GasPrice(),
			GasFeeCap:  tx.GasFeeCap(),
			GasTipCap:  tx.GasTipCap(),
			Value:      tx.Value(),
			Data:       tx.Data(),
			AccessList: tx.AccessList(),
		},
		nil,
	)
}
