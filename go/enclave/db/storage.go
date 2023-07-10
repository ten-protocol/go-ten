package db

import (
	"bytes"
	"crypto/ecdsa"
	sql2 "database/sql"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/rlp"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/db/orm"

	gethcore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/trie"

	"github.com/obscuronet/go-obscuro/go/common/syserr"

	"github.com/obscuronet/go-obscuro/go/enclave/db/sql"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/tracers"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	obscurorawdb "github.com/obscuronet/go-obscuro/go/enclave/db/rawdb"
)

// todo - this will require a dedicated table when updates are implemented
const master_seed_cfg = "MASTER_SEED"

// used only by sequencers
const current_sequence = "CURRENT_SEQ"

// ErrNoRollups is returned if no rollups have been published yet in the history of the network
// Note: this is not just "not found", we cache at every L1 block what rollup we are up to so we also record that we haven't seen one yet
var ErrNoRollups = errors.New("no rollups have been published")

// todo (#1551) - consistency around whether we assert the secret is available or not

type storageImpl struct {
	db          *sql.EnclaveDB
	stateDB     state.Database
	chainConfig *params.ChainConfig
	logger      gethlog.Logger
}

func NewStorage(backingDB *sql.EnclaveDB, chainConfig *params.ChainConfig, logger gethlog.Logger) Storage {
	cacheConfig := &gethcore.CacheConfig{
		TrieCleanLimit: 256,
		TrieDirtyLimit: 256,
		TrieTimeLimit:  5 * time.Minute,
		SnapshotLimit:  256,
		SnapshotWait:   true,
	}

	return &storageImpl{
		db: backingDB,
		stateDB: state.NewDatabaseWithConfig(backingDB, &trie.Config{
			Cache:     cacheConfig.TrieCleanLimit,
			Journal:   cacheConfig.TrieCleanJournal,
			Preimages: cacheConfig.Preimages,
		}),
		chainConfig: chainConfig,
		logger:      logger,
	}
}

func (s *storageImpl) TrieDB() *trie.Database {
	return s.stateDB.TrieDB()
}

func (s *storageImpl) OpenBatch() *sql.Batch {
	return s.db.NewSQLBatch()
}

func (s *storageImpl) CommitBatch(dbBatch *sql.Batch) error {
	return dbBatch.Write()
}

func (s *storageImpl) Close() error {
	return s.db.GetSQLDB().Close()
}

func (s *storageImpl) FetchHeadBatch() (*core.Batch, error) {
	return orm.FetchHeadBatch(s.db.GetSQLDB())
}

func (s *storageImpl) FetchCurrentSequencerNo() (*big.Int, error) {
	return orm.FetchCurrentSequencerNo(s.db.GetSQLDB())
}

func (s *storageImpl) FetchBatch(hash common.L2BatchHash) (*core.Batch, error) {
	return orm.FetchBatch(s.db.GetSQLDB(), hash)
}

func (s *storageImpl) FetchBatchHeader(hash common.L2BatchHash) (*common.BatchHeader, error) {
	return orm.ReadBatchHeader(s.db.GetSQLDB(), hash)
}

func (s *storageImpl) FetchBatchByHeight(height uint64) (*core.Batch, error) {
	return orm.FetchCanonicalBatchByHeight(s.db.GetSQLDB(), height)
}

func (s *storageImpl) StoreBlock(b *types.Block) {
	rawdb.WriteBlock(s.db, b)
}

func (s *storageImpl) FetchBlock(blockHash common.L1BlockHash) (*types.Block, error) {
	height := rawdb.ReadHeaderNumber(s.db, blockHash)
	if height == nil {
		return nil, errutil.ErrNotFound
	}
	b := rawdb.ReadBlock(s.db, blockHash, *height)
	if b == nil {
		return nil, errutil.ErrNotFound
	}
	return b, nil
}

func (s *storageImpl) FetchHeadBlock() (*types.Block, error) {
	block, err := s.FetchBlock(rawdb.ReadHeadHeaderHash(s.db))
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (s *storageImpl) StoreSecret(secret crypto.SharedEnclaveSecret) error {
	enc, err := rlp.EncodeToBytes(secret)
	if err != nil {
		return fmt.Errorf("could not encode shared secret. Cause: %w", err)
	}
	_, err = orm.WriteConfig(s.db.GetSQLDB(), master_seed_cfg, enc)
	if err != nil {
		return fmt.Errorf("could not shared secret in DB. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) FetchSecret() (*crypto.SharedEnclaveSecret, error) {
	var ss crypto.SharedEnclaveSecret

	cfg, err := orm.FetchConfig(s.db.GetSQLDB(), master_seed_cfg)
	if err != nil {
		return nil, err
	}
	if err := rlp.DecodeBytes(cfg, &ss); err != nil {
		return nil, fmt.Errorf("could not decode shared secret")
	}

	return &ss, nil
}

func (s *storageImpl) IsAncestor(block *types.Block, maybeAncestor *types.Block) bool {
	if bytes.Equal(maybeAncestor.Hash().Bytes(), block.Hash().Bytes()) {
		return true
	}

	if maybeAncestor.NumberU64() >= block.NumberU64() {
		return false
	}

	p, err := s.FetchBlock(block.ParentHash())
	if err != nil {
		return false
	}

	return s.IsAncestor(p, maybeAncestor)
}

func (s *storageImpl) IsBlockAncestor(block *types.Block, maybeAncestor common.L1BlockHash) bool {
	resolvedBlock, err := s.FetchBlock(maybeAncestor)
	if err != nil {
		return false
	}
	return s.IsAncestor(block, resolvedBlock)
}

func (s *storageImpl) HealthCheck() (bool, error) {
	headBatch, err := s.FetchHeadBatch()
	if err != nil {
		s.logger.Error("unable to HealthCheck storage", log.ErrKey, err)
		return false, err
	}
	return headBatch != nil, nil
}

func (s *storageImpl) FetchHeadBatchForBlock(blockHash common.L1BlockHash) (*core.Batch, error) {
	l2HeadBatch, err := obscurorawdb.ReadL2HeadBatchForBlock(s.db, blockHash)
	if err != nil {
		return nil, fmt.Errorf("could not read L2 head batch for block. Cause: %w", err)
	}
	return obscurorawdb.ReadBatch(s.db, *l2HeadBatch)
}

func (s *storageImpl) FetchHeadRollupForBlock(blockHash *common.L1BlockHash) (*common.RollupHeader, error) {
	l2HeadBatch, err := obscurorawdb.ReadL2HeadRollup(s.db, blockHash)
	if err != nil {
		return nil, fmt.Errorf("could not read L2 head rollup for block. Cause: %w", err)
	}
	if *l2HeadBatch == (gethcommon.Hash{}) { // empty hash ==> no rollups yet up to this block
		return nil, ErrNoRollups
	}
	return obscurorawdb.ReadRollupHeader(s.db, *l2HeadBatch)
}

func (s *storageImpl) UpdateHeadBatch(l1Head common.L1BlockHash, l2Head *core.Batch, receipts []*types.Receipt, dbBatch *sql.Batch) error {
	if dbBatch == nil {
		panic("UpdateHeadBatch called without an instance of sql.Batch")
	}

	if err := obscurorawdb.SetL2HeadBatch(dbBatch, l2Head.Hash()); err != nil {
		return fmt.Errorf("could not write block state. Cause: %w", err)
	}
	if err := obscurorawdb.WriteL1ToL2BatchMapping(dbBatch, l1Head, l2Head.Hash()); err != nil {
		return fmt.Errorf("could not write block state. Cause: %w", err)
	}

	// We update the canonical hash of the batch at this height.
	if err := obscurorawdb.WriteCanonicalHash(dbBatch, l2Head); err != nil {
		return fmt.Errorf("could not write canonical hash. Cause: %w", err)
	}

	if l2Head.Number().Int64() > 1 {
		err2 := s.writeLogs(l2Head.Header.ParentHash, receipts, dbBatch)
		if err2 != nil {
			return fmt.Errorf("could not save logs %w", err2)
		}
	}
	return nil
}

func (s *storageImpl) writeLogs(l2Head common.L2BatchHash, receipts []*types.Receipt, dbBatch *sql.Batch) error {
	stateDB, err := s.CreateStateDB(l2Head)
	if err != nil {
		return fmt.Errorf("could not create state DB to filter logs. Cause: %w", err)
	}

	// We update the block's logs, based on the batch's logs.
	for _, receipt := range receipts {
		for _, l := range receipt.Logs {
			s.writeLog(l, stateDB, dbBatch)
		}
	}
	return nil
}

// This method stores a log entry together with relevancy metadata
// Each types.Log has 5 indexable topics, where the first one is the event signature hash
// The other 4 topics are set by the programmer
// According to the data relevancy rules, an event is relevant to accounts referenced directly in topics
// If the event is not referring any user address, it is considered a "lifecycle event", and is relevant to everyone
func (s *storageImpl) writeLog(l *types.Log, stateDB *state.StateDB, dbBatch *sql.Batch) {
	// The topics are stored in an array with a maximum of 5 entries, but usually less
	var t0, t1, t2, t3, t4 *gethcommon.Hash

	// these are the addresses to which this event might be relevant to.
	var addr1, addr2, addr3, addr4 *gethcommon.Address

	// start with true, and as soon as a user address is discovered, it becomes false
	isLifecycle := true

	// internal variable
	var isUserAccount bool

	n := len(l.Topics)
	if n > 0 {
		t0 = &l.Topics[0]
	}

	// for every indexed topic, check whether it is an end user account
	// if yes, then mark it as relevant for that account
	if n > 1 {
		t1 = &l.Topics[1]
		isUserAccount, addr1 = s.isEndUserAccount(*t1, stateDB)
		isLifecycle = isLifecycle && !isUserAccount
	}
	if n > 2 {
		t2 = &l.Topics[2]
		isUserAccount, addr2 = s.isEndUserAccount(*t2, stateDB)
		isLifecycle = isLifecycle && !isUserAccount
	}
	if n > 3 {
		t3 = &l.Topics[3]
		isUserAccount, addr3 = s.isEndUserAccount(*t3, stateDB)
		isLifecycle = isLifecycle && !isUserAccount
	}
	if n > 4 {
		t4 = &l.Topics[4]
		isUserAccount, addr4 = s.isEndUserAccount(*t4, stateDB)
		isLifecycle = isLifecycle && !isUserAccount
	}

	// normalise the data field to nil to avoid duplicates
	data := l.Data
	if len(data) == 0 {
		data = nil
	}

	dbBatch.ExecuteSQL("insert into events values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		t0, t1, t2, t3, t4,
		data, l.BlockHash, l.BlockNumber, l.TxHash, l.TxIndex, l.Index, l.Address,
		isLifecycle, addr1, addr2, addr3, addr4,
	)
}

// Of the log's topics, returns those that are (potentially) user addresses. A topic is considered a user address if:
//   - It has at least 12 leading zero bytes (since addresses are 20 bytes long, while hashes are 32) and at most 22 leading zero bytes
//   - It does not have associated code (meaning it's a smart-contract address)
//   - It has a non-zero nonce (to prevent accidental or malicious creation of the address matching a given topic,
//     forcing its events to become permanently private (this is not implemented for now)
func (s *storageImpl) isEndUserAccount(topic gethcommon.Hash, db *state.StateDB) (bool, *gethcommon.Address) {
	potentialAddr := common.ExtractPotentialAddress(topic)
	if potentialAddr == nil {
		return false, nil
	}
	addrBytes := potentialAddr.Bytes()
	// Check the database if there are already entries for this address
	var count int
	query := "select count(*) from events where relAddress1=? OR relAddress2=? OR relAddress3=? OR relAddress4=?"
	err := s.db.GetSQLDB().QueryRow(query, addrBytes, addrBytes, addrBytes, addrBytes).Scan(&count)
	if err != nil {
		// exit here
		s.logger.Crit("Could not execute query", log.ErrKey, err)
	}

	if count > 0 {
		return true, potentialAddr
	}

	// TODO A user address must have a non-zero nonce. This prevents accidental or malicious sending of funds to an
	// address matching a topic, forcing its events to become permanently private.
	// if db.GetNonce(potentialAddr) != 0

	// If the address has code, it's a smart contract address instead.
	if db.GetCode(*potentialAddr) == nil {
		return true, potentialAddr
	}

	return false, nil
}

func (s *storageImpl) SetHeadBatchPointer(l2Head *core.Batch, dbBatch *sql.Batch) error {
	if dbBatch == nil {
		panic("SetHeadBatchPointer called without an instance of sql.Batch")
	}

	// We update the canonical hash of the batch at this height.
	if err := obscurorawdb.SetL2HeadBatch(dbBatch, l2Head.Hash()); err != nil {
		return fmt.Errorf("could not write canonical hash. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) UpdateHeadRollup(l1Head *common.L1BlockHash, l2Head *common.L2BatchHash) error {
	dbBatch := s.db.NewBatch()
	if err := obscurorawdb.WriteL2HeadRollup(dbBatch, l1Head, l2Head); err != nil {
		return fmt.Errorf("could not write block state. Cause: %w", err)
	}
	if err := dbBatch.Write(); err != nil {
		return fmt.Errorf("could not save new head. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) UpdateL1Head(l1Head common.L1BlockHash) error {
	dbBatch := s.db.NewBatch()
	rawdb.WriteHeadHeaderHash(dbBatch, l1Head)
	if err := dbBatch.Write(); err != nil {
		return fmt.Errorf("could not save new L1 head. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) CreateStateDB(hash common.L2BatchHash) (*state.StateDB, error) {
	batch, err := s.FetchBatch(hash)
	if err != nil {
		return nil, err
	}

	statedb, err := state.New(batch.Header.Root, s.stateDB, nil)
	if err != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("could not create state DB. Cause: %w", err))
	}

	return statedb, nil
}

func (s *storageImpl) EmptyStateDB() (*state.StateDB, error) {
	statedb, err := state.New(gethcommon.BigToHash(big.NewInt(0)), s.stateDB, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create state DB. Cause: %w", err)
	}
	return statedb, nil
}

// GetReceiptsByHash retrieves the receipts for all transactions in a given batch.
func (s *storageImpl) GetReceiptsByHash(hash gethcommon.Hash) (types.Receipts, error) {
	number, err := obscurorawdb.ReadBatchNumber(s.db, hash)
	if err != nil {
		return nil, err
	}
	return obscurorawdb.ReadReceipts(s.db, hash, *number, s.chainConfig)
}

func (s *storageImpl) GetTransaction(txHash gethcommon.Hash) (*types.Transaction, gethcommon.Hash, uint64, uint64, error) {
	tx, blockHash, blockNumber, index, err := obscurorawdb.ReadTransaction(s.db, txHash)
	if err != nil {
		return nil, gethcommon.Hash{}, 0, 0, err
	}
	return tx, blockHash, blockNumber, index, nil
}

func (s *storageImpl) GetSender(txHash gethcommon.Hash) (gethcommon.Address, error) {
	tx, _, _, _, err := s.GetTransaction(txHash) //nolint:dogsled
	if err != nil {
		return gethcommon.Address{}, err
	}
	// todo - make the signer a field of the rollup chain
	msg, err := tx.AsMessage(types.NewLondonSigner(tx.ChainId()), nil)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("could not convert transaction to message to retrieve sender address in eth_getTransactionReceipt request. Cause: %w", err)
	}
	return msg.From(), nil
}

func (s *storageImpl) GetContractCreationTx(address gethcommon.Address) (*gethcommon.Hash, error) {
	return obscurorawdb.ReadContractTransaction(s.db, address)
}

func (s *storageImpl) GetTransactionReceipt(txHash gethcommon.Hash) (*types.Receipt, error) {
	_, blockHash, _, index, err := s.GetTransaction(txHash)
	if err != nil {
		return nil, err
	}

	receipts, err := s.GetReceiptsByHash(blockHash)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve receipts for transaction. Cause: %w", err)
	}

	if len(receipts) <= int(index) {
		return nil, fmt.Errorf("receipt index not matching the transactions in block: %s", blockHash.Hex())
	}
	receipt := receipts[index]

	return receipt, nil
}

func (s *storageImpl) FetchAttestedKey(address gethcommon.Address) (*ecdsa.PublicKey, error) {
	key, err := orm.FetchAttKey(s.db.GetSQLDB(), address)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve attestation key for address %s. Cause: %w", address, err)
	}

	publicKey, err := ethcrypto.DecompressPubkey(key)
	if err != nil {
		return nil, fmt.Errorf("could not parse key from db. Cause: %w", err)
	}

	return publicKey, nil
}

func (s *storageImpl) StoreAttestedKey(aggregator gethcommon.Address, key *ecdsa.PublicKey) error {
	_, err := orm.WriteAttKey(s.db.GetSQLDB(), aggregator, ethcrypto.CompressPubkey(key))
	return err
}

func (s *storageImpl) FetchBatchBySeqNo(seqNum uint64) (*core.Batch, error) {
	return orm.FindBatchBySeqNo(s.db.GetSQLDB(), seqNum)
}

func (s *storageImpl) StoreBatch(batch *core.Batch, receipts []*types.Receipt, dbBatch *sql.Batch) error {
	if dbBatch == nil {
		panic("StoreBatch called without an instance of sql.Batch")
	}

	if _, err := s.FetchBatchBySeqNo(batch.SeqNo().Uint64()); err == nil {
		return nil
		// return fmt.Errorf("batch with same sequence number already exists: %d", batch.SeqNo())
	}

	if err := orm.WriteBatch(dbBatch, batch); err != nil {
		return fmt.Errorf("could not write batch. Cause: %w", err)
	}

	if err := obscurorawdb.WriteBatch(dbBatch, batch); err != nil {
		return fmt.Errorf("could not write batch. Cause: %w", err)
	}
	if err := obscurorawdb.WriteTxLookupEntriesByBatch(dbBatch, batch); err != nil {
		return fmt.Errorf("could not write transaction lookup entries by batch. Cause: %w", err)
	}
	if err := obscurorawdb.WriteReceipts(dbBatch, batch.Hash(), receipts); err != nil {
		return fmt.Errorf("could not write transaction receipts. Cause: %w", err)
	}
	if err := obscurorawdb.WriteContractCreationTxs(dbBatch, receipts); err != nil {
		return fmt.Errorf("could not save contract creation transaction. Cause: %w", err)
	}
	// orm.UpdateConfigToBatch(dbBatch, current_sequence, batch.Header.SequencerOrderNo.Bytes())
	if err := obscurorawdb.WriteBatchBySequenceNum(dbBatch, batch); err != nil {
		return fmt.Errorf("could not save the current seqencer number. Cause: %w", err)
	}
	// todo fix this as batches always stored even if not canonical
	if err := obscurorawdb.IncrementContractCreationCount(s.db, dbBatch, receipts); err != nil {
		return fmt.Errorf("unable to increment contract count")
	}
	return nil
}

func (s *storageImpl) StoreL1Messages(blockHash common.L1BlockHash, messages common.CrossChainMessages) error {
	return obscurorawdb.StoreL1Messages(s.db, blockHash, messages, s.logger)
}

func (s *storageImpl) GetL1Messages(blockHash common.L1BlockHash) (common.CrossChainMessages, error) {
	return obscurorawdb.GetL1Messages(s.db, blockHash, s.logger)
}

func (s *storageImpl) StoreEnclaveKey(enclaveKey *ecdsa.PrivateKey) error {
	return obscurorawdb.StoreEnclaveKey(s.db, enclaveKey, s.logger)
}

func (s *storageImpl) GetEnclaveKey() (*ecdsa.PrivateKey, error) {
	return obscurorawdb.GetEnclaveKey(s.db, s.logger)
}

func (s *storageImpl) StoreRollup(rollup *common.ExtRollup) error {
	dbBatch := s.db.NewBatch()

	if err := obscurorawdb.WriteRollup(dbBatch, rollup); err != nil {
		return fmt.Errorf("could not write rollup. Cause: %w", err)
	}

	if err := dbBatch.Write(); err != nil {
		return fmt.Errorf("could not write rollup to storage. Cause: %w", err)
	}
	return nil
}

// utility function that knows how to load relevant logs from the database
// todo always pass in the actual batch hashes because of reorgs, or make sure to clean up log entries from discarded batches
func (s *storageImpl) loadLogs(requestingAccount *gethcommon.Address, whereCondition string, whereParams []any) ([]*types.Log, error) {
	if requestingAccount == nil {
		return nil, fmt.Errorf("logs can only be requested for an account")
	}

	result := make([]*types.Log, 0)
	// todo - remove the "distinct" once the fast-finality work is completed
	// currently the events seem to be stored twice because of some weird logic in the rollup/batch processing.
	// Note: the where 1=1 clauses allows for an easier query building
	query := "select distinct topic0, topic1, topic2, topic3, topic4, datablob, blockHash, blockNumber, txHash, txIdx, logIdx, address from events where 1=1 "
	var queryParams []any

	// Add relevancy rules
	//  An event is considered relevant to all account owners whose addresses are used as topics in the event.
	//	In case there are no account addresses in an event's topics, then the event is considered relevant to everyone (known as a "lifecycle event").
	query += " AND (lifecycleEvent OR (relAddress1=? OR relAddress2=? OR relAddress3=? OR relAddress4=?)) "
	queryParams = append(queryParams, requestingAccount.Bytes())
	queryParams = append(queryParams, requestingAccount.Bytes())
	queryParams = append(queryParams, requestingAccount.Bytes())
	queryParams = append(queryParams, requestingAccount.Bytes())

	query += whereCondition
	queryParams = append(queryParams, whereParams...)

	rows, err := s.db.GetSQLDB().Query(query, queryParams...) //nolint: rowserrcheck
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		l := types.Log{
			Topics: []gethcommon.Hash{},
		}
		var t0, t1, t2, t3, t4 sql2.NullString
		err = rows.Scan(&t0, &t1, &t2, &t3, &t4, &l.Data, &l.BlockHash, &l.BlockNumber, &l.TxHash, &l.TxIndex, &l.Index, &l.Address)
		if err != nil {
			return nil, fmt.Errorf("could not load log entry from db: %w", err)
		}

		for _, topic := range []sql2.NullString{t0, t1, t2, t3, t4} {
			if topic.Valid {
				l.Topics = append(l.Topics, stringToHash(topic))
			}
		}

		result = append(result, &l)
	}

	if err = rows.Close(); err != nil { //nolint: sqlclosecheck
		return nil, syserr.NewInternalError(err)
	}

	return result, nil
}

func (s *storageImpl) DebugGetLogs(txHash common.TxHash) ([]*tracers.DebugLogs, error) {
	var queryParams []any

	query := `select distinct 
    			relAddress1, relAddress2, relAddress3, relAddress4,
    			lifecycleEvent,
    			topic0, topic1, topic2, topic3, topic4, 
    			datablob, blockHash, blockNumber, txHash, txIdx, logIdx, address 
				from events 
				where txHash = ?`

	queryParams = append(queryParams, txHash.Bytes())

	result := make([]*tracers.DebugLogs, 0)

	rows, err := s.db.GetSQLDB().Query(query, queryParams...) //nolint: rowserrcheck
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		l := tracers.DebugLogs{
			Log: types.Log{
				Topics: []gethcommon.Hash{},
			},
			LifecycleEvent: false,
		}

		var t0, t1, t2, t3, t4 sql2.NullString
		var relAddress1, relAddress2, relAddress3, relAddress4 sql2.NullByte
		err = rows.Scan(
			&relAddress1,
			&relAddress2,
			&relAddress3,
			&relAddress4,
			&l.LifecycleEvent,
			&t0, &t1, &t2, &t3, &t4,
			&l.Data,
			&l.BlockHash,
			&l.BlockNumber,
			&l.TxHash,
			&l.TxIndex,
			&l.Index,
			&l.Address,
		)
		if err != nil {
			return nil, fmt.Errorf("could not load log entry from db: %w", err)
		}

		for _, topic := range []sql2.NullString{t0, t1, t2, t3, t4} {
			if topic.Valid {
				l.Topics = append(l.Topics, stringToHash(topic))
			}
		}

		l.RelAddress1 = bytesToHash(relAddress1)
		l.RelAddress2 = bytesToHash(relAddress2)
		l.RelAddress3 = bytesToHash(relAddress3)
		l.RelAddress4 = bytesToHash(relAddress4)

		result = append(result, &l)
	}

	if err = rows.Close(); err != nil { //nolint: sqlclosecheck
		return nil, err
	}

	return result, nil
}

func (s *storageImpl) FilterLogs(
	requestingAccount *gethcommon.Address,
	fromBlock, toBlock *big.Int,
	blockHash *common.L2BatchHash,
	addresses []gethcommon.Address,
	topics [][]gethcommon.Hash,
) ([]*types.Log, error) {
	queryParams := []any{}
	query := ""
	if blockHash != nil {
		query += " AND blockHash = ?"
		queryParams = append(queryParams, blockHash.Bytes())
	}

	// ignore negative numbers
	if fromBlock != nil && fromBlock.Sign() > 0 {
		query += " AND blockNumber >= ?"
		queryParams = append(queryParams, fromBlock.Int64())
	}
	if toBlock != nil && toBlock.Sign() > 0 {
		query += " AND blockNumber <= ?"
		queryParams = append(queryParams, toBlock.Int64())
	}

	if len(addresses) > 0 {
		query += " AND address in (?" + strings.Repeat(",?", len(addresses)-1) + ")"
		for _, address := range addresses {
			queryParams = append(queryParams, address.Bytes())
		}
	}
	if len(topics) > 5 {
		return nil, fmt.Errorf("invalid filter. Too many topics")
	}
	if len(topics) > 0 {
		for i, sub := range topics {
			// empty rule set == wildcard
			if len(sub) > 0 {
				column := fmt.Sprintf("topic%d", i)
				query += " AND " + column + " in (?" + strings.Repeat(",?", len(sub)-1) + ")"
				for _, topic := range sub {
					queryParams = append(queryParams, topic.Bytes())
				}
			}
		}
	}

	return s.loadLogs(requestingAccount, query, queryParams)
}

func (s *storageImpl) GetContractCount() (*big.Int, error) {
	return obscurorawdb.ReadContractCreationCount(s.db)
}

func stringToHash(ns sql2.NullString) gethcommon.Hash {
	value, err := ns.Value()
	if err != nil {
		return [32]byte{}
	}
	s := value.(string)
	result := gethcommon.Hash{}
	result.SetBytes([]byte(s))
	return result
}

func bytesToHash(b sql2.NullByte) *gethcommon.Hash {
	result := gethcommon.Hash{}

	if !b.Valid {
		return nil
	}

	value, err := b.Value()
	if err != nil {
		return nil
	}
	s := value.(string)

	result.SetBytes([]byte(s))
	return &result
}
