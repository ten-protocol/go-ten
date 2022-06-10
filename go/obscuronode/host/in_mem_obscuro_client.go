package host

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
)

// An in-memory implementation of `obscuroclient.Client` that speaks directly to the node.
type inMemObscuroClient struct {
	obscuroAPI ObscuroAPI
	ethAPI     EthereumAPI
}

func NewInMemObscuroClient(host *Node) obscuroclient.Client {
	return &inMemObscuroClient{
		obscuroAPI: *NewObscuroAPI(host),
		ethAPI:     *NewEthereumAPI(host),
	}
}

// Call bypasses RPC, and invokes methods on the node directly.
func (c *inMemObscuroClient) Call(result interface{}, method string, args ...interface{}) error {
	switch method {
	case obscuroclient.RPCGetID:
		*result.(*common.Address) = c.obscuroAPI.GetID()

	case obscuroclient.RPCSendTransactionEncrypted:
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg to %s, got %d", obscuroclient.RPCSendTransactionEncrypted, len(args))
		}
		tx, ok := args[0].(nodecommon.EncryptedTx)
		if !ok {
			return fmt.Errorf("arg to %s was not of expected type nodecommon.EncryptedTx", obscuroclient.RPCSendTransactionEncrypted)
		}

		c.obscuroAPI.SendTransactionEncrypted(tx)

	case obscuroclient.RPCGetCurrentBlockHead:
		*result.(**types.Header) = c.obscuroAPI.GetCurrentBlockHead()

	case obscuroclient.RPCGetCurrentRollupHead:
		*result.(**nodecommon.Header) = c.obscuroAPI.GetCurrentRollupHead()

	case obscuroclient.RPCGetRollupHeader:
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg to %s, got %d", obscuroclient.RPCGetRollupHeader, len(args))
		}
		hash, ok := args[0].(common.Hash)
		if !ok {
			return fmt.Errorf("arg to %s was not of expected type common.Hash", obscuroclient.RPCGetRollupHeader)
		}

		*result.(**nodecommon.Header) = c.obscuroAPI.GetRollupHeader(hash)

	case obscuroclient.RPCGetRollup:
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg to %s, got %d", obscuroclient.RPCGetRollup, len(args))
		}
		hash, ok := args[0].(common.Hash)
		if !ok {
			return fmt.Errorf("arg to %s was not of expected type common.Hash", obscuroclient.RPCGetRollup)
		}

		*result.(**nodecommon.ExtRollup) = c.obscuroAPI.GetRollup(hash)

	case obscuroclient.RPCGetTransaction:
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg to %s, got %d", obscuroclient.RPCGetTransaction, len(args))
		}
		hash, ok := args[0].(common.Hash)
		if !ok {
			return fmt.Errorf("arg to %s was not of expected type common.Hash", obscuroclient.RPCGetTransaction)
		}

		*result.(**nodecommon.L2Tx) = c.obscuroAPI.GetTransaction(hash)

	case obscuroclient.RPCCall:
		if len(args) != 3 {
			return fmt.Errorf("expected 3 arg to %s, got %d", obscuroclient.RPCCall, len(args))
		}
		fromAddress, ok := args[0].(common.Address)
		if !ok {
			return fmt.Errorf("arg 0 to %s was not of expected type common.Address", obscuroclient.RPCCall)
		}
		contractAddress, ok := args[1].(common.Address)
		if !ok {
			return fmt.Errorf("arg 1 to %s was not of expected type common.Address", obscuroclient.RPCCall)
		}
		data, ok := args[2].([]byte)
		if !ok {
			return fmt.Errorf("arg 2 to %s was not of expected type []byte", obscuroclient.RPCCall)
		}

		convertedData := (hexutil.Bytes)(data)
		txArgs := TransactionArgs{From: &fromAddress, To: &contractAddress, Data: &convertedData}
		encryptedResponse, err := c.ethAPI.Call(context.Background(), txArgs, rpc.BlockNumberOrHash{}, nil)

		*result.(*OffChainResponse) = OffChainResponse{
			Response: common.Hex2Bytes(encryptedResponse),
			Error:    err,
		}

	case obscuroclient.RPCNonce:
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg to %s, got %d", obscuroclient.RPCNonce, len(args))
		}
		address, ok := args[0].(common.Address)
		if !ok {
			return fmt.Errorf("arg to %s was not of expected type common.Address", obscuroclient.RPCNonce)
		}

		*result.(*uint64) = c.obscuroAPI.Nonce(address)

	case obscuroclient.RPCStopHost:
		c.obscuroAPI.StopHost()

	default:
		return fmt.Errorf("RPC method %s is unknown", method)
	}

	return nil
}

func (c *inMemObscuroClient) Stop() {
	// There is no RPC connection to close.
}
