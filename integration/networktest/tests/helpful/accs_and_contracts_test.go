package helpful

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/actions"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
)

/*
 * This file contains helpful tests for funding accounts, transferring funds, interacting with contracts, etc.
 */

// Run this test to fund an account with native L2 funds
func TestSendFaucetFunds(t *testing.T) {
	// Set the account to fund here
	accountToFund := common.HexToAddress("<account address to fund>")

	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"send-faucet-funds",
		t,
		env.LongRunningLocalNetwork(""),
		&actions.AllocateFaucetFunds{Account: &accountToFund},
	)
}

// Run this test to send native L1 ETH from one account to another
func TestTransferL1Funds(t *testing.T) {
	// Set the accounts addresses and amount to send here
	fromPK := "<pk string with no 0x prefix>"
	to := common.HexToAddress("<account address to send to>")
	// amount to send in wei
	amt := big.NewInt(0).Mul(big.NewInt(1), big.NewInt(int64(params.Ether)))

	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"send-native-funds",
		t,
		env.SepoliaTestnet(),
		actions.RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
			l1Wallet := wallet.NewInMemoryWalletFromConfig(fromPK, _sepoliaChainID, testlog.Logger())
			cli, err := network.GetL1Client()
			if err != nil {
				panic(err)
			}
			// in sepolia if you have issues, you may need a more reliable RPC endpoint, e.g. infura with an api-key:
			// cli, err := ethadapter.NewEthClientFromURL("https://sepolia.infura.io/v3/<api-key>", 10*time.Second, common.HexToAddress("0x0"), testlog.Logger())
			nonce, err := cli.Nonce(l1Wallet.Address())
			if err != nil {
				panic(err)
			}
			l1Wallet.SetNonce(nonce)

			gasPrice, err := cli.EthClient().SuggestGasPrice(context.Background())
			if err != nil {
				panic(err)
			}
			// apply multiplier to the gas price here if you want to guarantee it goes through quickly
			// gasPrice = big.NewInt(0).Mul(gasPrice, big.NewInt(2))

			// create transaction from l1Wallet to toAddress
			tx := &types.LegacyTx{
				Nonce:    l1Wallet.GetNonce(),
				Value:    amt,
				Gas:      uint64(25_000),
				GasPrice: gasPrice,
				To:       &to,
			}
			signedTx, err := l1Wallet.SignTransaction(tx)
			if err != nil {
				panic(err)
			}

			err = cli.SendTransaction(signedTx)
			if err != nil {
				panic(err)
			}

			// await receipt
			err = retry.Do(func() error {
				receipt, err := cli.TransactionReceipt(signedTx.Hash())
				if err != nil {
					return err
				}
				if receipt == nil {
					return fmt.Errorf("no receipt yet")
				}
				if receipt.Status != types.ReceiptStatusSuccessful {
					return retry.FailFast(fmt.Errorf("receipt had status failed for transaction %s", signedTx.Hash().Hex()))
				}
				return nil
			}, retry.NewTimeoutStrategy(70*time.Second, 20*time.Second))

			if err != nil {
				panic(err)
			}

			return ctx, nil
		}),
	)
}
