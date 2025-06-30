package crypto

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// CreateNetworkSecretResponseTypedMessage builds the EIP-712 typed data for NetworkSecretResponse
func CreateNetworkSecretResponseTypedMessage(
	requesterID common.Address,
	responseSecret []byte,
	chainID int64,
	contractAddress common.Address,
) (apitypes.TypedData, error) {
	// Validate secret length
	if len(responseSecret) != 145 {
		return apitypes.TypedData{}, fmt.Errorf("invalid secret response length: expected 145, got %d", len(responseSecret))
	}

	// EIP-712 domain definition must match Solidity
	domain := apitypes.TypedDataDomain{
		Name:              "NetworkEnclaveRegistry",
		Version:           "1",
		ChainId:           (*math.HexOrDecimal256)(big.NewInt(chainID)),
		VerifyingContract: contractAddress.Hex(),
	}

	// Message fields: dynamic bytes must be hex-encoded
	message := map[string]interface{}{
		"requesterID":    requesterID.Hex(),
		"responseSecret": hexutil.Encode(crypto.Keccak256(responseSecret)),
	}

	// Types definition aligning with Solidity struct
	types := apitypes.Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
		"NetworkSecretResponse": {
			{Name: "requesterID", Type: "address"},
			{Name: "responseSecret", Type: "bytes32"},
		},
	}

	typedData := apitypes.TypedData{
		Types:       types,
		PrimaryType: "NetworkSecretResponse",
		Domain:      domain,
		Message:     message,
	}
	return typedData, nil
}

// CreateNetworkSecretResponseHash computes the EIP-712 digest for NetworkSecretResponse
func CreateNetworkSecretResponseHash(
	requesterID common.Address,
	responseSecret []byte,
	chainID int64,
	contractAddress common.Address,
) (common.Hash, error) {
	// Build typed data
	typedData, err := CreateNetworkSecretResponseTypedMessage(requesterID, responseSecret, chainID, contractAddress)
	if err != nil {
		return common.Hash{}, err
	}
	// Compute full EIP-712 hash: \x19\x01 || domainSeparator || structHash
	hash, err := HashTypedData(typedData)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to hash typed data: %w", err)
	}
	return hash, nil
}

// CreateRollupTypedMessage builds the EIP-712 typed data for Rollup
func CreateRollupTypedMessage(
	firstSequenceNumber *big.Int,
	lastSequenceNumber *big.Int,
	lastBatchHash common.Hash,
	blockBindingHash common.Hash,
	blockBindingNumber *big.Int,
	crossChainRoot common.Hash,
	blobHash common.Hash,
	chainID int64,
	contractAddress common.Address,
) apitypes.TypedData {
	domain := apitypes.TypedDataDomain{
		Name:              "DataAvailabilityRegistry",
		Version:           "1",
		ChainId:           (*math.HexOrDecimal256)(big.NewInt(chainID)),
		VerifyingContract: contractAddress.Hex(),
	}

	message := map[string]interface{}{
		"firstSequenceNumber": firstSequenceNumber.String(),
		"lastSequenceNumber":  lastSequenceNumber.String(),
		"lastBatchHash":       lastBatchHash.Hex(),
		"blockBindingHash":    blockBindingHash.Hex(),
		"blockBindingNumber":  blockBindingNumber.String(),
		"crossChainRoot":      crossChainRoot.Hex(),
		"blobHash":            blobHash.Hex(),
	}

	types := apitypes.Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
		"Rollup": {
			{Name: "firstSequenceNumber", Type: "uint256"},
			{Name: "lastSequenceNumber", Type: "uint256"},
			{Name: "lastBatchHash", Type: "bytes32"},
			{Name: "blockBindingHash", Type: "bytes32"},
			{Name: "blockBindingNumber", Type: "uint256"},
			{Name: "crossChainRoot", Type: "bytes32"},
			{Name: "blobHash", Type: "bytes32"},
		},
	}

	return apitypes.TypedData{
		Types:       types,
		PrimaryType: "Rollup",
		Domain:      domain,
		Message:     message,
	}
}

// HashTypedData computes the EIP-712 final digest (\x19\x01 || domainSep || structHash)
func HashTypedData(typedData apitypes.TypedData) (common.Hash, error) {
	// Domain separator
	domainSep, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to hash domain separator: %w", err)
	}
	// Struct hash
	structHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to hash typed data struct: %w", err)
	}
	// EIP-191 prefix and final keccak
	raw := append([]byte("\x19\x01"), append(domainSep, structHash...)...)
	return crypto.Keccak256Hash(raw), nil
}

// SignTypedData signs the EIP-712 typed data and returns a [R||S||V] signature
func SignTypedData(typedData apitypes.TypedData, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	hash, err := HashTypedData(typedData)
	if err != nil {
		return nil, err
	}
	sig, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign typed data: %w", err)
	}
	return sig, nil
}

// VerifyTypedDataSignature recovers the address that signed the typed data
func VerifyTypedDataSignature(typedData apitypes.TypedData, signature []byte) (common.Address, error) {
	hash, err := HashTypedData(typedData)
	if err != nil {
		return common.Address{}, err
	}
	pubKeyBytes, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to recover public key: %w", err)
	}
	pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to unmarshal public key: %w", err)
	}
	return crypto.PubkeyToAddress(*pubKey), nil
}

// CreateRollupHash computes the EIP-712 hash for Rollup
func CreateRollupHash(
	firstSequenceNumber *big.Int,
	lastSequenceNumber *big.Int,
	lastBatchHash common.Hash,
	blockBindingHash common.Hash,
	blockBindingNumber *big.Int,
	crossChainRoot common.Hash,
	blobHash common.Hash,
	chainID int64,
	contractAddress common.Address,
) (common.Hash, error) {
	// Reuse typed data
	typedData := CreateRollupTypedMessage(
		firstSequenceNumber,
		lastSequenceNumber,
		lastBatchHash,
		blockBindingHash,
		blockBindingNumber,
		crossChainRoot,
		blobHash,
		chainID,
		contractAddress,
	)

	// Compute the domain separator and struct hash, then final EIP-712 hash
	h, err := HashTypedData(typedData)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to hash typed data: %w", err)
	}
	return h, nil
}
