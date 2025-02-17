package async

import (
	"sync"
	"time"
)

// Timestamp is a thread safe timestamp
type Timestamp struct {
	lastTimestamp time.Time
	mutex         sync.RWMutex
}

func NewAsyncTimestamp(lastTimestamp time.Time) *Timestamp {
	return &Timestamp{
		lastTimestamp: lastTimestamp,
		mutex:         sync.RWMutex{},
	}
}

// Mark sets the timestamp with the current time
func (at *Timestamp) Mark() {
	at.mutex.Lock()
	defer at.mutex.Unlock()
	at.lastTimestamp = time.Now()
}

// LastTimestamp returns the last set timestamp
func (at *Timestamp) LastTimestamp() time.Time {
	at.mutex.RLock()
	defer at.mutex.RUnlock()

	newTimestamp := at.lastTimestamp
	return newTimestamp
}
