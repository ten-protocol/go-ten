package services

import (
	"sync"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// SessionKeyActivity holds last-activity metadata for a session key
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
	byKey map[gethcommon.Address]SessionKeyActivity
}

func NewSessionKeyActivityTracker() SessionKeyActivityTracker {
	return &sessionKeyActivityTracker{
		byKey: make(map[gethcommon.Address]SessionKeyActivity),
	}
}

func (t *sessionKeyActivityTracker) MarkActive(userID []byte, addr gethcommon.Address) {
	now := time.Now()
	t.mu.Lock()
	existing, ok := t.byKey[addr]
	if ok {
		existing.LastActive = now
		// keep original user association
		t.byKey[addr] = existing
	} else {
		t.byKey[addr] = SessionKeyActivity{Addr: addr, UserID: userID, LastActive: now}
	}
	t.mu.Unlock()
}

func (t *sessionKeyActivityTracker) ListOlderThan(cutoff time.Time) []SessionKeyActivity {
	t.mu.RLock()
	// preallocate with current size upper bound; filter below
	result := make([]SessionKeyActivity, 0, len(t.byKey))
	for _, entry := range t.byKey {
		if entry.LastActive.Before(cutoff) {
			result = append(result, entry)
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
