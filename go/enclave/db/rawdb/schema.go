package rawdb

import (
	"encoding/binary"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/common"
)

var (
	sharedSecret = []byte("SharedSecret")

	attestationKeyPrefix           = []byte("oAK")  // attestationKeyPrefix + address -> key
	syntheticTransactionsKeyPrefix = []byte("oSTX") // attestationKeyPrefix + address -> key

	batchHeaderPrefix       = []byte("oh")  // batchHeaderPrefix + num (uint64 big endian) + hash -> header
	headerHashSuffix        = []byte("on")  // batchHeaderPrefix + num (uint64 big endian) + headerHashSuffix -> hash
	batchBodyPrefix         = []byte("ob")  // batchBodyPrefix + num (uint64 big endian) + hash -> batch body
	batchHeaderNumberPrefix = []byte("oH")  // batchHeaderNumberPrefix + hash -> num (uint64 big endian)
	headsAfterL1BlockPrefix = []byte("och") // headsAfterL1BlockPrefix + hash -> num (uint64 big endian)
	logsPrefix              = []byte("olg") // logsPrefix + hash -> block logs
	batchReceiptsPrefix     = []byte("or")  // batchReceiptsPrefix + num (uint64 big endian) + hash -> batch receipts
	contractReceiptPrefix   = []byte("ocr") // contractReceiptPrefix + address -> tx hash
	txLookupPrefix          = []byte("ol")  // txLookupPrefix + hash -> transaction/receipt lookup metadata
)

// encodeBatchNumber encodes a batch number as big endian uint64
func encodeBatchNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

// For storing and fetching a batch header by batch hash.
func batchHeaderKey(hash common.L2RootHash) []byte {
	return append(batchHeaderPrefix, hash.Bytes()...)
}

// For storing and fetching a batch body by batch hash.
func batchBodyKey(hash common.L2RootHash) []byte {
	return append(batchBodyPrefix, hash.Bytes()...)
}

// For storing and fetching a batch number by batch hash.
func batchHeaderNumberKey(hash common.L2RootHash) []byte {
	return append(batchHeaderNumberPrefix, hash.Bytes()...)
}

// For storing and fetching the L2 head hash by L1 block hash.
func headsAfterL1BlockKey(hash common.L1RootHash) []byte {
	return append(headsAfterL1BlockPrefix, hash.Bytes()...)
}

// For storing and fetching the canonical L2 head hash by height.
func batchHeaderHashKey(number uint64) []byte {
	return append(append(batchHeaderPrefix, encodeBatchNumber(number)...), headerHashSuffix...)
}

// For storing and fetching a batch's receipts by batch hash.
func batchReceiptsKey(hash common.L2RootHash) []byte {
	return append(batchReceiptsPrefix, hash.Bytes()...)
}

// logsPrefix + hash
func logsKey(hash common.L1RootHash) []byte {
	return append(logsPrefix, hash.Bytes()...)
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

func crossChainMessagesKey(blockHash common.L1RootHash) []byte {
	return append(syntheticTransactionsKeyPrefix, blockHash.Bytes()...)
}
