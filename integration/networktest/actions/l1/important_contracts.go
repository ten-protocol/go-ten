package l1

import (
	"context"
	"fmt"

	"github.com/ten-protocol/go-ten/integration/networktest/actions"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/integration/networktest"
)

func VerifyL2MessageBusAddressAvailable() networktest.Action {
	// VerifyOnly because the L2MessageBus should be deployed automatically by the node as a synthetic tx
	return actions.VerifyOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) error {
		cli, err := obsclient.Dial(network.ValidatorRPCAddress(0))
		if err != nil {
			return errors.Wrap(err, "failed to dial obsClient")
		}
		networkCfg, err := cli.GetConfig()
		if err != nil {
			return errors.Wrap(err, "failed to get network config")
		}

		var _emptyAddress common.Address
		if networkCfg.L2MessageBusAddress == _emptyAddress {
			return errors.New("L2MessageBusAddress not set")
		}
		fmt.Println("L2MessageBusAddress: ", networkCfg.L2MessageBusAddress)
		return nil
	})
}
