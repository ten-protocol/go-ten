package common

import (
	"context"

	"github.com/edgelesssys/ego/enclave"
)

type GatewayAttestationReport struct {
	Report []byte // the signed bytes of the report which includes some encrypted identifying data
	PubKey []byte // a public key that can be used to send encrypted data back to the TEE securely (should only be used once Report has been verified)
	// EnclaveID   common.Address // address identifying the owner of the TEE which signed this report, can also be verified from the encrypted Report data  // TODO @ziga - ask if we need something like that.
	// HostAddress string // the IP address on which the host can be contacted by other Obscuro hosts for peer-to-peer communication // TODO @ziga - ask if we need something like that.
}

type AttestationProvider interface {
	// GetReport returns the verifiable attestation report
	GetReport(ctx context.Context, pubKey []byte) (*GatewayAttestationReport, error)
	// VerifyReport returns the embedded report data
	VerifyReport(att *GatewayAttestationReport) ([]byte, error)
}

// EgoAttestationProvider is a provider for attestation reports from the Ego TEE
type EgoAttestationProvider struct{}

func (e *EgoAttestationProvider) GetReport(ctx context.Context, pubKey []byte) (*GatewayAttestationReport, error) {
	// TODO @ziga - what do we need to pass to GetRemoteReport?
	// We need to pass the hash of what we want (check getIDHash in node attestation.go and do something similar here)
	report, err := enclave.GetRemoteReport(nil)
	if err != nil {
		return nil, err
	}

	return &GatewayAttestationReport{
		Report: report,
		PubKey: pubKey,
		// TODO add owner
	}, nil
}

// todo (#1059) - we need to verify the hash is a recognized enclave - figure out how we solve for upgradability
func (e *EgoAttestationProvider) VerifyReport(att *GatewayAttestationReport) ([]byte, error) {
	remoteReport, err := enclave.VerifyRemoteReport(att.Report)
	if err != nil {
		return []byte{}, err
	}
	return remoteReport.Data, nil
}
