package common

import (
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
