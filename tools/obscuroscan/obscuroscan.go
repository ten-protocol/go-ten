package obscuroscan

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"strconv"
	"strings"

	"github.com/obscuronet/go-obscuro/go/enclave/crypto"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/go-obscuro/go/common"

	"github.com/obscuronet/go-obscuro/go/rpcclientlib"
)

const (
	pathNumRollups    = "/numrollups/"
	pathBlock         = "/block/"
	pathRollup        = "/rollup/"
	pathDecryptTxBlob = "/decrypttxblob/"
	staticDir         = "static"
	pathRoot          = "/"
	httpCodeErr       = 500
)

//go:embed static
var staticFiles embed.FS

// Obscuroscan is a server that allows the monitoring of a running Obscuro network.
type Obscuroscan struct {
	server      *http.Server
	client      rpcclientlib.Client
	contractABI abi.ABI
}

func NewObscuroscan(address string) *Obscuroscan {
	client := rpcclientlib.NewClient(address)
	contractABI, err := abi.JSON(strings.NewReader(mgmtcontractlib.MgmtContractABI))
	if err != nil {
		panic("could not parse management contract ABI to decrypt rollups")
	}
	return &Obscuroscan{
		client:      client,
		contractABI: contractABI,
	}
}

// Serve listens for and serves Obscuroscan requests.
func (o *Obscuroscan) Serve(hostAndPort string) {
	serveMux := http.NewServeMux()

	serveMux.HandleFunc(pathNumRollups, o.getNumRollups)    // Get the number of published rollups.
	serveMux.HandleFunc(pathBlock, o.getBlock)              // Get the L1 block with the given number.
	serveMux.HandleFunc(pathRollup, o.getRollup)            // Get the rollup with the given number.
	serveMux.HandleFunc(pathDecryptTxBlob, o.decryptTxBlob) // Decrypt a transaction blob.

	// Serves the web assets for the user interface.
	noPrefixStaticFiles, err := fs.Sub(staticFiles, staticDir)
	if err != nil {
		panic(fmt.Sprintf("could not serve static files. Cause: %s", err))
	}
	serveMux.Handle(pathRoot, http.FileServer(http.FS(noPrefixStaticFiles)))

	o.server = &http.Server{Addr: hostAndPort, Handler: serveMux}

	err = o.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (o *Obscuroscan) Shutdown() {
	if o.server != nil {
		err := o.server.Shutdown(context.Background())
		if err != nil {
			fmt.Printf("could not shut down Obscuroscan. Cause: %s", err)
		}
	}
}

// Retrieves the number of published rollups.
func (o *Obscuroscan) getNumRollups(resp http.ResponseWriter, _ *http.Request) {
	var rollupHeader *common.Header
	err := o.client.Call(&rollupHeader, rpcclientlib.RPCGetCurrentRollupHead)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not retrieve head rollup. Cause: %s", err))
		return
	}

	numOfRollups := rollupHeader.Number.Int64()
	numOfRollupsStr := strconv.Itoa(int(numOfRollups))
	_, err = resp.Write([]byte(numOfRollupsStr))
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not return number of rollups to client. Cause: %s", err))
		return
	}
}

// Retrieves the L1 block header with the given number.
func (o *Obscuroscan) getBlock(resp http.ResponseWriter, req *http.Request) {
	body := req.Body
	defer body.Close()
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(body)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not read request body. Cause: %s", err))
		return
	}
	blockHashStr := buffer.String()
	blockHash := gethcommon.HexToHash(blockHashStr)

	var blockHeader *types.Header
	err = o.client.Call(&blockHeader, rpcclientlib.RPCGetBlockHeaderByHash, blockHash)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not retrieve block with hash %s. Cause: %s", blockHash, err))
		return
	}

	jsonBlock, err := json.Marshal(blockHeader)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not return block to client. Cause: %s", err))
		return
	}
	_, err = resp.Write(jsonBlock)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not return block to client. Cause: %s", err))
		return
	}
}

// Retrieves the rollup with the given number.
func (o *Obscuroscan) getRollup(resp http.ResponseWriter, req *http.Request) {
	body := req.Body
	defer body.Close()
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(body)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not read request body. Cause: %s", err))
		return
	}
	bufferStr := buffer.String()
	number, err := strconv.Atoi(bufferStr)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not parse \"%s\" as an integer", bufferStr))
		return
	}

	// TODO - If required, consolidate the two calls below into a single RPCGetRollupByNumber call to minimise round trips.
	var rollupHeader *common.Header
	err = o.client.Call(&rollupHeader, rpcclientlib.RPCGetRollupHeaderByNumber, number)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not retrieve rollup with number %d. Cause: %s", number, err))
		return
	}

	rollupHash := rollupHeader.Hash()
	if rollupHash == (gethcommon.Hash{}) {
		logAndSendErr(resp, fmt.Sprintf("rollup was retrieved but hash was nil"))
		return
	}

	var rollup *common.ExtRollup
	err = o.client.Call(&rollup, rpcclientlib.RPCGetRollup, rollupHash)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not retrieve rollup. Cause: %s", err))
		return
	}

	jsonRollup, err := json.Marshal(rollup)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not return rollup to client. Cause: %s", err))
		return
	}
	_, err = resp.Write(jsonRollup)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not return rollup to client. Cause: %s", err))
		return
	}
}

// Decrypts the provided transaction blob using the provided key.
// TODO - Use the passed-in key, rather than a hardcoded enclave key.
func (o *Obscuroscan) decryptTxBlob(resp http.ResponseWriter, req *http.Request) {
	body := req.Body
	defer body.Close()
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(body)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not read request body. Cause: %s", err))
		return
	}

	jsonTxs, err := decryptTxBlob(buffer.Bytes())
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not decrypt transaction blob. Cause: %s", err))
		return
	}

	_, err = resp.Write(jsonTxs)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not write decrypted transactions to client. Cause: %s", err))
		return
	}
}

// Decrypts the transaction blob and returns it as JSON.
func decryptTxBlob(encryptedTxBytesBase64 []byte) ([]byte, error) {
	encryptedTxBytes, err := base64.StdEncoding.DecodeString(string(encryptedTxBytesBase64))
	if err != nil {
		return nil, fmt.Errorf("could not decode encrypted transaction blob from Base64. Cause: %w", err)
	}

	key := gethcommon.Hex2Bytes(crypto.RollupEncryptionKeyHex)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not initialise AES cipher for enclave rollup key. Cause: %w", err)
	}
	transactionCipher, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("could not initialise wrapper for AES cipher for enclave rollup key. Cause: %w", err)
	}

	encodedTxs, err := transactionCipher.Open(nil, []byte(crypto.RollupCipherNonce), encryptedTxBytes, nil)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt encrypted L2 transactions. Cause: %w", err)
	}

	var cleartextTxs []*common.L2Tx
	if err = rlp.DecodeBytes(encodedTxs, &cleartextTxs); err != nil {
		return nil, fmt.Errorf("could not decode encoded L2 transactions. Cause: %w", err)
	}

	jsonRollup, err := json.Marshal(cleartextTxs)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt transaction blob. Cause: %w", err)
	}

	return jsonRollup, nil
}

// Logs the error message and sends it as an HTTP error.
func logAndSendErr(resp http.ResponseWriter, msg string) {
	fmt.Println(msg)
	http.Error(resp, msg, httpCodeErr)
}
