package faucet

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/go/wallet"
)

const (
	_timeout    = 60 * time.Second
	NativeToken = "eth"
	// DeprecatedNativeToken is left in temporarily for tooling that is getting native funds using `/ten` URL
	DeprecatedNativeToken = "ten" // todo (@matt) remove this once we have fixed the /ten usages
	WrappedOBX            = "wobx"
	WrappedEth            = "weth"
	WrappedUSDC           = "usdc"
)

type Faucet struct {
	client    *obsclient.AuthObsClient
	fundMutex sync.Mutex
	wallet    wallet.Wallet
	Logger    log.Logger
}

func NewFaucet(rpcURL string, chainID int64, pkString string) (*Faucet, error) {
	logger := log.New()
	w := wallet.NewInMemoryWalletFromConfig(pkString, chainID, logger)
	obsClient, err := obsclient.DialWithAuth(rpcURL, w, logger)
	if err != nil {
		return nil, fmt.Errorf("unable to connect with the node: %w", err)
	}

	return &Faucet{
		client: obsClient,
		wallet: w,
		Logger: logger,
	}, nil
}

func (f *Faucet) Fund(address *common.Address, token string, amount *big.Int) (string, error) {
	var err error
	var signedTx *types.Transaction

	if token == NativeToken || token == DeprecatedNativeToken {
		signedTx, err = f.fundNativeToken(address, amount)
	} else {
		return "", fmt.Errorf("token not fundable atm")
		// todo implement this when contracts are deployable somewhere
	}
	if err != nil {
		return "", err
	}

	// the faucet should be the only user of the faucet pk
	txMarshal, err := json.Marshal(signedTx)
	if err != nil {
		return "", err
	}
	f.Logger.Info(fmt.Sprintf("Funded address: %s - tx: %+v\n", address.Hex(), string(txMarshal)))
	// todo handle tx receipt

	if err := f.validateTx(signedTx); err != nil {
		return "", fmt.Errorf("unable to validate tx %s: %w", signedTx.Hash(), err)
	}

	return signedTx.Hash().Hex(), nil
}

func (f *Faucet) validateTx(tx *types.Transaction) error {
	for now := time.Now(); time.Since(now) < _timeout; time.Sleep(time.Second) {
		receipt, err := f.client.TransactionReceipt(context.Background(), tx.Hash())
		// end eagerly for unexpected errors
		if err != nil && !errors.Is(err, ethereum.NotFound) {
			return fmt.Errorf("could not retrieve transaction receipt in eth_getTransactionReceipt request. Cause: %w", err)
		}

		// try again until timeout
		if receipt == nil {
			continue
		}

		txReceiptBytes, err := receipt.MarshalJSON()
		if err != nil {
			return fmt.Errorf("could not marshal transaction receipt to JSON in eth_getTransactionReceipt request. Cause: %w", err)
		}
		fmt.Println(string(txReceiptBytes))

		if receipt.Status != 1 {
			return fmt.Errorf("tx status is not 0x1")
		}
		return nil
	}
	return fmt.Errorf("unable to fetch tx receipt after %s", _timeout)
}

func (f *Faucet) fundNativeToken(address *common.Address, amount *big.Int) (*types.Transaction, error) {
	// only one funding at the time
	f.fundMutex.Lock()
	defer f.fundMutex.Unlock()

	nonce, err := f.client.NonceAt(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch %s nonce: %w", f.wallet.Address(), err)
	}
	// this isn't great as the tx count might be incremented in between calls
	// but only after removing the pk from other apps can we use a proper counter

	tx := &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: big.NewInt(225),
		To:       address,
		Value:    amount,
	}

	estimatedTx := f.client.EstimateGasAndGasPrice(tx)

	signedTx, err := f.wallet.SignTransaction(estimatedTx)
	if err != nil {
		return nil, err
	}

	if err = f.client.SendTransaction(context.Background(), signedTx); err != nil {
		return signedTx, err
	}

	return signedTx, nil
}

func (f *Faucet) Balance(ctx context.Context) (*big.Int, error) {
	return f.client.BalanceAt(ctx, nil)
}
