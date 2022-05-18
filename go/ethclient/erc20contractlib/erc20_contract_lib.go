package erc20contractlib

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/contracts"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

const methodBytesLen = 4

// ERC20ContractLib provides methods for handling erc20 contracts
type ERC20ContractLib interface {
	// DecodeTx receives a *types.Transaction and converts it to an obscurocommon.L1Transaction
	// returns nil if the transaction is not convertible
	DecodeTx(tx *types.Transaction) obscurocommon.L1Transaction
}

// tokenContractLibImpl takes a mgmtContractAddr and processes multiple erc20ContractAddrs
// Watches for contract executions that might be deposits towards the Management Contract
type tokenContractLibImpl struct {
	mgmtContractAddr   *common.Address
	erc20ContractAddrs []*common.Address
}

func (t *tokenContractLibImpl) DecodeTx(tx *types.Transaction) obscurocommon.L1Transaction {
	if !t.isRelevant(tx) {
		return nil
	}
	method, err := contracts.StableTokenERC20ContractABIJSON.MethodById(tx.Data()[:methodBytesLen])
	if err != nil {
		panic(err)
	}

	contractCallData := map[string]interface{}{}
	if err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[methodBytesLen:]); err != nil {
		panic(err)
	}

	to, found := contractCallData[contracts.ToCallData]
	if !found {
		panic("to address not found for transfer")
	}

	// only process transfers made to the management contract
	if toAddr, ok := to.(common.Address); !ok || toAddr.Hex() != t.mgmtContractAddr.Hex() {
		return nil
	}

	amount, found := contractCallData[contracts.AmountCallData]
	if !found {
		panic("amount not found for transfer")
	}

	signer := types.NewEIP155Signer(tx.ChainId())
	sender, err := signer.Sender(tx)
	if err != nil {
		panic(err)
	}

	return &obscurocommon.L1DepositTx{
		Amount:        amount.(*big.Int).Uint64(),
		To:            sender,
		TokenContract: tx.To(),
	}
}

func (t *tokenContractLibImpl) isRelevant(tx *types.Transaction) bool {
	if tx.To() == nil || len(tx.Data()) == 0 {
		return false
	}
	for _, addr := range t.erc20ContractAddrs {
		if tx.To().Hex() == addr.Hex() {
			return true
		}
	}
	return false
}

func NewERC20ContractLib(mgmtContractAddr *common.Address, contractAddrs ...*common.Address) ERC20ContractLib {
	return &tokenContractLibImpl{
		mgmtContractAddr:   mgmtContractAddr,
		erc20ContractAddrs: contractAddrs,
	}
}
