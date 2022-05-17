package txdecoder

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/contracts"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

const methodBytesLen = 4

// TxDecoder handles converting eth transactions in obscuro transactions
type TxDecoder interface {
	// DecodeTx receives a *types.Transaction and converts it to an obscurocommon.L1Transaction
	DecodeTx(tx *types.Transaction) obscurocommon.L1Transaction
}

type txDecoderImpl struct {
	mgmtContractAddr  *common.Address
	erc20ContractAddr *common.Address
}

func NewTxDecoder(mgmtContractAddr *common.Address, erc20ContractAddr *common.Address) TxDecoder {
	return &txDecoderImpl{
		mgmtContractAddr:  mgmtContractAddr,
		erc20ContractAddr: erc20ContractAddr,
	}
}

func (t *txDecoderImpl) DecodeTx(tx *types.Transaction) obscurocommon.L1Transaction {
	if tx.To() == nil || len(tx.Data()) == 0 {
		return nil
	}

	// route the transaction based on the address
	switch tx.To().Hex() {
	case t.mgmtContractAddr.Hex():
		return t.unpackMgmt(tx)
	case t.erc20ContractAddr.Hex():
		return t.unpackERC20(tx)
	default:
		return nil
	}
}

// unpackMgmt converts eth transaction to obscurocommon.L1Transaction
func (t *txDecoderImpl) unpackMgmt(tx *types.Transaction) obscurocommon.L1Transaction {
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
		zipped := base64DecodeFromString(callData.(string))
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
			Secret: base64DecodeFromString(callData.(string)),
		}

	case contracts.RequestSecretMethod:
		return &obscurocommon.L1RequestSecretTx{}
	}

	return nil
}

// unpackERC20 converts eth transaction to obscurocommon.L1Transaction
func (t *txDecoderImpl) unpackERC20(tx *types.Transaction) obscurocommon.L1Transaction {
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
		TokenContract: t.erc20ContractAddr,
	}
}

// base64DecodeFromString decodes a string to a byte array
func base64DecodeFromString(in string) []byte {
	bytesStr, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}
	return bytesStr
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
