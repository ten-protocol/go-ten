package helpful

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/obscuronet/go-obscuro/integration/networktest"
	"github.com/obscuronet/go-obscuro/integration/networktest/env"
)

func TestRunLocalNetwork(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.EnsureTestLogsSetUp("run-local-network")
	networkConnector, cleanUp, err := env.LocalDevNetwork().Prepare()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanUp()

	fmt.Println("----")
	fmt.Println("Sequencer RPC", networkConnector.SequencerRPCAddress())
	for i := 0; i < networkConnector.NumValidators(); i++ {
		fmt.Println("Validator  ", i, networkConnector.ValidatorRPCAddress(i))
	}
	fmt.Println("----")
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Network running until test is stopped...")
	<-done // Will block here until user hits ctrl+c
}
