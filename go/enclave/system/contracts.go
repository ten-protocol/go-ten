package system

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/log"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/contracts/generated/SystemDeployer"
	"github.com/ten-protocol/go-ten/go/common"
)

func GenerateDeploymentTransaction(initCode []byte, logger gethlog.Logger) (*common.L2Tx, error) {
	tx := &types.LegacyTx{
		Nonce:    0, // The first transaction of the owner identity should always be deploying the contract
		Value:    gethcommon.Big0,
		Gas:      20_000_000,      // It's quite the expensive contract.
		GasPrice: gethcommon.Big0, // Synthetic transactions are on the house. Or the house.
		Data:     initCode,        // gethcommon.FromHex(SystemDeployer.SystemDeployerMetaData.Bin),
		To:       nil,             // Geth requires nil instead of gethcommon.Address{} which equates to zero address in order to return receipt.
	}

	stx := types.NewTx(tx)

	logger.Info("Generated synthetic deployment transaction for the SystemDeployer contract")

	return stx, nil
}

func SystemDeployerInitTransaction(logger gethlog.Logger, eoaOwner gethcommon.Address, l1BridgeAddress gethcommon.Address) (*common.L2Tx, error) {
	abi, _ := SystemDeployer.SystemDeployerMetaData.GetAbi()
	args, err := abi.Constructor.Inputs.Pack(eoaOwner, l1BridgeAddress)
	if err != nil {
		logger.Crit("This error is fatal. If the system contracts can't be initialized the network cannot bootstrap.", log.ErrKey, err)
	}

	bytecode := gethcommon.FromHex(SystemDeployer.SystemDeployerMetaData.Bin)
	initCode := append(bytecode, args...)

	return GenerateDeploymentTransaction(initCode, logger)
}

func VerifyLogs(receipt *types.Receipt) error {
	return nil
}

func DeriveAddresses(receipt *types.Receipt) (common.SystemContractAddresses, error) {
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("cannot derive system contract addresses from failed receipt")
	}

	addresses := make(map[string]*gethcommon.Address)

	eventName := "SystemContractDeployed"
	abi, err := SystemDeployer.SystemDeployerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	eventID := abi.Events[eventName].ID

	for _, log := range receipt.Logs {
		if log.Topics[0] != eventID {
			continue
		}

		var event SystemDeployer.SystemDeployerSystemContractDeployed
		err := abi.UnpackIntoInterface(&event, eventName, log.Data)
		if err != nil {
			return nil, err
		}

		addresses[event.Name] = &event.ContractAddress
	}

	return addresses, nil
}
