package host

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto/ecies"

	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/obscuro-playground/go/common"
	"github.com/obscuronet/obscuro-playground/go/rpcclientlib"
)

// An in-memory implementation of `rpcclientlib.Client` that speaks directly to the node.
type inMemObscuroClient struct {
	obscuroAPI ObscuroAPI
	ethAPI     EthereumAPI
}

func NewInMemObscuroClient(host *Node) rpcclientlib.Client {
	return &inMemObscuroClient{
		obscuroAPI: *NewObscuroAPI(host),
		ethAPI:     *NewEthereumAPI(host),
	}
}

// Call bypasses RPC, and invokes methods on the node directly.
func (c *inMemObscuroClient) Call(result interface{}, method string, args ...interface{}) error { //nolint:gocognit
	switch method {
	case rpcclientlib.RPCGetID:
		*result.(*gethcommon.Address) = c.obscuroAPI.GetID()

	case rpcclientlib.RPCSendRawTransaction:
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg to %s, got %d", rpcclientlib.RPCSendRawTransaction, len(args))
		}
		bytes, err := json.Marshal(args)
		if err != nil {
			return fmt.Errorf("failed to marshal the rpc args for %s - %w", rpcclientlib.RPCSendRawTransaction, err)
		}
		_, err = c.ethAPI.SendRawTransaction(context.Background(), bytes)
		return err

	case rpcclientlib.RPCGetCurrentBlockHead:
		*result.(**types.Header) = c.obscuroAPI.GetCurrentBlockHead()

	case rpcclientlib.RPCGetCurrentRollupHead:
		*result.(**common.Header) = c.obscuroAPI.GetCurrentRollupHead()

	case rpcclientlib.RPCGetRollupHeader:
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg to %s, got %d", rpcclientlib.RPCGetRollupHeader, len(args))
		}
		// we expect a hex string representation of the hash, since that's what gets sent over RPC
		hashStr, ok := args[0].(string)
		if !ok {
			return fmt.Errorf("arg to %s was not of expected type string", rpcclientlib.RPCGetRollupHeader)
		}
		hash := gethcommon.HexToHash(hashStr)

		*result.(**common.Header) = c.obscuroAPI.GetRollupHeader(hash)

	case rpcclientlib.RPCGetRollup:
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg to %s, got %d", rpcclientlib.RPCGetRollup, len(args))
		}
		hash, ok := args[0].(gethcommon.Hash)
		if !ok {
			return fmt.Errorf("arg to %s was not of expected type common.Hash", rpcclientlib.RPCGetRollup)
		}

		extRollup, err := c.obscuroAPI.GetRollup(hash)
		if err != nil {
			return fmt.Errorf("`obscuro_getRollup` call failed. Cause: %w", err)
		}
		*result.(**common.ExtRollup) = extRollup

	case rpcclientlib.RPCGetTransactionByHash:
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg to %s, got %d", rpcclientlib.RPCGetTransactionByHash, len(args))
		}
		bytes, err := json.Marshal(args)
		if err != nil {
			return fmt.Errorf("failed to marshal the rpc args for %s - %w", rpcclientlib.RPCGetTransactionByHash, err)
		}

		encryptedTx, err := c.ethAPI.GetTransactionByHash(context.Background(), bytes)
		if err != nil {
			return fmt.Errorf("`eth_getTransactionByHash` call failed. Cause: %w", err)
		}
		*result.(*string) = encryptedTx

	case rpcclientlib.RPCCall:
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg to %s, got %d", rpcclientlib.RPCCall, len(args))
		}
		bytes, err := json.Marshal(args)
		if err != nil {
			return fmt.Errorf("failed to marshal the rpc args for %s - %w", rpcclientlib.RPCCall, err)
		}

		encryptedResponse, err := c.ethAPI.Call(context.Background(), bytes)
		if err != nil {
			return fmt.Errorf("`eth_call` call failed. Cause: %w", err)
		}
		*result.(*string) = encryptedResponse

	case rpcclientlib.RPCNonce:
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg to %s, got %d", rpcclientlib.RPCNonce, len(args))
		}
		address, ok := args[0].(gethcommon.Address)
		if !ok {
			return fmt.Errorf("arg to %s was not of expected type common.Address", rpcclientlib.RPCNonce)
		}

		*result.(*uint64) = c.obscuroAPI.Nonce(address)

	case rpcclientlib.RPCGetTxReceipt:
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg to %s, got %d", rpcclientlib.RPCGetTxReceipt, len(args))
		}
		bytes, err := json.Marshal(args)
		if err != nil {
			return fmt.Errorf("failed to marshal the rpc args for %s - %w", rpcclientlib.RPCGetTxReceipt, err)
		}

		encryptedResponse, err := c.ethAPI.GetTransactionReceipt(context.Background(), bytes)
		if err != nil {
			return fmt.Errorf("`obscuro_getTransactionReceipt` call failed. Cause: %w", err)
		}
		*result.(*string) = encryptedResponse

	case rpcclientlib.RPCStopHost:
		c.obscuroAPI.StopHost()

	default:
		return fmt.Errorf("RPC method %s is unknown", method)
	}

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
