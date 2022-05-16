package noderunner

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/enclaverunner"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/hostrunner"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"
)

// A smoke test to check that we can stand up a standalone Obscuro host and enclave.
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

	// We sleep to give the network time to produce some blocks.
	time.Sleep(3 * time.Second)

	obscuroClient := obscuroclient.NewClient(common.BytesToAddress([]byte(hostConfig.NodeID)), hostConfig.ClientServerAddr)
	var result uint64
	err = obscuroClient.Call(&result, obscuroclient.RPCGetCurrentBlockHeadHeight)
	if err != nil {
		t.Fatal(err)
	}

	if !(result > 0) {
		t.Fatal("Zero blocks have been produced. Something is wrong.")
	}
}
