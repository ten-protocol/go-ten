package crosschain

import (
	"bytes"
	"errors"
	"strings"

	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave/core"

	smt "github.com/FantasyJony/openzeppelin-merkle-tree-go/standard_merkle_tree"
	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/contracts/generated/MessageBus"
	"github.com/ten-protocol/go-ten/go/common"
	"golang.org/x/crypto/sha3"
)

var (
	MessageBusABI, _       = abi.JSON(strings.NewReader(MessageBus.MessageBusMetaData.ABI))
	CrossChainEventName    = "LogMessagePublished"
	CrossChainEventID      = MessageBusABI.Events[CrossChainEventName].ID
	ValueTransferEventName = "ValueTransfer"
	ValueTransferEventID   = MessageBusABI.Events["ValueTransfer"].ID
)

func lazilyLogReceiptChecksum(block *common.L1Block, receipts types.Receipts, logger gethlog.Logger) {
	logger.Trace("Processing block", log.BlockHashKey, block.Hash(), "nr_rec", len(receipts), "Hash",
		gethlog.Lazy{Fn: func() string {
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
		}})
}

/*
func lazilyLogChecksum(msg string, transactions types.Transactions, logger gethlog.Logger) {
	logger.Trace(msg, "Hash",
		gethlog.Lazy{Fn: func() string {
			hasher := sha3.NewLegacyKeccak256().(crypto.KeccakState)
			hasher.Reset()
			for _, tx := range transactions {
				var buffer bytes.Buffer
				err := tx.EncodeRLP(&buffer)
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
		}})
}
*/

// filterLogsFromReceipts - filters the receipts for logs matching address, if provided and topic if provided.
func filterLogsFromReceipts(receipts types.Receipts, address *gethcommon.Address, topics []*gethcommon.Hash) ([]types.Log, error) {
	logs := make([]types.Log, 0)

	for _, receipt := range receipts {
		if receipt.Status == 0 {
			continue
		}

		logsForReceipt, err := filterLogsFromReceipt(receipt, address, topics)
		if err != nil {
			return logs, err
		}

		logs = append(logs, logsForReceipt...)
	}

	return logs, nil
}

// filterLogsFromReceipt - filters the receipt for logs matching address, if provided and topic if provided.
func filterLogsFromReceipt(receipt *types.Receipt, address *gethcommon.Address, topics []*gethcommon.Hash) ([]types.Log, error) {
	logs := make([]types.Log, 0)

	if receipt == nil {
		return logs, errors.New("null receipt")
	}

	for _, log := range receipt.Logs {
		shouldSkip := false

		if address != nil && log.Address.Hex() != address.Hex() {
			shouldSkip = true
		}

		for _, topic := range topics {
			if topic != nil && log.Topics[0] != *topic {
				shouldSkip = true
			}
		}

		if shouldSkip {
			continue
		}

		logs = append(logs, *log)
	}

	return logs, nil
}

// convertLogsToMessages - converts the logs of the event to messages. The logs should be filtered, otherwise fails.
func convertLogsToMessages(logs []types.Log, eventName string, messageBusABI abi.ABI) (common.CrossChainMessages, error) {
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

// convertLogsToMessages - converts the logs of the event to messages. The logs should be filtered, otherwise fails.
func convertLogsToValueTransfers(logs []types.Log, eventName string, messageBusABI abi.ABI) (common.ValueTransferEvents, error) {
	messages := make(common.ValueTransferEvents, 0)

	for _, log := range logs {
		var event MessageBus.MessageBusValueTransfer
		err := messageBusABI.UnpackIntoInterface(&event, eventName, log.Data)
		if err != nil {
			return nil, err
		}

		messages = append(messages, common.ValueTransferEvent{
			Sender:   event.Sender,
			Receiver: event.Receiver,
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
	if err := rlp.Encode(w, batch.Header.TransfersTree); err != nil {
		panic(err)
	}
}

func (mb MerkleBatches) ForMerkleTree() [][]interface{} {
	values := make([][]interface{}, 0)
	for _, batch := range mb {
		val := []interface{}{
			batch.Header.TransfersTree,
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
	for idx, _ := range ms {
		hashedVal := ms.HashPacked(idx)
		val := []interface{}{
			"message",
			hashedVal,
		}
		values = append(values, val)
	}
	return values
}

func (ms MessageStructs) HashPacked(index int) gethcommon.Hash {
	messageStruct := ms[index]
	/*	Sender           common.Address
		Sequence         uint64
		Nonce            uint32
		Topic            uint32
		Payload          []byte
		ConsistencyLevel uint8 */

	addrType, _ := abi.NewType("address", "", nil)
	uint64Type, _ := abi.NewType("uint64", "", nil)
	uint32Type, _ := abi.NewType("uint32", "", nil)
	uint8Type, _ := abi.NewType("uint32", "", nil)
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

	//todo @siliev: err
	packed, _ := args.Pack(messageStruct.Sender, messageStruct.Sequence, messageStruct.Nonce, messageStruct.Topic, messageStruct.Payload, messageStruct.ConsistencyLevel)
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
	for idx, _ := range vt {
		hashedVal := vt.HashPacked(idx)
		val := []interface{}{
			"value",
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

	bytes, _ := args.Pack(valueTransfer.Sender, valueTransfer.Receiver, valueTransfer.Amount, valueTransfer.Sequence)

	hash := crypto.Keccak256Hash(bytes)
	return hash
}

var CrossChainEncodings = []string{smt.SOL_STRING, smt.SOL_BYTES32}
