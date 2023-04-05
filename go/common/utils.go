package common

import (
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"
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

func MaxInt(x, y uint32) uint32 {
	if x < y {
		return y
	}
	return x
}

// ShortHash converts the hash to a shorter uint64 for printing.
func ShortHash(hash gethcommon.Hash) uint64 {
	return hash.Big().Uint64()
}

// ShortAddress converts the address to a shorter uint64 for printing.
func ShortAddress(address gethcommon.Address) uint64 {
	return ShortHash(address.Hash())
}

// ShortNonce converts the nonce to a shorter uint64 for printing.
func ShortNonce(nonce types.BlockNonce) uint64 {
	return new(big.Int).SetBytes(nonce[4:]).Uint64()
}

// ExtractPotentialAddress - given a 32 byte hash , it checks whether it can be an address and extracts that
func ExtractPotentialAddress(hash gethcommon.Hash) *gethcommon.Address {
	bitlen := hash.Big().BitLen()
	// Addresses have 20 bytes. If the field has more, it means it is clearly not an address
	// Discovering addresses with more than 20 leading 0s is very unlikely, so we assume that
	// any topic that has less than 80 bits of data to not be an address for sure
	if bitlen < 80 || bitlen > 160 {
		return nil
	}
	a := gethcommon.BytesToAddress(hash.Bytes())
	return &a
}
