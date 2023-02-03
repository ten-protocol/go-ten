package helpful

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/obscuronet/go-obscuro/integration/networktest"
	"github.com/obscuronet/go-obscuro/integration/networktest/env"
	"github.com/obscuronet/go-obscuro/integration/networktest/userwallet"
)

// Smoke tests are useful for checking a network is live or checking basic functionality is not broken

func TestCanExecuteNativeFundsTransfer(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(t, env.LocalDevNetwork(), networktest.CreateTest("native-funds-smoketest", nativeFundsSmoketest))
}

func nativeFundsSmoketest(network networktest.NetworkConnector) error {
	ctx := context.Background()
	transferAmt := big.NewInt(10000000000)
	wal1 := userwallet.GenerateRandomWallet(network)
	wal2 := userwallet.GenerateRandomWallet(network)

	err := network.AllocateFaucetFunds(wal1.Address())
	if err != nil {
		return fmt.Errorf("unable to allocate faucet funds - %w", err)
	}

	fmt.Println("Sending funds...")
	err = wal1.SendFunds(ctx, wal2.Address(), transferAmt)
	if err != nil {
		return fmt.Errorf("unable to send funds to address - %w", err)
	}
	fmt.Println("Checking balance...")
	balance, err := wal2.NativeBalance(ctx)
	if err != nil {
		return fmt.Errorf("unable to fetch balance of target acc - %w", err)
	}
	if balance.Cmp(transferAmt) != 0 {
		return fmt.Errorf("expected transfer target's balance to be %d but it was %d", transferAmt, balance)
	}
	return nil
}
