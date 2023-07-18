package rawdb

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
)

func ReadBatch(db ethdb.KeyValueReader, hash common.L2BatchHash) (*core.Batch, error) {
	header, err := ReadBatchHeader(db, hash)
	if err != nil {
		return nil, fmt.Errorf("could not read header. Cause: %w", err)
	}

	body, err := readBatchBody(db, hash)
	if err != nil {
		return nil, fmt.Errorf("could not read body. Cause: %w", err)
	}

	return &core.Batch{
		Header:       header,
		Transactions: body,
	}, nil
}

// ReadBatchNumber returns the number of a batch.
func ReadBatchNumber(db ethdb.KeyValueReader, hash common.L2BatchHash) (*uint64, error) {
	data, err := db.Get(batchNumberKey(hash))
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
	if err := writeBatchHeader(db, batch.Header); err != nil {
		return fmt.Errorf("could not write header. Cause: %w", err)
	}
	if err := writeBatchBody(db, batch.Hash(), batch.Transactions); err != nil {
		return fmt.Errorf("could not write body. Cause: %w", err)
	}
	return nil
}

func ReadRollupHeader(db ethdb.KeyValueReader, hash common.L2BatchHash) (*common.RollupHeader, error) {
	return readRollupHeader(db, hash)
}

func WriteRollup(db ethdb.KeyValueWriter, rollup *common.ExtRollup) error {
	if err := writeRollupHeader(db, rollup.Header); err != nil {
		return fmt.Errorf("could not write header. Cause: %w", err)
	}
	return nil
}

// Stores a batch header into the database and also stores the hash-to-number mapping.
func writeBatchHeader(db ethdb.KeyValueWriter, header *common.BatchHeader) error {
	// Write the hash -> number mapping
	h := header.Hash()
	err := writeBatchHeaderNumber(db, h, header.Number.Uint64())
	if err != nil {
		return fmt.Errorf("could not write header number. Cause: %w", err)
	}

	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		return fmt.Errorf("could not encode batch header. Cause: %w", err)
	}
	key := batchHeaderKey(h)
	if err = db.Put(key, data); err != nil {
		return fmt.Errorf("could not put header in DB. Cause: %w", err)
	}
	return nil
}

// Stores the hash->number mapping.
func writeBatchHeaderNumber(db ethdb.KeyValueWriter, hash gethcommon.Hash, number uint64) error {
	key := batchNumberKey(hash)
	enc := encodeNumber(number)
	if err := db.Put(key, enc); err != nil {
		return fmt.Errorf("could not put header number in DB. Cause: %w", err)
	}
	return nil
}

// ReadBatchHeader Retrieves the batch header corresponding to the hash.
func ReadBatchHeader(db ethdb.KeyValueReader, hash common.L2BatchHash) (*common.BatchHeader, error) {
	data, err := readBatchHeaderRLP(db, hash)
	if err != nil {
		return nil, fmt.Errorf("could not read header. Cause: %w", err)
	}
	header := new(common.BatchHeader)
	if err := rlp.Decode(bytes.NewReader(data), header); err != nil {
		return nil, fmt.Errorf("could not decode batch header. Cause: %w", err)
	}
	return header, nil
}

// Retrieves a batch header in its raw RLP database encoding.
func readBatchHeaderRLP(db ethdb.KeyValueReader, hash gethcommon.Hash) (rlp.RawValue, error) {
	data, err := db.Get(batchHeaderKey(hash))
	if err != nil {
		return nil, fmt.Errorf("could not retrieve batch header. Cause: %w", err)
	}
	return data, nil
}

func writeBatchBody(db ethdb.KeyValueWriter, hash gethcommon.Hash, body []*common.L2Tx) error {
	data, err := rlp.EncodeToBytes(body)
	if err != nil {
		return fmt.Errorf("could not encode L2 transactions. Cause: %w", err)
	}
	if err = writeBatchBodyRLP(db, hash, data); err != nil {
		return fmt.Errorf("could not write L2 transactions. Cause: %w", err)
	}
	return nil
}

// Retrieves the batch body corresponding to the hash.
func readBatchBody(db ethdb.KeyValueReader, hash common.L2BatchHash) ([]*common.L2Tx, error) {
	data, err := readBatchBodyRLP(db, hash)
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
func writeBatchBodyRLP(db ethdb.KeyValueWriter, hash common.L2BatchHash, rlp rlp.RawValue) error {
	if err := db.Put(batchBodyKey(hash), rlp); err != nil {
		return fmt.Errorf("could not put batch body into DB. Cause: %w", err)
	}
	return nil
}

// Retrieves the batch body (transactions and uncles) in RLP encoding.
func readBatchBodyRLP(db ethdb.KeyValueReader, hash common.L2BatchHash) (rlp.RawValue, error) {
	data, err := db.Get(batchBodyKey(hash))
	if err != nil {
		return nil, fmt.Errorf("could not retrieve batch body from DB. Cause: %w", err)
	}
	return data, nil
}

func SetL2HeadBatch(db ethdb.KeyValueWriter, l2Head common.L2BatchHash) error {
	if err := db.Put(headBatchHash, l2Head.Bytes()); err != nil {
		return fmt.Errorf("could not put chain heads in DB. Cause: %w", err)
	}
	return nil
}

func WriteL1ToL2BatchMapping(db ethdb.KeyValueWriter, l1Head common.L1BlockHash, l2Head common.L2BatchHash) error {
	if err := db.Put(headBatchAfterL1BlockKey(l1Head), l2Head.Bytes()); err != nil {
		return fmt.Errorf("could not put chain heads in DB. Cause: %w", err)
	}
	return nil
}

func WriteL2HeadRollup(db ethdb.KeyValueWriter, l1Head *common.L1BlockHash, l2Head *common.L2BatchHash) error {
	if err := db.Put(headRollupAfterL1BlockKey(l1Head), l2Head.Bytes()); err != nil {
		return fmt.Errorf("could not put chain heads in DB. Cause: %w", err)
	}
	return nil
}

func ReadL2HeadBatch(kv ethdb.KeyValueReader) (*common.L2BatchHash, error) {
	data, err := kv.Get(headBatchHash)
	if err != nil {
		return nil, errutil.ErrNotFound
	}
	l2Head := gethcommon.BytesToHash(data)
	return &l2Head, nil
}

func ReadL2HeadBatchForBlock(kv ethdb.KeyValueReader, l1Head common.L1BlockHash) (*common.L2BatchHash, error) {
	data, err := kv.Get(headBatchAfterL1BlockKey(l1Head))
	if err != nil {
		return nil, errutil.ErrNotFound
	}
	l2Head := gethcommon.BytesToHash(data)
	return &l2Head, nil
}

func ReadL2HeadRollup(kv ethdb.KeyValueReader, l1Head *common.L1BlockHash) (*common.L2BatchHash, error) {
	data, err := kv.Get(headRollupAfterL1BlockKey(l1Head))
	if err != nil {
		return nil, errutil.ErrNotFound
	}
	l2Head := gethcommon.BytesToHash(data)
	return &l2Head, nil
}

// ReadCanonicalBatchHash retrieves the hash of the canonical batch at a given height.
func ReadBatchBySequenceNum(db ethdb.Reader, number uint64) (*common.L2BatchHash, error) {
	// Get it by hash from leveldb
	data, err := db.Get(batchSeqHashKey(number))
	if err != nil {
		return nil, errutil.ErrNotFound
	}
	hash := gethcommon.BytesToHash(data)
	return &hash, nil
}

// WriteCanonicalHash stores the hash assigned to a canonical batch number.
func WriteBatchBySequenceNum(db ethdb.KeyValueWriter, l2Head *core.Batch) error {
	if err := db.Put(batchSeqHashKey(l2Head.Header.SequencerOrderNo.Uint64()), l2Head.Hash().Bytes()); err != nil {
		return fmt.Errorf("failed to store number to hash mapping. Cause: %w", err)
	}
	return nil
}

// ReadCanonicalBatchHash retrieves the hash of the canonical batch at a given height.
func ReadCanonicalBatchHash(db ethdb.Reader, number uint64) (*common.L2BatchHash, error) {
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

// Stores a rollup header into the database and also stores the hash-to-number mapping.
func writeRollupHeader(db ethdb.KeyValueWriter, header *common.RollupHeader) error {
	// Write the hash -> number mapping
	err := writeRollupHeaderNumber(db, header.Hash(), header.Number.Uint64())
	if err != nil {
		return fmt.Errorf("could not write header number. Cause: %w", err)
	}

	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		return fmt.Errorf("could not encode batch header. Cause: %w", err)
	}
	key := rollupHeaderKey(header.Hash())
	if err = db.Put(key, data); err != nil {
		return fmt.Errorf("could not put header in DB. Cause: %w", err)
	}
	return nil
}

// Stores the hash->number mapping.
func writeRollupHeaderNumber(db ethdb.KeyValueWriter, hash gethcommon.Hash, number uint64) error {
	key := rollupNumberKey(hash)
	enc := encodeNumber(number)
	if err := db.Put(key, enc); err != nil {
		return fmt.Errorf("could not put rollup header number in DB. Cause: %w", err)
	}
	return nil
}

// Retrieves the rollup header corresponding to the hash.
func readRollupHeader(db ethdb.KeyValueReader, hash common.L2BatchHash) (*common.RollupHeader, error) {
	data, err := readRollupHeaderRLP(db, hash)
	if err != nil {
		return nil, fmt.Errorf("could not read header. Cause: %w", err)
	}
	header := new(common.RollupHeader)
	if err := rlp.Decode(bytes.NewReader(data), header); err != nil {
		return nil, fmt.Errorf("could not decode rollup header. Cause: %w", err)
	}
	return header, nil
}

// Retrieves a rollup header in its raw RLP database encoding.
func readRollupHeaderRLP(db ethdb.KeyValueReader, hash gethcommon.Hash) (rlp.RawValue, error) {
	data, err := db.Get(rollupHeaderKey(hash))
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rollup header. Cause: %w", err)
	}
	return data, nil
}
