package storage

import (
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage/database/cosmosdb"
)

// SessionKeyActivityStorage defines persistence for session key activity tracker
type SessionKeyActivityStorage interface {
	Load() ([]wecommon.SessionKeyActivity, error)
	Save([]wecommon.SessionKeyActivity) error
}

// NewSessionKeyActivityStorage is a factory that returns a concrete storage based on dbType
func NewSessionKeyActivityStorage(dbType, dbConnectionURL string) (SessionKeyActivityStorage, error) {
	if dbType == "cosmosDB" {
		return cosmosdb.NewSessionKeyActivityStorage(dbConnectionURL)
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
