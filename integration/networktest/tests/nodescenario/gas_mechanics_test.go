package nodescenario

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/actions"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
)

// TestGasMechanics verifies gas mechanics including L1 publishing costs and base fee calculations
func TestGasMechanics(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)

	// Create test context to store transaction results
	type txResult struct {
		hash    common.Hash
		receipt *types.Receipt
		gasUsed uint64
	}
	var transactions []txResult

	networktest.Run(
		"gas-mechanics",
		t,
		env.LocalDevNetwork(),
		actions.Series(
			// Create test users with different balances
			actions.CreateAndFundTestUsers(2),

			// Test 1: Normal transaction with standard gas limit
			actions.RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
				sender, err := actions.FetchTestUser(ctx, 0)
				if err != nil {
					return ctx, err
				}
				receiver, err := actions.FetchTestUser(ctx, 1)
				if err != nil {
					return ctx, err
				}

				oneEth := big.NewInt(1e14)
				initialBalance, err := sender.NativeBalance(ctx)
				if err != nil {
					return ctx, err
				}

				// Send 1 ETH
				txHash, err := sender.SendFunds(ctx, receiver.Wallet().Address(), oneEth)
				if err != nil {
					return ctx, err
				}

				receipt, err := sender.AwaitReceipt(ctx, txHash)
				if err != nil {
					return ctx, err
				}

				afterBalance, err := sender.NativeBalance(ctx)
				if err != nil {
					return ctx, err
				}
				differenceInBalance := new(big.Int).Sub(initialBalance, afterBalance)
				differenceInBalance = new(big.Int).Sub(differenceInBalance, oneEth)

				gasUsedBalance := new(big.Int).Mul(big.NewInt(int64(receipt.GasUsed)), receipt.EffectiveGasPrice)

				if differenceInBalance.Cmp(gasUsedBalance) != 0 {
					return ctx, fmt.Errorf("balance difference (%v) does not match the gas cost + value transfer (%v)", differenceInBalance, gasUsedBalance)
				}

				// Store transaction result
				transactions = append(transactions, txResult{
					hash:    *txHash,
					receipt: receipt,
					gasUsed: receipt.GasUsed,
					// TODO: Add L1Cost tracking once available
				})
				return ctx, nil
			}),
		),
	)
}
