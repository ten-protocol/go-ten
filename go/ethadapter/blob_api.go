package ethadapter

import "github.com/ethereum/go-ethereum/crypto/kzg4844"

type BlobSidecar struct {
	Blob          Blob               `json:"blob"`
	Index         uint64             `json:"index"`
	KZGCommitment kzg4844.Commitment `json:"kzg_commitment"`
	KZGProof      kzg4844.Proof      `json:"kzg_proof"`
}

type APIBlobSidecar struct {
	Index             uint64                  `json:"index"`
	Blob              Blob                    `json:"blob"`
	KZGCommitment     kzg4844.Commitment      `json:"kzg_commitment"`
	KZGProof          kzg4844.Proof           `json:"kzg_proof"`
	SignedBlockHeader SignedBeaconBlockHeader `json:"signed_block_header"`
	// The inclusion-proof of the blob-sidecar into the beacon-block is ignored,
	// since we verify blobs by their versioned hashes against the execution-layer block instead.
}

func (sc *APIBlobSidecar) BlobSidecar() *BlobSidecar {
	return &BlobSidecar{
		Blob:          sc.Blob,
		Index:         sc.Index,
		KZGCommitment: sc.KZGCommitment,
		KZGProof:      sc.KZGProof,
	}
}

type SignedBeaconBlockHeader struct {
	Message BeaconBlockHeader `json:"message"`
	// signature is ignored, since we verify blobs against EL versioned-hashes
}

type BeaconBlockHeader struct {
	Slot          uint64 `json:"slot"`
	ProposerIndex uint64 `json:"proposer_index"`
	ParentRoot    []byte `json:"parent_root"`
	StateRoot     []byte `json:"state_root"`
	BodyRoot      []byte `json:"body_root"`
}

type APIGetBlobSidecarsResponse struct {
	Data []*APIBlobSidecar `json:"data"`
}

type ReducedGenesisData struct {
	GenesisTime uint64 `json:"genesis_time"`
}

type APIGenesisResponse struct {
	Data ReducedGenesisData `json:"data"`
}

type ReducedConfigData struct {
	SecondsPerSlot uint64 `json:"SECONDS_PER_SLOT"`
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
