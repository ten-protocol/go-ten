package p2p

import (
	"sync"
	"time"
)

// peerTracker tracks the last message received from different peers
type peerTracker struct {
	lock                      sync.RWMutex
	lastReceivedMessageByPeer map[string]time.Time
}

func newPeerTracker() *peerTracker {
	return &peerTracker{
		lock:                      sync.RWMutex{},
		lastReceivedMessageByPeer: map[string]time.Time{},
	}
}

func (s *peerTracker) receivedPeerMsg(peer string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.lastReceivedMessageByPeer[peer] = time.Now()
}

func (s *peerTracker) receivedMessagesByPeer() map[string]time.Time {
	s.lock.RLock()
	defer s.lock.RUnlock()
	newMap := map[string]time.Time{}

	for k, v := range s.lastReceivedMessageByPeer {
		newMap[k] = v
	}
	return newMap
}
