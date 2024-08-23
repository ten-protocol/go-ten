package ethadapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

const (
	versionMethod        = "eth/v1/node/version"
	specMethod           = "eth/v1/config/spec"
	genesisMethod        = "eth/v1/beacon/genesis"
	sidecarsMethodPrefix = "eth/v1/beacon/blob_sidecars/"
)

// L1BeaconClient is a high level golang client for the Beacon API.
type L1BeaconClient struct {
	cl           BeaconClient
	pool         *ClientPool[BlobSideCarsFetcher]
	initLock     sync.Mutex
	timeToSlotFn TimeToSlotFn
}

// BeaconClient is a thin wrapper over the Beacon APIs.
type BeaconClient interface {
	NodeVersion(ctx context.Context) (string, error)
	ConfigSpec(ctx context.Context) (APIConfigResponse, error)
	BeaconGenesis(ctx context.Context) (APIGenesisResponse, error)
	BeaconBlobSideCars(ctx context.Context, slot uint64, hashes []IndexedBlobHash) (APIGetBlobSidecarsResponse, error)
}

// BlobSideCarsFetcher is a thin wrapper over the Beacon APIs.
type BlobSideCarsFetcher interface {
	BeaconBlobSideCars(ctx context.Context, slot uint64, hashes []IndexedBlobHash) (APIGetBlobSidecarsResponse, error)
}

// BeaconHTTPClient implements BeaconClient. It provides golang types over the basic Beacon API.
type BeaconHTTPClient struct {
	client  *http.Client
	baseURL string
}

func NewBeaconHTTPClient(client *http.Client, baseURL string) *BeaconHTTPClient {
	return &BeaconHTTPClient{client: client, baseURL: baseURL}
}

func (bc *BeaconHTTPClient) request(ctx context.Context, dest any, reqPath string, reqQuery url.Values) error {
	baseURL, err := url.Parse(bc.baseURL)
	if err != nil {
		return fmt.Errorf("failed to parse base URL: %w", err)
	}

	reqURL, err := baseURL.Parse(reqPath)
	if err != nil {
		return fmt.Errorf("failed to parse request path: %w", err)
	}

	reqURL.RawQuery = reqQuery.Encode()

	headers := http.Header{}
	headers.Add("Accept", "application/json")

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header = headers

	resp, err := bc.client.Do(req)
	if err != nil {
		return fmt.Errorf("http Get failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		errMsg, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed request with status %d: %s: %w", resp.StatusCode, string(errMsg), ethereum.NotFound)
	} else if resp.StatusCode != http.StatusOK {
		errMsg, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed request with status %d: %s", resp.StatusCode, string(errMsg))
	}

	err = json.NewDecoder(resp.Body).Decode(dest)
	if err != nil {
		return err
	}

	return nil
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

func (bc *BeaconHTTPClient) BeaconBlobSideCars(ctx context.Context, slot uint64, hashes []IndexedBlobHash) (APIGetBlobSidecarsResponse, error) {
	return testBeaconBlobSideCars(bc, ctx, slot, hashes)
}

func testBeaconBlobSideCars(bc *BeaconHTTPClient, ctx context.Context, slot uint64, hashes []IndexedBlobHash) (APIGetBlobSidecarsResponse, error) {
	reqPath := path.Join(sidecarsMethodPrefix, strconv.FormatUint(slot, 10))
	var reqQuery url.Values
	var resp APIGetBlobSidecarsResponse

	err := bc.request(ctx, &resp, reqPath, reqQuery)

	if err != nil {
		return APIGetBlobSidecarsResponse{}, err
	}

	response := resp.Data
	if len(response) < len(hashes) {
		return APIGetBlobSidecarsResponse{}, fmt.Errorf("expected at least %d blobs for slot %d but only got %d", len(hashes), slot, len(response))
	}

	outputsFound := make([]bool, len(hashes))
	filteredResponse := APIGetBlobSidecarsResponse{}

	for _, blobItem := range response {
		commitment := kzg4844.Commitment(blobItem.KZGCommitment)
		proof := kzg4844.Proof(blobItem.KZGProof)
		blobPtr := &blobItem.Blob
		versionedHash := KZGToVersionedHash(commitment)

		// Try to match the versioned hash with one of the provided hashes
		var outputIdx int
		var found bool
		for outputIdx = range hashes {
			if hashes[outputIdx].Hash == versionedHash {
				if outputsFound[outputIdx] {
					// Duplicate, skip this one
					break
				}
				found = true
				outputsFound[outputIdx] = true
				break
			}
		}

		if !found {
			continue
		}

		// Verify the blob proof
		err = kzg4844.VerifyBlobProof(blobPtr, commitment, proof)
		if err != nil {
			return APIGetBlobSidecarsResponse{}, fmt.Errorf("failed to verify blob proof for blob at slot(%d) with index(%d)", slot, blobItem.Index)
		}

		// Add the matching item to the filtered response
		filteredResponse.Data = append(filteredResponse.Data, blobItem)
	}

	// Ensure all expected hashes were found
	for i, found := range outputsFound {
		if !found {
			return APIGetBlobSidecarsResponse{}, fmt.Errorf("missing blob %v in slot %v, can't reconstruct rollup payload", hashes[i], slot)
		}
	}

	return filteredResponse, nil
}

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

func (p *ClientPool[T]) MoveToNext() {
	p.index += 1
	if p.index == len(p.clients) {
		p.index = 0
	}
}

// NewL1BeaconClient returns a client for making requests to an L1 consensus layer node.
// Fallbacks are optional clients that will be used for fetching blobs. L1BeaconClient will rotate between
// the `cl` and the fallbacks whenever a client runs into an error while fetching blobs.
func NewL1BeaconClient(cl BeaconClient, fallbacks ...BlobSideCarsFetcher) *L1BeaconClient {
	cs := append([]BlobSideCarsFetcher{cl}, fallbacks...)
	return &L1BeaconClient{
		cl:   cl,
		pool: NewClientPool[BlobSideCarsFetcher](cs...)}
}

type TimeToSlotFn func(timestamp uint64) (uint64, error)

// GetTimeToSlotFn returns a function that converts a timestamp to a slot number.
func (cl *L1BeaconClient) GetTimeToSlotFn(ctx context.Context) (TimeToSlotFn, error) {
	cl.initLock.Lock()
	defer cl.initLock.Unlock()
	if cl.timeToSlotFn != nil {
		return cl.timeToSlotFn, nil
	}

	genesis, err := cl.cl.BeaconGenesis(ctx)
	if err != nil {
		return nil, err
	}

	config, err := cl.cl.ConfigSpec(ctx)
	if err != nil {
		return nil, err
	}

	genesisTime := uint64(genesis.Data.GenesisTime)
	secondsPerSlot := uint64(config.Data.SecondsPerSlot)
	if secondsPerSlot == 0 {
		return nil, fmt.Errorf("got bad value for seconds per slot: %v", config.Data.SecondsPerSlot)
	}
	cl.timeToSlotFn = func(timestamp uint64) (uint64, error) {
		if timestamp < genesisTime {
			return 0, fmt.Errorf("provided timestamp (%v) precedes genesis time (%v)", timestamp, genesisTime)
		}
		return (timestamp - genesisTime) / secondsPerSlot, nil
	}
	return cl.timeToSlotFn, nil
}

func (cl *L1BeaconClient) fetchSidecars(ctx context.Context, slot uint64, hashes []IndexedBlobHash) (APIGetBlobSidecarsResponse, error) {
	var errs []error
	for i := 0; i < cl.pool.Len(); i++ {
		f := cl.pool.Get()
		println("FETCHING BeaconBlobSideCars with hashes: ", hashes[0].Hash.Hex())
		resp, err := f.BeaconBlobSideCars(ctx, slot, hashes)
		if err != nil {
			cl.pool.MoveToNext()
			errs = append(errs, err)
		} else {
			return resp, nil
		}
	}
	return APIGetBlobSidecarsResponse{}, errors.Join(errs...)
}

//// GetBlobSidecars fetches blob sidecars that were confirmed in the specified
//// L1 block with the given indexed hashes.
//// Order of the returned sidecars is guaranteed to be that of the hashes.
//// Blob data is not checked for validity.
//func (cl *L1BeaconClient) GetBlobSidecars(ctx context.Context, b *types.Header, hashes []IndexedBlobHash) ([]*BlobSidecar, error) {
//	if len(hashes) == 0 {
//		return []*BlobSidecar{}, nil
//	}
//	slotFn, err := cl.GetTimeToSlotFn(ctx)
//	if err != nil {
//		return nil, fmt.Errorf("failed to get time to slot function: %w", err)
//	}
//	slot, err := slotFn(b.Time)
//	if err != nil {
//		return nil, fmt.Errorf("error in converting ref.Time to slot: %w", err)
//	}
//
//	resp, err := cl.fetchSidecars(ctx, slot, hashes)
//	if err != nil {
//		return nil, fmt.Errorf("failed to fetch blob sidecars for slot %v block %v: %w", slot, b, err)
//	}
//
//	apiSidecar := make([]*APIBlobSidecar, 0, len(hashes))
//	// filter and order by versioned hashes
//	for _, h := range hashes {
//		for _, apisc := range resp.Data {
//			commitment := kzg4844.Commitment(apisc.KZGCommitment)
//			versionedHash := KZGToVersionedHash(commitment)
//
//			if h.Hash == versionedHash {
//				apiSidecar = append(apiSidecar, apisc)
//				break
//			}
//		}
//	}
//
//	if len(hashes) != len(apiSidecar) {
//		fmt.Printf("expected %v sidecars but got %v", len(hashes), len(apiSidecar))
//		return nil, fmt.Errorf("expected %v sidecars but got %v", len(hashes), len(apiSidecar))
//	}
//
//	blobSidecars := make([]*BlobSidecar, 0, len(hashes))
//	for _, apisc := range apiSidecar {
//		blobSidecars = append(blobSidecars, apisc.BlobSidecar())
//	}
//
//	return blobSidecars, nil
//}

// GetBlobSidecars fetches blob sidecars that were confirmed in the specified
// L1 block with the given indexed hashes.
// Order of the returned sidecars is guaranteed to be that of the hashes.
// Blob data is not checked for validity.
func (cl *L1BeaconClient) GetBlobSidecars(ctx context.Context, b *types.Header, hashes []IndexedBlobHash) ([]*kzg4844.Blob, error) {
	return testGetBlobSidecars(cl, ctx, b, hashes)
}

// GetBlobSidecars fetches blob sidecars that were confirmed in the specified
// L1 block with the given indexed hashes.
// Order of the returned sidecars is guaranteed to be that of the hashes.
// Blob data is not checked for validity.
func testGetBlobSidecars(cl *L1BeaconClient, ctx context.Context, b *types.Header, hashes []IndexedBlobHash) ([]*kzg4844.Blob, error) {
	if len(hashes) == 0 {
		return []*kzg4844.Blob{}, nil
	}
	slotFn, err := cl.GetTimeToSlotFn(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get time to slot function: %w", err)
	}
	slot, err := slotFn(b.Time)
	if err != nil {
		return nil, fmt.Errorf("error in converting ref.Time to slot: %w", err)
	}

	resp, err := cl.fetchSidecars(ctx, slot, hashes)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch blob sidecars for slot %v block %v: %w", slot, b, err)
	}

	response := resp.Data
	if len(hashes) != len(response) {
		return nil, fmt.Errorf("expected %v sidecars but got %v", len(hashes), len(response))
	}

	blobs := make([]*kzg4844.Blob, 0, len(hashes))
	for _, sidecars := range response {
		blobPtr := sidecars.Blob
		blobs = append(blobs, &blobPtr)
	}

	return blobs, nil
}

// FetchBlobs fetches blobs that were confirmed in the specified L1 block with the given indexed
// hashes. The order of the returned blobs will match the order of `hashes`.  Confirms each
// blob's validity by checking its proof against the commitment, and confirming the commitment
// hashes to the expected value. Returns error if any blob is found invalid.
func (cl *L1BeaconClient) FetchBlobs(ctx context.Context, b *types.Header, hashes []IndexedBlobHash) ([]*kzg4844.Blob, error) {
	blobs, err := cl.GetBlobSidecars(ctx, b, hashes)
	if err != nil {
		return nil, fmt.Errorf("failed to get blob sidecars for Block Header %s: %w", b.Hash().Hex(), err)
	}
	return blobs, nil
}

func blobsFromSidecars(blobSidecars []*BlobSidecar, hashes []IndexedBlobHash) ([]*kzg4844.Blob, error) {
	if len(blobSidecars) != len(hashes) {
		return nil, fmt.Errorf("number of hashes and blobSidecars mismatch, %d != %d", len(hashes), len(blobSidecars))
	}

	out := make([]*kzg4844.Blob, len(hashes))
	matchedSidecars := make([]bool, len(hashes))

	for i, ih := range hashes {
		var matchedSidecar *BlobSidecar
		var found bool

		for _, sidecar := range blobSidecars {
			if uint64(sidecar.Index) == ih.Index {
				// Ensure the blob's KZG commitment hashes to the expected value
				hash := KZGToVersionedHash(kzg4844.Commitment(sidecar.KZGCommitment))
				if hash == ih.Hash {
					matchedSidecar = sidecar
					found = true
					break
				}
			}
		}

		if !found {
			return nil, fmt.Errorf("expected sidecar for hash %s at index %d but did not find it", ih.Hash.Hex(), ih.Index)
		}

		println("SUCCESSFUL blob hash to commitment hash comparison: ", ih.Hash.Hex())

		// Confirm blob data is valid by verifying its proof against the commitment
		if err := VerifyBlobProof(&matchedSidecar.Blob, kzg4844.Commitment(matchedSidecar.KZGCommitment), kzg4844.Proof(matchedSidecar.KZGProof)); err != nil {
			return nil, fmt.Errorf("blob at index %d failed verification: %w", i, err)
		}
		out[i] = &matchedSidecar.Blob
		matchedSidecars[i] = true
	}

	// Ensure all sidecars have been matched
	for i, matched := range matchedSidecars {
		if !matched {
			return nil, fmt.Errorf("missing blob sidecar for hash %s at index %d", hashes[i].Hash.Hex(), hashes[i].Index)
		}
	}

	return out, nil
}

// GetVersion fetches the version of the Beacon-node.
func (cl *L1BeaconClient) GetVersion(ctx context.Context) (string, error) {
	return cl.cl.NodeVersion(ctx)
}
