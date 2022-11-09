package rollupchain

import (
	"math/big"
	"testing"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"

	gethlog "github.com/ethereum/go-ethereum/log"
)

func TestInvalidBlocksAreRejected(t *testing.T) {
	// There are no tests of acceptance of valid chains of blocks. This is because the logic to generate a valid block
	// is non-trivial.
	genesisJSON, err := core.DefaultGenesisBlock().MarshalJSON()
	if err != nil {
		t.Errorf("could not parse genesis JSON: %v", err)
	}
	logger := log.New(log.DeployerCmp, int(gethlog.LvlDebug), log.SysOut)
	chain := RollupChain{l1Blockchain: NewL1Blockchain(genesisJSON, logger)}

	invalidHeaders := []types.Header{
		{ParentHash: common.HexToHash("0x0")},                                                            // Unknown ancestor.
		{ParentHash: core.DefaultGenesisBlock().ToBlock(nil).Hash(), Number: big.NewInt(999)},            // Wrong block number.
		{ParentHash: core.DefaultGenesisBlock().ToBlock(nil).Hash(), Number: big.NewInt(1), GasLimit: 1}, // Wrong gas limit.
	}

	for _, header := range invalidHeaders {
		loopHeader := header
		_, err := chain.insertBlockIntoL1Chain(types.NewBlock(&loopHeader, nil, nil, nil, &trie.StackTrie{}), false)
		if err == nil {
			t.Errorf("expected block with invalid header to be rejected but was accepted")
		}
	}
}
