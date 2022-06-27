package egoutils

import (
	"fmt"
	"os"

	"github.com/edgelesssys/ego/ecrypto"
)

// SealAndPersist uses SGX's ProductKey to encrypt contents string to filepath string
//	(note: filepath location must be accessible via ego mounts config in enclave.json)
func SealAndPersist(contents string, filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s - %w", filepath, err)
	}
	defer func() {
		_ = f.Close()
	}()

	// todo: do we prefer to seal with product key for upgradability or unique key to require fresh db with every code change
	enc, err := ecrypto.SealWithProductKey([]byte(contents), nil)
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
		return nil, err
	}

	data, err := ecrypto.Unseal(b, nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}
