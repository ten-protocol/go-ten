package metrics

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ten-protocol/go-ten/tools/walletextension/storage/database/cosmosdb"
)

const (
	// Persistence intervals (how often metrics are saved to CosmosDB)
	MetricsPersistInterval = 10 * time.Minute

	// Cleanup intervals (how often inactive users are cleaned up)
	InactiveUserCleanupInterval = 1 * time.Hour

	// Activity thresholds
	UserInactivityThreshold = 30 * 24 * time.Hour // 30 days
	MonthlyActiveUserWindow = 30 * 24 * time.Hour // 30 days

	// Batch size for user activity updates
	ActivityBatchSize = 30
)

// Metrics interface defines the metrics operations
type Metrics interface {
	RecordNewUser()
	RecordAccountRegistered()
	RecordUserActivity(anonymousID string)
	GetTotalUsers() uint64
	GetTotalAccountsRegistered() uint64
	GetMonthlyActiveUsers() int
	Stop()
}

type MetricsTracker struct {
	totalUsers         atomic.Uint64
	accountsRegistered atomic.Uint64

	// In-memory cache of recent activities
	activityCache     map[string]time.Time // key is double-hashed userID
	activityCacheLock sync.RWMutex
	activityBatch     map[string]time.Time // Batch for efficient storage updates
	activityBatchLock sync.Mutex

	storage           *cosmosdb.MetricsStorageCosmosDB
	persistTicker     *time.Ticker
	batchUpdateTicker *time.Ticker
}

func NewMetricsTracker(storage *cosmosdb.MetricsStorageCosmosDB) Metrics {
	mt := &MetricsTracker{
		activityCache:     make(map[string]time.Time),
		activityBatch:     make(map[string]time.Time),
		storage:           storage,
		persistTicker:     time.NewTicker(MetricsPersistInterval),
		batchUpdateTicker: time.NewTicker(1 * time.Minute), // Process batches more frequently
	}

	// Load existing metrics
	if metrics, err := storage.LoadMetrics(); err == nil {
		mt.totalUsers.Store(metrics.TotalUsers)
		mt.accountsRegistered.Store(metrics.AccountsRegistered)

		// We no longer load all active users into memory at startup
		// Instead, we'll populate the cache as users become active
	}

	// Start cleanup routine for inactive users
	go mt.cleanupInactiveUsers()
	go mt.persistMetrics()
	go mt.processBatchUpdates()

	return mt
}

// hashUserID creates a double-hashed version of the userID
func (mt *MetricsTracker) hashUserID(userID []byte) string {
	// First hash
	firstHash := sha256.Sum256(userID)
	// Second hash
	secondHash := sha256.Sum256(firstHash[:])
	return hex.EncodeToString(secondHash[:])
}

func (mt *MetricsTracker) RecordNewUser() {
	mt.totalUsers.Add(1)
}

// RecordAccountRegistered increments the total number of registered accounts
func (mt *MetricsTracker) RecordAccountRegistered() {
	mt.accountsRegistered.Add(1)
}

// RecordUserActivity updates the last activity timestamp for a user
func (mt *MetricsTracker) RecordUserActivity(anonymousID string) {
	hashedUserID := mt.hashUserID([]byte(anonymousID))
	now := time.Now()

	// Update in-memory cache
	mt.activityCacheLock.Lock()
	mt.activityCache[hashedUserID] = now
	mt.activityCacheLock.Unlock()

	// Add to batch for efficient storage updates
	mt.activityBatchLock.Lock()
	mt.activityBatch[hashedUserID] = now
	batchSize := len(mt.activityBatch)
	mt.activityBatchLock.Unlock()

	// If batch size exceeds threshold, trigger immediate processing
	if batchSize >= ActivityBatchSize {
		go mt.processBatchUpdates()
	}
}

// GetTotalUsers returns the total number of registered users
func (mt *MetricsTracker) GetTotalUsers() uint64 {
	return mt.totalUsers.Load()
}

// GetTotalAccountsRegistered returns the total number of registered accounts
func (mt *MetricsTracker) GetTotalAccountsRegistered() uint64 {
	return mt.accountsRegistered.Load()
}

// GetMonthlyActiveUsers returns the number of users active in the last 30 days
func (mt *MetricsTracker) GetMonthlyActiveUsers() int {
	// For accurate counts, query the database directly
	activeThreshold := time.Now().Add(-MonthlyActiveUserWindow)
	count, err := mt.storage.CountActiveUsers(activeThreshold)
	if err != nil {
		log.Printf("Failed to count active users: %v", err)

		// Fall back to in-memory cache if database query fails
		mt.activityCacheLock.RLock()
		defer mt.activityCacheLock.RUnlock()

		cacheCount := 0
		for _, lastActive := range mt.activityCache {
			if lastActive.After(activeThreshold) {
				cacheCount++
			}
		}
		return cacheCount
	}

	return count
}

// processBatchUpdates handles batch updates to the database
func (mt *MetricsTracker) processBatchUpdates() {
	// Get current batch and reset
	mt.activityBatchLock.Lock()
	if len(mt.activityBatch) == 0 {
		mt.activityBatchLock.Unlock()
		return
	}

	currentBatch := make(map[string]time.Time)
	for userID, timestamp := range mt.activityBatch {
		currentBatch[userID] = timestamp
	}
	mt.activityBatch = make(map[string]time.Time)
	mt.activityBatchLock.Unlock()

	// Process each user activity update
	for userID, timestamp := range currentBatch {
		if err := mt.storage.UpdateUserActivity(userID, timestamp); err != nil {
			log.Printf("Failed to update user activity: %v", err)

			// Put back in batch on failure
			mt.activityBatchLock.Lock()
			mt.activityBatch[userID] = timestamp
			mt.activityBatchLock.Unlock()
		}
	}
}

// persistMetrics periodically saves global metrics to CosmosDB
func (mt *MetricsTracker) persistMetrics() {
	for range mt.persistTicker.C {
		mt.saveMetrics()
	}
}

func (mt *MetricsTracker) saveMetrics() {
	// Create metrics document with global metrics only
	metrics := &cosmosdb.MetricsDocument{
		ID:                 cosmosdb.METRICS_DOC_ID,
		TotalUsers:         mt.totalUsers.Load(),
		AccountsRegistered: mt.accountsRegistered.Load(),
		ActiveUsers:        make(map[string]string), // Empty, will be populated from shards
		LastUpdated:        time.Now().UTC().Format(time.RFC3339),
	}

	// Process active users from the database
	if currentMetrics, err := mt.storage.LoadMetrics(); err == nil {
		metrics.ActiveUsers = currentMetrics.ActiveUsers
		metrics.ActiveUsersCount = currentMetrics.ActiveUsersCount
	}

	if err := mt.storage.SaveMetrics(metrics); err != nil {
		log.Printf("Failed to persist metrics: %v", err)
	}
}

func (mt *MetricsTracker) cleanupInactiveUsers() {
	ticker := time.NewTicker(InactiveUserCleanupInterval)
	for range ticker.C {
		// Clean up in-memory cache
		mt.activityCacheLock.Lock()
		inactiveThreshold := time.Now().Add(-UserInactivityThreshold)

		for userID, lastActive := range mt.activityCache {
			if lastActive.Before(inactiveThreshold) {
				delete(mt.activityCache, userID)
			}
		}
		mt.activityCacheLock.Unlock()

		// Global metrics will be saved during the regular persist cycle
		// Inactive users in the database will be handled when loading/saving metrics
	}
}

// Stop cleanly stops the metrics tracker
func (mt *MetricsTracker) Stop() {
	mt.persistTicker.Stop()
	mt.batchUpdateTicker.Stop()

	// Process any remaining batch updates
	mt.processBatchUpdates()

	// Final save before stopping
	mt.saveMetrics()
}

// NoOpMetricsTracker implements Metrics interface but does nothing
type NoOpMetricsTracker struct{}

func NewNoOpMetricsTracker() Metrics {
	return &NoOpMetricsTracker{}
}

func (mt *NoOpMetricsTracker) RecordNewUser()                     {}
func (mt *NoOpMetricsTracker) RecordAccountRegistered()           {}
func (mt *NoOpMetricsTracker) RecordUserActivity(string)          {}
func (mt *NoOpMetricsTracker) GetTotalUsers() uint64              { return 0 }
func (mt *NoOpMetricsTracker) GetTotalAccountsRegistered() uint64 { return 0 }
func (mt *NoOpMetricsTracker) GetMonthlyActiveUsers() int         { return 0 }
func (mt *NoOpMetricsTracker) Stop()                              {}
