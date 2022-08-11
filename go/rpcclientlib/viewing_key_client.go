package rpcclientlib

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

const (
	http                = "http://"
	reqJSONKeyFrom      = "from"
	reqJSONKeyData      = "data"
	ethCallPaddedArgLen = 64
	ethCallAddrPadding  = "000000000000000000000000"

	// todo: this is a convenience for testnet testing and will eventually be retrieved from the L1
	enclavePublicKeyHex = "034d3b7e63a8bcd532ee3d1d6ecad9d67fca7821981a044551f0f0cbec74d0bc5e"
)

// for these methods, the RPC method's requests and responses should be encrypted
var sensitiveMethods = []string{RPCCall, RPCGetBalance, RPCGetTxReceipt, RPCSendRawTransaction, RPCGetTransactionByHash}

func NewViewingKeyClient(client Client) (*ViewingKeyClient, error) {
	// todo: this is a convenience for testnet but needs to replaced by a parameter and/or retrieved from the target host
	enclPubECDSA, err := crypto.DecompressPubkey(common.Hex2Bytes(enclavePublicKeyHex))
	if err != nil {
		return nil, fmt.Errorf("failed to decompress key for RPC client: %w", err)
	}
	enclavePublicKey := ecies.ImportECDSAPublic(enclPubECDSA)

	return &ViewingKeyClient{
		obscuroClient:    client,
		enclavePublicKey: enclavePublicKey,
	}, nil
}

// NewViewingKeyNetworkClient returns a network RPC client with Viewing Key encryption/decryption
func NewViewingKeyNetworkClient(rpcAddress string) (*ViewingKeyClient, error) {
	rpcClient, err := NewNetworkClient(rpcAddress)
	if err != nil {
		return nil, err
	}
	vkClient, err := NewViewingKeyClient(rpcClient)
	if err != nil {
		return nil, err
	}
	return vkClient, nil
}

// ViewingKeyClient is a Client wrapper that implements Client but also has extra functionality for managing viewing key registration and decryption
type ViewingKeyClient struct {
	obscuroClient Client

	enclavePublicKey *ecies.PublicKey
	// todo: add support for multiple keys on the same client?
	viewingPrivKey *ecies.PrivateKey // private viewing key to use for decrypting sensitive requests
	viewingPubKey  []byte            // public viewing key, submitted to the enclave
	viewingKeyAddr common.Address    // address that generated the public key
}

// Call handles JSON rpc requests - if the method is sensitive it will encrypt the args before sending the request and
// then decrypts the response before returning.
// The result must be a pointer so that package json can unmarshal into it. You can also pass nil, in which case the result is ignored.
func (c *ViewingKeyClient) Call(result interface{}, method string, args ...interface{}) error {
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

	// if caller not interested in response or there was a nil response, we're done
	if result == nil || rawResult == nil {
		return nil
	}

	// method is sensitive, so we decrypt it before unmarshalling the result
	decrypted, err := c.decryptResponse(rawResult)
	if err != nil {
		return fmt.Errorf("failed to decrypt args for %s call - %w", method, err)
	}

	// process the decrypted result to get the desired type and set it on the result pointer
	err = c.setResult(decrypted, result)
	if err != nil {
		return fmt.Errorf("failed to extract result from %s response: %w", method, err)
	}

	return nil
}

func (c *ViewingKeyClient) encryptArgs(args ...interface{}) ([]byte, error) {
	if len(args) == 0 {
		return nil, nil
	}

	paramsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("could not json encode request params: %w", err)
	}

	return c.encryptParamBytes(paramsJSON)
}

func (c *ViewingKeyClient) encryptParamBytes(params []byte) ([]byte, error) {
	encryptedParams, err := ecies.Encrypt(rand.Reader, c.enclavePublicKey, params, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt the following request params with enclave public key: %s. Cause: %w", params, err)
	}
	return encryptedParams, nil
}

func (c *ViewingKeyClient) decryptResponse(resultBlob interface{}) ([]byte, error) {
	if c.viewingPrivKey == nil {
		return nil, fmt.Errorf("cannot decrypt response, viewing key has not been setup")
	}
	resultStr, ok := resultBlob.(string)
	if !ok {
		return nil, fmt.Errorf("expected hex string but result was of type %t instead, with value %s", resultBlob, resultBlob)
	}
	encryptedResult := common.Hex2Bytes(resultStr)
	decryptedResult, err := c.viewingPrivKey.Decrypt(encryptedResult, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt the following response with viewing key: %s. Cause: %w", resultStr, err)
	}

	return decryptedResult, nil
}

// setResult tries to cast/unmarshal data into the result pointer, based on its type
func (c *ViewingKeyClient) setResult(data []byte, result interface{}) error {
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

func (c *ViewingKeyClient) Stop() {
	c.obscuroClient.Stop()
}

func (c *ViewingKeyClient) SetViewingKey(viewingKey *ecies.PrivateKey, signerAddress common.Address, viewingPubKeyBytes []byte) {
	c.viewingPrivKey = viewingKey
	c.viewingPubKey = viewingPubKeyBytes
	c.viewingKeyAddr = signerAddress
}

func (c *ViewingKeyClient) RegisterViewingKey(signature []byte) error {
	// We encrypt the viewing key bytes
	encryptedViewingKeyBytes, err := ecies.Encrypt(rand.Reader, c.enclavePublicKey, c.viewingPubKey, nil, nil)
	if err != nil {
		return fmt.Errorf("could not encrypt viewing key with enclave public key: %w", err)
	}

	var rpcErr error
	err = c.Call(&rpcErr, RPCAddViewingKey, encryptedViewingKeyBytes, signature)
	if err != nil {
		return fmt.Errorf("could not add viewing key: %w", err)
	}
	return nil
}

// The enclave requires the `from` field to be set so that it can encrypt the response, but sources like MetaMask often
// don't set it. So we check whether it's present; if absent, we walk through the arguments in the request's `data`
// field, and if any of the arguments match our viewing key address, we set the `from` field to that address.
func (c *ViewingKeyClient) setFromFieldIfMissing(args []interface{}) ([]interface{}, error) {
	callParams, err := parseCallParams(args)
	if err != nil {
		return nil, fmt.Errorf("could not parse eth_call params. Cause: %w", err)
	}

	// We only modify `eth_call` requests where the `from` field is not set.
	if callParams[reqJSONKeyFrom] != nil {
		return args, nil
	}

	// TODO - Once we support multiple viewing keys, set the `from` field to the single viewing key if there's exactly one.

	// We attempt to set the `from` field based on the `data` field.
	fromAddress, err := searchDataFieldForFrom(callParams, &c.viewingKeyAddr)
	if err != nil {
		return nil, fmt.Errorf("could not process data field in eth_call params. Cause: %w", err)
	}
	if fromAddress != nil {
		callParams[reqJSONKeyFrom] = c.viewingKeyAddr
		args[0] = callParams
		return args, nil
	}

	// TODO - Consider defining an additional fallback to set the `from` field if the above all fail.

	return nil, fmt.Errorf("eth_call request did not have its `from` field set, and its `data` field " +
		"did not contain an address matching a viewing key. Aborting request as it will not be possible to " +
		"encrypt the response")
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

// Parses the eth_call params into a map.
func parseCallParams(args []interface{}) (map[string]interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("expected %s params to have a 'from' field but no params found", RPCCall)
	}

	callParams, ok := args[0].(map[string]interface{})
	if !ok {
		callParamsJSON, ok := args[0].([]byte)
		if !ok {
			return nil, fmt.Errorf("expected eth_call first param to be a map or json encoded bytes but "+
				"was %t", args[0])
		}

		err := json.Unmarshal(callParamsJSON, &callParams)
		if err != nil {
			return nil, fmt.Errorf("expected eth_call first param to be a map or json encoded bytes, "+
				"failed to decode: %w", err)
		}
	}

	return callParams, nil
}

// Extracts the arguments from the request's `data` field. If any of them, after removing padding, match the viewing
// key address, we return that address. Otherwise, we return nil.
func searchDataFieldForFrom(callParams map[string]interface{}, viewingKeyAddress *common.Address) (*common.Address, error) {
	// We ensure that the `data` field is present.
	data := callParams[reqJSONKeyData]
	if data == nil {
		return nil, fmt.Errorf("eth_call request did not have its `data` field set")
	}
	dataString, ok := data.(string)
	if !ok {
		return nil, fmt.Errorf("eth_call request's `data` field was not of the expected type `string`")
	}

	// We check that the data field is long enough before removing the leading "0x" (1 bytes/2 chars) and the method ID
	// (4 bytes/8 chars).
	if len(dataString) < 10 {
		return nil, nil //nolint:nilnil
	}
	dataString = dataString[10:]

	// We split up the arguments in the `data` field.
	var dataArgs []string
	for i := 0; i < len(dataString); i += ethCallPaddedArgLen {
		if i+ethCallPaddedArgLen > len(dataString) {
			break
		}
		dataArgs = append(dataArgs, dataString[i:i+ethCallPaddedArgLen])
	}

	// We iterate over the arguments, looking for an argument that matches the viewing key address. If we find one, we
	// set the `from` field to that address.
	for _, dataArg := range dataArgs {
		// If the argument doesn't have the correct padding, it's not an address.
		if !strings.HasPrefix(dataArg, ethCallAddrPadding) {
			continue
		}

		maybeAddress := common.HexToAddress(dataArg[len(ethCallAddrPadding):])
		if maybeAddress == *viewingKeyAddress {
			return viewingKeyAddress, nil
		}
	}

	return nil, nil //nolint:nilnil
}
