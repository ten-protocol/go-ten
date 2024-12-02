package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/crypto/ecies"
)

func generateSecureEntropy(nrBytes int) ([]byte, error) {
	rndBytes := make([]byte, nrBytes)
	if _, err := io.ReadFull(rand.Reader, rndBytes); err != nil {
		return nil, err
	}
	return rndBytes, nil
}

func encryptWithPublicKey(msg []byte, pub *ecdsa.PublicKey) ([]byte, error) {
	ciphertext, err := ecies.Encrypt(rand.Reader, ecies.ImportECDSAPublic(pub), msg, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt with public key. %w", err)
	}
	return ciphertext, nil
}

func decryptWithPrivateKey(ciphertext []byte, priv *ecdsa.PrivateKey) ([]byte, error) {
	plaintext, err := ecies.ImportECDSA(priv).Decrypt(ciphertext, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt with private key. %w", err)
	}
	return plaintext, nil
}
