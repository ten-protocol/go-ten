package rawdb

import (
	"bytes"
	"fmt"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
)

func ReadRollupHeader(db ethdb.KeyValueReader, hash common.L2BatchHash) (*common.RollupHeader, error) {
	return readRollupHeader(db, hash)
}

func WriteRollup(db ethdb.KeyValueWriter, rollup *common.ExtRollup) error {
	if err := writeRollupHeader(db, rollup.Header); err != nil {
		return fmt.Errorf("could not write header. Cause: %w", err)
	}
	return nil
}

func WriteL2HeadRollup(db ethdb.KeyValueWriter, l1Head *common.L1BlockHash, l2Head *common.L2BatchHash) error {
	if err := db.Put(headRollupAfterL1BlockKey(l1Head), l2Head.Bytes()); err != nil {
		return fmt.Errorf("could not put chain heads in DB. Cause: %w", err)
	}
	return nil
}

func ReadL2HeadRollup(kv ethdb.KeyValueReader, l1Head *common.L1BlockHash) (*common.L2BatchHash, error) {
	data, err := kv.Get(headRollupAfterL1BlockKey(l1Head))
	if err != nil {
		return nil, errutil.ErrNotFound
	}
	l2Head := gethcommon.BytesToHash(data)
	return &l2Head, nil
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
