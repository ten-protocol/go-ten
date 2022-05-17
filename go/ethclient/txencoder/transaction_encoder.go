package txencoder

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/contracts"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

var (
	// TODO review gas estimation - these should not be static values
	// Gas should be calculated so to not overpay what the operation requires
	// The values are hardcoded at the moment to guarantee the txs will be minted
	// It's using large gas values because rollups can be very expensive
	defaultGasPrice = big.NewInt(20000000000)
	defaultGas      = uint64(1024_000_000)
)

// TxEncoder encodes a L1 transaction into an eth transaction
type TxEncoder interface {
	CreateRollup(t *obscurocommon.L1RollupTx, nonce uint64) types.TxData
	CreateRequestSecret(tx *obscurocommon.L1RequestSecretTx, nonce uint64) types.TxData
	CreateStoreSecret(tx *obscurocommon.L1StoreSecretTx, nonce uint64) types.TxData
	CreateDepositTx(tx *obscurocommon.L1DepositTx, nonce uint64) types.TxData
}

// mgmtTxEncoderImpl implements the TxEncoder by converting obscurocommon.L1Transactions to ethereum transactions
// It does this be calling the Management Contract
type mgmtTxEncoderImpl struct {
	addr *common.Address
}

func NewEncoder(addr *common.Address) TxEncoder {
	return &mgmtTxEncoderImpl{
		addr: addr,
	}
}

func (h *mgmtTxEncoderImpl) CreateRollup(t *obscurocommon.L1RollupTx, nonce uint64) types.TxData {
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
		To:       h.addr,
		Data:     data,
	}
}

func (h *mgmtTxEncoderImpl) CreateRequestSecret(tx *obscurocommon.L1RequestSecretTx, nonce uint64) types.TxData {
	data, err := contracts.MgmtContractABIJSON.Pack(contracts.RequestSecretMethod)
	if err != nil {
		panic(err)
	}

	return &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: defaultGasPrice,
		Gas:      defaultGas,
		To:       h.addr,
		Data:     data,
	}
}

func (h *mgmtTxEncoderImpl) CreateStoreSecret(tx *obscurocommon.L1StoreSecretTx, nonce uint64) types.TxData {
	data, err := contracts.MgmtContractABIJSON.Pack(contracts.StoreSecretMethod, base64EncodeToString(tx.Secret))
	if err != nil {
		panic(err)
	}
	return &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: defaultGasPrice,
		Gas:      defaultGas,
		To:       h.addr,
		Data:     data,
	}
}

func (h *mgmtTxEncoderImpl) CreateDepositTx(tx *obscurocommon.L1DepositTx, nonce uint64) types.TxData {
	data, err := contracts.MgmtContractABIJSON.Pack(contracts.DepositMethod, tx.To)
	if err != nil {
		panic(err)
	}

	return &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: defaultGasPrice,
		Gas:      defaultGas,
		To:       h.addr,
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
