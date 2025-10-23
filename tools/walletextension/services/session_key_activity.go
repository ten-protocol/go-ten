package services

import (
	"sync"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

// SessionKeyActivity holds last-activity metadata for a session key (public return type)
type SessionKeyActivity struct {
	Addr       gethcommon.Address
	UserID     []byte
	LastActive time.Time
}

// SessionKeyActivityTracker exposes a minimal API for tracking activity
type SessionKeyActivityTracker interface {
	MarkActive(userID []byte, addr gethcommon.Address)
	ListOlderThan(cutoff time.Time) []SessionKeyActivity
	Delete(addr gethcommon.Address) bool
}

type sessionKeyActivityTracker struct {
	mu    sync.RWMutex
	byKey map[gethcommon.Address]sessionKeyActivityState
	// maxEntries bounds memory usage; when full, oldest entry is evicted upon new insert
	maxEntries int
	logger     gethlog.Logger
}

// sessionKeyActivityState is the internal storage value; address is the map key
type sessionKeyActivityState struct {
	UserID     []byte
	LastActive time.Time
}

// defaultMaxActivityEntries defines an upper bound to avoid unbounded memory growth
const defaultMaxActivityEntries = 100000

func NewSessionKeyActivityTracker(logger gethlog.Logger) SessionKeyActivityTracker {
	return &sessionKeyActivityTracker{
		byKey:      make(map[gethcommon.Address]sessionKeyActivityState),
		maxEntries: defaultMaxActivityEntries,
		logger:     logger,
	}
}

func (t *sessionKeyActivityTracker) MarkActive(userID []byte, addr gethcommon.Address) {
	now := time.Now()
	t.mu.Lock()
	// if the address is already in the map, update the last active time
	if state, ok := t.byKey[addr]; ok {
		state.LastActive = now
		t.byKey[addr] = state
	} else {
		// check if the map is at capacity
		if len(t.byKey) >= t.maxEntries {
			if t.logger != nil {
				t.logger.Warn("SessionKeyActivityTracker capacity reached; dropping new activity", "capacity", t.maxEntries, "addr", addr.Hex())
			}
		} else {
			// if the map is not at capacity, add the address to the map
			t.byKey[addr] = sessionKeyActivityState{UserID: userID, LastActive: now}
		}
	}
	t.mu.Unlock()
}

func (t *sessionKeyActivityTracker) ListOlderThan(cutoff time.Time) []SessionKeyActivity {
	t.mu.RLock()
	// preallocate with current size upper bound; filter below
	result := make([]SessionKeyActivity, 0, len(t.byKey))
	for addr, state := range t.byKey {
		if state.LastActive.Before(cutoff) {
			result = append(result, SessionKeyActivity{Addr: addr, UserID: state.UserID, LastActive: state.LastActive})
		}
	}
	t.mu.RUnlock()
	return result
}

func (t *sessionKeyActivityTracker) Delete(addr gethcommon.Address) bool {
	t.mu.Lock()
	_, existed := t.byKey[addr]
	if existed {
		delete(t.byKey, addr)
	}
	t.mu.Unlock()
	return existed
}
