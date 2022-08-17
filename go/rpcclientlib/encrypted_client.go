package rpcclientlib

import (
	"crypto/rand"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

const (
	http           = "http://"
	reqJSONKeyFrom = "from"

	// todo: this is a convenience for testnet testing and will eventually be retrieved from the L1
	enclavePublicKeyHex = "034d3b7e63a8bcd532ee3d1d6ecad9d67fca7821981a044551f0f0cbec74d0bc5e"
)

// for these methods, the RPC method's requests and responses should be encrypted
var sensitiveMethods = []string{RPCCall, RPCGetBalance, RPCGetTxReceipt, RPCSendRawTransaction, RPCGetTransactionByHash}

type ViewingKey struct {
	Account    *common.Address   // Account address that this private key is bound to
	PrivateKey *ecies.PrivateKey // private viewing key
	PublicKey  []byte            // public viewing key in bytes to share with enclave
	SignedKey  []byte            // public viewing key signed by the Account's private key
}

// NewEncRPCClient sets up a client with a viewing key for encrypted communication (this submits the VK to the enclave)
func NewEncRPCClient(client Client, viewingKey *ViewingKey) (*EncRPCClient, error) {
	// todo: this is a convenience for testnet but needs to replaced by a parameter and/or retrieved from the target host
	enclPubECDSA, err := crypto.DecompressPubkey(common.Hex2Bytes(enclavePublicKeyHex))
	if err != nil {
		return nil, fmt.Errorf("failed to decompress key for RPC client: %w", err)
	}
	enclavePublicKey := ecies.ImportECDSAPublic(enclPubECDSA)

	vkClient := &EncRPCClient{
		obscuroClient:    client,
		enclavePublicKey: enclavePublicKey,
		viewingKey:       viewingKey,
	}
	err = vkClient.RegisterViewingKey()
	if err != nil {
		return nil, err
	}

	return vkClient, nil
}

// NewEncNetworkClient returns a network RPC client with Viewing Key encryption/decryption
func NewEncNetworkClient(rpcAddress string, viewingKey *ViewingKey) (*EncRPCClient, error) {
	rpcClient, err := NewNetworkClient(rpcAddress)
	if err != nil {
		return nil, err
	}
	vkClient, err := NewEncRPCClient(rpcClient, viewingKey)
	if err != nil {
		return nil, err
	}
	return vkClient, nil
}

// EncRPCClient is a Client wrapper that implements Client but also has extra functionality for managing viewing key registration and decryption
type EncRPCClient struct {
	obscuroClient    Client
	enclavePublicKey *ecies.PublicKey // Used to encrypt messages destined to the enclave.
	viewingKey       *ViewingKey
}

// Call handles JSON rpc requests - if the method is sensitive it will encrypt the args before sending the request and
// then decrypts the response before returning.
// The result must be a pointer so that package json can unmarshal into it. You can also pass nil, in which case the result is ignored.
func (c *EncRPCClient) Call(result interface{}, method string, args ...interface{}) error {
	if !isSensitive(method) {
		// for non-sensitive methods or when viewing keys are disabled we just delegate directly to the geth RPC client
		return c.obscuroClient.Call(result, method, args...)
	}

	var err error
	if method == RPCCall {
		// RPCCall is a sensitive method that requires a viewing key lookup but the 'from' field is not mandatory in geth
		//	and is often not included from metamask etc. So we ensure it is populated here.
		args, err = c.setFromFieldIfMissing(args)
		if err != nil {
			return err
		}
	}

	// encode the params into a json blob and encrypt them
	var encryptedParams []byte
	encryptedParams, err = c.encryptArgs(args...)
	if err != nil {
		return fmt.Errorf("failed to encrypt args for %s call - %w", method, err)
	}

	// we set up a generic rawResult to receive the response (then we can decrypt it as necessary into the requested result type)
	var rawResult interface{}
	err = c.obscuroClient.Call(&rawResult, method, encryptedParams)
	if err != nil {
		return fmt.Errorf("%s rpc call failed - %w", method, err)
	}

	// if caller not interested in response, we're done
	if result == nil {
		return nil
	}

	if rawResult == nil {
		// note: some methods return nil for 'not found', caller can check for this Error type to verify
		return ErrNilResponse
	}

	// method is sensitive, so we decrypt it before unmarshalling the result
	decrypted, err := c.decryptResponse(rawResult)
	if err != nil {
		return fmt.Errorf("failed to decrypt response for %s call - %w", method, err)
	}

	// process the decrypted result to get the desired type and set it on the result pointer
	err = c.setResult(decrypted, result)
	if err != nil {
		return fmt.Errorf("failed to extract result from %s response: %w", method, err)
	}

	return nil
}

func (c *EncRPCClient) encryptArgs(args ...interface{}) ([]byte, error) {
	if len(args) == 0 {
		return nil, nil
	}

	paramsJSON, err := json.Marshal(args)
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

func (c *EncRPCClient) decryptResponse(resultBlob interface{}) ([]byte, error) {
	resultStr, ok := resultBlob.(string)
	if !ok {
		return nil, fmt.Errorf("expected hex string but result was of type %t instead, with value %s", resultBlob, resultBlob)
	}
	encryptedResult := common.Hex2Bytes(resultStr)

	decryptedResult, err := c.viewingKey.PrivateKey.Decrypt(encryptedResult, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt result with viewing key - %w", err)
	}

	return decryptedResult, nil
}

// setResult tries to cast/unmarshal data into the result pointer, based on its type
func (c *EncRPCClient) setResult(data []byte, result interface{}) error {
	switch result := result.(type) {
	case *string:
		*result = string(data)
		return nil

	case *interface{}:
		err := json.Unmarshal(data, result)
		if err != nil {
			// if unmarshal failed with generic return we can try to send it back as a string
			*result = string(data)
		}
		return nil

	default:
		// for any other type we attempt to json unmarshal it
		return json.Unmarshal(data, result)
	}
}

func (c *EncRPCClient) Stop() {
	c.obscuroClient.Stop()
}

// RegisterViewingKey submits the viewing key with signature to the enclave, this must be called before the viewing key is usable
func (c *EncRPCClient) RegisterViewingKey() error {
	// TODO: Store signatures to be able to resubmit keys if they are evicted by the node?
	// We encrypt the viewing key bytes
	encryptedViewingKeyBytes, err := ecies.Encrypt(rand.Reader, c.enclavePublicKey, c.viewingKey.PublicKey, nil, nil)
	if err != nil {
		return fmt.Errorf("could not encrypt viewing key with enclave public key: %w", err)
	}

	var rpcErr error
	err = c.Call(&rpcErr, RPCAddViewingKey, encryptedViewingKeyBytes, c.viewingKey.SignedKey)
	if err != nil {
		return fmt.Errorf("could not add viewing key: %w", err)
	}
	return nil
}

// The enclave requires the `from` field to be set so that it can encrypt the response, but sources like MetaMask often
// don't set it. So we check whether it's present; if absent, we walk through the arguments in the request's `data`
// field, and if any of the arguments match our viewing key address, we set the `from` field to that address.
func (c *EncRPCClient) setFromFieldIfMissing(args []interface{}) ([]interface{}, error) {
	callParams, err := parseCallParams(args)
	if err != nil {
		return nil, fmt.Errorf("could not parse eth_call params. Cause: %w", err)
	}

	// We only modify `eth_call` requests where the `from` field is not set.
	if callParams[reqJSONKeyFrom] != nil {
		return args, nil
	}

	callParams[reqJSONKeyFrom] = c.viewingKey.Account
	args[0] = callParams
	return args, nil
}

// isSensitive indicates whether the RPC method's requests and responses should be encrypted.
func isSensitive(method interface{}) bool {
	for _, m := range sensitiveMethods {
		if m == method {
			return true
		}
	}
	return false
}

// Many eth RPC requests provide params in a json map with similar fields (e.g. a `from` field)
func parseCallParams(args []interface{}) (map[string]interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no params found to unmarshal")
	}

	params, ok := args[0].(map[string]interface{})
	if !ok {
		callParamsJSON, ok := args[0].([]byte)
		if !ok {
			return nil, fmt.Errorf("expected eth_call first param to be a map or json encoded bytes but "+
				"was %t", args[0])
		}

		err := json.Unmarshal(callParamsJSON, &params)
		if err != nil {
			return nil, fmt.Errorf("expected eth_call first param to be a map or json encoded bytes, "+
				"failed to decode: %w", err)
		}
	}

	return params, nil
}
