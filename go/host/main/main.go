package main

import (
	"bytes"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/status-im/keycard-go/hexutils"

	"github.com/obscuronet/go-obscuro/go/host/hostrunner"
)

// Runs an Obscuro host as a standalone process.
func main() {
	config, err := hostrunner.ParseConfig()
	if err != nil {
		panic(fmt.Errorf("could not parse config. Cause: %w", err))
	}
	addr := toAddress(config.PrivateKeyString)
	if !bytes.Equal(hexutils.HexToBytes(config.PKAddress), addr.Bytes()) {
		panic(fmt.Errorf("the address %s does not match the private key %s", config.PKAddress, config.PrivateKeyString))
	}
	hostrunner.RunHost(config)
}

func toAddress(privateKey string) gethcommon.Address {
	privateKeyA, err := crypto.ToECDSA(hexutils.HexToBytes(privateKey))
	if err != nil {
		panic(err)
	}
	pubKeyA := privateKeyA.PublicKey
	return crypto.PubkeyToAddress(pubKeyA)
}
