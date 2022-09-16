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
	successMsg           = "success"
	l2ChainIDHex         = "0x309"
	enclavePrivateKeyHex = "81acce9620f0adf1728cb8df7f6b8b8df857955eb9e8b7aed6ef8390c09fc207"
)

// DummyAPI provides dummies for the RPC operations defined in the `eth_` namespace. For each sensitive RPC
// operation, it decrypts the parameters using the enclave's private key, then echoes them back to the caller encrypted
// with the viewing key set using the `setViewingKey` method, mimicking the privacy behaviour of the host.
type DummyAPI struct {
	enclavePrivateKey *ecies.PrivateKey
	viewingKey        *ecies.PublicKey
}

func NewDummyAPI() *DummyAPI {
	enclavePrivateKey, err := crypto.HexToECDSA(enclavePrivateKeyHex)
	if err != nil {
		panic(fmt.Errorf("failed to create enclave private key. Cause: %s", err))
	}

	return &DummyAPI{
		enclavePrivateKey: ecies.ImportECDSA(enclavePrivateKey),
	}
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

func (api *DummyAPI) Call(_ context.Context, encryptedParams common.EncryptedParamsCall) (string, error) {
	return api.reEncryptParams(encryptedParams)
}

func (api *DummyAPI) GetBalance(_ context.Context, encryptedParams common.EncryptedParamsGetBalance) (string, error) {
	return api.reEncryptParams(encryptedParams)
}

func (api *DummyAPI) GetTransactionByHash(_ context.Context, encryptedParams common.EncryptedParamsGetTxByHash) (*string, error) {
	reEncryptParams, err := api.reEncryptParams(encryptedParams)
	return &reEncryptParams, err
}

func (api *DummyAPI) GetTransactionCount(_ context.Context, encryptedParams common.EncryptedParamsGetTxCount) (string, error) {
	return api.reEncryptParams(encryptedParams)
}

func (api *DummyAPI) GetTransactionReceipt(_ context.Context, encryptedParams common.EncryptedParamsGetTxReceipt) (*string, error) {
	reEncryptParams, err := api.reEncryptParams(encryptedParams)
	return &reEncryptParams, err
}

func (api *DummyAPI) SendRawTransaction(_ context.Context, encryptedParams common.EncryptedParamsSendRawTx) (string, error) {
	return api.reEncryptParams(encryptedParams)
}

func (api *DummyAPI) EstimateGas(_ context.Context, encryptedParams common.EncryptedParamsEstimateGas, _ *rpc.BlockNumberOrHash) (*string, error) {
	reEncryptParams, err := api.reEncryptParams(encryptedParams)
	return &reEncryptParams, err
}

// Returns the message `successMsg`, encrypted with the viewing key set via `setViewingKey`.
func (api *DummyAPI) reEncryptParams(encryptedParams []byte) (string, error) {
	params, err := api.enclavePrivateKey.Decrypt(encryptedParams, nil, nil)
	if err != nil {
		return "", fmt.Errorf("could not decrypt params with enclave private key")
	}

	encryptedBytes, err := ecies.Encrypt(rand.Reader, api.viewingKey, params, nil, nil)
	if err != nil {
		return "", fmt.Errorf("could not encrypt params with viewing key")
	}

	return gethcommon.Bytes2Hex(encryptedBytes), err
}
