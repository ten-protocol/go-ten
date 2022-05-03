package mgmtcontractlib

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/obscuronet/obscuro-playground/contracts"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
	// PackTx receives an obscurocommon.L1TxData object and packs it into a types.TxData object
	// Nonce generation, transaction signature and any other operations are responsibility of the caller
	PackTx(tx *obscurocommon.L1TxData, from common.Address, nonce uint64) (types.TxData, error)

	// UnPackTx receives a *types.Transaction and converts it to an obscurocommon.L1TxData pointer
	// Any transaction that is not calling the management contract is purposefully ignored
	UnPackTx(tx *types.Transaction) *obscurocommon.L1TxData
}

type mgmtContractTxHandler struct {
	contractAddr common.Address
}

func NewEthMgmtContractTxHandler(contractAddress common.Address) TxHandler {
	return &mgmtContractTxHandler{
		contractAddr: contractAddress,
	}
}

func (h *mgmtContractTxHandler) PackTx(tx *obscurocommon.L1TxData, fromAddr common.Address, nonce uint64) (types.TxData, error) {
	ethTx := &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: defaultGasPrice,
		Gas:      defaultGas,
		To:       &h.contractAddr,
	}

	switch tx.TxType {
	case obscurocommon.DepositTx:
		ethTx.Value = big.NewInt(int64(tx.Amount))
		data, err := contracts.MgmtContractABIJSON.Pack(contracts.DepositMethod, tx.Dest)
		if err != nil {
			log.Panic("could not pack transaction. Cause: %s", err)
		}
		ethTx.Data = data
		log.Info(fmt.Sprintf("- Broadcasting - Issuing DepositTx - Addr: %s deposited %d to %s ",
			fromAddr, tx.Amount, tx.Dest))

	case obscurocommon.RollupTx:
		r, err := nodecommon.DecodeRollup(tx.Rollup)
		if err != nil {
			log.Panic("could not decode rollup. Cause: %s", err)
		}
		zipped, err := Compress(tx.Rollup)
		if err != nil {
			log.Panic("could not compress rollup. Cause: %s", err)
		}
		encRollupData := EncodeToString(zipped)
		data, err := contracts.MgmtContractABIJSON.Pack(contracts.AddRollupMethod, encRollupData)
		if err != nil {
			log.Panic("could not pack transaction. Cause: %s", err)
		}

		ethTx.Data = data
		log.Info(fmt.Sprintf("- Broadcasting - Issuing Rollup: r_%d - %d txs - datasize: %d - gas: %d",
			obscurocommon.ShortHash(r.Hash()), len(r.Transactions), len(data), ethTx.Gas))

	case obscurocommon.StoreSecretTx:
		data, err := contracts.MgmtContractABIJSON.Pack(contracts.StoreSecretMethod, EncodeToString(tx.Secret), EncodeToString(tx.Attestation))
		if err != nil {
			log.Panic("could not pack transaction. Cause: %s", err)
		}
		ethTx.Data = data
		log.Info(fmt.Sprintf("- Broadcasting - Issuing StoreSecretTx: encoded as %s", EncodeToString(tx.Secret)))
	case obscurocommon.RequestSecretTx:
		data, err := contracts.MgmtContractABIJSON.Pack(contracts.RequestSecretMethod, EncodeToString(tx.Attestation))
		if err != nil {
			log.Panic("could not pack transaction. Cause: %s", err)
		}
		ethTx.Data = data
		log.Info("- Broadcasting - Issuing RequestSecret")
	}

	return ethTx, nil
}

func (h *mgmtContractTxHandler) UnPackTx(tx *types.Transaction) *obscurocommon.L1TxData {
	// ignore transactions that are not calling the contract
	if tx.To() == nil || tx.To().Hex() != h.contractAddr.Hex() || len(tx.Data()) == 0 {
		log.Info(fmt.Sprintf("UnpackTx: Ignoring transaction %+v", tx))
		return nil
	}

	method, err := contracts.MgmtContractABIJSON.MethodById(tx.Data()[:methodBytesLen])
	if err != nil {
		log.Panic("could not get contract method. Cause: %s", err)
	}

	l1txData := obscurocommon.L1TxData{
		TxType:      0,
		Attestation: obscurocommon.EncodedAttestationReport{},
		Amount:      0,
		Dest:        common.Address{},
	}
	contractCallData := map[string]interface{}{}
	switch method.Name {
	case contracts.DepositMethod:
		unpackDepositTx(tx, method, contractCallData, &l1txData)

	case contracts.AddRollupMethod:
		unpackRollupTx(tx, method, contractCallData, &l1txData)

	case contracts.StoreSecretMethod:
		unpackStoreSecreTx(tx, method, contractCallData, &l1txData)

	case contracts.RequestSecretMethod:
		unpackRequestSecretTx(tx, method, contractCallData, &l1txData)
	}

	return &l1txData
}

func unpackRequestSecretTx(tx *types.Transaction, method *abi.Method, contractCallData map[string]interface{}, l1txData *obscurocommon.L1TxData) {
	err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[methodBytesLen:])
	if err != nil {
		panic(err)
	}
	callData, found := contractCallData["requestReport"]
	if !found {
		panic("call data not found for requestReport")
	}

	l1txData.Attestation, err = DecodeFromString(callData.(string))
	if err != nil {
		log.Panic("could not decode attestation request. Cause: %s", err)
	}
	l1txData.TxType = obscurocommon.RequestSecretTx
}

func unpackStoreSecreTx(tx *types.Transaction, method *abi.Method, contractCallData map[string]interface{}, l1txData *obscurocommon.L1TxData) {
	err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[methodBytesLen:])
	if err != nil {
		log.Panic("could not unpack transaction. Cause: %s", err)
	}
	secretData, found := contractCallData["inputSecret"]
	if !found {
		log.Panic("call data not found for inputSecret")
	}
	l1txData.Secret, err = DecodeFromString(secretData.(string))
	if err != nil {
		log.Panic("could not decode secret data. Cause: %s", err)
	}

	reportData, found := contractCallData["requestReport"]
	if !found {
		panic("call data not found for requestReport")
	}
	l1txData.Attestation, err = DecodeFromString(reportData.(string))
	if err != nil {
		log.Panic("could not decode report data. Cause: %s", err)
	}

	l1txData.TxType = obscurocommon.StoreSecretTx
}

func unpackRollupTx(tx *types.Transaction, method *abi.Method, contractCallData map[string]interface{}, l1txData *obscurocommon.L1TxData) {
	if err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[methodBytesLen:]); err != nil {
		log.Panic("could not unpack transaction. Cause: %s", err)
	}
	callData, found := contractCallData["rollupData"]
	if !found {
		log.Panic("call data not found for rollupData")
	}
	zipped, err := DecodeFromString(callData.(string))
	if err != nil {
		log.Panic("could not decode call data. Cause: %s", err)
	}

	l1txData.Rollup, err = Decompress(zipped)
	if err != nil {
		log.Panic("could not decompress rollup. Cause: %s", err)
	}
	l1txData.TxType = obscurocommon.RollupTx
}

func unpackDepositTx(tx *types.Transaction, method *abi.Method, contractCallData map[string]interface{}, l1txData *obscurocommon.L1TxData) {
	if err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[methodBytesLen:]); err != nil {
		log.Panic("could not unpack transaction. Cause: %s", err)
	}
	callData, found := contractCallData["dest"]
	if !found {
		log.Panic("call data not found for dest")
	}

	l1txData.TxType = obscurocommon.DepositTx
	l1txData.Amount = tx.Value().Uint64()
	l1txData.Dest = callData.(common.Address)
}
