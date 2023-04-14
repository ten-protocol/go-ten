package components

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
)

func TestInvalidBlocksAreRejected(t *testing.T) {
	// todo - how does this test even work, storage is never set and we attempt to fetch head block?
	blockConsumer := blockConsumer{}

	invalidHeaders := []types.Header{
		{ParentHash: common.HexToHash("0x0")},                                                            // Unknown ancestor.
		{ParentHash: core.DefaultGenesisBlock().ToBlock(nil).Hash(), Number: big.NewInt(999)},            // Wrong block number.
		{ParentHash: core.DefaultGenesisBlock().ToBlock(nil).Hash(), Number: big.NewInt(1), GasLimit: 1}, // Wrong gas limit.
	}

	for _, header := range invalidHeaders {
		loopHeader := header
		_, err := blockConsumer.ingestBlock(types.NewBlock(&loopHeader, nil, nil, nil, &trie.StackTrie{}), false)
		if err == nil {
			t.Errorf("expected block with invalid header to be rejected but was accepted")
		}
	}
}
