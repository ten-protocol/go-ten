package sql

// todo: was hoping to use the edgeless tool ERA as a dependency for verifying the edgeless db attestations but I'm hitting
//   dependency issues from `eclient` so progressing by duplicating relevant part of ERA using their `enclave` client
//   library instead of `eclient`

// ERA (Edgeless remote attestation) is a simple protocol for verifying edgeless tools. It's basically just a small json schema
// that they use as a standard data blob to encrypt into their attestation reports, includes signerID, security version etc.

// The only change from https://github.com/edgelesssys/era/blob/master/era/era.go is the use of enclave.VerifyRemoteReport
// in place of eclient.VerifyRemoteReport

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/edgelesssys/ego/attestation"
	"github.com/edgelesssys/ego/attestation/tcbstatus"
	"github.com/edgelesssys/ego/enclave"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
)

type certQuoteResp struct {
	Cert  string
	Quote []byte
}

// ErrEmptyQuote defines an error type when no quote was received. This likely occurs when the host is running in OE Simulation mode.
var ErrEmptyQuote = errors.New("no quote received")

// GetCertificate gets the TLS certificate from the server in PEM format. It performs remote attestation
// to verify the certificate. A config file must be provided that contains the attestation metadata.
func GetCertificate(host, configFilename string) ([]*pem.Block, tcbstatus.Status, error) {
	config, err := ioutil.ReadFile(configFilename)
	if err != nil {
		return nil, tcbstatus.Unknown, err
	}
	return getCertificate(host, config, enclave.VerifyRemoteReport)
}

// InsecureGetCertificate gets the TLS certificate from the server in PEM format, but does not perform remote attestation.
func InsecureGetCertificate(host string) ([]*pem.Block, error) {
	certs, _, err := getCertificate(host, nil, nil)
	return certs, err
}

type verifyFunc func([]byte) (attestation.Report, error)

func getCertificate(host string, config []byte, verifyRemoteReport verifyFunc) ([]*pem.Block, tcbstatus.Status, error) {
	cert, quote, err := httpGetCertQuote(&tls.Config{InsecureSkipVerify: true}, host, "quote")
	if err != nil {
		return nil, tcbstatus.Unknown, err
	}

	var certs []*pem.Block
	block, rest := pem.Decode([]byte(cert))
	if block == nil {
		return nil, tcbstatus.Unknown, errors.New("could not parse certificate")
	}
	certs = append(certs, block)

	// If we get more than one certificate, append it to the slice
	for len(rest) > 0 {
		block, rest = pem.Decode([]byte(rest))
		if block == nil {
			return nil, tcbstatus.Unknown, errors.New("could not parse certificate chain")
		}
		certs = append(certs, block)
	}

	if verifyRemoteReport == nil {
		return certs, tcbstatus.Unknown, nil
	}

	if len(quote) == 0 {
		return nil, tcbstatus.Unknown, ErrEmptyQuote
	}

	report, verifyErr := verifyRemoteReport(quote)
	if verifyErr != nil && verifyErr != attestation.ErrTCBLevelInvalid {
		return nil, tcbstatus.Unknown, verifyErr
	}

	// Use Root CA (last entry in certs) for attestation
	certRaw := certs[len(certs)-1].Bytes

	if err := verifyReport(report, certRaw, config); err != nil {
		return nil, tcbstatus.Unknown, err
	}

	return certs, report.TCBStatus, verifyErr
}

func verifyReport(report attestation.Report, cert []byte, config []byte) error {
	hash := sha256.Sum256(cert)
	if !bytes.Equal(report.Data[:len(hash)], hash[:]) {
		return errors.New("report data does not match the certificate's hash")
	}

	var cfg struct {
		SecurityVersion uint
		UniqueID        string
		SignerID        string
		ProductID       uint16
		Debug           bool
	}
	if err := json.Unmarshal(config, &cfg); err != nil {
		return err
	}
	if cfg.UniqueID == "" {
		if cfg.SecurityVersion == 0 {
			return errors.New("missing securityVersion in config")
		}
		if cfg.ProductID == 0 {
			return errors.New("missing productID in config")
		}
	}

	if cfg.SecurityVersion != 0 && report.SecurityVersion < cfg.SecurityVersion {
		return errors.New("invalid security version")
	}
	if cfg.ProductID != 0 && binary.LittleEndian.Uint16(report.ProductID) != cfg.ProductID {
		return errors.New("invalid product")
	}
	if report.Debug && !cfg.Debug {
		return errors.New("debug enclave not allowed")
	}
	if err := verifyID(cfg.UniqueID, report.UniqueID, "unqiueID"); err != nil {
		return err
	}
	if err := verifyID(cfg.SignerID, report.SignerID, "signerID"); err != nil {
		return err
	}
	if cfg.UniqueID == "" && cfg.SignerID == "" {
		fmt.Println("Warning: Configuration contains neither uniqueID nor signerID!")
	}

	return nil
}

func verifyID(expected string, actual []byte, name string) error {
	if expected == "" {
		return nil
	}
	expectedBytes, err := hex.DecodeString(expected)
	if err != nil {
		return err
	}
	if !bytes.Equal(expectedBytes, actual) {
		return errors.New("invalid " + name)
	}
	return nil
}

func httpGetCertQuote(tlsConfig *tls.Config, host, path string) (string, []byte, error) {
	client := http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
	url := url.URL{Scheme: "https", Host: host, Path: path}
	resp, err := client.Get(url.String())
	if err != nil {
		return "", nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	var certquote certQuoteResp
	if err != nil {
		return "", nil, err
	}

	/* Newer versions of Marblerun use a common JSON output format in which the quote
	is embedded into "data" and the error messages are stored inside "message".
	To keep compability with older versions, check if this block exists or
	if we get the data back directly. */

	if resp.StatusCode != http.StatusOK {
		errorMessage := gjson.GetBytes(body, "message")
		if errorMessage.Exists() {
			return "", nil, errors.New(resp.Status + ": " + errorMessage.String())
		}
		return "", nil, errors.New(resp.Status + ": " + string(body))
	}

	quoteData := gjson.GetBytes(body, "data")
	if quoteData.Exists() {
		err = json.Unmarshal([]byte(quoteData.String()), &certquote)
	} else {
		err = json.Unmarshal(body, &certquote)
	}

	if err != nil {
		return "", nil, err
	}
	resp.Body.Close()
	return certquote.Cert, certquote.Quote, nil
}
