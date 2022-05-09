package rawdb

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
)

var (
	sharedSecret      = []byte("SharedSecret")
	genesisRollupHash = []byte("GenesisRollupHash")
	headerPrefix      = []byte("h") // headerPrefix + num (uint64 big endian) + hash -> header
	rollupBodyPrefix  = []byte("b") // rollupBodyPrefix + num (uint64 big endian) + hash -> rollup body

)

// encodeRollupNumber encodes a rollup number as big endian uint64
func encodeRollupNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

// headerKey = headerPrefix + num (uint64 big endian) + hash
func headerKey(number uint64, hash common.Hash) []byte {
	return append(append(headerPrefix, encodeRollupNumber(number)...), hash.Bytes()...)
}

// rollupBodyKey = rollupBodyPrefix + num (uint64 big endian) + hash
func rollupBodyKey(number uint64, hash common.Hash) []byte {
	return append(append(rollupBodyPrefix, encodeRollupNumber(number)...), hash.Bytes()...)
}
