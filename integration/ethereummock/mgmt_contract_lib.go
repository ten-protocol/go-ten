package ethereummock

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/integration/datagenerator"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"
)

var (
	depositTxAddr          = datagenerator.RandomAddress()
	rollupTxAddr           = datagenerator.RandomAddress()
	storeSecretTxAddr      = datagenerator.RandomAddress()
	requestSecretTxAddr    = datagenerator.RandomAddress()
	initializeSecretTxAddr = datagenerator.RandomAddress()
	// MgmtContractAddresses make all these addresses available for the host to know what receipts will be forwarded to the enclave
	MgmtContractAddresses = []gethcommon.Address{
		depositTxAddr,
		rollupTxAddr,
		storeSecretTxAddr,
		requestSecretTxAddr,
		initializeSecretTxAddr,
	}
)

// mockContractLib is an implementation of the mgmtcontractlib.MgmtContractLib
// it creates ethereum mocked transactions from common.L1Transaction
// and converts ethereum mocked transactions to common.L1Transaction
type mockContractLib struct{}

func NewMgmtContractLibMock() mgmtcontractlib.MgmtContractLib {
	return &mockContractLib{}
}

func (m *mockContractLib) IsMock() bool {
	return true
}

func (m *mockContractLib) GetContractAddr() *gethcommon.Address {
	return &rollupTxAddr
}

func (m *mockContractLib) DecodeTx(tx *types.Transaction) ethadapter.L1Transaction {
	// Do not decode erc20 transactions, this is the responsibility
	// of the erc20 contract lib.
	if tx.To().Hex() == depositTxAddr.Hex() {
		return nil
	}

	return decodeTx(tx)
}

func (m *mockContractLib) CreateRollup(tx *ethadapter.L1RollupTx) types.TxData {
	return encodeTx(tx, rollupTxAddr)
}

func (m *mockContractLib) CreateBlobRollup(t *ethadapter.L1RollupTx) (types.TxData, error) {
	var err error
	blobs := encodeBlobs(t.Rollup)
	if err != nil {
		return nil, fmt.Errorf("failed to convert rollup to blobs: %w", err)
	}

	var blobHashes []gethcommon.Hash
	var sidecar *types.BlobTxSidecar

	if sidecar, blobHashes, err = makeSidecar(blobs); err != nil {
		return nil, fmt.Errorf("failed to make sidecar: %w", err)
	}
	return &types.BlobTx{
		To:         rollupTxAddr,
		BlobHashes: blobHashes,
		Sidecar:    sidecar,
	}, nil
}

func (m *mockContractLib) CreateRequestSecret(tx *ethadapter.L1RequestSecretTx) types.TxData {
	return encodeTx(tx, requestSecretTxAddr)
}

func (m *mockContractLib) CreateRespondSecret(tx *ethadapter.L1RespondSecretTx, _ bool) types.TxData {
	return encodeTx(tx, storeSecretTxAddr)
}

func (m *mockContractLib) CreateInitializeSecret(tx *ethadapter.L1InitializeSecretTx) types.TxData {
	return encodeTx(tx, initializeSecretTxAddr)
}

func (m *mockContractLib) GetHostAddressesMsg() (ethereum.CallMsg, error) {
	return ethereum.CallMsg{}, nil
}

func (m *mockContractLib) DecodeHostAddressesResponse([]byte) ([]string, error) {
	return []string{""}, nil
}

func (m *mockContractLib) GetImportantContractKeysMsg() (ethereum.CallMsg, error) {
	return ethereum.CallMsg{}, nil
}

func (m *mockContractLib) DecodeImportantContractKeysResponse([]byte) ([]string, error) {
	return []string{""}, nil
}

func (m *mockContractLib) SetImportantContractMsg(string, gethcommon.Address) (ethereum.CallMsg, error) {
	return ethereum.CallMsg{}, nil
}

func (m *mockContractLib) GetImportantAddressCallMsg(string) (ethereum.CallMsg, error) {
	return ethereum.CallMsg{}, nil
}

func (m *mockContractLib) DecodeImportantAddressResponse([]byte) (gethcommon.Address, error) {
	return gethcommon.Address{}, nil
}

func decodeTx(tx *types.Transaction) ethadapter.L1Transaction {
	if len(tx.Data()) == 0 {
		panic("Data cannot be 0 in the mock implementation")
	}

	// prepare byte buffer
	buf := bytes.NewBuffer(tx.Data())
	dec := gob.NewDecoder(buf)

	// in the mock implementation we use the To address field to specify the L1 operation (rollup/storesecret/requestsecret)
	// the mock implementation does not process contracts
	// so this is a way that we can differentiate different contract calls
	var t ethadapter.L1Transaction
	switch tx.To().Hex() {
	case rollupTxAddr.Hex():
		t = &ethadapter.L1RollupTx{}
	case storeSecretTxAddr.Hex():
		t = &ethadapter.L1RespondSecretTx{}
	case depositTxAddr.Hex():
		t = &ethadapter.L1DepositTx{}
	case requestSecretTxAddr.Hex():
		t = &ethadapter.L1RequestSecretTx{}
	case initializeSecretTxAddr.Hex():
		t = &ethadapter.L1InitializeSecretTx{}
	default:
		panic("unexpected type")
	}

	// decode to interface implementation
	if err := dec.Decode(t); err != nil {
		panic(err)
	}
	return t
}

func encodeTx(tx ethadapter.L1Transaction, opType gethcommon.Address) types.TxData {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(tx); err != nil {
		panic(err)
	}

	// the mock implementation does not process contract calls
	// this uses the To address to distinguish between different contract calls / different l1 transactions
	return &types.LegacyTx{
		Data: buf.Bytes(),
		To:   &opType,
	}
}

func encodeBlobs(data []byte) []kzg4844.Blob {
	blobs := []kzg4844.Blob{{}}
	blobIndex := 0
	fieldIndex := -1
	for i := 0; i < len(data); i += 31 {
		fieldIndex++
		if fieldIndex == params.BlobTxFieldElementsPerBlob {
			blobs = append(blobs, kzg4844.Blob{})
			blobIndex++
			fieldIndex = 0
		}
		max := i + 31
		if max > len(data) {
			max = len(data)
		}
		copy(blobs[blobIndex][fieldIndex*32+1:], data[i:max])
	}
	return blobs
}

// MakeSidecar builds & returns the BlobTxSidecar and corresponding blob hashes from the raw blob
// data.
func makeSidecar(blobs []kzg4844.Blob) (*types.BlobTxSidecar, []gethcommon.Hash, error) {
	sidecar := &types.BlobTxSidecar{}
	blobHashes := []gethcommon.Hash{}
	for i, blob := range blobs {
		sidecar.Blobs = append(sidecar.Blobs, blob)
		commitment, err := kzg4844.BlobToCommitment(&blob)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot compute KZG commitment of blob %d in tx candidate: %w", i, err)
		}
		sidecar.Commitments = append(sidecar.Commitments, commitment)
		proof, err := kzg4844.ComputeBlobProof(&blob, commitment)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot compute KZG proof for fast commitment verification of blob %d in tx candidate: %w", i, err)
		}
		sidecar.Proofs = append(sidecar.Proofs, proof)
		blobHashes = append(blobHashes, kzgToVersionedHash(commitment))
	}
	return sidecar, blobHashes, nil
}

// kzgToVersionedHash computes the versioned hash of a blob-commitment. Implemented here as it's not exposed by geth.
func kzgToVersionedHash(commitment kzg4844.Commitment) (out gethcommon.Hash) {
	hasher := sha256.New()
	return kzg4844.CalcBlobHashV1(hasher, &commitment)
}
