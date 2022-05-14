package mgmtcontractlib

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/contracts"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

const methodBytesLen = 4

var (
	// TODO review estimating gas - these should not be static values
	defaultGasPrice = big.NewInt(20000000000)
	defaultGas      = uint64(1024_000_000)
)

type TxHandler interface {
	// PackTx receives an obscurocommon.L1Transaction object and packs it into a types.TxData object
	// Nonce generation, transaction signature and any other operations are responsibility of the caller
	PackTx(tx obscurocommon.L1Transaction, from common.Address, nonce uint64) (types.TxData, error)

	// UnPackTx receives a *types.Transaction and converts it to an obscurocommon.L1Transaction
	// Any transaction that is not calling the management contract is purposefully ignored
	UnPackTx(tx *types.Transaction) obscurocommon.L1Transaction
}

type mgmtContractTxHandler struct {
	contractAddr         *common.Address
	erc20ContractAddress *common.Address
}

func NewEthMgmtContractTxHandler(contractAddress *common.Address) TxHandler {
	return &mgmtContractTxHandler{
		contractAddr: contractAddress,
	}
}

func NewEthMgmtContractTxHandlerWithERC20(contractAddress *common.Address, erc20ContractAddress *common.Address) TxHandler {
	return &mgmtContractTxHandler{
		contractAddr:         contractAddress,
		erc20ContractAddress: erc20ContractAddress,
	}
}

func (h *mgmtContractTxHandler) PackTx(t obscurocommon.L1Transaction, fromAddr common.Address, nonce uint64) (types.TxData, error) {
	ethTx := &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: defaultGasPrice,
		Gas:      defaultGas,
		To:       h.contractAddr,
	}

	// using (obj) type instead of t.Type() to immediately fetch the cast object
	switch tx := t.(type) {
	case *obscurocommon.L1DepositTx:
		ethTx.Value = big.NewInt(int64(tx.Amount))
		data, err := contracts.MgmtContractABIJSON.Pack(contracts.DepositMethod, tx.To)
		if err != nil {
			panic(err)
		}
		ethTx.Data = data
		log.Log(fmt.Sprintf("- Broadcasting - Issuing DepositTx - Addr: %s deposited %d to %s ",
			fromAddr, tx.Amount, tx.To))

	case *obscurocommon.L1RollupTx:
		r, err := nodecommon.DecodeRollup(tx.Rollup)
		if err != nil {
			panic(err)
		}
		zipped, err := Compress(tx.Rollup)
		if err != nil {
			panic(err)
		}
		encRollupData := EncodeToString(zipped)
		data, err := contracts.MgmtContractABIJSON.Pack(contracts.AddRollupMethod, encRollupData)
		if err != nil {
			panic(err)
		}

		ethTx.Data = data
		log.Log(fmt.Sprintf("- Broadcasting - Issuing Rollup: r_%d - %d txs - datasize: %d - gas: %d",
			obscurocommon.ShortHash(r.Hash()), len(r.Transactions), len(data), ethTx.Gas))

	case *obscurocommon.L1StoreSecretTx:
		data, err := contracts.MgmtContractABIJSON.Pack(contracts.StoreSecretMethod, EncodeToString(tx.Secret))
		if err != nil {
			panic(err)
		}
		ethTx.Data = data
		log.Log(fmt.Sprintf("- Broadcasting - Issuing StoreSecretTx: encoded as %s", EncodeToString(tx.Secret)))
	case *obscurocommon.L1RequestSecretTx:
		data, err := contracts.MgmtContractABIJSON.Pack(contracts.RequestSecretMethod)
		if err != nil {
			panic(err)
		}
		ethTx.Data = data
		log.Log("- Broadcasting - Issuing RequestSecret")
	}

	return ethTx, nil
}

func (h *mgmtContractTxHandler) UnPackTx(tx *types.Transaction) obscurocommon.L1Transaction {
	// ignore transactions that are not calling the contract
	if tx.To() == nil || (tx.To().Hex() != h.contractAddr.Hex() && tx.To().Hex() != h.erc20ContractAddress.Hex()) || len(tx.Data()) == 0 {
		log.Log(fmt.Sprintf("UnpackTx: Ignoring transaction %+v", tx))
		return nil
	}

	// TODO review this
	if tx.To().Hex() == h.erc20ContractAddress.Hex() {
		return h.unpackPedroERC20(tx)
	}

	method, err := contracts.MgmtContractABIJSON.MethodById(tx.Data()[:methodBytesLen])
	if err != nil {
		panic(err)
	}

	contractCallData := map[string]interface{}{}
	switch method.Name {
	case contracts.DepositMethod:
		if err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[4:]); err != nil {
			panic(err)
		}
		callData, found := contractCallData["dest"]
		if !found {
			panic("call data not found for dest")
		}

		return &obscurocommon.L1DepositTx{
			Amount:        tx.Value().Uint64(),
			To:            callData.(common.Address),
			TokenContract: nil, // TODO have fixed Token contract for Eth deposits ?
		}

	case contracts.AddRollupMethod:
		if err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[4:]); err != nil {
			panic(err)
		}
		callData, found := contractCallData["rollupData"]
		if !found {
			panic("call data not found for rollupData")
		}
		zipped := DecodeFromString(callData.(string))
		rollup, err := Decompress(zipped)
		if err != nil {
			panic(err)
		}

		return &obscurocommon.L1RollupTx{
			Rollup: rollup,
		}

	case contracts.StoreSecretMethod:
		if err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[4:]); err != nil {
			panic(err)
		}
		callData, found := contractCallData["inputSecret"]
		if !found {
			panic("call data not found for inputSecret")
		}

		return &obscurocommon.L1StoreSecretTx{
			Secret: DecodeFromString(callData.(string)),
		}

	case contracts.RequestSecretMethod:
		return &obscurocommon.L1RequestSecretTx{}
	}

	return nil
}

func (h *mgmtContractTxHandler) unpackPedroERC20(tx *types.Transaction) obscurocommon.L1Transaction {
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
		TokenContract: h.erc20ContractAddress,
	}
}
