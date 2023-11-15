package l1

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common/retry"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/obsclient"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"
	"github.com/obscuronet/go-obscuro/integration/networktest"
	"github.com/pkg/errors"
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
	tx, err := l1Client.PrepareTransactionToSend(txData, networkCfg.ManagementContractAddress, mcOwner.GetNonceAndIncrement())
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
