package p2p

import (
	"sync"
	"time"
)

type item struct {
	value      int64
	lastUpdate time.Time
}

// ttlMap is a map in which the values that are not updated withing TTL are removed/reset
// this is effective for :
// - count errors, clears if no errors happen within TTL - > no errors for the TTL period
// - count messages, clears if no messages are received withing TTL -> error for the TTL period
type ttlMap struct {
	values map[string]*item
	lock   sync.RWMutex
}

func newTTLMap(ttl time.Duration) *ttlMap {
	retMap := &ttlMap{
		values: map[string]*item{},
		lock:   sync.RWMutex{},
	}

	// cleanup routine
	go func() {
		for now := range time.Tick(time.Second) {
			retMap.lock.Lock()
			for k, v := range retMap.values {
				if now.After(v.lastUpdate.Add(ttl)) {
					delete(retMap.values, k)
				}
			}
			retMap.lock.Unlock()
		}
	}()

	return retMap
}

func (t *ttlMap) increment(key string) {
	t.lock.Lock()
	defer t.lock.Unlock()

	newItem, ok := t.values[key]
	if !ok {
		newItem = &item{value: int64(0)}
		t.values[key] = newItem
	}
	newItem.value++
	newItem.lastUpdate = time.Now()
}

func (t *ttlMap) toMap() map[string]int64 {
	t.lock.RLock()
	defer t.lock.RUnlock()

	retMap := map[string]int64{}
	for k, v := range t.values {
		retMap[k] = v.value
	}
	return retMap
}
