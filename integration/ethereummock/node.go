package ethereummock

import (
	"bytes"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/obscuronet/go-obscuro/go/common/gethutil"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/common"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	ethclient_ethereum "github.com/ethereum/go-ethereum/ethclient"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
)

type L1Network interface {
	// BroadcastBlock - send the block and the parent to make sure there are no gaps
	BroadcastBlock(b common.EncodedL1Block, p common.EncodedL1Block)
	BroadcastTx(tx *types.Transaction)
}

type MiningConfig struct {
	PowTime common.Latency
	LogFile string
}

type TxDB interface {
	Txs(block *types.Block) (map[common.TxHash]*types.Transaction, bool)
	AddTxs(*types.Block, map[common.TxHash]*types.Transaction)
}

type StatsCollector interface {
	// L1Reorg registers when a miner has to process a reorg (a winning block from a fork)
	L1Reorg(id gethcommon.Address)
}

type NotifyNewBlock interface {
	MockedNewHead(b common.EncodedL1Block, p common.EncodedL1Block)
	MockedNewFork(b []common.EncodedL1Block)
}

type Node struct {
	l2ID     gethcommon.Address // the address of the Obscuro node this client is dedicated to
	cfg      MiningConfig
	clients  []NotifyNewBlock
	Network  L1Network
	mining   bool
	stats    StatsCollector
	Resolver db.BlockResolver
	db       TxDB

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

func (m *Node) SendTransaction(tx *types.Transaction) error {
	m.Network.BroadcastTx(tx)
	return nil
}

func (m *Node) TransactionReceipt(_ gethcommon.Hash) (*types.Receipt, error) {
	// all transactions are immediately processed
	return &types.Receipt{
		Status: types.ReceiptStatusSuccessful,
	}, nil
}

func (m *Node) Nonce(gethcommon.Address) (uint64, error) {
	return 0, nil
}

// BlockListener is not used in the mock
func (m *Node) BlockListener() (chan *types.Header, ethereum.Subscription) {
	return make(chan *types.Header), &mockSubscription{}
}

func (m *Node) BlockNumber() (uint64, error) {
	blk, found := m.Resolver.FetchHeadBlock()
	if !found {
		return 0, ethereum.NotFound
	}
	return blk.NumberU64(), nil
}

func (m *Node) BlockByNumber(n *big.Int) (*types.Block, error) {
	if n.Int64() == 0 {
		return common.GenesisBlock, nil
	}
	// TODO this should be a method in the resolver
	var f bool
	blk, found := m.Resolver.FetchHeadBlock()
	if !found {
		return nil, ethereum.NotFound
	}
	for !bytes.Equal(blk.ParentHash().Bytes(), common.GenesisHash.Bytes()) {
		if blk.NumberU64() == n.Uint64() {
			return blk, nil
		}

		blk, f = m.Resolver.FetchBlock(blk.ParentHash())
		if !f {
			return nil, fmt.Errorf("block in the chain without a parent")
		}
	}
	return nil, ethereum.NotFound
}

func (m *Node) BlockByHash(id gethcommon.Hash) (*types.Block, error) {
	blk, f := m.Resolver.FetchBlock(id)
	if !f {
		return nil, fmt.Errorf("blk not found")
	}
	return blk, nil
}

func (m *Node) FetchHeadBlock() (*types.Block, bool) {
	return m.Resolver.FetchHeadBlock()
}

func (m *Node) Info() ethadapter.Info {
	return ethadapter.Info{
		L2ID: m.l2ID,
	}
}

func (m *Node) IsBlockAncestor(block *types.Block, proof common.L1RootHash) bool {
	return m.Resolver.IsBlockAncestor(block, proof)
}

func (m *Node) BalanceAt(gethcommon.Address, *big.Int) (*big.Int, error) {
	panic("not implemented")
}

// Start runs an infinite loop that listens to the two block producing channels and processes them.
func (m *Node) Start() {
	if m.mining {
		// This starts the mining
		go m.startMining()
	}

	m.Resolver.StoreBlock(common.GenesisBlock)
	head := m.setHead(common.GenesisBlock)

	for {
		select {
		case p2pb := <-m.p2pCh: // Received from peers
			_, received := m.Resolver.FetchBlock(p2pb.Hash())
			// only process blocks if they haven't been processed before
			if !received {
				head = m.processBlock(p2pb, head)
			}

		case mb := <-m.miningCh: // Received from the local mining
			head = m.processBlock(mb, head)
			if bytes.Equal(head.Hash().Bytes(), mb.Hash().Bytes()) { // Ignore the locally produced block if someone else found one already
				p, found := m.Resolver.ParentBlock(mb)
				if !found {
					panic("noo")
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
	m.Resolver.StoreBlock(b)
	_, f := m.Resolver.FetchBlock(b.Header().ParentHash)

	// only proceed if the parent is available
	if !f {
		m.logger.Info(fmt.Sprintf("Parent block not found=b_%d", common.ShortHash(b.Header().ParentHash)))
		return head
	}

	// Ignore superseded blocks
	if b.NumberU64() <= head.NumberU64() {
		return head
	}

	// Check for Reorgs
	if !m.Resolver.IsAncestor(b, head) {
		m.stats.L1Reorg(m.l2ID)
		fork, err := gethutil.LCA(head, b, m.Resolver)
		if err != nil {
			panic(err)
		}
		m.logger.Info(
			fmt.Sprintf("L1Reorg new=b_%d(%d), old=b_%d(%d), fork=b_%d(%d)", common.ShortHash(b.Hash()), b.NumberU64(), common.ShortHash(head.Hash()), head.NumberU64(), common.ShortHash(fork.Hash()), fork.NumberU64()))
		return m.setFork(m.BlocksBetween(fork, b))
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

	// notify the clients
	for _, c := range m.clients {
		t := c
		encodedBlock, err := common.EncodeBlock(b)
		if err != nil {
			panic(fmt.Errorf("could not encode block. Cause: %w", err))
		}
		if b.NumberU64() == common.L1GenesisHeight {
			go t.MockedNewHead(encodedBlock, nil)
		} else {
			p, f := m.Resolver.ParentBlock(b)
			if !f {
				panic("This should not happen")
			}
			encodedParentBlock, err := common.EncodeBlock(p)
			if err != nil {
				panic(fmt.Errorf("could not encode parent block. Cause: %w", err))
			}
			go t.MockedNewHead(encodedBlock, encodedParentBlock)
		}
	}
	m.canonicalCh <- b

	return b
}

func (m *Node) setFork(blocks []*types.Block) *types.Block {
	head := blocks[len(blocks)-1]
	if atomic.LoadInt32(m.interrupt) == 1 {
		return head
	}

	fork := make([]common.EncodedL1Block, len(blocks))
	for i, block := range blocks {
		encodedBlock, err := common.EncodeBlock(block)
		if err != nil {
			panic(fmt.Errorf("could not encode block. Cause: %w", err))
		}
		fork[i] = encodedBlock
	}

	// notify the clients
	for _, c := range m.clients {
		go c.MockedNewFork(fork)
	}
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
			mempool = m.removeCommittedTransactions(canonicalBlock, mempool, m.Resolver, m.db)

			// notify the existing mining go routine to stop mining
			atomic.StoreInt32(interrupt, 1)
			c := int32(0)
			interrupt = &c

			// Generate a random number, and wait for that number of ms. Equivalent to PoW
			// Include all rollups received during this period.
			common.ScheduleInterrupt(m.cfg.PowTime(), interrupt, func() {
				toInclude := findNotIncludedTxs(canonicalBlock, mempool, m.Resolver, m.db)
				// todo - iterate through the rollup transactions and include only the ones with the proof on the canonical chain
				if atomic.LoadInt32(m.interrupt) == 1 {
					return
				}

				m.miningCh <- common.NewBlock(canonicalBlock, m.l2ID, toInclude)
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

func (m *Node) AddClient(client NotifyNewBlock) {
	m.clients = append(m.clients, client)
}

func (m *Node) BlocksBetween(blockA *types.Block, blockB *types.Block) []*types.Block {
	if bytes.Equal(blockA.Hash().Bytes(), blockB.Hash().Bytes()) {
		return []*types.Block{blockA}
	}
	blocks := make([]*types.Block, 0)
	tempBlock := blockB
	var found bool
	for {
		blocks = append(blocks, tempBlock)
		if bytes.Equal(tempBlock.Hash().Bytes(), blockA.Hash().Bytes()) {
			break
		}
		tempBlock, found = m.Resolver.ParentBlock(tempBlock)
		if !found {
			panic("should not happen")
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

func NewMiner(
	id gethcommon.Address,
	cfg MiningConfig,
	network L1Network,
	statsCollector StatsCollector,
) *Node {
	return &Node{
		l2ID:             id,
		mining:           true,
		cfg:              cfg,
		stats:            statsCollector,
		Resolver:         NewResolver(),
		db:               NewTxDB(),
		Network:          network,
		exitCh:           make(chan bool),
		exitMiningCh:     make(chan bool),
		interrupt:        new(int32),
		p2pCh:            make(chan *types.Block),
		miningCh:         make(chan *types.Block),
		canonicalCh:      make(chan *types.Block),
		mempoolCh:        make(chan *types.Transaction),
		headInCh:         make(chan bool),
		headOutCh:        make(chan *types.Block),
		erc20ContractLib: NewERC20ContractLibMock(),
		mgmtContractLib:  NewMgmtContractLibMock(),
		logger:           log.New(log.EthereumL1Cmp, int(gethlog.LvlInfo), cfg.LogFile, log.NodeIDKey, id),
	}
}

// implements the ethereum.Subscription
type mockSubscription struct{}

func (sub *mockSubscription) Err() <-chan error {
	c := make(chan error, 2) // size 2, so that we don't block
	// drop an error in to this channel to remind callers that this client does not stream.
	c <- ethadapter.ErrSubscriptionNotSupported
	return c
}

func (sub *mockSubscription) Unsubscribe() {
}
