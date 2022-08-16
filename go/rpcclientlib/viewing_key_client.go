package rpcclientlib

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/obscuronet/go-obscuro/go/wallet"
)

const (
	http = "http://"

	// todo: this is a convenience for testnet testing and will eventually be retrieved from the L1
	enclavePublicKeyHex = "034d3b7e63a8bcd532ee3d1d6ecad9d67fca7821981a044551f0f0cbec74d0bc5e"
)

// for these methods, the RPC method's requests and responses should be encrypted
var sensitiveMethods = []string{RPCCall, RPCGetBalance, RPCGetTxReceipt, RPCSendRawTransaction, RPCGetTransactionByHash}

func NewViewingKeyClient(client Client, vkManager ViewingKeyManager) (*ViewingKeyClient, error) {
	// todo: this is a convenience for testnet but needs to replaced by a parameter and/or retrieved from the target host
	enclPubECDSA, err := crypto.DecompressPubkey(common.Hex2Bytes(enclavePublicKeyHex))
	if err != nil {
		return nil, fmt.Errorf("failed to decompress key for RPC client: %w", err)
	}
	enclavePublicKey := ecies.ImportECDSAPublic(enclPubECDSA)

	vkClient := &ViewingKeyClient{
		obscuroClient:     client,
		enclavePublicKey:  enclavePublicKey,
		viewingKeyManager: vkManager,
	}

	return vkClient, nil
}

// NewViewingKeyNetworkClient returns a network RPC client with Viewing Key encryption/decryption for a single account
func NewViewingKeyNetworkClient(rpcAddress string, wallet wallet.Wallet, viewingKey *ecies.PrivateKey, vkPublic []byte, signedVK []byte) (*ViewingKeyClient, error) {
	rpcClient, err := NewNetworkClient(rpcAddress)
	if err != nil {
		return nil, err
	}
	vkManager := SingleAccountVKManager{
		Address:    wallet.Address(),
		ViewingKey: viewingKey,
	}
	vkClient, err := NewViewingKeyClient(rpcClient, vkManager)
	if err != nil {
		return nil, err
	}

	// todo: should we be registering in here or should that be a separate step for calling code?
	err = vkClient.RegisterViewingKeyWithEnclave(vkPublic, signedVK)
	if err != nil {
		return nil, fmt.Errorf("failed to register viewing key with enclave - %w", err)
	}

	return vkClient, nil
}

// ViewingKeyClient is a Client wrapper that implements Client but also has extra functionality for managing viewing key registration and decryption
type ViewingKeyClient struct {
	obscuroClient     Client
	enclavePublicKey  *ecies.PublicKey  // Used to encrypt messages destined to the enclave.
	viewingKeyManager ViewingKeyManager // handles the Address and viewing keys for the account or accounts using this client
}

// Call handles JSON rpc requests - if the method is sensitive it will encrypt the args before sending the request and
// then decrypts the response before returning.
// The result must be a pointer so that package json can unmarshal into it. You can also pass nil, in which case the result is ignored.
func (c *ViewingKeyClient) Call(result interface{}, method string, args ...interface{}) error {
	if !isSensitive(method) {
		// for non-sensitive methods or when viewing keys are disabled we just delegate directly to the geth RPC client
		return c.obscuroClient.Call(result, method, args...)
	}

	// encode the params into a json blob and encrypt them
	encryptedParams, err := c.encryptArgs(args...)
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
	if !c.viewingKeyManager.IsReady() {
		return nil, fmt.Errorf("cannot decrypt response as viewing key not set up")
	}

	resultStr, ok := resultBlob.(string)
	if !ok {
		return nil, fmt.Errorf("expected hex string but result was of type %t instead, with value %s", resultBlob, resultBlob)
	}
	encryptedResult := common.Hex2Bytes(resultStr)

	return c.viewingKeyManager.DecryptBytes(encryptedResult)
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

func (c *ViewingKeyClient) RegisterViewingKeyWithEnclave(publicViewingKey []byte, signature []byte) error {
	// We encrypt the viewing key bytes
	encryptedViewingKeyBytes, err := ecies.Encrypt(rand.Reader, c.enclavePublicKey, publicViewingKey, nil, nil)
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

// isSensitive indicates whether the RPC method's requests and responses should be encrypted.
func isSensitive(method interface{}) bool {
	for _, m := range sensitiveMethods {
		if m == method {
			return true
		}
	}
	return false
}
