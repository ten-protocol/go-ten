package metrics

import (
	"crypto/sha256"
	"sync"
	"sync/atomic"
	"time"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage/database/cosmosdb"
)

const (
	// Persistence intervals (how often metrics are saved to CosmosDB)
	MetricsPersistInterval = 10 * time.Minute

	// Cleanup intervals (how often inactive users are cleaned up)
	InactiveUserCleanupInterval = 1 * time.Hour

	// Update intervals for daily stats
	DailyStatsUpdateInterval = 1 * time.Hour

	// Activity thresholds
	UserInactivityThreshold = 30 * 24 * time.Hour // 30 days
	MonthlyActiveUserWindow = 30 * 24 * time.Hour // 30 days

	// Batch size for user activity updates
	ActivityBatchSize = 100
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
	dailyStatsTicker  *time.Ticker
	logger            gethlog.Logger
}

func NewMetricsTracker(storage *cosmosdb.MetricsStorageCosmosDB, logger gethlog.Logger) Metrics {
	mt := &MetricsTracker{
		activityCache:     make(map[string]time.Time),
		activityBatch:     make(map[string]time.Time),
		storage:           storage,
		persistTicker:     time.NewTicker(MetricsPersistInterval),
		batchUpdateTicker: time.NewTicker(1 * time.Minute), // Process batches more frequently
		dailyStatsTicker:  time.NewTicker(DailyStatsUpdateInterval),
		logger:            logger,
	}

	// Load existing metrics
	if global, err := storage.LoadGlobalMetrics(); err == nil {
		mt.totalUsers.Store(global.TotalUsers)
		mt.accountsRegistered.Store(global.AccountsRegistered)
	}

	// Start background routines
	go mt.cleanupInactiveUsers()
	go mt.persistMetrics()
	go mt.processBatchUpdates()
	go mt.updateDailyStats()

	// Update stats immediately on startup
	go mt.updateDailyStats()

	return mt
}

// hashUserID creates a double-hashed version of the userID, using only the first 8 bytes of the hash
func (mt *MetricsTracker) hashUserID(userID []byte) []byte {
	// First hash
	firstHash := sha256.Sum256(userID)
	// Second hash
	secondHash := sha256.Sum256(firstHash[:])
	// Return only first 8 bytes of the hash (sufficient for activity tracking)
	return secondHash[:8]
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
	mt.activityCache[string(hashedUserID)] = now
	mt.activityCacheLock.Unlock()

	// Add to batch for efficient storage updates
	mt.activityBatchLock.Lock()
	mt.activityBatch[string(hashedUserID)] = now
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
		mt.logger.Error("Failed to count active users", "error", err)

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
			mt.logger.Error("Failed to update user activity", "error", err)

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
	// Load the current global metrics
	global, err := mt.storage.LoadGlobalMetrics()
	if err != nil {
		mt.logger.Error("Failed to load global metrics", "error", err)
		return
	}

	// Update with latest values
	global.TotalUsers = mt.totalUsers.Load()
	global.AccountsRegistered = mt.accountsRegistered.Load()

	// Update active user count
	activeThreshold := time.Now().Add(-MonthlyActiveUserWindow)
	count, err := mt.storage.CountActiveUsers(activeThreshold)
	if err == nil {
		global.ActiveUsersCount = count
	}

	// Save the updated global metrics
	if err := mt.storage.SaveGlobalMetrics(global); err != nil {
		mt.logger.Error("Failed to persist global metrics", "error", err)
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
		// Inactive users in the database will be cleaned up when each shard is loaded
		// and written back during regular operations
	}
}

// updateDailyStats periodically updates the daily activity statistics
func (mt *MetricsTracker) updateDailyStats() {
	for range mt.dailyStatsTicker.C {
		if err := mt.storage.UpdateDailyStats(); err != nil {
			mt.logger.Error("Failed to update daily stats", "error", err)
		}
	}
}

// Stop cleanly stops the metrics tracker
func (mt *MetricsTracker) Stop() {
	mt.persistTicker.Stop()
	mt.batchUpdateTicker.Stop()
	mt.dailyStatsTicker.Stop()

	// Process any remaining batch updates
	mt.processBatchUpdates()

	// Final save before stopping
	mt.saveMetrics()

	// Final update of daily stats
	if err := mt.storage.UpdateDailyStats(); err != nil {
		mt.logger.Error("Failed to update daily stats during shutdown", "error", err)
	}
}

// NoOpMetricsTracker implements Metrics interface but does nothing
type NoOpMetricsTracker struct {
	logger gethlog.Logger
}

func NewNoOpMetricsTracker(logger gethlog.Logger) Metrics {
	return &NoOpMetricsTracker{
		logger: logger,
	}
}

func (mt *NoOpMetricsTracker) RecordNewUser()                     {}
func (mt *NoOpMetricsTracker) RecordAccountRegistered()           {}
func (mt *NoOpMetricsTracker) RecordUserActivity(string)          {}
func (mt *NoOpMetricsTracker) GetTotalUsers() uint64              { return 0 }
func (mt *NoOpMetricsTracker) GetTotalAccountsRegistered() uint64 { return 0 }
func (mt *NoOpMetricsTracker) GetMonthlyActiveUsers() int         { return 0 }
func (mt *NoOpMetricsTracker) Stop()                              {}
