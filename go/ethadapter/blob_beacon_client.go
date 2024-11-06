package ethadapter

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"sync"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

const (
	versionMethod        = "eth/v1/node/version"
	specMethod           = "eth/v1/config/spec"
	genesisMethod        = "eth/v1/beacon/genesis"
	sidecarsMethodPrefix = "eth/v1/beacon/blob_sidecars/"
)

// BeaconClient is a thin wrapper over the Beacon APIs.
type BeaconClient interface {
	NodeVersion(ctx context.Context) (string, error)
	ConfigSpec(ctx context.Context) (APIConfigResponse, error)
	BeaconGenesis(ctx context.Context) (APIGenesisResponse, error)
	BeaconBlobSidecars(ctx context.Context, slot uint64, hashes []gethcommon.Hash) (APIGetBlobSidecarsResponse, error)
}

// BlobRetrievalService is a wrapper for clients that can fetch blobs from different sources.
type BlobRetrievalService interface {
	BeaconBlobSidecars(ctx context.Context, slot uint64, hashes []gethcommon.Hash) (APIGetBlobSidecarsResponse, error)
}

type L1BeaconClient struct {
	cl             BeaconClient
	pool           *ClientPool[BlobRetrievalService]
	initLock       sync.Mutex
	genesisTime    uint64
	secondsPerSlot uint64
	timeToSlotFn   TimeToSlot
}

// TimeToSlot cache the function to avoid recomputing it for every block.
type TimeToSlot func(timestamp uint64) (uint64, error)

// BeaconHTTPClient implements BeaconClient. It provides golang types over the basic Beacon API.
type BeaconHTTPClient struct {
	httpClient *BaseHTTPClient
}

func NewBeaconHTTPClient(client *http.Client, baseURL string) *BeaconHTTPClient {
	return &BeaconHTTPClient{
		httpClient: NewBaseHTTPClient(client, baseURL),
	}
}

func (bc *BeaconHTTPClient) request(ctx context.Context, dest any, reqPath string, reqQuery url.Values) error {
	return bc.httpClient.Request(ctx, dest, reqPath, reqQuery)
}

func (bc *BeaconHTTPClient) NodeVersion(ctx context.Context) (string, error) {
	var resp APIVersionResponse
	if err := bc.request(ctx, &resp, versionMethod, nil); err != nil {
		return "", err
	}
	return resp.Data.Version, nil
}

func (bc *BeaconHTTPClient) ConfigSpec(ctx context.Context) (APIConfigResponse, error) {
	var configResp APIConfigResponse
	if err := bc.request(ctx, &configResp, specMethod, nil); err != nil {
		return APIConfigResponse{}, err
	}
	return configResp, nil
}

func (bc *BeaconHTTPClient) BeaconGenesis(ctx context.Context) (APIGenesisResponse, error) {
	var genesisResp APIGenesisResponse
	if err := bc.request(ctx, &genesisResp, genesisMethod, nil); err != nil {
		return APIGenesisResponse{}, err
	}
	return genesisResp, nil
}

func (bc *BeaconHTTPClient) BeaconBlobSidecars(ctx context.Context, slot uint64, _ []gethcommon.Hash) (APIGetBlobSidecarsResponse, error) {
	reqPath := path.Join(sidecarsMethodPrefix, strconv.FormatUint(slot, 10))
	var reqQuery url.Values
	var resp APIGetBlobSidecarsResponse

	err := bc.request(ctx, &resp, reqPath, reqQuery)
	if err != nil {
		return APIGetBlobSidecarsResponse{}, err
	}
	return resp, nil
}

// ClientPool is a simple round-robin client pool
type ClientPool[T any] struct {
	clients []T
	index   int
}

func NewClientPool[T any](clients ...T) *ClientPool[T] {
	return &ClientPool[T]{
		clients: clients,
		index:   0,
	}
}

func (p *ClientPool[T]) Len() int {
	return len(p.clients)
}

func (p *ClientPool[T]) Get() T {
	return p.clients[p.index]
}

func (p *ClientPool[T]) Next() {
	p.index += 1
	if p.index == len(p.clients) {
		p.index = 0
	}
}

// NewL1BeaconClient returns a client for making requests to an L1 consensus layer node.
// Fallbacks are optional clients that will be used for fetching blobs. L1BeaconClient will rotate between
// the `cl` and the fallbacks whenever a client runs into an error while fetching blobs.
func NewL1BeaconClient(cl BeaconClient, fallbacks ...BlobRetrievalService) *L1BeaconClient {
	cs := append([]BlobRetrievalService{cl}, fallbacks...)
	return &L1BeaconClient{
		cl:   cl,
		pool: NewClientPool[BlobRetrievalService](cs...),
	}
}

func (cl *L1BeaconClient) Init(ctx context.Context) error {
	genesis, err := cl.cl.BeaconGenesis(ctx)
	if err != nil {
		return err
	}

	config, err := cl.cl.ConfigSpec(ctx)
	if err != nil {
		return err
	}
	secondsPerSlot := config.Data.SecondsPerSlot
	if secondsPerSlot == 0 {
		return fmt.Errorf("bad value for seconds per slot: %v", config.Data.SecondsPerSlot)
	}
	cl.genesisTime = uint64(genesis.Data.GenesisTime)
	cl.secondsPerSlot = uint64(secondsPerSlot)
	return nil
}

// GetTimeToSlot returns a function that converts a timestamp to a slot number.
func (cl *L1BeaconClient) GetTimeToSlot(ctx context.Context) (TimeToSlot, error) {
	cl.initLock.Lock()
	defer cl.initLock.Unlock()

	if cl.genesisTime == 0 || cl.secondsPerSlot == 0 {
		if err := cl.Init(ctx); err != nil {
			return nil, fmt.Errorf("failed to initialize beacon client: %w", err)
		}
	}

	if cl.timeToSlotFn != nil {
		return cl.timeToSlotFn, nil
	}

	cl.timeToSlotFn = func(timestamp uint64) (uint64, error) {
		if timestamp < cl.genesisTime {
			return 0, fmt.Errorf("provided timestamp (%v) precedes genesis time (%v)", timestamp, cl.genesisTime)
		}
		return (timestamp - cl.genesisTime) / cl.secondsPerSlot, nil
	}
	return cl.timeToSlotFn, nil
}

func (cl *L1BeaconClient) fetchSidecars(ctx context.Context, slot uint64, hashes []gethcommon.Hash) (APIGetBlobSidecarsResponse, error) {
	var errs []error
	for i := 0; i < cl.pool.Len(); i++ {
		f := cl.pool.Get()
		resp, err := f.BeaconBlobSidecars(ctx, slot, hashes)
		if err != nil || len(resp.Data) == 0 {
			cl.pool.Next()
			errs = append(errs, err)
		} else {
			return resp, nil
		}
	}
	return APIGetBlobSidecarsResponse{}, errors.Join(errs...)
}

// GetBlobSidecars fetches blob sidecars that were confirmed in the specified
// L1 block. If hashes are provided, only returns sidecars matching those hashes.
// If no hashes are provided, returns all sidecars for the block.
func (cl *L1BeaconClient) GetBlobSidecars(ctx context.Context, b *types.Header, hashes []gethcommon.Hash) ([]*BlobSidecar, error) {
	slotFn, err := cl.GetTimeToSlot(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get time to slot function: %w", err)
	}
	slot, err := slotFn(b.Time)
	if err != nil {
		return nil, fmt.Errorf("error in converting b.Time to slot: %w", err)
	}

	resp, err := cl.fetchSidecars(ctx, slot, hashes)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch blob sidecars for slot %v block %v: %w", slot, b, err)
	}

	// return all sidecars for block if no hashes provided
	if len(hashes) == 0 {
		return resp.Data, nil
	}

	// match sidecars with provided hashes
	sidecars, err := MatchSidecarsWithHashes(resp.Data, hashes)
	if err != nil {
		return nil, err
	}

	return sidecars, nil
}

// FetchBlobs fetches blobs that were confirmed in the specified L1 block with the
// hashes. Confirms each blob's validity by checking its proof against the commitment, and confirming the commitment
// hashes to the expected value. Returns error if any blob is found invalid.
func (cl *L1BeaconClient) FetchBlobs(ctx context.Context, b *types.Header, hashes []gethcommon.Hash) ([]*kzg4844.Blob, error) {
	blobSidecars, err := cl.GetBlobSidecars(ctx, b, hashes)
	if err != nil {
		return nil, fmt.Errorf("failed to get blob sidecars for Block Header %s: %w", b.Hash().Hex(), err)
	}

	// no hashes were provided, create slice of all hashes from sidecars
	if len(hashes) == 0 {
		hashes = make([]gethcommon.Hash, len(blobSidecars))
		for i, sidecar := range blobSidecars {
			hashes[i] = KZGToVersionedHash(kzg4844.Commitment(sidecar.KZGCommitment))
		}
	}

	return BlobsFromSidecars(blobSidecars, hashes)
}
func BlobsFromSidecars(blobSidecars []*BlobSidecar, hashes []gethcommon.Hash) ([]*kzg4844.Blob, error) {
	if len(blobSidecars) != len(hashes) {
		return nil, fmt.Errorf("number of hashes and blobSidecars mismatch, %d != %d", len(hashes), len(blobSidecars))
	}

	out := make([]*kzg4844.Blob, len(hashes))

	for i, hash := range hashes {
		var matchedSidecar *BlobSidecar
		for _, sidecar := range blobSidecars {
			versionedHash := KZGToVersionedHash(kzg4844.Commitment(sidecar.KZGCommitment))
			if versionedHash == hash {
				matchedSidecar = sidecar
				break
			}
		}

		if matchedSidecar == nil {
			return nil, fmt.Errorf("no matching BlobSidecar found for hash %s", hash.Hex())
		}

		if err := VerifyBlobProof(&matchedSidecar.Blob, kzg4844.Commitment(matchedSidecar.KZGCommitment), kzg4844.Proof(matchedSidecar.KZGProof)); err != nil {
			return nil, fmt.Errorf("blob for hash %s failed verification: %w", hash.Hex(), err)
		}

		out[i] = &matchedSidecar.Blob
	}

	return out, nil
}

// MatchSidecarsWithHashes matches the fetched sidecars with the provided hashes.
func MatchSidecarsWithHashes(fetchedSidecars []*BlobSidecar, hashes []gethcommon.Hash) ([]*BlobSidecar, error) {
	sidecarMap := make(map[gethcommon.Hash]*BlobSidecar)
	for _, sidecar := range fetchedSidecars {
		versionedHash := KZGToVersionedHash(kzg4844.Commitment(sidecar.KZGCommitment))
		sidecarMap[versionedHash] = sidecar
	}

	blobSidecars := make([]*BlobSidecar, len(hashes))
	for i, h := range hashes {
		sidecar, exists := sidecarMap[h]
		if !exists {
			return nil, fmt.Errorf("no matching BlobSidecar found for hash %s", h.Hex())
		}
		blobSidecars[i] = sidecar
	}

	return blobSidecars, nil
}
