package noderunner

import (
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/enclaverunner"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/hostrunner"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"
	"testing"
)

func TestCanStartStandaloneObscuroHostAndEnclave(t *testing.T) {
	gethBinaryPath, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
	if err != nil {
		t.Fatal(err)
	}

	network := gethnetwork.NewGethNetwork(8446, gethBinaryPath, 1, 1, nil)
	defer network.StopNodes()

	go enclaverunner.RunEnclave(enclaverunner.DefaultEnclaveConfig())
	go hostrunner.RunHost(hostrunner.DefaultHostConfig())

	select {}

	// todo - joel - get the above running - currently there's an insufficient gas issue
	// todo - joel - check can make request against running host
}
