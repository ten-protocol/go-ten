package ethadapter

import (
	"context"
	"errors"
	"math/big"
	"net/http"
	"testing"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"

	"github.com/ethereum/go-ethereum/rlp"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/stretchr/testify/require"
	"github.com/ten-protocol/go-ten/go/common"
)

const (
	vHash1 = "0x012b7a6a22399aa9eecd8eda6ec658679e81be21af6ff296116aee205e2218f2"
	vHash2 = "0x012374e04a848591844b75bc2f500318cf640552379b5e3a1a77bb828620690e"
)

func TestBlobsFromSidecars(t *testing.T) {
	indices := []uint64{5, 7, 2}

	// blobs should be returned in order of their indices in the hashes array regardless
	// of the sidecar ordering
	hash0, sidecar0 := makeTestBlobSidecar(indices[0])
	hash1, sidecar1 := makeTestBlobSidecar(indices[1])
	hash2, sidecar2 := makeTestBlobSidecar(indices[2])

	hashes := []gethcommon.Hash{hash0, hash1, hash2}

	// too few sidecars should error
	sidecars := []*BlobSidecar{sidecar0, sidecar1}
	_, err := BlobsFromSidecars(sidecars, hashes)
	require.Error(t, err)

	// correct order should work
	sidecars = []*BlobSidecar{sidecar0, sidecar1, sidecar2}
	blobs, err := BlobsFromSidecars(sidecars, hashes)
	require.NoError(t, err)
	for i := range blobs {
		require.Equal(t, byte(indices[i]), blobs[i][0])
	}

	badProof := *sidecar0
	badProof.KZGProof[11]++
	sidecars[1] = &badProof
	_, err = BlobsFromSidecars(sidecars, hashes)
	require.Error(t, err)

	badCommitment := *sidecar0
	badCommitment.KZGCommitment[13]++
	sidecars[1] = &badCommitment
	_, err = BlobsFromSidecars(sidecars, hashes)
	require.Error(t, err)

	sidecars[1] = sidecar0
	hashes[2][17]++
	_, err = BlobsFromSidecars(sidecars, hashes)
	require.Error(t, err)
}

func TestEmptyBlobSidecars(t *testing.T) {
	var hashes []gethcommon.Hash
	var sidecars []*BlobSidecar
	blobs, err := BlobsFromSidecars(sidecars, hashes)
	require.NoError(t, err)
	require.Empty(t, blobs, "blobs should be empty when no sidecars are provided")
}

func TestClientPoolSingle(t *testing.T) {
	p := NewClientPool[int](1)
	for i := 0; i < 10; i++ {
		require.Equal(t, 1, p.Get())
		p.Next()
	}
}

func TestClientPoolSeveral(t *testing.T) {
	p := NewClientPool[int](0, 1, 2, 3)
	for i := 0; i < 25; i++ {
		require.Equal(t, i%4, p.Get())
		p.Next()
	}
}

func TestBlobEncoding(t *testing.T) {
	extRlp := createRollup(4444)
	encRollup, err := common.EncodeRollup(&extRlp)
	if err != nil {
		t.Errorf("error encoding rollup")
	}

	// Encode data into blobs
	blobs, err := EncodeBlobs(encRollup, gethlog.New())
	if err != nil {
		t.Errorf("error encoding blobs: %s", err)
	}

	rollup, err := ReconstructRollup(blobs)
	if err != nil {
		t.Errorf("error reconstructing rollup: %s", err)
	}

	if big.NewInt(int64(rollup.Header.LastBatchSeqNo)).Cmp(big.NewInt(4444)) != 0 {
		t.Errorf("rollup was not decoded correctly")
	}
}

// Write a test that will create a rollup that exceeds 128kb and ensure that it is split into multiple blobs
// and then reassembled correctly
func TestBlobEncodingLarge(t *testing.T) {
	// make this rollup larger than 128kb
	extRlp := createLargeRollup(4445)
	encRollup, _ := common.EncodeRollup(&extRlp)
	_, err := EncodeBlobs(encRollup, gethlog.New())
	require.Error(t, err)
}

func TestBlobArchiveClient(t *testing.T) {
	t.Skipf("TODO need to fix this")
	client := NewArchivalHTTPClient(new(http.Client), "https://eth-beacon-chain.drpc.org/rest/")
	vHashes := []gethcommon.Hash{gethcommon.HexToHash(vHash1), gethcommon.HexToHash(vHash2)}
	ctx := context.Background()

	resp, err := client.BeaconBlobSidecars(ctx, 1, vHashes)
	require.NoError(t, err)

	require.Len(t, resp.Data, 2)
	require.NotNil(t, client)
}

func TestBeaconClientFallback(t *testing.T) {
	indices := []uint64{5, 7, 2}
	hash0, sidecar0 := makeTestBlobSidecar(indices[0])
	hash1, sidecar1 := makeTestBlobSidecar(indices[1])
	hash2, sidecar2 := makeTestBlobSidecar(indices[2])

	hashes := []gethcommon.Hash{hash0, hash1, hash2}
	sidecars := []*BlobSidecar{sidecar0, sidecar1, sidecar2}

	ctx := context.Background()

	mockPrimary := &MockBeaconClient{}
	mockFallback := &MockBlobRetrievalService{}

	client := NewL1BeaconClient(mockPrimary, mockFallback)

	mockPrimary.On("BeaconGenesis", ctx).Return(APIGenesisResponse{Data: ReducedGenesisData{GenesisTime: 10}}, nil)
	mockPrimary.On("ConfigSpec", ctx).Return(APIConfigResponse{Data: ReducedConfigData{SecondsPerSlot: 2}}, nil)
	mockPrimary.On("BeaconBlobSidecars", ctx, uint64(1), hashes).Return(APIGetBlobSidecarsResponse{}, errors.New("404 not found"))
	mockFallback.On("BeaconBlobSidecars", ctx, uint64(1), hashes).Return(APIGetBlobSidecarsResponse{Data: sidecars}, nil)

	header := &types.Header{Time: 12}
	resp, err := client.GetBlobSidecars(ctx, header, hashes)
	require.NoError(t, err)
	require.Equal(t, sidecars, resp)

	mockFallback.On("BeaconBlobSidecars", ctx, uint64(2), hashes).Return(APIGetBlobSidecarsResponse{}, errors.New("404 not found"))
	mockPrimary.On("BeaconBlobSidecars", ctx, uint64(2), hashes).Return(APIGetBlobSidecarsResponse{Data: sidecars}, nil)

	header = &types.Header{Time: 14}
	resp, err = client.GetBlobSidecars(ctx, header, hashes)
	require.NoError(t, err)
	require.Equal(t, sidecars, resp)

	mockPrimary.AssertExpectations(t)
	mockFallback.AssertExpectations(t)
}

// MockBeaconClient is a mock implementation used only in these tests
type MockBeaconClient struct {
	mock.Mock
}

func (m *MockBeaconClient) NodeVersion(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockBeaconClient) ConfigSpec(ctx context.Context) (APIConfigResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(APIConfigResponse), args.Error(1)
}

func (m *MockBeaconClient) BeaconGenesis(ctx context.Context) (APIGenesisResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(APIGenesisResponse), args.Error(1)
}

func (m *MockBeaconClient) BeaconBlobSidecars(ctx context.Context, slot uint64, hashes []gethcommon.Hash) (APIGetBlobSidecarsResponse, error) {
	args := m.Called(ctx, slot, hashes)
	return args.Get(0).(APIGetBlobSidecarsResponse), args.Error(1)
}

// MockBlobRetrievalService is a mock implementation of the BlobRetrievalService interface
type MockBlobRetrievalService struct {
	mock.Mock
}

func (m *MockBlobRetrievalService) BeaconBlobSidecars(ctx context.Context, slot uint64, hashes []gethcommon.Hash) (APIGetBlobSidecarsResponse, error) {
	args := m.Called(ctx, slot, hashes)
	return args.Get(0).(APIGetBlobSidecarsResponse), args.Error(1)
}

func makeTestBlobSidecar(index uint64) (gethcommon.Hash, *BlobSidecar) {
	blob := kzg4844.Blob{}
	// make first byte of test blob match its index so we can easily verify if is returned in the
	// expected order
	blob[0] = byte(index)
	commit, _ := kzg4844.BlobToCommitment(&blob)
	proof, _ := kzg4844.ComputeBlobProof(&blob, commit)
	hash := KZGToVersionedHash(commit)

	sidecar := BlobSidecar{
		Index:         Uint64String(index),
		Blob:          blob,
		KZGCommitment: Bytes48(commit),
		KZGProof:      Bytes48(proof),
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

func createLargeRollup(seqNo int64) common.ExtRollup {
	header := common.RollupHeader{
		LastBatchSeqNo: uint64(seqNo),
	}
	largeData := make([]byte, 130*1024) // 130KB
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	bytes, _ := rlp.EncodeToBytes(largeData)
	rollup := common.ExtRollup{
		Header:        &header,
		BatchPayloads: bytes,
	}
	return rollup
}
