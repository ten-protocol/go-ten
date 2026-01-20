package storage

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/tools/walletextension/cache"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

// UserStorageWithCache implements the UserStorage interface with caching
type UserStorageWithCache struct {
	storage UserStorage
	cache   cache.Cache
}

const UserCacheSize = 10_000

// NewUserStorageWithCache creates a new UserStorageWithCache instance
func NewUserStorageWithCache(storage UserStorage, logger log.Logger) (*UserStorageWithCache, error) {
	c, err := cache.NewCache(UserCacheSize, logger)
	if err != nil {
		return nil, err
	}
	return &UserStorageWithCache{
		storage: storage,
		cache:   c,
	}, nil
}

// AddUser adds a new user and invalidates the cache for the userID
func (s *UserStorageWithCache) AddUser(userID []byte, privateKey []byte) error {
	// Invalidate cache before operation to prevent race conditions
	s.cache.Remove(userID)
	return s.storage.AddUser(userID, privateKey)
}

// DeleteUser deletes a user and invalidates the cache for the userID
func (s *UserStorageWithCache) DeleteUser(userID []byte) error {
	// Invalidate cache before operation to prevent race conditions
	s.cache.Remove(userID)
	return s.storage.DeleteUser(userID)
}

func (s *UserStorageWithCache) AddSessionKey(userID []byte, key wecommon.GWSessionKey) error {
	// Invalidate cache before operation to prevent race conditions
	s.cache.Remove(userID)
	return s.storage.AddSessionKey(userID, key)
}

func (s *UserStorageWithCache) RemoveSessionKey(userID []byte, sessionKeyAddr *gethcommon.Address) error {
	// Invalidate cache before operation to prevent race conditions
	s.cache.Remove(userID)
	return s.storage.RemoveSessionKey(userID, sessionKeyAddr)
}

// AddAccount adds an account to a user and invalidates the cache for the userID
// Cache is invalidated BEFORE the operation to prevent race conditions where concurrent
// GetUser calls might read stale cached data during the storage transaction.
// This ensures any concurrent reads will get a cache miss and read fresh data
// from the database after the transaction commits.
func (s *UserStorageWithCache) AddAccount(userID []byte, accountAddress []byte, signature []byte, signatureType viewingkey.SignatureType) error {
	// Invalidate cache BEFORE operation to prevent race conditions
	s.cache.Remove(userID)
	return s.storage.AddAccount(userID, accountAddress, signature, signatureType)
}

// GetUser retrieves a user from the cache or underlying storage
func (s *UserStorageWithCache) GetUser(userID []byte) (*wecommon.GWUser, error) {
	return cache.WithCache(s.cache, &cache.Cfg{Type: cache.LongLiving}, userID, func() (*wecommon.GWUser, error) {
		return s.storage.GetUser(userID)
	})
}

// GetEncryptionKey delegates to the underlying storage
func (s *UserStorageWithCache) GetEncryptionKey() []byte {
	return s.storage.GetEncryptionKey()
}
