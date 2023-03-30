package rawdb

import (
	"encoding/binary"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/common"
)

var (
	sharedSecret  = []byte("SharedSecret")
	headBatchHash = []byte("HeadBatch") // headBatchHashPrefix -> curr L2 head batch hash

	attestationKeyPrefix           = []byte("oAK")  // attestationKeyPrefix + address -> key
	syntheticTransactionsKeyPrefix = []byte("oSTX") // attestationKeyPrefix + address -> key

	batchHeaderPrefix            = []byte("oh")  // batchHeaderPrefix + num (uint64 big endian) + hash -> header
	batchHashSuffix              = []byte("on")  // batchHeaderPrefix + num (uint64 big endian) + headerHashSuffix -> hash
	batchBodyPrefix              = []byte("ob")  // batchBodyPrefix + num (uint64 big endian) + hash -> batch body
	batchNumberPrefix            = []byte("oH")  // batchNumberPrefix + hash -> num (uint64 big endian)
	rollupHeaderPrefix           = []byte("rh")  // rollupHeaderPrefix + num (uint64 big endian) + hash -> header
	rollupBodyPrefix             = []byte("rb")  // rollupBodyPrefix + num (uint64 big endian) + hash -> batch body
	rollupNumberPrefix           = []byte("rn")  // rollupNumberPrefix + hash -> num (uint64 big endian)
	headBatchAfterL1BlockPrefix  = []byte("hb")  // headBatchAfterL1BlockPrefix + hash -> num (uint64 big endian)
	headRollupAfterL1BlockPrefix = []byte("hr")  // headRollupAfterL1BlockPrefix + hash -> num (uint64 big endian)
	batchReceiptsPrefix          = []byte("or")  // batchReceiptsPrefix + num (uint64 big endian) + hash -> batch receipts
	contractReceiptPrefix        = []byte("ocr") // contractReceiptPrefix + address -> tx hash
	txLookupPrefix               = []byte("ol")  // txLookupPrefix + hash -> transaction/receipt lookup metadata
)

// encodeNumber encodes a number as big endian uint64
func encodeNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

// For storing and fetching a batch header by batch hash.
func batchHeaderKey(hash common.L2BatchHash) []byte {
	return append(batchHeaderPrefix, hash.Bytes()...)
}

// For storing and fetching a batch body by batch hash.
func batchBodyKey(hash common.L2BatchHash) []byte {
	return append(batchBodyPrefix, hash.Bytes()...)
}

// For storing and fetching a batch number by batch hash.
func batchNumberKey(hash common.L2BatchHash) []byte {
	return append(batchNumberPrefix, hash.Bytes()...)
}

// For storing and fetching the L2 head batch hash by L1 block hash.
func headBatchAfterL1BlockKey(hash common.L1BlockHash) []byte {
	return append(headBatchAfterL1BlockPrefix, hash.Bytes()...)
}

// For storing and fetching the canonical L2 head batch hash by height.
func batchHeaderHashKey(number uint64) []byte {
	return append(append(batchHeaderPrefix, encodeNumber(number)...), batchHashSuffix...)
}

// For storing and fetching a batch's receipts by batch hash.
func batchReceiptsKey(hash common.L2BatchHash) []byte {
	return append(batchReceiptsPrefix, hash.Bytes()...)
}

// For storing and fetching the L2 head rollup hash by L1 block hash.
func headRollupAfterL1BlockKey(hash *common.L1BlockHash) []byte {
	return append(headRollupAfterL1BlockPrefix, hash.Bytes()...)
}

func contractReceiptKey(contractAddress gethcommon.Address) []byte {
	return append(contractReceiptPrefix, contractAddress.Bytes()...)
}

// txLookupKey = txLookupPrefix + hash
func txLookupKey(hash common.TxHash) []byte {
	return append(txLookupPrefix, hash.Bytes()...)
}

func attestationPkKey(aggregator gethcommon.Address) []byte {
	return append(attestationKeyPrefix, aggregator.Bytes()...)
}

func crossChainMessagesKey(blockHash common.L1BlockHash) []byte {
	return append(syntheticTransactionsKeyPrefix, blockHash.Bytes()...)
}

// For storing and fetching a rollup header by batch hash.
func rollupHeaderKey(hash common.L2BatchHash) []byte {
	return append(rollupHeaderPrefix, hash.Bytes()...)
}

// For storing and fetching a rollup body by batch hash.
func rollupBodyKey(hash common.L2BatchHash) []byte {
	return append(rollupBodyPrefix, hash.Bytes()...)
}

// For storing and fetching a rollup number by batch hash.
func rollupNumberKey(hash common.L2BatchHash) []byte {
	return append(rollupNumberPrefix, hash.Bytes()...)
}
