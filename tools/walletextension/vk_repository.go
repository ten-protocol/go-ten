package walletextension

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

func NewViewingKeyRepository() *ViewingKeyRepository {
	return &ViewingKeyRepository{
		viewingKeysPrivate: make(map[common.Address]*ecies.PrivateKey),
		viewingKeysPublic:  make(map[common.Address][]byte),
	}
}

type ViewingKeyRepository struct {
	viewingKeysPrivate map[common.Address]*ecies.PrivateKey // Maps an address to its private viewing key.
	viewingKeysPublic  map[common.Address][]byte            // Maps an address to its public viewing key bytes.
}

func (v *ViewingKeyRepository) IsReady() bool {
	return len(v.viewingKeysPrivate) > 0
}

func (v *ViewingKeyRepository) DecryptBytes(respBytes []byte) ([]byte, error) {
	// We attempt to decrypt the response with each viewing key in turn.
	// TODO - Avoid trying all keys for certain requests by inspecting the request body?
	for _, privateKey := range v.viewingKeysPrivate {
		decryptedResult, err := privateKey.Decrypt(respBytes, nil, nil)
		if err == nil {
			// The decryption did not error, which means we successfully decrypted the result.
			return decryptedResult, nil
		}
	}

	return nil, fmt.Errorf("could not decrypt the response with any of the registered viewing keys")
}

func (v *ViewingKeyRepository) suggestFromAddressForEthCall(callParams map[string]interface{}) (*common.Address, error) {
	var fromAddress *common.Address
	// If there's only one viewing key, we use that to set the `from` field.
	if len(v.viewingKeysPrivate) == 1 {
		for address := range v.viewingKeysPrivate {
			foundAddress := address
			fromAddress = &foundAddress
			break
		}
	} else {
		var err error
		// Otherwise, we search the `data` field for an address matching a registered viewing key.
		fromAddress, err = searchDataFieldForFrom(callParams, v.viewingKeysPrivate)
		if err != nil {
			return nil, fmt.Errorf("could not process data field in eth_call params. Cause: %w", err)
		}
	}

	// TODO - Consider defining an additional fallback to set the `from` field if the above all fail.

	// error if no suitable address found
	if fromAddress == nil {
		return nil, fmt.Errorf("eth_call request did not have its `from` field set, and its `data` field " +
			"did not contain an address matching a viewing key. Aborting request as it will not be possible to " +
			"encrypt the response")
	}

	return fromAddress, nil
}

func (v *ViewingKeyRepository) SetViewingKey(address common.Address, viewingKey *ecies.PrivateKey, viewingKeyPublic []byte) {
	v.viewingKeysPrivate[address] = viewingKey
	v.viewingKeysPublic[address] = viewingKeyPublic
}
