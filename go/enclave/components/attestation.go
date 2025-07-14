package components

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/enclave/crypto"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/edgelesssys/ego/enclave"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

type IDData struct {
	EnclaveID   gethcommon.Address
	PubKey      []byte
	HostAddress string
}

// AttestationProvider creates and verifies attestation reports
type AttestationProvider interface {
	// CreateAttestationReport returns the verifiable attestation report
	CreateAttestationReport(ctx context.Context, hostAddress string) (*common.AttestationReport, error)
	// VerifyReport returns the embedded report data
	VerifyReport(att *common.AttestationReport) ([]byte, error)
}

func NewAttestationProvider(enclaveKeyService *crypto.EnclaveAttestedKeyService, willAttest bool, logger gethlog.Logger) AttestationProvider {
	if willAttest {
		return &EgoAttestationProvider{enclaveKeyService: enclaveKeyService, logger: logger}
	}
	logger.Warn("WARNING - Attestation is not enabled, enclave will not create a verified attestation report.")
	return &DummyAttestationProvider{enclaveKeyService: enclaveKeyService}
}

type EgoAttestationProvider struct {
	enclaveKeyService *crypto.EnclaveAttestedKeyService
	logger            gethlog.Logger
}

func (e *EgoAttestationProvider) CreateAttestationReport(ctx context.Context, hostAddress string) (*common.AttestationReport, error) {
	e.logger.Info("CreateAttestationReport: Starting attestation report creation", "hostAddress", hostAddress)
	e.logger.Info("CreateAttestationReport: Enclave info", "enclaveID", e.enclaveKeyService.EnclaveID(), "pubKeySize", len(e.enclaveKeyService.PublicKeyBytes()))
	
	idHash, err := getIDHash(e.enclaveKeyService.EnclaveID(), e.enclaveKeyService.PublicKeyBytes(), hostAddress)
	if err != nil {
		e.logger.Error("CreateAttestationReport: Failed to create ID hash", "error", err)
		return nil, err
	}
	e.logger.Info("CreateAttestationReport: Generated ID hash", "hashSize", len(idHash), "hash", fmt.Sprintf("%x", idHash[:8]))
	
	e.logger.Info("CreateAttestationReport: Calling enclave.GetRemoteReport()")
	report, err := enclave.GetRemoteReport(idHash)
	if err != nil {
		e.logger.Error("CreateAttestationReport: Failed to get remote report from SGX enclave", "error", err)
		e.logger.Error("CreateAttestationReport: SGX remote report failure indicates hardware/driver issues, PCCS problems, or missing SGX services")
		return nil, err
	}
	e.logger.Info("CreateAttestationReport: Successfully obtained remote report", "reportSize", len(report))

	attestationReport := &common.AttestationReport{
		Report:      report,
		PubKey:      e.enclaveKeyService.PublicKeyBytes(),
		EnclaveID:   e.enclaveKeyService.EnclaveID(),
		HostAddress: hostAddress,
	}
	e.logger.Info("CreateAttestationReport: Attestation report created successfully")
	return attestationReport, nil
}

// todo (#1059) - we need to verify the hash is a recognized enclave - figure out how we solve for upgradability
func (e *EgoAttestationProvider) VerifyReport(att *common.AttestationReport) ([]byte, error) {
	e.logger.Info("VerifyReport: Starting attestation report verification", "reportSize", len(att.Report))
	e.logger.Info("VerifyReport: Report details", "enclaveID", att.EnclaveID, "hostAddress", att.HostAddress, "pubKeySize", len(att.PubKey))
	
	e.logger.Info("VerifyReport: Calling enclave.VerifyRemoteReport()")
	remoteReport, err := enclave.VerifyRemoteReport(att.Report)
	if err != nil {
		e.logger.Error("VerifyReport: Failed to verify remote report", "error", err)
		e.logger.Error("VerifyReport: Verification failure indicates invalid report, SGX issues, or PCCS problems")
		return []byte{}, err
	}
	e.logger.Info("VerifyReport: Successfully verified remote report", "dataSize", len(remoteReport.Data))
	return remoteReport.Data, nil
}

type DummyAttestationProvider struct {
	enclaveKeyService *crypto.EnclaveAttestedKeyService
}

func (e *DummyAttestationProvider) CreateAttestationReport(ctx context.Context, hostAddress string) (*common.AttestationReport, error) {
	return &common.AttestationReport{
		Report:      []byte("MOCK REPORT"),
		PubKey:      e.enclaveKeyService.PublicKeyBytes(),
		EnclaveID:   e.enclaveKeyService.EnclaveID(),
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
