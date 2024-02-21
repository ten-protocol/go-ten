package actions

import (
	"context"
	"fmt"

	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/datagenerator"
	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/userwallet"
)

type CreateTestUser struct {
	UserID     int
	UseGateway bool
}

func (c *CreateTestUser) String() string {
	return fmt.Sprintf("CreateTestUser [id: %d]", c.UserID)
}

func (c *CreateTestUser) Run(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
	logger := testlog.Logger()

	wal := datagenerator.RandomWallet(integration.TenChainID)
	var user userwallet.User
	if c.UseGateway {
		gwURL, err := network.GetGatewayURL()
		if err != nil {
			return ctx, fmt.Errorf("failed to get required gateway URL: %w", err)
		}
		user, err = userwallet.NewGatewayUser(wal, gwURL, logger)
		if err != nil {
			return ctx, fmt.Errorf("failed to create gateway user: %w", err)
		}
	} else {
		// traffic sim users are round robin-ed onto the validators for now (todo (@matt) - make that overridable)
		user = userwallet.NewUserWallet(wal, network.ValidatorRPCAddress(c.UserID%network.NumValidators()), logger)
	}
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
	return ctx, network.AllocateFaucetFunds(ctx, user.Wallet().Address())
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
