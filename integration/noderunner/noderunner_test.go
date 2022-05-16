package noderunner

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/enclaverunner"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/hostrunner"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"
	"testing"
	"time"
)

// todo - joel - describe
func TestCanStartStandaloneObscuroHostAndEnclave(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	hostConfig := hostrunner.DefaultHostConfig()
	hostConfig.PrivateKeyString = hex.EncodeToString(crypto.FromECDSA(privateKey))
	enclaveConfig := enclaverunner.DefaultEnclaveConfig()

	gethBinaryPath, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
	if err != nil {
		t.Fatal(err)
	}
	network := gethnetwork.NewGethNetwork(8446, gethBinaryPath, 1, 1, []string{address.String()})
	defer network.StopNodes()
	go enclaverunner.RunEnclave(enclaveConfig)
	go hostrunner.RunHost(hostConfig)

	// todo - joel - avoid sleep
	time.Sleep(3 * time.Second)

	// todo - joel - use proper node ID
	// todo - joel - too brittle that must not add http
	// todo - joel - use 127... in host runner too
	obscuroClient := obscuroclient.NewClient(0, "127.0.0.1:12000")
	var result uint64
	err = obscuroClient.Call(&result, obscuroclient.RPCGetCurrentBlockHeadHeight)
	if err != nil {
		t.Fatal(err)
	}

	if !(result > 0) {
		t.Fatal("Zero blocks have been produced. Something is wrong.")
	}
}
