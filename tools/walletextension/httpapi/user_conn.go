package httpapi

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	gethlog "github.com/ethereum/go-ethereum/log"
)

// UserConn represents a connection to a user.
type UserConn interface {
	ReadRequest() ([]byte, error)
	ReadRequestParams() map[string]string
	WriteResponse(msg []byte) error
	SupportsSubscriptions() bool
	IsClosed() bool
	GetHTTPRequest() *http.Request
}

// Represents a user's connection over HTTP.
type userConnHTTP struct {
	resp   http.ResponseWriter
	req    *http.Request
	logger gethlog.Logger
}

func NewUserConnHTTP(resp http.ResponseWriter, req *http.Request, logger gethlog.Logger) UserConn {
	return &userConnHTTP{resp: resp, req: req, logger: logger}
}

func (h *userConnHTTP) ReadRequest() ([]byte, error) {
	body, err := io.ReadAll(h.req.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read request body: %w", err)
	}
	return body, nil
}

func (h *userConnHTTP) WriteResponse(msg []byte) error {
	_, err := h.resp.Write(msg)
	if err != nil {
		return fmt.Errorf("could not write response: %w", err)
	}
	return nil
}

func (h *userConnHTTP) SupportsSubscriptions() bool {
	return false
}

func (h *userConnHTTP) IsClosed() bool {
	return false
}

func (h *userConnHTTP) ReadRequestParams() map[string]string {
	return getQueryParams(h.req.URL.Query())
}

func (h *userConnHTTP) GetHTTPRequest() *http.Request {
	return h.req
}

func getQueryParams(query url.Values) map[string]string {
	params := make(map[string]string)
	queryParams := query
	for key, value := range queryParams {
		params[key] = value[0]
	}
	return params
}
