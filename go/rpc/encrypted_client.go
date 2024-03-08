package rpc

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ten-protocol/go-ten/go/common/rpc"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rlp"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/go/responses"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

const (
	// todo: this is a convenience for testnet testing and will eventually be retrieved from the L1
	enclavePublicKeyHex = "034d3b7e63a8bcd532ee3d1d6ecad9d67fca7821981a044551f0f0cbec74d0bc5e"
	emptyFilterCriteria = "[]" // This is the value that gets passed for an empty filter criteria.
)

// SensitiveMethods for which the RPC requests and responses should be encrypted
var SensitiveMethods = []string{
	Call,
	GetBalance,
	GetTransactionByHash,
	GetTransactionCount,
	GetTransactionReceipt,
	SendRawTransaction,
	EstimateGas,
	GetLogs,
	GetStorageAt,
}

// EncRPCClient is a Client wrapper that implements Client but also has extra functionality for managing viewing key registration and decryption
type EncRPCClient struct {
	obscuroClient    Client
	enclavePublicKey *ecies.PublicKey // Used to encrypt messages destined to the enclave.
	viewingKey       *viewingkey.ViewingKey
	logger           gethlog.Logger
}

// NewEncRPCClient sets up a client with a viewing key for encrypted communication
func NewEncRPCClient(client Client, viewingKey *viewingkey.ViewingKey, logger gethlog.Logger) (*EncRPCClient, error) {
	// todo: this is a convenience for testnet but needs to replaced by a parameter and/or retrieved from the target host
	enclPubECDSA, err := crypto.DecompressPubkey(gethcommon.Hex2Bytes(enclavePublicKeyHex))
	if err != nil {
		return nil, fmt.Errorf("failed to decompress key for RPC client: %w", err)
	}
	enclavePublicKey := ecies.ImportECDSAPublic(enclPubECDSA)

	encClient := &EncRPCClient{
		obscuroClient:    client,
		enclavePublicKey: enclavePublicKey,
		viewingKey:       viewingKey,
		logger:           logger,
	}

	return encClient, nil
}

// Call handles JSON rpc requests without a context - see CallContext for details
func (c *EncRPCClient) Call(result interface{}, method string, args ...interface{}) error {
	return c.CallContext(nil, result, method, args...) //nolint:staticcheck
}

// CallContext is the main logic to execute JSON-RPC requests, the context can be nil.
// - if the method is sensitive it will encrypt the args before sending the request and then decrypts the response before returning
// - result must be a pointer so that package json can unmarshal into it. You can also pass nil, in which case the result is ignored.
// - callExec handles the delegated call, allows EncClient to use the same code for calling with or without a context
func (c *EncRPCClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	assertResultIsPointer(result)
	if !IsSensitiveMethod(method) {
		// for non-sensitive methods or when viewing keys are disabled we just delegate directly to the geth RPC client
		return c.executeRPCCall(ctx, result, method, args...)
	}

	return c.executeSensitiveCall(ctx, result, method, args...)
}

func (c *EncRPCClient) Subscribe(ctx context.Context, _ interface{}, namespace string, ch interface{}, args ...interface{}) (*gethrpc.ClientSubscription, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("subscription did not specify its type")
	}

	subscriptionType := args[0]
	if subscriptionType != SubscriptionTypeLogs {
		return nil, fmt.Errorf("only subscriptions of type %s are supported", SubscriptionTypeLogs)
	}

	logSubscription, err := c.createAuthenticatedLogSubscription(args)
	if err != nil {
		return nil, err
	}

	// We use RLP instead of JSON marshaling here, as for some reason the filter criteria doesn't unmarshal correctly from JSON.
	encodedLogSubscription, err := rlp.EncodeToBytes(logSubscription)
	if err != nil {
		return nil, err
	}

	encryptedParams, err := c.encryptParamBytes(encodedLogSubscription)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt args for subscription in namespace %s - %w", namespace, err)
	}

	logCh, ok := ch.(chan common.IDAndLog)
	if !ok {
		return nil, fmt.Errorf("expected a channel of type `chan types.Log`, got %T", ch)
	}
	clientChannel := make(chan common.IDAndEncLog)
	subscriptionToObscuro, err := c.obscuroClient.Subscribe(ctx, nil, namespace, clientChannel, subscriptionType, encryptedParams)
	if err != nil {
		return nil, err
	}

	go c.forwardLogs(clientChannel, logCh, subscriptionToObscuro)

	return subscriptionToObscuro, nil
}

func (c *EncRPCClient) forwardLogs(clientChannel chan common.IDAndEncLog, logCh chan common.IDAndLog, subscription *gethrpc.ClientSubscription) {
	for {
		select {
		case idAndEncLog := <-clientChannel:
			jsonLogs, err := c.decryptResponse(idAndEncLog.EncLog)
			if err != nil {
				c.logger.Error("could not decrypt logs received from subscription.", log.ErrKey, err)
				continue
			}

			var logs []*types.Log
			err = json.Unmarshal(jsonLogs, &logs)
			if err != nil {
				c.logger.Error(fmt.Sprintf("could not unmarshal log from JSON. Received data: %s.", string(jsonLogs)), log.ErrKey, err)
				continue
			}

			for _, decryptedLog := range logs {
				idAndLog := common.IDAndLog{
					SubID: idAndEncLog.SubID,
					Log:   decryptedLog,
				}
				logCh <- idAndLog
			}

		case err := <-subscription.Err():
			if err != nil {
				c.logger.Info("subscription to obscuro node closed with error", log.ErrKey, err)
			} else {
				c.logger.Info("subscription to obscuro node closed")
			}
			return
		}
	}
}

func (c *EncRPCClient) createAuthenticatedLogSubscription(args []interface{}) (*common.LogSubscription, error) {
	logSubscription := &common.LogSubscription{
		ViewingKey: &viewingkey.RPCSignedViewingKey{
			PublicKey:               c.viewingKey.PublicKey,
			SignatureWithAccountKey: c.viewingKey.SignatureWithAccountKey,
			SignatureType:           c.viewingKey.SignatureType,
		},
	}

	// If there are less than two arguments, it means no filter criteria was passed.
	if len(args) < 2 {
		logSubscription.Filter = &filters.FilterCriteria{}
		return logSubscription, nil
	}

	// TODO - Consider switching to using the common.FilterCriteriaJSON type. Should allow us to avoid RLP serialisation.
	// We marshal the filter criteria from a map to JSON, then back from JSON into a FilterCriteria. This is
	// because the filter criteria arrives as a map, and there is no way to convert it from a map directly into a
	// FilterCriteria.
	filterCriteriaJSON, err := json.Marshal(args[1])
	if err != nil {
		return nil, fmt.Errorf("could not marshal filter criteria to JSON. Cause: %w", err)
	}

	filterCriteria := filters.FilterCriteria{}
	if string(filterCriteriaJSON) != emptyFilterCriteria {
		err = filterCriteria.UnmarshalJSON(filterCriteriaJSON)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal filter criteria from the following JSON: `%s`. Cause: %w", string(filterCriteriaJSON), err)
		}
	}

	// If we do not override a nil block hash to an empty one, RLP decoding will fail on the enclave side.
	if filterCriteria.BlockHash == nil {
		filterCriteria.BlockHash = &gethcommon.Hash{}
	}

	logSubscription.Filter = &filterCriteria
	return logSubscription, nil
}

func (c *EncRPCClient) executeSensitiveCall(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	// encode the params into a json blob and encrypt them
	encryptedParams, err := c.encryptArgs(args...)
	if err != nil {
		return fmt.Errorf("failed to encrypt args for %s call - %w", method, err)
	}

	// We setup the rawResult to receive an EnclaveResponse. All sensitive methods should return this
	var rawResult responses.EnclaveResponse
	err = c.executeRPCCall(ctx, &rawResult, method, encryptedParams)
	if err != nil {
		return err
	}

	// if caller not interested in response, we're done
	if result == nil {
		return nil
	}

	// If the enclave has produced a plaintext error we give the
	// plaintext error back
	if rawResult.Error() != nil {
		return rawResult.Error()
	}

	// If there is no encrypted response then this is equivalent to nil response
	if rawResult.EncUserResponse == nil || len(rawResult.EncUserResponse) == 0 {
		return ErrNilResponse
	}

	// We decrypt the user response from the enclave response.
	decrypted, err := c.decryptResponse(rawResult.EncUserResponse)
	if err != nil {
		return fmt.Errorf("could not decrypt response for %s call - %w", method, err)
	}

	// We decode the UserResponse but keep the result as a json object
	// this method returns the user error if any and the result encoded as json.
	decodedResult, decodedError := responses.DecodeResponse[json.RawMessage](decrypted)

	// If there is a user error that was decrypted we return it
	if decodedError != nil {
		// EstimateGas and Call methods return EVM Errors that are json objects
		// and contain multiple keys that normally do not get serialized
		if method == EstimateGas || method == Call {
			var result errutil.EVMSerialisableError
			err = json.Unmarshal([]byte(decodedError.Error()), &result)
			if err != nil {
				return err
			}
			// Return the evm user error.
			return result
		}

		// Return the user error.
		return decodedError
	}

	if decodedResult == nil {
		return nil
	}

	// We get the bytes behind the raw json object.
	// note that RawJson messages simply return the bytes
	// and never error.
	resultBytes, _ := decodedResult.MarshalJSON()

	// We put the raw json in the passed result object.
	// This works for structs, strings, integers and interface types.
	err = json.Unmarshal(resultBytes, result)
	if err != nil {
		return fmt.Errorf("could not populate the response object with the json_rpc result. Cause: %w", err)
	}

	return nil
}

func (c *EncRPCClient) executeRPCCall(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	if ctx == nil {
		return c.obscuroClient.Call(result, method, args...)
	}
	return c.obscuroClient.CallContext(ctx, result, method, args...)
}

func (c *EncRPCClient) Stop() {
	c.obscuroClient.Stop()
}

func (c *EncRPCClient) Account() *gethcommon.Address {
	return c.viewingKey.Account
}

func (c *EncRPCClient) encryptArgs(args ...interface{}) ([]byte, error) {
	if len(args) == 0 {
		return nil, nil
	}
	vk := viewingkey.RPCSignedViewingKey{
		PublicKey:               c.viewingKey.PublicKey,
		SignatureWithAccountKey: c.viewingKey.SignatureWithAccountKey,
		SignatureType:           c.viewingKey.SignatureType,
	}
	argsWithVK := &rpc.RequestWithVk{VK: &vk, Params: args}

	paramsJSON, err := json.Marshal(argsWithVK)
	if err != nil {
		return nil, fmt.Errorf("could not json encode request params: %w", err)
	}

	return c.encryptParamBytes(paramsJSON)
}

func (c *EncRPCClient) encryptParamBytes(params []byte) ([]byte, error) {
	encryptedParams, err := ecies.Encrypt(rand.Reader, c.enclavePublicKey, params, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt the following request params with enclave public key: %s. Cause: %w", params, err)
	}
	return encryptedParams, nil
}

func (c *EncRPCClient) decryptResponse(encryptedBytes []byte) ([]byte, error) {
	decryptedResult, err := c.viewingKey.PrivateKey.Decrypt(encryptedBytes, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt bytes with viewing key. Cause: %w. Bytes: %s", err, string(encryptedBytes))
	}
	return decryptedResult, nil
}

// IsSensitiveMethod indicates whether the RPC method's requests and responses should be encrypted.
func IsSensitiveMethod(method string) bool {
	for _, m := range SensitiveMethods {
		if m == method {
			return true
		}
	}
	return false
}

func assertResultIsPointer(result interface{}) {
	// result MUST be an initialized pointer else call won't be able to return it
	if result != nil {
		// todo: replace these panics with an error for invalid usage (same behaviour as json.Unmarshal())
		if reflect.ValueOf(result).Kind() != reflect.Ptr {
			// we panic if result is not a pointer, this is a coding mistake and we want to fail fast during development
			panic("result MUST be a pointer else Call cannot populate it")
		}
		if reflect.ValueOf(result).IsNil() {
			// we panic if result is a nil pointer, cannot unmarshal json to it. Pointer must be initialized.
			// if you see this then the calling code probably used: `var resObj *ResType` instead of: `var resObj ResType`
			panic("result pointer must be initialized else Call cannot populate it")
		}
	}
}
