//nolint:contextcheck
package p2p

import (
	"context"
	"fmt"

	gethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/obscuronet/go-obscuro/go/host"

	"github.com/obscuronet/go-obscuro/go/host/rpc/clientapi"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/crypto/ecies"

	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/rpc"
)

const (
	// todo: this is a convenience for testnet testing and will eventually be retrieved from the L1
	enclavePublicKeyHex = "034d3b7e63a8bcd532ee3d1d6ecad9d67fca7821981a044551f0f0cbec74d0bc5e"
)

// An in-memory implementation of `rpc.Client` that speaks directly to the node.
type inMemObscuroClient struct {
	obscuroAPI       *clientapi.ObscuroAPI
	ethAPI           *clientapi.EthereumAPI
	filterAPI        *clientapi.FilterAPI
	obscuroScanAPI   *clientapi.ObscuroScanAPI
	testAPI          *clientapi.TestAPI
	enclavePublicKey *ecies.PublicKey
}

func NewInMemObscuroClient(nodeHost host.Host) rpc.Client {
	// todo: this is a convenience for testnet but needs to replaced by a parameter and/or retrieved from the target host
	enclPubECDSA, err := crypto.DecompressPubkey(gethcommon.Hex2Bytes(enclavePublicKeyHex))
	if err != nil {
		panic(err)
	}
	enclPubKey := ecies.ImportECDSAPublic(enclPubECDSA)

	return &inMemObscuroClient{
		obscuroAPI:       clientapi.NewObscuroAPI(nodeHost),
		ethAPI:           clientapi.NewEthereumAPI(nodeHost),
		filterAPI:        clientapi.NewFilterAPI(nodeHost),
		obscuroScanAPI:   clientapi.NewObscuroScanAPI(nodeHost),
		testAPI:          clientapi.NewTestAPI(nodeHost),
		enclavePublicKey: enclPubKey,
	}
}

func NewInMemoryEncRPCClient(host host.Host, viewingKey *rpc.ViewingKey) *rpc.EncRPCClient {
	inMemClient := NewInMemObscuroClient(host)
	encClient, err := rpc.NewEncRPCClient(inMemClient, viewingKey)
	if err != nil {
		panic(err)
	}
	return encClient
}

// Call bypasses RPC, and invokes methods on the node directly.
func (c *inMemObscuroClient) Call(result interface{}, method string, args ...interface{}) error {
	switch method {
	case rpc.RPCGetID:
		*result.(*gethcommon.Address) = c.testAPI.GetID()
		return nil

	case rpc.RPCSendRawTransaction:
		return c.sendRawTransaction(args)

	case rpc.RPCGetCurrentBlockHead:
		*result.(**types.Header) = c.testAPI.GetCurrentBlockHead()
		return nil

	case rpc.RPCGetCurrentRollupHead:
		*result.(**common.Header) = c.obscuroScanAPI.GetCurrentRollupHead()
		return nil

	case rpc.RPCGetRollupHeader:
		return c.getRollupHeader(result, args)

	case rpc.RPCGetRollup:
		return c.getRollup(result, args)

	case rpc.RPCGetTransactionByHash:
		return c.getTransactionByHash(result, args)

	case rpc.RPCCall:
		return c.rpcCall(result, args)

	case rpc.RPCGetTransactionCount:
		return c.getTransactionCount(result, args)

	case rpc.RPCGetTxReceipt:
		return c.getTransactionReceipt(result, args)

	case rpc.RPCStopHost:
		c.testAPI.StopHost()
		return nil

	case rpc.RPCAddViewingKey:
		return c.addViewingKey(args)

	case rpc.RPCGetBalance:
		return c.getBalance(result, args)

	case rpc.RPCGetLogs:
		return c.getLogs(result, args)

	default:
		return fmt.Errorf("RPC method %s is unknown", method)
	}
}

// CallContext not currently supported by in-memory obscuro client, the context will be ignored.
func (c *inMemObscuroClient) CallContext(_ context.Context, result interface{}, method string, args ...interface{}) error {
	return c.Call(result, method, args...)
}

func (c *inMemObscuroClient) Subscribe(context.Context, interface{}, string, interface{}, ...interface{}) (*gethrpc.ClientSubscription, error) {
	panic("not implemented")
}

func (c *inMemObscuroClient) sendRawTransaction(args []interface{}) error {
	encBytes, err := getEncryptedBytes(args, rpc.RPCSendRawTransaction)
	if err != nil {
		return err
	}

	_, err = c.ethAPI.SendRawTransaction(context.Background(), encBytes)
	return err
}

func (c *inMemObscuroClient) getTransactionByHash(result interface{}, args []interface{}) error {
	enc, err := getEncryptedBytes(args, rpc.RPCGetTransactionByHash)
	if err != nil {
		return err
	}
	encryptedResponse, err := c.ethAPI.GetTransactionByHash(context.Background(), enc)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.RPCGetTransactionByHash, err)
	}

	// GetTransactionByHash returns string pointer, we want string
	if encryptedResponse != nil {
		*result.(*interface{}) = *encryptedResponse
	}
	return nil
}

func (c *inMemObscuroClient) rpcCall(result interface{}, args []interface{}) error {
	enc, err := getEncryptedBytes(args, rpc.RPCCall)
	if err != nil {
		return err
	}
	encryptedResponse, err := c.ethAPI.Call(context.Background(), enc)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.RPCCall, err)
	}
	*result.(*interface{}) = encryptedResponse
	return nil
}

func (c *inMemObscuroClient) getTransactionReceipt(result interface{}, args []interface{}) error {
	enc, err := getEncryptedBytes(args, rpc.RPCGetTxReceipt)
	if err != nil {
		return err
	}
	encryptedResponse, err := c.ethAPI.GetTransactionReceipt(context.Background(), enc)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.RPCGetTxReceipt, err)
	}

	// GetTransactionReceipt returns string pointer, we want string
	if encryptedResponse != nil {
		*result.(*interface{}) = *encryptedResponse
	}
	return nil
}

func (c *inMemObscuroClient) getRollupHeader(result interface{}, args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg to %s, got %d", rpc.RPCGetRollupHeader, len(args))
	}
	// we expect a hex string representation of the hash, since that's what gets sent over RPC
	hashStr, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("arg to %s was not of expected type string", rpc.RPCGetRollupHeader)
	}
	hash := gethcommon.HexToHash(hashStr)

	*result.(**common.Header) = c.testAPI.GetRollupHeader(hash)
	return nil
}

func (c *inMemObscuroClient) getRollup(result interface{}, args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg to %s, got %d", rpc.RPCGetRollup, len(args))
	}
	hash, ok := args[0].(gethcommon.Hash)
	if !ok {
		return fmt.Errorf("arg to %s was not of expected type common.Hash", rpc.RPCGetRollup)
	}

	extRollup, err := c.obscuroScanAPI.GetRollup(hash)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.RPCGetRollup, err)
	}
	*result.(**common.ExtRollup) = extRollup
	return nil
}

func (c *inMemObscuroClient) getTransactionCount(result interface{}, args []interface{}) error {
	enc, err := getEncryptedBytes(args, rpc.RPCGetTransactionCount)
	if err != nil {
		return err
	}
	encryptedResponse, err := c.ethAPI.GetTransactionCount(context.Background(), enc)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.RPCGetTransactionCount, err)
	}

	*result.(*interface{}) = encryptedResponse
	return nil
}

func (c *inMemObscuroClient) getBalance(result interface{}, args []interface{}) error {
	enc, err := getEncryptedBytes(args, rpc.RPCGetBalance)
	if err != nil {
		return err
	}
	encryptedResponse, err := c.ethAPI.GetBalance(context.Background(), enc)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.RPCGetBalance, err)
	}
	*result.(*interface{}) = encryptedResponse
	return nil
}

func (c *inMemObscuroClient) getLogs(result interface{}, args []interface{}) error {
	enc, err := getEncryptedBytes(args, rpc.RPCGetLogs)
	if err != nil {
		return err
	}
	encryptedResponse, err := c.filterAPI.GetLogs(context.Background(), enc)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.RPCGetLogs, err)
	}
	*result.(*interface{}) = encryptedResponse
	return nil
}

func (c *inMemObscuroClient) Stop() {
	// There is no RPC connection to close.
}

func (c *inMemObscuroClient) SetViewingKey(_ *ecies.PrivateKey, _ []byte) {
	panic("viewing key encryption/decryption is not currently supported by in-memory obscuro-client")
}

func (c *inMemObscuroClient) RegisterViewingKey(_ gethcommon.Address, _ []byte) error {
	panic("viewing key encryption/decryption is not currently supported by in-memory obscuro-client")
}

func (c *inMemObscuroClient) addViewingKey(args []interface{}) error {
	if len(args) != 2 {
		return fmt.Errorf("expected 2 args to %s, got %d", rpc.RPCAddViewingKey, len(args))
	}

	vk, ok := args[0].([]byte)
	if !ok {
		return fmt.Errorf("expected first arg to %s containing viewing key bytes but it had type %t", rpc.RPCAddViewingKey, args[0])
	}

	sig, ok := args[1].([]byte)
	if !ok {
		return fmt.Errorf("expected second arg to %s containing signature bytes but it had type %t", rpc.RPCAddViewingKey, args[1])
	}
	return c.obscuroAPI.AddViewingKey(vk, sig)
}

// getEncryptedBytes expects args to have a single element and it to be of type bytes (client doesn't know anything about what's getting passed through on sensitive methods)
func getEncryptedBytes(args []interface{}, methodName string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("expected 1 arg to %s, got %d", methodName, len(args))
	}
	encBytes, ok := args[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("expected single arg to %s containing bytes but it had type %t", methodName, args[0])
	}
	return encBytes, nil
}
