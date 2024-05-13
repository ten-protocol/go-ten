package common

import (
	"math/big"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

type (
	Latency func() time.Duration
)

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
	return address.Big().Uint64()
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

// Generates a random string n characters long.
func RandomStr(n int) string {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	suffix := make([]rune, n)
	for i := range suffix {
		suffix[i] = letters[randGen.Intn(len(letters))]
	}
	return string(suffix)
}
