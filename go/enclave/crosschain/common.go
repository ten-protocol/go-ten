package crosschain

import (
	"bytes"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/contracts/generated/MessageBus"
	"github.com/obscuronet/go-obscuro/go/common"
	"golang.org/x/crypto/sha3"
)

var (
	MessageBusABI, _       = abi.JSON(strings.NewReader(MessageBus.MessageBusMetaData.ABI))
	CrossChainEventName    = "LogMessagePublished"
	CrossChainEventID      = MessageBusABI.Events[CrossChainEventName].ID
	ValueTransferEventName = "ValueTransfer"
	ValueTransferEventID   = MessageBusABI.Events["ValueTransfer"].ID
)

func lazilyLogReceiptChecksum(msg string, receipts types.Receipts, logger gethlog.Logger) {
	logger.Trace(msg, "Hash",
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
func filterLogsFromReceipts(receipts types.Receipts, address *gethcommon.Address, topic *gethcommon.Hash) ([]types.Log, error) {
	logs := make([]types.Log, 0)

	for _, receipt := range receipts {
		if receipt.Status == 0 {
			continue
		}

		logsForReceipt, err := filterLogsFromReceipt(receipt, address, topic)
		if err != nil {
			return logs, err
		}

		logs = append(logs, logsForReceipt...)
	}

	return logs, nil
}

// filterLogsFromReceipt - filters the receipt for logs matching address, if provided and topic if provided.
func filterLogsFromReceipt(receipt *types.Receipt, address *gethcommon.Address, topic *gethcommon.Hash) ([]types.Log, error) {
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
		})
	}

	return messages, nil
}
