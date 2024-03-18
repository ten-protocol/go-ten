package test

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/vkhandler"
	"github.com/ten-protocol/go-ten/go/responses"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	l2ChainIDHex         = "0x309"
	l2ChainIDDecimal     = 443
	enclavePrivateKeyHex = "81acce9620f0adf1728cb8df7f6b8b8df857955eb9e8b7aed6ef8390c09fc207"
)

// DummyAPI provides dummies for the RPC operations defined in the `eth_` namespace. For each sensitive RPC
// operation, it decrypts the parameters using the enclave's private key, then echoes them back to the caller encrypted
// with the viewing key set using the `setViewingKey` method, mimicking the privacy behaviour of the host.
type DummyAPI struct {
	enclavePrivateKey *ecies.PrivateKey
	viewingKey        []byte
	signature         []byte
	address           *gethcommon.Address
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

// Determines which key the API will encrypt responses with.
func (api *DummyAPI) setViewingKey(address *gethcommon.Address, compressedVKKeyHexBytes, signature []byte) {
	api.viewingKey = compressedVKKeyHexBytes
	api.address = address
	api.signature = signature
}

func (api *DummyAPI) ChainId() (*hexutil.Big, error) { //nolint:stylecheck,revive
	chainID, err := hexutil.DecodeBig(l2ChainIDHex)
	return (*hexutil.Big)(chainID), err
}

func (api *DummyAPI) Call(_ context.Context, encryptedParams common.EncryptedParamsCall) (*responses.EnclaveResponse, error) {
	return api.reEncryptParams(encryptedParams)
}

func (api *DummyAPI) GetBalance(_ context.Context, encryptedParams common.EncryptedParamsGetBalance) (*responses.EnclaveResponse, error) {
	return api.reEncryptParams(encryptedParams)
}

func (api *DummyAPI) GetTransactionByHash(_ context.Context, encryptedParams common.EncryptedParamsGetTxByHash) (*responses.EnclaveResponse, error) {
	reEncryptParams, err := api.reEncryptParams(encryptedParams)
	return reEncryptParams, err
}

func (api *DummyAPI) GetTransactionCount(_ context.Context, encryptedParams common.EncryptedParamsGetTxCount) (*responses.EnclaveResponse, error) {
	return api.reEncryptParams(encryptedParams)
}

func (api *DummyAPI) GetTransactionReceipt(_ context.Context, encryptedParams common.EncryptedParamsGetTxReceipt) (*responses.EnclaveResponse, error) {
	reEncryptParams, err := api.reEncryptParams(encryptedParams)
	return reEncryptParams, err
}

func (api *DummyAPI) SendRawTransaction(_ context.Context, encryptedParams common.EncryptedParamsSendRawTx) (*responses.EnclaveResponse, error) {
	return api.reEncryptParams(encryptedParams)
}

func (api *DummyAPI) EstimateGas(_ context.Context, encryptedParams common.EncryptedParamsEstimateGas, _ *rpc.BlockNumberOrHash) (*responses.EnclaveResponse, error) {
	reEncryptParams, err := api.reEncryptParams(encryptedParams)
	return reEncryptParams, err
}

func (api *DummyAPI) Logs(ctx context.Context, encryptedParams common.EncryptedParamsLogSubscription) (*rpc.Subscription, error) {
	// We decrypt and decode the params.
	encodedParams, err := api.enclavePrivateKey.Decrypt(encryptedParams, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt params with enclave private key. Cause: %w", err)
	}
	var params common.LogSubscription
	if err = rlp.DecodeBytes(encodedParams, &params); err != nil {
		return nil, fmt.Errorf("could not decocde log subscription request from RLP. Cause: %w", err)
	}

	// We set up the subscription.
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return nil, fmt.Errorf("creation of subscriptions is not supported")
	}
	subscription := notifier.CreateSubscription()
	err = notifier.Notify(subscription.ID, common.IDAndEncLog{
		SubID: subscription.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("could not send subscription ID to client on subscription %s", subscription.ID)
	}

	// We emit a unique log every ten milliseconds.
	go func() {
		idx := big.NewInt(0)
		for {
			// We create the logs
			logs := []*types.Log{{Topics: []gethcommon.Hash{
				// We set the topic from the filter as a topic in the response logs, so that we can check in the tests
				// that we are a) decrypting the params correctly, and b) returning the logs with the correct contents
				// via the wallet extension.
				params.Filter.Topics[0][0],
				// We also add an incrementing integer as a topic, so we can detect duplicate logs.
				gethcommon.BigToHash(idx),
			}}}
			jsonLogs, err := json.Marshal(logs)
			if err != nil {
				panic("could not marshal log to JSON")
			}

			// We send the encrypted log via the subscription.
			pubkey, err := crypto.DecompressPubkey(api.viewingKey)
			if err != nil {
				panic("could not decompress Pub key")
			}

			encryptedBytes, err := ecies.Encrypt(rand.Reader, ecies.ImportECDSAPublic(pubkey), jsonLogs, nil, nil)
			if err != nil {
				panic("could not encrypt logs with viewing key")
			}
			idAndEncLog := common.IDAndEncLog{
				SubID:  subscription.ID,
				EncLog: encryptedBytes,
			}
			notifier.Notify(subscription.ID, idAndEncLog) //nolint:errcheck

			time.Sleep(10 * time.Millisecond)
			idx = idx.Add(idx, big.NewInt(1))
		}
	}()
	return subscription, nil
}

func (api *DummyAPI) GetLogs(_ context.Context, encryptedParams common.EncryptedParamsGetLogs) (*responses.EnclaveResponse, error) {
	reEncryptParams, err := api.reEncryptParams(encryptedParams)
	return reEncryptParams, err
}

func (api *DummyAPI) GetStorageAt(_ context.Context, encryptedParams common.EncryptedParamsSendRawTx) (*responses.EnclaveResponse, error) {
	return api.reEncryptParams(encryptedParams)
}

// Decrypts the params with the enclave key, and returns them encrypted with the viewing key set via `setViewingKey`.
func (api *DummyAPI) reEncryptParams(encryptedParams []byte) (*responses.EnclaveResponse, error) {
	params, err := api.enclavePrivateKey.Decrypt(encryptedParams, nil, nil)
	if err != nil {
		return responses.AsEmptyResponse(), fmt.Errorf("could not decrypt params with enclave private key. Cause: %w", err)
	}

	encryptor, err := vkhandler.VerifyViewingKey(&viewingkey.RPCSignedViewingKey{
		PublicKey:               api.viewingKey,
		SignatureWithAccountKey: api.signature,
		SignatureType:           viewingkey.Legacy, // todo - is this correct
	}, l2ChainIDDecimal)
	if err != nil {
		return nil, fmt.Errorf("unable to create vk encryption for request - %w", err)
	}

	strParams := string(params)

	return responses.AsEncryptedResponse(&strParams, encryptor), nil
}
