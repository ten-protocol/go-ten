package txhandler

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/buildhelper/buildconstants"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

const methodBytesLen = 4

type TxHandler interface {
	PackTx(tx *obscurocommon.L1TxData, from common.Address, nonce uint64) (types.TxData, error)
	UnPackTx(tx *types.Transaction) *obscurocommon.L1TxData
}

type EthTxHandler struct {
	contractABI abi.ABI
}

func NewEthTxHandler() TxHandler {
	contractABI, err := abi.JSON(strings.NewReader(buildconstants.ContractAbi))
	if err != nil {
		panic(err)
	}

	return &EthTxHandler{
		contractABI: contractABI,
	}
}

func (h *EthTxHandler) PackTx(tx *obscurocommon.L1TxData, fromAddr common.Address, nonce uint64) (types.TxData, error) {
	ethTx := &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: big.NewInt(20000000000),
		Gas:      1024_000_000,
		To:       &buildconstants.ContractAddress,
	}

	// TODO each of these cases should be a function:
	// TODO like: func createRollupTx() or func createDepositTx()
	// TODO And then eventually, these functions would be called directly, when we get rid of our special format. (we'll have to change the mock thing as well for that)
	switch tx.TxType {
	case obscurocommon.DepositTx:
		ethTx.Value = big.NewInt(int64(tx.Amount))
		data, err := h.contractABI.Pack("Deposit", tx.Dest)
		if err != nil {
			panic(err)
		}
		ethTx.Data = data
		log.Log(fmt.Sprintf("BROADCAST TX: Issuing DepositTx - Addr: %s deposited %d to %s ",
			fromAddr, tx.Amount, tx.Dest))

	case obscurocommon.RollupTx:
		r, err := nodecommon.DecodeRollup(tx.Rollup)
		if err != nil {
			panic(err)
		}
		zipped := Compress(tx.Rollup)
		encRollupData := EncodeToString(zipped)
		data, err := h.contractABI.Pack("AddRollup", encRollupData)
		if err != nil {
			panic(err)
		}

		ethTx.Data = data
		derolled, _ := nodecommon.DecodeRollup(tx.Rollup)

		log.Log(fmt.Sprintf("BROADCAST TX - Issuing Rollup: %s - %d txs - datasize: %d - gas: %d \n", r.Hash(), len(derolled.Transactions), len(data), ethTx.Gas))

	case obscurocommon.StoreSecretTx:
		data, err := h.contractABI.Pack("StoreSecret", EncodeToString(tx.Secret))
		if err != nil {
			panic(err)
		}
		ethTx.Data = data
		log.Log(fmt.Sprintf("BROADCAST TX: Issuing StoreSecretTx: encoded as %s", EncodeToString(tx.Secret)))
	case obscurocommon.RequestSecretTx:
		data, err := h.contractABI.Pack("RequestSecret")
		if err != nil {
			panic(err)
		}
		ethTx.Data = data
		log.Log("BROADCAST TX: Issuing RequestSecret")
	}

	return ethTx, nil
}

func (h *EthTxHandler) UnPackTx(tx *types.Transaction) *obscurocommon.L1TxData {
	//if !IsRealEth {
	//	t := obscurocommon.TxData(tx)
	//	return &t
	//}

	// ignore transactions that are not calling the contract
	if tx.To() == nil || tx.To().Hex() != buildconstants.ContractAddress.Hex() || len(tx.Data()) == 0 {
		return nil
	}

	contractABI, err := abi.JSON(strings.NewReader(buildconstants.ContractAbi))
	if err != nil {
		panic(err)
	}

	method, err := contractABI.MethodById(tx.Data()[:methodBytesLen])
	if err != nil {
		panic(err)
	}

	l1txData := obscurocommon.L1TxData{
		TxType:      0,
		Rollup:      nil,
		Secret:      nil,
		Attestation: obscurocommon.AttestationReport{},
		Amount:      0,
		Dest:        common.Address{},
	}
	switch method.Name {
	case "Deposit":
		contractCallData := map[string]interface{}{}
		if err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[4:]); err != nil {
			panic(err)
		}
		callData, found := contractCallData["dest"]
		if !found {
			panic("call data not found for dest")
		}

		l1txData.TxType = obscurocommon.DepositTx
		l1txData.Amount = tx.Value().Uint64()
		l1txData.Dest = callData.(common.Address)

	case "AddRollup":
		contractCallData := map[string]interface{}{}
		if err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[4:]); err != nil {
			panic(err)
		}
		callData, found := contractCallData["rollupData"]
		if !found {
			panic("call data not found for rollupData")
		}
		zipped := DecodeFromString(callData.(string))
		l1txData.Rollup = Decompress(zipped)
		l1txData.TxType = obscurocommon.RollupTx

	case "StoreSecret":
		contractCallData := map[string]interface{}{}
		if err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[4:]); err != nil {
			panic(err)
		}
		callData, found := contractCallData["inputSecret"]
		if !found {
			panic("call data not found for inputSecret")
		}
		l1txData.Secret = DecodeFromString(callData.(string))
		l1txData.TxType = obscurocommon.StoreSecretTx
	}

	return &l1txData
}
