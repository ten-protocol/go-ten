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
	err := obscurorawdb.WriteGenesisHash(s.db, rol.Hash())
	if err != nil {
		return fmt.Errorf("could not write genesis hash. Cause: %w", err)
	}
	return s.StoreRollup(rol)
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
	hash, err := obscurorawdb.ReadHeadRollupHash(s.db)
	if err != nil {
		return nil, err
	}
	r, err := s.FetchRollup(*hash)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s *storageImpl) StoreRollup(rollup *core.Rollup) error {
	s.assertSecretAvailable()

	batch := s.db.NewBatch()
	if err := obscurorawdb.WriteRollup(batch, rollup); err != nil {
		return fmt.Errorf("could not write rollup to storage. Cause: %w", err)
	}
	if err := batch.Write(); err != nil {
		return fmt.Errorf("could not write rollup to storage. Cause: %w", err)
	}
	return nil
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

func (s *storageImpl) FetchRollups(height uint64) ([]*core.Rollup, error) {
	s.assertSecretAvailable()
	return obscurorawdb.ReadRollupsForHeight(s.db, height)
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
	s.assertSecretAvailable()
	if bytes.Equal(maybeAncestor.Bytes(), block.Hash().Bytes()) {
		return true
	}

	if bytes.Equal(maybeAncestor.Bytes(), common.GenesisBlock.Hash().Bytes()) {
		return true
	}

	if block.NumberU64() == common.L1GenesisHeight {
		return false
	}

	resolvedBlock, err := s.FetchBlock(maybeAncestor)
	if err == nil {
		if resolvedBlock.NumberU64() >= block.NumberU64() {
			return false
		}
	}

	p, err := s.ParentBlock(block)
	if err != nil {
		// TODO - If error is not `errutil.ErrNotFound`, throw.
		return false
	}

	return s.IsBlockAncestor(p, maybeAncestor)
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

// todo - find a better way. This is a workaround to handle rollups created with proofs that haven't propagated yet
func (s *storageImpl) ProofHeight(r *core.Rollup) int64 {
	v, err := s.FetchBlock(r.Header.L1Proof)
	if err != nil {
		return -1
	}
	return int64(v.NumberU64())
}

func (s *storageImpl) Proof(r *core.Rollup) (*types.Block, error) {
	block, err := s.FetchBlock(r.Header.L1Proof)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (s *storageImpl) FetchBlockState(hash common.L1RootHash) (*core.BlockState, error) {
	return obscurorawdb.ReadBlockState(s.db, hash)
}

func (s *storageImpl) FetchLogs(hash common.L1RootHash) ([]*types.Log, error) {
	logs, err := obscurorawdb.ReadBlockLogs(s.db, hash)
	if err != nil {
		// TODO - Return the error itself, once we move from `errutil.ErrNotFound` to `ethereum.NotFound`
		return nil, errutil.ErrNotFound
	}
	return logs, nil
}

func (s *storageImpl) StoreNewHead(state *core.BlockState, rollup *core.Rollup, receipts []*types.Receipt, logs []*types.Log) error {
	batch := s.db.NewBatch()

	if state.FoundNewRollup {
		err := s.storeNewRollup(batch, rollup, receipts)
		if err != nil {
			return err
		}
	}

	if err := obscurorawdb.WriteBlockState(batch, state); err != nil {
		return fmt.Errorf("could not write block state. Cause: %w", err)
	}
	if err := obscurorawdb.WriteBlockLogs(batch, state.Block, logs); err != nil {
		return fmt.Errorf("could not write block logs. Cause: %w", err)
	}
	rawdb.WriteHeadHeaderHash(batch, state.Block)

	if err := batch.Write(); err != nil {
		return fmt.Errorf("could not save new head. Cause: %w", err)
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

func (s *storageImpl) FetchHeadState() (*core.BlockState, error) {
	h := rawdb.ReadHeadHeaderHash(s.db)
	if (bytes.Equal(h.Bytes(), gethcommon.Hash{}.Bytes())) {
		return nil, errutil.ErrNotFound
	}

	blockState, err := obscurorawdb.ReadBlockState(s.db, h)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve block state for head. Cause: %w", err)
	}

	return blockState, nil
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
	if err := obscurorawdb.WriteHeadHeaderHash(batch, rollup.Hash()); err != nil {
		return fmt.Errorf("could not write head header hash. Cause: %w", err)
	}
	if err := obscurorawdb.WriteCanonicalHash(batch, rollup.Hash(), rollup.NumberU64()); err != nil {
		return fmt.Errorf("could not write canonical hash. Cause: %w", err)
	}
	if err := obscurorawdb.WriteTxLookupEntriesByBlock(batch, rollup); err != nil {
		return fmt.Errorf("could not write transaction lookup entries by block. Cause: %w", err)
	}
	if err := obscurorawdb.WriteHeadRollupHash(batch, rollup.Hash()); err != nil {
		return fmt.Errorf("could not write head rollup hash. Cause: %w", err)
	}
	if err := obscurorawdb.WriteReceipts(batch, rollup.Hash(), rollup.NumberU64(), receipts); err != nil {
		return fmt.Errorf("could not write transaction receipts. Cause: %w", err)
	}
	if err := obscurorawdb.WriteContractCreationTx(batch, receipts); err != nil {
		return fmt.Errorf("could not save contract creation transaction. Cause: %w", err)
	}
	return nil
}
