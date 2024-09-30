package ethadapter

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

const (
	versionedHashPrefix = "/v1/blob/"
)

type ArchivalBlobResponse struct {
	Blob struct {
		VersionedHash string `json:"versionedHash"`
		Commitment    string `json:"commitment"`
		Proof         string `json:"proof"`
		Data          string `json:"data"`
	} `json:"blob"`
}

type ArchivalHTTPClient struct {
	client  *http.Client
	baseURL string
}

func NewArchivalHTTPClient(client *http.Client, baseURL string) *ArchivalHTTPClient {
	return &ArchivalHTTPClient{client: client, baseURL: baseURL}
}

func (ac *ArchivalHTTPClient) BeaconBlobSidecars(ctx context.Context, _ uint64, hashes []gethcommon.Hash) (APIGetBlobSidecarsResponse, error) {
	var resp APIGetBlobSidecarsResponse
	resp.Data = make([]*APIBlobSidecar, 0, len(hashes))

	for i, hash := range hashes {
		var archivalResp ArchivalBlobResponse
		reqPath := path.Join(versionedHashPrefix, hash.Hex())
		err := ac.request(ctx, &archivalResp, reqPath)
		if err != nil {
			return APIGetBlobSidecarsResponse{}, fmt.Errorf("failed to fetch blob for hash %s: %w", hash.Hex(), err)
		}

		blobSidecar, err := convertToAPIBlobSidecar(&archivalResp, i)
		if err != nil {
			return APIGetBlobSidecarsResponse{}, fmt.Errorf("failed to convert blob for hash %s: %w", hash.Hex(), err)
		}

		resp.Data = append(resp.Data, blobSidecar)
	}

	return resp, nil
}

func (ac *ArchivalHTTPClient) request(ctx context.Context, dest any, reqPath string) error {
	base := ac.baseURL
	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		base = "http://" + base
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return fmt.Errorf("failed to parse base URL: %w", err)
	}

	reqURL, err := baseURL.Parse(reqPath)
	if err != nil {
		return fmt.Errorf("failed to parse request path: %w", err)
	}

	headers := http.Header{}
	headers.Add("Accept", "application/json")

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header = headers

	resp, err := ac.client.Do(req)
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

func convertToAPIBlobSidecar(archivalResp *ArchivalBlobResponse, index int) (*APIBlobSidecar, error) {
	blobData, err := hex.DecodeString(strings.TrimPrefix(archivalResp.Blob.Data, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to decode blob data: %w", err)
	}

	var blob kzg4844.Blob
	copy(blob[:], blobData)

	commitment, err := hexToBytes48(archivalResp.Blob.Commitment)
	if err != nil {
		return nil, fmt.Errorf("failed to decode commitment: %w", err)
	}

	proof, err := hexToBytes48(archivalResp.Blob.Proof)
	if err != nil {
		return nil, fmt.Errorf("failed to decode proof: %w", err)
	}

	return &APIBlobSidecar{
		Blob:          blob,
		KZGCommitment: commitment,
		KZGProof:      proof,
		Index:         Uint64String(index),
	}, nil
}

func hexToBytes48(s string) (Bytes48, error) {
	b, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	if err != nil {
		return Bytes48{}, err
	}
	if len(b) != 48 {
		return Bytes48{}, fmt.Errorf("expected 48 bytes, got %d", len(b))
	}
	var result Bytes48
	copy(result[:], b)
	return result, nil
}
