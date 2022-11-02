package main

import (
	"bytes"
	"fmt"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/status-im/keycard-go/hexutils"
	"strings"

	"github.com/obscuronet/go-obscuro/go/host/hostrunner"
)

// Runs an Obscuro host as a standalone process.
func main() {
	config, err := hostrunner.ParseConfig()
	if err != nil {
		panic(fmt.Errorf("could not parse config. Cause: %w", err))
	}
	addr := toAddress(config.PrivateKeyString)
	pkadd := config.PKAddress
	if strings.HasPrefix(config.PKAddress, "0x") {
		pkadd = config.PKAddress[2:]
	}

	if config.PKAddress != "" && !bytes.Equal(hexutils.HexToBytes(pkadd), addr.Bytes()) {
		fmt.Printf("WARN: the address: %s does not match the private key %s\n", config.PKAddress, config.PrivateKeyString)
	}
	hostrunner.RunHost(config)
}

func toAddress(privateKey string) gethcommon.Address {
	k := privateKey
	if strings.HasPrefix(privateKey, "0x") {
		k = privateKey[2:]
	}
	privateKeyA, err := crypto.ToECDSA(hexutils.HexToBytes(k))
	if err != nil {
		panic(err)
	}
	pubKeyA := privateKeyA.PublicKey
	return crypto.PubkeyToAddress(pubKeyA)
}
