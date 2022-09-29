package userconn

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
	WriteResponse(msg []byte) error
	HandleError(msg string)
	SupportsSubscriptions() bool
	IsClosed() bool
}

// Represents a user's connection over HTTP.
type userConnHTTP struct {
	resp http.ResponseWriter
	req  *http.Request
}

// Represents a user's connection websockets.
type userConnWS struct {
	conn     *websocket.Conn
	isClosed bool
}

func NewUserConnHTTP(resp http.ResponseWriter, req *http.Request) UserConn {
	return &userConnHTTP{resp: resp, req: req}
}

func NewUserConnWS(resp http.ResponseWriter, req *http.Request) (UserConn, error) {
	// We search all the request's headers. If there's a websocket upgrade header, we upgrade to a websocket connection.
	conn, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		err = fmt.Errorf("unable to upgrade to websocket connection")
		httpLogAndSendErr(resp, err.Error()) // TODO - Handle error properly for websockets.
		return nil, err
	}

	return &userConnWS{
		conn: conn,
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

func (h *userConnHTTP) HandleError(msg string) {
	httpLogAndSendErr(h.resp, msg)
}

func (h *userConnHTTP) SupportsSubscriptions() bool {
	return false
}

func (h *userConnHTTP) IsClosed() bool {
	return false
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
		if websocket.IsCloseError(err) {
			w.isClosed = true
		}
		return fmt.Errorf("could not write response: %w", err)
	}
	return nil
}

// HandleError logs and prints the error, and writes it to the websocket as a JSON object with a single key, "error".
func (w *userConnWS) HandleError(msg string) {
	log.Error(msg)
	fmt.Println(msg)

	errMsg, err := json.Marshal(map[string]interface{}{
		common.JSONKeyErr: msg,
	})
	if err != nil {
		log.Error("could not marshal websocket error message to JSON")
		return
	}

	err = w.conn.WriteMessage(websocket.TextMessage, errMsg)
	if err != nil {
		if websocket.IsCloseError(err) {
			w.isClosed = true
		}
		log.Error("could not write error message to websocket")
		return
	}
}

func (w *userConnWS) SupportsSubscriptions() bool {
	return true
}

func (w *userConnWS) IsClosed() bool {
	return w.isClosed
}

// Logs the error, prints it to the console, and returns the error over HTTP.
func httpLogAndSendErr(resp http.ResponseWriter, msg string) {
	log.Error(msg)
	fmt.Println(msg)
	http.Error(resp, msg, httpCodeErr)
}
