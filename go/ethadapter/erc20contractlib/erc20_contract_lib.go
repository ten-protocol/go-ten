package erc20contractlib

import (
	"fmt"
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
	DecodeTx(tx *types.Transaction) (common.L1TenTransaction, error)

	// CreateDepositTx receives an common.L1Transaction and converts it to an eth transaction
	CreateDepositTx(tx *common.L1DepositTx) (types.TxData, error)
}

// erc20ContractLibImpl takes a mgmtContractAddr and processes multiple erc20ContractAddrs
// Watches for contract executions that might be deposits towards the Management Contract
type erc20ContractLibImpl struct {
	crossChainContractAddr *gethcommon.Address
	erc20ContractAddrs     []*gethcommon.Address
	contractABI            abi.ABI
}

func NewERC20ContractLib(crossChainContractAddr *gethcommon.Address, contractAddrs ...*gethcommon.Address) ERC20ContractLib {
	contractABI, err := abi.JSON(strings.NewReader(ERC20ContractABI))
	if err != nil {
		panic(err)
	}

	return &erc20ContractLibImpl{
		crossChainContractAddr: crossChainContractAddr,
		erc20ContractAddrs:     contractAddrs,
		contractABI:            contractABI,
	}
}

func (c *erc20ContractLibImpl) CreateDepositTx(tx *common.L1DepositTx) (types.TxData, error) {
	data, err := c.contractABI.Pack("transfer", &tx.To, tx.Amount)
	if err != nil {
		return nil, fmt.Errorf("failed to pack transfer data: %w", err)
	}

	return &types.LegacyTx{
		To:   tx.TokenContract,
		Data: data,
	}, nil
}

func (c *erc20ContractLibImpl) DecodeTx(tx *types.Transaction) (common.L1TenTransaction, error) {
	if !c.isRelevant(tx) {
		return nil, nil
	}
	method, err := c.contractABI.MethodById(tx.Data()[:methodBytesLen])
	if err != nil {
		return nil, fmt.Errorf("failed to extract method from tx data: %w", err)
	}

	contractCallData := map[string]interface{}{}
	if err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[methodBytesLen:]); err != nil {
		return nil, fmt.Errorf("failed to unpack contract call data: %w", err)
	}

	to, found := contractCallData[ToCallData]
	if !found {
		return nil, fmt.Errorf("to not found for transfer")
	}

	// only process transfers made to the management contract
	toAddr, ok := to.(gethcommon.Address)
	if !ok {
		return nil, nil
	}
	//FIXME not sure if this is correct
	//println("ERC20 TO ADDR: ", toAddr.Hex())

	amount, found := contractCallData[AmountCallData]
	if !found {
		return nil, fmt.Errorf("amount not found for transfer")
	}

	signer := types.NewLondonSigner(tx.ChainId())
	sender, err := signer.Sender(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to extract sender from tx: %w", err)
	}

	return &common.L1DepositTx{
		Amount:        amount.(*big.Int),
		To:            &toAddr,
		TokenContract: tx.To(),
		Sender:        &sender,
	}, nil
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
