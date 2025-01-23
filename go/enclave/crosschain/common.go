package crosschain

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"golang.org/x/crypto/sha3"

	smt "github.com/FantasyJony/openzeppelin-merkle-tree-go/standard_merkle_tree"
	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/contracts/generated/ManagementContract"
	"github.com/ten-protocol/go-ten/contracts/generated/MessageBus"
	"github.com/ten-protocol/go-ten/go/common"
)

var (
	MessageBusABI, _                  = abi.JSON(strings.NewReader(MessageBus.MessageBusMetaData.ABI))
	MgmtContractABI, _                = abi.JSON(strings.NewReader(ManagementContract.ManagementContractMetaData.ABI))
	CrossChainEventName               = "LogMessagePublished"
	CrossChainEventID                 = MessageBusABI.Events[CrossChainEventName].ID
	ValueTransferEventName            = "ValueTransfer"
	ValueTransferEventID              = MessageBusABI.Events["ValueTransfer"].ID
	SequencerEnclaveGrantedEventID    = MgmtContractABI.Events["SequencerEnclaveGranted"].ID
	SequencerEnclaveRevokedEventID    = MgmtContractABI.Events["SequencerEnclaveRevoked"].ID
	NetworkSecretRequestedID          = MgmtContractABI.Events["NetworkSecretRequested"].ID
	NetworkSecretRespondedID          = MgmtContractABI.Events["NetworkSecretResponded"].ID
	RollupAddedID                     = MgmtContractABI.Events["RollupAdded"].ID
	DebugID                           = MgmtContractABI.Events["Debug"].ID
	ImportantContractAddressUpdatedID = MgmtContractABI.Events["ImportantContractAddressUpdated"].ID
)

func lazilyLogReceiptChecksum(block *types.Header, receipts types.Receipts, logger gethlog.Logger) {
	if logger.Enabled(context.Background(), gethlog.LevelTrace) {
		logger.Trace("Processing block", log.BlockHashKey, block.Hash(), "nr_rec", len(receipts), "Hash", receiptsHash(receipts))
	}
}

func receiptsHash(receipts types.Receipts) string {
	hasher := sha3.NewLegacyKeccak256().(crypto.KeccakState)
	hasher.Reset()
	for _, receipt := range receipts {
		var buffer bytes.Buffer
		err := receipt.EncodeRLP(&buffer)
		if err != nil {
			return err.Error()
		}
		hasher.Write(buffer.Bytes())
	}
	var hash gethcommon.Hash
	_, err := hasher.Read(hash[:])
	if err != nil {
		return err.Error()
	}
	return hash.Hex()
}

// filterLogsFromReceipts - filters the receipts for logs matching address, if provided and topic if provided.
func filterLogsFromReceipts(receipts types.Receipts, address *gethcommon.Address, topic *gethcommon.Hash) ([]types.Log, error) {
	logs := make([]types.Log, 0)

	for _, receipt := range receipts {
		if receipt.Status == 0 {
			continue
		}

		logsForReceipt, err := FilterLogsFromReceipt(receipt, address, topic)
		if err != nil {
			return logs, err
		}

		logs = append(logs, logsForReceipt...)
	}

	return logs, nil
}

// FilterLogsFromReceipt - filters the receipt for logs matching address, if provided and matching any of the provided topics.
func FilterLogsFromReceipt(receipt *types.Receipt, address *gethcommon.Address, topic *gethcommon.Hash) ([]types.Log, error) {
	logs := make([]types.Log, 0)

	if receipt == nil {
		return logs, errors.New("null receipt")
	}

	for _, log := range receipt.Logs {
		shouldSkip := false

		if address != nil && log.Address.Hex() != address.Hex() {
			shouldSkip = true
		}

		if topic != nil && log.Topics[0] != *topic {
			shouldSkip = true
		}

		if shouldSkip {
			continue
		}

		logs = append(logs, *log)
	}

	return logs, nil
}

// ConvertLogsToMessages - converts the logs of the event to messages. The logs should be filtered, otherwise fails.
func ConvertLogsToMessages(logs []types.Log, eventName string, messageBusABI abi.ABI) (common.CrossChainMessages, error) {
	messages := make(common.CrossChainMessages, 0)

	for _, log := range logs {
		var event MessageBus.MessageBusLogMessagePublished
		err := messageBusABI.UnpackIntoInterface(&event, eventName, log.Data)
		if err != nil {
			return nil, err
		}

		msg := createCrossChainMessage(event)
		messages = append(messages, msg)
	}

	return messages, nil
}

// createCrossChainMessage - Uses the logged event by the message bus to produce a cross chain message struct
func createCrossChainMessage(event MessageBus.MessageBusLogMessagePublished) MessageBus.StructsCrossChainMessage {
	return MessageBus.StructsCrossChainMessage{
		Sender:   event.Sender,
		Sequence: event.Sequence,
		Nonce:    event.Nonce,
		Topic:    event.Topic,
		Payload:  event.Payload,
	}
}

// ConvertLogsToValueTransfers - converts the logs of the event to messages. The logs should be filtered, otherwise fails.
func ConvertLogsToValueTransfers(logs []types.Log, eventName string, messageBusABI abi.ABI) (common.ValueTransferEvents, error) {
	messages := make(common.ValueTransferEvents, 0)

	for _, log := range logs {
		if len(log.Topics) != 3 {
			return nil, fmt.Errorf("invalid number of topics in log: %d", len(log.Topics))
		}
		var event MessageBus.MessageBusValueTransfer
		err := messageBusABI.UnpackIntoInterface(&event, eventName, log.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to unpack event data: %w\nData length: %d\nEvent: %s",
				err, len(log.Data), eventName)
		}

		messages = append(messages, common.ValueTransferEvent{
			Sender:   gethcommon.BytesToAddress(log.Topics[1].Bytes()),
			Receiver: gethcommon.BytesToAddress(log.Topics[2].Bytes()),
			Amount:   event.Amount,
			Sequence: event.Sequence,
		})
	}

	return messages, nil
}

type MerkleBatches []*core.Batch

func (mb MerkleBatches) Len() int {
	return len(mb)
}

func (mb MerkleBatches) EncodeIndex(index int, w *bytes.Buffer) {
	batch := mb[index]
	if err := rlp.Encode(w, batch.Header.CrossChainRoot); err != nil {
		panic(err)
	}
}

func (mb MerkleBatches) ForMerkleTree() [][]interface{} {
	values := make([][]interface{}, 0)
	for _, batch := range mb {
		val := []interface{}{
			batch.Header.CrossChainRoot,
		}
		values = append(values, val)
	}

	return values
}

type MessageStructs []MessageBus.StructsCrossChainMessage

func (ms MessageStructs) Len() int {
	return len(ms)
}

func (ms MessageStructs) EncodeIndex(index int, w *bytes.Buffer) {
	message := ms[index]
	if err := rlp.Encode(w, message); err != nil {
		panic(err)
	}
}

func (ms MessageStructs) ForMerkleTree() [][]interface{} {
	values := make([][]interface{}, 0)
	for idx := range ms {
		hashedVal := ms.HashPacked(idx)
		val := []interface{}{
			"m",
			hashedVal,
		}
		values = append(values, val)
	}
	return values
}

func (ms MessageStructs) HashPacked(index int) gethcommon.Hash {
	messageStruct := ms[index]

	addrType, _ := abi.NewType("address", "", nil)
	uint64Type, _ := abi.NewType("uint64", "", nil)
	uint32Type, _ := abi.NewType("uint32", "", nil)
	uint8Type, _ := abi.NewType("uint8", "", nil)
	bytesType, _ := abi.NewType("bytes", "", nil)
	args := abi.Arguments{
		{
			Type: addrType,
		},
		{
			Type: uint64Type,
		},
		{
			Type: uint32Type,
		},
		{
			Type: uint32Type,
		},
		{
			Type: bytesType,
		},
		{
			Type: uint8Type,
		},
	}

	// todo @siliev: err
	packed, err := args.Pack(messageStruct.Sender, messageStruct.Sequence, messageStruct.Nonce, messageStruct.Topic, messageStruct.Payload, messageStruct.ConsistencyLevel)
	if err != nil {
		panic(err)
	}
	hash := crypto.Keccak256Hash(packed)
	return hash
}

type ValueTransfers []common.ValueTransferEvent

func (vt ValueTransfers) Len() int {
	return len(vt)
}

func (vt ValueTransfers) EncodeIndex(index int, w *bytes.Buffer) {
	transfer := vt[index]
	if err := rlp.Encode(w, transfer); err != nil {
		panic(err)
	}
}

func (vt ValueTransfers) ForMerkleTree() [][]interface{} {
	values := make([][]interface{}, 0)
	for idx := range vt {
		hashedVal := vt.HashPacked(idx)
		val := []interface{}{
			"v", // [v, "0xblabla"]
			hashedVal,
		}
		values = append(values, val)
	}
	return values
}

func (vt ValueTransfers) HashPacked(index int) gethcommon.Hash {
	valueTransfer := vt[index]

	uint256Type, _ := abi.NewType("uint256", "", nil)
	uint64Type, _ := abi.NewType("uint64", "", nil)
	addrType, _ := abi.NewType("address", "", nil)

	args := abi.Arguments{
		{
			Type: addrType,
		},
		{
			Type: addrType,
		},
		{
			Type: uint256Type,
		},
		{
			Type: uint64Type,
		},
	}

	bytes, err := args.Pack(valueTransfer.Sender, valueTransfer.Receiver, valueTransfer.Amount, valueTransfer.Sequence)
	if err != nil {
		panic(err)
	}

	hash := crypto.Keccak256Hash(bytes)
	return hash
}

var CrossChainEncodings = []string{smt.SOL_STRING, smt.SOL_BYTES32}
