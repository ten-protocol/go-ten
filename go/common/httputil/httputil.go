package httputil

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/valyala/fasthttp"
)

const (
	// CORS-related constants.
	CorsAllowOrigin      = "Access-Control-Allow-Origin"
	CorsAllowCredentials = "Access-Control-Allow-Credentials"
	OriginAll            = "*"
	CorsAllowMethods     = "Access-Control-Allow-Methods"
	ReqOptions           = "OPTIONS"
	CorsAllowHeaders     = "Access-Control-Allow-Headers"
	CorsHeaders          = "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"
)

// CreateTLSHTTPClient provides a basic http client prepared with a trusted CA cert
func CreateTLSHTTPClient(caCertPEM string) (*http.Client, error) {
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM([]byte(caCertPEM)); !ok {
		return nil, fmt.Errorf("failed to append to CA cert from caCertPEM=%s", caCertPEM)
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:    caCertPool,
				MinVersion: tls.VersionTLS12,
			},
		},
	}, nil
}

// ExecuteHTTPReq executes an HTTP request:
// * returns an error if request fails or if the response code was outside the range 200-299
// * returns response body as bytes if there was a response body
func ExecuteHTTPReq(client *http.Client, req *http.Request) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed url=%s - %w", req.URL.String(), err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var msg []byte
		_, err := resp.Body.Read(msg)
		if err != nil {
			return nil, fmt.Errorf("req failed url=%s, statusCode=%d, failed to read status text", req.URL.String(), resp.StatusCode)
		}
		return nil, fmt.Errorf("req failed url=%s status: %d %s", req.URL.String(), resp.StatusCode, msg)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// success status code but no body, ignoring error and returning no bytes
		return []byte{}, nil //nolint:nilerr
	}
	return body, nil
}

// EnableCORS Allows Tenscan and WalletExtension APIs to serve other web apps via CORS.
func EnableCORS(resp http.ResponseWriter, req *http.Request) bool {
	origin := req.Header.Get("Origin")
	
	// For cookie-based authentication, we need to allow specific origins, not "*"
	// Allow requests from .ten.xyz domains and localhost for development
	if origin != "" && (strings.HasSuffix(origin, ".ten.xyz") || origin == "https://ten.xyz" || strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1")) {
		resp.Header().Set(CorsAllowOrigin, origin)
		resp.Header().Set(CorsAllowCredentials, "true")
	} else {
		// Fallback for non-cookie requests
		resp.Header().Set(CorsAllowOrigin, OriginAll)
	}
	
	if (*req).Method == ReqOptions {
		// Returns true if the request was a pre-flight, e.g. OPTIONS, to stop further processing.
		resp.Header().Set(CorsAllowMethods, ReqOptions)
		resp.Header().Set(CorsAllowHeaders, CorsHeaders)
		return true
	}
	return false
}

// PostDataJSON sends a JSON payload to the given URL and returns the status, response body, and any error encountered.
func PostDataJSON(url string, data []byte) (int, []byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetBody(data)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := fasthttp.Do(req, resp); err != nil {
		return 0, nil, fmt.Errorf("error while sending request: %w", err)
	}

	return resp.StatusCode(), resp.Body(), nil
}
