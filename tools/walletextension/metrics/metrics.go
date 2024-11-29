package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

type MetricsTracker struct {
	totalUsers         atomic.Uint64
	accountsRegistered atomic.Uint64
	activeUsers        map[string]time.Time
	activeUserLock     sync.RWMutex
}

func NewMetricsTracker() *MetricsTracker {
	mt := &MetricsTracker{
		activeUsers: make(map[string]time.Time),
	}

	// Start cleanup routine for inactive users
	go mt.cleanupInactiveUsers()
	return mt
}

// RecordNewUser increments the total user count
func (mt *MetricsTracker) RecordNewUser() {
	mt.totalUsers.Add(1)
}

// RecordAccountRegistered increments the total number of registered accounts
func (mt *MetricsTracker) RecordAccountRegistered() {
	mt.accountsRegistered.Add(1)
}

// RecordUserActivity updates the last activity timestamp for a user
func (mt *MetricsTracker) RecordUserActivity(anonymousID string) {
	mt.activeUserLock.Lock()
	defer mt.activeUserLock.Unlock()
	mt.activeUsers[anonymousID] = time.Now()
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
	mt.activeUserLock.RLock()
	defer mt.activeUserLock.RUnlock()

	count := 0
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	for _, lastActive := range mt.activeUsers {
		if lastActive.After(thirtyDaysAgo) {
			count++
		}
	}
	return count
}

// cleanupInactiveUsers removes users that haven't been active for more than 30 days
func (mt *MetricsTracker) cleanupInactiveUsers() {
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		mt.activeUserLock.Lock()
		thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

		for userID, lastActive := range mt.activeUsers {
			if lastActive.Before(thirtyDaysAgo) {
				delete(mt.activeUsers, userID)
			}
		}
		mt.activeUserLock.Unlock()
	}
}
