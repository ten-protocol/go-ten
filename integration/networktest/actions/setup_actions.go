package actions

import (
	"context"
	"fmt"

	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
	"github.com/obscuronet/go-obscuro/integration/networktest"
	"github.com/obscuronet/go-obscuro/integration/networktest/userwallet"
)

type CreateTestUser struct {
	UserID int
}

func (c *CreateTestUser) String() string {
	return fmt.Sprintf("CreateTestUser [id: %d]", c.UserID)
}

func (c *CreateTestUser) Run(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
	logger := testlog.Logger()
	wal := datagenerator.RandomWallet(integration.ObscuroChainID)
	// traffic sim users are round robin-ed onto the validators for now (todo (@matt) - make that overridable)
	user := userwallet.NewUserWallet(wal.PrivateKey(), network.ValidatorRPCAddress(c.UserID%network.NumValidators()), logger)
	return storeTestUser(ctx, c.UserID, user), nil
}

func (c *CreateTestUser) Verify(_ context.Context, _ networktest.NetworkConnector) error {
	return nil
}

type AllocateFaucetFunds struct {
	UserID int
}

func (a *AllocateFaucetFunds) String() string {
	return fmt.Sprintf("AllocateFaucetFunds [user: %d]", a.UserID)
}

func (a *AllocateFaucetFunds) Run(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
	user, err := FetchTestUser(ctx, a.UserID)
	if err != nil {
		return ctx, err
	}
	return ctx, network.AllocateFaucetFunds(ctx, user.Address())
}

func (a *AllocateFaucetFunds) Verify(_ context.Context, _ networktest.NetworkConnector) error {
	return nil
}

// action generators (create series of actions for convenient test setups)

func CreateAndFundTestUsers(numUsers int) *MultiAction {
	var newUserActions []networktest.Action
	for i := 0; i < numUsers; i++ {
		newUserActions = append(newUserActions, &CreateTestUser{UserID: i}, &AllocateFaucetFunds{
			UserID: i,
		})
	}
	// set number of users on the context so downstream know how many test users to access if they want all of them
	newUserActions = append(newUserActions, SetContextValue(KeyNumberOfTestUsers, numUsers))
	newUserActions = append(newUserActions, SnapshotUserBalances(SnapAfterAllocation))
	return Series(newUserActions...)
}

func AuthenticateAllUsers() networktest.Action {
	return RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
		numUsers, err := FetchNumberOfTestUsers(ctx)
		if err != nil {
			return nil, fmt.Errorf("expected number of test users to be set on the context")
		}
		for i := 0; i < numUsers; i++ {
			user, err := FetchTestUser(ctx, i)
			if err != nil {
				return nil, err
			}
			err = user.ResetClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("unable to (re)authenticate client %d - %w", i, err)
			}
		}
		return ctx, nil
	})
}
