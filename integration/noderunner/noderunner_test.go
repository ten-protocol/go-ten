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

	time.Sleep(3 * time.Second)

	obscuroClient := obscuroclient.NewClient(0, hostConfig.ClientServerAddr)
	var result uint64
	err = obscuroClient.Call(&result, obscuroclient.RPCBalance, address)
	if err != nil {
		t.Fatal(err)
	}

	println(result)

	select {}

	// todo - joel - use node client to check can make request against running host
}
