package obscuronode

import (
	"sync"

	common2 "github.com/obscuronet/obscuro-playground/go/obscuronode/common"

	"github.com/ethereum/go-ethereum/common"
)

// BlockHeader is holds the header info for the l1 blocks
type BlockHeader struct {
	ID     common.Hash
	Parent common.Hash
	Height uint
}

// RollupHeader is holds the header info for the l2 rollups
type RollupHeader struct {
	ID     common.Hash
	Parent common.Hash
	Height uint

	Withdrawals []common2.Withdrawal
	IsRollup    bool
}

// NodeHeader allows to access the nodes current l1, l2 headers
type NodeHeader struct {
	blockLock        sync.RWMutex
	currentBlockHead common.Hash
	blockDB          map[string]*BlockHeader

	rollupLock        sync.RWMutex
	currentRollupHead common.Hash
	rollupDB          map[string]*RollupHeader
}

// NewNodeHeader returns a new instance of the NodeHeaders
func NewNodeHeader() *NodeHeader {
	return &NodeHeader{
		blockDB:  map[string]*BlockHeader{},
		rollupDB: map[string]*RollupHeader{},
	}
}

// GetCurrentBlockHead returns the current block header (head) of the Node
func (n *NodeHeader) GetCurrentBlockHead() *BlockHeader {
	n.blockLock.RLock()
	current := n.currentBlockHead
	n.blockLock.RUnlock()

	return n.GetBlockHeader(current)
}

// GetBlockHeader returns the block header given the Hash
func (n *NodeHeader) GetBlockHeader(hash common.Hash) *BlockHeader {
	n.blockLock.Lock()
	defer n.blockLock.Unlock()
	return n.blockDB[hash.Hex()]
}

// SetCurrentBlockHead sets the node block head
func (n *NodeHeader) SetCurrentBlockHead(hash common.Hash) {
	n.blockLock.Lock()
	defer n.blockLock.Unlock()
	n.currentBlockHead = hash
}

// AddBlockHeader adds a BlockHeader to the known headers
func (n *NodeHeader) AddBlockHeader(blockHeader *BlockHeader) {
	n.blockLock.Lock()
	defer n.blockLock.Unlock()

	n.blockDB[blockHeader.ID.Hex()] = blockHeader
}

// GetCurrentRollupHead returns the current rollup header (head) of the Node
func (n *NodeHeader) GetCurrentRollupHead() *RollupHeader {
	n.rollupLock.RLock()
	current := n.currentRollupHead
	n.rollupLock.RUnlock()

	return n.GetRollupHeader(current)
}

// GetRollupHeader returns the rollup header given the Hash
func (n *NodeHeader) GetRollupHeader(hash common.Hash) *RollupHeader {
	n.rollupLock.Lock()
	defer n.rollupLock.Unlock()
	return n.rollupDB[hash.Hex()]
}

// SetCurrentRollupHead sets the node rollup head
func (n *NodeHeader) SetCurrentRollupHead(hash common.Hash) {
	n.rollupLock.Lock()
	defer n.rollupLock.Unlock()
	n.currentRollupHead = hash
}

// AddRollupHeader adds a RollupHeader to the known headers
func (n *NodeHeader) AddRollupHeader(header *RollupHeader) {
	n.rollupLock.Lock()
	defer n.rollupLock.Unlock()

	n.rollupDB[header.ID.Hex()] = header
}
