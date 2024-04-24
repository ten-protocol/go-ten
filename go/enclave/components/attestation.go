package components

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/edgelesssys/ego/enclave"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

type IDData struct {
	EnclaveID   gethcommon.Address
	PubKey      []byte
	HostAddress string
}

type AttestationProvider interface {
	// GetReport returns the verifiable attestation report
	GetReport(ctx context.Context, pubKey []byte, enclaveID gethcommon.Address, hostAddress string) (*common.AttestationReport, error)
	// VerifyReport returns the embedded report data
	VerifyReport(att *common.AttestationReport) ([]byte, error)
}

type EgoAttestationProvider struct{}

func (e *EgoAttestationProvider) GetReport(ctx context.Context, pubKey []byte, enclaveID gethcommon.Address, hostAddress string) (*common.AttestationReport, error) {
	idHash, err := getIDHash(enclaveID, pubKey, hostAddress)
	if err != nil {
		return nil, err
	}
	report, err := enclave.GetRemoteReport(idHash)
	if err != nil {
		return nil, err
	}

	return &common.AttestationReport{
		Report:      report,
		PubKey:      pubKey,
		EnclaveID:   enclaveID,
		HostAddress: hostAddress,
	}, nil
}

// todo (#1059) - we need to verify the hash is a recognized enclave - figure out how we solve for upgradability
func (e *EgoAttestationProvider) VerifyReport(att *common.AttestationReport) ([]byte, error) {
	remoteReport, err := enclave.VerifyRemoteReport(att.Report)
	if err != nil {
		return []byte{}, err
	}
	return remoteReport.Data, nil
}

type DummyAttestationProvider struct{}

func (e *DummyAttestationProvider) GetReport(ctx context.Context, pubKey []byte, enclaveID gethcommon.Address, hostAddress string) (*common.AttestationReport, error) {
	return &common.AttestationReport{
		Report:      []byte("MOCK REPORT"),
		PubKey:      pubKey,
		EnclaveID:   enclaveID,
		HostAddress: hostAddress,
	}, nil
}

func (e *DummyAttestationProvider) VerifyReport(att *common.AttestationReport) ([]byte, error) {
	return getIDHash(att.EnclaveID, att.PubKey, att.HostAddress)
}

// getIDHash provides a hash of identifying data to be included in an attestation report (or verified against the contents of an attestation report)
func getIDHash(enclaveID gethcommon.Address, pubKey []byte, hostAddress string) ([]byte, error) {
	idData := IDData{
		EnclaveID:   enclaveID,
		PubKey:      pubKey,
		HostAddress: hostAddress,
	}
	idJSON, err := json.Marshal(idData)
	if err != nil {
		return nil, fmt.Errorf("failed to format ID data as JSON. Cause: %w", err)
	}
	hash := sha256.Sum256(idJSON)
	return hash[:], nil
}

func VerifyIdentity(data []byte, att *common.AttestationReport) error {
	expectedIDHash, err := getIDHash(att.EnclaveID, att.PubKey, att.HostAddress)
	if err != nil {
		return fmt.Errorf("failed to create ID data to check attestation report with enclaveID: %s. Cause: %w", att.EnclaveID, err)
	}
	// we trim the actual data because data extracted from the verified attestation is always 64 bytes long (padded with zeroes at the end)
	if !bytes.Equal(expectedIDHash, data[:len(expectedIDHash)]) {
		return fmt.Errorf("failed to verify hash for attestation report with enclaveID: %s", att.EnclaveID)
	}
	return nil
}
