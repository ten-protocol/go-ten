package actions

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/userwallet"
)

type SendNativeFunds struct {
	FromUser   int
	ToUser     int
	Amount     *big.Int
	GasLimit   *big.Int
	SkipVerify bool

	user   userwallet.User
	txHash *common.Hash
}

func (s *SendNativeFunds) String() string {
	return fmt.Sprintf("SendNativeFunds [from:%d, to:%d]", s.FromUser, s.ToUser)
}

func (s *SendNativeFunds) Run(ctx context.Context, _ networktest.NetworkConnector) (context.Context, error) {
	user, err := FetchTestUser(ctx, s.FromUser)
	if err != nil {
		return ctx, err
	}
	target, err := FetchTestUser(ctx, s.ToUser)
	if err != nil {
		return ctx, err
	}
	txHash, err := user.SendFunds(ctx, target.Wallet().Address(), s.Amount)
	if err != nil {
		return nil, err
	}
	s.user = user
	s.txHash = txHash
	return ctx, nil
}

func (s *SendNativeFunds) Verify(ctx context.Context, _ networktest.NetworkConnector) error {
	if s.SkipVerify {
		return nil
	}

	receipt, err := s.user.AwaitReceipt(ctx, s.txHash)
	if err != nil {
		return fmt.Errorf("failed to fetch receipt - %w", err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("receipt status not successful, state=%v", receipt.Status)
	}
	return nil
}

type VerifyBalanceAfterTest struct {
	UserID          int
	ExpectedBalance *big.Int
}

func (c *VerifyBalanceAfterTest) String() string {
	return fmt.Sprintf("**Verify Only** VerifyBalanceAfterTest [user:%d]", c.UserID)
}

func (c *VerifyBalanceAfterTest) Run(ctx context.Context, _ networktest.NetworkConnector) (context.Context, error) {
	// this is a verifier action, nothing to do during run
	return ctx, nil
}

func (c *VerifyBalanceAfterTest) Verify(ctx context.Context, _ networktest.NetworkConnector) error {
	user, err := FetchTestUser(ctx, c.UserID)
	if err != nil {
		return err
	}
	bal, err := user.NativeBalance(ctx)
	if err != nil {
		return err
	}
	if bal.Cmp(c.ExpectedBalance) != 0 {
		return fmt.Errorf("unexpected balance, expected=%d, actual=%d", c.ExpectedBalance, bal)
	}
	return nil
}

// VerifyBalanceDiffAfterTest compares the post-test user balance with the balance at the given snapshot
//
// This only checks the balance against an upper bound because of unknown gas spend, it expects `currBal < snapshotBal + expectedDiff`
type VerifyBalanceDiffAfterTest struct {
	UserID       int
	ExpectedDiff *big.Int
	Snapshot     string
}

func (v *VerifyBalanceDiffAfterTest) Run(ctx context.Context, _ networktest.NetworkConnector) (context.Context, error) {
	return ctx, nil
}

func (v *VerifyBalanceDiffAfterTest) Verify(ctx context.Context, _ networktest.NetworkConnector) error {
	snapshotBalance, err := FetchBalanceAtSnapshot(ctx, v.UserID, v.Snapshot)
	if err != nil {
		return fmt.Errorf("failed to fetch balance - %w", err)
	}
	user, err := FetchTestUser(ctx, v.UserID)
	if err != nil {
		return err
	}
	bal, err := user.NativeBalance(ctx)
	if err != nil {
		return err
	}
	expectedMaxBal := big.NewInt(0).Add(snapshotBalance, v.ExpectedDiff)
	if bal.Cmp(expectedMaxBal) >= 0 {
		return fmt.Errorf("expected balance to be lower than %d but was %d", expectedMaxBal, bal)
	}
	return nil
}
