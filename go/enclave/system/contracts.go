package system

import (
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/contracts/generated/MessageBus"
	"github.com/ten-protocol/go-ten/contracts/generated/SystemDeployer"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/wallet"
)

// todo (#1549) - implement with cryptography epic
const ( // DO NOT USE OR CHANGE THIS KEY IN THE REST OF THE CODEBASE
	ownerKeyHex = "6e384a07a01263518a18a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682"
)

func GenerateDeploymentTransaction(initCode []byte, wallet wallet.Wallet, logger gethlog.Logger) (*common.L2Tx, error) {
	tx := &types.LegacyTx{
		Nonce:    wallet.GetNonceAndIncrement(), // The first transaction of the owner identity should always be deploying the contract
		Value:    gethcommon.Big0,
		Gas:      500_000_000,     // It's quite the expensive contract.
		GasPrice: gethcommon.Big0, // Synthetic transactions are on the house. Or the house.
		Data:     initCode,        //gethcommon.FromHex(SystemDeployer.SystemDeployerMetaData.Bin),
		To:       nil,             // Geth requires nil instead of gethcommon.Address{} which equates to zero address in order to return receipt.
	}

	types.NewEIP155Signer(big.NewInt(1))
	stx, err := wallet.SignTransaction(tx)
	if err != nil {
		return nil, err
	}

	logger.Info("Generated synthetic deployment transaction for the SystemDeployer contract")

	return stx, nil
}

func SystemDeployerInitTransaction(wallet wallet.Wallet, logger gethlog.Logger, eoaOwner gethcommon.Address) (*common.L2Tx, error) {
	abi, _ := SystemDeployer.SystemDeployerMetaData.GetAbi()
	args, err := abi.Constructor.Inputs.Pack(eoaOwner)
	if err != nil {
		panic(err) // This error is fatal. If the system contracts can't be initialized the network cannot bootstrap.
	}

	bytecode := gethcommon.FromHex(SystemDeployer.SystemDeployerMetaData.Bin)
	initCode := append(bytecode, args...)

	return GenerateDeploymentTransaction(initCode, wallet, logger)
}

func MessageBusInitTransaction(wallet wallet.Wallet, logger gethlog.Logger) (*common.L2Tx, error) {
	return GenerateDeploymentTransaction(gethcommon.FromHex(MessageBus.MessageBusMetaData.Bin), wallet, logger)
}

func GetPlaceholderWallet(chainID *big.Int, logger gethlog.Logger) wallet.Wallet {
	key, _ := crypto.HexToECDSA(ownerKeyHex)
	return wallet.NewInMemoryWalletFromPK(chainID, key, logger)
}

func VerifyLogs(receipt *types.Receipt) error {
	return nil
}

func DeriveAddresses(receipt *types.Receipt) (SystemContractAddresses, error) {
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

type SystemContractAddresses map[string]*gethcommon.Address

func (s *SystemContractAddresses) ToString() string {
	var str string
	for name, addr := range *s {
		str += fmt.Sprintf("%s: %s; ", name, addr.Hex())
	}
	return str
}
