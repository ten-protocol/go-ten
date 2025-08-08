package main

import (
	"fmt"
	"os"

	"github.com/ten-protocol/go-ten/go/enclave/core/egoutils"
)

func main() {
	keyPath := "/data/encryption-key.json"

	data, err := egoutils.ReadAndUnseal(keyPath)
	if err != nil {
		fmt.Println("Error unsealing key:", err)
		return
	}

	fmt.Println(string(data))
	err = writeDataToFile("/data/encryption-key-unsealed.json", data)
	if err != nil {
		fmt.Println("Error writing data to file:", err)
		return
	}
}

func writeDataToFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0o644)
}
