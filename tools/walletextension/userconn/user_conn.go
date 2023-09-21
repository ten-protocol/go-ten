package userconn

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/tools/walletextension/common"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/gorilla/websocket"
)

const (
	httpCodeErr = 500
)

var upgrader = websocket.Upgrader{} // Used to upgrade connections to websocket connections.

// UserConn represents a connection to a user.
type UserConn interface {
	ReadRequest() ([]byte, error)
	ReadRequestParams() map[string]string
	WriteResponse(msg []byte) error
	HandleError(msg string)
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
		err = fmt.Errorf("unable to upgrade to websocket connection")
		logger.Error("unable to upgrade to websocket connection")
		httpLogAndSendErr(resp, err.Error()) // todo (@ziga) - Handle error properly for websockets.
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
		wrappedErr := fmt.Errorf("could not read request body: %w", err)
		h.HandleError(wrappedErr.Error())
		return nil, wrappedErr
	}
	return body, nil
}

func (h *userConnHTTP) WriteResponse(msg []byte) error {
	_, err := h.resp.Write(msg)
	if err != nil {
		wrappedErr := fmt.Errorf("could not write response: %w", err)
		h.HandleError(wrappedErr.Error())
		return wrappedErr
	}
	return nil
}

func (h *userConnHTTP) HandleError(msg string) {
	h.logger.Error(msg)
	httpLogAndSendErr(h.resp, msg)
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
		wrappedErr := fmt.Errorf("could not read request: %w", err)
		w.HandleError(wrappedErr.Error())
		return nil, wrappedErr
	}
	return msg, nil
}

func (w *userConnWS) WriteResponse(msg []byte) error {
	err := w.conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		if websocket.IsCloseError(err) {
			w.isClosed = true
		}
		wrappedErr := fmt.Errorf("could not write response: %w", err)
		w.HandleError(wrappedErr.Error())
		return wrappedErr
	}
	return nil
}

// HandleError logs and prints the error, and writes it to the websocket as a JSON object with a single key, "error".
func (w *userConnWS) HandleError(msg string) {
	w.logger.Error(msg)

	errMsg, err := json.Marshal(map[string]interface{}{
		common.JSONKeyErr: msg,
	})
	if err != nil {
		w.logger.Error("could not marshal websocket error message to JSON", log.ErrKey, err)
		return
	}

	err = w.conn.WriteMessage(websocket.TextMessage, errMsg)
	if err != nil {
		if websocket.IsCloseError(err) || strings.Contains(msg, "EOF") {
			w.isClosed = true
		}
		w.logger.Error("could not write error message to websocket", log.ErrKey, err)
		return
	}
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

// Logs the error, prints it to the console, and returns the error over HTTP.
func httpLogAndSendErr(resp http.ResponseWriter, msg string) {
	http.Error(resp, msg, httpCodeErr)
}

func getQueryParams(query url.Values) map[string]string {
	params := make(map[string]string)
	queryParams := query
	for key, value := range queryParams {
		params[key] = value[0]
	}
	return params
}
