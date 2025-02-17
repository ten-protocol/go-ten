package main

import (
	"encoding/hex"
	"fmt"

	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/datagenerator"
)

// Noddy tool to generate and test account PKs for manual testing etc.

func main() {
	fmt.Println("Generating account for local dev testing...")

	// generate random PK
	pk := hex.EncodeToString(datagenerator.RandomBytes(32))

	// OR set PK here to check acc number
	// pk := "b56aa8059fe42fd702f32a1055cf06f3f7ca851f9da5cc89541b2fcf9a3c654b"

	fmt.Println("Private key:", pk)

	w1 := wallet.NewInMemoryWalletFromConfig(pk, 0, testlog.Logger())
	fmt.Println("Wallet addr:", w1.Address())

	// print it as a snippet that can be pasted into code to define the PK as a string
	fmt.Printf("\nCode snippet:\naccPK := \"%s\" // account %s\n", pk, w1.Address())
}
