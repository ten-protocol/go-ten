package ethadapter

import (
	"context"
	"encoding/json"
	"fmt"
	gethlog "github.com/ethereum/go-ethereum/log"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ethereum/go-ethereum"
)

// BaseHTTPClient provides common HTTP functionality for different clients
type BaseHTTPClient struct {
	client  *http.Client
	logger  gethlog.Logger
	baseURL string
}

func NewBaseHTTPClient(client *http.Client, logger gethlog.Logger, baseURL string) *BaseHTTPClient {
	return &BaseHTTPClient{client: client, logger: logger, baseURL: baseURL}
}

func (chc *BaseHTTPClient) Request(ctx context.Context, dest any, reqPath string, reqQuery url.Values) error {
	base := chc.baseURL
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

	reqURL.RawQuery = reqQuery.Encode()

	chc.logger.Debug("Beacon client GET: %s", "url", reqURL.String())

	headers := http.Header{}
	headers.Add("Accept", "application/json")

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header = headers

	resp, err := chc.client.Do(req)
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
