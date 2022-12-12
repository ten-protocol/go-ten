package p2p

import (
	"sync"
	"time"
)

const (
	_failedMessageRead        = "Failed Message Reads"
	_failedMessageDecode      = "Failed Message Decodes"
	_failedConnectSendMessage = "Failed Peer Connects"
	_failedWriteSendMessage   = "Failed Socket Writes"
)

type status struct {
	lock      sync.RWMutex
	timestamp time.Time
	failures  map[string]map[string]int64
}

func newStatus() *status {
	return &status{
		lock:     sync.RWMutex{},
		failures: map[string]map[string]int64{},
	}
}

func (s *status) increment(failType string, host string) {
	// don't wait for locks - fire and forget
	go func() {
		s.lock.Lock()
		defer s.lock.Unlock()

		// only keep track for 10 min periods at a time
		if time.Now().After(s.timestamp.Add(10 * time.Minute)) {
			s.timestamp = time.Now()
			s.failures = map[string]map[string]int64{}
		}

		if _, ok := s.failures[failType]; !ok {
			s.failures[failType] = map[string]int64{}
		}

		s.failures[failType][host] += 1
	}()
}

func (s *status) status() map[string]map[string]int64 {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.failures
}
