package mgmtcontractlib

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"math/big"
	"strings"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

const methodBytesLen = 4

var (
	// TODO review gas estimation - these should not be static values
	// Gas should be calculated so to not overpay what the operation requires
	// The values are hardcoded at the moment to guarantee the txs will be minted
	// It's using large gas values because rollups can be very expensive
	defaultGasPrice = big.NewInt(20000000000)
	defaultGas      = uint64(1024_000_000)
)

// MgmtContractLib provide methods for creating ethereum transactions by providing a L1Transaction
// Also provides a method to convert ethereum transactions into a L1Transaction
type MgmtContractLib interface {
	CreateRollup(t *obscurocommon.L1RollupTx, nonce uint64) types.TxData
	CreateRequestSecret(tx *obscurocommon.L1RequestSecretTx, nonce uint64) types.TxData
	CreateRespondSecret(tx *obscurocommon.L1RespondSecretTx, nonce uint64) types.TxData
	CreateInitializeSecret(tx *obscurocommon.L1InitializeSecretTx, nonce uint64) types.TxData

	// DecodeTx receives a *types.Transaction and converts it to an obscurocommon.L1Transaction
	DecodeTx(tx *types.Transaction) obscurocommon.L1Transaction
}

type contractLibImpl struct {
	addr        *common.Address
	contractABI abi.ABI
}

func NewMgmtContractLib(addr *common.Address) MgmtContractLib {
	contractABI, err := abi.JSON(strings.NewReader(MgmtContractABI))
	if err != nil {
		panic(err)
	}

	return &contractLibImpl{
		addr:        addr,
		contractABI: contractABI,
	}
}

func (c *contractLibImpl) DecodeTx(tx *types.Transaction) obscurocommon.L1Transaction {
	if tx.To() == nil || tx.To().Hex() != c.addr.Hex() || len(tx.Data()) == 0 {
		return nil
	}
	method, err := c.contractABI.MethodById(tx.Data()[:methodBytesLen])
	if err != nil {
		panic(err)
	}

	contractCallData := map[string]interface{}{}
	switch method.Name {
	case AddRollupMethod:
		if err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[4:]); err != nil {
			panic(err)
		}
		callData, found := contractCallData["rollupData"]
		if !found {
			panic("call data not found for rollupData")
		}
		zipped := Base64DecodeFromString(callData.(string))
		rollup, err := Decompress(zipped)
		if err != nil {
			panic(err)
		}

		return &obscurocommon.L1RollupTx{
			Rollup: rollup,
		}

	case RespondSecretMethod:
		return unpackRespondSecretTx(tx, method, contractCallData)

	case RequestSecretMethod:
		return unpackRequestSecretTx(tx, method, contractCallData)
	}

	return nil
}

func (c *contractLibImpl) CreateRollup(t *obscurocommon.L1RollupTx, nonce uint64) types.TxData {
	decodedRollup, err := nodecommon.DecodeRollup(t.Rollup)
	if err != nil {
		panic(err)
	}

	zipped, err := compress(t.Rollup)
	if err != nil {
		panic(err)
	}
	encRollupData := base64EncodeToString(zipped)

	data, err := c.contractABI.Pack(
		AddRollupMethod,
		decodedRollup.Header.ParentHash,
		decodedRollup.Hash(),
		decodedRollup.Header.Agg,
		decodedRollup.Header.L1Proof,
		decodedRollup.Header.Number,
		encRollupData,
	)
	if err != nil {
		panic(err)
	}

	return &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: defaultGasPrice,
		Gas:      defaultGas,
		To:       c.addr,
		Data:     data,
	}
}

func (c *contractLibImpl) CreateRequestSecret(tx *obscurocommon.L1RequestSecretTx, nonce uint64) types.TxData {
	data, err := c.contractABI.Pack(RequestSecretMethod, base64EncodeToString(tx.Attestation))
	if err != nil {
		panic(err)
	}

	return &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: defaultGasPrice,
		Gas:      defaultGas,
		To:       c.addr,
		Data:     data,
	}
}

func (c *contractLibImpl) CreateRespondSecret(tx *obscurocommon.L1RespondSecretTx, nonce uint64) types.TxData {
	data, err := c.contractABI.Pack(
		RespondSecretMethod,
		tx.AttesterID,
		tx.RequesterID,
		tx.AttesterSig,
		tx.Secret,
	)
	if err != nil {
		panic(err)
	}
	return &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: defaultGasPrice,
		Gas:      defaultGas,
		To:       c.addr,
		Data:     data,
	}
}

func (c *contractLibImpl) CreateInitializeSecret(tx *obscurocommon.L1InitializeSecretTx, nonce uint64) types.TxData {
	data, err := c.contractABI.Pack(
		InitializeSecretMethod,
		tx.AggregatorID,
		tx.InitialSecret,
	)
	if err != nil {
		panic(err)
	}
	return &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: defaultGasPrice,
		Gas:      defaultGas,
		To:       c.addr,
		Data:     data,
	}
}

func unpackRequestSecretTx(tx *types.Transaction, method *abi.Method, contractCallData map[string]interface{}) *obscurocommon.L1RequestSecretTx {
	err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[methodBytesLen:])
	if err != nil {
		panic(err)
	}
	callData, found := contractCallData["requestReport"]
	if !found {
		panic("call data not found for requestReport")
	}

	att := Base64DecodeFromString(callData.(string))
	if err != nil {
		log.Panic("could not decode attestation request. Cause: %s", err)
	}
	return &obscurocommon.L1RequestSecretTx{
		Attestation: att,
	}
}

func unpackRespondSecretTx(tx *types.Transaction, method *abi.Method, contractCallData map[string]interface{}) *obscurocommon.L1RespondSecretTx {
	err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[methodBytesLen:])
	if err != nil {
		log.Panic("could not unpack transaction. Cause: %s", err)
	}
	requesterData, found := contractCallData["requesterID"]
	if !found {
		log.Panic("call data not found for requesterID")
	}

	requesterAddr, ok := requesterData.(common.Address)
	if !ok {
		log.Panic("could not decode requester data")
	}

	attesterData, found := contractCallData["attesterID"]
	if !found {
		log.Panic("call data not found for attesterID")
	}

	attesterAddr, ok := attesterData.(common.Address)
	if !ok {
		log.Panic("could not decode attester data")
	}

	responseSecretData, found := contractCallData["responseSecret"]
	if !found {
		log.Panic("call data not found for inputSecret")
	}
	responseSecretBytes, ok := responseSecretData.([]uint8)
	if !ok {
		log.Panic("could not decode requester responseSecret data")
	}

	return &obscurocommon.L1RespondSecretTx{
		AttesterID:  attesterAddr,
		RequesterID: requesterAddr,
		Secret:      responseSecretBytes[:],
	}
}

// base64EncodeToString encodes a byte array to a string
func base64EncodeToString(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
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

// Base64DecodeFromString decodes a string to a byte array
func Base64DecodeFromString(in string) []byte {
	bytesStr, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}
	return bytesStr
}

// Decompress the byte array using gzip
func Decompress(in []byte) ([]byte, error) {
	reader := bytes.NewReader(in)
	gz, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	return ioutil.ReadAll(gz)
}
