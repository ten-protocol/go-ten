package rollupchain

import (
	gethlogs "github.com/ethereum/go-ethereum/log"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
)

func TestInvalidBlocksAreRejected(t *testing.T) {
	// disables most of the geth logs, useful for CI
	handler := gethlogs.StreamHandler(os.Stderr, gethlogs.TerminalFormat(true))
	filteredHandler := gethlogs.LvlFilterHandler(gethlogs.LvlCrit, handler)
	gethlogs.Root().SetHandler(filteredHandler)

	// There are no tests of acceptance of valid chains of blocks. This is because the logic to generate a valid block
	// is non-trivial.
	genesisJSON, err := core.DefaultGenesisBlock().MarshalJSON()
	if err != nil {
		t.Errorf("could not parse genesis JSON: %v", err)
	}
	chain := RollupChain{l1Blockchain: NewL1Blockchain(genesisJSON)}

	invalidHeaders := []types.Header{
		{ParentHash: common.HexToHash("0x0")},                                                            // Unknown ancestor.
		{ParentHash: core.DefaultGenesisBlock().ToBlock(nil).Hash(), Number: big.NewInt(999)},            // Wrong block number.
		{ParentHash: core.DefaultGenesisBlock().ToBlock(nil).Hash(), Number: big.NewInt(1), GasLimit: 1}, // Wrong gas limit.
	}

	for _, header := range invalidHeaders {
		loopHeader := header
		ingestionFailedResponse := chain.insertBlockIntoL1Chain(types.NewBlock(&loopHeader, nil, nil, nil, &trie.StackTrie{}))
		if ingestionFailedResponse == nil {
			t.Errorf("expected block with invalid header to be rejected but was accepted")
		}
	}
}
