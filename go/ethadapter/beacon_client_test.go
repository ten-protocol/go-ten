package ethadapter

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"github.com/ten-protocol/go-ten/go/common"
	"testing"
)

const spareBlobBits = 6 // = math.floor(math.log2(BLS_MODULUS)) % 8

func TestBlobsFromSidecars(t *testing.T) {
	indices := []uint64{5, 7, 2}

	// blobs should be returned in order of their indices in the hashes array regardless
	// of the sidecar ordering
	index0, sidecar0 := makeTestBlobSidecar(indices[0])
	index1, sidecar1 := makeTestBlobSidecar(indices[1])
	index2, sidecar2 := makeTestBlobSidecar(indices[2])

	hashes := []IndexedBlobHash{index0, index1, index2}

	// put the sidecars in scrambled order to confirm error
	sidecars := []*BlobSidecar{sidecar2, sidecar0, sidecar1}
	_, err := blobsFromSidecars(sidecars, hashes)
	require.Error(t, err)

	// too few sidecars should error
	sidecars = []*BlobSidecar{sidecar0, sidecar1}
	_, err = blobsFromSidecars(sidecars, hashes)
	require.Error(t, err)

	// correct order should work
	sidecars = []*BlobSidecar{sidecar0, sidecar1, sidecar2}
	blobs, err := blobsFromSidecars(sidecars, hashes)
	require.NoError(t, err)
	// confirm order by checking first blob byte against expected index
	for i := range blobs {
		require.Equal(t, byte(indices[i]), blobs[i][0])
	}

	// mangle a proof to make sure it's detected
	badProof := *sidecar0
	badProof.KZGProof[11]++
	sidecars[1] = &badProof
	_, err = blobsFromSidecars(sidecars, hashes)
	require.Error(t, err)

	// mangle a commitment to make sure it's detected
	badCommitment := *sidecar0
	badCommitment.KZGCommitment[13]++
	sidecars[1] = &badCommitment
	_, err = blobsFromSidecars(sidecars, hashes)
	require.Error(t, err)

	// mangle a hash to make sure it's detected
	sidecars[1] = sidecar0
	hashes[2].Hash[17]++
	_, err = blobsFromSidecars(sidecars, hashes)
	require.Error(t, err)
}

func TestBlobsFromSidecars_EmptySidecarList(t *testing.T) {
	hashes := []IndexedBlobHash{}
	sidecars := []*BlobSidecar{}
	blobs, err := blobsFromSidecars(sidecars, hashes)
	require.NoError(t, err)
	require.Empty(t, blobs, "blobs should be empty when no sidecars are provided")
}

func TestClientPoolSingle(t *testing.T) {
	p := NewClientPool[int](1)
	for i := 0; i < 10; i++ {
		require.Equal(t, 1, p.Get())
		p.MoveToNext()
	}
}
func TestClientPoolSeveral(t *testing.T) {
	p := NewClientPool[int](0, 1, 2, 3)
	for i := 0; i < 25; i++ {
		require.Equal(t, i%4, p.Get())
		p.MoveToNext()
	}
}

func TestBlobEncoding(t *testing.T) {
	// Example data
	extRlp := createRollup(4444)
	encRollup, err := common.EncodeRollup(&extRlp)

	// Encode data into blobs
	blobs, err := EncodeBlobs(encRollup)
	if err != nil {
		fmt.Println("Error encoding rollup blob:", err)
	}

	blobsPtr := make([]*kzg4844.Blob, len(blobs))
	for i := range blobs {
		blobsPtr[i] = &blobs[i]
	}

	rollup, err := reconstructRollup(blobsPtr)
	if err != nil {
		fmt.Println("Error reconstructing rollup:", err)
		return
	}

	fmt.Println("Reconstructed rollup:", rollup)
}

// Function to reconstruct rollup from blobs
func reconstructRollup(blobs []*kzg4844.Blob) (*common.ExtRollup, error) {
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

//// Function to reconstruct rollup from blobs
//func reconstructRollup(blobs []kzg4844.Blob) (*common.ExtRollup, error) {
//	var serializedData []byte
//	for _, blob := range blobs {
//		for i := 0; i < len(blob); i += 32 {
//			// We need to make sure to not go beyond the actual length of the blob
//			if i+32 <= len(blob) {
//				serializedData = append(serializedData, blob[i+1:i+32]...)
//			}
//		}
//	}
//
//	var rollup common.ExtRollup
//	if err := rlp.DecodeBytes(serializedData, &rollup); err != nil {
//		return nil, fmt.Errorf("could not decode rollup. Cause: %w", err)
//	}
//
//	return &rollup, nil
//}

func makeTestBlobSidecar(index uint64) (IndexedBlobHash, *BlobSidecar) {
	blob := kzg4844.Blob{}
	// make first byte of test blob match its index so we can easily verify if is returned in the
	// expected order
	blob[0] = byte(index)
	commit, _ := kzg4844.BlobToCommitment(&blob)
	proof, _ := kzg4844.ComputeBlobProof(&blob, commit)
	hash := KZGToVersionedHash(commit)

	idh := IndexedBlobHash{
		Index: index,
		Hash:  hash,
	}
	sidecar := BlobSidecar{
		Index:         Uint64String(index),
		Blob:          blob,
		KZGCommitment: Bytes48(commit),
		KZGProof:      Bytes48(proof),
	}
	return idh, &sidecar
}

// Function to encode data into blobs
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

func createRollup(lastBatch int64) common.ExtRollup {
	header := common.RollupHeader{
		LastBatchSeqNo: uint64(lastBatch),
	}

	rollup := common.ExtRollup{
		Header: &header,
	}

	return rollup
}
