package rawdb

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

func ReadRollup(db ethdb.KeyValueReader, hash common.Hash) *core.Rollup {
	height := rawdb.ReadHeaderNumber(db, hash)
	if height == nil {
		return nil
	}
	return &core.Rollup{
		Header:       ReadHeader(db, hash, *height),
		Transactions: ReadBody(db, hash, *height),
	}
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
	rawdb.WriteHeaderNumber(db, hash, number)

	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		panic(err)
	}
	key := headerKey(number, hash)
	if err := db.Put(key, data); err != nil {
		panic(err)
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
		panic(err)
	}
	return header
}

// ReadHeaderRLP retrieves a block header in its raw RLP database encoding.
func ReadHeaderRLP(db ethdb.KeyValueReader, hash common.Hash, number uint64) rlp.RawValue {
	data, err := db.Get(headerKey(number, hash))
	if err != nil {
		panic(err)
	}
	return data
}

func WriteBody(db ethdb.KeyValueWriter, hash common.Hash, number uint64, body core.L2Txs) {
	data, err := rlp.EncodeToBytes(body)
	if err != nil {
		panic(err)
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
		panic(err)
	}
	return *body
}

// WriteBodyRLP stores an RLP encoded block body into the database.
func WriteBodyRLP(db ethdb.KeyValueWriter, hash common.Hash, number uint64, rlp rlp.RawValue) {
	if err := db.Put(rollupBodyKey(number, hash), rlp); err != nil {
		panic(err)
	}
}

// ReadBodyRLP retrieves the block body (transactions and uncles) in RLP encoding.
func ReadBodyRLP(db ethdb.KeyValueReader, hash common.Hash, number uint64) rlp.RawValue {
	data, err := db.Get(rollupBodyKey(number, hash))
	if err != nil {
		panic(err)
	}
	return data
}

func ReadRollupsForHeight(db ethdb.Database, number uint64) []*core.Rollup {
	hashes := rawdb.ReadAllHashes(db, number)
	rollups := make([]*core.Rollup, len(hashes))
	for i, hash := range hashes {
		rollups[i] = ReadRollup(db, hash)
	}
	return rollups
}
