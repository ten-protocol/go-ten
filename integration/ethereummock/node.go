package ethereummock

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ten-protocol/go-ten/go/host/l1"

	"github.com/ten-protocol/go-ten/go/common/async"

	"github.com/google/uuid"

	"github.com/ten-protocol/go-ten/go/common/errutil"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ten-protocol/go-ten/go/common"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	ethclient_ethereum "github.com/ethereum/go-ethereum/ethclient"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/ethadapter/erc20contractlib"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"
)

const SecondsPerSlot = uint64(20)

type L1Network interface {
	// BroadcastBlock - send the block and the parent to make sure there are no gaps
	BroadcastBlock(b common.EncodedL1Block, p common.EncodedL1Block)
	BroadcastTx(tx *types.Transaction)
}

type MiningConfig struct {
	PowTime      common.Latency
	LogFile      string
	L1BeaconPort int
}

type TxDB interface {
	Txs(block *types.Block) (map[common.TxHash]*types.Transaction, bool)
	AddTxs(*types.Block, map[common.TxHash]*types.Transaction)
}

type StatsCollector interface {
	// L1Reorg registers when a miner has to process a reorg (a winning block from a fork)
	L1Reorg(id gethcommon.Address)
}

type BlockWithBlobs struct {
	Block *types.Block
	Blobs []*kzg4844.Blob
}

type Node struct {
	l2ID          gethcommon.Address // the address of the Obscuro node this client is dedicated to
	cfg           MiningConfig
	Network       L1Network
	mining        bool
	stats         StatsCollector
	BlockResolver *blockResolverInMem
	BlobResolver  l1.BlobResolver
	db            TxDB
	subs          map[uuid.UUID]*mockSubscription // active subscription for mock blocks
	subMu         sync.Mutex

	// Channels
	exitCh       chan bool // the Node stops
	exitMiningCh chan bool // the mining loop is notified to stop
	interrupt    *int32

	p2pCh       chan *types.Block       // this is where blocks received from peers are dropped
	miningCh    chan *types.Block       // this is where blocks created by the mining setup of the current node are dropped
	canonicalCh chan *types.Block       // this is where the main processing routine drops blocks that are canonical
	mempoolCh   chan *types.Transaction // where l1 transactions to be published in the next block are added

	// internal
	headInCh         chan bool
	headOutCh        chan *types.Block
	erc20ContractLib erc20contractlib.ERC20ContractLib
	mgmtContractLib  mgmtcontractlib.MgmtContractLib

	logger gethlog.Logger
}

func (m *Node) PrepareTransactionToSend(_ context.Context, txData types.TxData, _ gethcommon.Address) (types.TxData, error) {
	switch tx := txData.(type) {
	case *types.LegacyTx:
		return createLegacyTx(txData)
	case *types.BlobTx:
		return createBlobTx(txData)
	default:
		return nil, fmt.Errorf("unsupported transaction type: %T", tx)
	}
}

func (m *Node) PrepareTransactionToRetry(ctx context.Context, txData types.TxData, from gethcommon.Address, _ uint64, _ int) (types.TxData, error) {
	return m.PrepareTransactionToSend(ctx, txData, from)
}

func createLegacyTx(txData types.TxData) (types.TxData, error) {
	tx := types.NewTx(txData)
	return &types.LegacyTx{
		Nonce:    123,
		GasPrice: tx.GasPrice(),
		Gas:      tx.Gas(),
		To:       tx.To(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}, nil
}

func createBlobTx(txData types.TxData) (types.TxData, error) {
	tx := types.NewTx(txData)
	return &types.BlobTx{
		To:         *tx.To(),
		Data:       tx.Data(),
		BlobHashes: tx.BlobHashes(),
		Sidecar:    tx.BlobTxSidecar(),
	}, nil
}

func (m *Node) SendTransaction(tx *types.Transaction) error {
	m.Network.BroadcastTx(tx)
	return nil
}

func (m *Node) TransactionReceipt(_ gethcommon.Hash) (*types.Receipt, error) {
	// all transactions are immediately processed
	return &types.Receipt{
		BlockNumber: big.NewInt(1),
		Status:      types.ReceiptStatusSuccessful,
	}, nil
}

func (m *Node) Nonce(gethcommon.Address) (uint64, error) {
	return 0, nil
}

func (m *Node) getRollupFromBlock(block *types.Block) *common.ExtRollup {
	for _, tx := range block.Transactions() {
		decodedTx := m.mgmtContractLib.DecodeTx(tx)
		if decodedTx == nil {
			continue
		}
		switch l1tx := decodedTx.(type) {
		case *ethadapter.L1RollupHashes:
			ctx := context.TODO()
			blobs, _ := m.BlobResolver.FetchBlobs(ctx, block.Header(), l1tx.BlobHashes)
			r, err := ethadapter.ReconstructRollup(blobs)
			if err != nil {
				m.logger.Error("could not recreate rollup from blobs. Cause: %w", err)
				return nil
			}
			return r
		}
	}
	return nil
}

func (m *Node) FetchLastBatchSeqNo(gethcommon.Address) (*big.Int, error) {
	startingBlock, err := m.FetchHeadBlock()
	if err != nil {
		return nil, err
	}

	for currentBlock := startingBlock; currentBlock.NumberU64() != 0; {
		currentBlock, err = m.BlockByHash(currentBlock.Header().ParentHash)
		if err != nil {
			m.logger.Error("Error fetching block by hash", "error", err)
			break
		}
		rollup := m.getRollupFromBlock(currentBlock)
		if rollup != nil {
			return big.NewInt(int64(rollup.Header.LastBatchSeqNo)), nil
		}
	}
	// the first batch is number 1
	return big.NewInt(int64(common.L2GenesisSeqNo)), nil
}

// BlockListener provides stream of latest mock head headers as they are created
func (m *Node) BlockListener() (chan *types.Header, ethereum.Subscription) {
	id := uuid.New()
	mockSub := &mockSubscription{
		node:   m,
		id:     id,
		headCh: make(chan *types.Header),
	}
	m.subMu.Lock()
	defer m.subMu.Unlock()
	m.subs[id] = mockSub
	return mockSub.headCh, mockSub
}

func (m *Node) BlockNumber() (uint64, error) {
	blk, err := m.BlockResolver.FetchHeadBlock(context.Background())
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			return 0, ethereum.NotFound
		}
		return 0, fmt.Errorf("could not retrieve head block. Cause: %w", err)
	}
	return blk.NumberU64(), nil
}

func (m *Node) BlockByNumber(n *big.Int) (*types.Block, error) {
	if n.Int64() == 0 {
		return MockGenesisBlock, nil
	}
	// TODO this should be a method in the resolver
	blk, err := m.BlockResolver.FetchHeadBlock(context.Background())
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			return nil, ethereum.NotFound
		}
		return nil, fmt.Errorf("could not retrieve head block. Cause: %w", err)
	}
	for !bytes.Equal(blk.ParentHash().Bytes(), (common.L1BlockHash{}).Bytes()) {
		if blk.NumberU64() == n.Uint64() {
			return blk, nil
		}

		blk, err = m.BlockResolver.FetchBlock(context.Background(), blk.ParentHash())
		if err != nil {
			return nil, fmt.Errorf("could not retrieve parent for block in chain. Cause: %w", err)
		}
	}
	return nil, ethereum.NotFound
}

func (m *Node) BlockByHash(id gethcommon.Hash) (*types.Block, error) {
	blk, err := m.BlockResolver.FetchBlock(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("block could not be retrieved. Cause: %w", err)
	}
	return blk, nil
}

func (m *Node) FetchHeadBlock() (*types.Block, error) {
	block, err := m.BlockResolver.FetchHeadBlock(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not retrieve head block. Cause: %w", err)
	}
	return block, nil
}

func (m *Node) Info() ethadapter.Info {
	return ethadapter.Info{
		L2ID: m.l2ID,
	}
}

func (m *Node) IsBlockAncestor(block *types.Block, proof common.L1BlockHash) bool {
	return m.BlockResolver.IsBlockAncestor(context.Background(), block, proof)
}

func (m *Node) BalanceAt(gethcommon.Address, *big.Int) (*big.Int, error) {
	panic("not implemented")
}

// GetLogs is a mock method - we don't really have logs on the mock transactions, so it returns a basic log for every tx
// so the host recognises them as relevant
func (m *Node) GetLogs(fq ethereum.FilterQuery) ([]types.Log, error) {
	logs := make([]types.Log, 0)
	if fq.BlockHash == nil {
		return logs, nil
	}
	blk, err := m.BlockByHash(*fq.BlockHash)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve block. Cause: %w", err)
	}
	for _, tx := range blk.Transactions() {
		dummyLog := types.Log{
			BlockHash: blk.Hash(),
			TxHash:    tx.Hash(),
		}
		logs = append(logs, dummyLog)
	}
	return logs, nil
}

func (m *Node) Start() {
	if m.mining {
		// This starts the mining
		go m.startMining()
	}

	err := m.BlockResolver.StoreBlock(context.Background(), MockGenesisBlock, nil)
	if err != nil {
		m.logger.Crit("Failed to store block")
	}
	head := m.setHead(MockGenesisBlock)

	for {
		select {
		case p2pb := <-m.p2pCh: // Received from peers
			_, err := m.BlockResolver.FetchBlock(context.Background(), p2pb.Hash())
			// only process blocks if they haven't been processed before
			if err != nil {
				if errors.Is(err, errutil.ErrNotFound) {
					head = m.processBlock(p2pb, head)
				} else {
					panic(fmt.Errorf("could not retrieve parent block. Cause: %w", err))
				}
			}

		case mb := <-m.miningCh: // Received from the local mining
			head = m.processBlock(mb, head)
			if bytes.Equal(head.Hash().Bytes(), mb.Hash().Bytes()) { // Only broadcast if it's the new head
				p, err := m.BlockResolver.FetchBlock(context.Background(), mb.ParentHash())
				if err != nil {
					panic(fmt.Errorf("could not retrieve parent. Cause: %w", err))
				}
				encodedBlock, err := common.EncodeBlock(mb)
				if err != nil {
					panic(fmt.Errorf("could not encode block. Cause: %w", err))
				}
				encodedParentBlock, err := common.EncodeBlock(p)
				if err != nil {
					panic(fmt.Errorf("could not encode parent block. Cause: %w", err))
				}
				m.Network.BroadcastBlock(encodedBlock, encodedParentBlock)
			}
		case <-m.headInCh:
			m.headOutCh <- head
		case <-m.exitCh:
			return
		}
	}
}

func (m *Node) processBlock(b *types.Block, head *types.Block) *types.Block {
	err := m.BlockResolver.StoreBlock(context.Background(), b, nil)
	if err != nil {
		m.logger.Crit("Failed to store block. Cause: %w", err)
	}

	_, err = m.BlockResolver.FetchBlock(context.Background(), b.Header().ParentHash)
	// only proceed if the parent is available
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			m.logger.Info(fmt.Sprintf("Parent block not found=b_%d", common.ShortHash(b.Header().ParentHash)))
			return head
		}
		m.logger.Crit("Could not fetch block parent. Cause: %w", err)
	}

	// Ignore superseded blocks
	if b.NumberU64() <= head.NumberU64() {
		return head
	}

	// Check for Reorgs
	if !m.BlockResolver.IsAncestor(context.Background(), b, head) {
		m.stats.L1Reorg(m.l2ID)
		fork, err := LCA(context.Background(), head, b, m.BlockResolver)
		if err != nil {
			panic(err)
		}
		m.logger.Info(
			fmt.Sprintf("L1Reorg new=b_%d(%d), old=b_%d(%d), fork=b_%d(%d)", common.ShortHash(b.Hash()), b.NumberU64(), common.ShortHash(head.Hash()), head.NumberU64(), common.ShortHash(fork.CommonAncestor.Hash()), fork.CommonAncestor.Number.Uint64()))
		return m.setFork(m.BlocksBetween(fork.CommonAncestor, b))
	}
	if b.NumberU64() > (head.NumberU64() + 1) {
		m.logger.Crit("Should not happen")
	}

	return m.setHead(b)
}

// Notifies the Miner to start mining on the new block and the aggregator to produce rollups
func (m *Node) setHead(b *types.Block) *types.Block {
	if atomic.LoadInt32(m.interrupt) == 1 {
		return b
	}

	// notify the client subscriptions
	m.subMu.Lock()
	for _, s := range m.subs {
		sub := s
		go sub.publish(b)
	}
	m.subMu.Unlock()
	m.canonicalCh <- b

	return b
}

func (m *Node) setFork(blocks []*types.Block) *types.Block {
	head := blocks[len(blocks)-1]
	if atomic.LoadInt32(m.interrupt) == 1 {
		return head
	}

	// notify the client subs
	m.subMu.Lock()
	for _, s := range m.subs {
		sub := s
		go sub.publishAll(blocks)
	}
	m.subMu.Unlock()

	m.canonicalCh <- head

	return head
}

// P2PReceiveBlock is called by counterparties when there is a block to broadcast
// All it does is drop the blocks in a channel for processing.
func (m *Node) P2PReceiveBlock(b common.EncodedL1Block, p common.EncodedL1Block) {
	if atomic.LoadInt32(m.interrupt) == 1 {
		return
	}
	decodedBlock, err := b.DecodeBlock()
	if err != nil {
		panic(fmt.Errorf("could not decode block. Cause: %w", err))
	}
	decodedParentBlock, err := p.DecodeBlock()
	if err != nil {
		panic(fmt.Errorf("could not decode parent block. Cause: %w", err))
	}
	m.p2pCh <- decodedParentBlock
	m.p2pCh <- decodedBlock
}

// startMining - listens on the canonicalCh and schedule a go routine that produces a block after a PowTime and drop it
// on the miningCh channel
func (m *Node) startMining() {
	m.logger.Info(" starting miner...")
	// stores all transactions seen from the beginning of time.
	mempool := make([]*types.Transaction, 0)
	z := int32(0)
	interrupt := &z
	for {
		select {
		case <-m.exitMiningCh:
			return
		case tx := <-m.mempoolCh:
			mempool = append(mempool, tx)

		case canonicalBlock := <-m.canonicalCh:
			// A new canonical block was found. Start a new round based on that block.

			// remove transactions that are already considered committed
			mempool = m.removeCommittedTransactions(context.Background(), canonicalBlock, mempool, m.BlockResolver, m.db)

			// notify the existing mining go routine to stop mining
			atomic.StoreInt32(interrupt, 1)
			c := int32(0)
			interrupt = &c

			// Generate a random number, and wait for that number of ms. Equivalent to PoW
			// Include all rollups received during this period.
			blockTime := uint64(time.Now().Unix())
			async.Schedule(m.cfg.PowTime(), func() {
				toInclude := findNotIncludedTxs(canonicalBlock, mempool, m.BlockResolver, m.db)
				// todo - iterate through the rollup transactions and include only the ones with the proof on the canonical chain
				if atomic.LoadInt32(m.interrupt) == 1 {
					return
				}

				block, blobs := NewBlock(canonicalBlock, m.l2ID, toInclude, blockTime)
				blobPointers := make([]*kzg4844.Blob, len(blobs))
				copy(blobPointers, blobs)

				slot, _ := ethadapter.CalculateSlot(block.Time(), MockGenesisBlock.Time(), SecondsPerSlot)
				if len(blobs) > 0 {
					_ = m.BlobResolver.StoreBlobs(slot, blobs)
				}
				m.miningCh <- block
			})
		}
	}
}

// P2PGossipTx receive rollups to publish from the linked aggregators
func (m *Node) P2PGossipTx(tx *types.Transaction) {
	if atomic.LoadInt32(m.interrupt) == 1 {
		return
	}

	m.mempoolCh <- tx
}

func (m *Node) BroadcastTx(tx types.TxData) {
	m.Network.BroadcastTx(types.NewTx(tx))
}

func (m *Node) Stop() {
	// block all requests
	atomic.StoreInt32(m.interrupt, 1)
	time.Sleep(time.Millisecond * 100)
	m.exitMiningCh <- true
	m.exitCh <- true
}

func (m *Node) BlocksBetween(blockA *types.Header, blockB *types.Block) []*types.Block {
	if bytes.Equal(blockA.Hash().Bytes(), blockB.Hash().Bytes()) {
		return []*types.Block{blockB}
	}
	blocks := make([]*types.Block, 0)
	tempBlock := blockB
	var err error
	for {
		blocks = append(blocks, tempBlock)
		if bytes.Equal(tempBlock.Hash().Bytes(), blockA.Hash().Bytes()) {
			break
		}
		tempBlock, err = m.BlockResolver.FetchBlock(context.Background(), tempBlock.ParentHash())
		if err != nil {
			panic(fmt.Errorf("could not retrieve parent block. Cause: %w", err))
		}
	}
	n := len(blocks)
	result := make([]*types.Block, n)
	for i, block := range blocks {
		result[n-i-1] = block
	}
	return result
}

func (m *Node) CallContract(ethereum.CallMsg) ([]byte, error) {
	return nil, nil
}

func (m *Node) EthClient() *ethclient_ethereum.Client {
	return nil
}

func (m *Node) RemoveSubscription(id uuid.UUID) {
	m.subMu.Lock()
	defer m.subMu.Unlock()
	delete(m.subs, id)
}

func (m *Node) ReconnectIfClosed() error {
	return nil
}

func (m *Node) Alive() bool {
	return true
}

func NewMiner(
	id gethcommon.Address,
	cfg MiningConfig,
	network L1Network,
	statsCollector StatsCollector,
	blobResolver l1.BlobResolver,
	logger gethlog.Logger,
) *Node {
	return &Node{
		l2ID:          id,
		mining:        true,
		cfg:           cfg,
		stats:         statsCollector,
		BlockResolver: NewResolver(),
		BlobResolver:  blobResolver,
		db:            NewTxDB(),
		Network:       network,
		exitCh:        make(chan bool),
		exitMiningCh:  make(chan bool),
		interrupt:     new(int32),
		p2pCh:         make(chan *types.Block),
		miningCh:      make(chan *types.Block),
		// miningCh:         make(chan *BlockWithBlobs),
		canonicalCh:      make(chan *types.Block),
		mempoolCh:        make(chan *types.Transaction),
		headInCh:         make(chan bool),
		headOutCh:        make(chan *types.Block),
		erc20ContractLib: NewERC20ContractLibMock(),
		mgmtContractLib:  NewMgmtContractLibMock(),
		logger:           logger,
		subs:             map[uuid.UUID]*mockSubscription{},
		subMu:            sync.Mutex{},
	}
}

// implements the ethereum.Subscription
type mockSubscription struct {
	id     uuid.UUID
	headCh chan *types.Header
	node   *Node // we hold a reference to the node to unsubscribe ourselves - not ideal but this is just a mock
}

func (sub *mockSubscription) Err() <-chan error {
	return make(chan error)
}

func (sub *mockSubscription) Unsubscribe() {
	sub.node.RemoveSubscription(sub.id)
}

func (sub *mockSubscription) publish(b *types.Block) {
	sub.headCh <- b.Header()
}

func (sub *mockSubscription) publishAll(blocks []*types.Block) {
	for _, b := range blocks {
		sub.publish(b)
	}
}
