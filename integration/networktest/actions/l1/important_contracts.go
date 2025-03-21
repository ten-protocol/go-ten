package l1

import (
	"context"
	"fmt"
	"time"

	"github.com/ten-protocol/go-ten/go/ethadapter/contractlib"

	"github.com/ten-protocol/go-ten/go/ethadapter"

	"github.com/ten-protocol/go-ten/integration/networktest/actions"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/networktest"
)

type setImportantContract struct {
	contractName    string
	contractAddress common.Address
}

func SetImportantContract(contractName string, contractAddress common.Address) networktest.Action {
	return &setImportantContract{
		contractName:    contractName,
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

	networkContract, err := contractlib.NewNetworkConfigLib(networkCfg.NetworkConfig, *l1Client.EthClient())
	if err != nil {
		return ctx, errors.Wrap(err, "failed to get L1 client")
	}

	msg, err := networkContract.AddAdditionalAddress(s.contractName, s.contractAddress)
	if err != nil {
		return ctx, errors.Wrap(err, "failed to create SetImportantContractMsg")
	}

	txData := &types.LegacyTx{
		To:   &networkCfg.NetworkConfig,
		Data: msg.Data,
	}
	contractOwner, err := network.GetContractOwnerWallet()
	if err != nil {
		return ctx, errors.Wrap(err, "failed to get MC owner wallet")
	}
	// !! Important note !!
	// The ownerOnly check in the contract doesn't like the gas estimate in here, to test you may need to hardcode
	// the gas value when the estimate errors
	nonce, err := l1Client.Nonce(networkCfg.NetworkConfig)
	if err != nil {
		return nil, err
	}
	tx, err := ethadapter.SetTxGasPrice(ctx, l1Client, txData, networkCfg.NetworkConfig, nonce, 0, nil, testlog.Logger())
	if err != nil {
		return ctx, errors.Wrap(err, "failed to prepare tx")
	}
	signedTx, err := contractOwner.SignTransaction(tx)
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

	if len(networkCfg.AdditionalContracts) == 0 {
		return errors.New("no important contracts set")
	}
	found := false
	for _, contract := range networkCfg.AdditionalContracts {
		if contract.Name == s.contractName && contract.Addr == s.contractAddress {
			found = true
			break
		}
	}
	if !found {
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
		if networkCfg.L2MessageBus == _emptyAddress {
			return errors.New("L2MessageBus not set")
		}
		fmt.Println("L2MessageBusAddress: ", networkCfg.L2MessageBus)
		return nil
	})
}
