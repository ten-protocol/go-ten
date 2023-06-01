package useraccountmanager

import (
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/tools/walletextension/accountmanager"
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

func (m *UserAccountManager) AddUserAccountManager(userID string) {
	_, exists := m.userAccountManager[userID]
	if exists {
		m.logger.Warn("UserAccountManager already exists for ", userID)
		return
	}
	newAccountManager := accountmanager.NewAccountManager(m.unauthenticatedClient, m.logger)
	m.userAccountManager[userID] = newAccountManager
}

func (m *UserAccountManager) GetUserAccountManager(userID string) (*accountmanager.AccountManager, error) {
	userAccManager, exists := m.userAccountManager[userID]
	if exists {
		return userAccManager, nil
	}
	return nil, fmt.Errorf("UserAccountManager doesn't exist for user: %s", userID)
}
