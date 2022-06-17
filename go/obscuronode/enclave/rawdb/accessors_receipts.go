package rawdb

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/log"
)

// HasReceipts verifies the existence of all the transaction receipts belonging
// to a block.
func HasReceipts(db ethdb.Reader, hash common.Hash, number uint64) bool {
	if has, err := db.Has(rollupReceiptsKey(number, hash)); !has || err != nil {
		return false
	}
	return true
}

// ReadReceiptsRLP retrieves all the transaction receipts belonging to a block in RLP encoding.
func ReadReceiptsRLP(db ethdb.Reader, hash common.Hash, number uint64) rlp.RawValue {
	data, err := db.Get(rollupReceiptsKey(number, hash))
	if err != nil {
		log.Panic("Could not read receipts.Cause: %s", err)
	}
	return data
}

// ReadRawReceipts retrieves all the transaction receipts belonging to a block.
// The receipt metadata fields are not guaranteed to be populated, so they
// should not be used. Use ReadReceipts instead if the metadata is needed.
func ReadRawReceipts(db ethdb.Reader, hash common.Hash, number uint64) types.Receipts {
	// Retrieve the flattened receipt slice
	data := ReadReceiptsRLP(db, hash, number)
	if len(data) == 0 {
		return nil
	}
	// Convert the receipts from their storage form to their internal representation
	storageReceipts := []*types.ReceiptForStorage{}
	if err := rlp.DecodeBytes(data, &storageReceipts); err != nil {
		log.Error("Invalid receipt array RLP", "hash", hash, "err", err)
		return nil
	}
	receipts := make(types.Receipts, len(storageReceipts))
	for i, storageReceipt := range storageReceipts {
		receipts[i] = (*types.Receipt)(storageReceipt)
	}
	return receipts
}

// ReadReceipts retrieves all the transaction receipts belonging to a block, including
// its correspoinding metadata fields. If it is unable to populate these metadata
// fields then nil is returned.
//
// The current implementation populates these metadata fields by reading the receipts'
// corresponding block body, so if the block body is not found it will return nil even
// if the receipt itself is stored.
func ReadReceipts(db ethdb.Reader, hash common.Hash, number uint64, config *params.ChainConfig) types.Receipts {
	// We're deriving many fields from the block body, retrieve beside the receipt
	receipts := ReadRawReceipts(db, hash, number)
	if receipts == nil {
		return nil
	}
	body := ReadBody(db, hash, number)
	if body == nil {
		log.Error("Missing body but have receipt", "hash", hash, "number", number)
		return nil
	}

	if err := receipts.DeriveFields(config, hash, number, types.Transactions(body)); err != nil {
		log.Error("Failed to derive block receipts fields", "hash", hash, "number", number, "err", err)
		return nil
	}
	return receipts
}

// WriteReceipts stores all the transaction receipts belonging to a block.
func WriteReceipts(db ethdb.KeyValueWriter, hash common.Hash, number uint64, receipts types.Receipts) {
	// Convert the receipts into their storage form and serialize them
	storageReceipts := make([]*types.ReceiptForStorage, len(receipts))
	for i, receipt := range receipts {
		storageReceipts[i] = (*types.ReceiptForStorage)(receipt)
	}
	bytes, err := rlp.EncodeToBytes(storageReceipts)
	if err != nil {
		log.Panic("Failed to encode block receipts. Cause: %s", err)
	}
	// Store the flattened receipt slice
	if err := db.Put(rollupReceiptsKey(number, hash), bytes); err != nil {
		log.Panic("Failed to store block receipts. Cause: %s", err)
	}
}

// DeleteReceipts removes all receipt data associated with a block hash.
func DeleteReceipts(db ethdb.KeyValueWriter, hash common.Hash, number uint64) {
	if err := db.Delete(rollupReceiptsKey(number, hash)); err != nil {
		log.Panic("Failed to delete block receipts. Cause: %s", err)
	}
}

// storedReceiptRLP is the storage encoding of a receipt.
// Re-definition in core/types/receipt.go.
type storedReceiptRLP struct {
	PostStateOrStatus []byte
	CumulativeGasUsed uint64
	Logs              []*types.LogForStorage
}

// ReceiptLogs is a barebone version of ReceiptForStorage which only keeps
// the list of logs. When decoding a stored receipt into this object we
// avoid creating the bloom filter.
type receiptLogs struct {
	Logs []*types.Log
}

// DecodeRLP implements rlp.Decoder.
func (r *receiptLogs) DecodeRLP(s *rlp.Stream) error {
	var stored storedReceiptRLP
	if err := s.Decode(&stored); err != nil {
		return err
	}
	r.Logs = make([]*types.Log, len(stored.Logs))
	for i, log := range stored.Logs {
		r.Logs[i] = (*types.Log)(log)
	}
	return nil
}

// DeriveLogFields fills the logs in receiptLogs with information such as block number, txhash, etc.
func deriveLogFields(receipts []*receiptLogs, hash common.Hash, number uint64, txs types.Transactions) error {
	logIndex := uint(0)
	if len(txs) != len(receipts) {
		return errors.New("transaction and receipt count mismatch")
	}
	for i := 0; i < len(receipts); i++ {
		txHash := txs[i].Hash()
		// The derived log fields can simply be set from the block and transaction
		for j := 0; j < len(receipts[i].Logs); j++ {
			receipts[i].Logs[j].BlockNumber = number
			receipts[i].Logs[j].BlockHash = hash
			receipts[i].Logs[j].TxHash = txHash
			receipts[i].Logs[j].TxIndex = uint(i)
			receipts[i].Logs[j].Index = logIndex
			logIndex++
		}
	}
	return nil
}

// ReadLogs retrieves the logs for all transactions in a block. The log fields
// are populated with metadata. In case the receipts or the block body
// are not found, a nil is returned.
func ReadLogs(db ethdb.Reader, hash common.Hash, number uint64, config *params.ChainConfig) [][]*types.Log {
	// Retrieve the flattened receipt slice
	data := ReadReceiptsRLP(db, hash, number)
	if len(data) == 0 {
		return nil
	}
	receipts := []*receiptLogs{}
	if err := rlp.DecodeBytes(data, &receipts); err != nil {
		// Receipts might be in the legacy format, try decoding that.
		// TODO: to be removed after users migrated
		if logs := readLegacyLogs(db, hash, number, config); logs != nil {
			return logs
		}
		log.Error("Invalid receipt array RLP", "hash", hash, "err", err)
		return nil
	}

	body := ReadBody(db, hash, number)
	if body == nil {
		log.Error("Missing body but have receipt", "hash", hash, "number", number)
		return nil
	}
	if err := deriveLogFields(receipts, hash, number, types.Transactions(body)); err != nil {
		log.Error("Failed to derive block receipts fields", "hash", hash, "number", number, "err", err)
		return nil
	}
	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs
}

// readLegacyLogs is a temporary workaround for when trying to read logs
// from a block which has its receipt stored in the legacy format. It'll
// be removed after users have migrated their freezer databases.
func readLegacyLogs(db ethdb.Reader, hash common.Hash, number uint64, config *params.ChainConfig) [][]*types.Log {
	receipts := ReadReceipts(db, hash, number, config)
	if receipts == nil {
		return nil
	}
	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs
}
