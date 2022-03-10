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
}

// NodeHeader allows to access the nodes current l1, l2 headers
type NodeHeader struct {
	blockLock        sync.RWMutex
	currentBlockHead common.Hash
	blockDB          map[common.Hash]*BlockHeader

	rollupLock        sync.RWMutex
	currentRollupHead common.Hash
	rollupDB          map[common.Hash]*RollupHeader
}

// NewNodeHeader returns a new instance of the Node Headers
func NewNodeHeader() *NodeHeader {
	return &NodeHeader{
		blockDB:  map[common.Hash]*BlockHeader{},
		rollupDB: map[common.Hash]*RollupHeader{},
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
	n.blockLock.RLock()
	defer n.blockLock.RUnlock()
	return n.blockDB[hash]
}

// AddBlockHeader adds a BlockHeader to the known headers
func (n *NodeHeader) AddBlockHeader(header *BlockHeader) {
	n.blockLock.Lock()
	defer n.blockLock.Unlock()

	n.blockDB[header.ID] = header

	// update the head if the new height is greater than the existing one
	currentBlockHead := n.blockDB[n.currentBlockHead]
	if currentBlockHead == nil || currentBlockHead.Height <= header.Height {
		n.currentBlockHead = header.ID
	}
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
	n.rollupLock.RLock()
	defer n.rollupLock.RUnlock()
	return n.rollupDB[hash]
}

// AddRollupHeader adds a RollupHeader to the known headers
func (n *NodeHeader) AddRollupHeader(header *RollupHeader) {
	n.rollupLock.Lock()
	defer n.rollupLock.Unlock()

	n.rollupDB[header.ID] = header

	// update the head if the new height is greater than the existing one
	currentRollupHead := n.rollupDB[n.currentRollupHead]
	if currentRollupHead == nil || currentRollupHead.Height <= header.Height {
		n.currentRollupHead = header.ID
	}
}
