package mgmtcontractlib

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"math/big"
	"strings"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/contracts/managementcontract/generated/ManagementContract"
	"github.com/obscuronet/go-obscuro/contracts/messagebuscontract/generated/MessageBus"
	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/ethadapter"

	"github.com/ethereum/go-ethereum"

	"github.com/obscuronet/go-obscuro/go/common"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const methodBytesLen = 4

var (
	// TODO review gas estimation - these should not be static values
	// Gas should be calculated so to not overpay what the operation requires
	// The values are hardcoded at the moment to guarantee the txs will be minted
	// It's using large gas values because rollups can be very expensive
	defaultGasPrice = big.NewInt(20000000000)
	defaultGas      = uint64(1024_000_000)
)

// MgmtContractLib provides methods for creating ethereum transactions by providing an L1Transaction, creating call
// messages for call requests, and converting ethereum transactions into L1Transactions.
type MgmtContractLib interface {
	CreateRollup(t *ethadapter.L1RollupTx, nonce uint64) types.TxData
	CreateRequestSecret(tx *ethadapter.L1RequestSecretTx, nonce uint64) types.TxData
	CreateRespondSecret(tx *ethadapter.L1RespondSecretTx, nonce uint64, verifyAttester bool) types.TxData
	CreateInitializeSecret(tx *ethadapter.L1InitializeSecretTx, nonce uint64) types.TxData
	GetHostAddresses() (ethereum.CallMsg, error)

	// DecodeTx receives a *types.Transaction and converts it to an common.L1Transaction
	DecodeTx(tx *types.Transaction) ethadapter.L1Transaction
	// DecodeCallResponse unpacks a call response into a slice of strings.
	DecodeCallResponse(callResponse []byte) ([][]string, error)
	GetContractAddr() *gethcommon.Address
}

type contractLibImpl struct {
	Addr        *gethcommon.Address
	contractABI abi.ABI
	logger      gethlog.Logger
}

func NewMgmtContractLib(addr *gethcommon.Address, logger gethlog.Logger) MgmtContractLib {
	contractABI, err := abi.JSON(strings.NewReader(MgmtContractABI))
	if err != nil {
		panic(err)
	}

	return &contractLibImpl{
		Addr:        addr,
		contractABI: contractABI,
		logger:      logger,
	}
}

func (c *contractLibImpl) GetContractAddr() *gethcommon.Address {
	return c.Addr
}

func (c *contractLibImpl) DecodeTx(tx *types.Transaction) ethadapter.L1Transaction {
	if tx.To() == nil || tx.To().Hex() != c.Addr.Hex() || len(tx.Data()) == 0 {
		return nil
	}
	method, err := c.contractABI.MethodById(tx.Data()[:methodBytesLen])
	if err != nil {
		panic(err)
	}

	contractCallData := map[string]interface{}{}
	switch method.Name {
	case AddRollupMethod:
		if err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[4:]); err != nil {
			panic(err)
		}
		callData, found := contractCallData["_rollupData"]
		if !found {
			panic("call data not found for rollupData")
		}
		zipped := Base64DecodeFromString(callData.(string))
		rollup, err := Decompress(zipped)
		if err != nil {
			panic(err)
		}

		return &ethadapter.L1RollupTx{
			Rollup: rollup,
		}

	case RespondSecretMethod:
		return c.unpackRespondSecretTx(tx, method, contractCallData)

	case RequestSecretMethod:
		return c.unpackRequestSecretTx(tx, method, contractCallData)

	case InitializeSecretMethod:
		return c.unpackInitSecretTx(tx, method, contractCallData)
	}

	return nil
}

func (c *contractLibImpl) CreateRollup(t *ethadapter.L1RollupTx, nonce uint64) types.TxData {
	decodedRollup, err := common.DecodeRollup(t.Rollup)
	if err != nil {
		panic(err)
	}

	zipped, err := compress(t.Rollup)
	if err != nil {
		panic(err)
	}
	encRollupData := base64EncodeToString(zipped)

	metaRollup := ManagementContract.StructsMetaRollup{
		ParentHash:   decodedRollup.Header.ParentHash,
		Hash:         decodedRollup.Hash(),
		AggregatorID: decodedRollup.Header.Agg,
		L1Block:      decodedRollup.Header.L1Proof,
		Number:       decodedRollup.Header.Number,
	}

	crossChain := ManagementContract.StructsHeaderCrossChainData{
		BlockNumber: decodedRollup.Header.LatestInboundCrossChainHeight,
		BlockHash:   decodedRollup.Header.LatestInboudCrossChainHash,
		Messages:    convertCrossChainMessages(decodedRollup.Header.CrossChainMessages),
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
		Nonce:    nonce,
		GasPrice: defaultGasPrice,
		Gas:      defaultGas,
		To:       c.Addr,
		Data:     data,
	}
}

func (c *contractLibImpl) CreateRequestSecret(tx *ethadapter.L1RequestSecretTx, nonce uint64) types.TxData {
	data, err := c.contractABI.Pack(RequestSecretMethod, base64EncodeToString(tx.Attestation))
	if err != nil {
		panic(err)
	}

	return &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: defaultGasPrice,
		Gas:      defaultGas,
		To:       c.Addr,
		Data:     data,
	}
}

func (c *contractLibImpl) CreateRespondSecret(tx *ethadapter.L1RespondSecretTx, nonce uint64, verifyAttester bool) types.TxData {
	data, err := c.contractABI.Pack(
		RespondSecretMethod,
		tx.AttesterID,
		tx.RequesterID,
		tx.AttesterSig,
		tx.Secret,
		tx.HostAddress,
		verifyAttester,
	)
	if err != nil {
		panic(err)
	}
	return &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: defaultGasPrice,
		Gas:      defaultGas,
		To:       c.Addr,
		Data:     data,
	}
}

func (c *contractLibImpl) CreateInitializeSecret(tx *ethadapter.L1InitializeSecretTx, nonce uint64) types.TxData {
	data, err := c.contractABI.Pack(
		InitializeSecretMethod,
		tx.AggregatorID,
		tx.InitialSecret,
		tx.HostAddress,
		base64EncodeToString(tx.Attestation),
	)
	if err != nil {
		panic(err)
	}
	return &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: defaultGasPrice,
		Gas:      defaultGas,
		To:       c.Addr,
		Data:     data,
	}
}

func (c *contractLibImpl) GetHostAddresses() (ethereum.CallMsg, error) {
	data, err := c.contractABI.Pack(GetHostAddressesMethod)
	if err != nil {
		return ethereum.CallMsg{}, fmt.Errorf("could not pack the call data. Cause: %w", err)
	}
	return ethereum.CallMsg{To: c.Addr, Data: data}, nil
}

func (c *contractLibImpl) DecodeCallResponse(callResponse []byte) ([][]string, error) {
	unpackedResponse, err := c.contractABI.Unpack(GetHostAddressesMethod, callResponse)
	if err != nil {
		return nil, fmt.Errorf("could not unpack call response. Cause: %w", err)
	}

	// We convert the returned interfaces to strings.
	unpackedResponseStrings := make([][]string, 0, len(unpackedResponse))
	for _, obj := range unpackedResponse {
		str, ok := obj.([]string)
		if !ok {
			return nil, fmt.Errorf("could not convert interface in call response to string")
		}
		unpackedResponseStrings = append(unpackedResponseStrings, str)
	}

	return unpackedResponseStrings, nil
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

	// todo - add the other fields
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

	hostAddressData, found := contractCallData["hostAddress"]
	if !found {
		c.logger.Crit("call data not found for hostAddress")
	}
	hostAddressString, ok := hostAddressData.(string)
	if !ok {
		c.logger.Crit("could not decode hostAddress data")
	}

	return &ethadapter.L1RespondSecretTx{
		AttesterID:  attesterAddr,
		RequesterID: requesterAddr,
		Secret:      responseSecretBytes[:],
		HostAddress: hostAddressString,
	}
}

// base64EncodeToString encodes a byte array to a string
func base64EncodeToString(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

// compress the byte array using gzip
func compress(in []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(in); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Base64DecodeFromString decodes a string to a byte array
func Base64DecodeFromString(in string) []byte {
	bytesStr, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}
	return bytesStr
}

// Decompress the byte array using gzip
func Decompress(in []byte) ([]byte, error) {
	reader := bytes.NewReader(in)
	gz, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	return io.ReadAll(gz)
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
