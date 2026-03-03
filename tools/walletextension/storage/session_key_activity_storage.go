package storage

import (
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage/database/cosmosdb"
)

// SessionKeyActivityStorage defines persistence for session key activity tracker
type SessionKeyActivityStorage interface {
	// Load retrieves all stored session key activities (used on startup)
	Load() ([]wecommon.SessionKeyActivity, error)
	// Save performs a full replacement of all stored activities
	Save([]wecommon.SessionKeyActivity) error
	// SaveBatch upserts a batch of activities (used for async persistence of evicted entries)
	SaveBatch([]wecommon.SessionKeyActivity) error
	// ListOlderThan returns activities with LastActive before the given cutoff
	ListOlderThan(cutoff time.Time) ([]wecommon.SessionKeyActivity, error)
	// Delete removes a specific activity by address
	Delete(addr gethcommon.Address) error
}

// NewSessionKeyActivityStorage is a factory that returns a concrete storage based on dbType
func NewSessionKeyActivityStorage(dbType, dbConnectionURL string, encryptionKey []byte) (SessionKeyActivityStorage, error) {
	if dbType == "cosmosDB" {
		return cosmosdb.NewSessionKeyActivityStorage(dbConnectionURL, encryptionKey)
	}
	return NewNoOpSessionKeyActivityStorage(), nil
}

// noOpSessionKeyActivityStorage is a no-op implementation used for non-production DBs
type noOpSessionKeyActivityStorage struct{}

func NewNoOpSessionKeyActivityStorage() SessionKeyActivityStorage {
	return &noOpSessionKeyActivityStorage{}
}

func (n *noOpSessionKeyActivityStorage) Load() ([]wecommon.SessionKeyActivity, error) {
	return nil, nil
}

func (n *noOpSessionKeyActivityStorage) Save([]wecommon.SessionKeyActivity) error { return nil }

func (n *noOpSessionKeyActivityStorage) SaveBatch([]wecommon.SessionKeyActivity) error { return nil }

func (n *noOpSessionKeyActivityStorage) ListOlderThan(time.Time) ([]wecommon.SessionKeyActivity, error) {
	return nil, nil
}

func (n *noOpSessionKeyActivityStorage) Delete(gethcommon.Address) error { return nil }
