package host

import (
	"sync"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/ethereum/go-ethereum/common"
)

// DB allows to access the nodes public nodeDB
type DB struct {
	blockLock        sync.RWMutex
	currentBlockHead common.Hash
	blockDB          map[common.Hash]*types.Header

	rollupLock        sync.RWMutex
	currentRollupHead common.Hash
	rollupDB          map[common.Hash]*nodecommon.Header
}

// NewDB returns a new instance of the Node DB
func NewDB() *DB {
	return &DB{
		blockDB:  map[common.Hash]*types.Header{},
		rollupDB: map[common.Hash]*nodecommon.Header{},
	}
}

// GetCurrentBlockHead returns the current block header (head) of the Node
func (n *DB) GetCurrentBlockHead() *types.Header {
	n.blockLock.RLock()
	current := n.currentBlockHead
	n.blockLock.RUnlock()

	return n.GetBlockHeader(current)
}

// GetBlockHeader returns the block header given the Hash
func (n *DB) GetBlockHeader(hash common.Hash) *types.Header {
	n.blockLock.RLock()
	defer n.blockLock.RUnlock()
	return n.blockDB[hash]
}

// AddBlockHeader adds a types.Header to the known headers
func (n *DB) AddBlockHeader(header *types.Header) {
	n.blockLock.Lock()
	defer n.blockLock.Unlock()

	n.blockDB[header.Hash()] = header

	// update the head if the new height is greater than the existing one
	currentBlockHead := n.blockDB[n.currentBlockHead]
	if currentBlockHead == nil || currentBlockHead.Number.Int64() <= header.Number.Int64() {
		n.currentBlockHead = header.Hash()
	}
}

// GetCurrentRollupHead returns the current rollup header (head) of the Node
func (n *DB) GetCurrentRollupHead() *nodecommon.Header {
	n.rollupLock.RLock()
	current := n.currentRollupHead
	n.rollupLock.RUnlock()

	return n.GetRollupHeader(current)
}

// GetRollupHeader returns the rollup header given the Hash
func (n *DB) GetRollupHeader(hash common.Hash) *nodecommon.Header {
	n.rollupLock.RLock()
	defer n.rollupLock.RUnlock()
	return n.rollupDB[hash]
}

// AddRollupHeader adds a RollupHeader to the known headers
func (n *DB) AddRollupHeader(header *nodecommon.Header) {
	n.rollupLock.Lock()
	defer n.rollupLock.Unlock()

	n.rollupDB[header.Hash()] = header

	// update the head if the new height is greater than the existing one
	currentRollupHead := n.rollupDB[n.currentRollupHead]
	if currentRollupHead == nil || currentRollupHead.Number <= header.Number {
		n.currentRollupHead = header.Hash()
	}
}
