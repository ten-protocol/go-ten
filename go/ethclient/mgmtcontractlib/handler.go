package mgmtcontractlib

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/contracts"
	"github.com/obscuronet/obscuro-playground/go/ethclient/txhandler"
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

type mgmtContractTxHandler struct {
	addr *common.Address
}

func NewHandler(addr *common.Address) txhandler.ContractHandler {
	return &mgmtContractTxHandler{
		addr: addr,
	}
}

func (h *mgmtContractTxHandler) Address() *common.Address {
	return h.addr
}

func (h *mgmtContractTxHandler) Pack(t obscurocommon.L1Transaction, fromAddr common.Address, nonce uint64) (types.TxData, error) {
	ethTx := &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: defaultGasPrice,
		Gas:      defaultGas,
		To:       h.addr,
	}

	switch tx := t.(type) {
	case *obscurocommon.L1DepositTx:
		ethTx.Value = big.NewInt(int64(tx.Amount))
		data, err := contracts.MgmtContractABIJSON.Pack(contracts.DepositMethod, tx.To)
		if err != nil {
			panic(err)
		}
		ethTx.Data = data
		log.Info("- Broadcasting - Issuing DepositTx - Addr: %s deposited %d to %s ", fromAddr, tx.Amount, tx.To)

	case *obscurocommon.L1RollupTx:
		r, err := nodecommon.DecodeRollup(tx.Rollup)
		if err != nil {
			panic(err)
		}
		zipped, err := compress(tx.Rollup)
		if err != nil {
			panic(err)
		}
		encRollupData := encodeToString(zipped)
		data, err := contracts.MgmtContractABIJSON.Pack(contracts.AddRollupMethod, encRollupData)
		if err != nil {
			panic(err)
		}

		ethTx.Data = data
		log.Info(fmt.Sprintf("- Broadcasting - Issuing Rollup: r_%d - %d txs - datasize: %d - gas: %d",
			obscurocommon.ShortHash(r.Hash()), len(r.Transactions), len(data), ethTx.Gas))

	case *obscurocommon.L1StoreSecretTx:
		data, err := contracts.MgmtContractABIJSON.Pack(contracts.StoreSecretMethod, encodeToString(tx.Secret))
		if err != nil {
			panic(err)
		}
		ethTx.Data = data
		log.Info(fmt.Sprintf("- Broadcasting - Issuing StoreSecretTx: encoded as %s", encodeToString(tx.Secret)))
	case *obscurocommon.L1RequestSecretTx:
		data, err := contracts.MgmtContractABIJSON.Pack(contracts.RequestSecretMethod)
		if err != nil {
			panic(err)
		}
		ethTx.Data = data
		log.Info("- Broadcasting - Issuing RequestSecret")
	}

	return ethTx, nil
}

func (h *mgmtContractTxHandler) UnPack(tx *types.Transaction) obscurocommon.L1Transaction {
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
		zipped := decodeFromString(callData.(string))
		rollup, err := decompress(zipped)
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
			Secret: decodeFromString(callData.(string)),
		}

	case contracts.RequestSecretMethod:
		return &obscurocommon.L1RequestSecretTx{}
	}

	return nil
}

// encodeToString encodes a byte array to a string
func encodeToString(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

// decodeFromString decodes a string to a byte array
func decodeFromString(in string) []byte {
	bytesStr, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}
	return bytesStr
}

// compress the byte array using gzip
func compress(in []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(in); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// decompress the byte array using gzip
func decompress(in []byte) ([]byte, error) {
	reader := bytes.NewReader(in)
	gz, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	return ioutil.ReadAll(gz)
}
