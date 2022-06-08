package rawdb

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
)

var (
	sharedSecret             = []byte("SharedSecret")
	genesisRollupHash        = []byte("GenesisRollupHash")
	rollupHeaderPrefix       = []byte("oh")  // rollupHeaderPrefix + num (uint64 big endian) + hash -> header
	rollupBodyPrefix         = []byte("or")  // rollupBodyPrefix + num (uint64 big endian) + hash -> rollup body
	rollupHeaderNumberPrefix = []byte("oH")  // headerNumberPrefix + hash -> num (uint64 big endian)
	blockStatePrefix         = []byte("obs") // headerNumberPrefix + hash -> num (uint64 big endian)

)

// encodeRollupNumber encodes a rollup number as big endian uint64
func encodeRollupNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

// headerKey = rollupHeaderPrefix + num (uint64 big endian) + hash
func headerKey(number uint64, hash common.Hash) []byte {
	return append(append(rollupHeaderPrefix, encodeRollupNumber(number)...), hash.Bytes()...)
}

// headerKeyPrefix = headerPrefix + num (uint64 big endian)
func headerKeyPrefix(number uint64) []byte {
	return append(rollupHeaderPrefix, encodeRollupNumber(number)...)
}

// headerNumberKey = headerNumberPrefix + hash
func headerNumberKey(hash common.Hash) []byte {
	return append(rollupHeaderNumberPrefix, hash.Bytes()...)
}

// rollupBodyKey = rollupBodyPrefix + num (uint64 big endian) + hash
func rollupBodyKey(number uint64, hash common.Hash) []byte {
	return append(append(rollupBodyPrefix, encodeRollupNumber(number)...), hash.Bytes()...)
}

func blockStateKey(hash common.Hash) []byte {
	return append(blockStatePrefix, hash.Bytes()...)
}
