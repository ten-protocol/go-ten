package stabletokencontractlib

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/contracts"
	"github.com/obscuronet/obscuro-playground/go/ethclient/txhandler"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

const methodBytesLen = 4

type stableTokenHandler struct {
	addr *common.Address
}

func NewHandler(addr *common.Address) txhandler.ContractHandler {
	return &stableTokenHandler{
		addr: addr,
	}
}

func (h *stableTokenHandler) Address() *common.Address {
	return h.addr
}

func (h *stableTokenHandler) Pack(t obscurocommon.L1Transaction, fromAddr common.Address, nonce uint64) (types.TxData, error) {
	panic("Not implemented")
}

func (h *stableTokenHandler) UnPack(tx *types.Transaction) obscurocommon.L1Transaction {
	method, err := contracts.StableTokenERC20ContractABIJSON.MethodById(tx.Data()[:methodBytesLen])
	if err != nil {
		panic(err)
	}

	contractCallData := map[string]interface{}{}
	if err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[4:]); err != nil {
		panic(err)
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
		TokenContract: h.addr,
	}
}
