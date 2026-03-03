package services

import (
	"container/list"
	"sync"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage"
)

// SessionKeyActivityTracker exposes a minimal API for tracking activity
type SessionKeyActivityTracker interface {
	MarkActive(userID []byte, addr gethcommon.Address)
	ListOlderThan(cutoff time.Time) []common.SessionKeyActivity
	ListAll() []common.SessionKeyActivity
	Load(items []common.SessionKeyActivity)
	Delete(addr gethcommon.Address) bool
	// Stop gracefully shuts down the tracker, flushing pending writes
	Stop()
}

// lruEntry represents an entry in the LRU cache
type lruEntry struct {
	addr       gethcommon.Address
	userID     []byte
	lastActive time.Time
}

type sessionKeyActivityTracker struct {
	mu sync.RWMutex

	// LRU cache: doubly-linked list for O(1) eviction of oldest entry
	// Front = most recently used, Back = least recently used (oldest)
	lruList *list.List
	// Map for O(1) lookup by address
	byKey map[gethcommon.Address]*list.Element

	// maxEntries bounds memory usage; when full, oldest entry is evicted
	maxEntries int
	logger     gethlog.Logger

	// Async write queue for persisting evicted entries to DB
	persistQueue   chan common.SessionKeyActivity
	persistStorage storage.SessionKeyActivityStorage
	stopChan       chan struct{}
	stopOnce       sync.Once
	wg             sync.WaitGroup
}

// Configuration constants
const (
	defaultMaxActivityEntries = 100000
	persistQueueSize          = 10000
	persistBatchSize          = 100
	persistFlushInterval      = 5 * time.Second
)

func NewSessionKeyActivityTracker(logger gethlog.Logger) SessionKeyActivityTracker {
	return NewSessionKeyActivityTrackerWithStorage(logger, nil)
}

// NewSessionKeyActivityTrackerWithStorage creates a tracker with async DB persistence
func NewSessionKeyActivityTrackerWithStorage(logger gethlog.Logger, persistStorage storage.SessionKeyActivityStorage) SessionKeyActivityTracker {
	t := &sessionKeyActivityTracker{
		lruList:        list.New(),
		byKey:          make(map[gethcommon.Address]*list.Element),
		maxEntries:     defaultMaxActivityEntries,
		logger:         logger,
		persistStorage: persistStorage,
		stopChan:       make(chan struct{}),
	}

	// Start async persistence worker if storage is provided
	if persistStorage != nil {
		t.persistQueue = make(chan common.SessionKeyActivity, persistQueueSize)
		t.wg.Add(1)
		go t.persistWorker()
	}

	return t
}

// persistWorker runs in the background and batches writes to CosmosDB
func (t *sessionKeyActivityTracker) persistWorker() {
	defer t.wg.Done()

	batch := make([]common.SessionKeyActivity, 0, persistBatchSize)
	ticker := time.NewTicker(persistFlushInterval)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		if err := t.persistStorage.SaveBatch(batch); err != nil {
			if t.logger != nil {
				t.logger.Warn("Failed to persist evicted session key activities", "count", len(batch), "error", err)
			}
		} else {
			if t.logger != nil {
				t.logger.Debug("Persisted evicted session key activities", "count", len(batch))
			}
		}
		batch = batch[:0]
	}

	for {
		select {
		case item, ok := <-t.persistQueue:
			if !ok {
				// Channel closed, flush remaining and exit
				flush()
				return
			}
			batch = append(batch, item)
			if len(batch) >= persistBatchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		case <-t.stopChan:
			// Drain remaining items from queue, checking for closed channel
			for {
				select {
				case item, ok := <-t.persistQueue:
					if !ok {
						// Channel closed, flush and exit
						flush()
						return
					}
					batch = append(batch, item)
				default:
					flush()
					return
				}
			}
		}
	}
}

// Stop gracefully shuts down the tracker, flushing pending writes
func (t *sessionKeyActivityTracker) Stop() {
	t.stopOnce.Do(func() {
		close(t.stopChan)
		if t.persistQueue != nil {
			close(t.persistQueue)
		}
	})
	t.wg.Wait()
}

func (t *sessionKeyActivityTracker) MarkActive(userID []byte, addr gethcommon.Address) {
	now := time.Now()
	t.mu.Lock()
	defer t.mu.Unlock()

	// If the address already exists, update and move to front (most recently used)
	if elem, ok := t.byKey[addr]; ok {
		entry := elem.Value.(*lruEntry)
		entry.lastActive = now
		t.lruList.MoveToFront(elem)
		return
	}

	// New entry: check capacity
	if len(t.byKey) >= t.maxEntries {
		// Evict the oldest entry (back of the list)
		t.evictOldest()
	}

	// Add new entry at front (most recently used)
	entry := &lruEntry{
		addr:       addr,
		userID:     userID,
		lastActive: now,
	}
	elem := t.lruList.PushFront(entry)
	t.byKey[addr] = elem
}

// evictOldest removes the least recently used entry and queues it for DB persistence
// Must be called with lock held
func (t *sessionKeyActivityTracker) evictOldest() {
	back := t.lruList.Back()
	if back == nil {
		return
	}

	entry := back.Value.(*lruEntry)

	// Queue for async DB persistence before removing from memory
	if t.persistQueue != nil {
		activity := common.SessionKeyActivity{
			Addr:       entry.addr,
			UserID:     entry.userID,
			LastActive: entry.lastActive,
		}
		select {
		case t.persistQueue <- activity:
			// Successfully queued
		default:
			// Queue full, log warning but continue with eviction
			if t.logger != nil {
				t.logger.Warn("Persist queue full, evicted entry may be lost", "addr", entry.addr.Hex())
			}
		}
	}

	// Remove from cache
	t.lruList.Remove(back)
	delete(t.byKey, entry.addr)

	if t.logger != nil {
		t.logger.Debug("Evicted oldest session key activity", "addr", entry.addr.Hex(), "lastActive", entry.lastActive)
	}
}

func (t *sessionKeyActivityTracker) ListOlderThan(cutoff time.Time) []common.SessionKeyActivity {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make([]common.SessionKeyActivity, 0)
	for elem := t.lruList.Back(); elem != nil; elem = elem.Prev() {
		entry := elem.Value.(*lruEntry)
		if entry.lastActive.Before(cutoff) {
			result = append(result, common.SessionKeyActivity{
				Addr:       entry.addr,
				UserID:     entry.userID,
				LastActive: entry.lastActive,
			})
		} else {
			// Since list is ordered by last access time (oldest at back),
			// once we hit an entry newer than cutoff, all remaining entries
			// will also be newer, so we can stop early
			break
		}
	}
	return result
}

func (t *sessionKeyActivityTracker) ListAll() []common.SessionKeyActivity {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make([]common.SessionKeyActivity, 0, len(t.byKey))
	for elem := t.lruList.Front(); elem != nil; elem = elem.Next() {
		entry := elem.Value.(*lruEntry)
		result = append(result, common.SessionKeyActivity{
			Addr:       entry.addr,
			UserID:     entry.userID,
			LastActive: entry.lastActive,
		})
	}
	return result
}

func (t *sessionKeyActivityTracker) Load(items []common.SessionKeyActivity) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Clear existing data
	t.lruList = list.New()
	t.byKey = make(map[gethcommon.Address]*list.Element)

	// Sort items by LastActive (oldest first) so we can build the LRU list correctly
	// Items at the front will be most recent, items at back will be oldest
	// We'll add them in reverse order (oldest first) so oldest ends up at back
	sorted := make([]common.SessionKeyActivity, len(items))
	copy(sorted, items)

	// Simple insertion sort by LastActive (ascending = oldest first)
	for i := 1; i < len(sorted); i++ {
		j := i
		for j > 0 && sorted[j].LastActive.Before(sorted[j-1].LastActive) {
			sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
			j--
		}
	}

	// Enforce capacity limit
	startIdx := 0
	if len(sorted) > t.maxEntries {
		startIdx = len(sorted) - t.maxEntries
		if t.logger != nil {
			t.logger.Warn("Load truncated due to capacity, oldest entries dropped",
				"total", len(sorted), "dropped", startIdx, "loaded", t.maxEntries)
		}
	}

	// Add entries: oldest first (will be at back of list), newest last (will be at front)
	for i := startIdx; i < len(sorted); i++ {
		item := sorted[i]
		entry := &lruEntry{
			addr:       item.Addr,
			userID:     item.UserID,
			lastActive: item.LastActive,
		}
		elem := t.lruList.PushFront(entry)
		t.byKey[item.Addr] = elem
	}
}

func (t *sessionKeyActivityTracker) Delete(addr gethcommon.Address) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if elem, ok := t.byKey[addr]; ok {
		t.lruList.Remove(elem)
		delete(t.byKey, addr)
		return true
	}
	return false
}
