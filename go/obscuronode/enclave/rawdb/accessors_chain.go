package rawdb

import (
	"bytes"
	"encoding/binary"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

func ReadRollup(db ethdb.KeyValueReader, hash common.Hash) *core.Rollup {
	height := ReadHeaderNumber(db, hash)
	if height == nil {
		return nil
	}
	return &core.Rollup{
		Header:       ReadHeader(db, hash, *height),
		Transactions: ReadBody(db, hash, *height),
	}
}

// ReadHeaderNumber returns the header number assigned to a hash.
func ReadHeaderNumber(db ethdb.KeyValueReader, hash common.Hash) *uint64 {
	data, _ := db.Get(headerNumberKey(hash))
	if len(data) != 8 {
		return nil
	}
	number := binary.BigEndian.Uint64(data)
	return &number
}

func WriteRollup(db ethdb.KeyValueWriter, rollup *core.Rollup) {
	WriteHeader(db, rollup.Header)
	WriteBody(db, rollup.Hash(), rollup.Header.Number, rollup.Transactions)
}

// WriteHeader stores a rollup header into the database and also stores the hash-
// to-number mapping.
func WriteHeader(db ethdb.KeyValueWriter, header *nodecommon.Header) {
	var (
		hash   = header.Hash()
		number = header.Number
	)
	// Write the hash -> number mapping
	WriteHeaderNumber(db, hash, number)

	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		log.Panic("could not encode rollup header. Cause: %s", err)
	}
	key := headerKey(number, hash)
	if err := db.Put(key, data); err != nil {
		log.Panic("could not put header in DB. Cause: %s", err)
	}
}

// WriteHeaderNumber stores the hash->number mapping.
func WriteHeaderNumber(db ethdb.KeyValueWriter, hash common.Hash, number uint64) {
	key := headerNumberKey(hash)
	enc := encodeRollupNumber(number)
	if err := db.Put(key, enc); err != nil {
		log.Panic("could not put header number in DB. Cause: %s", err)
	}
}

// ReadHeader retrieves the rollup header corresponding to the hash.
func ReadHeader(db ethdb.KeyValueReader, hash common.Hash, number uint64) *nodecommon.Header {
	data := ReadHeaderRLP(db, hash, number)
	if len(data) == 0 {
		return nil
	}
	header := new(nodecommon.Header)
	if err := rlp.Decode(bytes.NewReader(data), header); err != nil {
		log.Panic("could not decode rollup header. Cause: %s", err)
	}
	return header
}

// ReadHeaderRLP retrieves a block header in its raw RLP database encoding.
func ReadHeaderRLP(db ethdb.KeyValueReader, hash common.Hash, number uint64) rlp.RawValue {
	data, err := db.Get(headerKey(number, hash))
	if err != nil {
		log.Panic("could not retrieve block header. Cause: %s", err)
	}
	return data
}

func WriteBody(db ethdb.KeyValueWriter, hash common.Hash, number uint64, body core.L2Txs) {
	data, err := rlp.EncodeToBytes(body)
	if err != nil {
		log.Panic("could not encode L2 transactions. Cause: %s", err)
	}
	WriteBodyRLP(db, hash, number, data)
}

// ReadBody retrieves the rollup body corresponding to the hash.
func ReadBody(db ethdb.KeyValueReader, hash common.Hash, number uint64) core.L2Txs {
	data := ReadBodyRLP(db, hash, number)
	if len(data) == 0 {
		return nil
	}
	body := new(core.L2Txs)
	if err := rlp.Decode(bytes.NewReader(data), body); err != nil {
		log.Panic("could not decode L2 transactions. Cause: %s", err)
	}
	return *body
}

// WriteBodyRLP stores an RLP encoded block body into the database.
func WriteBodyRLP(db ethdb.KeyValueWriter, hash common.Hash, number uint64, rlp rlp.RawValue) {
	if err := db.Put(rollupBodyKey(number, hash), rlp); err != nil {
		log.Panic("could not put block body into DB. Cause: %s", err)
	}
}

// ReadBodyRLP retrieves the block body (transactions and uncles) in RLP encoding.
func ReadBodyRLP(db ethdb.KeyValueReader, hash common.Hash, number uint64) rlp.RawValue {
	data, err := db.Get(rollupBodyKey(number, hash))
	if err != nil {
		log.Panic("could not retrieve block body from DB. Cause: %s", err)
	}
	return data
}

func ReadRollupsForHeight(db ethdb.Database, number uint64) []*core.Rollup {
	hashes := ReadAllHashes(db, number)
	rollups := make([]*core.Rollup, len(hashes))
	for i, hash := range hashes {
		rollups[i] = ReadRollup(db, hash)
	}
	return rollups
}

// ReadAllHashes retrieves all the hashes assigned to blocks at a certain heights,
// both canonical and reorged forks included.
func ReadAllHashes(db ethdb.Iteratee, number uint64) []common.Hash {
	prefix := headerKeyPrefix(number)

	hashes := make([]common.Hash, 0, 1)
	it := db.NewIterator(prefix, nil)
	defer it.Release()

	for it.Next() {
		if key := it.Key(); len(key) == len(prefix)+32 {
			hashes = append(hashes, common.BytesToHash(key[len(key)-32:]))
		}
	}
	return hashes
}

func WriteBlockState(db ethdb.KeyValueWriter, bs *core.BlockState) {
	bytes, err := rlp.EncodeToBytes(bs)
	if err != nil {
		log.Panic("could not encode block state. Cause: %s", err)
	}
	if err := db.Put(blockStateKey(bs.Block), bytes); err != nil {
		log.Panic("could not put block state in DB. Cause: %s", err)
	}
}

func ReadBlockState(kv ethdb.KeyValueReader, hash common.Hash) *core.BlockState {
	data, _ := kv.Get(blockStateKey(hash))
	if data == nil {
		return nil
	}
	bs := new(core.BlockState)
	if err := rlp.Decode(bytes.NewReader(data), bs); err != nil {
		log.Panic("could not decode block state. Cause: %s", err)
	}
	return bs
}
