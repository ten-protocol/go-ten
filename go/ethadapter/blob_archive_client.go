package ethadapter

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"path"
	"strings"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

const (
	versionedHashPrefix = "/v1/blob/"
)

type ArchivalBlobResponse struct {
	VersionedHash string `json:"versionedHash"`
	Commitment    string `json:"commitment"`
	Proof         string `json:"proof"`
	Data          string `json:"data"`
}

type ArchivalHTTPClient struct {
	httpClient *BaseHTTPClient
}

func NewArchivalHTTPClient(client *http.Client, baseURL string) *ArchivalHTTPClient {
	return &ArchivalHTTPClient{
		httpClient: NewBaseHTTPClient(client, baseURL),
	}
}

func (ac *ArchivalHTTPClient) BeaconBlobSidecars(ctx context.Context, _ uint64, hashes []gethcommon.Hash) (APIGetBlobSidecarsResponse, error) {
	var resp APIGetBlobSidecarsResponse
	resp.Data = make([]*BlobSidecar, 0, len(hashes))

	for i, hash := range hashes {
		var archivalResp ArchivalBlobResponse
		reqPath := path.Join(versionedHashPrefix, hash.Hex())
		err := ac.request(ctx, &archivalResp, reqPath)
		if err != nil {
			return APIGetBlobSidecarsResponse{}, fmt.Errorf("failed to fetch blob for hash %s: %w", hash.Hex(), err)
		}

		blobSidecar, err := convertToSidecar(&archivalResp, i)
		if err != nil {
			return APIGetBlobSidecarsResponse{}, fmt.Errorf("failed to convert blob for hash %s: %w", hash.Hex(), err)
		}

		resp.Data = append(resp.Data, blobSidecar)
	}

	return resp, nil
}

func (ac *ArchivalHTTPClient) request(ctx context.Context, dest any, reqPath string) error {
	return ac.httpClient.Request(ctx, dest, reqPath, nil)
}

func convertToSidecar(archivalResp *ArchivalBlobResponse, index int) (*BlobSidecar, error) {
	blobData, err := hex.DecodeString(strings.TrimPrefix(archivalResp.Data, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to decode blob data: %w", err)
	}

	var blob kzg4844.Blob
	copy(blob[:], blobData)

	commitment, err := hexToBytes48(archivalResp.Commitment)
	if err != nil {
		return nil, fmt.Errorf("failed to decode commitment: %w", err)
	}

	proof, err := hexToBytes48(archivalResp.Proof)
	if err != nil {
		return nil, fmt.Errorf("failed to decode proof: %w", err)
	}

	return &BlobSidecar{
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
