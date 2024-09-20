package ethadapter

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
)

// The number of bits in a BLS scalar that aren't part of a whole byte.
const excessBlobBits = 6 // = math.floor(math.log2(BLS_MODULUS)) % 8

// MakeSidecar builds & returns the BlobTxSidecar and corresponding blob hashes from the raw blob
// data.
func MakeSidecar(blobs []*kzg4844.Blob) (*types.BlobTxSidecar, []gethcommon.Hash, error) {
	sidecar := &types.BlobTxSidecar{}
	var blobHashes []gethcommon.Hash
	for i, blob := range blobs {
		sidecar.Blobs = append(sidecar.Blobs, *blob)
		commitment, err := kzg4844.BlobToCommitment(blob)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot compute KZG commitment of blob %d in tx candidate: %w", i, err)
		}
		sidecar.Commitments = append(sidecar.Commitments, commitment)
		proof, err := kzg4844.ComputeBlobProof(blob, commitment)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot compute KZG proof for fast commitment verification of blob %d in tx candidate: %w", i, err)
		}
		sidecar.Proofs = append(sidecar.Proofs, proof)
		// blobHash := KZGToVersionedHash(commitment)
		// println("created blob hash at ", call, blobHash.Hex())
		blobHashes = append(blobHashes, KZGToVersionedHash(commitment))
	}
	return sidecar, blobHashes, nil
}

// EncodeBlobs takes converts bytes into blobs used for KZG commitment EIP-4844
// transactions on Ethereum.
func EncodeBlobs(data []byte) ([]*kzg4844.Blob, error) {
	data, err := rlp.EncodeToBytes(data)
	if err != nil {
		return nil, err
	}
	var blobs []*kzg4844.Blob
	for len(data) > 0 {
		var b kzg4844.Blob
		data = fillBlobBytes(b[:], data)
		data, err = fillBlobBits(b[:], data)
		if err != nil {
			return nil, err
		}
		blobs = append(blobs, &b)
	}
	return blobs, nil
}

func fillBlobBytes(blob []byte, data []byte) []byte {
	for fieldElement := 0; fieldElement < params.BlobTxFieldElementsPerBlob; fieldElement++ {
		startIdx := fieldElement*32 + 1
		copy(blob[startIdx:startIdx+31], data)
		if len(data) <= 31 {
			return nil
		}
		data = data[31:]
	}
	return data
}

func fillBlobBits(blob []byte, data []byte) ([]byte, error) {
	var acc uint16
	accBits := 0
	for fieldElement := 0; fieldElement < params.BlobTxFieldElementsPerBlob; fieldElement++ {
		if accBits < excessBlobBits && len(data) > 0 {
			acc |= uint16(data[0]) << accBits
			accBits += 8
			data = data[1:]
		}
		blob[fieldElement*32] = uint8(acc & ((1 << excessBlobBits) - 1))
		accBits -= excessBlobBits
		if accBits < 0 {
			// no more data
			break
		}
		acc >>= excessBlobBits
	}
	if accBits > 0 {
		return nil, fmt.Errorf("somehow ended up with %v spare accBits", accBits)
	}
	return data, nil
}

// DecodeBlobs decodes blobs into the data encoded in them accounting for excess blob bits
func DecodeBlobs(blobs []*kzg4844.Blob) ([]byte, error) {
	var rlpData []byte
	for _, blob := range blobs {
		for fieldIndex := 0; fieldIndex < params.BlobTxFieldElementsPerBlob; fieldIndex++ {
			rlpData = append(rlpData, blob[fieldIndex*32+1:(fieldIndex+1)*32]...)
		}
		var acc uint16
		cumulativeBits := 0
		for fieldIndex := 0; fieldIndex < params.BlobTxFieldElementsPerBlob; fieldIndex++ {
			acc |= uint16(blob[fieldIndex*32]) << cumulativeBits
			cumulativeBits += excessBlobBits
			if cumulativeBits >= 8 {
				rlpData = append(rlpData, uint8(acc))
				acc >>= 8
				cumulativeBits -= 8
			}
		}
		if cumulativeBits != 0 {
			return nil, fmt.Errorf("somehow ended up with %v spare cumulative bits", cumulativeBits)
		}
	}
	var outputData []byte
	err := rlp.Decode(bytes.NewReader(rlpData), &outputData)
	return outputData, err
}

// ReconstructRollup decodes and returns the ExtRollup in the blob
func ReconstructRollup(blobs []*kzg4844.Blob) (*common.ExtRollup, error) {
	data, err := DecodeBlobs(blobs)
	if err != nil {
		fmt.Println("Error decoding rollup blob:", err)
	}
	var rollup common.ExtRollup
	if err := rlp.DecodeBytes(data, &rollup); err != nil {
		return nil, fmt.Errorf("could not decode rollup. Cause: %w", err)
	}
	return &rollup, nil
}
