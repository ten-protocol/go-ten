package erc20contractlib

import (
	"math/big"
	"strings"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const methodBytesLen = 4

// ERC20ContractLib provides methods for handling erc20 contracts
type ERC20ContractLib interface {
	// DecodeTx receives a *types.Transaction and converts it to an common.L1Transaction
	// returns nil if the transaction is not convertible
	DecodeTx(tx *types.Transaction) common.TenTransaction

	// CreateDepositTx receives an common.L1Transaction and converts it to an eth transaction
	CreateDepositTx(tx *common.L1DepositTx) types.TxData
}

// erc20ContractLibImpl takes a mgmtContractAddr and processes multiple erc20ContractAddrs
// Watches for contract executions that might be deposits towards the Management Contract
type erc20ContractLibImpl struct {
	mgmtContractAddr   *gethcommon.Address
	erc20ContractAddrs []*gethcommon.Address
	contractABI        abi.ABI
}

func NewERC20ContractLib(mgmtContractAddr *gethcommon.Address, contractAddrs ...*gethcommon.Address) ERC20ContractLib {
	contractABI, err := abi.JSON(strings.NewReader(ERC20ContractABI))
	if err != nil {
		panic(err)
	}

	return &erc20ContractLibImpl{
		mgmtContractAddr:   mgmtContractAddr,
		erc20ContractAddrs: contractAddrs,
		contractABI:        contractABI,
	}
}

func (c *erc20ContractLibImpl) CreateDepositTx(tx *common.L1DepositTx) types.TxData {
	data, err := c.contractABI.Pack("transfer", &tx.To, tx.Amount)
	if err != nil {
		panic(err)
	}

	return &types.LegacyTx{
		To:   tx.TokenContract,
		Data: data,
	}
}

func (c *erc20ContractLibImpl) DecodeTx(tx *types.Transaction) common.TenTransaction {
	if !c.isRelevant(tx) {
		return nil
	}
	method, err := c.contractABI.MethodById(tx.Data()[:methodBytesLen])
	if err != nil {
		panic(err)
	}

	contractCallData := map[string]interface{}{}
	if err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[methodBytesLen:]); err != nil {
		panic(err)
	}

	to, found := contractCallData[ToCallData]
	if !found {
		panic("to address not found for transfer")
	}

	// only process transfers made to the management contract
	toAddr, ok := to.(gethcommon.Address)
	if !ok || toAddr.Hex() != c.mgmtContractAddr.Hex() {
		return nil
	}

	amount, found := contractCallData[AmountCallData]
	if !found {
		panic("amount not found for transfer")
	}

	signer := types.NewLondonSigner(tx.ChainId())
	sender, err := signer.Sender(tx)
	if err != nil {
		panic(err)
	}

	return &common.L1DepositTx{
		Amount:        amount.(*big.Int),
		To:            &toAddr,
		TokenContract: tx.To(),
		Sender:        &sender,
	}
}

func (c *erc20ContractLibImpl) isRelevant(tx *types.Transaction) bool {
	if tx.To() == nil || len(tx.Data()) == 0 {
		return false
	}
	for _, addr := range c.erc20ContractAddrs {
		if tx.To().Hex() == addr.Hex() {
			return true
		}
	}
	return false
}
