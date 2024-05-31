package helpful

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/actions"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
)

// Smoke tests are useful for checking a network is live or checking basic functionality is not broken

var _transferAmount = big.NewInt(100_000_000)

func TestExecuteNativeFundsTransfer(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"native-funds-smoketest",
		t,
		env.LocalDevNetwork(),
		actions.Series(
			&actions.CreateTestUser{UserID: 0},
			&actions.CreateTestUser{UserID: 1},
			actions.SetContextValue(actions.KeyNumberOfTestUsers, 2),

			&actions.AllocateFaucetFunds{UserID: 0},
			actions.SnapshotUserBalances(actions.SnapAfterAllocation), // record user balances (we have no guarantee on how much the network faucet allocates)

			&actions.SendNativeFunds{FromUser: 0, ToUser: 1, Amount: _transferAmount},

			&actions.VerifyBalanceAfterTest{UserID: 1, ExpectedBalance: _transferAmount},
			&actions.VerifyBalanceDiffAfterTest{UserID: 0, Snapshot: actions.SnapAfterAllocation, ExpectedDiff: big.NewInt(0).Neg(_transferAmount)},
		),
	)
}

// util test that transfers funds from a sepolia account with a known PK to another account
func TestExecuteSepoliaFundsTransfer(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"sepolia-funds-smoketest",
		t,
		env.SepoliaTestnet(),
		actions.RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
			// create wallet for DEPLOYER_PK
			DEPLOYER_PK := ""
			l1Wallet := wallet.NewInMemoryWalletFromConfig(DEPLOYER_PK, _sepoliaChainID, testlog.Logger())
			cli, err := ethadapter.NewEthClientFromURL("https://sepolia.infura.io/v3/187f3057b39849d091727cc3dcfae011", 10*time.Second, gethcommon.HexToAddress("0x0"), testlog.Logger())
			if err != nil {
				panic(err)
			}
			nonce, err := cli.Nonce(l1Wallet.Address())
			if err != nil {
				panic(err)
			}

			l1Wallet.SetNonce(nonce)

			amt := big.NewInt(0).Mul(big.NewInt(0), big.NewInt(int64(params.GWei)))

			destAddr := gethcommon.HexToAddress("0x563EAc5dfDFebA3C53c2160Bf1Bd62941E3D0005")

			gasPrice, err := cli.EthClient().SuggestGasPrice(context.Background())
			if err != nil {
				panic(err)
			}
			// multiply gas price by 2
			gasPrice = big.NewInt(0).Mul(gasPrice, big.NewInt(3))

			// create transaction from l1Wallet to toAddress
			tx := &types.LegacyTx{
				Nonce:    l1Wallet.GetNonce(),
				Value:    amt,
				Gas:      uint64(25_000),
				GasPrice: gasPrice,
				To:       &destAddr,
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
