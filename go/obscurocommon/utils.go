package obscurocommon

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"
)

type (
	Latency       func() time.Duration
	ScheduledFunc func()
)

func RndBtwTime(min time.Duration, max time.Duration) time.Duration {
	if min <= 0 || max <= 0 {
		panic("invalid durations")
	}
	return time.Duration(RndBtw(uint64(min.Nanoseconds()), uint64(max.Nanoseconds()))) * time.Nanosecond
}

func RndBtw(min uint64, max uint64) uint64 {
	if min >= max {
		panic(fmt.Sprintf("RndBtw requires min (%d) to be greater than max (%d)", min, max))
	}
	return uint64(rand.Int63n(int64(max-min))) + min //nolint:gosec
}

// ScheduleInterrupt runs the function after the delay and can be interrupted
func ScheduleInterrupt(delay time.Duration, interrupt *int32, fun ScheduledFunc) {
	ticker := time.NewTicker(delay)

	go func() {
		<-ticker.C
		if atomic.LoadInt32(interrupt) == 1 {
			return
		}

		fun()
		ticker.Stop()
	}()
}

// Schedule runs the function after the delay
func Schedule(delay time.Duration, fun ScheduledFunc) {
	ticker := time.NewTicker(delay)
	go func() {
		<-ticker.C
		ticker.Stop()
		fun()
	}()
}

func GenerateNonce() Nonce {
	return uint64(rand.Int63n(math.MaxInt64)) //nolint:gosec
}

func MaxInt(x, y uint32) uint32 {
	if x < y {
		return y
	}
	return x
}

// TODO - Refactor the three duplicate-finding methods below to reduce duplicated code

// FindHashDups - returns a map of all hashes that appear multiple times, and how many times
func FindHashDups(list []common.Hash) map[common.Hash]int {
	elementCount := make(map[common.Hash]int)

	for _, item := range list {
		// check if the item/element exist in the duplicate_frequency map
		_, exist := elementCount[item]
		if exist {
			elementCount[item]++ // increase counter by 1 if already in the map
		} else {
			elementCount[item] = 1 // else start counting from 1
		}
	}
	dups := make(map[common.Hash]int)
	for u, i := range elementCount {
		if i > 1 {
			dups[u] = i
			fmt.Printf("Dup: %d\n", u)
		}
	}
	return dups
}

// FindRollupDups - returns a map of all L2 root hashes that appear multiple times, and how many times
func FindRollupDups(list []L2RootHash) map[L2RootHash]int {
	elementCount := make(map[L2RootHash]int)

	for _, item := range list {
		// check if the item/element exist in the duplicate_frequency map
		_, exist := elementCount[item]
		if exist {
			elementCount[item]++ // increase counter by 1 if already in the map
		} else {
			elementCount[item] = 1 // else start counting from 1
		}
	}
	dups := make(map[L2RootHash]int)
	for u, i := range elementCount {
		if i > 1 {
			dups[u] = i
			fmt.Printf("Dup: %d\n", u)
		}
	}
	return dups
}

// ShortHash converts the hash to a shorter uint64 for printing.
func ShortHash(hash common.Hash) uint64 {
	return hash.Big().Uint64()
}

// ShortAddress converts the address to a shorter uint64 for printing.
func ShortAddress(address common.Address) uint64 {
	return ShortHash(address.Hash())
}

// ShortNonce converts the nonce to a shorter uint64 for printing.
func ShortNonce(nonce types.BlockNonce) uint64 {
	return new(big.Int).SetBytes(nonce[4:]).Uint64()
}
