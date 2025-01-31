package edgelessdb

// An obscuro implementation of the ERA (Edgeless remote attestation) tool (which is basically just a small json schema
// that Edgeless use as a standard data blob to encrypt into their attestation reports, includes signerID, security version etc.)

// Initially forked from https://github.com/edgelesssys/era/blob/master/era/era.go but reworked to have an obscuro-friendly
// 	 API that just included the things we were using.

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
	"io"
	"net/http"
	"net/url"

	"github.com/edgelesssys/ego/attestation"
	"github.com/edgelesssys/ego/attestation/tcbstatus"
	"github.com/edgelesssys/ego/enclave"
	gethlog "github.com/ethereum/go-ethereum/log"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"
	"github.com/tidwall/gjson"
)

const (
	quoteEndpoint = "quote"
	msgJSONField  = "message"
	dataJSONField = "data"
)

type certQuoteResp struct {
	Cert  string
	Quote []byte
}

type EdgelessAttestationConstraints struct {
	// This triplet of fields is typically used to attest an instance of an edgeless product (again, see ERA docs for more info)
	SecurityVersion uint   // Min required security version of the Edgeless product
	SignerID        string // corresponds to MRSIGNER SGX data, the expected fingerprint of Edgeless System's signing key
	ProductID       uint16 // The ID for the edgeless product, ProductID = 16 for Edgeless DB

	// Alternative to the triplet above you can specify a UniqueID which corresponds to a specific enclave package
	UniqueID string // This corresponds to the MRENCLAVE field in the SGX attestation data, it is stricter than the triplet above

	Debug bool // while debugging this can be set to true to permit debug attestations to pass verification
}

// The values here were the latest values from: https://github.com/edgelesssys/edgelessdb/releases/latest/download/edgelessdb-sgx.json
// This should probably be configurable rather than relying on this hardcoded snapshot
var defaultEDBConstraints = &EdgelessAttestationConstraints{
	SecurityVersion: 2,
	SignerID:        "67d7b00741440d29922a15a9ead427b6faf1d610238ae9826da345cea4fee0fe",
	ProductID:       16,
}

// performEDBRemoteAttestation perform the SGX enclave attestation to verify edb running in a legit enclave and with expected edb version etc.
func performEDBRemoteAttestation(config enclaveconfig.EnclaveConfig, edbHost string, constraints *EdgelessAttestationConstraints, logger gethlog.Logger) (string, error) {
	logger.Info("Verifying attestation from edgeless DB...")
	edbHTTPAddr := fmt.Sprintf("%s:%s", edbHost, edbHTTPPort)
	certs, tcbStatus, err := performRAAndFetchTLSCert(config, edbHTTPAddr, constraints)
	if err != nil {
		// todo (#1550) - should we check the error type with: err == attestation.ErrTCBLevelInvalid?
		// for now it's maximum strictness (we can revisit this and permit some tcbStatuses if desired)
		return "", fmt.Errorf("attestation failed, host=%s, tcbStatus=%s, err=%w", edbHTTPAddr, tcbStatus, err)
	}
	if len(certs) == 0 {
		return "", fmt.Errorf("no certificates found from edgeless db attestation process")
	}

	logger.Info("Successfully verified edb attestation and retrieved certificate.")
	// the last cert in the list is the CA
	return string(pem.EncodeToMemory(certs[len(certs)-1])), nil
}

// performRAAndFetchTLSCert gets the TLS certificate from the Edgeless DB server in PEM format. It performs remote attestation
// to verify the certificate. Attestation constraints must be provided to validate against.
func performRAAndFetchTLSCert(enclaveConfig enclaveconfig.EnclaveConfig, host string, constraints *EdgelessAttestationConstraints) ([]*pem.Block, tcbstatus.Status, error) {
	// we don't need to verify the TLS because we will be verifying the attestation report and that can't be faked
	cert, quote, err := httpGetCertQuote(&tls.Config{InsecureSkipVerify: true}, host, quoteEndpoint) //nolint:gosec
	if err != nil {
		return nil, tcbstatus.Unknown, err
	}

	if len(quote) == 0 && enclaveConfig.WillAttest {
		return nil, tcbstatus.Unknown, errors.New("no quote found, attestation failed")
	}

	block, rest := pem.Decode([]byte(cert))
	if block == nil {
		return nil, tcbstatus.Unknown, errors.New("could not parse certificate")
	}
	certs := []*pem.Block{block}

	// If we get more than one certificate, append it to the slice
	for len(rest) > 0 {
		block, rest = pem.Decode(rest)
		if block == nil {
			return nil, tcbstatus.Unknown, errors.New("could not parse certificate chain")
		}
		certs = append(certs, block)
	}

	if !enclaveConfig.WillAttest {
		return certs, tcbstatus.Unknown, nil
	}

	report, verifyErr := enclave.VerifyRemoteReport(quote)
	// depending on how strict you are being, some invalid TCBLevels would be acceptable (e.g. something might need an
	//		upgrade but not have any known vulnerabilities). That's why we proceed when TCBLevelInvalid and let caller decide
	if verifyErr != nil && !errors.Is(verifyErr, attestation.ErrTCBLevelInvalid) {
		return nil, tcbstatus.Unknown, verifyErr
	}

	// Use Root CA (last entry in certs) for attestation
	certRaw := certs[len(certs)-1].Bytes

	if err = checkAttestationConstraints(report, certRaw, constraints); err != nil {
		return nil, tcbstatus.Unknown, err
	}

	return certs, report.TCBStatus, verifyErr
}

func checkAttestationConstraints(report attestation.Report, cert []byte, constraints *EdgelessAttestationConstraints) error {
	hash := sha256.Sum256(cert)
	if !bytes.Equal(report.Data[:len(hash)], hash[:]) {
		return errors.New("report data does not match the certificate's hash")
	}

	if constraints.UniqueID == "" {
		if constraints.SecurityVersion == 0 {
			return errors.New("missing securityVersion in config")
		}
		if constraints.ProductID == 0 {
			return errors.New("missing productID in config")
		}
	}

	if constraints.SecurityVersion != 0 && report.SecurityVersion < constraints.SecurityVersion {
		return errors.New("invalid security version")
	}
	if constraints.ProductID != 0 && binary.LittleEndian.Uint16(report.ProductID) != constraints.ProductID {
		return errors.New("invalid product")
	}
	if report.Debug && !constraints.Debug {
		return errors.New("debug enclave not allowed")
	}
	if err := verifyID(constraints.UniqueID, report.UniqueID, "uniqueID"); err != nil {
		return err
	}
	if err := verifyID(constraints.SignerID, report.SignerID, "signerID"); err != nil {
		return err
	}
	if constraints.UniqueID == "" && constraints.SignerID == "" {
		fmt.Println("Warning: Configured constraints contains neither uniqueID nor signerID! " +
			"This will not provide validation of the code running in the remote enclave.")
	}

	return nil
}

func verifyID(expected string, actual []byte, name string) error {
	if expected == "" {
		// we don't verify every ID, if expected is empty then don't check the value of actual
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
	resp, err := client.Get(url.String()) //nolint:noctx
	if err != nil {
		return "", nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	/* Newer versions of Marblerun use a common JSON output format in which the quote
	is embedded into "data" and the error messages are stored inside "message".
	To keep compatibility with older versions, check if this block exists or
	if we get the data back directly. */

	if resp.StatusCode != http.StatusOK {
		errorMessage := gjson.GetBytes(body, msgJSONField)
		if errorMessage.Exists() {
			return "", nil, errors.New(resp.Status + ": " + errorMessage.String())
		}
		return "", nil, errors.New(resp.Status + ": " + string(body))
	}

	var certQuote certQuoteResp
	quoteData := gjson.GetBytes(body, dataJSONField)
	if quoteData.Exists() {
		err = json.Unmarshal([]byte(quoteData.String()), &certQuote)
	} else {
		err = json.Unmarshal(body, &certQuote)
	}

	if err != nil {
		return "", nil, err
	}
	resp.Body.Close()
	return certQuote.Cert, certQuote.Quote, nil
}
