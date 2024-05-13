package publicdata

import (
	"context"
	"fmt"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/actions"
)

// VerifyBatchesDataAction tests the batches data RPC endpoint
func VerifyBatchesDataAction() networktest.Action {
	return actions.VerifyOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) error {
		client, err := obsclient.Dial(network.ValidatorRPCAddress(0))
		if err != nil {
			return err
		}

		pagination := common.QueryPagination{
			Offset: 0,
			Size:   20,
		}
		batchListing, err := client.GetBatchesListing(&pagination)
		if err != nil {
			return err
		}
		if len(batchListing.BatchesData) != 20 {
			return fmt.Errorf("expected 20 batches, got %d", len(batchListing.BatchesData))
		}
		if batchListing.Total <= 10 {
			return fmt.Errorf("expected more than 10 batches, got %d", batchListing.Total)
		}
		if batchListing.BatchesData[0].Height.Cmp(batchListing.BatchesData[1].Height) < 0 {
			return fmt.Errorf("expected batches to be sorted by height descending")
		}

		return nil
	})
}
