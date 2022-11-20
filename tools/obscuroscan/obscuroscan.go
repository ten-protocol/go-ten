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
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/obscuronet/go-obscuro/go/obsclient"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/edgelesssys/ego/enclave"

	"github.com/obscuronet/go-obscuro/go/enclave/crypto"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/go-obscuro/go/common"

	"github.com/obscuronet/go-obscuro/go/rpc"
)

const (
	pathAPI               = "/api"
	pathNumRollups        = "/numrollups/"
	pathNumTxs            = "/numtxs/"
	pathGetRollupTime     = "/rolluptime/"
	pathLatestRollups     = "/latestrollups/"
	pathLatestTxs         = "/latesttxs/"
	pathBlock             = "/block/"
	pathRollup            = "/rollup/"
	pathDecryptTxBlob     = "/decrypttxblob/"
	pathAttestation       = "/attestation/"
	pathAttestationReport = "/attestationreport/"
	pathRoot              = "/"

	staticDir   = "static"
	extDivider  = "."
	extHTML     = ".html"
	httpCodeErr = 500
)

//go:embed static
var staticFiles embed.FS

// Obscuroscan is a server that allows the monitoring of a running Obscuro network.
type Obscuroscan struct {
	server      *http.Server
	client      rpc.Client
	obsClient   *obsclient.ObsClient
	contractABI abi.ABI
	logger      gethlog.Logger
}

// Identical to attestation.Report, but with the status mapped to a user-friendly string.
type attestationReportExternal struct {
	Data            []byte
	SecurityVersion uint
	Debug           bool
	UniqueID        []byte
	SignerID        []byte
	ProductID       []byte
	TCBStatus       string
}

func NewObscuroscan(address string, logger gethlog.Logger) *Obscuroscan {
	client, err := rpc.NewNetworkClient(address)
	if err != nil {
		panic(err)
	}
	obsClient := obsclient.NewObsClient(client)

	contractABI, err := abi.JSON(strings.NewReader(mgmtcontractlib.MgmtContractABI))
	if err != nil {
		panic("could not parse management contract ABI to decrypt rollups")
	}

	return &Obscuroscan{
		client:      client,
		obsClient:   obsClient,
		contractABI: contractABI,
		logger:      logger,
	}
}

// Serve listens for and serves Obscuroscan requests.
func (o *Obscuroscan) Serve(hostAndPort string) {
	serveMux := http.NewServeMux()

	serveMux.HandleFunc(pathAPI+pathNumRollups, o.getNumRollups)            // Get the number of published rollups.
	serveMux.HandleFunc(pathAPI+pathNumTxs, o.getNumTransactions)           // Get the number of rolled-up transactions.
	serveMux.HandleFunc(pathAPI+pathGetRollupTime, o.getRollupTime)         // Get the average rollup time.
	serveMux.HandleFunc(pathAPI+pathLatestRollups, o.getLatestRollups)      // Get the latest rollup numbers.
	serveMux.HandleFunc(pathAPI+pathLatestTxs, o.getLatestTxs)              // Get the latest transaction hashes.
	serveMux.HandleFunc(pathAPI+pathRollup, o.getRollupByNumOrTxHash)       // Get the rollup given its number or the hash of a transaction it contains.
	serveMux.HandleFunc(pathAPI+pathBlock, o.getBlock)                      // Get the L1 block with the given number.
	serveMux.HandleFunc(pathAPI+pathDecryptTxBlob, o.decryptTxBlob)         // Decrypt a transaction blob.
	serveMux.HandleFunc(pathAPI+pathAttestation, o.attestation)             // Retrieve the node's attestation.
	serveMux.HandleFunc(pathAPI+pathAttestationReport, o.attestationReport) // Retrieve the node's attestation report.

	// Serves the web assets for the user interface.
	staticFileFS, err := fs.Sub(staticFiles, staticDir)
	if err != nil {
		o.logger.Crit("could not serve static files.", log.ErrKey, err)
	}
	staticFileFilesystem := http.FileServer(http.FS(staticFileFS))
	serveMux.Handle(pathRoot, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If we get a request without an extension (other than for the root), we tack on ".html".
		if r.URL.Path != pathRoot && !strings.Contains(r.URL.Path, extDivider) {
			r.URL.Path += extHTML
		}
		staticFileFilesystem.ServeHTTP(w, r)
	}))

	o.server = &http.Server{Addr: hostAndPort, Handler: serveMux, ReadHeaderTimeout: 10 * time.Second}
	err = o.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (o *Obscuroscan) Shutdown() {
	if o.server != nil {
		err := o.server.Shutdown(context.Background())
		if err != nil {
			o.logger.Error("could not shut down Obscuroscan.", log.ErrKey, err)
		}
	}
}

// Retrieves the number of published rollups.
func (o *Obscuroscan) getNumRollups(resp http.ResponseWriter, _ *http.Request) {
	numOfRollups, err := o.obsClient.RollupNumber()
	if err != nil {
		o.logger.Error("Could not fetch number of rollups.", log.ErrKey, err)
		logAndSendErr(resp, "Could not fetch number of rollups.")
		return
	}

	numOfRollupsStr := strconv.Itoa(int(numOfRollups))
	_, err = resp.Write([]byte(numOfRollupsStr))
	if err != nil {
		o.logger.Error("could not return number of rollups to client.", log.ErrKey, err)
		logAndSendErr(resp, "Could not fetch number of rollups.")
		return
	}
}

// Retrieves the total number of transactions.
func (o *Obscuroscan) getNumTransactions(resp http.ResponseWriter, _ *http.Request) {
	var numTransactions *big.Int
	err := o.client.Call(&numTransactions, rpc.GetTotalTxs)
	if err != nil {
		o.logger.Error("Could not fetch total transactions.", log.ErrKey, err)
		logAndSendErr(resp, "Could not fetch total transactions.")
		return
	}

	_, err = resp.Write([]byte(numTransactions.String()))
	if err != nil {
		o.logger.Error("could not return total number of transactions to client.", log.ErrKey, err)
		logAndSendErr(resp, "Could not fetch total transactions.")
		return
	}
}

// Retrieves the average rollup time, as (time last rollup - time first rollup)/number of rollups
func (o *Obscuroscan) getRollupTime(resp http.ResponseWriter, _ *http.Request) {
	numLatestRollup, err := o.obsClient.RollupNumber()
	if err != nil {
		o.logger.Error("Could not fetch latest rollup number.", log.ErrKey, err)
		logAndSendErr(resp, "Could not fetch average rollup time.")
		return
	}

	firstRollupHeader, err := o.obsClient.RollupHeaderByNumber(big.NewInt(0))
	if err != nil || firstRollupHeader == nil {
		o.logger.Error("Could not fetch first rollup.", log.ErrKey, err)
		logAndSendErr(resp, "Could not fetch average rollup time.")
		return
	}

	latestRollupHeader, err := o.obsClient.RollupHeaderByNumber(big.NewInt(int64(numLatestRollup)))
	if err != nil || latestRollupHeader == nil {
		o.logger.Error("Could not fetch latest rollup.", log.ErrKey, err)
		logAndSendErr(resp, "Could not fetch average rollup time.")
		return
	}

	avgRollupTime := float64(latestRollupHeader.Time-firstRollupHeader.Time) / float64(numLatestRollup)
	_, err = resp.Write([]byte(fmt.Sprintf("%.2f", avgRollupTime)))
	if err != nil {
		o.logger.Error("could not return average rollup time to client.", log.ErrKey, err)
		logAndSendErr(resp, "Could not fetch average rollup time.")
		return
	}
}

// Retrieves the last five rollup numbers.
func (o *Obscuroscan) getLatestRollups(resp http.ResponseWriter, _ *http.Request) {
	latestRollupNum, err := o.obsClient.RollupNumber()
	if err != nil {
		o.logger.Error("Could not fetch latest rollups.", log.ErrKey, err)
		logAndSendErr(resp, "Could not fetch latest rollups.")
		return
	}

	// We walk the chain of rollups, getting the number for the most recent five.
	rollupNums := make([]string, 5)
	for idx := 0; idx < 5; idx++ {
		rollupNum := int(latestRollupNum) - idx
		if rollupNum < 0 {
			// If there are less than five rollups, we return an N/A.
			rollupNums[idx] = "N/A"
		} else {
			rollupNums[idx] = strconv.Itoa(rollupNum)
		}
	}

	jsonRollupNums, err := json.Marshal(rollupNums)
	if err != nil {
		o.logger.Error("could not return latest rollups to client.", log.ErrKey, err)
		logAndSendErr(resp, "Could not fetch latest rollups.")
		return
	}
	_, err = resp.Write(jsonRollupNums)
	if err != nil {
		o.logger.Error("could not return latest rollups to client.", log.ErrKey, err)
		logAndSendErr(resp, "Could not fetch latest rollups.")
		return
	}
}

// Retrieves the last five transaction hashes.
func (o *Obscuroscan) getLatestTxs(resp http.ResponseWriter, _ *http.Request) {
	numTransactions := 5

	var txHashes []gethcommon.Hash
	err := o.client.Call(&txHashes, rpc.GetLatestTxs, numTransactions)
	if err != nil {
		o.logger.Error("Could not fetch latest transactions.", log.ErrKey, err)
		logAndSendErr(resp, "Could not fetch latest transactions.")
	}

	// We convert the hashes to strings and pad with N/As as needed.
	txHashStrings := make([]string, numTransactions)
	for idx := 0; idx < numTransactions; idx++ {
		if idx < len(txHashes) {
			txHashStrings[idx] = txHashes[idx].String()
		} else {
			txHashStrings[idx] = "N/A"
		}
	}

	jsonTxHashes, err := json.Marshal(txHashStrings)
	if err != nil {
		o.logger.Error("could not return latest transaction hashes to client.", log.ErrKey, err)
		logAndSendErr(resp, "Could not fetch latest transactions.")
		return
	}
	_, err = resp.Write(jsonTxHashes)
	if err != nil {
		o.logger.Error("could not return latest transaction hashes to client.", log.ErrKey, err)
		logAndSendErr(resp, "Could not fetch latest transactions.")
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
		o.logger.Error("could not read request body.", log.ErrKey, err)
		logAndSendErr(resp, "Could not fetch block.")
		return
	}
	blockHashStr := buffer.String()
	blockHash := gethcommon.HexToHash(blockHashStr)

	var blockHeader *types.Header
	err = o.client.Call(&blockHeader, rpc.GetBlockHeaderByHash, blockHash)
	if err != nil {
		o.logger.Error(fmt.Sprintf("could not retrieve block with hash %s", blockHash), log.ErrKey, err)
		logAndSendErr(resp, "Could not fetch block.")
		return
	}

	jsonBlock, err := json.Marshal(blockHeader)
	if err != nil {
		o.logger.Error("could not return block to client.", log.ErrKey, err)
		logAndSendErr(resp, "Could not fetch block.")
		return
	}
	_, err = resp.Write(jsonBlock)
	if err != nil {
		o.logger.Error("could not return block to client.", log.ErrKey, err)
		logAndSendErr(resp, "Could not fetch block.")
		return
	}
}

// Retrieves a rollup given its number or the hash of a transaction it contains.
func (o *Obscuroscan) getRollupByNumOrTxHash(resp http.ResponseWriter, req *http.Request) {
	body := req.Body
	defer body.Close()
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(body)
	if err != nil {
		o.logger.Error("could not read request body.", log.ErrKey, err)
		logAndSendErr(resp, "Could not fetch rollup.")
		return
	}

	var rollup *common.ExtRollup
	if strings.HasPrefix(buffer.String(), "0x") {
		// A "0x" prefix indicates that we should retrieve the rollup by transaction hash.
		txHash := gethcommon.HexToHash(buffer.String())

		err = o.client.Call(&rollup, rpc.GetRollupForTx, txHash)
		if err != nil || rollup.Header == nil {
			o.logger.Error("could not retrieve rollup.", log.ErrKey, err)
			logAndSendErr(resp, fmt.Sprintf("Could not fetch rollup for transaction %s.", txHash))
			return
		}
	} else {
		// Otherwise, we treat the input as a rollup number.
		rollupNumber, err := strconv.Atoi(buffer.String())
		if err != nil {
			o.logger.Error(fmt.Sprintf("could not parse \"%s\" as an integer", buffer.String()))
			logAndSendErr(resp, fmt.Sprintf("Could not parse number %s.", buffer.String()))
			return
		}
		rollup, err = o.getRollupByNumber(int64(rollupNumber))
		if err != nil {
			o.logger.Error("Could not fetch rollup.", log.ErrKey, err)
			logAndSendErr(resp, fmt.Sprintf("Could not fetch rollup for number %d.", rollupNumber))
			return
		}
	}

	jsonRollup, err := json.Marshal(rollup)
	if err != nil {
		o.logger.Error("could not return rollup to client.", log.ErrKey, err)
		logAndSendErr(resp, "Could not marshal rollup to JSON.")
		return
	}
	_, err = resp.Write(jsonRollup)
	if err != nil {
		o.logger.Error("could not return rollup to client.", log.ErrKey, err)
		logAndSendErr(resp, "Could not return rollup to client.")
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
		o.logger.Error("could not read request body.", log.ErrKey, err)
		logAndSendErr(resp, "Could not decrypt transaction blob.")
		return
	}

	jsonTxs, err := decryptTxBlob(buffer.Bytes())
	if err != nil {
		o.logger.Error("could not decrypt transaction blob.", log.ErrKey, err)
		logAndSendErr(resp, "Could not decrypt transaction blob.")
		return
	}

	_, err = resp.Write(jsonTxs)
	if err != nil {
		o.logger.Error("could not write decrypted transactions to client.", log.ErrKey, err)
		logAndSendErr(resp, "Could not decrypt transaction blob.")
		return
	}
}

// Retrieves the node's attestation.
func (o *Obscuroscan) attestation(resp http.ResponseWriter, _ *http.Request) {
	var attestation *common.AttestationReport
	err := o.client.Call(&attestation, rpc.Attestation)
	if err != nil {
		o.logger.Error("could not retrieve node's attestation.", log.ErrKey, err)
		logAndSendErr(resp, "Could not retrieve node's attestation.")
		return
	}

	jsonAttestation, err := json.Marshal(attestation)
	if err != nil {
		o.logger.Error("could not convert node's attestation to JSON.", log.ErrKey, err)
		logAndSendErr(resp, "Could not retrieve node's attestation.")
		return
	}
	_, err = resp.Write(jsonAttestation)
	if err != nil {
		o.logger.Error("could not return JSON attestation to client.", log.ErrKey, err)
		logAndSendErr(resp, "Could not retrieve node's attestation.")
		return
	}
}

// Retrieves the node's attestation report.
func (o *Obscuroscan) attestationReport(resp http.ResponseWriter, _ *http.Request) {
	var attestation *common.AttestationReport
	err := o.client.Call(&attestation, rpc.Attestation)
	if err != nil {
		o.logger.Error("could not retrieve node's attestation.", log.ErrKey, err)
		logAndSendErr(resp, "Could not verify node's attestation.")
		return
	}

	// If DCAP isn't set up, verifying the report will send a SIGSYS signal. We catch this, so that it doesn't crash the program.
	sigChannel := make(chan os.Signal, 1)
	defer signal.Stop(sigChannel)
	signal.Notify(sigChannel, syscall.SIGSYS)
	attestationReport, err := enclave.VerifyRemoteReport(attestation.Report)
	signal.Stop(sigChannel)

	if err != nil {
		o.logger.Error("could not verify node's attestation.", log.ErrKey, err)
		logAndSendErr(resp, "Could not verify node's attestation.")
		return
	}

	attestationReportExt := attestationReportExternal{
		Data:            attestationReport.Data,
		SecurityVersion: attestationReport.SecurityVersion,
		Debug:           attestationReport.Debug,
		UniqueID:        attestationReport.UniqueID,
		SignerID:        attestationReport.SignerID,
		ProductID:       attestationReport.ProductID,
		TCBStatus:       attestationReport.TCBStatus.String(),
	}

	jsonAttestationReport, err := json.Marshal(attestationReportExt)
	if err != nil {
		o.logger.Error("could not convert node's attestation report to JSON.", log.ErrKey, err)
		logAndSendErr(resp, "Could not verify node's attestation.")
		return
	}
	_, err = resp.Write(jsonAttestationReport)
	if err != nil {
		o.logger.Error("could not return JSON attestation report to client.", log.ErrKey, err)
		logAndSendErr(resp, "Could not verify node's attestation.")
		return
	}
}

// Returns the rollup with the given number.
func (o *Obscuroscan) getRollupByNumber(rollupNumber int64) (*common.ExtRollup, error) {
	// TODO - If required, consolidate the two calls below into a single RPCGetRollupByNumber call to minimise round trips.
	rollupHeader, err := o.obsClient.RollupHeaderByNumber(big.NewInt(rollupNumber))
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rollup with number %d. Cause: %w", rollupNumber, err)
	}

	var rollup *common.ExtRollup
	err = o.client.Call(&rollup, rpc.GetRollup, rollupHeader.Hash())
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rollup. Cause: %w", err)
	}

	if rollup.Header == nil {
		return nil, fmt.Errorf("retrieved rollup had a nil header")
	}

	return rollup, nil
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

	// The nonce is prepended to the ciphertext.
	nonce := encryptedTxBytes[0:crypto.NonceLength]
	ciphertext := encryptedTxBytes[crypto.NonceLength:]
	encodedTxs, err := transactionCipher.Open(nil, nonce, ciphertext, nil)
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
