package system

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/log"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/contracts/generated/SystemDeployerPhase1"
	"github.com/ten-protocol/go-ten/contracts/generated/SystemDeployerPhase2"
	"github.com/ten-protocol/go-ten/go/common"
)

func GenerateDeploymentTransaction(initCode []byte, nonce uint64, logger gethlog.Logger) (*common.L2Tx, error) {
	tx := &types.LegacyTx{
		Nonce:    nonce,
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

func SystemDeployerPhase1InitTransaction(logger gethlog.Logger, eoaOwner gethcommon.Address) (*common.L2Tx, error) {
	// Phase 1 constructor signature: constructor(eoaAdmin)
	abi, _ := SystemDeployerPhase1.SystemDeployerPhase1MetaData.GetAbi()
	args, err := abi.Constructor.Inputs.Pack(eoaOwner)
	if err != nil {
		logger.Crit("This error is fatal. If the system contracts can't be initialized the network cannot bootstrap.", log.ErrKey, err)
	}
	bytecode := gethcommon.FromHex(SystemDeployerPhase1.SystemDeployerPhase1MetaData.Bin)
	initCode := append(bytecode, args...)

	logger.Info("Generated synthetic deployment transaction for SystemDeployerPhase1")
	return GenerateDeploymentTransaction(initCode, 0, logger)
}

func SystemDeployerPhase2InitTransaction(logger gethlog.Logger, eoaOwner gethcommon.Address, feesAddress gethcommon.Address, l1BridgeAddress gethcommon.Address) (*common.L2Tx, error) {
	// Phase 1 constructor signature: constructor(eoaAdmin)
	abi, _ := SystemDeployerPhase2.SystemDeployerPhase2MetaData.GetAbi()
	args, err := abi.Constructor.Inputs.Pack(eoaOwner, feesAddress, l1BridgeAddress)
	if err != nil {
		logger.Crit("This error is fatal. If the system contracts can't be initialized the network cannot bootstrap.", log.ErrKey, err)
	}
	bytecode := gethcommon.FromHex(SystemDeployerPhase2.SystemDeployerPhase2MetaData.Bin)
	initCode := append(bytecode, args...)

	logger.Info("Generated synthetic deployment transaction for SystemDeployerPhase1")
	return GenerateDeploymentTransaction(initCode, 1, logger)
}

func VerifyLogs(receipt *types.Receipt) error {
	return nil
}

func DeriveAddresses(receipt *types.Receipt) (common.SystemContractAddresses, error) {
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("cannot derive system contract addresses from failed receipt")
	}

	addresses := make(map[string]*gethcommon.Address)

	// Try Phase 1 events first
	phase1Addresses, err := derivePhase1Addresses(receipt)
	if err == nil {
		for name, addr := range phase1Addresses {
			addresses[name] = addr
		}
	}

	// Try Phase 2 events
	phase2Addresses, err := derivePhase2Addresses(receipt)
	if err == nil {
		for name, addr := range phase2Addresses {
			addresses[name] = addr
		}
	}

	return addresses, nil
}

func derivePhase1Addresses(receipt *types.Receipt) (common.SystemContractAddresses, error) {
	addresses := make(map[string]*gethcommon.Address)
	eventName := "SystemContractDeployed"

	abi, err := SystemDeployerPhase1.SystemDeployerPhase1MetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	eventID := abi.Events[eventName].ID

	for _, log := range receipt.Logs {
		if log.Topics[0] != eventID {
			continue
		}

		var event SystemDeployerPhase1.SystemDeployerPhase1SystemContractDeployed
		err := abi.UnpackIntoInterface(&event, eventName, log.Data)
		if err != nil {
			continue // Skip if parsing fails
		}

		addresses[event.Name] = &event.ContractAddress
	}

	return addresses, nil
}

func derivePhase2Addresses(receipt *types.Receipt) (common.SystemContractAddresses, error) {
	addresses := make(map[string]*gethcommon.Address)
	eventName := "SystemContractDeployed"

	abi, err := SystemDeployerPhase2.SystemDeployerPhase2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	eventID := abi.Events[eventName].ID

	for _, log := range receipt.Logs {
		if log.Topics[0] != eventID {
			continue
		}

		var event SystemDeployerPhase2.SystemDeployerPhase2SystemContractDeployed
		err := abi.UnpackIntoInterface(&event, eventName, log.Data)
		if err != nil {
			continue // Skip if parsing fails
		}

		addresses[event.Name] = &event.ContractAddress
	}

	return addresses, nil
}
