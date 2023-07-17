package rawdb

import (
	"encoding/binary"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/common"
)

var (
	rollupHeaderPrefix           = []byte("rh")  // rollupHeaderPrefix + num (uint64 big endian) + hash -> header
	rollupNumberPrefix           = []byte("rn")  // rollupNumberPrefix + hash -> num (uint64 big endian)
	headRollupAfterL1BlockPrefix = []byte("hr")  // headRollupAfterL1BlockPrefix + hash -> num (uint64 big endian)
	contractReceiptPrefix        = []byte("ocr") // contractReceiptPrefix + address -> tx hash
)

// encodeNumber encodes a number as big endian uint64
func encodeNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

// For storing and fetching the L2 head rollup hash by L1 block hash.
func headRollupAfterL1BlockKey(hash *common.L1BlockHash) []byte {
	return append(headRollupAfterL1BlockPrefix, hash.Bytes()...)
}

func contractReceiptKey(contractAddress gethcommon.Address) []byte {
	return append(contractReceiptPrefix, contractAddress.Bytes()...)
}

// For storing and fetching a rollup header by batch hash.
func rollupHeaderKey(hash common.L2BatchHash) []byte {
	return append(rollupHeaderPrefix, hash.Bytes()...)
}

// For storing and fetching a rollup number by batch hash.
func rollupNumberKey(hash common.L2BatchHash) []byte {
	return append(rollupNumberPrefix, hash.Bytes()...)
}

// for storing the contract creation count
func contractCreationCountKey() []byte {
	return []byte("contractCreationCountKey")
}
