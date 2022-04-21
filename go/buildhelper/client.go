package buildhelper

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/l1client"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

type HelperClient struct {
	node *EthNode
}

func (b *HelperClient) FetchBlockByNumber(n *big.Int) (*types.Block, error) {
	return b.node.FetchBlockByNumber(n)
}

func (b *HelperClient) FetchBlock(id common.Hash) (*types.Block, error) {
	return b.node.FetchBlock(id)
}

func (b *HelperClient) FetchHeadBlock() (*types.Block, uint64) {
	blk, err := b.node.FetchBlockByNumber(nil)
	if err != nil {
		panic(err)
	}
	return blk, blk.Number().Uint64()
}

func (b *HelperClient) IssueTx(tx obscurocommon.EncodedL1Tx) {
	b.node.BroadcastTx(tx)
}

func (b *HelperClient) Info() l1client.Info {
	return l1client.Info{ID: b.node.id}
}

func (b *HelperClient) BlocksBetween(startingBlock *types.Block, lastBlock *types.Block) []*types.Block {
	// TODO this should be a stream
	var blocksBetween []*types.Block
	var err error

	for currentBlk := lastBlock; currentBlk != nil && currentBlk.Hash() != startingBlock.Hash() && currentBlk.ParentHash() != common.HexToHash(""); {
		currentBlk, err = b.FetchBlock(currentBlk.ParentHash())
		if err != nil {
			panic(err)
		}
		blocksBetween = append(blocksBetween, currentBlk)
	}

	return blocksBetween
}

func (b *HelperClient) IsBlockAncestor(block *types.Block, maybeAncestor obscurocommon.L1RootHash) bool {
	if maybeAncestor == block.Hash() || maybeAncestor == obscurocommon.GenesisBlock.Hash() {
		return true
	}

	if block.Number().Int64() == int64(obscurocommon.L1GenesisHeight) {
		return false
	}

	resolvedBlock, err := b.FetchBlock(maybeAncestor)
	if err != nil {
		panic(err)
	}
	if resolvedBlock == nil {
		if resolvedBlock.Number().Int64() >= block.Number().Int64() {
			return false
		}
	}

	p, err := b.FetchBlock(block.ParentHash())
	if err != nil {
		panic(err)
	}
	if p == nil {
		return false
	}

	return b.IsBlockAncestor(p, maybeAncestor)
}

func NewClient(node obscurocommon.L1Node) l1client.Client {
	return &HelperClient{
		node: node.(*EthNode),
	}
}
