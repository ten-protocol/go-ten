package nodetype

import (
	"context"
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
)

func ExportCrossChainData(ctx context.Context, storage storage.Storage, fromSeqNo uint64, toSeqNo uint64) (*common.ExtCrossChainBundle, error) {
	canonicalBatches, err := storage.FetchCanonicalBatchesBetween((ctx), fromSeqNo, toSeqNo)
	if err != nil {
		return nil, err
	}

	if len(canonicalBatches) == 0 {
		return nil, fmt.Errorf("no batches found for export of cross chain data")
	}

	blockHash := canonicalBatches[len(canonicalBatches)-1].Header.L1Proof
	batchHash := canonicalBatches[len(canonicalBatches)-1].Header.Hash()

	block, err := storage.FetchBlock(ctx, blockHash)
	if err != nil {
		return nil, err
	}

	crossChainHashes := make([][]byte, 0)
	for _, batch := range canonicalBatches {
		if batch.Header.CrossChainRoot != gethcommon.BigToHash(gethcommon.Big0) {
			crossChainHashes = append(crossChainHashes, batch.Header.CrossChainRoot.Bytes())
		}
	}

	bundle := &common.ExtCrossChainBundle{
		LastBatchHash:        batchHash, // unused for now.
		L1BlockHash:          block.Hash(),
		L1BlockNum:           big.NewInt(0).Set(block.Header().Number),
		CrossChainRootHashes: crossChainHashes,
	} // todo: check fromSeqNo
	return bundle, nil
}
