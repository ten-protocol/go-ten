package egoutils

import (
	"fmt"
	"os"

	"github.com/edgelesssys/ego/enclave"

	"github.com/edgelesssys/ego/ecrypto"
)

// SealAndPersist uses SGX's Unique measurement key to encrypt the contents string to filepath string
// (note: filepath location must be accessible via ego mounts config in enclave.json)
func SealAndPersist(contents string, filepath string, testEnvSealOnly bool) error {
	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s - %w", filepath, err)
	}
	defer func() {
		_ = f.Close()
	}()

	sealMethod := ecrypto.SealWithUniqueKey
	if testEnvSealOnly {
		// todo (#1377) - remove this option - this is a stop-gap solution for upgradability in testnet while we implement the final solution
		// In prod this must not be used in this way, it would make the secret vulnerable to anyone that manages to get
		// access to the product signing key
		sealMethod = ecrypto.SealWithProductKey
	}
	// We need to seal with a key derived from the measurement of the enclave to prevent the signer from decrypting the secret.
	enc, err := sealMethod([]byte(contents), nil)
	if err != nil {
		return fmt.Errorf("failed to seal contents bytes with enclave key to persist in %s - %w", filepath, err)
	}
	_, err = f.Write(enc)
	if err != nil {
		return fmt.Errorf("failed to write manifest json file - %w", err)
	}
	return nil
}

// ReadAndUnseal reverses SealAndPersist, uses SGX product key it attempts to decrypt the file
func ReadAndUnseal(filepath string) ([]byte, error) {
	b, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filepath, err)
	}

	data, err := ecrypto.Unseal(b, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to unseal data from %s (file size: %d bytes): %w", filepath, len(b), err)
	}
	return data, nil
}

func GetEnclaveSignerPublicKey() ([]byte, error) {
	// Get a local report with empty report data
	report, err := enclave.GetLocalReport(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get local report: %w", err)
	}

	// Verify the local report to get the attestation report with parsed claims
	attestationReport, err := enclave.VerifyLocalReport(report)
	if err != nil {
		return nil, fmt.Errorf("failed to verify local report: %w", err)
	}

	// The SignerID field contains the public key hash of the enclave signer
	return attestationReport.SignerID, nil
}
