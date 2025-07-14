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
	fmt.Printf("SealAndPersist: Starting to seal data to file %s\n", filepath)
	fmt.Printf("SealAndPersist: Content size: %d bytes, testEnvSealOnly: %t\n", len(contents), testEnvSealOnly)
	
	f, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("SealAndPersist: Failed to create file %s - %v\n", filepath, err)
		return fmt.Errorf("failed to create file %s - %w", filepath, err)
	}
	defer func() {
		_ = f.Close()
	}()
	
	fmt.Printf("SealAndPersist: Successfully created file %s\n", filepath)

	sealMethod := ecrypto.SealWithUniqueKey
	if testEnvSealOnly {
		// todo (#1377) - remove this option - this is a stop-gap solution for upgradability in testnet while we implement the final solution
		// In prod this must not be used in this way, it would make the secret vulnerable to anyone that manages to get
		// access to the product signing key
		sealMethod = ecrypto.SealWithProductKey
		fmt.Printf("SealAndPersist: Using ProductKey sealing method (testEnvSealOnly=true)\n")
	} else {
		fmt.Printf("SealAndPersist: Using UniqueKey sealing method\n")
	}
	
	// We need to seal with a key derived from the measurement of the enclave to prevent the signer from decrypting the secret.
	fmt.Printf("SealAndPersist: Attempting to seal %d bytes using SGX enclave key\n", len(contents))
	enc, err := sealMethod([]byte(contents), nil)
	if err != nil {
		fmt.Printf("SealAndPersist: SGX sealing failed - %v\n", err)
		fmt.Printf("SealAndPersist: This indicates SGX enclave issues or missing SGX support\n")
		return fmt.Errorf("failed to seal contents bytes with enclave key to persist in %s - %w", filepath, err)
	}
	fmt.Printf("SealAndPersist: Successfully sealed data, encrypted size: %d bytes\n", len(enc))
	_, err = f.Write(enc)
	if err != nil {
		fmt.Printf("SealAndPersist: Failed to write sealed data to file %s - %v\n", filepath, err)
		return fmt.Errorf("failed to write manifest json file - %w", err)
	}
	
	fmt.Printf("SealAndPersist: Successfully wrote sealed data to file %s\n", filepath)
	return nil
}

// ReadAndUnseal reverses SealAndPersist, uses SGX product key it attempts to decrypt the file
func ReadAndUnseal(filepath string) ([]byte, error) {
	fmt.Printf("ReadAndUnseal: Starting to read and unseal file %s\n", filepath)
	
	b, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("ReadAndUnseal: Failed to read file %s - %v\n", filepath, err)
		return nil, err
	}
	fmt.Printf("ReadAndUnseal: Successfully read %d bytes from file %s\n", len(b), filepath)

	fmt.Printf("ReadAndUnseal: Attempting to unseal %d bytes using SGX enclave key\n", len(b))
	data, err := ecrypto.Unseal(b, nil)
	if err != nil {
		fmt.Printf("ReadAndUnseal: SGX unsealing failed - %v\n", err)
		fmt.Printf("ReadAndUnseal: This indicates SGX enclave issues, key mismatch, or corrupted sealed data\n")
		return nil, err
	}
	fmt.Printf("ReadAndUnseal: Successfully unsealed data, decrypted size: %d bytes\n", len(data))
	return data, nil
}

func GetEnclaveSignerPublicKey() ([]byte, error) {
	fmt.Printf("GetEnclaveSignerPublicKey: Starting to get enclave signer public key\n")
	
	// Get a local report with empty report data
	fmt.Printf("GetEnclaveSignerPublicKey: Requesting local report from SGX enclave\n")
	report, err := enclave.GetLocalReport(nil, nil)
	if err != nil {
		fmt.Printf("GetEnclaveSignerPublicKey: Failed to get local report - %v\n", err)
		fmt.Printf("GetEnclaveSignerPublicKey: This indicates SGX enclave initialization issues\n")
		return nil, fmt.Errorf("failed to get local report: %w", err)
	}
	fmt.Printf("GetEnclaveSignerPublicKey: Successfully obtained local report, size: %d bytes\n", len(report))

	// Verify the local report to get the attestation report with parsed claims
	fmt.Printf("GetEnclaveSignerPublicKey: Verifying local report\n")
	attestationReport, err := enclave.VerifyLocalReport(report)
	if err != nil {
		fmt.Printf("GetEnclaveSignerPublicKey: Failed to verify local report - %v\n", err)
		fmt.Printf("GetEnclaveSignerPublicKey: This indicates corrupted report or SGX verification issues\n")
		return nil, fmt.Errorf("failed to verify local report: %w", err)
	}
	fmt.Printf("GetEnclaveSignerPublicKey: Successfully verified local report\n")

	// The SignerID field contains the public key hash of the enclave signer
	fmt.Printf("GetEnclaveSignerPublicKey: Extracting signer ID, size: %d bytes\n", len(attestationReport.SignerID))
	return attestationReport.SignerID, nil
}
