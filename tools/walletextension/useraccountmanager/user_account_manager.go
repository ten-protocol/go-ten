package useraccountmanager

import (
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/tools/walletextension/accountmanager"
)

type UserAccountManager struct {
	userAccountManager    map[string]*accountmanager.AccountManager
	unauthenticatedClient rpc.Client
	logger                gethlog.Logger
}

func NewUserAccountManager(unauthenticatedClient rpc.Client, logger gethlog.Logger) UserAccountManager {
	return UserAccountManager{
		userAccountManager:    make(map[string]*accountmanager.AccountManager),
		unauthenticatedClient: unauthenticatedClient,
		logger:                logger,
	}
}

// AddAndReturnAccountManager adds new UserAccountManager if it doesn't exist and returns it, if UserAccountManager already exists for that user just return it
func (m *UserAccountManager) AddAndReturnAccountManager(userID string) *accountmanager.AccountManager {
	existingUserAccountManager, exists := m.userAccountManager[userID]
	if exists {
		return existingUserAccountManager
	}
	newAccountManager := accountmanager.NewAccountManager(m.unauthenticatedClient, m.logger)
	m.userAccountManager[userID] = newAccountManager
	return newAccountManager
}

// GetUserAccountManager retrieves the UserAccountManager associated with the given userID.
// It returns the UserAccountManager and nil error if one exists.
// If a UserAccountManager does not exist for the userID, it returns nil and an error.
func (m *UserAccountManager) GetUserAccountManager(userID string) (*accountmanager.AccountManager, error) {
	userAccManager, exists := m.userAccountManager[userID]
	if exists {
		return userAccManager, nil
	}
	return nil, fmt.Errorf("UserAccountManager doesn't exist for user: %s", userID)
}

// DeleteUserAccountManager removes the UserAccountManager associated with the given userID.
// It returns an error if no UserAccountManager exists for that userID.
func (m *UserAccountManager) DeleteUserAccountManager(userID string) error {
	_, exists := m.userAccountManager[userID]
	if !exists {
		return fmt.Errorf("no UserAccountManager exists for userID %s", userID)
	}
	delete(m.userAccountManager, userID)
	return nil
}
