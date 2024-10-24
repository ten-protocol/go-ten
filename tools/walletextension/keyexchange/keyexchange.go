package keyexchange

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
)

// GetEncryptionKey initializes the KeyExchange module.
func GetEncryptionKey(config common.Config, logger gethlog.Logger) ([]byte, error) {
	// If the database type is not sqlite or the existing gateway URL is not set, generate a new encryption key and use it for encrypting values in the database
	if config.DBType != "sqlite" || config.ExistingGatewayURL == "" {
		logger.Info("Generating a new encryption key")
		// Generate a new encryption key
		key, err := common.GenerateRandomKey()
		if err != nil {
			return nil, err
		}
		return key, nil
	}

	logger.Info("Starting key exschange process with existing gateway URL", "url", config.ExistingGatewayURL)

	fmt.Println("Generated a new key pair")
	// Generate a new ECDSA key pair
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	// Derive the public key from the private key
	publicKey := privateKey.Public()

	logger.Info("Generated new key pair", "publicKey", publicKey)

	// Get Attestation from Intel

	fmt.Println("Using existing gateway URL for encryption key exchange")
	return nil, nil
}
