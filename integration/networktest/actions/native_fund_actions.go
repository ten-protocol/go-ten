package actions

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/integration/networktest"
	"github.com/obscuronet/go-obscuro/integration/networktest/userwallet"
)

type SendNativeFunds struct {
	FromUser int
	ToUser   int
	Amount   *big.Int

	user   *userwallet.UserWallet
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
	txHash, err := user.SendFunds(ctx, target.Address(), s.Amount)
	if err != nil {
		return nil, err
	}
	s.user = user
	s.txHash = txHash
	return ctx, nil
}

func (s *SendNativeFunds) Verify(ctx context.Context, _ networktest.NetworkConnector) error {
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

func (c *VerifyBalanceAfterTest) Verify(ctx context.Context, network networktest.NetworkConnector) error {
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
