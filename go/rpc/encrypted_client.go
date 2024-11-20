package rpc

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/ten-protocol/go-ten/go/common/rpc"
	"github.com/ten-protocol/go-ten/go/common/subscription"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/go/responses"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

const (
	// todo: this is a convenience for testnet testing and will eventually be retrieved from the L1
	enclavePublicKeyHex = "034d3b7e63a8bcd532ee3d1d6ecad9d67fca7821981a044551f0f0cbec74d0bc5e"
)

// EncRPCClient is a Client wrapper that implements Client but also has extra functionality for managing viewing key registration and decryption
type EncRPCClient struct {
	obscuroClient    Client
	enclavePublicKey *ecies.PublicKey // Used to encrypt messages destined to the enclave.
	viewingKey       *viewingkey.ViewingKey
	logger           gethlog.Logger
}

// NewEncRPCClient wrapper over rpc clients with a viewing key for encrypted communication
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

func (c *EncRPCClient) BackingClient() Client {
	return c.obscuroClient
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
	// borrowed from geth
	if result != nil && reflect.TypeOf(result).Kind() != reflect.Ptr {
		return fmt.Errorf("call result parameter must be pointer or nil interface: %v", result)
	}

	if rpc.IsEncryptedMethod(method) {
		return c.executeEncryptedCall(ctx, result, method, args...)
	}

	// for non-sensitive methods or when viewing keys are disabled we just delegate directly to the geth RPC client
	return c.executeRPCCall(ctx, result, method, args...)
}

func (c *EncRPCClient) Subscribe(ctx context.Context, namespace string, ch interface{}, args ...interface{}) (*gethrpc.ClientSubscription, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("missing subscription type")
	}

	switch args[0] {
	case SubscriptionTypeLogs:
		return c.logSubscription(ctx, namespace, ch, args...)
	case SubscriptionTypeNewHeads:
		return c.newHeadSubscription(ctx, namespace, ch, args...)
	default:
		return nil, fmt.Errorf("only subscriptions of type %s and %s are supported", SubscriptionTypeLogs, SubscriptionTypeNewHeads)
	}
}

func (c *EncRPCClient) executeEncryptedCall(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	// encode the params into a json blob and encrypt them
	encryptedParams, err := c.encryptArgs(method, args...)
	if err != nil {
		return fmt.Errorf("failed to encrypt args for %s call - %w", method, err)
	}

	// We setup the rawResult to receive an EnclaveResponse. All sensitive methods should return this
	var rawResult responses.EnclaveResponse
	err = c.executeRPCCall(ctx, &rawResult, "ten_encryptedRPC", encryptedParams)
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
	if len(rawResult.EncUserResponse) == 0 {
		return nil
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
		return decodedError
	}

	if decodedResult == nil {
		return nil
	}

	// We get the bytes behind the raw json object.
	// note that RawJson messages simply return the bytes
	// and never error.
	resultBytes, _ := decodedResult.MarshalJSON()

	// if expected result type is bytes, we return the bytes
	if _, ok := result.(*[]byte); ok {
		*result.(*[]byte) = resultBytes
		return nil
	}

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

func (c *EncRPCClient) encryptArgs(method string, args ...interface{}) ([]byte, error) {
	if len(args) == 0 {
		return nil, nil
	}
	vk := viewingkey.RPCSignedViewingKey{
		PublicKey:               c.viewingKey.PublicKey,
		SignatureWithAccountKey: c.viewingKey.SignatureWithAccountKey,
		SignatureType:           c.viewingKey.SignatureType,
	}
	argsWithVK := &rpc.RequestWithVk{VK: &vk, Method: method, Params: args}

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

// creates a subscription to the TEN node, but decrypts the messages from that channel and forwards them to the `ch`
func (c *EncRPCClient) logSubscription(ctx context.Context, namespace string, ch interface{}, args ...any) (*gethrpc.ClientSubscription, error) {
	outboundChannel, ok := ch.(chan types.Log)
	if !ok {
		return nil, fmt.Errorf("expected a channel of type `chan types.Log`, got %T", ch)
	}

	logSubscription, err := common.CreateAuthenticatedLogSubscriptionPayload(args, c.viewingKey)
	if err != nil {
		return nil, err
	}

	encodedLogSubscription, err := json.Marshal(logSubscription)
	if err != nil {
		return nil, err
	}

	encryptedParams, err := c.encryptParamBytes(encodedLogSubscription)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt args for subscription in namespace %s - %w", namespace, err)
	}

	// the node sends encrypted logs
	inboundChannel := make(chan []byte)
	backendSub, err := c.obscuroClient.Subscribe(ctx, namespace, inboundChannel, SubscriptionTypeLogs, encryptedParams)
	if err != nil {
		return nil, err
	}

	backendDisconnected := &atomic.Bool{}
	go subscription.HandleUnsubscribeErrChan([]<-chan error{backendSub.Err()}, func() {
		backendDisconnected.Store(true)
	})
	go subscription.ForwardFromChannels(
		[]chan []byte{inboundChannel},
		func(encLog []byte) error {
			return c.onMessage(encLog, outboundChannel)
		},
		nil,
		backendDisconnected,
		nil,
		12*time.Hour,
		c.logger,
	)

	return backendSub, nil
}

func (c *EncRPCClient) onMessage(encLog []byte, outboundChannel chan types.Log) error {
	jsonLogs, err := c.decryptResponse(encLog)
	if err != nil {
		c.logger.Error("could not decrypt logs received from subscription.", log.ErrKey, err)
		return err
	}

	var logs []*types.Log
	err = json.Unmarshal(jsonLogs, &logs)
	if err != nil {
		c.logger.Error(fmt.Sprintf("could not unmarshal log from JSON. Received data: %s.", string(jsonLogs)), log.ErrKey, err)
		return err
	}

	for _, decryptedLog := range logs {
		outboundChannel <- *decryptedLog
	}
	return nil
}

func (c *EncRPCClient) newHeadSubscription(ctx context.Context, namespace string, ch interface{}, args ...any) (*gethrpc.ClientSubscription, error) {
	return nil, fmt.Errorf("not implemented")
}
