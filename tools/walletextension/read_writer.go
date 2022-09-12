package walletextension

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
)

// ReadWriter handles reading and writing Ethereum JSON RPC requests.
type ReadWriter interface {
	ReadRequest() ([]byte, error)
	WriteResponse([]byte) error
	SupportsSubscriptions() bool
}

// HTTPReadWriter is a ReadWriter over HTTP.
type HTTPReadWriter struct {
	resp http.ResponseWriter
	req  *http.Request
}

// WSReadWriter is a ReadWriter over websockets.
type WSReadWriter struct {
	conn *websocket.Conn
}

func NewReadWriter(resp http.ResponseWriter, req *http.Request) (ReadWriter, error) {
	for _, header := range req.Header["Upgrade"] { // todo - joel - use constant
		if header == "websocket" { // todo - joel - use constant
			conn, err := upgrader.Upgrade(resp, req, nil)
			if err != nil {
				return nil, fmt.Errorf("attempted to subscribe, but was unable to create websocket connection")
			}
			return &WSReadWriter{
				conn: conn,
			}, nil
		}
	}

	return &HTTPReadWriter{resp: resp, req: req}, nil
}

func (h HTTPReadWriter) ReadRequest() ([]byte, error) {
	body, err := io.ReadAll(h.req.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read JSON-RPC request body: %w", err)
	}
	return body, nil
}

func (h HTTPReadWriter) WriteResponse(responseBytes []byte) error {
	_, err := h.resp.Write(responseBytes)
	if err != nil {
		return fmt.Errorf("could not write JSON-RPC response: %w", err)
	}
	return nil
}

func (h HTTPReadWriter) SupportsSubscriptions() bool {
	return false
}

func (w WSReadWriter) ReadRequest() ([]byte, error) {
	panic("todo - joel")
}

func (w WSReadWriter) WriteResponse(responseBytes []byte) error {
	panic("todo - joel")
}

func (w WSReadWriter) SupportsSubscriptions() bool {
	return true
}
