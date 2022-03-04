package obscuronode

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	common2 "github.com/obscuronet/obscuro-playground/go/obscuronode/common"
)

type NodeHeaderElement struct {
	Parent      common.Hash
	ID          common.Hash
	Height      uint
	Withdrawals []common2.Withdrawal
	IsRollup    bool
}

type NodeHeader struct {
	currentHeader common.Hash
	lock          sync.RWMutex
	db            map[string]*NodeHeaderElement
}

// NewNodeHead returns a new Head for the node
func NewNodeHead() *NodeHeader {
	return &NodeHeader{db: map[string]*NodeHeaderElement{}}
}

// GetCurrentHead returns the current header of the Node
func (n *NodeHeader) GetCurrentHead() *NodeHeaderElement {
	n.lock.Lock()
	current := n.currentHeader
	n.lock.Unlock()

	return n.GetHeader(current)
}

// GetHeader returns a header given the Hash
func (n *NodeHeader) GetHeader(hash common.Hash) *NodeHeaderElement {
	n.lock.Lock()
	defer n.lock.Unlock()
	return n.db[hash.Hex()]
}

// SetCurrent sets the node head
func (n *NodeHeader) SetCurrent(hash common.Hash) {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.currentHeader = hash
}

// AddHeader adds a NodeHeaderElement to the known headers
func (n *NodeHeader) AddHeader(element *NodeHeaderElement) {
	n.lock.Lock()
	defer n.lock.Unlock()

	// added block headers always point to the latest head
	if !element.IsRollup {
		n.currentHeader = element.ID
	}

	n.db[element.ID.Hex()] = element
}
