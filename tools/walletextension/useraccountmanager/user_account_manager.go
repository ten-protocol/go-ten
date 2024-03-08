package useraccountmanager

import (
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/tools/walletextension/accountmanager"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage"
)

// UserAccountManager is a struct that stores one account manager per user and other required data
type UserAccountManager struct {
	userAccountManager    map[string]*accountmanager.AccountManager
	unauthenticatedClient rpc.Client
	storage               storage.Storage
	hostRPCBinAddrHTTP    string
	hostRPCBinAddrWS      string
	logger                gethlog.Logger
	mu                    sync.Mutex
}

func NewUserAccountManager(unauthenticatedClient rpc.Client, logger gethlog.Logger, storage storage.Storage, hostRPCBindAddrHTTP string, hostRPCBindAddrWS string) UserAccountManager {
	return UserAccountManager{
		userAccountManager:    make(map[string]*accountmanager.AccountManager),
		unauthenticatedClient: unauthenticatedClient,
		storage:               storage,
		hostRPCBinAddrHTTP:    hostRPCBindAddrHTTP,
		hostRPCBinAddrWS:      hostRPCBindAddrWS,
		logger:                logger,
	}
}

// AddAndReturnAccountManager adds new UserAccountManager if it doesn't exist and returns it, if UserAccountManager already exists for that user just return it
func (m *UserAccountManager) AddAndReturnAccountManager(userID string) *accountmanager.AccountManager {
	m.mu.Lock()
	defer m.mu.Unlock()
	existingUserAccountManager, exists := m.userAccountManager[userID]
	if exists {
		return existingUserAccountManager
	}
	newAccountManager := accountmanager.NewAccountManager(userID, m.unauthenticatedClient, m.hostRPCBinAddrWS, m.storage, m.logger)
	m.userAccountManager[userID] = newAccountManager
	return newAccountManager
}

// GetUserAccountManager retrieves the UserAccountManager associated with the given userID.
// it returns the UserAccountManager and nil error if one exists.
// before returning it checks the database and creates all missing clients for that userID
// (we are not loading all of them at startup to limit the number of established connections)
// If a UserAccountManager does not exist for the userID, it returns nil and an error.
func (m *UserAccountManager) GetUserAccountManager(userID string) (*accountmanager.AccountManager, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	userAccManager, exists := m.userAccountManager[userID]
	if !exists {
		return nil, fmt.Errorf("UserAccountManager doesn't exist for user: %s", userID)
	}

	// we have userAccountManager as expected.
	// now we need to create all clients that don't exist there yet
	addressesWithClients := userAccManager.GetAllAddressesWithClients()

	// get all addresses for current userID
	userIDbytes, err := hex.DecodeString(userID)
	if err != nil {
		return nil, err
	}

	// log that we don't have a storage, but still return existing userAccountManager
	// this should never happen, but is useful for tests
	if m.storage == nil {
		m.logger.Error("storage is nil in UserAccountManager")
		return userAccManager, nil
	}

	databaseAccounts, err := m.storage.GetAccounts(userIDbytes)
	if err != nil {
		return nil, err
	}

	userPrivateKey, err := m.storage.GetUserPrivateKey(userIDbytes)
	if err != nil {
		return nil, err
	}

	for _, account := range databaseAccounts {
		addressHexString := common.BytesToAddress(account.AccountAddress).Hex()
		// check if a client for the current address already exists (and skip it if it does)
		if addressAlreadyExists(addressHexString, addressesWithClients) {
			continue
		}

		// create a new client
		encClient, err := wecommon.CreateEncClient(m.hostRPCBinAddrWS, account.AccountAddress, userPrivateKey, account.Signature, account.SignatureType, m.logger)
		if err != nil {
			m.logger.Error(fmt.Errorf("error creating new client, %w", err).Error())
		}

		// add a client to requested userAccountManager
		userAccManager.AddClient(common.BytesToAddress(account.AccountAddress), encClient)
		addressesWithClients = append(addressesWithClients, addressHexString)
	}

	return userAccManager, nil
}

// DeleteUserAccountManager removes the UserAccountManager associated with the given userID.
// It returns an error if no UserAccountManager exists for that userID.
func (m *UserAccountManager) DeleteUserAccountManager(userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, exists := m.userAccountManager[userID]
	if !exists {
		return fmt.Errorf("no UserAccountManager exists for userID %s", userID)
	}
	delete(m.userAccountManager, userID)
	return nil
}

// addressAlreadyExists is a helper function to check if an address is already present in a list of existing addresses
func addressAlreadyExists(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
