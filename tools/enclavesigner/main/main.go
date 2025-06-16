package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
)

const (
	// SGX SIGSTRUCT format constants
	SigStructSize    = 1808 // Size of SIGSTRUCT in bytes
	RSASignatureSize = 384
)

// SGX SIGSTRUCT magic bytes
var (
	Header  = common.Hex2Bytes("06000000E10000000000010000000000")
	Header2 = common.Hex2Bytes("01010000600000006000000001000000")
	Vendor  = common.Hex2Bytes("00008086")
)

// SIGSTRUCT field offsets (based on Intel SGX Programming Reference)
type SigStructOffsets struct {
	Magic         int // 16 bytes
	Vendor        int // 4 bytes
	Date          int // 4 bytes
	Header        int // 16 bytes
	SwDefined     int // 4 bytes
	Reserved1     int // 84 bytes
	Modulus       int // 384 bytes (RSA-3072)
	Exponent      int // 4 bytes
	Signature     int // 384 bytes (RSA-3072)
	MiscSelect    int // 4 bytes
	MiscMask      int // 4 bytes
	Reserved2     int // 20 bytes
	Attributes    int // 16 bytes
	AttributeMask int // 16 bytes
	EnclaveHash   int // 32 bytes (MRENCLAVE)
	Reserved3     int // 32 bytes
	ISVProdID     int // 2 bytes
	ISVSVN        int // 2 bytes
	Reserved4     int // 12 bytes
	Q1            int // 384 bytes
	Q2            int // 384 bytes
}

var offsets = SigStructOffsets{
	Magic:         0,
	Vendor:        16,
	Date:          20,
	Header:        24,
	SwDefined:     40,
	Reserved1:     44,
	Modulus:       128,
	Exponent:      512,
	Signature:     516,
	MiscSelect:    900,
	MiscMask:      904,
	Reserved2:     908,
	Attributes:    928,
	AttributeMask: 944,
	EnclaveHash:   960,
	Reserved3:     992,
	ISVProdID:     1024,
	ISVSVN:        1026,
	Reserved4:     1028,
	Q1:            1040,
	Q2:            1424,
}

// SGXSignatureReplacer handles SGX enclave signature operations
type SGXSignatureReplacer struct{}

// NewSGXSignatureReplacer creates a new signature replacer instance
func NewSGXSignatureReplacer() *SGXSignatureReplacer {
	return &SGXSignatureReplacer{}
}

// findSigStruct locates the SIGSTRUCT in binary data
func (r *SGXSignatureReplacer) findSigStruct(data []byte) (int, error) {
	// Search in chunks to avoid memory issues with large binaries
	chunkSize := 1024 * 1024 // 1MB chunks
	magicLen := len(Header)

	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize + magicLen // Overlap to avoid missing magic at chunk boundary
		if end > len(data) {
			end = len(data)
		}

		chunk := data[i:end]
		for j := 0; j <= len(chunk)-magicLen; j++ {
			if bytesEqual(chunk[j:j+magicLen], Header) {
				absoluteOffset := i + j
				// Check HEADER field (offset 0, 16 bytes)
				if !bytesEqual(chunk[j:j+magicLen], Header) {
					continue
				}

				// Check HEADER2 field (offset 24, 16 bytes)
				if !bytesEqual(chunk[j+24:j+24+magicLen], Header2) {
					continue
				}

				// All validation passed - verify we have enough data for complete SIGSTRUCT
				if len(data) < absoluteOffset+SigStructSize {
					return -1, fmt.Errorf("incomplete SIGSTRUCT found at offset %d", absoluteOffset)
				}

				return absoluteOffset, nil
			}
		}
	}

	return -1, fmt.Errorf("SIGSTRUCT not found in binary")
}

// bytesEqual compares two byte slices
func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ExtractOriginalHash extracts the hash that was originally signed
func (r *SGXSignatureReplacer) ExtractOriginalHash(binaryPath, keyPath string) (string, error) {
	// Read the enclave binary
	data, err := os.ReadFile(binaryPath)
	if err != nil {
		return "", fmt.Errorf("failed to read binary: %w", err)
	}

	// Find SIGSTRUCT
	sigStructOffset, err := r.findSigStruct(data)
	if err != nil {
		return "", err
	}

	sigStruct := data[sigStructOffset : sigStructOffset+SigStructSize]

	vendor := sigStruct[offsets.Vendor:offsets.Date]
	fmt.Fprintf(os.Stderr, "Vendor: %x\n", vendor)

	date := sigStruct[offsets.Date:offsets.Header]
	fmt.Fprintf(os.Stderr, "Date: %x\n", date)

	header2 := sigStruct[offsets.Header:offsets.SwDefined]
	if !bytes.Equal(header2, Header2) {
		return "", fmt.Errorf("invalid header 2")
	}

	// Extract signature
	// signature := sigStruct[offsets.Signature:offsets.MiscSelect]

	// Read and parse the public key
	//publicKey, err := r.loadPublicKey(keyPath)
	//if err != nil {
	//	return "", fmt.Errorf("failed to load public key: %w", err)
	//}

	// Verify signature and extract hash
	//hash, err := r.verifyAndExtractHash(signature, publicKey)
	//if err != nil {
	//	return "", fmt.Errorf("failed to extract hash: %w", err)
	//}

	hashExtr := sigStruct[offsets.EnclaveHash:offsets.Reserved3]
	//if !bytes.Equal(hashExtr, hash) {
	//	return "", fmt.Errorf("hash mismatch")
	//}

	return base64.StdEncoding.EncodeToString(hashExtr), nil
}

// loadPublicKey loads an RSA public key from a PEM file
func (r *SGXSignatureReplacer) loadPublicKey(keyPath string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	var publicKey *rsa.PublicKey

	switch block.Type {
	case "PUBLIC KEY":
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		var ok bool
		publicKey, ok = pub.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("not an RSA public key")
		}
	case "RSA PUBLIC KEY":
		pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		publicKey = pub
	case "RSA PRIVATE KEY":
		// Extract public key from private key
		priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		publicKey = &priv.PublicKey
	default:
		return nil, fmt.Errorf("unsupported key type: %s", block.Type)
	}

	return publicKey, nil
}

// verifyAndExtractHash performs RSA verification to extract the original hash
func (r *SGXSignatureReplacer) verifyAndExtractHash(signature []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	// Convert signature to big integer (reverse byte order for little-endian)
	sigBig := new(big.Int).SetBytes(signature)

	// Perform RSA verification: signature^e mod n
	decrypted := new(big.Int).Exp(sigBig, big.NewInt(int64(publicKey.E)), publicKey.N)

	// Convert back to bytes
	decryptedBytes := decrypted.Bytes()

	// Ensure proper length (pad with leading zeros if necessary)
	keySize := (publicKey.N.BitLen() + 7) / 8
	if len(decryptedBytes) < keySize {
		padded := make([]byte, keySize)
		copy(padded[keySize-len(decryptedBytes):], decryptedBytes)
		decryptedBytes = padded
	}

	// Debug output
	fmt.Fprintf(os.Stderr, "Debug: Signature length: %d\n", len(signature))
	fmt.Fprintf(os.Stderr, "Debug: Key size: %d\n", keySize)
	fmt.Fprintf(os.Stderr, "Debug: Decrypted length: %d\n", len(decryptedBytes))
	fmt.Fprintf(os.Stderr, "Debug: First 32 bytes: %x\n", decryptedBytes[:min(32, len(decryptedBytes))])
	fmt.Fprintf(os.Stderr, "Debug: Last 32 bytes: %x\n", decryptedBytes[max(0, len(decryptedBytes)-32):])

	// Try to remove PKCS#1 v1.5 padding
	hash, err := r.removePKCS1v15Padding(decryptedBytes)
	if err != nil {
		// If standard padding fails, try alternative approaches
		fmt.Fprintf(os.Stderr, "Debug: Standard padding failed, trying alternatives\n")

		// Check if the signature might be in a different format
		// Sometimes the hash is at the end without proper PKCS#1 padding
		if len(decryptedBytes) >= 32 {
			// Try last 32 bytes (SHA-256 hash size)
			possibleHash := decryptedBytes[len(decryptedBytes)-32:]
			fmt.Fprintf(os.Stderr, "Debug: Trying last 32 bytes as hash: %x\n", possibleHash)
			return possibleHash, nil
		}

		return nil, fmt.Errorf("padding removal failed: %w", err)
	}

	return hash, nil
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// removePKCS1v15Padding removes PKCS#1 v1.5 padding from decrypted signature
func (r *SGXSignatureReplacer) removePKCS1v15Padding(data []byte) ([]byte, error) {
	if len(data) < 11 {
		return nil, fmt.Errorf("invalid padding length: %d", len(data))
	}

	fmt.Fprintf(os.Stderr, "Debug: Padding check - first 4 bytes: %x\n", data[:4])

	// PKCS#1 v1.5 padding format: 0x00 0x01 0xFF...0xFF 0x00 <hash>
	if data[0] != 0x00 || data[1] != 0x01 {
		// Try alternative padding patterns or formats
		if data[0] == 0x01 {
			// Maybe missing leading zero
			fmt.Fprintf(os.Stderr, "Debug: Trying alternative padding (missing leading zero)\n")
			return r.removePKCS1v15PaddingAlt(data[1:])
		}
		return nil, fmt.Errorf("invalid padding header: %02x %02x", data[0], data[1])
	}

	// Find the 0x00 separator
	var separatorIndex int = -1
	for i := 2; i < len(data); i++ {
		if data[i] == 0x00 {
			separatorIndex = i
			break
		}
		if data[i] != 0xFF {
			fmt.Fprintf(os.Stderr, "Debug: Invalid padding byte at position %d: %02x\n", i, data[i])
			return nil, fmt.Errorf("invalid padding byte at position %d: %02x", i, data[i])
		}
	}

	if separatorIndex == -1 {
		return nil, fmt.Errorf("padding separator not found")
	}

	if separatorIndex < 10 {
		return nil, fmt.Errorf("padding too short: separator at %d", separatorIndex)
	}

	fmt.Fprintf(os.Stderr, "Debug: Found separator at position %d\n", separatorIndex)
	fmt.Fprintf(os.Stderr, "Debug: Hash length: %d\n", len(data)-separatorIndex-1)

	// Return the hash part
	return data[separatorIndex+1:], nil
}

// Alternative padding removal for cases where format is slightly different
func (r *SGXSignatureReplacer) removePKCS1v15PaddingAlt(data []byte) ([]byte, error) {
	if len(data) < 10 {
		return nil, fmt.Errorf("invalid alternative padding length: %d", len(data))
	}

	// Look for 0x01 followed by 0xFF bytes and then 0x00
	if data[0] != 0x01 {
		return nil, fmt.Errorf("invalid alternative padding header: %02x", data[0])
	}

	var separatorIndex int = -1
	for i := 1; i < len(data); i++ {
		if data[i] == 0x00 {
			separatorIndex = i
			break
		}
		if data[i] != 0xFF {
			return nil, fmt.Errorf("invalid alternative padding byte at position %d: %02x", i, data[i])
		}
	}

	if separatorIndex == -1 || separatorIndex < 9 {
		return nil, fmt.Errorf("alternative padding separator not found or too short")
	}

	return data[separatorIndex+1:], nil
}

// ReplaceSignature replaces the signature in the enclave binary
func (r *SGXSignatureReplacer) ReplaceSignature(binaryPath, signatureB64, outputPath string) error {
	// Decode the base64 signature
	signatureBytes, err := base64.StdEncoding.DecodeString(signatureB64)
	if err != nil {
		return fmt.Errorf("invalid base64 signature: %w", err)
	}

	// Ensure signature is the right size (pad or truncate if necessary)
	if len(signatureBytes) != RSASignatureSize {
		return fmt.Errorf("invalid signature size %s: %d", signatureB64, len(signatureBytes))
	}

	// Read the original binary
	data, err := os.ReadFile(binaryPath)
	if err != nil {
		return fmt.Errorf("failed to read binary: %w", err)
	}

	// Find SIGSTRUCT and replace signature
	sigStructOffset, err := r.findSigStruct(data)
	if err != nil {
		return err
	}

	signatureOffset := sigStructOffset + offsets.Signature

	// Replace the signature
	copy(data[signatureOffset:signatureOffset+RSASignatureSize], signatureBytes)

	// Write the modified binary
	err = os.WriteFile(outputPath, data, 0o755)
	if err != nil {
		return fmt.Errorf("failed to write output binary: %w", err)
	}

	fmt.Printf("Signature replaced in %s\n", outputPath)
	return nil
}

// VerifyStructure verifies the SIGSTRUCT is valid
func (r *SGXSignatureReplacer) VerifyStructure(binaryPath string) error {
	data, err := os.ReadFile(binaryPath)
	if err != nil {
		return fmt.Errorf("failed to read binary: %w", err)
	}

	sigStructOffset, err := r.findSigStruct(data)
	if err != nil {
		return err
	}

	sigStruct := data[sigStructOffset : sigStructOffset+SigStructSize]

	// Basic structure verification
	magic := sigStruct[offsets.Magic:offsets.Vendor]
	if !bytesEqual(magic, Header) {
		return fmt.Errorf("invalid SIGSTRUCT magic")
	}

	// Check if we have non-zero signature
	signature := sigStruct[offsets.Signature : offsets.Signature+RSASignatureSize]
	allZeros := true
	for _, b := range signature {
		if b != 0 {
			allZeros = false
			break
		}
	}
	if allZeros {
		fmt.Println("Warning: Signature appears to be all zeros")
	}

	// Extract some key information
	vendor := binary.LittleEndian.Uint32(sigStruct[offsets.Vendor : offsets.Vendor+4])
	date := binary.LittleEndian.Uint32(sigStruct[offsets.Date : offsets.Date+4])
	isvProdID := binary.LittleEndian.Uint16(sigStruct[offsets.ISVProdID : offsets.ISVProdID+2])
	isvSVN := binary.LittleEndian.Uint16(sigStruct[offsets.ISVSVN : offsets.ISVSVN+2])

	fmt.Printf("SIGSTRUCT found at offset %d\n", sigStructOffset)
	fmt.Printf("Magic: %x\n", magic)
	fmt.Printf("Vendor: %d\n", vendor)
	fmt.Printf("Date: %d\n", date)
	fmt.Printf("ISV Prod ID: %d\n", isvProdID)
	fmt.Printf("ISV SVN: %d\n", isvSVN)
	fmt.Printf("Signature (first 32 bytes): %x\n", signature[:32])
	fmt.Printf("SIGSTRUCT structure appears valid\n")

	return nil
}

// printUsage prints the usage information
func printUsage() {
	fmt.Printf("Usage:\n")
	fmt.Printf("  %s extract_hash <enclave_binary> <public_key_file>\n", filepath.Base(os.Args[0]))
	fmt.Printf("  %s replace <enclave_binary> <signature_file> <output_binary>\n", filepath.Base(os.Args[0]))
	fmt.Printf("  %s verify <enclave_binary>\n", filepath.Base(os.Args[0]))
}

func main() {
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	replacer := NewSGXSignatureReplacer()
	command := os.Args[1]

	switch command {
	case "extract_hash":
		if len(os.Args) != 4 {
			fmt.Printf("Usage: %s extract_hash <enclave_binary> <public_key_file>\n", filepath.Base(os.Args[0]))
			os.Exit(1)
		}

		binaryPath := os.Args[2]
		keyPath := os.Args[3]

		hash, err := replacer.ExtractOriginalHash(binaryPath, keyPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "Hash: %s\n", hash)

		fmt.Println(hash)

	case "replace":
		if len(os.Args) != 5 {
			fmt.Printf("Usage: %s replace <input_binary> <signature_file> <output_binary>\n", filepath.Base(os.Args[0]))
			os.Exit(1)
		}

		binaryPath := os.Args[2]
		signatureFile := os.Args[3]
		outputPath := os.Args[4]

		signatureData, err := os.ReadFile(signatureFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading signature file: %v\n", err)
			os.Exit(1)
		}

		signatureB64 := string(signatureData)
		// Remove any trailing whitespace
		signatureB64 = string(bytes.TrimSpace([]byte(signatureB64)))

		err = replacer.ReplaceSignature(binaryPath, signatureB64, outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "verify":
		binaryPath := os.Args[2]
		err := replacer.VerifyStructure(binaryPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}
