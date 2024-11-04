package storage

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/tools/walletextension/cache"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

// StorageWithCache implements the Storage interface with caching
type StorageWithCache struct {
	storage Storage
	cache   cache.Cache
	mu      sync.RWMutex
}

// NewStorageWithCache creates a new StorageWithCache instance
func NewStorageWithCache(storage Storage, logger log.Logger) (*StorageWithCache, error) {
	c, err := cache.NewCache(logger)
	if err != nil {
		return nil, err
	}
	return &StorageWithCache{
		storage: storage,
		cache:   c,
	}, nil
}

// AddUser adds a new user and invalidates the cache for the userID
func (s *StorageWithCache) AddUser(userID []byte, privateKey []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.storage.AddUser(userID, privateKey)
	if err != nil {
		return err
	}
	s.cache.Remove(userID)
	return nil
}

// DeleteUser deletes a user and invalidates the cache for the userID
func (s *StorageWithCache) DeleteUser(userID []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.storage.DeleteUser(userID)
	if err != nil {
		return err
	}
	s.cache.Remove(userID)
	return nil
}

// AddAccount adds an account to a user and invalidates the cache for the userID
func (s *StorageWithCache) AddAccount(userID []byte, accountAddress []byte, signature []byte, signatureType viewingkey.SignatureType) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.storage.AddAccount(userID, accountAddress, signature, signatureType)
	if err != nil {
		return err
	}
	s.cache.Remove(userID)
	return nil
}

// GetUser retrieves a user from the cache or underlying storage
func (s *StorageWithCache) GetUser(userID []byte) (wecommon.GWUserDB, error) {
	s.mu.RLock()
	// Check if the user is in the cache
	if cachedUser, found := s.cache.Get(userID); found {
		s.mu.RUnlock()
		return cachedUser.(wecommon.GWUserDB), nil
	}
	s.mu.RUnlock()

	// If not in cache, retrieve from storage
	user, err := s.storage.GetUser(userID)
	if err != nil {
		return wecommon.GWUserDB{}, err
	}

	// Store the retrieved user in the cache
	s.mu.Lock()
	s.cache.Set(userID, user, 5*time.Minute)
	s.mu.Unlock()

	return user, nil
}
