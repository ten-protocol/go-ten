package obscuroscan

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
)

const (
	pathBlockHead  = "/blockhead/"
	pathHeadRollup = "/headrollup/"
	staticDir      = "./tools/obscuroscan/static"
	pathRoot       = "/"
	httpCodeErr    = 500
)

// Obscuroscan is a server that allows the monitoring of a running Obscuro network.
type Obscuroscan struct {
	server *http.Server
	client *obscuroclient.Client
}

func NewObscuroscan(nodeID common.Address, address string) *Obscuroscan {
	client := obscuroclient.NewClient(nodeID, address)
	return &Obscuroscan{
		client: &client,
	}
}

// Serve listens for and serves Obscuroscan requests.
func (o *Obscuroscan) Serve(hostAndPort string) {
	serveMux := http.NewServeMux()
	// Serves the web interface.
	serveMux.Handle(pathRoot, http.FileServer(http.Dir(staticDir)))
	// Handle requests for block head height.
	serveMux.HandleFunc(pathBlockHead, o.getBlockHead)
	// Handle requests for the head rollup.
	serveMux.HandleFunc(pathHeadRollup, o.getHeadRollup)
	o.server = &http.Server{Addr: hostAndPort, Handler: serveMux}

	err := o.server.ListenAndServe()
	if err != http.ErrServerClosed {
		panic(err)
	}
}

func (o *Obscuroscan) Shutdown() {
	if o.server != nil {
		err := o.server.Shutdown(context.Background())
		if err != nil {
			fmt.Printf("could not shut down Obscuroscan: %s", err)
		}
	}
}

// Retrieves the current block header for the Obscuro network.
func (o *Obscuroscan) getBlockHead(resp http.ResponseWriter, _ *http.Request) {
	var headBlock *types.Header
	err := (*o.client).Call(&headBlock, obscuroclient.RPCGetCurrentBlockHead)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not retrieve head block: %s", err))
		return
	}

	jsonRollup, err := json.Marshal(headBlock)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not return head block to client: %s", err))
		return
	}
	_, err = resp.Write(jsonRollup)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not return head block to client: %s", err))
		return
	}
}

// Retrieves the head rollup for the Obscuro network.
func (o *Obscuroscan) getHeadRollup(resp http.ResponseWriter, _ *http.Request) {
	var headRollup *nodecommon.Header
	err := (*o.client).Call(&headRollup, obscuroclient.RPCGetCurrentRollupHead)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not retrieve head rollup: %s", err))
		return
	}

	jsonRollup, err := json.Marshal(headRollup)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not return head rollup to client: %s", err))
		return
	}
	_, err = resp.Write(jsonRollup)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not return head rollup to client: %s", err))
		return
	}
}

// Logs the error message and sends it as an HTTP error.
func logAndSendErr(resp http.ResponseWriter, msg string) {
	fmt.Println(msg)
	http.Error(resp, msg, httpCodeErr)
}
