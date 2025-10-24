package keymanager

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/edgelesssys/ego/enclave"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	tencommon "github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave/core/egoutils"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
)

const (
	dataDir         = "/data"
	RSAKeySize      = 2048
	AzureHSMKeyName = "database-encryption-key" // Hardcoded HSM key name
)

var encryptionKeyFile = filepath.Join(dataDir, "encryption-key.json")

// KeyExchangeRequest represents the structure of the data sent from KeyRequester to KeyProvider
type KeyExchangeRequest struct {
	PublicKey   []byte `json:"public_key"`
	Attestation []byte `json:"attestation"`
}

// KeyExchangeResponse represents the structure of the data sent from KeyProvider to KeyRequester
type KeyExchangeResponse struct {
	EncryptedKey string `json:"encrypted_key"` // Base64 encoded encrypted encryption key
}

// GetEncryptionKey returns the encryption key for the database
// - If we use an SQLite database, no encryption key is needed as SQLite typically runs in development or testing environments.
// - First try to unseal the existing key if there is one
// - If no existing key and HSM recovery is enabled, try to recover from HSM
// - If no existing key is available, we look at `encryptionKeySource`, if it is `new` we generate a new one, if it is a URL we attempt to perform key exchange
// - Finally, we seal the key locally (and if there was an existing key file that could not be unsealed, we store it as a backup)
// - If HSM backup is enabled, we always backup the key to Azure HSM (whether existing or newly generated)
func GetEncryptionKey(config common.Config, logger gethlog.Logger) ([]byte, error) {
	// check if we are using sqlite database and no encryption key needed
	if config.DBType == "sqlite" {
		logger.Info("using sqlite database, no encryption key needed - exiting key exchange process")
		return nil, nil
	}

	// Validate HSM configuration if HSM features are enabled
	if config.AzureHSMBackupEnabled || config.AzureHSMRecoveryEnabled {
		if config.AzureHSMURL == "" {
			logger.Crit("Azure HSM features are enabled but HSM URL is not configured")
			return nil, fmt.Errorf("Azure HSM URL is required when HSM features are enabled")
		}
		logger.Info("Azure HSM configuration", "backup_enabled", config.AzureHSMBackupEnabled, "recovery_enabled", config.AzureHSMRecoveryEnabled, "url", config.AzureHSMURL)
	}

	// Step 1: First try to unseal an existing encryption key
	var encryptionKey []byte
	var found bool
	var err error

	encryptionKey, found, err = tryUnsealKey(encryptionKeyFile, config.InsideEnclave)
	if err != nil {
		logger.Warn("failed to unseal existing encryption key", "error", err)
	}

	if found {
		logger.Info("successfully unsealed existing encryption key")

		// Always backup to HSM if backup is enabled, even for existing keys
		if config.AzureHSMBackupEnabled {
			logger.Info("backing up existing key to Azure HSM")
			err = backupKeyToHSM(encryptionKey, config, logger)
			if err != nil {
				// Backup failure is FATAL - fail completely
				logger.Crit("failed to backup existing encryption key to Azure HSM", log.ErrKey, err)
				return nil, fmt.Errorf("failed to backup existing encryption key to Azure HSM: %w", err)
			}
			logger.Info("successfully backed up existing encryption key to Azure HSM")
		}

		return encryptionKey, nil
	}

	// Step 2: Try to recover from Azure HSM if recovery is enabled
	if config.AzureHSMRecoveryEnabled {
		logger.Info("attempting HSM recovery")
		encryptionKey, found, err = recoverKeyFromHSM(config, logger)
		if err != nil {
			// HSM recovery failure is not fatal - log warning and continue to generate new key
			logger.Warn("failed to recover key from HSM, will generate new key", "error", err)
		}

		if found {
			logger.Info("successfully recovered encryption key from Azure HSM")

			// Seal the recovered key locally
			err = trySealKey(encryptionKey, encryptionKeyFile, config.InsideEnclave, logger)
			if err != nil {
				logger.Crit("unable to seal recovered encryption key", log.ErrKey, err)
				return nil, err
			}
			logger.Info("sealed recovered encryption key to local storage")

			return encryptionKey, nil
		}
	}

	// Step 3: No existing key found (neither local nor HSM), generate or exchange
	// If no encryptionKeySource is provided, fail since we couldn't unseal an existing key or recover from HSM
	if config.EncryptionKeySource == "" {
		logger.Crit("no sealed encryption key found, HSM recovery failed/disabled, and no key source provided", log.ErrKey, err)
		return nil, fmt.Errorf("no sealed encryption key found and no key source provided: %w", err)
	}

	// If the "new" keyword is used for the encryptionKeySource, generate a new random encryption key
	if config.EncryptionKeySource == "new" {
		logger.Info("encryptionKeySource set to 'new' -> generating new random encryption key")
		encryptionKey, err = common.GenerateRandomKey()
		if err != nil {
			logger.Crit("unable to generate random encryption key", log.ErrKey, err)
			return nil, err
		}
	} else {
		// If encryptionKeySource is a URL, attempt to perform key exchange with the specified key provider
		logger.Info(fmt.Sprintf("encryptionKeySource set to '%s', trying to get encryption key from key provider", config.EncryptionKeySource))
		encryptionKey, err = HandleKeyExchange(config, logger)
		if err != nil {
			logger.Crit("unable to get encryption key from key provider", log.ErrKey, err)
			return nil, err
		}
	}

	// Step 4: Seal the key locally (that we generated or got from key exchange)
	err = trySealKey(encryptionKey, encryptionKeyFile, config.InsideEnclave, logger)
	if err != nil {
		logger.Crit("unable to seal encryption key", log.ErrKey, err)
		return nil, err
	}
	logger.Info("sealed new encryption key to local storage")

	// Step 5: Backup to Azure HSM if backup is enabled
	if config.AzureHSMBackupEnabled {
		logger.Info("backing up key to Azure HSM")
		err = backupKeyToHSM(encryptionKey, config, logger)
		if err != nil {
			// Backup failure is FATAL - fail completely
			logger.Crit("failed to backup encryption key to Azure HSM", log.ErrKey, err)
			return nil, fmt.Errorf("failed to backup encryption key to Azure HSM: %w", err)
		}
		logger.Info("successfully backed up encryption key to Azure HSM")
	}

	return encryptionKey, nil
}

// tryUnsealKey attempts to unseal an encryption key from disk
// Returns (key, found, error)
func tryUnsealKey(keyPath string, isEnclave bool) ([]byte, bool, error) {
	// Only attempt unsealing if we're in an SGX enclave
	if !isEnclave {
		return nil, false, nil
	}

	// Read the key and unseal if possible
	data, err := egoutils.ReadAndUnseal(keyPath)
	if err != nil {
		return nil, false, err
	}

	return data, true, nil
}

// backupExistingKey creates a timestamped backup of an existing key file
func backupExistingKey(keyPath string, logger gethlog.Logger) error {
	// Check if the key file exists
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		// No existing key file to backup
		return nil
	}

	// Create backup filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.backup_%s", keyPath, timestamp)

	// Copy the existing file to backup location
	if err := copyFile(keyPath, backupPath); err != nil {
		logger.Warn("failed to backup existing encryption key", "error", err)
		return err
	}

	logger.Info("backed up existing encryption key", "backup_path", backupPath)
	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// trySealKey attempts to seal an encryption key to disk
// Only seals if running in an SGX enclave
// If an existing key file exists but failed to be unsealed, then it will be backed up before overwriting
func trySealKey(key []byte, keyPath string, isEnclave bool, logger gethlog.Logger) error {
	// Only attempt sealing if we're in an SGX enclave
	if !isEnclave {
		return nil
	}

	// Backup existing key file if it exists
	if err := backupExistingKey(keyPath, logger); err != nil {
		// Log warning but continue - backup failure shouldn't prevent key sealing
		logger.Warn("failed to backup existing key, proceeding with new key", "error", err)
	}

	// Seal and persist the key to /data/encryption.key
	if err := egoutils.SealAndPersist(string(key), keyPath, true); err != nil {
		return fmt.Errorf("failed to seal and persist key: %w", err)
	}
	return nil
}

// HandleKeyExchange handles the key exchange process from KeyRequester side.
func HandleKeyExchange(config common.Config, logger gethlog.Logger) ([]byte, error) {
	// Step 1: Generate RSA key pair
	privkey, err := GenerateKeyPair(RSAKeySize)
	if err != nil {
		logger.Error("KeyRequester: Unable to generate RSA key pair", "error", err)
		return nil, fmt.Errorf("unable to generate RSA key pair: %w", err)
	}
	pubkey := &privkey.PublicKey
	logger.Info("KeyRequester: Generated RSA key pair for key exchange")

	// Step 2: Serialize and encode the public key (needed for sending it over the network)
	serializedPubKey, err := SerializePublicKey(pubkey)
	if err != nil {
		logger.Error("KeyRequester: Failed to serialize public key", "error", err)
		return nil, fmt.Errorf("failed to serialize public key: %w", err)
	}

	// Step 4: Get the attestation report
	// Hash the serialized public key
	pubKeyHash := sha256.Sum256(serializedPubKey)
	attestationReport, err := GetReport(pubKeyHash[:])
	if err != nil {
		logger.Error("KeyRequester: Failed to get attestation report", "error", err)
		return nil, fmt.Errorf("failed to get attestation report: %w", err)
	}

	marshalledAttestation, err := json.Marshal(attestationReport)
	if err != nil {
		logger.Crit("unable to marshal attestation report", log.ErrKey, err)
		return nil, err
	}

	// Step 6: Create the message to send (PublicKey and Attestation)
	messageRequester := KeyExchangeRequest{
		PublicKey:   serializedPubKey,
		Attestation: marshalledAttestation,
	}

	// Step 7: Serialize the message to JSON for transmission
	messageBytesRequester, err := json.Marshal(messageRequester)
	if err != nil {
		logger.Error("KeyRequester: Failed to serialize message", "error", err)
		return nil, fmt.Errorf("failed to serialize message: %w", err)
	}

	// Step 8: Send the message to KeyProvider via HTTP POST
	resp, err := http.Post(config.EncryptionKeySource+"/v1"+common.PathKeyExchange, "application/json", bytes.NewBuffer(messageBytesRequester))
	if err != nil {
		logger.Error("KeyRequester: Failed to send message to KeyProvider", "error", err)
		return nil, fmt.Errorf("failed to send message to KeyProvider: %w", err)
	}
	defer resp.Body.Close()

	// Step 9: Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("KeyRequester: Failed to read response body from KeyProvider", "error", err)
		return nil, fmt.Errorf("failed to read response body from KeyProvider: %w", err)
	}

	// Check the HTTP response status
	if resp.StatusCode != http.StatusOK {
		logger.Error("KeyRequester: Received non-OK response from KeyProvider", "status", resp.Status, "body", string(bodyBytes))
		return nil, fmt.Errorf("received non-OK response from KeyProvider: %s", resp.Status)
	}

	// Step 10: Deserialize the received message
	var receivedMessageRequester KeyExchangeResponse
	err = json.Unmarshal(bodyBytes, &receivedMessageRequester)
	if err != nil {
		logger.Error("KeyRequester: Failed to deserialize received message", "error", err)
		return nil, fmt.Errorf("failed to deserialize received message: %w", err)
	}

	// Step 11: Extract and decode the encrypted encryption key from Base64
	encryptedKeyBytesRequester, err := DecodeBase64(receivedMessageRequester.EncryptedKey)
	if err != nil {
		logger.Error("KeyRequester: Failed to decode encrypted encryption key", "error", err)
		return nil, fmt.Errorf("failed to decode encrypted encryption key: %w", err)
	}

	// Step 12: Decrypt the encryption key using KeyRequester's private key
	decryptedKeyRequester, err := DecryptWithPrivateKey(encryptedKeyBytesRequester, privkey)
	if err != nil {
		logger.Error("KeyRequester: Decryption failed", "error", err)
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return decryptedKeyRequester, nil
}

// GetReport returns the attestation report for the given public key
func GetReport(pubKey []byte) (*tencommon.AttestationReport, error) {
	report, err := enclave.GetRemoteReport(pubKey)
	if err != nil {
		return nil, err
	}
	return &tencommon.AttestationReport{
		Report:      report,
		PubKey:      pubKey,
		EnclaveID:   gethcommon.Address{}, // this field is not needed for the key exchange
		HostAddress: "",                   // this field is not needed for the key exchange
	}, nil
}

// VerifyReport verifies the attestation report and returns the embedded data
func VerifyReport(att *tencommon.AttestationReport) ([]byte, error) {
	remoteReport, err := enclave.VerifyRemoteReport(att.Report)
	if err != nil {
		return []byte{}, err
	}
	return remoteReport.Data, nil
}

// GenerateKeyPair generates an RSA key pair of a given bit size.
func GenerateKeyPair(bits int) (*rsa.PrivateKey, error) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return privkey, nil
}

// SerializePublicKey serializes an RSA public key to DER-encoded bytes.
func SerializePublicKey(pubkey *rsa.PublicKey) ([]byte, error) {
	pubkeyBytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return nil, err
	}
	return pubkeyBytes, nil
}

// DeserializePublicKey deserializes DER-encoded bytes to an RSA public key.
func DeserializePublicKey(data []byte) (*rsa.PublicKey, error) {
	pubInterface, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, err
	}
	pubkey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}
	return pubkey, nil
}

// EncryptWithPublicKey encrypts data using RSA-OAEP and a public key.
func EncryptWithPublicKey(msg []byte, pubkey *rsa.PublicKey) ([]byte, error) {
	label := []byte("") // OAEP label is optional
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pubkey, msg, label)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// DecryptWithPrivateKey decrypts data using RSA-OAEP and a private key.
func DecryptWithPrivateKey(ciphertext []byte, privkey *rsa.PrivateKey) ([]byte, error) {
	label := []byte("") // OAEP label is optional
	hash := sha256.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, privkey, ciphertext, label)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// EncodeBase64 encodes data to a Base64 string.
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeBase64 decodes a Base64 string to data.
func DecodeBase64(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SerializeAttestationReport serializes an AttestationReport to JSON bytes.
func SerializeAttestationReport(report *tencommon.AttestationReport) ([]byte, error) {
	return json.Marshal(report)
}

// DeserializeAttestationReport deserializes JSON bytes to an AttestationReport.
func DeserializeAttestationReport(data []byte) (*tencommon.AttestationReport, error) {
	var report tencommon.AttestationReport
	err := json.Unmarshal(data, &report)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

// createHSMClient creates an Azure Managed HSM client using Managed Identity
func createHSMClient(hsmURL string) (*azkeys.Client, error) {
	if hsmURL == "" {
		return nil, fmt.Errorf("HSM URL is empty")
	}

	// Use DefaultAzureCredential which automatically uses AKS Managed Identity
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure credential: %w", err)
	}

	client, err := azkeys.NewClient(hsmURL, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HSM client: %w", err)
	}

	return client, nil
}

// backupKeyToHSM backs up the encryption key to Azure Managed HSM
func backupKeyToHSM(key []byte, config common.Config, logger gethlog.Logger) error {
	logger.Info("backing up encryption key to Azure HSM", "hsm_url", config.AzureHSMURL)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := createHSMClient(config.AzureHSMURL)
	if err != nil {
		return fmt.Errorf("failed to create HSM client: %w", err)
	}

	// Check if key already exists in HSM
	_, err = client.GetKey(ctx, AzureHSMKeyName, "", nil)
	if err == nil {
		logger.Warn("key already exists in HSM, skipping backup", "key_name", AzureHSMKeyName)
		return nil
	}

	// Define key operations
	encryptOp := azkeys.KeyOperationEncrypt
	decryptOp := azkeys.KeyOperationDecrypt
	wrapOp := azkeys.KeyOperationWrapKey
	unwrapOp := azkeys.KeyOperationUnwrapKey
	keyOps := []*azkeys.KeyOperation{
		&encryptOp,
		&decryptOp,
		&wrapOp,
		&unwrapOp,
	}

	// Import key to HSM
	keyType := azkeys.KeyTypeOctHSM
	enabled := true
	params := azkeys.ImportKeyParameters{
		Key: &azkeys.JSONWebKey{
			Kty:    &keyType,
			K:      key,
			KeyOps: keyOps,
		},
		KeyAttributes: &azkeys.KeyAttributes{
			Enabled: &enabled,
		},
	}

	_, err = client.ImportKey(ctx, AzureHSMKeyName, params, nil)
	if err != nil {
		return fmt.Errorf("failed to import key to HSM: %w", err)
	}

	logger.Info("successfully backed up encryption key to Azure HSM", "key_name", AzureHSMKeyName)
	return nil
}

// recoverKeyFromHSM recovers the encryption key from Azure Managed HSM
// Returns (key, found, error)
func recoverKeyFromHSM(config common.Config, logger gethlog.Logger) ([]byte, bool, error) {
	logger.Info("attempting to recover encryption key from Azure HSM", "hsm_url", config.AzureHSMURL)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := createHSMClient(config.AzureHSMURL)
	if err != nil {
		return nil, false, fmt.Errorf("failed to create HSM client: %w", err)
	}

	// Get the key from HSM
	resp, err := client.GetKey(ctx, AzureHSMKeyName, "", nil)
	if err != nil {
		// Key not found in HSM is not a fatal error during recovery
		logger.Warn("key not found in HSM", "key_name", AzureHSMKeyName, "error", err)
		return nil, false, nil
	}

	if resp.Key.K == nil || len(resp.Key.K) == 0 {
		return nil, false, fmt.Errorf("key material not available in HSM response")
	}

	// The key material is already in bytes format
	key := resp.Key.K

	logger.Info("successfully recovered encryption key from Azure HSM", "key_name", AzureHSMKeyName)
	return key, true, nil
}
