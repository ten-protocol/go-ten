package helpful

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/integration/networktest"
	"github.com/obscuronet/go-obscuro/integration/networktest/actions"
	"github.com/obscuronet/go-obscuro/integration/networktest/env"
)

// super basic test that verifies it can connect to the L1 client and sees block numbers increasing (useful to sanity check testnet issues etc.)
func TestL1IsAvailableAndProducingBlocks(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"l1-availability",
		t,
		env.DevTestnet(),
		actions.RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
			client, err := network.GetL1Client()
			if err != nil {
				return nil, err
			}
			blockStart, err := client.BlockNumber()
			if err != nil {
				return nil, err
			}
			fmt.Println("current block", blockStart)
			time.Sleep(20 * time.Second)

			blockEnd, err := client.BlockNumber()
			if err != nil {
				return nil, err
			}
			fmt.Println("current block", blockEnd)
			if blockEnd <= blockStart {
				return nil, fmt.Errorf("expected blocks to be increasing but went from %d to %d after 20secs", blockStart, blockEnd)
			}
			return ctx, nil
		}),
	)
}
