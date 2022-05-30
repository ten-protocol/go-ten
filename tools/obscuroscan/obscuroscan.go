package obscuroscan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
)

const (
	pathHeadBlock     = "/headblock/"
	pathHeadRollup    = "/headrollup/"
	pathDecryptRollup = "/decryptrollup/"
	staticDir         = "./tools/obscuroscan/static"
	pathRoot          = "/"
	httpCodeErr       = 500
	methodBytesLen    = 4
)

// Obscuroscan is a server that allows the monitoring of a running Obscuro network.
type Obscuroscan struct {
	server      *http.Server
	client      *obscuroclient.Client
	contractABI abi.ABI
}

func NewObscuroscan(address string) *Obscuroscan {
	client := obscuroclient.NewClient(address)
	contractABI, err := abi.JSON(strings.NewReader(mgmtcontractlib.MgmtContractABI))
	if err != nil {
		panic("could not parse management contract ABI to decrypt rollups")
	}
	return &Obscuroscan{
		client:      &client,
		contractABI: contractABI,
	}
}

// Serve listens for and serves Obscuroscan requests.
func (o *Obscuroscan) Serve(hostAndPort string) {
	serveMux := http.NewServeMux()
	// Serves the web interface.
	serveMux.Handle(pathRoot, http.FileServer(http.Dir(staticDir)))
	// Handle requests for block head height.
	serveMux.HandleFunc(pathHeadBlock, o.getBlockHead)
	// Handle requests for the head rollup.
	serveMux.HandleFunc(pathHeadRollup, o.getHeadRollup)
	// Handle requests to decrypt rollup.
	serveMux.HandleFunc(pathDecryptRollup, o.decryptRollup)
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

	jsonBlock, err := json.Marshal(headBlock)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not return head block to client: %s", err))
		return
	}
	_, err = resp.Write(jsonBlock)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not return head block to client: %s", err))
		return
	}
}

// Retrieves the head rollup for the Obscuro network.
func (o *Obscuroscan) getHeadRollup(resp http.ResponseWriter, _ *http.Request) {
	// TODO - Update logic here once rollups are encrypted.
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

// Decrypts the provided rollup using the provided key.
func (o *Obscuroscan) decryptRollup(resp http.ResponseWriter, req *http.Request) {
	// TODO - Update logic here once rollups are encrypted.
	body := req.Body
	defer body.Close()
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(body)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not read request body: %s", err))
		return
	}
	encryptedRollupHex := buffer.Bytes()
	encryptedRollupBytes := common.Hex2Bytes(string(encryptedRollupHex))

	method, err := o.contractABI.MethodById(encryptedRollupBytes[:methodBytesLen])
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not read ABI method for encrypted rollup. Cause: %s", err))
		return
	}
	if method.Name != mgmtcontractlib.AddRollupMethod {
		logAndSendErr(resp, fmt.Sprintf("encrypted rollup did not have correct ABI method name. Expected %s, got %s", mgmtcontractlib.AddRollupMethod, method.Name))
		return
	}

	contractCallData := map[string]interface{}{}
	if err = method.Inputs.UnpackIntoMap(contractCallData, encryptedRollupBytes[4:]); err != nil {
		logAndSendErr(resp, fmt.Sprintf("encrypted rollup could not be unpacked using ABI. Cause: %s", err))
		return
	}
	callData, found := contractCallData["rollupData"]
	if !found {
		logAndSendErr(resp, "encrypted rollup did not contain call data for rollupData")
		return
	}
	zippedRollup := mgmtcontractlib.Base64DecodeFromString(callData.(string))
	encodedRollup, err := mgmtcontractlib.Decompress(zippedRollup)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("decrypted rollup could not be decompressed. Cause: %s", err))
		return
	}
	cleartextRollup, err := nodecommon.DecodeRollup(encodedRollup)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not decode decompressed rollup. Cause: %s", err))
		return
	}

	jsonRollup, err := json.Marshal(cleartextRollup)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not decrypt rollup: %s", err))
		return
	}
	_, err = resp.Write(jsonRollup)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not decrypt rollup: %s", err))
		return
	}
}

// Logs the error message and sends it as an HTTP error.
func logAndSendErr(resp http.ResponseWriter, msg string) {
	fmt.Println(msg)
	http.Error(resp, msg, httpCodeErr)
}
