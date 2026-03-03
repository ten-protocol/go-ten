package services

import (
	"sync"
	"testing"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
)

// mockActivityStorage is a mock implementation of SessionKeyActivityStorage for testing
type mockActivityStorage struct {
	mu     sync.Mutex
	items  map[gethcommon.Address]common.SessionKeyActivity
	saved  []common.SessionKeyActivity
	errors map[string]error
}

func newMockActivityStorage() *mockActivityStorage {
	return &mockActivityStorage{
		items:  make(map[gethcommon.Address]common.SessionKeyActivity),
		saved:  make([]common.SessionKeyActivity, 0),
		errors: make(map[string]error),
	}
}

func (m *mockActivityStorage) Load() ([]common.SessionKeyActivity, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]common.SessionKeyActivity, 0, len(m.items))
	for _, item := range m.items {
		result = append(result, item)
	}
	return result, m.errors["Load"]
}

func (m *mockActivityStorage) Save(items []common.SessionKeyActivity) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items = make(map[gethcommon.Address]common.SessionKeyActivity)
	for _, item := range items {
		m.items[item.Addr] = item
	}
	return m.errors["Save"]
}

func (m *mockActivityStorage) SaveBatch(items []common.SessionKeyActivity) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.saved = append(m.saved, items...)
	for _, item := range items {
		m.items[item.Addr] = item
	}
	return m.errors["SaveBatch"]
}

func (m *mockActivityStorage) ListOlderThan(cutoff time.Time) ([]common.SessionKeyActivity, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]common.SessionKeyActivity, 0)
	for _, item := range m.items {
		if item.LastActive.Before(cutoff) {
			result = append(result, item)
		}
	}
	return result, m.errors["ListOlderThan"]
}

func (m *mockActivityStorage) Delete(addr gethcommon.Address) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.items, addr)
	return m.errors["Delete"]
}

func (m *mockActivityStorage) getSaved() []common.SessionKeyActivity {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]common.SessionKeyActivity{}, m.saved...)
}

func TestSessionKeyActivityTracker_MarkActive_Basic(t *testing.T) {
	tracker := NewSessionKeyActivityTracker(nil)
	defer tracker.Stop()

	addr := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	userID := []byte("test-user-1")

	tracker.MarkActive(userID, addr)

	all := tracker.ListAll()
	require.Len(t, all, 1)
	assert.Equal(t, addr, all[0].Addr)
	assert.Equal(t, userID, all[0].UserID)
}

func TestSessionKeyActivityTracker_MarkActive_UpdateExisting(t *testing.T) {
	tracker := NewSessionKeyActivityTracker(nil)
	defer tracker.Stop()

	addr := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	userID := []byte("test-user-1")

	// First activation
	tracker.MarkActive(userID, addr)
	firstAll := tracker.ListAll()
	firstTime := firstAll[0].LastActive

	// Wait a bit and activate again
	time.Sleep(10 * time.Millisecond)
	tracker.MarkActive(userID, addr)

	all := tracker.ListAll()
	require.Len(t, all, 1)
	assert.True(t, all[0].LastActive.After(firstTime), "LastActive should be updated")
}

func TestSessionKeyActivityTracker_ListOlderThan(t *testing.T) {
	tracker := NewSessionKeyActivityTracker(nil)
	defer tracker.Stop()

	addr1 := gethcommon.HexToAddress("0x1111111111111111111111111111111111111111")
	addr2 := gethcommon.HexToAddress("0x2222222222222222222222222222222222222222")

	// Add first entry
	tracker.MarkActive([]byte("user1"), addr1)

	// Wait and add second entry
	time.Sleep(20 * time.Millisecond)
	cutoff := time.Now()
	time.Sleep(20 * time.Millisecond)

	tracker.MarkActive([]byte("user2"), addr2)

	// Only addr1 should be older than cutoff
	older := tracker.ListOlderThan(cutoff)
	require.Len(t, older, 1)
	assert.Equal(t, addr1, older[0].Addr)
}

func TestSessionKeyActivityTracker_Delete(t *testing.T) {
	tracker := NewSessionKeyActivityTracker(nil)
	defer tracker.Stop()

	addr := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")

	tracker.MarkActive([]byte("user"), addr)
	require.Len(t, tracker.ListAll(), 1)

	deleted := tracker.Delete(addr)
	assert.True(t, deleted)
	assert.Empty(t, tracker.ListAll())

	// Delete non-existent
	deleted = tracker.Delete(addr)
	assert.False(t, deleted)
}

func TestSessionKeyActivityTracker_Load(t *testing.T) {
	tracker := NewSessionKeyActivityTracker(nil)
	defer tracker.Stop()

	items := []common.SessionKeyActivity{
		{
			Addr:       gethcommon.HexToAddress("0x1111111111111111111111111111111111111111"),
			UserID:     []byte("user1"),
			LastActive: time.Now().Add(-2 * time.Hour),
		},
		{
			Addr:       gethcommon.HexToAddress("0x2222222222222222222222222222222222222222"),
			UserID:     []byte("user2"),
			LastActive: time.Now().Add(-1 * time.Hour),
		},
	}

	tracker.Load(items)

	all := tracker.ListAll()
	require.Len(t, all, 2)
}

func TestSessionKeyActivityTracker_LRU_EvictsOldest(t *testing.T) {
	storage := newMockActivityStorage()
	tracker := NewSessionKeyActivityTrackerWithStorage(nil, storage).(*sessionKeyActivityTracker)
	// Override max entries for testing
	tracker.maxEntries = 3

	addr1 := gethcommon.HexToAddress("0x1111111111111111111111111111111111111111")
	addr2 := gethcommon.HexToAddress("0x2222222222222222222222222222222222222222")
	addr3 := gethcommon.HexToAddress("0x3333333333333333333333333333333333333333")
	addr4 := gethcommon.HexToAddress("0x4444444444444444444444444444444444444444")

	// Add 3 entries (at capacity)
	tracker.MarkActive([]byte("user1"), addr1)
	time.Sleep(5 * time.Millisecond)
	tracker.MarkActive([]byte("user2"), addr2)
	time.Sleep(5 * time.Millisecond)
	tracker.MarkActive([]byte("user3"), addr3)

	require.Len(t, tracker.ListAll(), 3)

	// Add 4th entry - should evict addr1 (oldest)
	time.Sleep(5 * time.Millisecond)
	tracker.MarkActive([]byte("user4"), addr4)

	// Should still have 3 entries
	all := tracker.ListAll()
	require.Len(t, all, 3)

	// addr1 should be gone, addr2, addr3, addr4 should remain
	addrs := make(map[gethcommon.Address]bool)
	for _, a := range all {
		addrs[a.Addr] = true
	}

	assert.False(t, addrs[addr1], "addr1 should have been evicted")
	assert.True(t, addrs[addr2], "addr2 should remain")
	assert.True(t, addrs[addr3], "addr3 should remain")
	assert.True(t, addrs[addr4], "addr4 should remain")

	// Stop the tracker to flush pending writes
	tracker.Stop()

	// Check that evicted entry was queued for persistence
	saved := storage.getSaved()
	require.Len(t, saved, 1)
	assert.Equal(t, addr1, saved[0].Addr)
}

func TestSessionKeyActivityTracker_LRU_UpdateMovesToFront(t *testing.T) {
	storage := newMockActivityStorage()
	tracker := NewSessionKeyActivityTrackerWithStorage(nil, storage).(*sessionKeyActivityTracker)
	tracker.maxEntries = 3

	addr1 := gethcommon.HexToAddress("0x1111111111111111111111111111111111111111")
	addr2 := gethcommon.HexToAddress("0x2222222222222222222222222222222222222222")
	addr3 := gethcommon.HexToAddress("0x3333333333333333333333333333333333333333")
	addr4 := gethcommon.HexToAddress("0x4444444444444444444444444444444444444444")

	// Add 3 entries
	tracker.MarkActive([]byte("user1"), addr1) // oldest
	time.Sleep(5 * time.Millisecond)
	tracker.MarkActive([]byte("user2"), addr2)
	time.Sleep(5 * time.Millisecond)
	tracker.MarkActive([]byte("user3"), addr3)

	// Re-activate addr1 - this should move it to front (most recently used)
	time.Sleep(5 * time.Millisecond)
	tracker.MarkActive([]byte("user1"), addr1)

	// Now add addr4 - addr2 should be evicted (it's now the oldest)
	time.Sleep(5 * time.Millisecond)
	tracker.MarkActive([]byte("user4"), addr4)

	all := tracker.ListAll()
	require.Len(t, all, 3)

	addrs := make(map[gethcommon.Address]bool)
	for _, a := range all {
		addrs[a.Addr] = true
	}

	assert.True(t, addrs[addr1], "addr1 should remain (was re-activated)")
	assert.False(t, addrs[addr2], "addr2 should have been evicted (was oldest after addr1 re-activation)")
	assert.True(t, addrs[addr3], "addr3 should remain")
	assert.True(t, addrs[addr4], "addr4 should remain")

	tracker.Stop()
}

func TestSessionKeyActivityTracker_ConcurrentAccess(t *testing.T) {
	tracker := NewSessionKeyActivityTracker(nil)
	defer tracker.Stop()

	var wg sync.WaitGroup
	numGoroutines := 10
	numOpsPerGoroutine := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < numOpsPerGoroutine; j++ {
				addr := gethcommon.HexToAddress("0x" + string(rune('0'+idx)) + "234567890123456789012345678901234567890")
				tracker.MarkActive([]byte("user"), addr)
				tracker.ListAll()
				tracker.ListOlderThan(time.Now())
			}
		}(i)
	}

	wg.Wait()
	// Should not panic or deadlock
}

func TestSessionKeyActivityTracker_LoadWithCapacityLimit(t *testing.T) {
	tracker := NewSessionKeyActivityTracker(nil).(*sessionKeyActivityTracker)
	tracker.maxEntries = 3
	defer tracker.Stop()

	// Create more items than capacity
	items := []common.SessionKeyActivity{
		{
			Addr:       gethcommon.HexToAddress("0x1111111111111111111111111111111111111111"),
			UserID:     []byte("user1"),
			LastActive: time.Now().Add(-4 * time.Hour), // oldest - should be dropped
		},
		{
			Addr:       gethcommon.HexToAddress("0x2222222222222222222222222222222222222222"),
			UserID:     []byte("user2"),
			LastActive: time.Now().Add(-3 * time.Hour), // second oldest - should be dropped
		},
		{
			Addr:       gethcommon.HexToAddress("0x3333333333333333333333333333333333333333"),
			UserID:     []byte("user3"),
			LastActive: time.Now().Add(-2 * time.Hour),
		},
		{
			Addr:       gethcommon.HexToAddress("0x4444444444444444444444444444444444444444"),
			UserID:     []byte("user4"),
			LastActive: time.Now().Add(-1 * time.Hour),
		},
		{
			Addr:       gethcommon.HexToAddress("0x5555555555555555555555555555555555555555"),
			UserID:     []byte("user5"),
			LastActive: time.Now(), // newest
		},
	}

	tracker.Load(items)

	all := tracker.ListAll()
	require.Len(t, all, 3, "should only load up to capacity")

	// Should have the 3 newest entries
	addrs := make(map[gethcommon.Address]bool)
	for _, a := range all {
		addrs[a.Addr] = true
	}

	assert.False(t, addrs[gethcommon.HexToAddress("0x1111111111111111111111111111111111111111")], "oldest should be dropped")
	assert.False(t, addrs[gethcommon.HexToAddress("0x2222222222222222222222222222222222222222")], "second oldest should be dropped")
	assert.True(t, addrs[gethcommon.HexToAddress("0x3333333333333333333333333333333333333333")])
	assert.True(t, addrs[gethcommon.HexToAddress("0x4444444444444444444444444444444444444444")])
	assert.True(t, addrs[gethcommon.HexToAddress("0x5555555555555555555555555555555555555555")])
}
