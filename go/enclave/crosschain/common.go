package crosschain

import (
	"bytes"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/obscuronet/go-obscuro/contracts/messagebuscontract/generated/MessageBus"
	"github.com/obscuronet/go-obscuro/go/common"
	"golang.org/x/crypto/sha3"
)

func lazilyLogReceiptChecksum(msg string, receipts types.Receipts, logger gethlog.Logger) {
	logger.Trace(msg, "Hash",
		gethlog.Lazy{Fn: func() string {
			hasher := sha3.NewLegacyKeccak256().(crypto.KeccakState)
			hasher.Reset()
			for _, receipt := range receipts {
				var buffer bytes.Buffer
				receipt.EncodeRLP(&buffer)
				hasher.Write(buffer.Bytes())
			}
			var hash gethcommon.Hash
			hasher.Read(hash[:])
			return hash.Hex()
		}})
}

func lazilyLogChecksum(msg string, transactions types.Transactions, logger gethlog.Logger) {
	logger.Trace(msg, "Hash",
		gethlog.Lazy{Fn: func() string {
			hasher := sha3.NewLegacyKeccak256().(crypto.KeccakState)
			hasher.Reset()
			for _, tx := range transactions {
				var buffer bytes.Buffer
				tx.EncodeRLP(&buffer)
				hasher.Write(buffer.Bytes())
			}
			var hash gethcommon.Hash
			hasher.Read(hash[:])
			return hash.Hex()
		}})
}

func filterLogsFromReceipts(receipts types.Receipts, address *gethcommon.Address, topic *gethcommon.Hash) ([]types.Log, error) {
	logs := make([]types.Log, 0)

	for _, receipt := range receipts {
		logsForReceipt, err := filterLogsFromReceipt(receipt, address, topic)
		if err != nil {
			return logs, err
		}

		logs = append(logs, logsForReceipt...)
	}

	return logs, nil
}

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

func convertLogsToMessages(logs []types.Log, eventName string, contractABI abi.ABI) (common.CrossChainMessages, error) {
	messages := make(common.CrossChainMessages, 0)

	for _, log := range logs {
		var event MessageBus.MessageBusLogMessagePublished
		err := contractABI.UnpackIntoInterface(&event, eventName, log.Data)
		if err != nil {
			return messages, err
		}

		msg := createCrossChainMessage(event)
		messages = append(messages, msg)
	}

	return messages, nil
}

func createCrossChainMessage(event MessageBus.MessageBusLogMessagePublished) MessageBus.StructsCrossChainMessage {
	return MessageBus.StructsCrossChainMessage{
		Sender:   event.Sender,
		Sequence: event.Sequence,
		Nonce:    event.Nonce,
		Topic:    event.Topic,
		Payload:  event.Payload,
	}
}

func VerifyReceiptHash(block *common.L1Block, receipts common.L1Receipts) bool {
	if len(receipts) == 0 {
		return block.ReceiptHash() == types.EmptyRootHash
	}

	calculatedHash := types.DeriveSha(receipts, &trie.StackTrie{})
	expectedHash := block.ReceiptHash()

	return calculatedHash == expectedHash
}
