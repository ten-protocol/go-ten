package merkle

import (
	"encoding/json"
	"fmt"

	smt "github.com/FantasyJony/openzeppelin-merkle-tree-go/standard_merkle_tree"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/crosschain"
)

func UnmarshalCrossChainTree(serializedTree common.SerializedCrossChainTree) ([][]interface{}, error) {
	xchainTree := make([][]interface{}, 0) // ["v", "0xblablablabla"]
	err := json.Unmarshal(serializedTree, &xchainTree)
	if err != nil {
		return nil, err
	}

	for k, value := range xchainTree {
		xchainTree[k][1] = gethcommon.HexToHash(value[1].(string))
	}
	return xchainTree, nil
}

func ComputeCrossChainRootFromBatches(batches []*common.BatchHeader) (gethcommon.Hash, common.SerializedCrossChainTree, error) {
	xchainTrees := make([][]interface{}, 0)
	for _, batch := range batches {
		if len(batch.CrossChainTree) == 0 {
			// Batch with no outbound messages; nothing to do.
			continue
		}
		xchainTree, err := UnmarshalCrossChainTree(batch.CrossChainTree)
		if err != nil {
			return gethcommon.MaxHash, nil, fmt.Errorf("failed to unmarshal cross chain tree: %w", err)
		}
		xchainTrees = append(xchainTrees, xchainTree...)
	}

	if len(xchainTrees) == 0 {
		return gethcommon.MaxHash, nil, nil
	}

	tree, err := smt.Of(xchainTrees, crosschain.CrossChainEncodings)
	if err != nil {
		return gethcommon.MaxHash, nil, fmt.Errorf("failed to create tree: %w", err)
	}

	serializedTree, err := json.Marshal(xchainTrees)
	if err != nil {
		return gethcommon.MaxHash, nil, err
	}

	return gethcommon.Hash(tree.GetRoot()), serializedTree, nil
}
