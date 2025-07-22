package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ten-protocol/go-ten/go/enclave/core/egoutils"
)

func main() {
	keyPath := "/data/encryption-key.json"

	// Log environment info
	fmt.Printf("Starting unsealing process...\n")
	fmt.Printf("Key path: %s\n", keyPath)

	// Check if file exists and get info
	if fileInfo, err := os.Stat(keyPath); err != nil {
		fmt.Printf("Error accessing file %s: %v\n", keyPath, err)
		return
	} else {
		fmt.Printf("File exists - Size: %d bytes, ModTime: %v\n", fileInfo.Size(), fileInfo.ModTime())
	}

	// Read raw file content for debugging
	rawData, err := os.ReadFile(keyPath)
	if err != nil {
		fmt.Printf("Error reading raw file: %v\n", err)
		return
	}
	fmt.Printf("Raw file size: %d bytes\n", len(rawData))
	fmt.Printf("First 32 bytes (hex): %x\n", rawData[:min(32, len(rawData))])

	// Attempt unsealing with detailed error info
	fmt.Printf("Attempting to unseal...\n")
	data, err := egoutils.ReadAndUnseal(keyPath)
	if err != nil {
		fmt.Printf("Error unsealing key: %v\n", err)
		fmt.Printf("Error type: %T\n", err)

		// Additional SGX-specific error info
		if strings.Contains(err.Error(), "OE_INVALID_PARAMETER") {
			fmt.Printf("SGX unsealing failed - possible causes:\n")
			fmt.Printf("  - Data was sealed by different enclave\n")
			fmt.Printf("  - Corrupted sealed data\n")
			fmt.Printf("  - Wrong enclave context\n")
		}
		return
	}

	fmt.Printf("Successfully unsealed data - Size: %d bytes\n", len(data))
	fmt.Println("Unsealed content:")
	fmt.Println(string(data))

	// Write unsealed data
	fmt.Printf("Writing unsealed data to file...\n")
	err = writeDataToFile("/data/encryption-key-unsealed.json", data)
	if err != nil {
		fmt.Printf("Error writing data to file: %v\n", err)
		return
	}

	fmt.Printf("Successfully wrote unsealed data to /data/encryption-key-unsealed.json\n")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func writeDataToFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}
