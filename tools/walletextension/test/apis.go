package test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
)

const (
	successMsg   = "success"
	l2ChainIDHex = "0x309"
)

// DummyAPI provides dummies for the RPC operations defined in the `eth_` namespace. For each sensitive RPC
// operation, it decrypts the parameters using the enclave's private key, then echoes them back to the caller encrypted
// with the viewing key set using the `setViewingKey` method.
type DummyAPI struct {
	viewingKey *ecies.PublicKey
}

func (api *DummyAPI) AddViewingKey([]byte, []byte) error {
	return nil
}

// Determines which key the API will encrypt responses with.
func (api *DummyAPI) setViewingKey(viewingKeyHexBytes []byte) error {
	viewingKeyBytes, err := hex.DecodeString(string(viewingKeyHexBytes))
	if err != nil {
		return err
	}

	viewingKey, err := crypto.DecompressPubkey(viewingKeyBytes)
	if err != nil {
		return fmt.Errorf("received viewing key bytes but could not decompress them. Cause: %w", err)
	}
	api.viewingKey = ecies.ImportECDSAPublic(viewingKey)
	return nil
}

func (api *DummyAPI) ChainId() (*hexutil.Big, error) { //nolint:stylecheck,revive
	chainID, err := hexutil.DecodeBig(l2ChainIDHex)
	return (*hexutil.Big)(chainID), err
}

func (api *DummyAPI) Call(context.Context, common.EncryptedParamsCall) (string, error) {
	return api.encryptedSuccess()
}

func (api *DummyAPI) GetBalance(context.Context, common.EncryptedParamsGetBalance) (string, error) {
	return api.encryptedSuccess()
}

func (api *DummyAPI) GetTransactionByHash(context.Context, common.EncryptedParamsGetTxByHash) (*string, error) {
	encryptedSuccess, err := api.encryptedSuccess()
	return &encryptedSuccess, err
}

func (api *DummyAPI) GetTransactionCount(context.Context, common.EncryptedParamsGetTxCount) (string, error) {
	return api.encryptedSuccess()
}

func (api *DummyAPI) GetTransactionReceipt(context.Context, common.EncryptedParamsGetTxReceipt) (*string, error) {
	encryptedSuccess, err := api.encryptedSuccess()
	return &encryptedSuccess, err
}

func (api *DummyAPI) SendRawTransaction(context.Context, common.EncryptedParamsSendRawTx) (string, error) {
	return api.encryptedSuccess()
}

func (api *DummyAPI) EstimateGas(context.Context, common.EncryptedParamsEstimateGas, *rpc.BlockNumberOrHash) (*string, error) {
	encryptedSuccess, err := api.encryptedSuccess()
	return &encryptedSuccess, err
}

// Returns the message `successMsg`, encrypted with the viewing key set via `setViewingKey`.
func (api *DummyAPI) encryptedSuccess() (string, error) {
	encryptedBytes, err := ecies.Encrypt(rand.Reader, api.viewingKey, []byte(successMsg), nil, nil)
	return gethcommon.Bytes2Hex(encryptedBytes), err
}
