package mgmtcontractlib

import (
	"encoding/base64"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/contracts/generated/ManagementContract"
	"github.com/ten-protocol/go-ten/contracts/generated/MessageBus"
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
	CreateRollup(t *ethadapter.L1RollupTx) types.TxData
	CreateBlobRollup(t *ethadapter.L1RollupTx) (types.TxData, error)
	CreateRequestSecret(tx *ethadapter.L1RequestSecretTx) types.TxData
	CreateRespondSecret(tx *ethadapter.L1RespondSecretTx, verifyAttester bool) types.TxData
	CreateInitializeSecret(tx *ethadapter.L1InitializeSecretTx) types.TxData

	// DecodeTx receives a *types.Transaction and converts it to a common.L1Transaction
	DecodeTx(tx *types.Transaction) ethadapter.L1Transaction
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

func (c *contractLibImpl) GetContractAddr() *gethcommon.Address {
	return c.addr
}

func (c *contractLibImpl) DecodeTx(tx *types.Transaction) ethadapter.L1Transaction {
	if tx.To() == nil || tx.To().Hex() != c.addr.Hex() || len(tx.Data()) == 0 {
		return nil
	}
	method, err := c.contractABI.MethodById(tx.Data()[:methodBytesLen])
	if err != nil {
		panic(err)
	}

	contractCallData := map[string]interface{}{}
	switch method.Name {
	case AddRollupMethod:
		if tx.Type() == types.BlobTxType {
			blobHashes := ethadapter.ToIndexedBlobHashes(tx.BlobHashes()...)
			return &ethadapter.L1RollupHashes{
				BlobHashes: blobHashes,
			}
		} else {
			return nil
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
			return nil
		}
		return tx
	}

	return nil
}

func (c *contractLibImpl) CreateRollup(t *ethadapter.L1RollupTx) types.TxData {
	decodedRollup, err := common.DecodeRollup(t.Rollup)
	if err != nil {
		panic(err)
	}

	encRollupData := base64EncodeToString(t.Rollup)

	metaRollup := ManagementContract.StructsMetaRollup{
		Hash:               decodedRollup.Hash(),
		Signature:          decodedRollup.Header.Signature,
		LastSequenceNumber: big.NewInt(int64(decodedRollup.Header.LastBatchSeqNo)),
	}

	crossChain := ManagementContract.StructsHeaderCrossChainData{
		Messages: convertCrossChainMessages(decodedRollup.Header.CrossChainMessages),
	}

	data, err := c.contractABI.Pack(
		AddRollupMethod,
		metaRollup,
		encRollupData,
		crossChain,
	)
	if err != nil {
		panic(err)
	}

	return &types.LegacyTx{
		To:   c.addr,
		Data: data,
	}
}

func (c *contractLibImpl) CreateBlobRollup(t *ethadapter.L1RollupTx) (types.TxData, error) {
	decodedRollup, err := common.DecodeRollup(t.Rollup)
	if err != nil {
		panic(err)
	}

	metaRollup := ManagementContract.StructsMetaRollup{
		Hash:               decodedRollup.Hash(),
		Signature:          decodedRollup.Header.Signature,
		LastSequenceNumber: big.NewInt(int64(decodedRollup.Header.LastBatchSeqNo)),
	}

	crossChain := ManagementContract.StructsHeaderCrossChainData{
		Messages: convertCrossChainMessages(decodedRollup.Header.CrossChainMessages),
	}

	data, err := c.contractABI.Pack(
		AddRollupMethod,
		metaRollup,
		crossChain,
	)
	if err != nil {
		panic(err)
	}

	//TODO handle when blobs exceed 1Mb
	blobs, err := ethadapter.EncodeBlobs(t.Rollup)
	if err != nil {
		return nil, fmt.Errorf("failed to encode rollup to blobs: %w", err)
	}

	var blobHashes []gethcommon.Hash
	var sidecar *types.BlobTxSidecar

	if sidecar, blobHashes, err = makeSidecar(blobs); err != nil {
		return nil, fmt.Errorf("failed to make sidecar: %w", err)
	}

	return &types.BlobTx{
		To:         *c.addr,
		Data:       data,
		BlobHashes: blobHashes,
		Sidecar:    sidecar,
	}, nil
}

func (c *contractLibImpl) CreateRequestSecret(tx *ethadapter.L1RequestSecretTx) types.TxData {
	data, err := c.contractABI.Pack(RequestSecretMethod, base64EncodeToString(tx.Attestation))
	if err != nil {
		panic(err)
	}

	return &types.LegacyTx{
		To:   c.addr,
		Data: data,
	}
}

func (c *contractLibImpl) CreateRespondSecret(tx *ethadapter.L1RespondSecretTx, verifyAttester bool) types.TxData {
	data, err := c.contractABI.Pack(
		RespondSecretMethod,
		tx.AttesterID,
		tx.RequesterID,
		tx.AttesterSig,
		tx.Secret,
		verifyAttester,
	)
	if err != nil {
		panic(err)
	}
	return &types.LegacyTx{
		To:   c.addr,
		Data: data,
	}
}

func (c *contractLibImpl) CreateInitializeSecret(tx *ethadapter.L1InitializeSecretTx) types.TxData {
	data, err := c.contractABI.Pack(
		InitializeSecretMethod,
		tx.EnclaveID,
		tx.InitialSecret,
		base64EncodeToString(tx.Attestation),
	)
	if err != nil {
		panic(err)
	}
	return &types.LegacyTx{
		To:   c.addr,
		Data: data,
	}
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

func (c *contractLibImpl) unpackInitSecretTx(tx *types.Transaction, method *abi.Method, contractCallData map[string]interface{}) *ethadapter.L1InitializeSecretTx {
	err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[methodBytesLen:])
	if err != nil {
		panic(err)
	}
	callData, found := contractCallData["_genesisAttestation"]
	if !found {
		panic("call data not found for requestReport")
	}

	att := Base64DecodeFromString(callData.(string))
	if err != nil {
		c.logger.Crit("could not decode genesis attestation request.", log.ErrKey, err)
	}

	// todo (#1275) - add the other fields
	return &ethadapter.L1InitializeSecretTx{
		Attestation: att,
	}
}

func (c *contractLibImpl) unpackRequestSecretTx(tx *types.Transaction, method *abi.Method, contractCallData map[string]interface{}) *ethadapter.L1RequestSecretTx {
	err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[methodBytesLen:])
	if err != nil {
		panic(err)
	}
	callData, found := contractCallData["requestReport"]
	if !found {
		panic("call data not found for requestReport")
	}

	att := Base64DecodeFromString(callData.(string))
	if err != nil {
		c.logger.Crit("could not decode attestation request.", log.ErrKey, err)
	}
	return &ethadapter.L1RequestSecretTx{
		Attestation: att,
	}
}

func (c *contractLibImpl) unpackRespondSecretTx(tx *types.Transaction, method *abi.Method, contractCallData map[string]interface{}) *ethadapter.L1RespondSecretTx {
	err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[methodBytesLen:])
	if err != nil {
		c.logger.Crit("could not unpack transaction.", log.ErrKey, err)
	}

	requesterData, found := contractCallData["requesterID"]
	if !found {
		c.logger.Crit("call data not found for requesterID")
	}
	requesterAddr, ok := requesterData.(gethcommon.Address)
	if !ok {
		c.logger.Crit("could not decode requester data")
	}

	attesterData, found := contractCallData["attesterID"]
	if !found {
		c.logger.Crit("call data not found for attesterID")
	}
	attesterAddr, ok := attesterData.(gethcommon.Address)
	if !ok {
		c.logger.Crit("could not decode attester data")
	}

	responseSecretData, found := contractCallData["responseSecret"]
	if !found {
		c.logger.Crit("call data not found for responseSecret")
	}
	responseSecretBytes, ok := responseSecretData.([]uint8)
	if !ok {
		c.logger.Crit("could not decode responseSecret data")
	}

	return &ethadapter.L1RespondSecretTx{
		AttesterID:  attesterAddr,
		RequesterID: requesterAddr,
		Secret:      responseSecretBytes[:],
	}
}

func (c *contractLibImpl) unpackSetImportantContractsTx(tx *types.Transaction, method *abi.Method, contractCallData map[string]interface{}) (*ethadapter.L1SetImportantContractsTx, error) {
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

	return &ethadapter.L1SetImportantContractsTx{
		Key:        keyString,
		NewAddress: contractAddress,
	}, nil
}

// base64EncodeToString encodes a byte array to a string
func base64EncodeToString(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

// Base64DecodeFromString decodes a string to a byte array
func Base64DecodeFromString(in string) []byte {
	bytesStr, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}
	return bytesStr
}

func convertCrossChainMessages(messages []MessageBus.StructsCrossChainMessage) []ManagementContract.StructsCrossChainMessage {
	msgs := make([]ManagementContract.StructsCrossChainMessage, 0)

	for _, message := range messages {
		msgs = append(msgs, ManagementContract.StructsCrossChainMessage{
			Sender:   message.Sender,
			Sequence: message.Sequence,
			Nonce:    message.Nonce,
			Topic:    message.Topic,
			Payload:  message.Payload,
		})
	}

	return msgs
}

// MakeSidecar builds & returns the BlobTxSidecar and corresponding blob hashes from the raw blob
// data.
func makeSidecar(blobs []kzg4844.Blob) (*types.BlobTxSidecar, []gethcommon.Hash, error) {
	sidecar := &types.BlobTxSidecar{}
	var blobHashes []gethcommon.Hash
	for i, blob := range blobs {
		sidecar.Blobs = append(sidecar.Blobs, blob)
		commitment, err := kzg4844.BlobToCommitment(&blob)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot compute KZG commitment of blob %d in tx candidate: %w", i, err)
		}
		sidecar.Commitments = append(sidecar.Commitments, commitment)
		proof, err := kzg4844.ComputeBlobProof(&blob, commitment)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot compute KZG proof for fast commitment verification of blob %d in tx candidate: %w", i, err)
		}
		sidecar.Proofs = append(sidecar.Proofs, proof)
		blobHashes = append(blobHashes, ethadapter.KZGToVersionedHash(commitment))
	}
	return sidecar, blobHashes, nil
}
