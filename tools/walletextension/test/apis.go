package test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
)

const (
	l2ChainIDHex         = "0x309"
	enclavePrivateKeyHex = "81acce9620f0adf1728cb8df7f6b8b8df857955eb9e8b7aed6ef8390c09fc207"
	filter               = "Filter"
	filterKeyTopics      = "Topics"
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
		panic(fmt.Errorf("failed to create enclave private key. Cause: %w", err))
	}

	return &DummyAPI{
		enclavePrivateKey: ecies.ImportECDSA(enclavePrivateKey),
	}
}

func (api *DummyAPI) AddViewingKey([]byte, []byte) error {
	return nil
}

// Determines which key the API will encrypt responses with.
func (api *DummyAPI) setViewingKey(viewingKeyHexBytes []byte) {
	viewingKeyBytes, err := hex.DecodeString(string(viewingKeyHexBytes))
	if err != nil {
		panic(err)
	}

	viewingKey, err := crypto.DecompressPubkey(viewingKeyBytes)
	if err != nil {
		panic(fmt.Errorf("received viewing key bytes but could not decompress them. Cause: %w", err))
	}
	api.viewingKey = ecies.ImportECDSAPublic(viewingKey)
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

func (api *DummyAPI) Logs(ctx context.Context, encryptedParams common.EncryptedParamsLogSubscription) (*rpc.Subscription, error) {
	params, err := api.enclavePrivateKey.Decrypt(encryptedParams, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt params with enclave private key. Cause: %w", err)
	}

	// We set the topic from the filter as a topic in the response logs, to provide additional assurance that we are
	// a) decrypting the params correctly, and b) returning the logs correctly via the wallet extension.
	var paramsMap []map[string]interface{}
	err = json.Unmarshal(params, &paramsMap)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal params from JSON")
	}
	paramsTopic := paramsMap[0][filter].(map[string]interface{})[filterKeyTopics].([]interface{})[0].([]interface{})[0].(string)

	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return nil, fmt.Errorf("creation of subscriptions is not supported")
	}

	sub := notifier.CreateSubscription()
	go func() {
		for {
			jsonLogs, err := json.Marshal([]*types.Log{{Topics: []gethcommon.Hash{gethcommon.HexToHash(paramsTopic)}}})
			if err != nil {
				panic("could not marshal log to JSON")
			}

			encryptedBytes, err := ecies.Encrypt(rand.Reader, api.viewingKey, jsonLogs, nil, nil)
			if err != nil {
				panic("could not encrypt logs with viewing key")
			}

			// Like the host, we wrap the encrypted log bytes in the data field of an unencrypted "wrapper" log object.
			// todo - joel - update this
			wrapperLog := types.Log{
				Topics: []gethcommon.Hash{},
				Data:   encryptedBytes,
			}

			notifier.Notify(sub.ID, &wrapperLog) //nolint:errcheck
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return sub, nil
}

// Decrypts the params with the enclave key, and returns them encrypted with the viewing key set via `setViewingKey`.
func (api *DummyAPI) reEncryptParams(encryptedParams []byte) (string, error) {
	params, err := api.enclavePrivateKey.Decrypt(encryptedParams, nil, nil)
	if err != nil {
		return "", fmt.Errorf("could not decrypt params with enclave private key. Cause: %w", err)
	}

	encryptedBytes, err := ecies.Encrypt(rand.Reader, api.viewingKey, params, nil, nil)
	if err != nil {
		return "", fmt.Errorf("could not encrypt params with viewing key")
	}

	return gethcommon.Bytes2Hex(encryptedBytes), err
}
