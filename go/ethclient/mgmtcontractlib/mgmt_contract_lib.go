package mgmtcontractlib

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
	CreateStoreSecret(tx *obscurocommon.L1StoreSecretTx, nonce uint64) types.TxData
	CreateDepositTx(tx *obscurocommon.L1DepositTx, nonce uint64) types.TxData

	// DecodeTx receives a *types.Transaction and converts it to an obscurocommon.L1Transaction
	DecodeTx(tx *types.Transaction) obscurocommon.L1Transaction
}

type contractLibImpl struct {
	addr *common.Address
}

func NewMgmtContractLib(addr *common.Address) MgmtContractLib {
	return &contractLibImpl{
		addr: addr,
	}
}

func (c *contractLibImpl) DecodeTx(tx *types.Transaction) obscurocommon.L1Transaction {
	if tx.To() == nil || tx.To().Hex() != c.addr.Hex() || len(tx.Data()) == 0 {
		return nil
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

func (c *contractLibImpl) CreateRollup(t *obscurocommon.L1RollupTx, nonce uint64) types.TxData {
	zipped, err := compress(t.Rollup)
	if err != nil {
		panic(err)
	}
	encRollupData := base64EncodeToString(zipped)
	data, err := contracts.MgmtContractABIJSON.Pack(contracts.AddRollupMethod, encRollupData)
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
	data, err := contracts.MgmtContractABIJSON.Pack(contracts.RequestSecretMethod)
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

func (c *contractLibImpl) CreateStoreSecret(tx *obscurocommon.L1StoreSecretTx, nonce uint64) types.TxData {
	data, err := contracts.MgmtContractABIJSON.Pack(contracts.StoreSecretMethod, base64EncodeToString(tx.Secret))
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

func (c *contractLibImpl) CreateDepositTx(tx *obscurocommon.L1DepositTx, nonce uint64) types.TxData {
	data, err := contracts.MgmtContractABIJSON.Pack(contracts.DepositMethod, tx.To)
	if err != nil {
		panic(err)
	}

	return &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: defaultGasPrice,
		Gas:      defaultGas,
		To:       c.addr,
		Data:     data,
		Value:    big.NewInt(int64(tx.Amount)),
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
