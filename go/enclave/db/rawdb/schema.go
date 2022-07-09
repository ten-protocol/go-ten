package rawdb

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
)

var (
	sharedSecret = []byte("SharedSecret")

	attestationKeyPrefix = []byte("oAK") // attestationKeyPrefix + address -> key

	genesisRollupHash        = []byte("GenesisRollupHash")
	rollupHeaderPrefix       = []byte("oh")          // rollupHeaderPrefix + num (uint64 big endian) + hash -> header
	headerHashSuffix         = []byte("on")          // headerPrefix + num (uint64 big endian) + headerHashSuffix -> hash
	rollupBodyPrefix         = []byte("ob")          // rollupBodyPrefix + num (uint64 big endian) + hash -> rollup body
	rollupHeaderNumberPrefix = []byte("oH")          // headerNumberPrefix + hash -> num (uint64 big endian)
	blockStatePrefix         = []byte("obs")         // headerNumberPrefix + hash -> num (uint64 big endian)
	rollupReceiptsPrefix     = []byte("or")          // rollupReceiptsPrefix + num (uint64 big endian) + hash -> block receipts
	txLookupPrefix           = []byte("ol")          // txLookupPrefix + hash -> transaction/receipt lookup metadata
	bloomBitsPrefix          = []byte("oB")          // bloomBitsPrefix + bit (uint16 big endian) + section (uint64 big endian) + hash -> bloom bits
	headRollupKey            = []byte("oLastBlock")  // headRollupKey tracks the latest known full block's hash.
	headHeaderKey            = []byte("oLastHeader") // headHeaderKey tracks the latest known header's hash.
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

// rollupReceiptsKey = rollupReceiptsPrefix + num (uint64 big endian) + hash
func rollupReceiptsKey(number uint64, hash common.Hash) []byte {
	return append(append(rollupReceiptsPrefix, encodeRollupNumber(number)...), hash.Bytes()...)
}

// txLookupKey = txLookupPrefix + hash
func txLookupKey(hash common.Hash) []byte {
	return append(txLookupPrefix, hash.Bytes()...)
}

// bloomBitsKey = bloomBitsPrefix + bit (uint16 big endian) + section (uint64 big endian) + hash
func bloomBitsKey(bit uint, section uint64, hash common.Hash) []byte {
	key := append(append(bloomBitsPrefix, make([]byte, 10)...), hash.Bytes()...)

	binary.BigEndian.PutUint16(key[1:], uint16(bit))
	binary.BigEndian.PutUint64(key[3:], section)

	return key
}

// headerHashKey = headerPrefix + num (uint64 big endian) + headerHashSuffix
func headerHashKey(number uint64) []byte {
	return append(append(rollupHeaderPrefix, encodeRollupNumber(number)...), headerHashSuffix...)
}

func attestationPkKey(aggregator common.Address) []byte {
	return append(attestationKeyPrefix, aggregator.Bytes()...)
}
