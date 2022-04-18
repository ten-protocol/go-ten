package helpertypes

import (
	"bytes"
	"compress/gzip"
	b64 "encoding/base64"
	"io/ioutil"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/buildhelper/buildconstants"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

var IsRealEth bool

func UnpackL1Tx(tx *types.Transaction) *obscurocommon.L1TxData {
	if !IsRealEth {
		t := obscurocommon.TxData(tx)
		return &t
	}
	// ignore transactions that are not calling the contract
	if tx.To() == nil || tx.To().Hex() != buildconstants.CONTRACT_ADDRESS.Hex() || len(tx.Data()) == 0 {
		return nil
	}

	contractABI, err := abi.JSON(strings.NewReader(buildconstants.CONTRACT_ABI))
	if err != nil {
		panic(err)
	}

	//log.Log(fmt.Sprintf("Unpacking data: %b", tx.Data()))
	method, err := contractABI.MethodById(tx.Data()[:4])
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
			panic("call data not found for rollup")
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
			panic("call data not found for rollup")
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
			panic("call data not found for rollup")
		}
		l1txData.Secret = DecodeFromString(callData.(string))
		l1txData.TxType = obscurocommon.StoreSecretTx

	}

	return &l1txData
}

func EncodeToString(bytes []byte) string {
	return b64.StdEncoding.EncodeToString(bytes)
}

func DecodeFromString(in string) []byte {
	bytes, err := b64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}
	return bytes
}

func Compress(in []byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(in); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}
	return b.Bytes()
}

func Decompress(in []byte) []byte {
	reader := bytes.NewReader(in)
	gz, err := gzip.NewReader(reader)
	if err != nil {
		panic(err)
	}

	output, err := ioutil.ReadAll(gz)
	if err != nil {
		panic(err)
	}

	return output
}
