package obscuroscan

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
)

const (
	pathGetBlockHeadHeight = "/blockheadheight/"
	staticDir              = "./tools/obscuroscan/static"
	pathRoot               = "/"
	httpCodeErr            = 500
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
	serveMux.HandleFunc(pathGetBlockHeadHeight, o.getBlockHeadHeight)
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

// Retrieves the current block head height for the Obscuro network.
func (o *Obscuroscan) getBlockHeadHeight(resp http.ResponseWriter, req *http.Request) {
	var currentBlockHeadHeight uint64
	err := (*o.client).Call(&currentBlockHeadHeight, obscuroclient.RPCGetCurrentBlockHeadHeight)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not retrieve current block head height: %s", err))
		return
	}

	_, err = resp.Write([]byte(strconv.FormatUint(currentBlockHeadHeight, 10)))
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not return current block head height to client: %s", err))
		return
	}
}

// Logs the error message and sends it as an HTTP error.
func logAndSendErr(resp http.ResponseWriter, msg string) {
	fmt.Println(msg)
	http.Error(resp, msg, httpCodeErr)
}
