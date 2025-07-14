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
	"os"
	"path/filepath"
	"strings"
	"time"

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

// GetEncryptionKey returns the encryption key for the database
// - If we use an SQLite database, no encryption key is needed as SQLite typically runs in development or testing environments.
// - If an encryptionKeySource is provided, attempt to obtain the encryption key from the specified source.
// - If the encryptionKeySource is set to "new", check for an existing encryption key and generate a new one if not found.
// - If no encryptionKeySource is provided, attempt to unseal an existing encryption key.
// - If a new key is generated or obtained, seal it for future use.
func GetEncryptionKey(config common.Config, logger gethlog.Logger) ([]byte, error) {
	logger.Info("GetEncryptionKey: Starting encryption key acquisition process")
	logger.Info("GetEncryptionKey: Configuration", "dbType", config.DBType, "encryptionKeySource", config.EncryptionKeySource, "insideEnclave", config.InsideEnclave)
	
	// check if we are using sqlite database and no encryption key needed
	if config.DBType == "sqlite" {
		logger.Info("using sqlite database, no encryption key needed - exiting key exchange process")
		return nil, nil
	}

	// If no encryptionKeySource is provided, attempt to unseal an existing encryption key and fail if no key is found
	// (in this case operator needs to provide a source for the encryption key or decide to generate a new one)
	if config.EncryptionKeySource == "" {
		logger.Info("no key exchange url set, try to unseal existing encryption key", "keyFile", encryptionKeyFile)
		encryptionKey, found, err := tryUnsealKey(encryptionKeyFile, config.InsideEnclave)
		if !found {
			logger.Crit("no sealed encryption key found", log.ErrKey, err)
			logger.Error("GetEncryptionKey: No existing key found - operator must provide encryptionKeySource or use 'new' to generate")
			return nil, fmt.Errorf("no sealed encryption key found: %w", err)
		}
		logger.Info("unsealed existing encryption key", "keySize", len(encryptionKey))
		return encryptionKey, nil
	}

	var encryptionKey []byte
	var err error

	// If the "new" keyword is used for the encryptionKeySource, we first check if there is an existing encryption key
	// that can be unsealed and used. If no such key is found, we proceed to generate a new random encryption key.
	// This ensures that we do not overwrite an existing key unless necessary, and a new key is only generated when
	// there is no existing key available.
	if config.EncryptionKeySource == "new" {
		logger.Info("encryptionKeySource set to 'new' -> checking if there is an existing encryption key that we can use", "keyFile", encryptionKeyFile)
		var found bool
		encryptionKey, found, _ = tryUnsealKey(encryptionKeyFile, config.InsideEnclave)
		logger.Info("Encryption key status", "found", found, "error", err)
		if !found {
			logger.Info("No existing encryption key found, generating new random encryption key")
			encryptionKey, err = common.GenerateRandomKey()
			if err != nil {
				logger.Crit("unable to generate random encryption key", log.ErrKey, err)
				return nil, err
			}
			logger.Info("Generated new random encryption key", "keySize", len(encryptionKey))
		} else {
			logger.Info("Using existing encryption key", "keySize", len(encryptionKey))
		}
	} else {
		// Attempt to perform key exchange with the specified key provider.
		// This step is crucial, and the process should fail if the key exchange is not successful.
		logger.Info(fmt.Sprintf("encryptionKeySource set to '%s', trying to get encryption key from key provider", config.EncryptionKeySource))
		logger.Info("GetEncryptionKey: Initiating key exchange with remote provider")
		encryptionKey, err = HandleKeyExchange(config, logger)
		if err != nil {
			logger.Crit("unable to get encryption key from key provider", log.ErrKey, err)
			logger.Error("GetEncryptionKey: Key exchange failed - this is typically due to SGX attestation issues")
			return nil, err
		}
		logger.Info("Successfully obtained encryption key from key provider", "keySize", len(encryptionKey))
	}

	// Seal the key that we generated / got from the key exchange from another enclave
	logger.Info("GetEncryptionKey: Sealing encryption key for future use", "keyFile", encryptionKeyFile, "insideEnclave", config.InsideEnclave)
	err = trySealKey(encryptionKey, encryptionKeyFile, config.InsideEnclave)
	if err != nil {
		logger.Crit("unable to seal encryption key", log.ErrKey, err)
		logger.Error("GetEncryptionKey: Key sealing failed - future startups may require key re-acquisition")
		return nil, err
	}
	logger.Info("sealed new encryption key", "keySize", len(encryptionKey))
	logger.Info("GetEncryptionKey: Encryption key acquisition completed successfully")

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
	logger.Info("KeyRequester: Starting key exchange process", "keySource", config.EncryptionKeySource)
	
	// Step 1: Generate RSA key pair
	logger.Info("KeyRequester: Generating RSA key pair", "keySize", RSAKeySize)
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
	logger.Info("KeyRequester: Hashed public key for attestation", "hashLength", len(pubKeyHash), "hash", fmt.Sprintf("%x", pubKeyHash[:8]))
	
	logger.Info("KeyRequester: Attempting to get attestation report from SGX enclave")
	attestationReport, err := GetReport(pubKeyHash[:])
	if err != nil {
		logger.Error("KeyRequester: Failed to get attestation report", "error", err)
		logger.Error("KeyRequester: Attestation failure indicates SGX issues - check SGX setup, PCCS connectivity, and device availability")
		return nil, fmt.Errorf("failed to get attestation report: %w", err)
	}
	logger.Info("KeyRequester: Successfully obtained attestation report", "reportSize", len(attestationReport.Report))

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
	keyProviderURL := config.EncryptionKeySource + "/v1" + common.PathKeyExchange
	logger.Info("KeyRequester: Sending key exchange request to KeyProvider", "url", keyProviderURL, "messageSize", len(messageBytesRequester))
	
	resp, err := http.Post(keyProviderURL, "application/json", bytes.NewBuffer(messageBytesRequester))
	if err != nil {
		logger.Error("KeyRequester: Failed to send message to KeyProvider", "error", err, "url", keyProviderURL)
		logger.Error("KeyRequester: Network error - check KeyProvider availability and network connectivity")
		return nil, fmt.Errorf("failed to send message to KeyProvider: %w", err)
	}
	defer resp.Body.Close()
	
	logger.Info("KeyRequester: Received response from KeyProvider", "statusCode", resp.StatusCode, "status", resp.Status)

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
	logger.Info("KeyRequester: Decrypting encryption key with private key", "encryptedKeySize", len(encryptedKeyBytesRequester))
	decryptedKeyRequester, err := DecryptWithPrivateKey(encryptedKeyBytesRequester, privkey)
	if err != nil {
		logger.Error("KeyRequester: Decryption failed", "error", err)
		logger.Error("KeyRequester: Decryption failure indicates key mismatch or corrupted data")
		return nil, fmt.Errorf("decryption failed: %w", err)
	}
	
	logger.Info("KeyRequester: Successfully decrypted encryption key", "keySize", len(decryptedKeyRequester))
	logger.Info("KeyRequester: Key exchange process completed successfully")
	return decryptedKeyRequester, nil
}

// testPCCSConnectivity tests connectivity to the PCCS service
func testPCCSConnectivity() error {
	pccsURL := os.Getenv("PCCS_URL")
	if pccsURL == "" {
		pccsURL = "https://global.acccache.azure.net/sgx/certification/v4/"
	}
	
	gethlog.Info("PCCS Connectivity Test: Testing connection to PCCS service", "url", pccsURL)
	
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	// Try to connect to PCCS service
	req, err := http.NewRequest("GET", pccsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create PCCS request: %w", err)
	}
	
	resp, err := client.Do(req)
	if err != nil {
		gethlog.Error("PCCS Connectivity Test: Failed to connect to PCCS service", "error", err)
		gethlog.Error("PCCS Connectivity Test: This indicates network connectivity issues")
		gethlog.Error("PCCS Connectivity Test: Check if the container has outbound internet access")
		gethlog.Error("PCCS Connectivity Test: Verify DNS resolution and firewall rules")
		gethlog.Error("PCCS Connectivity Test: For Kubernetes, check network policies and egress rules")
		return fmt.Errorf("PCCS connectivity test failed: %w", err)
	}
	defer resp.Body.Close()
	
	gethlog.Info("PCCS Connectivity Test: Successfully connected to PCCS service", "statusCode", resp.StatusCode)
	
	// Read response body for additional debugging
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		gethlog.Warn("PCCS Connectivity Test: Failed to read response body", "error", err)
	} else {
		gethlog.Info("PCCS Connectivity Test: Response received", "bodySize", len(body))
	}
	
	return nil
}

// GetReport returns the attestation report for the given public key
func GetReport(pubKey []byte) (*tencommon.AttestationReport, error) {
	// Log SGX environment status before attempting attestation
	logSGXEnvironment()
	
	// Test PCCS connectivity before attempting SGX attestation
	gethlog.Info("GetReport: Testing PCCS connectivity before SGX attestation")
	if err := testPCCSConnectivity(); err != nil {
		gethlog.Error("GetReport: PCCS connectivity test failed - this will likely cause SGX attestation to fail")
		gethlog.Error("GetReport: Consider fixing network connectivity before proceeding")
		// Continue with attestation attempt to provide full error context
	} else {
		gethlog.Info("GetReport: PCCS connectivity test passed - proceeding with SGX attestation")
	}
	
	// Log the public key hash for debugging
	gethlog.Info("GetReport: Starting attestation report generation", "pubKeyHash", fmt.Sprintf("%x", pubKey))
	
	// Attempt to get the remote report from SGX enclave
	gethlog.Info("GetReport: Calling enclave.GetRemoteReport()")
	report, err := enclave.GetRemoteReport(pubKey)
	if err != nil {
		gethlog.Error("GetReport: Failed to get remote report from SGX enclave", "error", err)
		gethlog.Error("GetReport: This typically indicates SGX hardware/driver issues, PCCS connectivity problems, or missing SGX services")
		
		// Log additional SGX-specific error context
		logSGXErrorContext(err)
		return nil, err
	}
	
	gethlog.Info("GetReport: Successfully obtained remote report from SGX enclave", "reportSize", len(report))
	
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

// logSGXEnvironment logs detailed information about the SGX environment
func logSGXEnvironment() {
	gethlog.Info("SGX Environment Check: Starting detailed environment analysis")
	
	// Log environment variables
	envVars := []string{
		"OE_SIMULATION", "AESM_PATH", "PCCS_URL", "SGX_AESM_ADDR", 
		"SGX_SPID", "SGX_LINKABLE", "SGX_DEBUG", "SGX_MODE",
	}
	
	for _, envVar := range envVars {
		value := os.Getenv(envVar)
		if value != "" {
			gethlog.Info(fmt.Sprintf("SGX Environment: %s=%s", envVar, value))
		} else {
			gethlog.Info(fmt.Sprintf("SGX Environment: %s is not set", envVar))
		}
	}
	
	// Check SGX device files
	sgxDevices := []string{"/dev/sgx_enclave", "/dev/sgx_provision", "/dev/sgx/enclave", "/dev/sgx/provision"}
	for _, device := range sgxDevices {
		if _, err := os.Stat(device); err == nil {
			gethlog.Info(fmt.Sprintf("SGX Device: %s exists", device))
		} else {
			gethlog.Info(fmt.Sprintf("SGX Device: %s does not exist or is not accessible: %v", device, err))
		}
	}
	
	// Check AESM socket
	aesmSocket := os.Getenv("AESM_PATH")
	if aesmSocket == "" {
		aesmSocket = "/var/run/aesmd/aesm.socket"
	}
	if _, err := os.Stat(aesmSocket); err == nil {
		gethlog.Info(fmt.Sprintf("AESM Socket: %s exists", aesmSocket))
	} else {
		gethlog.Info(fmt.Sprintf("AESM Socket: %s does not exist or is not accessible: %v", aesmSocket, err))
	}
	
	// Check data directory permissions
	if stat, err := os.Stat(dataDir); err == nil {
		gethlog.Info(fmt.Sprintf("Data Directory: %s exists, mode: %v", dataDir, stat.Mode()))
	} else {
		gethlog.Info(fmt.Sprintf("Data Directory: %s does not exist or is not accessible: %v", dataDir, err))
	}
	
	// Log container environment indicators
	if _, err := os.Stat("/.dockerenv"); err == nil {
		gethlog.Info("Container Environment: Running inside Docker")
	}
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		gethlog.Info("Container Environment: Running inside Kubernetes")
	}
}

// logSGXErrorContext provides additional context for SGX-related errors
func logSGXErrorContext(err error) {
	errStr := err.Error()
	gethlog.Error("SGX Error Analysis: Detailed error context", "error", errStr)
	
	// Common SGX error patterns and their meanings
	errorPatterns := map[string]string{
		"OE_UNEXPECTED": "Unexpected error - usually indicates SGX hardware/driver issues or PCCS connectivity problems",
		"SGX_QL_NETWORK_ERROR": "Network error during quote generation - PCCS service unreachable",
		"Couldn't connect to server": "PCCS server connectivity issue - check network and PCCS_URL",
		"SGX_ERROR_DEVICE_BUSY": "SGX device is busy or in use by another process",
		"SGX_ERROR_OUT_OF_MEMORY": "Insufficient memory for SGX operations",
		"SGX_ERROR_INVALID_PARAMETER": "Invalid parameters passed to SGX functions",
		"SGX_ERROR_ENCLAVE_LOST": "Enclave has been destroyed or is not accessible",
		"0xe019": "Platform quote certificate data retrieval failed",
		"0xb006": "Quote provider library failed to get quote config",
	}
	
	for pattern, description := range errorPatterns {
		if strings.Contains(errStr, pattern) {
			gethlog.Error(fmt.Sprintf("SGX Error Pattern Detected: %s - %s", pattern, description))
		}
	}
	
	// Log troubleshooting suggestions
	if strings.Contains(errStr, "network") || strings.Contains(errStr, "connect") {
		gethlog.Info("SGX Troubleshooting: Network-related error detected")
		gethlog.Info("SGX Troubleshooting: 1. Check PCCS_URL environment variable")
		gethlog.Info("SGX Troubleshooting: 2. Verify network connectivity to PCCS service")
		gethlog.Info("SGX Troubleshooting: 3. Check firewall rules and DNS resolution")
		gethlog.Info("SGX Troubleshooting: 4. Verify Kubernetes network policies allow PCCS access")
	}
	
	if strings.Contains(errStr, "device") || strings.Contains(errStr, "/dev/sgx") {
		gethlog.Info("SGX Troubleshooting: Device-related error detected")
		gethlog.Info("SGX Troubleshooting: 1. Check if SGX devices are properly mounted in container")
		gethlog.Info("SGX Troubleshooting: 2. Verify SGX driver is loaded on host")
		gethlog.Info("SGX Troubleshooting: 3. Check container security context and privileged mode")
		gethlog.Info("SGX Troubleshooting: 4. Verify SGX hardware is enabled and available")
	}
}
