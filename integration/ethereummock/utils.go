package ethereummock

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/trie"
	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ethereum/go-ethereum/core/types"
)

// findNotIncludedTxs - given a list of transactions, it keeps only the ones that were not included in the block
// todo (#1491) - inefficient
func findNotIncludedTxs(head *types.Block, txs []*types.Transaction, r *blockResolverInMem, db TxDB) []*types.Transaction {
	included := allIncludedTransactions(head, r, db)
	return removeExisting(txs, included)
}

func allIncludedTransactions(b *types.Block, r *blockResolverInMem, db TxDB) map[common.TxHash]*types.Transaction {
	val, found := db.Txs(b)
	if found {
		return val
	}
	if b.NumberU64() == common.L1GenesisHeight {
		return makeMap(b.Transactions())
	}
	newMap := make(map[common.TxHash]*types.Transaction)
	p, err := r.FetchFullBlock(context.Background(), b.ParentHash())
	if err != nil {
		panic(fmt.Errorf("should not happen. Could not retrieve parent. Cause: %w", err))
	}
	for k, v := range allIncludedTransactions(p, r, db) {
		newMap[k] = v
	}
	for _, tx := range b.Transactions() {
		newMap[tx.Hash()] = tx
	}
	db.AddTxs(b, newMap)
	return newMap
}

func removeExisting(base []*types.Transaction, toRemove map[common.TxHash]*types.Transaction) (r []*types.Transaction) {
	for _, t := range base {
		_, f := toRemove[t.Hash()]
		if !f {
			r = append(r, t)
		}
	}
	return
}

func makeMap(txs types.Transactions) map[common.TxHash]*types.Transaction {
	m := make(map[common.TxHash]*types.Transaction)
	for _, tx := range txs {
		m[tx.Hash()] = tx
	}
	return m
}

// EncodedL1Block the encoded version of an L1 block.
type (
	EncodedL1Block []byte
	extblock       struct {
		Header      *types.Header
		Txs         []*types.Transaction
		Uncles      []*types.Header
		Withdrawals []*types.Withdrawal `rlp:"optional"`
	}
)

func EncodeBlock(b *types.Block) (EncodedL1Block, error) {
	return json.Marshal(&extblock{
		Header: b.Header(),
		Txs:    b.Transactions(),
	})
}

func (eb EncodedL1Block) DecodeBlock() (*types.Block, error) {
	var b extblock
	if err := json.Unmarshal(eb, &b); err != nil {
		return nil, fmt.Errorf("could not decode block from bytes. Cause: %w", err)
	}
	return types.NewBlock(b.Header, &types.Body{
		Transactions: b.Txs,
		Uncles:       nil,
		Withdrawals:  nil,
	}, nil, trie.NewStackTrie(nil)), nil
}
