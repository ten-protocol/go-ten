package main

import (
	"fmt"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/go/rpc"
	l1gs "github.com/ten-protocol/go-ten/testnet/launcher/l1grantsequencers"
)

func main() {
	if err := grantSequencerStatus("0x51D43a3Ca257584E770B6188232b199E76B022A2"); err != nil {
		fmt.Printf("Error granting sequencer status: %v\n", err)
		return
	}
}

func grantSequencerStatus(mgmtContractAddr string) error {
	// fetch enclaveIDs
	hostURL := fmt.Sprintf("http://localhost:%d", 80)
	client, err := rpc.NewNetworkClient(hostURL)
	if err != nil {
		return fmt.Errorf("failed to create network client: %w", err)
	}
	defer client.Stop()

	obsClient := obsclient.NewObsClient(client)
	health, err := obsClient.Health()
	if err != nil {
		return fmt.Errorf("failed to get health status: %w", err)
	}

	fmt.Printf("HEALTH: %v\n", health.OverallHealth)
	fmt.Printf("HEALTH ENCLAVES: %d\n", len(health.Enclaves))
	if len(health.Enclaves) > 0 {
		fmt.Printf("ENCLAVE ID: %s\n", health.Enclaves[0].EnclaveID.Hex())
	}

	if len(health.Enclaves) == 0 {
		return fmt.Errorf("could not retrieve enclave IDs from health endpoint")
	}

	var enclaveIDs []string
	for _, status := range health.Enclaves {
		enclaveIDs = append(enclaveIDs, status.EnclaveID.String())
	}
	//enclaveIDsStr := strings.Join(enclaveIDs, ",")

	enclaveIDsStr := "0xBD0D613bCbDbcC93abE025117564cc4435896A5F,0x5555E184dDC7de1A1fD0FF237CcA77338cE7162D,0xa00E66438600c5D104f842cBAf0D7E09fcB76555"

	fmt.Printf("enclaveIDsStr: %s\n", enclaveIDsStr)
	l1grantsequencers, err := l1gs.NewGrantSequencers(
		l1gs.NewGrantSequencerConfig(
			l1gs.WithL1HTTPURL("http://eth2network:8025"),
			l1gs.WithPrivateKey("f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb"),
			l1gs.WithDockerImage("testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest"),
			l1gs.WithMgmtContractAddress(mgmtContractAddr),
			l1gs.WithEnclaveIDs(enclaveIDsStr),
		),
	)
	if err != nil {
		return fmt.Errorf("unable to configure l1 grant sequencers - %w", err)
	}

	err = l1grantsequencers.Start()
	if err != nil {
		return fmt.Errorf("unable to start l1 grant sequencers - %w", err)
	}

	fmt.Println("Enclaves were successfully granted sequencer roles...")
	return nil
}
