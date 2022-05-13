package noderunner

import (
	"flag"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/enclaverunner"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/hostrunner"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"
	"os"
	"testing"
)

func TestCanStartStandaloneObscuroHostAndEnclave(t *testing.T) {
	gethBinaryPath, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
	if err != nil {
		t.Fatal(err)
	}

	network := gethnetwork.NewGethNetwork(8446, gethBinaryPath, 1, 1, nil)
	defer network.StopNodes()

	defaultEnclaveConfig := enclaverunner.ParseCLIArgs()
	go enclaverunner.RunEnclave(defaultEnclaveConfig)

	// todo - joel - need to reset flags. something like the below
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	defaultHostConfig := hostrunner.ParseCLIArgs()
	go hostrunner.RunHost(defaultHostConfig)

	select {}

	// todo - joel - check can make request against running host
}
