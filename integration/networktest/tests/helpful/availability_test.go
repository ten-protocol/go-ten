package helpful

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ten-protocol/go-ten/go/obsclient"

	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/actions"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
)

const _testTimeSpan = 30 * time.Second

// basic test that verifies it can connect the L1 client and L2 client and sees block numbers increasing (useful to sanity check testnet issues etc.)
func TestNetworkAvailability(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"network-availability",
		t,
		env.LocalDevNetwork(),
		actions.RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
			client, err := network.GetL1Client()
			if err != nil {
				return nil, err
			}
			blockStart, err := client.BlockNumber()
			if err != nil {
				return nil, err
			}
			fmt.Println("start block", blockStart)

			obsClient, err := obsclient.Dial(network.SequencerRPCAddress())
			if err != nil {
				return nil, err
			}
			batchStart, err := obsClient.BatchNumber()
			if err != nil {
				return nil, err
			}
			fmt.Println("start batch", batchStart)

			time.Sleep(_testTimeSpan)

			blockEnd, err := client.BlockNumber()
			if err != nil {
				return nil, err
			}
			fmt.Println("end block", blockEnd)

			batchEnd, err := obsClient.BatchNumber()
			if err != nil {
				return nil, err
			}
			fmt.Println("end batch", batchEnd)

			if blockEnd <= blockStart {
				return nil, fmt.Errorf("expected blocks to be increasing but went from %d to %d after %s", blockStart, blockEnd, _testTimeSpan)
			}
			if batchEnd <= batchStart {
				return nil, fmt.Errorf("expected batches to be increasing but went from %d to %d after %s", batchStart, batchEnd, _testTimeSpan)
			}
			// log block rate and batch rate
			fmt.Printf("block rate: %f seconds per block\n", _testTimeSpan.Seconds()/float64(blockEnd-blockStart))
			fmt.Printf("batch rate: %f seconds per batch\n", _testTimeSpan.Seconds()/float64(batchEnd-batchStart))

			return ctx, nil
		}),
	)
}
