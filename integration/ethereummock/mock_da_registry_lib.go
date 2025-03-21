package ethereummock

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/ethadapter/contractlib"
)

type mockDataAvailabilityRegistryLib struct{}

func NewDataAvailabilityRegistryLibMock() contractlib.DataAvailabilityRegistryLib {
	return &mockDataAvailabilityRegistryLib{}
}

func (m *mockDataAvailabilityRegistryLib) IsMock() bool {
	return true
}

func (m *mockDataAvailabilityRegistryLib) BlobHasher() ethadapter.BlobHasher {
	return MockBlobHasher{}
}

func (m *mockDataAvailabilityRegistryLib) GetContractAddr() *gethcommon.Address {
	return &RollupTxAddr
}

func (m *mockDataAvailabilityRegistryLib) DecodeTx(tx *types.Transaction) (common.L1TenTransaction, error) {
	if tx.To() == nil || tx.To().Hex() != RollupTxAddr.Hex() {
		return nil, nil
	}

	if tx.Type() == types.BlobTxType {
		return &common.L1RollupHashes{
			BlobHashes: tx.BlobHashes(),
		}, nil
	}
	return nil, nil
}

func (m *mockDataAvailabilityRegistryLib) PopulateAddRollup(_ *common.L1RollupTx, blobs []*kzg4844.Blob, _ common.RollupSignature) (types.TxData, error) {
	var err error
	var blobHashes []gethcommon.Hash
	var sidecar *types.BlobTxSidecar
	if sidecar, blobHashes, err = ethadapter.MakeSidecar(blobs, MockBlobHasher{}); err != nil {
		return nil, fmt.Errorf("failed to make sidecar: %w", err)
	}

	hashesTx := common.L1RollupHashes{BlobHashes: blobHashes}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(hashesTx); err != nil {
		panic(err)
	}
	blobTx := types.BlobTx{
		To:         RollupTxAddr,
		Data:       buf.Bytes(),
		BlobHashes: blobHashes,
		Sidecar:    sidecar,
	}
	// Force wait before publishing tx for in-mem test
	time.Sleep(time.Second * 1)
	return &blobTx, nil
}
