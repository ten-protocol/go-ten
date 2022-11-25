package rawdb

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/status-im/keycard-go/hexutils"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
)

// todo - all the function in this file should return an error, which must be handled by the caller
//  once that is done, the logger parameter should be removed

func ReadRollup(db ethdb.KeyValueReader, hash gethcommon.Hash, logger gethlog.Logger) (*core.Rollup, error) {
	height, err := ReadHeaderNumber(db, hash)
	if err != nil {
		return nil, err
	}
	return &core.Rollup{
		Header:       ReadHeader(db, hash, *height, logger),
		Transactions: ReadBody(db, hash, *height, logger),
	}, nil
}

// ReadHeaderNumber returns the header number assigned to a hash.
func ReadHeaderNumber(db ethdb.KeyValueReader, hash gethcommon.Hash) (*uint64, error) {
	data, err := db.Get(headerNumberKey(hash))
	if err != nil {
		return nil, err
	}
	if len(data) != 8 {
		return nil, fmt.Errorf("header number bytes had wrong length")
	}
	number := binary.BigEndian.Uint64(data)
	return &number, nil
}

func WriteRollup(db ethdb.KeyValueWriter, rollup *core.Rollup) error {
	if err := WriteHeader(db, rollup.Header); err != nil {
		return fmt.Errorf("could not write header. Cause: %w", err)
	}
	if err := WriteBody(db, rollup.Hash(), rollup.Header.Number.Uint64(), rollup.Transactions); err != nil {
		return fmt.Errorf("could not write body. Cause: %w", err)
	}
	return nil
}

// WriteHeader stores a rollup header into the database and also stores the hash-to-number mapping.
func WriteHeader(db ethdb.KeyValueWriter, header *common.Header) error {
	hash := header.Hash()
	number := header.Number.Uint64()

	// Write the hash -> number mapping
	err := WriteHeaderNumber(db, hash, number)
	if err != nil {
		return fmt.Errorf("could not write header number. Cause: %w", err)
	}

	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		return fmt.Errorf("could not encode rollup header. Cause: %w", err)
	}
	key := headerKey(number, hash)
	if err = db.Put(key, data); err != nil {
		return fmt.Errorf("could not put header in DB. Cause: %w", err)
	}
	return nil
}

// WriteHeaderNumber stores the hash->number mapping.
func WriteHeaderNumber(db ethdb.KeyValueWriter, hash gethcommon.Hash, number uint64) error {
	key := headerNumberKey(hash)
	enc := encodeRollupNumber(number)
	if err := db.Put(key, enc); err != nil {
		return fmt.Errorf("could not put header number in DB. Cause: %w", err)
	}
	return nil
}

// ReadHeader retrieves the rollup header corresponding to the hash.
func ReadHeader(db ethdb.KeyValueReader, hash gethcommon.Hash, number uint64, logger gethlog.Logger) *common.Header {
	data := ReadHeaderRLP(db, hash, number, logger)
	if len(data) == 0 {
		return nil
	}
	header := new(common.Header)
	if err := rlp.Decode(bytes.NewReader(data), header); err != nil {
		logger.Crit("could not decode rollup header. ", log.ErrKey, err)
	}
	return header
}

// ReadHeaderRLP retrieves a block header in its raw RLP database encoding.
func ReadHeaderRLP(db ethdb.KeyValueReader, hash gethcommon.Hash, number uint64, logger gethlog.Logger) rlp.RawValue {
	data, err := db.Get(headerKey(number, hash))
	if err != nil {
		logger.Crit("could not retrieve block header. ", log.ErrKey, err)
	}
	return data
}

func WriteBody(db ethdb.KeyValueWriter, hash gethcommon.Hash, number uint64, body []*common.L2Tx) error {
	data, err := rlp.EncodeToBytes(body)
	if err != nil {
		return fmt.Errorf("could not encode L2 transactions. Cause: %w", err)
	}
	if err = WriteBodyRLP(db, hash, number, data); err != nil {
		return fmt.Errorf("could not write L2 transactions. Cause: %w", err)
	}
	return nil
}

// ReadBody retrieves the rollup body corresponding to the hash.
func ReadBody(db ethdb.KeyValueReader, hash gethcommon.Hash, number uint64, logger gethlog.Logger) []*common.L2Tx {
	data := ReadBodyRLP(db, hash, number, logger)
	if len(data) == 0 {
		return nil
	}
	body := new([]*common.L2Tx)
	if err := rlp.Decode(bytes.NewReader(data), body); err != nil {
		logger.Crit("could not decode L2 transactions. ", log.ErrKey, err)
	}
	return *body
}

// WriteBodyRLP stores an RLP encoded block body into the database.
func WriteBodyRLP(db ethdb.KeyValueWriter, hash gethcommon.Hash, number uint64, rlp rlp.RawValue) error {
	if err := db.Put(rollupBodyKey(number, hash), rlp); err != nil {
		return fmt.Errorf("could not put rollup body into DB. Cause: %w", err)
	}
	return nil
}

// ReadBodyRLP retrieves the block body (transactions and uncles) in RLP encoding.
func ReadBodyRLP(db ethdb.KeyValueReader, hash gethcommon.Hash, number uint64, logger gethlog.Logger) rlp.RawValue {
	data, err := db.Get(rollupBodyKey(number, hash))
	if err != nil {
		logger.Crit(fmt.Sprintf("could not retrieve rollup body :r_%d from DB. ", common.ShortHash(hash)), "key", hexutils.BytesToHex(rollupBodyKey(number, hash)), log.ErrKey, err)
	}
	return data
}

func ReadRollupsForHeight(db ethdb.Database, number uint64, logger gethlog.Logger) ([]*core.Rollup, error) {
	hashes := ReadAllHashes(db, number)
	rollups := make([]*core.Rollup, len(hashes))
	for i, hash := range hashes {
		rollup, err := ReadRollup(db, hash, logger)
		if err != nil {
			return nil, err
		}
		rollups[i] = rollup
	}
	return rollups, nil
}

// ReadAllHashes retrieves all the hashes assigned to blocks at a certain heights,
// both canonical and reorged forks included.
func ReadAllHashes(db ethdb.Iteratee, number uint64) []gethcommon.Hash {
	prefix := headerKeyPrefix(number)

	hashes := make([]gethcommon.Hash, 0, 1)
	it := db.NewIterator(prefix, nil)
	defer it.Release()

	for it.Next() {
		if key := it.Key(); len(key) == len(prefix)+32 {
			hashes = append(hashes, gethcommon.BytesToHash(key[len(key)-32:]))
		}
	}
	return hashes
}

func WriteBlockState(db ethdb.KeyValueWriter, bs *core.BlockState, logger gethlog.Logger) {
	blockStateBytes, err := rlp.EncodeToBytes(bs)
	if err != nil {
		logger.Crit("could not encode block state. ", log.ErrKey, err)
	}
	if err := db.Put(blockStateKey(bs.Block), blockStateBytes); err != nil {
		logger.Crit("could not put block state in DB. ", log.ErrKey, err)
	}
}

func ReadBlockState(kv ethdb.KeyValueReader, hash gethcommon.Hash) (*core.BlockState, error) {
	// TODO - Handle error.
	data, _ := kv.Get(blockStateKey(hash))
	if data == nil {
		return nil, errutil.ErrNotFound
	}
	bs := new(core.BlockState)
	if err := rlp.Decode(bytes.NewReader(data), bs); err != nil {
		return nil, fmt.Errorf("could not decode block state. Cause: %w", err)
	}
	return bs, nil
}

func WriteBlockLogs(db ethdb.KeyValueWriter, blockHash gethcommon.Hash, logs []*types.Log, logger gethlog.Logger) {
	// Geth serialises its logs in a reduced form to minimise storage space. For now, it is more straightforward for us
	// to serialise all the fields by converting the logs to this type.
	logsForStorage := make([]*logForStorage, len(logs))
	for idx, fullFatLog := range logs {
		logsForStorage[idx] = toLogForStorage(fullFatLog)
	}

	logBytes, err := rlp.EncodeToBytes(logsForStorage)
	if err != nil {
		logger.Crit("could not encode logs. ", log.ErrKey, err)
	}

	if err := db.Put(logsKey(blockHash), logBytes); err != nil {
		logger.Crit("could not put logs in DB. ", log.ErrKey, err)
	}
}

func ReadBlockLogs(kv ethdb.KeyValueReader, blockHash gethcommon.Hash, logger gethlog.Logger) []*types.Log {
	data, _ := kv.Get(logsKey(blockHash))
	if data == nil {
		return nil
	}

	logsForStorage := new([]*logForStorage)
	if err := rlp.Decode(bytes.NewReader(data), logsForStorage); err != nil {
		logger.Crit("could not decode logs. ", log.ErrKey, err)
	}

	logs := make([]*types.Log, len(*logsForStorage))
	for idx, logToStore := range *logsForStorage {
		logs[idx] = logToStore.toLog()
	}

	return logs
}

// ReadCanonicalHash retrieves the hash assigned to a canonical block number.
func ReadCanonicalHash(db ethdb.Reader, number uint64) (*gethcommon.Hash, error) {
	// Get it by hash from leveldb
	data, err := db.Get(headerHashKey(number))
	if err != nil {
		return nil, err
	}
	hash := gethcommon.BytesToHash(data)
	return &hash, nil
}

// WriteCanonicalHash stores the hash assigned to a canonical block number.
func WriteCanonicalHash(db ethdb.KeyValueWriter, hash gethcommon.Hash, number uint64, logger gethlog.Logger) {
	if err := db.Put(headerHashKey(number), hash.Bytes()); err != nil {
		logger.Crit("Failed to store number to hash mapping. ", log.ErrKey, err)
	}
}

// DeleteCanonicalHash removes the number to hash canonical mapping.
func DeleteCanonicalHash(db ethdb.KeyValueWriter, number uint64) error {
	if err := db.Delete(headerHashKey(number)); err != nil {
		return fmt.Errorf("failed to delete number to hash mapping. Cause: %w", err)
	}
	return nil
}

// ReadHeadRollupHash retrieves the hash of the current canonical head block.
func ReadHeadRollupHash(db ethdb.KeyValueReader) (*gethcommon.Hash, error) {
	data, err := db.Get(headRollupKey)
	if err != nil {
		return nil, err
	}
	hash := gethcommon.BytesToHash(data)
	return &hash, nil
}

// WriteHeadRollupHash stores the head block's hash.
func WriteHeadRollupHash(db ethdb.KeyValueWriter, hash gethcommon.Hash) error {
	if err := db.Put(headRollupKey, hash.Bytes()); err != nil {
		return fmt.Errorf("failed to store last block's hash. Cause: %w", err)
	}
	return nil
}

// ReadHeadHeaderHash retrieves the hash of the current canonical head header.
func ReadHeadHeaderHash(db ethdb.KeyValueReader) (*gethcommon.Hash, error) {
	data, err := db.Get(headHeaderKey)
	if err != nil {
		return nil, err
	}
	hash := gethcommon.BytesToHash(data)
	return &hash, nil
}

// WriteHeadHeaderHash stores the hash of the current canonical head header.
func WriteHeadHeaderHash(db ethdb.KeyValueWriter, hash gethcommon.Hash, logger gethlog.Logger) {
	if err := db.Put(headHeaderKey, hash.Bytes()); err != nil {
		logger.Crit("Failed to store last header's hash. ", log.ErrKey, err)
	}
}
