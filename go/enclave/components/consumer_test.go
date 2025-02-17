package components

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestInvalidBlocksAreRejected(t *testing.T) {
	t.Skipf("skipping test that relied on l1blockchain previously which hasn't been used in a while")

	// todo - how does this test even work, storage is never set and we attempt to fetch head block?
	blockConsumer := l1BlockProcessor{}

	invalidHeaders := []types.Header{
		{ParentHash: common.HexToHash("0x0")},                                                         // Unknown ancestor.
		{ParentHash: core.DefaultGenesisBlock().ToBlock().Hash(), Number: big.NewInt(999)},            // Wrong block number.
		{ParentHash: core.DefaultGenesisBlock().ToBlock().Hash(), Number: big.NewInt(1), GasLimit: 1}, // Wrong gas limit.
	}

	for _, header := range invalidHeaders {
		loopHeader := header
		_, err := blockConsumer.ingestBlock(context.Background(), &loopHeader)
		if err == nil {
			t.Errorf("expected block with invalid header to be rejected but was accepted")
		}
	}
}
