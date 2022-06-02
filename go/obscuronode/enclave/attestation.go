package enclave

import (
	"crypto/sha256"
	"encoding/json"

	"github.com/edgelesssys/ego/enclave"
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/obscurocommon"
)

type IDData struct {
	Owner  common.Address
	PubKey []byte
}

type AttestationProvider interface {
	// GetReport returns the verifiable attestation report
	GetReport(pubKey []byte, owner common.Address) (*obscurocommon.AttestationReport, error)
	// VerifyReport returns the embedded report data
	VerifyReport(att *obscurocommon.AttestationReport) ([]byte, error)
}

type EgoAttestationProvider struct{}

func (e *EgoAttestationProvider) GetReport(pubKey []byte, owner common.Address) (*obscurocommon.AttestationReport, error) {
	idHash := getIDHash(owner, pubKey)
	report, err := enclave.GetRemoteReport(idHash)
	if err != nil {
		return nil, err
	}

	return &obscurocommon.AttestationReport{
		Report: report,
		PubKey: pubKey,
		Owner:  owner,
	}, nil
}

// TODO: we need to verify the hash is a recognized enclave - figure out how we solve for upgradability
// todo: we should probably return other properties for manual verification, not just the data (e.g. validate code hash)
func (e *EgoAttestationProvider) VerifyReport(att *obscurocommon.AttestationReport) ([]byte, error) {
	remoteReport, err := enclave.VerifyRemoteReport(att.Report)
	if err != nil {
		return []byte{}, err
	}
	return remoteReport.Data, nil
}

type DummyAttestationProvider struct{}

func (e *DummyAttestationProvider) GetReport(pubKey []byte, owner common.Address) (*obscurocommon.AttestationReport, error) {
	return &obscurocommon.AttestationReport{
		Report: []byte("MOCK REPORT"),
		PubKey: pubKey,
		Owner:  owner,
	}, nil
}

func (e *DummyAttestationProvider) VerifyReport(att *obscurocommon.AttestationReport) ([]byte, error) {
	return getIDHash(att.Owner, att.PubKey), nil
}

// getIDHash provides a hash of identifying data to be included in an attestation report (or verified against the contents of an attestation report)
func getIDHash(owner common.Address, pubKey []byte) []byte {
	idData := IDData{
		Owner:  owner,
		PubKey: pubKey,
	}
	idJSON, _ := json.Marshal(idData)
	hash := sha256.Sum256(idJSON)
	return hash[:]
}
