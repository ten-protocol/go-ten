package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// SGX SIGSTRUCT format constants
	SigStructSize    = 1808 // Size of SIGSTRUCT in bytes
	RSASignatureSize = 384  // RSA-3072 signature size
	RSAModulusSize   = 384  // RSA-3072 modulus size
	RSAExponentSize  = 4    // RSA exponent size
	MREnclaveSize    = 32   // MRENCLAVE size
)

// SGX SIGSTRUCT magic bytes
var SigStructMagic = []byte{0x06, 0x00, 0x00, 0x00, 0xE1, 0x00, 0x00, 0x00}

// SIGSTRUCT field offsets (based on Intel SGX Programming Reference)
type SigStructOffsets struct {
	Magic         int // 8 bytes
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
	Vendor:        8,
	Date:          12,
	Header:        16,
	SwDefined:     32,
	Reserved1:     36,
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

// SigningData contains information extracted from the enclave for signing
type SigningData struct {
	SigStructOffset int    `json:"sigstruct_offset"`
	MREnclave       string `json:"mrenclave"`
	HashToSign      string `json:"hash_to_sign"`
	ModulusOffset   int    `json:"modulus_offset"`
	SignatureOffset int    `json:"signature_offset"`
	ExponentOffset  int    `json:"exponent_offset"`
}

// SGXSignatureReplacer handles SGX enclave signature operations
type SGXSignatureReplacer struct{}

// NewSGXSignatureReplacer creates a new signature replacer instance
func NewSGXSignatureReplacer() *SGXSignatureReplacer {
	return &SGXSignatureReplacer{}
}

// findSigStruct locates the SIGSTRUCT in binary data
func (r *SGXSignatureReplacer) findSigStruct(data []byte) (int, error) {
	// Look for the SGX SIGSTRUCT magic bytes
	for i := 0; i <= len(data)-len(SigStructMagic); i++ {
		if bytesEqual(data[i:i+len(SigStructMagic)], SigStructMagic) {
			// Verify we have enough data for a complete SIGSTRUCT
			if len(data) < i+SigStructSize {
				return -1, fmt.Errorf("incomplete SIGSTRUCT found at offset %d", i)
			}
			return i, nil
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

// ExtractSigningData extracts the data that needs to be signed from the enclave
func (r *SGXSignatureReplacer) ExtractSigningData(binaryPath string) (*SigningData, error) {
	data, err := os.ReadFile(binaryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read binary: %w", err)
	}

	sigStructOffset, err := r.findSigStruct(data)
	if err != nil {
		return nil, err
	}

	sigStruct := data[sigStructOffset : sigStructOffset+SigStructSize]

	// Extract MRENCLAVE (enclave hash)
	mrEnclave := sigStruct[offsets.EnclaveHash : offsets.EnclaveHash+MREnclaveSize]

	// Extract the data that gets signed (everything except the signature itself)
	// This includes header, modulus, exponent, and enclave measurement
	var signingData []byte

	// Before signature
	signingData = append(signingData, sigStruct[offsets.Magic:offsets.Signature]...)
	// After signature, before Q1
	signingData = append(signingData, sigStruct[offsets.MiscSelect:offsets.Q1]...)

	// Calculate hash of the signing data
	hash := sha256.Sum256(signingData)

	return &SigningData{
		SigStructOffset: sigStructOffset,
		MREnclave:       base64.StdEncoding.EncodeToString(mrEnclave),
		HashToSign:      base64.StdEncoding.EncodeToString(hash[:]),
		ModulusOffset:   sigStructOffset + offsets.Modulus,
		SignatureOffset: sigStructOffset + offsets.Signature,
		ExponentOffset:  sigStructOffset + offsets.Exponent,
	}, nil
}

// ReplaceSignature replaces the signature in the enclave binary
func (r *SGXSignatureReplacer) ReplaceSignature(binaryPath, signatureB64, outputPath string) error {
	// Decode the base64 signature
	signatureBytes, err := base64.StdEncoding.DecodeString(signatureB64)
	if err != nil {
		return fmt.Errorf("invalid base64 signature: %w", err)
	}

	// Ensure signature is the right size (pad or truncate if necessary)
	if len(signatureBytes) > RSASignatureSize {
		signatureBytes = signatureBytes[:RSASignatureSize]
	} else if len(signatureBytes) < RSASignatureSize {
		// Pad with zeros
		padded := make([]byte, RSASignatureSize)
		copy(padded, signatureBytes)
		signatureBytes = padded
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
	magic := sigStruct[offsets.Magic : offsets.Magic+8]
	if !bytesEqual(magic, SigStructMagic) {
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
	fmt.Printf("  %s extract <enclave_binary>\n", filepath.Base(os.Args[0]))
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
	case "extract":
		binaryPath := os.Args[2]
		data, err := replacer.ExtractSigningData(binaryPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(jsonData))

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
