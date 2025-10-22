package ethadapter

import (
	"crypto/sha256"
	"fmt"
	"reflect"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

type BlobSidecar struct {
	// Blob is the actual blob data, a large byte array (up to 128 KB) that contains the off-chain data
	Blob kzg4844.Blob `json:"blob"`

	// Index is the position of this blob within the transaction's blob array
	// It's represented as a string to accommodate uint64 values in JSON
	Index Uint64String `json:"index"`

	// KZGCommitment is a cryptographic commitment to the blob content
	// This is what the EVM sees and uses to verify blob availability without accessing the full blob data
	// It's significantly smaller than the blob itself, reducing on-chain storage requirements
	KZGCommitment Bytes48 `json:"kzg_commitment"`

	// KZGProof is a zero-knowledge proof that allows verification of small portions of the blob
	// without needing the entire blob data
	// It's used to prove that the commitment corresponds to the actual blob data
	KZGProof Bytes48 `json:"kzg_proof"`

	//TODO
	SignedBlockHeader SignedBeaconBlockHeader `json:"signed_block_header"`

	//TODO
	InclusionProof []Bytes32 `json:"kzg_commitment_inclusion_proof"`
}

type SignedBeaconBlockHeader struct {
	Message   BeaconBlockHeader `json:"message"`
	Signature hexutil.Bytes     `json:"signature"`
}

type BeaconBlockHeader struct {
	Slot          Uint64String `json:"slot"`
	ProposerIndex Uint64String `json:"proposer_index"`
	ParentRoot    Bytes32      `json:"parent_root"`
	StateRoot     Bytes32      `json:"state_root"`
	BodyRoot      Bytes32      `json:"body_root"`
}
type APIGetBlobSidecarsResponse struct {
	Data []*BlobSidecar `json:"data"`
}

type ReducedGenesisData struct {
	GenesisTime Uint64String `json:"genesis_time"`
}

type APIGenesisResponse struct {
	Data ReducedGenesisData `json:"data"`
}

type ReducedConfigData struct {
	SecondsPerSlot Uint64String `json:"SECONDS_PER_SLOT"`
}

type APIConfigResponse struct {
	Data ReducedConfigData `json:"data"`
}

type APIVersionResponse struct {
	Data VersionInformation `json:"data"`
}

type VersionInformation struct {
	Version string `json:"version"`
}

// Uint64String is a decimal string representation of an uint64, for usage in the Beacon API JSON encoding
type Uint64String uint64

func (v Uint64String) MarshalText() (out []byte, err error) {
	out = strconv.AppendUint(out, uint64(v), 10)
	return
}

func (v *Uint64String) UnmarshalText(b []byte) error {
	n, err := strconv.ParseUint(string(b), 0, 64)
	if err != nil {
		return err
	}
	*v = Uint64String(n)
	return nil
}

type Bytes48 [48]byte

func (b *Bytes48) UnmarshalJSON(text []byte) error {
	return hexutil.UnmarshalFixedJSON(reflect.TypeOf(b), text, b[:])
}

func (b *Bytes48) UnmarshalText(text []byte) error {
	return hexutil.UnmarshalFixedText("Bytes32", text, b[:])
}

func (b Bytes48) MarshalText() ([]byte, error) {
	return hexutil.Bytes(b[:]).MarshalText()
}

func (b Bytes48) String() string {
	return hexutil.Encode(b[:])
}

// TerminalString implements log.TerminalStringer, formatting a string for console
// output during logging.
func (b Bytes48) TerminalString() string {
	return fmt.Sprintf("%x..%x", b[:3], b[45:])
}

type Bytes32 [32]byte

func (b *Bytes32) UnmarshalJSON(text []byte) error {
	return hexutil.UnmarshalFixedJSON(reflect.TypeOf(b), text, b[:])
}

func (b *Bytes32) UnmarshalText(text []byte) error {
	return hexutil.UnmarshalFixedText("Bytes32", text, b[:])
}

func (b Bytes32) MarshalText() ([]byte, error) {
	return hexutil.Bytes(b[:]).MarshalText()
}

func (b Bytes32) String() string {
	return hexutil.Encode(b[:])
}

// KZGToVersionedHash computes the versioned hash of a blob-commitment, as used in a blob-tx.
func KZGToVersionedHash(commitment kzg4844.Commitment) (out common.Hash) {
	hasher := sha256.New()
	return kzg4844.CalcBlobHashV1(hasher, &commitment)
}

// VerifyBlobProof verifies that the given blob and proof corresponds to the given commitment
func VerifyBlobProof(blob *kzg4844.Blob, commitment kzg4844.Commitment, proof kzg4844.Proof) error {
	return kzg4844.VerifyBlobProof(blob, commitment, proof)
}
