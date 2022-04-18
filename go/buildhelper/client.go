package buildhelper

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

type BuildHelperClient struct {
	node *EthNode
}

func (b *BuildHelperClient) FetchBlockByNumber(n *big.Int) (*types.Block, error) {
	return b.node.FetchBlockByNumber(n)
}

func (b *BuildHelperClient) FetchBlock(id common.Hash) (*types.Block, error) {
	return b.node.FetchBlock(id)
}

func (b *BuildHelperClient) FetchHeadBlock() (*types.Block, uint64) {
	blk, err := b.node.FetchBlockByNumber(nil)
	if err != nil {
		panic(err)
	}
	return blk, blk.Number().Uint64()
}

func (b *BuildHelperClient) IssueTx(tx obscurocommon.EncodedL1Tx) {
	b.node.BroadcastTx(tx)
}

func (b *BuildHelperClient) Info() ethclient.Info {
	return ethclient.Info{ID: b.node.id}
}

func (b *BuildHelperClient) BlocksBetween(startingBlock *types.Block, lastBlock *types.Block) []*types.Block {
	//TODO this should be a stream
	var blocksBetween []*types.Block
	var err error
	for currentBlk := lastBlock; currentBlk != nil || currentBlk.Hash() != startingBlock.Hash(); {
		if currentBlk.ParentHash() == common.HexToHash("") {
			break
		}
		currentBlk, err = b.FetchBlock(currentBlk.ParentHash())
		if err != nil {
			panic(err)
		}
		blocksBetween = append(blocksBetween, currentBlk)
	}

	return blocksBetween
}

func (b *BuildHelperClient) IsBlockAncestor(block *types.Block, maybeAncestor obscurocommon.L1RootHash) bool {
	if maybeAncestor == block.Hash() {
		return true
	}

	if maybeAncestor == obscurocommon.GenesisBlock.Hash() {
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

func NewClient(node obscurocommon.L1Node) ethclient.Client {
	return &BuildHelperClient{
		node: node.(*EthNode),
	}
}
