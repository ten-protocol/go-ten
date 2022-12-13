package rawdb

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
)

func ReadBatch(db ethdb.KeyValueReader, hash common.L2RootHash) (*core.Batch, error) {
	header, err := readHeader(db, hash)
	if err != nil {
		return nil, fmt.Errorf("could not read header. Cause: %w", err)
	}

	body, err := ReadBody(db, hash)
	if err != nil {
		return nil, fmt.Errorf("could not read body. Cause: %w", err)
	}

	return &core.Batch{
		Header:       header,
		Transactions: body,
	}, nil
}

// ReadHeaderNumber returns the header number assigned to a hash.
func ReadHeaderNumber(db ethdb.KeyValueReader, hash gethcommon.Hash) (*uint64, error) {
	data, err := db.Get(batchHeaderNumberKey(hash))
	if err != nil {
		return nil, errutil.ErrNotFound
	}
	if len(data) != 8 {
		return nil, fmt.Errorf("header number bytes had wrong length")
	}
	number := binary.BigEndian.Uint64(data)
	return &number, nil
}

func WriteBatch(db ethdb.KeyValueWriter, batch *core.Batch) error {
	if err := writeHeader(db, batch.Header); err != nil {
		return fmt.Errorf("could not write header. Cause: %w", err)
	}
	if err := writeBody(db, *batch.Hash(), batch.Transactions); err != nil {
		return fmt.Errorf("could not write body. Cause: %w", err)
	}
	return nil
}

// Stores a batch header into the database and also stores the hash-to-number mapping.
func writeHeader(db ethdb.KeyValueWriter, header *common.Header) error {
	// Write the hash -> number mapping
	err := writeHeaderNumber(db, header.Hash(), header.Number.Uint64())
	if err != nil {
		return fmt.Errorf("could not write header number. Cause: %w", err)
	}

	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		return fmt.Errorf("could not encode batch header. Cause: %w", err)
	}
	key := batchHeaderKey(header.Hash())
	if err = db.Put(key, data); err != nil {
		return fmt.Errorf("could not put header in DB. Cause: %w", err)
	}
	return nil
}

// Stores the hash->number mapping.
func writeHeaderNumber(db ethdb.KeyValueWriter, hash gethcommon.Hash, number uint64) error {
	key := batchHeaderNumberKey(hash)
	enc := encodeBatchNumber(number)
	if err := db.Put(key, enc); err != nil {
		return fmt.Errorf("could not put header number in DB. Cause: %w", err)
	}
	return nil
}

// Retrieves the batch header corresponding to the hash.
func readHeader(db ethdb.KeyValueReader, hash common.L2RootHash) (*common.Header, error) {
	data, err := readHeaderRLP(db, hash)
	if err != nil {
		return nil, fmt.Errorf("could not read header. Cause: %w", err)
	}
	header := new(common.Header)
	if err := rlp.Decode(bytes.NewReader(data), header); err != nil {
		return nil, fmt.Errorf("could not decode batch header. Cause: %w", err)
	}
	return header, nil
}

// Retrieves a block header in its raw RLP database encoding.
func readHeaderRLP(db ethdb.KeyValueReader, hash gethcommon.Hash) (rlp.RawValue, error) {
	data, err := db.Get(batchHeaderKey(hash))
	if err != nil {
		return nil, fmt.Errorf("could not retrieve block header. Cause: %w", err)
	}
	return data, nil
}

func writeBody(db ethdb.KeyValueWriter, hash gethcommon.Hash, body []*common.L2Tx) error {
	data, err := rlp.EncodeToBytes(body)
	if err != nil {
		return fmt.Errorf("could not encode L2 transactions. Cause: %w", err)
	}
	if err = writeBodyRLP(db, hash, data); err != nil {
		return fmt.Errorf("could not write L2 transactions. Cause: %w", err)
	}
	return nil
}

// ReadBody retrieves the batch body corresponding to the hash.
func ReadBody(db ethdb.KeyValueReader, hash common.L2RootHash) ([]*common.L2Tx, error) {
	data, err := readBodyRLP(db, hash)
	if err != nil {
		return nil, fmt.Errorf("could not read body. Cause: %w", err)
	}
	body := new([]*common.L2Tx)
	if err := rlp.Decode(bytes.NewReader(data), body); err != nil {
		return nil, fmt.Errorf("could not decode L2 transactions. Cause: %w", err)
	}
	return *body, nil
}

// Stores an RLP encoded batch body into the database.
func writeBodyRLP(db ethdb.KeyValueWriter, hash common.L2RootHash, rlp rlp.RawValue) error {
	if err := db.Put(batchBodyKey(hash), rlp); err != nil {
		return fmt.Errorf("could not put batch body into DB. Cause: %w", err)
	}
	return nil
}

// Retrieves the batch body (transactions and uncles) in RLP encoding.
func readBodyRLP(db ethdb.KeyValueReader, hash common.L2RootHash) (rlp.RawValue, error) {
	data, err := db.Get(batchBodyKey(hash))
	if err != nil {
		return nil, fmt.Errorf("could not retrieve batch body from DB. Cause: %w", err)
	}
	return data, nil
}

func WriteL2Head(db ethdb.KeyValueWriter, l1Head common.L1RootHash, l2Head common.L2RootHash) error {
	if err := db.Put(headsAfterL1BlockKey(l1Head), l2Head.Bytes()); err != nil {
		return fmt.Errorf("could not put chain heads in DB. Cause: %w", err)
	}
	return nil
}

func ReadL2Head(kv ethdb.KeyValueReader, l1Head common.L1RootHash) (*common.L2RootHash, error) {
	data, err := kv.Get(headsAfterL1BlockKey(l1Head))
	if err != nil {
		return nil, errutil.ErrNotFound
	}
	l2Head := gethcommon.BytesToHash(data)
	return &l2Head, nil
}

func WriteBlockLogs(db ethdb.KeyValueWriter, blockHash gethcommon.Hash, logs []*types.Log) error {
	// Geth serialises its logs in a reduced form to minimise storage space. For now, it is more straightforward for us
	// to serialise all the fields by converting the logs to this type.
	logsForStorage := make([]*logForStorage, len(logs))
	for idx, fullFatLog := range logs {
		logsForStorage[idx] = toLogForStorage(fullFatLog)
	}

	logBytes, err := rlp.EncodeToBytes(logsForStorage)
	if err != nil {
		return fmt.Errorf("could not encode logs. Cause: %w", err)
	}

	if err := db.Put(logsKey(blockHash), logBytes); err != nil {
		return fmt.Errorf("could not put logs in DB. Cause: %w", err)
	}
	return nil
}

func ReadBlockLogs(kv ethdb.KeyValueReader, blockHash gethcommon.Hash) ([]*types.Log, error) {
	data, err := kv.Get(logsKey(blockHash))
	if err != nil {
		return nil, err
	}

	logsForStorage := new([]*logForStorage)
	if err := rlp.Decode(bytes.NewReader(data), logsForStorage); err != nil {
		return nil, fmt.Errorf("could not decode logs. Cause: %w", err)
	}

	logs := make([]*types.Log, len(*logsForStorage))
	for idx, logToStore := range *logsForStorage {
		logs[idx] = logToStore.toLog()
	}

	return logs, nil
}

// ReadCanonicalHash retrieves the hash of the canonical batch at a given height.
func ReadCanonicalHash(db ethdb.Reader, number uint64) (*common.L2RootHash, error) {
	// Get it by hash from leveldb
	data, err := db.Get(batchHeaderHashKey(number))
	if err != nil {
		return nil, errutil.ErrNotFound
	}
	hash := gethcommon.BytesToHash(data)
	return &hash, nil
}

// WriteCanonicalHash stores the hash assigned to a canonical batch number.
func WriteCanonicalHash(db ethdb.KeyValueWriter, l2Head *core.Batch) error {
	if err := db.Put(batchHeaderHashKey(l2Head.NumberU64()), l2Head.Hash().Bytes()); err != nil {
		return fmt.Errorf("failed to store number to hash mapping. Cause: %w", err)
	}
	return nil
}
