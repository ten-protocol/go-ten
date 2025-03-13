package contractlib

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/contracts/generated/RollupContract"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/ethadapter"
)

type RollupContractLib interface {
	ContractLib
	PopulateAddRollup(t *common.L1RollupTx, blobs []*kzg4844.Blob, signature common.RollupSignature) (types.TxData, error)
	BlobHasher() ethadapter.BlobHasher
}

type rollupContractLibImpl struct {
	addr        *gethcommon.Address
	contractABI abi.ABI
	logger      gethlog.Logger
}

func NewRollupContractLib(addr *gethcommon.Address, logger gethlog.Logger) RollupContractLib {
	return &rollupContractLibImpl{
		addr:        addr,
		contractABI: ethadapter.RollupContractABI,
		logger:      logger,
	}
}

func (r *rollupContractLibImpl) IsMock() bool {
	return false
}

func (r *rollupContractLibImpl) BlobHasher() ethadapter.BlobHasher {
	return ethadapter.KZGToVersionedHasher{}
}

func (r *rollupContractLibImpl) PopulateAddRollup(t *common.L1RollupTx, blobs []*kzg4844.Blob, signature common.RollupSignature) (types.TxData, error) {
	decodedRollup, err := common.DecodeRollup(t.Rollup)
	if err != nil {
		return nil, fmt.Errorf("could not decode rollup. Cause: %w", err)
	}

	metaRollup := RollupContract.StructsMetaRollup{
		Hash:                decodedRollup.Hash(),
		Signature:           signature,
		FirstSequenceNumber: big.NewInt(int64(decodedRollup.Header.FirstBatchSeqNo)),
		LastSequenceNumber:  big.NewInt(int64(decodedRollup.Header.LastBatchSeqNo)),
		BlockBindingHash:    decodedRollup.Header.CompressionL1Head,
		BlockBindingNumber:  decodedRollup.Header.CompressionL1Number,
		CrossChainRoot:      decodedRollup.Header.CrossChainRoot,
		LastBatchHash:       decodedRollup.Header.LastBatchHash,
	}

	data, err := r.contractABI.Pack(
		ethadapter.AddRollupMethod,
		metaRollup,
	)
	if err != nil {
		return nil, fmt.Errorf("could not pack the call data. Cause: %w", err)
	}

	var blobHashes []gethcommon.Hash
	var sidecar *types.BlobTxSidecar

	// Using blobs created here (they are verified that the hash matches with the blobs from the enclave)
	if sidecar, blobHashes, err = ethadapter.MakeSidecar(blobs, r.BlobHasher()); err != nil {
		return nil, fmt.Errorf("failed to make sidecar: %w", err)
	}

	return &types.BlobTx{
		To:         *r.addr,
		Data:       data,
		BlobHashes: blobHashes,
		Sidecar:    sidecar,
	}, nil
}

func (r *rollupContractLibImpl) DecodeTx(tx *types.Transaction) (common.L1TenTransaction, error) {
	if tx.To() == nil || tx.To().Hex() != r.addr.Hex() || len(tx.Data()) == 0 {
		return nil, nil
	}

	if tx.Type() == types.BlobTxType {
		return &common.L1RollupHashes{
			BlobHashes: tx.BlobHashes(),
		}, nil
	}
	//else {
	//	return nil, fmt.Errorf("invalid transaction type: %v", tx.Type())
	//}
	return nil, nil
}

func (r *rollupContractLibImpl) GetContractAddr() *gethcommon.Address {
	return r.addr
}
