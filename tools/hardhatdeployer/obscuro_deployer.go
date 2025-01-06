package contractdeployer

import (
	"context"
	"fmt"
	"math/big"
	"time"

	testcommon "github.com/ten-protocol/go-ten/integration/common"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/go/obsclient/clientutil"
	"github.com/ten-protocol/go-ten/go/wallet"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

func prepareObscuroDeployer(cfg *Config, wal wallet.Wallet, logger gethlog.Logger) (contractDeployerClient, error) {
	client, err := connectClient(getURL(cfg), wal, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to setup obscuro client - %w", err)
	}

	// todo (#1357) - this step doesn't belong in the contract_deployer tool, script should fail for underfunded deployer account
	err = fundDeployerWithFaucet(cfg, client, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to fund deployer acc from faucet - %w", err)
	}

	return &obscuroDeployer{client: client}, nil
}

func fundDeployerWithFaucet(cfg *Config, client *obsclient.AuthObsClient, logger gethlog.Logger) error {
	// Create the L2 faucet wallet and client.
	faucetPrivKey, err := crypto.HexToECDSA(testcommon.TestnetPrefundedPK)
	if err != nil {
		panic("could not initialise L2 faucet private key")
	}
	faucetWallet := wallet.NewInMemoryWalletFromPK(cfg.ChainID, faucetPrivKey, logger)

	faucetClient, err := connectClient(getURL(cfg), faucetWallet, logger)
	if err != nil {
		return err
	}

	balance, err := client.BalanceAt(context.TODO(), nil)
	if err != nil {
		return fmt.Errorf("failed to fetch contract deployer balance. Cause: %w", err)
	}
	// We do not send more funds if the contract deployer already has enough.
	if balance.Cmp(big.NewInt(Prealloc)) == 1 {
		return nil
	}

	nonce, err := faucetClient.NonceAt(context.TODO(), nil)
	if err != nil {
		return fmt.Errorf("failed to fetch faucet nonce. Cause: %w", err)
	}

	destAddr := client.Address()
	txData := &types.LegacyTx{
		Nonce:    nonce,
		Value:    big.NewInt(Prealloc),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		To:       &destAddr,
	}
	tx := faucetClient.EstimateGasAndGasPrice(txData)
	signedTx, err := faucetWallet.SignTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to sign faucet transaction: %w", err)
	}

	err = faucetClient.SendTransaction(context.TODO(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send contract deploy transaction: %w", err)
	}

	rec, err := clientutil.AwaitTransactionReceipt(context.TODO(), faucetClient, signedTx.Hash(), timeoutWait)
	if err != nil {
		return fmt.Errorf("failed to complete faucet transaction. Cause: %w", err)
	}
	if rec.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("faucet transaction status unsuccessful: %v", rec)
	}

	return nil
}

func connectClient(url string, wal wallet.Wallet, logger gethlog.Logger) (*obsclient.AuthObsClient, error) {
	var client *obsclient.AuthObsClient
	var err error

	startConnectingTime := time.Now()
	// since the nodes we are connecting to may have only just started, we retry connection until it is successful
	for client == nil && time.Since(startConnectingTime) < timeoutWait {
		client, err = obsclient.DialWithAuth(url, wal, logger)
		if err == nil {
			break // success
		}
		// if there was an error we'll retry, if we timeout the last seen error will display
		time.Sleep(retryInterval)
	}
	if client == nil {
		return nil, fmt.Errorf("failed to initialise client connection after retrying for %s, %w", timeoutWait, err)
	}
	return client, nil
}

type obscuroDeployer struct {
	client *obsclient.AuthObsClient
}

func (o *obscuroDeployer) Nonce(address gethcommon.Address) (uint64, error) {
	if address != o.client.Address() {
		return 0, fmt.Errorf("nonce requested for a different address to the deploying client address - this shouldn't happen")
	}
	return o.client.NonceAt(context.TODO(), nil)
}

func (o *obscuroDeployer) SendTransaction(tx *types.Transaction) error {
	return o.client.SendTransaction(context.TODO(), tx)
}

func (o *obscuroDeployer) TransactionReceipt(hash gethcommon.Hash) (*types.Receipt, error) {
	return o.client.TransactionReceipt(context.TODO(), hash)
}

func getURL(cfg *Config) string {
	return fmt.Sprintf("ws://%s:%d", cfg.NodeHost, cfg.NodePort)
}
