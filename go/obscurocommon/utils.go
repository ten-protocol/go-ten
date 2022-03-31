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
	Latency       func() uint64
	ScheduledFunc func()
)

func RndBtw(min uint64, max uint64) uint64 {
	if min >= max {
		panic(fmt.Sprintf("RndBtw requires min (%d) to be greater than max (%d)", min, max))
	}
	return uint64(rand.Int63n(int64(max-min))) + min //nolint:gosec
}

// ScheduleInterrupt runs the function after the delay and can be interrupted
func ScheduleInterrupt(delay uint64, interrupt *int32, fun ScheduledFunc) {
	ticker := time.NewTicker(Duration(delay))

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
	return uint64(rand.Int63n(math.MaxInt64)) //nolint:gosec
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
