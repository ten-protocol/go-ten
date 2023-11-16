package userconn

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ten-protocol/go-ten/go/common/log"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // Used to upgrade connections to websocket connections.

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

// Represents a user's connection websockets.
type userConnWS struct {
	conn     *websocket.Conn
	isClosed bool
	logger   gethlog.Logger
	req      *http.Request
}

func NewUserConnHTTP(resp http.ResponseWriter, req *http.Request, logger gethlog.Logger) UserConn {
	return &userConnHTTP{resp: resp, req: req, logger: logger}
}

func NewUserConnWS(resp http.ResponseWriter, req *http.Request, logger gethlog.Logger) (UserConn, error) {
	// We search all the request's headers. If there's a websocket upgrade header, we upgrade to a websocket connection.
	conn, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		err = fmt.Errorf("unable to upgrade to websocket connection - %w", err)
		_, _ = resp.Write([]byte(err.Error()))
		logger.Error("unable to upgrade to websocket connection", log.ErrKey, err)
		return nil, err
	}

	return &userConnWS{
		conn:   conn,
		logger: logger,
		req:    req,
	}, nil
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

func (w *userConnWS) ReadRequest() ([]byte, error) {
	_, msg, err := w.conn.ReadMessage()
	if err != nil {
		if websocket.IsCloseError(err) {
			w.isClosed = true
		}
		return nil, fmt.Errorf("could not read request: %w", err)
	}
	return msg, nil
}

func (w *userConnWS) WriteResponse(msg []byte) error {
	err := w.conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		if websocket.IsCloseError(err) || strings.Contains(string(msg), "EOF") {
			w.isClosed = true
		}
		return fmt.Errorf("could not write response: %w", err)
	}
	return nil
}

func (w *userConnWS) SupportsSubscriptions() bool {
	return true
}

func (w *userConnWS) IsClosed() bool {
	return w.isClosed
}

func (w *userConnWS) ReadRequestParams() map[string]string {
	return getQueryParams(w.req.URL.Query())
}

func (w *userConnWS) GetHTTPRequest() *http.Request {
	return w.req
}

func getQueryParams(query url.Values) map[string]string {
	params := make(map[string]string)
	queryParams := query
	for key, value := range queryParams {
		params[key] = value[0]
	}
	return params
}
