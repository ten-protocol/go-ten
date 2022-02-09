package common

import (
	"fmt"
	"github.com/google/uuid"
	"math"
	"math/rand"
	"sync/atomic"
	"time"
)

type Latency func() uint64
type ScheduledFunc func()

func RndBtw(min uint64, max uint64) uint64 {
	r := uint64(rand.Int63n(int64(max-min))) + min
	return r
}

// ScheduleInterrupt runs the function after the delay and can be interrupted
func ScheduleInterrupt(delay uint64, interrupt *int32, fun ScheduledFunc) {
	ticker := time.NewTicker(Duration(delay))
	go func() {
		select {
		case <-ticker.C:
			if atomic.LoadInt32(interrupt) == 1 {
				return
			}
			fun()
		}
		ticker.Stop()
	}()
}

// Schedule runs the function after the delay
func Schedule(delay uint64, fun ScheduledFunc) {
	ticker := time.NewTicker(Duration(delay))
	go func() {
		<-ticker.C
		ticker.Stop()
		fun()
	}()
}

func Duration(us uint64) time.Duration {
	return time.Duration(us) * time.Microsecond
}

func GenerateNonce() Nonce {
	return uint64(rand.Int63n(math.MaxInt))
}

func Max(x, y uint64) uint64 {
	if x < y {
		return y
	}
	return x
}

func MaxInt(x, y uint32) uint32 {
	if x < y {
		return y
	}
	return x
}

// FindDups - returns a map of all elements that appear multiple times, and how many times
func FindDups(list []uuid.UUID) map[uuid.UUID]int {

	elementCount := make(map[uuid.UUID]int)

	for _, item := range list {
		// check if the item/element exist in the duplicate_frequency map
		_, exist := elementCount[item]
		if exist {
			elementCount[item] += 1 // increase counter by 1 if already in the map
		} else {
			elementCount[item] = 1 // else start counting from 1
		}
	}
	dups := make(map[uuid.UUID]int)
	for u, i := range elementCount {
		if i > 1 {
			dups[u] = i
			fmt.Printf("Dup: %d\n", u.ID())
		}
	}
	return dups
}

func FindTxDups(list []L1Tx) map[TxHash]int {

	elementCount := make(map[TxHash]int)

	for _, item := range list {
		// check if the item/element exist in the duplicate_frequency map
		_, exist := elementCount[item.Id]
		if exist {
			elementCount[item.Id] += 1 // increase counter by 1 if already in the map
		} else {
			elementCount[item.Id] = 1 // else start counting from 1
		}
	}
	dups := make(map[TxHash]int)
	for u, i := range elementCount {
		if i > 1 {
			dups[u] = i
			fmt.Printf(">>Dup: %d\n", u.ID())
		}
	}
	return dups
}
