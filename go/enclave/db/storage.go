package db

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/go/enclave/crypto"

	obscurorawdb "github.com/obscuronet/go-obscuro/go/enclave/db/rawdb"

	"github.com/ethereum/go-ethereum/params"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
)

// TODO - Consistency around whether we assert the secret is available or not.

type storageImpl struct {
	db          ethdb.Database
	stateDB     state.Database
	chainConfig *params.ChainConfig
	logger      gethlog.Logger
}

func NewStorage(backingDB ethdb.Database, chainConfig *params.ChainConfig, logger gethlog.Logger) Storage {
	return &storageImpl{
		db:          backingDB,
		stateDB:     state.NewDatabase(backingDB),
		chainConfig: chainConfig,
		logger:      logger,
	}
}

func (s *storageImpl) StoreGenesisRollup(rol *core.Rollup) error {
	batch := s.db.NewBatch()

	if err := obscurorawdb.WriteGenesisHash(s.db, rol.Hash()); err != nil {
		return fmt.Errorf("could not write genesis hash. Cause: %w", err)
	}
	if err := obscurorawdb.WriteRollup(batch, rol); err != nil {
		return fmt.Errorf("could not write rollup to storage. Cause: %w", err)
	}

	if err := batch.Write(); err != nil {
		return fmt.Errorf("could not write genesis rollup to storage. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) FetchGenesisRollup() (*core.Rollup, error) {
	hash, err := obscurorawdb.ReadGenesisHash(s.db)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve genesis rollup. Cause: %w", err)
	}
	rollup, err := s.FetchRollup(*hash)
	if err != nil {
		return nil, err
	}
	return rollup, nil
}

func (s *storageImpl) FetchHeadRollup() (*core.Rollup, error) {
	l2Head, err := s.FetchL2Head()
	if err != nil {
		return nil, fmt.Errorf("could not fetch L2 head hash")
	}
	return s.FetchRollup(*l2Head)
}

func (s *storageImpl) FetchRollup(hash common.L2RootHash) (*core.Rollup, error) {
	s.assertSecretAvailable()
	rollup, err := obscurorawdb.ReadRollup(s.db, hash)
	if err != nil {
		return nil, err
	}
	return rollup, nil
}

func (s *storageImpl) FetchRollupByHeight(height uint64) (*core.Rollup, error) {
	if height == 0 {
		genesisRollup, err := s.FetchGenesisRollup()
		if err != nil {
			return nil, fmt.Errorf("could not fetch genesis rollup. Cause: %w", err)
		}
		return genesisRollup, nil
	}

	hash, err := obscurorawdb.ReadCanonicalHash(s.db, height)
	if err != nil {
		return nil, err
	}
	return s.FetchRollup(*hash)
}

func (s *storageImpl) StoreBlock(b *types.Block) {
	s.assertSecretAvailable()
	rawdb.WriteBlock(s.db, b)
}

func (s *storageImpl) FetchBlock(blockHash common.L1RootHash) (*types.Block, error) {
	s.assertSecretAvailable()
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
	s.assertSecretAvailable()
	block, err := s.FetchBlock(rawdb.ReadHeadHeaderHash(s.db))
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (s *storageImpl) StoreSecret(secret crypto.SharedEnclaveSecret) error {
	return obscurorawdb.WriteSharedSecret(s.db, secret)
}

func (s *storageImpl) FetchSecret() (*crypto.SharedEnclaveSecret, error) {
	return obscurorawdb.ReadSharedSecret(s.db)
}

func (s *storageImpl) ParentRollup(r *core.Rollup) (*core.Rollup, error) {
	s.assertSecretAvailable()
	return s.FetchRollup(r.Header.ParentHash)
}

func (s *storageImpl) ParentBlock(b *types.Block) (*types.Block, error) {
	s.assertSecretAvailable()
	return s.FetchBlock(b.Header().ParentHash)
}

func (s *storageImpl) IsAncestor(block *types.Block, maybeAncestor *types.Block) bool {
	s.assertSecretAvailable()
	if bytes.Equal(maybeAncestor.Hash().Bytes(), block.Hash().Bytes()) {
		return true
	}

	if maybeAncestor.NumberU64() >= block.NumberU64() {
		return false
	}

	p, err := s.ParentBlock(block)
	if err != nil {
		return false
	}

	return s.IsAncestor(p, maybeAncestor)
}

func (s *storageImpl) IsBlockAncestor(block *types.Block, maybeAncestor common.L1RootHash) bool {
	resolvedBlock, err := s.FetchBlock(maybeAncestor)
	if err != nil {
		return false
	}
	return s.IsAncestor(block, resolvedBlock)
}

func (s *storageImpl) HealthCheck() (bool, error) {
	headRollup, err := s.FetchHeadRollup()
	if err != nil {
		s.logger.Error("unable to HealthCheck storage", "err", err)
		return false, err
	}
	return headRollup != nil, nil
}

func (s *storageImpl) assertSecretAvailable() {
	// TODO uncomment this
	//if s.FetchSecret() == nil {
	//	panic("Enclave not initialized")
	//}
}

func (s *storageImpl) Proof(r *core.Rollup) (*types.Block, error) {
	block, err := s.FetchBlock(r.Header.L1Proof)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (s *storageImpl) FetchHeadRollupForL1Block(blockHash common.L1RootHash) (*common.L2RootHash, error) {
	return obscurorawdb.ReadL2Head(s.db, blockHash)
}

func (s *storageImpl) FetchLogs(blockHash common.L1RootHash) ([]*types.Log, error) {
	logs, err := obscurorawdb.ReadBlockLogs(s.db, blockHash)
	if err != nil {
		// TODO - Return the error itself, once we move from `errutil.ErrNotFound` to `ethereum.NotFound`
		return nil, errutil.ErrNotFound
	}
	return logs, nil
}

// TODO - #718 - This method has behaviour that's too dependent on various flags. Decompose.
func (s *storageImpl) UpdateL2HeadForL1Block(l1Head common.L1RootHash, l2Head *core.Rollup, receipts []*types.Receipt, isNewRollup bool) error {
	batch := s.db.NewBatch()

	if err := obscurorawdb.WriteL2Head(batch, l1Head, l2Head.Hash()); err != nil {
		return fmt.Errorf("could not write block state. Cause: %w", err)
	}

	if isNewRollup {
		err := s.storeNewRollup(batch, l2Head, receipts)
		if err != nil {
			return err
		}
	}

	if err := obscurorawdb.WriteCanonicalHash(batch, l2Head); err != nil {
		return fmt.Errorf("could not write canonical hash. Cause: %w", err)
	}
	var logs []*types.Log
	for _, receipt := range receipts {
		logs = append(logs, receipt.Logs...)
	}
	if err := obscurorawdb.WriteBlockLogs(batch, l1Head, logs); err != nil {
		return fmt.Errorf("could not write block logs. Cause: %w", err)
	}

	if err := batch.Write(); err != nil {
		return fmt.Errorf("could not save new head. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) UpdateL1Head(l1Head common.L1RootHash) error {
	batch := s.db.NewBatch()
	rawdb.WriteHeadHeaderHash(batch, l1Head)
	if err := batch.Write(); err != nil {
		return fmt.Errorf("could not save new L1 head. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) CreateStateDB(hash common.L2RootHash) (*state.StateDB, error) {
	rollup, err := s.FetchRollup(hash)
	if err != nil {
		return nil, err
	}

	// todo - snapshots?
	statedb, err := state.New(rollup.Header.Root, s.stateDB, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create state DB. Cause: %w", err)
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

func (s *storageImpl) FetchL2Head() (*common.L2RootHash, error) {
	l1Head := rawdb.ReadHeadHeaderHash(s.db)
	if (bytes.Equal(l1Head.Bytes(), gethcommon.Hash{}.Bytes())) {
		return nil, errutil.ErrNotFound
	}

	l2Head, err := obscurorawdb.ReadL2Head(s.db, l1Head)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve block state for head. Cause: %w", err)
	}

	return l2Head, nil
}

// GetReceiptsByHash retrieves the receipts for all transactions in a given rollup.
func (s *storageImpl) GetReceiptsByHash(hash gethcommon.Hash) (types.Receipts, error) {
	number, err := obscurorawdb.ReadHeaderNumber(s.db, hash)
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

func (s *storageImpl) FetchAttestedKey(aggregator gethcommon.Address) (*ecdsa.PublicKey, error) {
	return obscurorawdb.ReadAttestationKey(s.db, aggregator)
}

func (s *storageImpl) StoreAttestedKey(aggregator gethcommon.Address, key *ecdsa.PublicKey) error {
	return obscurorawdb.WriteAttestationKey(s.db, aggregator, key)
}

func (s *storageImpl) storeNewRollup(batch ethdb.Batch, rollup *core.Rollup, receipts []*types.Receipt) error {
	if err := obscurorawdb.WriteRollup(batch, rollup); err != nil {
		return fmt.Errorf("could not write rollup. Cause: %w", err)
	}
	if err := obscurorawdb.WriteTxLookupEntriesByRollup(batch, rollup); err != nil {
		return fmt.Errorf("could not write transaction lookup entries by block. Cause: %w", err)
	}
	if err := obscurorawdb.WriteReceipts(batch, rollup.Hash(), receipts); err != nil {
		return fmt.Errorf("could not write transaction receipts. Cause: %w", err)
	}
	if err := obscurorawdb.WriteContractCreationTxs(batch, receipts); err != nil {
		return fmt.Errorf("could not save contract creation transaction. Cause: %w", err)
	}

	return nil
}

func (s *storageImpl) StoreL1Messages(blockHash gethcommon.Hash, messages common.CrossChainMessages) error {
	return obscurorawdb.StoreL1Messages(s.db, blockHash, messages, s.logger)
}

func (s *storageImpl) GetL1Messages(blockHash gethcommon.Hash) (common.CrossChainMessages, error) {
	return obscurorawdb.GetL1Messages(s.db, blockHash, s.logger)
}
