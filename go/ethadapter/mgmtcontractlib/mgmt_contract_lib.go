package mgmtcontractlib

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ten-protocol/go-ten/contracts/generated/ManagementContract"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/ethadapter"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

const methodBytesLen = 4

// MgmtContractLib provides methods for creating ethereum transactions by providing an L1Transaction, creating call
// messages for call requests, and converting ethereum transactions into L1Transactions.
type MgmtContractLib interface {
	IsMock() bool
	BlobHasher() ethadapter.BlobHasher
	PopulateAddRollup(t *common.L1RollupTx, blobs []*kzg4844.Blob, signature common.RollupSignature) (types.TxData, error)
	CreateRequestSecret(tx *common.L1RequestSecretTx) (types.TxData, error)
	CreateRespondSecret(tx *common.L1RespondSecretTx, verifyAttester bool) (types.TxData, error)
	CreateInitializeSecret(tx *common.L1InitializeSecretTx) (types.TxData, error)

	// DecodeTx receives a *types.Transaction and converts it to a common.L1Transaction
	DecodeTx(tx *types.Transaction) (common.L1TenTransaction, error)
	GetContractAddr() *gethcommon.Address

	// The methods below are used to create call messages for mgmt contract data and unpack the responses

	GetHostAddressesMsg() (ethereum.CallMsg, error)
	DecodeHostAddressesResponse(callResponse []byte) ([]string, error)

	SetImportantContractMsg(key string, address gethcommon.Address) (ethereum.CallMsg, error)

	GetImportantContractKeysMsg() (ethereum.CallMsg, error)
	DecodeImportantContractKeysResponse(callResponse []byte) ([]string, error)

	GetImportantAddressCallMsg(key string) (ethereum.CallMsg, error)
	DecodeImportantAddressResponse(callResponse []byte) (gethcommon.Address, error)
}

type contractLibImpl struct {
	addr        *gethcommon.Address
	contractABI abi.ABI
	logger      gethlog.Logger
}

func NewMgmtContractLib(addr *gethcommon.Address, logger gethlog.Logger) MgmtContractLib {
	contractABI, err := abi.JSON(strings.NewReader(MgmtContractABI))
	if err != nil {
		panic(err)
	}

	return &contractLibImpl{
		addr:        addr,
		contractABI: contractABI,
		logger:      logger,
	}
}

func (c *contractLibImpl) IsMock() bool {
	return false
}

func (c *contractLibImpl) BlobHasher() ethadapter.BlobHasher {
	return ethadapter.KZGToVersionedHasher{}
}

func (c *contractLibImpl) GetContractAddr() *gethcommon.Address {
	return c.addr
}

func (c *contractLibImpl) DecodeTx(tx *types.Transaction) (common.L1TenTransaction, error) {
	if tx.To() == nil || tx.To().Hex() != c.addr.Hex() || len(tx.Data()) == 0 {
		return nil, nil
	}
	method, err := c.contractABI.MethodById(tx.Data()[:methodBytesLen])
	if err != nil {
		return nil, fmt.Errorf("could not decode tx. Cause: %w", err)
	}

	contractCallData := map[string]interface{}{}
	switch method.Name {
	case AddRollupMethod:
		if tx.Type() == types.BlobTxType {
			return &common.L1RollupHashes{
				BlobHashes: tx.BlobHashes(),
			}, nil
		} else {
			return nil, nil
		}
	case RespondSecretMethod:
		return c.unpackRespondSecretTx(tx, method, contractCallData)

	case RequestSecretMethod:
		return c.unpackRequestSecretTx(tx, method, contractCallData)

	case InitializeSecretMethod:
		return c.unpackInitSecretTx(tx, method, contractCallData)

	case SetImportantContractsMethod:
		tx, err := c.unpackSetImportantContractsTx(tx, method, contractCallData)
		if err != nil {
			c.logger.Warn("could not unpack set important contracts tx", log.ErrKey, err)
			return nil, err
		}
		return tx, nil
	}

	return nil, nil
}

// PopulateAddRollup creates a BlobTx, encoding the rollup data into blobs.
func (c *contractLibImpl) PopulateAddRollup(t *common.L1RollupTx, blobs []*kzg4844.Blob, signature common.RollupSignature) (types.TxData, error) {
	decodedRollup, err := common.DecodeRollup(t.Rollup)
	if err != nil {
		return nil, fmt.Errorf("could not decode rollup. Cause: %w", err)
	}

	metaRollup := ManagementContract.StructsMetaRollup{
		Hash:               decodedRollup.Hash(),
		Signature:          signature,
		LastSequenceNumber: big.NewInt(int64(decodedRollup.Header.LastBatchSeqNo)),
		BlockBindingHash:   decodedRollup.Header.CompressionL1Head,
		BlockBindingNumber: decodedRollup.Header.CompressionL1Number,
		CrossChainRoot:     decodedRollup.Header.CrossChainRoot,
		LastBatchHash:      decodedRollup.Header.LastBatchHash,
	}

	data, err := c.contractABI.Pack(
		AddRollupMethod,
		metaRollup,
	)
	if err != nil {
		return nil, fmt.Errorf("could not pack the call data. Cause: %w", err)
	}

	var blobHashes []gethcommon.Hash
	var sidecar *types.BlobTxSidecar

	// Use se blobs created here (they are verified that the hash matches with the blobs from the enclave)
	if sidecar, blobHashes, err = ethadapter.MakeSidecar(blobs, c.BlobHasher()); err != nil {
		return nil, fmt.Errorf("failed to make sidecar: %w", err)
	}

	return &types.BlobTx{
		To:         *c.addr,
		Data:       data,
		BlobHashes: blobHashes,
		Sidecar:    sidecar,
	}, nil
}

func (c *contractLibImpl) CreateRequestSecret(tx *common.L1RequestSecretTx) (types.TxData, error) {
	data, err := c.contractABI.Pack(RequestSecretMethod, base64EncodeToString(tx.Attestation))
	if err != nil {
		return nil, fmt.Errorf("could not pack the call data. Cause: %w", err)
	}

	return &types.LegacyTx{
		To:   c.addr,
		Data: data,
	}, nil
}

func (c *contractLibImpl) CreateRespondSecret(tx *common.L1RespondSecretTx, verifyAttester bool) (types.TxData, error) {
	data, err := c.contractABI.Pack(
		RespondSecretMethod,
		tx.AttesterID,
		tx.RequesterID,
		tx.AttesterSig,
		tx.Secret,
		verifyAttester,
	)
	if err != nil {
		return nil, fmt.Errorf("could not pack the call data. Cause: %w", err)
	}
	return &types.LegacyTx{
		To:   c.addr,
		Data: data,
	}, nil
}

func (c *contractLibImpl) CreateInitializeSecret(tx *common.L1InitializeSecretTx) (types.TxData, error) {
	data, err := c.contractABI.Pack(
		InitializeSecretMethod,
		tx.EnclaveID,
		tx.InitialSecret,
		base64EncodeToString(tx.Attestation),
	)
	if err != nil {
		return nil, fmt.Errorf("could not pack the call data. Cause: %w", err)
	}
	return &types.LegacyTx{
		To:   c.addr,
		Data: data,
	}, nil
}

func (c *contractLibImpl) GetHostAddressesMsg() (ethereum.CallMsg, error) {
	data, err := c.contractABI.Pack(GetHostAddressesMethod)
	if err != nil {
		return ethereum.CallMsg{}, fmt.Errorf("could not pack the call data. Cause: %w", err)
	}
	return ethereum.CallMsg{To: c.addr, Data: data}, nil
}

func (c *contractLibImpl) DecodeHostAddressesResponse(callResponse []byte) ([]string, error) {
	unpackedResponse, err := c.contractABI.Unpack(GetHostAddressesMethod, callResponse)
	if err != nil {
		return nil, fmt.Errorf("could not unpack call response. Cause: %w", err)
	}

	// We expect the response to be a list containing one element, that element is a list of address strings
	if len(unpackedResponse) != 1 {
		return nil, fmt.Errorf("unexpected number of results (%d) returned from call, response: %s", len(unpackedResponse), unpackedResponse)
	}
	addresses, ok := unpackedResponse[0].([]string)
	if !ok {
		return nil, fmt.Errorf("could not convert element in call response to list of strings")
	}

	return addresses, nil
}

func (c *contractLibImpl) GetContractNamesMsg() (ethereum.CallMsg, error) {
	data, err := c.contractABI.Pack(GetImportantContractKeysMethod)
	if err != nil {
		return ethereum.CallMsg{}, fmt.Errorf("could not pack the call data. Cause: %w", err)
	}
	return ethereum.CallMsg{To: c.addr, Data: data}, nil
}

func (c *contractLibImpl) DecodeContractNamesResponse(callResponse []byte) ([]string, error) {
	unpackedResponse, err := c.contractABI.Unpack(GetImportantContractKeysMethod, callResponse)
	if err != nil {
		return nil, fmt.Errorf("could not unpack call response. Cause: %w", err)
	}

	// We expect the response to be a list containing one element, that element is a list of address strings
	if len(unpackedResponse) != 1 {
		return nil, fmt.Errorf("unexpected number of results (%d) returned from call, response: %s", len(unpackedResponse), unpackedResponse)
	}
	contractNames, ok := unpackedResponse[0].([]string)
	if !ok {
		return nil, fmt.Errorf("could not convert element in call response to list of strings")
	}

	return contractNames, nil
}

func (c *contractLibImpl) SetImportantContractMsg(key string, address gethcommon.Address) (ethereum.CallMsg, error) {
	data, err := c.contractABI.Pack(SetImportantContractsMethod, key, address)
	if err != nil {
		return ethereum.CallMsg{}, fmt.Errorf("could not pack the call data. Cause: %w", err)
	}
	return ethereum.CallMsg{To: c.addr, Data: data}, nil
}

func (c *contractLibImpl) GetImportantContractKeysMsg() (ethereum.CallMsg, error) {
	data, err := c.contractABI.Pack(GetImportantContractKeysMethod)
	if err != nil {
		return ethereum.CallMsg{}, fmt.Errorf("could not pack the call data. Cause: %w", err)
	}
	return ethereum.CallMsg{To: c.addr, Data: data}, nil
}

func (c *contractLibImpl) DecodeImportantContractKeysResponse(callResponse []byte) ([]string, error) {
	unpackedResponse, err := c.contractABI.Unpack(GetImportantContractKeysMethod, callResponse)
	if err != nil {
		return nil, fmt.Errorf("could not unpack call response. Cause: %w", err)
	}

	// We expect the response to be a list containing one element, that element is a list of address strings
	if len(unpackedResponse) != 1 {
		return nil, fmt.Errorf("unexpected number of results (%d) returned from call, response: %s", len(unpackedResponse), unpackedResponse)
	}
	contractNames, ok := unpackedResponse[0].([]string)
	if !ok {
		return nil, fmt.Errorf("could not convert element in call response to list of strings")
	}

	return contractNames, nil
}

func (c *contractLibImpl) GetImportantAddressCallMsg(key string) (ethereum.CallMsg, error) {
	data, err := c.contractABI.Pack(GetImportantAddressMethod, key)
	if err != nil {
		return ethereum.CallMsg{}, fmt.Errorf("could not pack the call data. Cause: %w", err)
	}
	return ethereum.CallMsg{To: c.addr, Data: data}, nil
}

func (c *contractLibImpl) DecodeImportantAddressResponse(callResponse []byte) (gethcommon.Address, error) {
	unpackedResponse, err := c.contractABI.Unpack(GetImportantAddressMethod, callResponse)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("could not unpack call response. Cause: %w", err)
	}

	// We expect the response to be a list containing one element, that element is a list of address strings
	if len(unpackedResponse) != 1 {
		return gethcommon.Address{}, fmt.Errorf("unexpected number of results (%d) returned from call, response: %s", len(unpackedResponse), unpackedResponse)
	}
	address, ok := unpackedResponse[0].(gethcommon.Address)
	if !ok {
		return gethcommon.Address{}, fmt.Errorf("could not convert element in call response to list of strings")
	}

	return address, nil
}

func (c *contractLibImpl) unpackInitSecretTx(tx *types.Transaction, method *abi.Method, contractCallData map[string]interface{}) (*common.L1InitializeSecretTx, error) {
	err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[methodBytesLen:])
	if err != nil {
		return nil, fmt.Errorf("could not unpack transaction. Cause: %w", err)
	}
	callData, found := contractCallData["_genesisAttestation"]
	if !found {
		return nil, fmt.Errorf("call data not found for _genesisAttestation")
	}

	att, err := Base64DecodeFromString(callData.(string))
	if err != nil {
		return nil, err
	}

	// todo (#1275) - add the other fields
	return &common.L1InitializeSecretTx{
		Attestation: att,
	}, nil
}

func (c *contractLibImpl) unpackRequestSecretTx(tx *types.Transaction, method *abi.Method, contractCallData map[string]interface{}) (*common.L1RequestSecretTx, error) {
	err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[methodBytesLen:])
	if err != nil {
		return nil, fmt.Errorf("could not unpack transaction. Cause: %w", err)
	}
	callData, found := contractCallData["requestReport"]
	if !found {
		return nil, fmt.Errorf("call data not found for requestReport")
	}

	att, err := Base64DecodeFromString(callData.(string))
	if err != nil {
		return nil, fmt.Errorf("could not decode attestation request. Cause: %w", err)
	}
	return &common.L1RequestSecretTx{
		Attestation: att,
	}, nil
}

func (c *contractLibImpl) unpackRespondSecretTx(tx *types.Transaction, method *abi.Method, contractCallData map[string]interface{}) (*common.L1RespondSecretTx, error) {
	err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[methodBytesLen:])
	if err != nil {
		return nil, fmt.Errorf("could not unpack transaction. Cause: %w", err)
	}

	requesterData, found := contractCallData["requesterID"]
	if !found {
		return nil, fmt.Errorf("call data not found for requesterID")
	}
	requesterAddr, ok := requesterData.(gethcommon.Address)
	if !ok {
		return nil, fmt.Errorf("could not decode requester data")
	}

	attesterData, found := contractCallData["attesterID"]
	if !found {
		return nil, fmt.Errorf("call data not found for attesterID")
	}
	attesterAddr, ok := attesterData.(gethcommon.Address)
	if !ok {
		return nil, fmt.Errorf("could not decode attester data")
	}

	responseSecretData, found := contractCallData["responseSecret"]
	if !found {
		return nil, fmt.Errorf("call data not found for responseSecret")
	}
	responseSecretBytes, ok := responseSecretData.([]uint8)
	if !ok {
		return nil, fmt.Errorf("could not decode responseSecret data")
	}

	return &common.L1RespondSecretTx{
		AttesterID:  attesterAddr,
		RequesterID: requesterAddr,
		Secret:      responseSecretBytes[:],
	}, nil
}

func (c *contractLibImpl) unpackSetImportantContractsTx(tx *types.Transaction, method *abi.Method, contractCallData map[string]interface{}) (*common.L1SetImportantContractsTx, error) {
	err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[methodBytesLen:])
	if err != nil {
		return nil, fmt.Errorf("could not unpack transaction. Cause: %w", err)
	}

	keyData, found := contractCallData["key"]
	if !found {
		return nil, fmt.Errorf("call data not found for key")
	}
	keyString, ok := keyData.(string)
	if !ok {
		return nil, fmt.Errorf("could not decode key data")
	}

	contractAddressData, found := contractCallData["newAddress"]
	if !found {
		return nil, fmt.Errorf("call data not found for newAddress")
	}
	contractAddress, ok := contractAddressData.(gethcommon.Address)
	if !ok {
		return nil, fmt.Errorf("could not decode newAddress data")
	}

	return &common.L1SetImportantContractsTx{
		Key:        keyString,
		NewAddress: contractAddress,
	}, nil
}

// base64EncodeToString encodes a byte array to a string
func base64EncodeToString(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

// Base64DecodeFromString decodes a string to a byte array
func Base64DecodeFromString(in string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(in)
}
