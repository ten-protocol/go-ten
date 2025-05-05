package keymanager

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/edgelesssys/ego/enclave"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	tencommon "github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave/core/egoutils"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
)

const (
	dataDir    = "/data"
	RSAKeySize = 2048
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

// GetEncryptionKey returns encryption key for the database
// 1.) If we use sqlite database we don't need to do anything since sqlite does not need encryption and is usually running in dev environments / for testing
// 2.) We need to check if the key is already sealed and unseal it if so
// 3.) If there is a URL to exchange the key with another enclave we need to get the key from there and seal it on this enclave
// 4.) If the key is not sealed and we don't have the URL to exchange the key we need to generate a new one and seal it
func GetEncryptionKey(config common.Config, logger gethlog.Logger) ([]byte, error) {
	// 1.) check if we are using sqlite database and no encryption key needed
	if config.DBType == "sqlite" {
		logger.Info("using sqlite database, no encryption key needed - exiting key exchange process")
		return nil, nil
	}

	var encryptionKey []byte

	// 2.) Check if we have a sealed encryption key and try to unseal it
	encryptionKey, found, err := tryUnsealKey(encryptionKeyFile, config.InsideEnclave)
	if err != nil {
		logger.Info("unable to unseal encryption key", log.ErrKey, err)
	}
	// If we found a sealed key we can return it
	if found {
		logger.Info("found sealed encryption key")
		return encryptionKey, nil
	}

	// 3.) We have to exchange the key with another enclave if we have a key exchange url set
	if config.KeyExchangeURL != "" {
		encryptionKey, err = HandleKeyExchange(config, logger)
		if err != nil {
			logger.Crit("unable to exchange key", log.ErrKey, err)
		} else {
			logger.Info("successfully exchanged key with another enclave")
		}
	}

	// 4.) If we don't have a key we need to generate a new one
	if len(encryptionKey) == 0 {
		encryptionKey, err = common.GenerateRandomKey()
		if err != nil {
			logger.Crit("unable to generate random encryption key", log.ErrKey, err)
			return nil, err
		} else {
			logger.Info("Successfully generated random encryption key")
		}
	}

	// Seal the key that we generated / got from the key exchange from another enclave
	err = trySealKey(encryptionKey, encryptionKeyFile, config.InsideEnclave)
	if err != nil {
		logger.Crit("unable to seal encryption key", log.ErrKey, err)
		return nil, err
	}
	logger.Info("sealed new encryption key")

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

// trySealKey attempts to seal an encryption key to disk
// Only seals if running in an SGX enclave
func trySealKey(key []byte, keyPath string, isEnclave bool) error {
	// Only attempt sealing if we're in an SGX enclave
	if !isEnclave {
		return nil
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
	resp, err := http.Post(config.KeyExchangeURL+"/v1"+common.PathKeyExchange, "application/json", bytes.NewBuffer(messageBytesRequester))
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
		return nil, fmt.Errorf("not RSA public key")
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
