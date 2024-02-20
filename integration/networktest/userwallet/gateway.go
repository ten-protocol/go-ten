package userwallet

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/tools/walletextension/lib"
)

type gatewayUser struct {
	privateKey     *ecdsa.PrivateKey
	publicKey      *ecdsa.PublicKey
	accountAddress gethcommon.Address
	chainID        *big.Int

	gwLib  *lib.TGLib // TenGateway utility
	client *ethclient.Client

	// state managed by the wallet
	nonce uint64

	logger gethlog.Logger
}

func NewGatewayUser(pk *ecdsa.PrivateKey, gatewayURL string, logger gethlog.Logger) (*gatewayUser, error) {
	publicKeyECDSA, ok := pk.Public().(*ecdsa.PublicKey)
	if !ok {
		// this shouldn't happen
		logger.Crit("error casting public key to ECDSA")
	}

	gwLib := lib.NewTenGatewayLibrary(gatewayURL, "") // not providing wsURL for now, add if we need it

	err := gwLib.Join()
	if err != nil {
		return nil, fmt.Errorf("failed to join TenGateway: %w", err)
	}
	err = gwLib.RegisterAccount(pk, crypto.PubkeyToAddress(*publicKeyECDSA))
	if err != nil {
		return nil, fmt.Errorf("failed to register account with TenGateway: %w", err)
	}

	client, err := ethclient.Dial(gwLib.HTTP())
	if err != nil {
		return nil, fmt.Errorf("failed to dial TenGateway HTTP: %w", err)
	}

	fmt.Printf("Registered acc with TenGateway: %s (%s)\n", crypto.PubkeyToAddress(*publicKeyECDSA).Hex(), gwLib.HTTP())

	wal := &gatewayUser{
		privateKey:     pk,
		publicKey:      publicKeyECDSA,
		accountAddress: crypto.PubkeyToAddress(*publicKeyECDSA),
		chainID:        big.NewInt(integration.TenChainID),
		gwLib:          gwLib,
		client:         client,
		logger:         logger,
	}

	return wal, nil
}

func (g *gatewayUser) SendFunds(ctx context.Context, addr gethcommon.Address, value *big.Int) (*gethcommon.Hash, error) {
	txData := &types.LegacyTx{
		Nonce: g.nonce,
		Value: value,
		To:    &addr,
	}
	gasPrice, err := g.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to suggest gas price - %w", err)
	}
	txData.GasPrice = gasPrice
	gasLimit, err := g.client.EstimateGas(ctx, ethereum.CallMsg{
		From: g.accountAddress,
		To:   &addr,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to estimate gas - %w", err)
	}
	txData.Gas = gasLimit
	signedTx, err := g.SignTransaction(txData)
	if err != nil {
		return nil, fmt.Errorf("unable to sign transaction - %w", err)
	}
	err = g.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, fmt.Errorf("unable to send transaction - %w", err)
	}
	txHash := signedTx.Hash()
	return &txHash, nil
}

func (g *gatewayUser) AwaitReceipt(ctx context.Context, txHash *gethcommon.Hash) (*types.Receipt, error) {
	var receipt *types.Receipt
	var err error
	err = retry.Do(func() error {
		receipt, err = g.client.TransactionReceipt(ctx, *txHash)
		if !errors.Is(err, rpc.ErrNilResponse) {
			return retry.FailFast(err)
		}
		return err
	}, retry.NewTimeoutStrategy(20*time.Second, 1*time.Second))
	if err != nil {
		return nil, fmt.Errorf("unable to get receipt - %w", err)
	}
	return receipt, nil
}

func (g *gatewayUser) NativeBalance(ctx context.Context) (*big.Int, error) {
	return g.client.BalanceAt(ctx, g.accountAddress, nil)
}

func (g *gatewayUser) Address() gethcommon.Address {
	return g.accountAddress
}

func (g *gatewayUser) SignTransaction(tx types.TxData) (*types.Transaction, error) {
	return g.SignTransactionForChainID(tx, g.chainID)
}

func (g *gatewayUser) SignTransactionForChainID(tx types.TxData, chainID *big.Int) (*types.Transaction, error) {
	return types.SignNewTx(g.privateKey, types.NewLondonSigner(chainID), tx)
}

func (g *gatewayUser) SetNonce(nonce uint64) {
	panic("gatewayUser is designed to manage its own nonce - this method exists to support legacy interface methods")
}

func (g *gatewayUser) GetNonceAndIncrement() uint64 {
	panic("gatewayUser is designed to manage its own nonce - this method exists to support legacy interface methods")
}

func (g *gatewayUser) GetNonce() uint64 {
	return g.nonce
}

func (g *gatewayUser) ChainID() *big.Int {
	return g.chainID
}

func (g *gatewayUser) PrivateKey() *ecdsa.PrivateKey {
	return g.privateKey
}
