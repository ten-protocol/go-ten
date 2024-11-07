package storage

import (
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
	return s.storage.AddUser(userID, privateKey)
}

// DeleteUser deletes a user and invalidates the cache for the userID
func (s *UserStorageWithCache) DeleteUser(userID []byte) error {
	err := s.storage.DeleteUser(userID)
	if err != nil {
		return err
	}
	s.cache.Remove(userID)
	return nil
}

// AddAccount adds an account to a user and invalidates the cache for the userID
func (s *UserStorageWithCache) AddAccount(userID []byte, accountAddress []byte, signature []byte, signatureType viewingkey.SignatureType) error {
	err := s.storage.AddAccount(userID, accountAddress, signature, signatureType)
	if err != nil {
		return err
	}
	s.cache.Remove(userID)
	return nil
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
