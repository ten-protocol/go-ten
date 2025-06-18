package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/big"
	"os"
)

const (
	SigStructSize = 1808 // offs.SGX_ARCH_SIGSTRUCT_SIZE
	RSAKeySize    = 384  // 3072 bits / 8
)

// SGX SIGSTRUCT field offsets based on Gramine's sigstruct.py
var offsets = struct {
	Header        int
	Vendor        int
	Date          int
	Header2       int
	SwDefined     int
	Modulus       int
	Exponent      int
	Signature     int
	MiscSelect    int
	MiscMask      int
	Attributes    int
	AttributeMask int
	EnclaveHash   int
	IsvProdID     int
	IsvSvn        int
	Q1            int
	Q2            int
}{
	Header:        0,    // SGX_ARCH_SIGSTRUCT_HEADER
	Vendor:        16,   // SGX_ARCH_SIGSTRUCT_VENDOR
	Date:          20,   // SGX_ARCH_SIGSTRUCT_DATE
	Header2:       24,   // SGX_ARCH_SIGSTRUCT_HEADER2
	SwDefined:     40,   // SGX_ARCH_SIGSTRUCT_SWDEFINED
	Modulus:       128,  // SGX_ARCH_SIGSTRUCT_MODULUS
	Exponent:      512,  // SGX_ARCH_SIGSTRUCT_EXPONENT
	Signature:     516,  // SGX_ARCH_SIGSTRUCT_SIGNATURE
	MiscSelect:    900,  // SGX_ARCH_SIGSTRUCT_MISC_SELECT (after_sig_offset)
	MiscMask:      904,  // SGX_ARCH_SIGSTRUCT_MISC_MASK
	Attributes:    928,  // SGX_ARCH_SIGSTRUCT_ATTRIBUTES
	AttributeMask: 944,  // SGX_ARCH_SIGSTRUCT_ATTRIBUTE_MASK
	EnclaveHash:   960,  // SGX_ARCH_SIGSTRUCT_ENCLAVE_HASH
	IsvProdID:     992,  // SGX_ARCH_SIGSTRUCT_ISV_PROD_ID
	IsvSvn:        994,  // SGX_ARCH_SIGSTRUCT_ISV_SVN
	Q1:            1040, // SGX_ARCH_SIGSTRUCT_Q1
	Q2:            1424, // SGX_ARCH_SIGSTRUCT_Q2
}

type SGXSignatureReplacer struct{}

func NewSGXSignatureReplacer() *SGXSignatureReplacer {
	return &SGXSignatureReplacer{}
}

// ExtractHash extracts the hash that needs to be signed from a SIGSTRUCT
// This follows Gramine's get_signing_data() method exactly
func (r *SGXSignatureReplacer) ExtractHash(sigStruct []byte) ([32]byte, error) {
	if len(sigStruct) < SigStructSize {
		return [32]byte{}, fmt.Errorf("invalid SIGSTRUCT size: %d", len(sigStruct))
	}

	// Gramine's get_signing_data() method:
	// return data[:128] + data[after_sig_offset:after_sig_offset+128]
	afterSigOffset := offsets.MiscSelect // 900

	// Extract the two 128-byte chunks as per Gramine
	signingData := make([]byte, 256) // 128 + 128
	copy(signingData[:128], sigStruct[:128])
	copy(signingData[128:], sigStruct[afterSigOffset:afterSigOffset+128])

	// Hash the signing data
	hash := sha256.Sum256(signingData)
	return hash, nil
}

// ReplaceSignature replaces the signature and modulus in a SIGSTRUCT
func (r *SGXSignatureReplacer) ReplaceSignature(sigStruct []byte, signatureB64, modulusB64 string) error {
	if len(sigStruct) < SigStructSize {
		return fmt.Errorf("invalid SIGSTRUCT size: %d", len(sigStruct))
	}

	// Decode base64url signature and modulus (Azure format)
	signatureBytesBE, err := base64.StdEncoding.DecodeString(signatureB64)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	modulusbytesBE, err := base64.StdEncoding.DecodeString(modulusB64)
	if err != nil {
		return fmt.Errorf("failed to decode modulus: %w", err)
	}

	fmt.Printf("Debug: Decoded signature length: %d bytes\n", len(signatureBytesBE))
	fmt.Printf("Debug: Decoded modulus length: %d bytes\n", len(modulusbytesBE))
	fmt.Printf("Debug: Expected RSA key size: %d bytes\n", RSAKeySize)

	// Pad signature to full RSA key size if needed (Azure may strip leading zeros)
	if len(signatureBytesBE) < RSAKeySize {
		paddedSig := make([]byte, RSAKeySize)
		copy(paddedSig[RSAKeySize-len(signatureBytesBE):], signatureBytesBE)
		signatureBytesBE = paddedSig
		fmt.Printf("Debug: Padded signature from %d to %d bytes\n", len(signatureBytesBE), RSAKeySize)
	}

	// Pad modulus to full RSA key size if needed
	if len(modulusbytesBE) < RSAKeySize {
		paddedMod := make([]byte, RSAKeySize)
		copy(paddedMod[RSAKeySize-len(modulusbytesBE):], modulusbytesBE)
		modulusbytesBE = paddedMod
		fmt.Printf("Debug: Padded modulus from %d to %d bytes\n", len(modulusbytesBE), RSAKeySize)
	}

	// Validate sizes after padding
	if len(signatureBytesBE) != RSAKeySize {
		return fmt.Errorf("signature size mismatch: expected %d, got %d", RSAKeySize, len(signatureBytesBE))
	}

	if len(modulusbytesBE) != RSAKeySize {
		return fmt.Errorf("modulus size mismatch: expected %d, got %d", RSAKeySize, len(modulusbytesBE))
	}

	// Convert from Azure's big-endian format to SGX's little-endian format
	signatureLE := bigEndianToLittleEndian(signatureBytesBE)
	modulusLE := bigEndianToLittleEndian(modulusbytesBE)

	// Update modulus and signature in SIGSTRUCT (both in little-endian format)
	copy(sigStruct[offsets.Modulus:offsets.Modulus+RSAKeySize], modulusLE)
	copy(sigStruct[offsets.Signature:offsets.Signature+RSAKeySize], signatureLE)

	// Set exponent to 3 (little-endian) - SGX always uses exponent 3
	binary.LittleEndian.PutUint32(sigStruct[offsets.Exponent:offsets.Exponent+4], 3)

	// Compute Q1 and Q2 after updating modulus and signature
	q1, q2 := r.computeQ1Q2(sigStruct)
	q1Bytes := bigIntToLittleEndian(q1, RSAKeySize)
	q2Bytes := bigIntToLittleEndian(q2, RSAKeySize)

	// Update Q1 and Q2 in SIGSTRUCT
	copy(sigStruct[offsets.Q1:offsets.Q1+RSAKeySize], q1Bytes)
	copy(sigStruct[offsets.Q2:offsets.Q2+RSAKeySize], q2Bytes)

	fmt.Printf("Debug: Successfully replaced signature and computed Q1/Q2\n")
	return nil
}

// VerifySignature verifies the SIGSTRUCT signature using SGX-style verification
func (r *SGXSignatureReplacer) VerifySignature(sigStruct []byte) error {
	if len(sigStruct) < SigStructSize {
		return fmt.Errorf("invalid SIGSTRUCT size: %d", len(sigStruct))
	}

	// Extract modulus and signature from SIGSTRUCT (both in little-endian)
	modulusLE := sigStruct[offsets.Modulus : offsets.Modulus+RSAKeySize]
	signatureLE := sigStruct[offsets.Signature : offsets.Signature+RSAKeySize]

	// Convert to big integers for RSA operations
	modulus := littleEndianToBigInt(modulusLE)
	signature := littleEndianToBigInt(signatureLE)

	// Get the hash that should have been signed
	expectedHash, err := r.ExtractHash(sigStruct)
	if err != nil {
		return fmt.Errorf("failed to extract expected hash: %w", err)
	}

	fmt.Printf("Debug: Expected hash: %x\n", expectedHash)
	fmt.Printf("Debug: Modulus size: %d bits\n", modulus.BitLen())
	fmt.Printf("Debug: Signature size: %d bits\n", signature.BitLen())

	// Try different verification approaches
	fmt.Printf("Debug: Attempting PKCS1v15 verification...\n")
	pubKey := &rsa.PublicKey{N: modulus, E: 3}
	signatureBytes := signature.Bytes()

	// Pad signature to RSA key size if needed
	if len(signatureBytes) < RSAKeySize {
		paddedSig := make([]byte, RSAKeySize)
		copy(paddedSig[RSAKeySize-len(signatureBytes):], signatureBytes)
		signatureBytes = paddedSig
	}

	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, expectedHash[:], signatureBytes)
	if err == nil {
		fmt.Printf("✓ PKCS1v15 verification successful!\n")
		return nil
	}
	fmt.Printf("Debug: PKCS1v15 failed: %v\n", err)

	// Try raw RSA verification
	fmt.Printf("Debug: Attempting raw RSA verification...\n")
	e := big.NewInt(3)
	decrypted := new(big.Int).Exp(signature, e, modulus)
	decryptedBytes := decrypted.Bytes()

	fmt.Printf("Debug: Raw RSA result: %d bytes\n", len(decryptedBytes))
	if len(decryptedBytes) >= 64 {
		fmt.Printf("Debug: Raw RSA (first 32 bytes): %x\n", decryptedBytes[:32])
		fmt.Printf("Debug: Raw RSA (last 32 bytes): %x\n", decryptedBytes[len(decryptedBytes)-32:])
	}

	// Check if our hash appears anywhere in the decrypted result
	hashBytes := expectedHash[:]
	for i := 0; i <= len(decryptedBytes)-len(hashBytes); i++ {
		if bytesEqual(decryptedBytes[i:i+len(hashBytes)], hashBytes) {
			fmt.Printf("✓ Found matching hash at offset %d in decrypted signature!\n", i)
			return nil
		}
	}

	// Try with different padding schemes
	fmt.Printf("Debug: Checking for standard RSA padding patterns...\n")

	// Look for PKCS#1 v1.5 padding pattern: 00 01 FF...FF 00 <hash>
	if len(decryptedBytes) >= 2 && decryptedBytes[0] == 0x00 && decryptedBytes[1] == 0x01 {
		fmt.Printf("Debug: Found PKCS#1 v1.5 padding pattern\n")
		// Find the 0x00 separator after the FF padding
		for i := 2; i < len(decryptedBytes)-len(hashBytes); i++ {
			if decryptedBytes[i] == 0x00 {
				hashStart := i + 1
				if hashStart+len(hashBytes) <= len(decryptedBytes) {
					actualHash := decryptedBytes[hashStart : hashStart+len(hashBytes)]
					if bytesEqual(actualHash, hashBytes) {
						fmt.Printf("✓ PKCS#1 v1.5 hash verification successful!\n")
						return nil
					}
				}
				break
			}
		}
	}

	return fmt.Errorf("signature verification failed: no matching hash found in decrypted signature")
}

// computeQ1Q2 computes Q1 and Q2 fields following Gramine's RSA math exactly
func (r *SGXSignatureReplacer) computeQ1Q2(sigStruct []byte) (*big.Int, *big.Int) {
	// Extract signature and modulus as big integers (from little-endian format)
	signatureLE := sigStruct[offsets.Signature : offsets.Signature+RSAKeySize]
	modulusLE := sigStruct[offsets.Modulus : offsets.Modulus+RSAKeySize]

	signature := littleEndianToBigInt(signatureLE)
	modulus := littleEndianToBigInt(modulusLE)

	fmt.Printf("Debug: Computing Q1/Q2 with signature bits: %d, modulus bits: %d\n",
		signature.BitLen(), modulus.BitLen())

	// Gramine's exact computation:
	// tmp1 = signature_int * signature_int  (signature squared)
	// q1_int = tmp1 // modulus_int          (integer division)
	// tmp2 = tmp1 % modulus_int             (remainder)
	// q2_int = tmp2 * signature_int // modulus_int

	tmp1 := new(big.Int).Mul(signature, signature)
	q1 := new(big.Int).Div(tmp1, modulus)
	tmp2 := new(big.Int).Mod(tmp1, modulus)
	tmp3 := new(big.Int).Mul(tmp2, signature)
	q2 := new(big.Int).Div(tmp3, modulus)

	fmt.Printf("Debug: Q1 bits: %d, Q2 bits: %d\n", q1.BitLen(), q2.BitLen())

	// Validate that Q1 and Q2 are within expected ranges
	// They should be smaller than the modulus
	if q1.Cmp(modulus) >= 0 {
		fmt.Printf("Warning: Q1 >= modulus (Q1 bits: %d, modulus bits: %d)\n", q1.BitLen(), modulus.BitLen())
	}
	if q2.Cmp(modulus) >= 0 {
		fmt.Printf("Warning: Q2 >= modulus (Q2 bits: %d, modulus bits: %d)\n", q2.BitLen(), modulus.BitLen())
	}

	return q1, q2
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <command> [args...]\n", os.Args[0])
		fmt.Println("Commands:")
		fmt.Println("  extract_hash <enclave_file>")
		fmt.Println("  replace <enclave_file> <signature_b64> <modulus_b64> <output_file>")
		fmt.Println("  verify <enclave_file>")
		os.Exit(1)
	}

	command := os.Args[1]
	replacer := NewSGXSignatureReplacer()

	switch command {
	case "extract_hash":
		if len(os.Args) != 3 {
			fmt.Println("Usage: extract_hash <enclave_file>")
			os.Exit(1)
		}

		data, err := os.ReadFile(os.Args[2])
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}

		sigStruct, err := extractSigStruct(data)
		if err != nil {
			fmt.Printf("Error finding SIGSTRUCT: %v\n", err)
			os.Exit(1)
		}

		//err = replacer.VerifySignature(sigStruct)
		//if err != nil {
		//	fmt.Printf("Error verifying original sig: %v\n", err)
		//	os.Exit(1)
		//}

		hash, err := replacer.ExtractHash(sigStruct)
		if err != nil {
			fmt.Printf("Error extracting hash: %v\n", err)
			os.Exit(1)
		}

		hashB64 := base64.StdEncoding.EncodeToString(hash[:])
		fmt.Printf("%s", hashB64) // Print without newline for script compatibility

	case "replace":
		if len(os.Args) != 6 {
			fmt.Println("Usage: replace <enclave_file> <signature_b64> <modulus_b64> <output_file>")
			os.Exit(1)
		}

		data, err := os.ReadFile(os.Args[2])
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}

		offset, err := findSigStructOffset(data)
		if err != nil {
			fmt.Printf("Error finding SIGSTRUCT: %v\n", err)
			os.Exit(1)
		}

		sigStruct := data[offset : offset+SigStructSize]
		err = replacer.ReplaceSignature(sigStruct, os.Args[3], os.Args[4])
		if err != nil {
			fmt.Printf("Error replacing signature: %v\n", err)
			os.Exit(1)
		}

		err = os.WriteFile(os.Args[5], data, 0644)
		if err != nil {
			fmt.Printf("Error writing file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Signature replaced successfully (SIGSTRUCT found at offset %d)\n", offset)

	case "verify":
		if len(os.Args) != 3 {
			fmt.Println("Usage: verify <enclave_file>")
			os.Exit(1)
		}

		data, err := os.ReadFile(os.Args[2])
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}

		sigStruct, err := extractSigStruct(data)
		if err != nil {
			fmt.Printf("Error finding SIGSTRUCT: %v\n", err)
			os.Exit(1)
		}

		err = replacer.VerifySignature(sigStruct)
		if err != nil {
			fmt.Printf("Verification failed: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

// findSigStructOffset searches for the SIGSTRUCT within the binary data
// SIGSTRUCT has a specific header signature: 0x06, 0x00, 0x00, 0x00, 0xE1, 0x00, 0x00, 0x00
func findSigStructOffset(data []byte) (int, error) {
	// SGX SIGSTRUCT header signature (little-endian)
	// Header1: 0x00000006, Header2: 0x000000E1
	sigstructHeader := []byte{0x06, 0x00, 0x00, 0x00, 0xE1, 0x00, 0x00, 0x00}

	// Search for the signature in the file
	for i := 0; i <= len(data)-len(sigstructHeader); i++ {
		if bytesEqual(data[i:i+len(sigstructHeader)], sigstructHeader) {
			// Verify we have enough space for the full SIGSTRUCT
			if i+SigStructSize > len(data) {
				continue // Not enough space, keep searching
			}

			// Additional validation: check if this looks like a real SIGSTRUCT
			// Verify the second header field at offset 24
			if i+24+8 <= len(data) {
				header2 := data[i+24 : i+24+8]
				expectedHeader2 := []byte{0x01, 0x01, 0x00, 0x00, 0x60, 0x00, 0x00, 0x00}
				if bytesEqual(header2, expectedHeader2) {
					return i, nil
				}
			}
		}
	}

	return -1, fmt.Errorf("SIGSTRUCT not found in binary")
}

// extractSigStruct extracts the SIGSTRUCT from binary data
func extractSigStruct(data []byte) ([]byte, error) {
	offset, err := findSigStructOffset(data)
	if err != nil {
		return nil, err
	}

	if offset+SigStructSize > len(data) {
		return nil, fmt.Errorf("not enough data for complete SIGSTRUCT at offset %d", offset)
	}

	return data[offset : offset+SigStructSize], nil
}

// Helper function to compare byte slices
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

// Helper function to convert big-endian bytes to little-endian bytes
func bigEndianToLittleEndian(data []byte) []byte {
	result := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		result[i] = data[len(data)-1-i]
	}
	return result
}

// Helper function to convert little-endian bytes to big.Int
func littleEndianToBigInt(data []byte) *big.Int {
	// Reverse the bytes to convert from little-endian to big-endian
	reversed := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		reversed[i] = data[len(data)-1-i]
	}
	return new(big.Int).SetBytes(reversed)
}

// Helper function to convert big.Int to little-endian bytes with proper padding
func bigIntToLittleEndian(val *big.Int, size int) []byte {
	// Get the big-endian bytes from the big.Int
	bigEndianBytes := val.Bytes()

	// Create result buffer filled with zeros (for left-padding)
	result := make([]byte, size)

	// If the big.Int is larger than our target size, truncate from the left
	if len(bigEndianBytes) > size {
		bigEndianBytes = bigEndianBytes[len(bigEndianBytes)-size:]
	}

	// Convert to little-endian by reversing the bytes
	// Place the least significant bytes at the beginning
	for i := 0; i < len(bigEndianBytes); i++ {
		result[i] = bigEndianBytes[len(bigEndianBytes)-1-i]
	}

	// The remaining bytes in result are already zero (proper padding)
	return result
}
