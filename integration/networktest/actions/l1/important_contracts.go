package l1

import (
	"context"
	"fmt"
	"time"

	"github.com/ten-protocol/go-ten/go/ethadapter"

	"github.com/ten-protocol/go-ten/integration/networktest/actions"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/networktest"
)

type setImportantContract struct {
	contractKey     string
	contractAddress common.Address
}

func SetImportantContract(contractKey string, contractAddress common.Address) networktest.Action {
	return &setImportantContract{
		contractKey:     contractKey,
		contractAddress: contractAddress,
	}
}

func (s *setImportantContract) Run(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
	obsClient, err := obsclient.Dial(network.ValidatorRPCAddress(0))
	if err != nil {
		return ctx, errors.Wrap(err, "failed to dial obsClient")
	}

	networkCfg, err := obsClient.GetConfig()
	if err != nil {
		return ctx, errors.Wrap(err, "failed to get network config")
	}

	l1Client, err := network.GetL1Client()
	if err != nil {
		return ctx, errors.Wrap(err, "failed to get L1 client")
	}

	mgmtContract := mgmtcontractlib.NewMgmtContractLib(&networkCfg.ManagementContractAddress, testlog.Logger())

	msg, err := mgmtContract.SetImportantContractMsg(s.contractKey, s.contractAddress)
	if err != nil {
		return ctx, errors.Wrap(err, "failed to create SetImportantContractMsg")
	}

	txData := &types.LegacyTx{
		To:   &networkCfg.ManagementContractAddress,
		Data: msg.Data,
	}
	mcOwner, err := network.GetMCOwnerWallet()
	if err != nil {
		return ctx, errors.Wrap(err, "failed to get MC owner wallet")
	}
	// !! Important note !!
	// The ownerOnly check in the contract doesn't like the gas estimate in here, to test you may need to hardcode a
	// the gas value when the estimate errors
	nonce, err := l1Client.Nonce(networkCfg.ManagementContractAddress)
	if err != nil {
		return nil, err
	}
	tx, err := ethadapter.SetTxGasPrice(ctx, l1Client, txData, networkCfg.ManagementContractAddress, nonce, 0, testlog.Logger())
	if err != nil {
		return ctx, errors.Wrap(err, "failed to prepare tx")
	}
	signedTx, err := mcOwner.SignTransaction(tx)
	if err != nil {
		return ctx, errors.Wrap(err, "failed to sign tx")
	}
	err = l1Client.SendTransaction(signedTx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send tx")
	}

	// wait for tx to be mined
	return ctx, retry.Do(func() error {
		receipt, err := l1Client.TransactionReceipt(signedTx.Hash())
		if err != nil {
			return err
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			return retry.FailFast(errors.New("tx failed"))
		}
		return nil
	}, retry.NewTimeoutStrategy(15*time.Second, 1*time.Second))
}

func (s *setImportantContract) Verify(_ context.Context, network networktest.NetworkConnector) error {
	cli, err := obsclient.Dial(network.ValidatorRPCAddress(0))
	if err != nil {
		return errors.Wrap(err, "failed to dial obsClient")
	}
	networkCfg, err := cli.GetConfig()
	if err != nil {
		return errors.Wrap(err, "failed to get network config")
	}

	if networkCfg.ImportantContracts == nil || len(networkCfg.ImportantContracts) == 0 {
		return errors.New("no important contracts set")
	}
	if addr, ok := networkCfg.ImportantContracts[s.contractKey]; !ok || addr != s.contractAddress {
		return errors.New("important contract not set")
	}
	return nil
}

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
