package erc20contractlib

import (
	"math/big"

	"github.com/obscuronet/obscuro-playground/go/ethclient/txhandler"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/contracts"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

const methodBytesLen = 4

type erc20Handler struct {
	addr *common.Address
}

func NewHandler(addr *common.Address) txhandler.ContractHandler {
	return &erc20Handler{
		addr: addr,
	}
}

func (h *erc20Handler) Address() *common.Address {
	return h.addr
}

func (h *erc20Handler) PackTx(t obscurocommon.L1Transaction, fromAddr common.Address, nonce uint64) (types.TxData, error) {
	panic("Not implemented")
}

func (h *erc20Handler) UnPackTx(tx *types.Transaction) obscurocommon.L1Transaction {
	method, err := contracts.PedroERC20ContractABIJSON.MethodById(tx.Data()[:methodBytesLen])
	if err != nil {
		panic(err)
	}

	contractCallData := map[string]interface{}{}
	if err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[4:]); err != nil {
		panic(err)
	}
	amount, found := contractCallData["amount"]
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
		TokenContract: h.addr,
	}
}
