package ethereummock

import (
	"context"
	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"math/big"
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/stretchr/testify/require"
	"github.com/ten-protocol/go-ten/go/common"
)

func TestGetVersion(t *testing.T) {
	l := testlog.Logger()

	beaconApi := NewBeaconMock(l, uint64(1), uint64(1), 8000)
	t.Cleanup(func() {
		_ = beaconApi.Close()
	})
	require.NoError(t, beaconApi.Start("127.0.0.1"))
}

func Test404NotFound(t *testing.T) {
	l := testlog.Logger()

	beaconApi := NewBeaconMock(l, uint64(1), uint64(12), 8000)
	t.Cleanup(func() {
		_ = beaconApi.Close()
	})
	require.NoError(t, beaconApi.Start("127.0.0.1"))

	cl := ethadapter.NewL1BeaconClient(ethadapter.NewBeaconHTTPClient(new(http.Client), beaconApi.BeaconAddr()))

	var hashes []gethcommon.Hash
	_, err := cl.FetchBlobs(context.Background(), &types.Header{Number: big.NewInt(10), Time: 120}, hashes)
	require.ErrorIs(t, err, ethereum.NotFound)
}

func TestBlobsFromSidecars(t *testing.T) {
	indices := []uint64{5, 7, 2}

	// blobs should be returned in order of their indices in the hashes array regardless
	// of the sidecar ordering
	hash0, sidecar0 := makeTestBlobSidecar(indices[0])
	hash1, sidecar1 := makeTestBlobSidecar(indices[1])
	hash2, sidecar2 := makeTestBlobSidecar(indices[2])

	hashes := []gethcommon.Hash{hash0, hash1, hash2}

	// put the sidecars in scrambled order to confirm error
	sidecars := []*ethadapter.BlobSidecar{sidecar2, sidecar0, sidecar1}
	_, err := ethadapter.BlobsFromSidecars(sidecars, hashes)
	require.Error(t, err)

	// too few sidecars should error
	sidecars = []*ethadapter.BlobSidecar{sidecar0, sidecar1}
	_, err = ethadapter.BlobsFromSidecars(sidecars, hashes)
	require.Error(t, err)

	// correct order should work
	sidecars = []*ethadapter.BlobSidecar{sidecar0, sidecar1, sidecar2}
	blobs, err := ethadapter.BlobsFromSidecars(sidecars, hashes)
	require.NoError(t, err)
	// confirm order by checking first blob byte against expected index
	for i := range blobs {
		require.Equal(t, byte(indices[i]), blobs[i][0])
	}

	// mangle a proof to make sure it's detected
	badProof := *sidecar0
	badProof.KZGProof[11]++
	sidecars[1] = &badProof
	_, err = ethadapter.BlobsFromSidecars(sidecars, hashes)
	require.Error(t, err)

	// mangle a commitment to make sure it's detected
	badCommitment := *sidecar0
	badCommitment.KZGCommitment[13]++
	sidecars[1] = &badCommitment
	_, err = ethadapter.BlobsFromSidecars(sidecars, hashes)
	require.Error(t, err)

	// mangle a hash to make sure it's detected
	sidecars[1] = sidecar0
	hashes[2][17]++
	_, err = ethadapter.BlobsFromSidecars(sidecars, hashes)
	require.Error(t, err)
}

func TestBlobsFromSidecars_EmptySidecarList(t *testing.T) {
	var hashes []gethcommon.Hash
	var sidecars []*ethadapter.BlobSidecar
	blobs, err := ethadapter.BlobsFromSidecars(sidecars, hashes)
	require.NoError(t, err)
	require.Empty(t, blobs, "blobs should be empty when no sidecars are provided")
}

func TestClientPoolSingle(t *testing.T) {
	p := ethadapter.NewClientPool[int](1)
	for i := 0; i < 10; i++ {
		require.Equal(t, 1, p.Get())
		p.MoveToNext()
	}
}

func TestClientPoolSeveral(t *testing.T) {
	p := ethadapter.NewClientPool[int](0, 1, 2, 3)
	for i := 0; i < 25; i++ {
		require.Equal(t, i%4, p.Get())
		p.MoveToNext()
	}
}

func TestBlobEncoding(t *testing.T) {
	extRlp := createRollup(4444)
	encRollup, err := common.EncodeRollup(&extRlp)
	if err != nil {
		t.Errorf("error encoding rollup")
	}

	// Encode data into blobs
	blobs, err := ethadapter.EncodeBlobs(encRollup)
	if err != nil {
		t.Errorf("error encoding blobs: %s", err)
	}

	rollup, err := ethadapter.ReconstructRollup(blobs)
	if err != nil {
		t.Errorf("error reconstructing rollup: %s", err)
	}

	if big.NewInt(int64(rollup.Header.LastBatchSeqNo)).Cmp(big.NewInt(4444)) != 0 {
		t.Errorf("rollup was not decoded correctly")
	}
}

func makeTestBlobSidecar(index uint64) (gethcommon.Hash, *ethadapter.BlobSidecar) {
	blob := kzg4844.Blob{}
	// make first byte of test blob match its index so we can easily verify if is returned in the
	// expected order
	blob[0] = byte(index)
	commit, _ := kzg4844.BlobToCommitment(&blob)
	proof, _ := kzg4844.ComputeBlobProof(&blob, commit)
	hash := ethadapter.KZGToVersionedHash(commit)

	sidecar := ethadapter.BlobSidecar{
		Index:         ethadapter.Uint64String(index),
		Blob:          blob,
		KZGCommitment: ethadapter.Bytes48(commit),
		KZGProof:      ethadapter.Bytes48(proof),
	}
	return hash, &sidecar
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
