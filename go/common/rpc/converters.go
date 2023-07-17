package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/contracts/generated/MessageBus"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/rpc/generated"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Functions to convert classes that need to be sent between the host and the enclave to and from their equivalent
// Protobuf message classes.

func ToAttestationReportMsg(report *common.AttestationReport) generated.AttestationReportMsg {
	return generated.AttestationReportMsg{Report: report.Report, PubKey: report.PubKey, Owner: report.Owner.Bytes(), HostAddress: report.HostAddress}
}

func FromAttestationReportMsg(msg *generated.AttestationReportMsg) *common.AttestationReport {
	return &common.AttestationReport{
		Report:      msg.Report,
		PubKey:      msg.PubKey,
		Owner:       gethcommon.BytesToAddress(msg.Owner),
		HostAddress: msg.HostAddress,
	}
}

func ToBlockSubmissionResponseMsg(response *common.BlockSubmissionResponse) (*generated.BlockSubmissionResponseMsg, error) {
	var subscribedLogBytes []byte
	var err error

	if response == nil {
		return nil, fmt.Errorf("no response that could be converted to a message")
	}

	producedBatchMsg := ToExtBatchMsg(response.ProducedBatch)
	producedRollupMsg := ToExtRollupMsg(response.ProducedRollup)

	msg := &generated.BlockSubmissionResponseMsg{
		ProducedBatch:           &producedBatchMsg,
		ProducedRollup:          &producedRollupMsg,
		SubscribedLogs:          subscribedLogBytes,
		ProducedSecretResponses: ToSecretRespMsg(response.ProducedSecretResponses),
	}

	if response.SubscribedLogs != nil {
		msg.SubscribedLogs, err = json.Marshal(response.SubscribedLogs)
		if err != nil {
			return &generated.BlockSubmissionResponseMsg{}, fmt.Errorf("could not marshal subscribed logs to JSON. Cause: %w", err)
		}
	}

	return msg, nil
}

func ToSecretRespMsg(responses []*common.ProducedSecretResponse) []*generated.SecretResponseMsg {
	respMsgs := make([]*generated.SecretResponseMsg, len(responses))

	for i, resp := range responses {
		msg := generated.SecretResponseMsg{
			Secret:      resp.Secret,
			RequesterID: resp.RequesterID.Bytes(),
			HostAddress: resp.HostAddress,
		}
		respMsgs[i] = &msg
	}

	return respMsgs
}

func FromSecretRespMsg(secretResponses []*generated.SecretResponseMsg) []*common.ProducedSecretResponse {
	respList := make([]*common.ProducedSecretResponse, len(secretResponses))

	for i, msgResp := range secretResponses {
		r := common.ProducedSecretResponse{
			Secret:      msgResp.Secret,
			RequesterID: gethcommon.BytesToAddress(msgResp.RequesterID),
			HostAddress: msgResp.HostAddress,
		}
		respList[i] = &r
	}
	return respList
}

func FromBlockSubmissionResponseMsg(msg *generated.BlockSubmissionResponseMsg) (*common.BlockSubmissionResponse, error) {
	if msg.Error != nil {
		return nil, &errutil.BlockRejectError{
			L1Head:  gethcommon.BytesToHash(msg.Error.L1Head),
			Wrapped: errors.New(msg.Error.Cause),
		}
	}
	var subscribedLogs map[rpc.ID][]byte
	if msg.SubscribedLogs != nil {
		if err := json.Unmarshal(msg.SubscribedLogs, &subscribedLogs); err != nil {
			return nil, fmt.Errorf("could not unmarshal subscribed logs from submission response JSON. Cause: %w", err)
		}
	}

	return &common.BlockSubmissionResponse{
		ProducedBatch:           FromExtBatchMsg(msg.ProducedBatch),
		ProducedRollup:          FromExtRollupMsg(msg.ProducedRollup),
		SubscribedLogs:          subscribedLogs,
		ProducedSecretResponses: FromSecretRespMsg(msg.ProducedSecretResponses),
	}, nil
}

func ToCrossChainMsgs(messages []MessageBus.StructsCrossChainMessage) []*generated.CrossChainMsg {
	generatedMessages := make([]*generated.CrossChainMsg, 0)

	for _, message := range messages {
		generatedMessages = append(generatedMessages, &generated.CrossChainMsg{
			Sender:   message.Sender.Bytes(),
			Sequence: message.Sequence,
			Nonce:    message.Nonce,
			Topic:    message.Topic,
			Payload:  message.Payload,
		})
	}

	return generatedMessages
}

func FromCrossChainMsgs(messages []*generated.CrossChainMsg) []MessageBus.StructsCrossChainMessage {
	outMessages := make([]MessageBus.StructsCrossChainMessage, 0)

	for _, message := range messages {
		outMessages = append(outMessages, MessageBus.StructsCrossChainMessage{
			Sender:   gethcommon.BytesToAddress(message.Sender),
			Sequence: message.Sequence,
			Nonce:    message.Nonce,
			Topic:    message.Topic,
			Payload:  message.Payload,
		})
	}

	return outMessages
}

func ToExtBatchMsg(batch *common.ExtBatch) generated.ExtBatchMsg {
	if batch == nil || batch.Header == nil {
		return generated.ExtBatchMsg{}
	}

	txHashBytes := make([][]byte, len(batch.TxHashes))
	for idx, txHash := range batch.TxHashes {
		txHashBytes[idx] = txHash.Bytes()
	}

	return generated.ExtBatchMsg{Header: ToBatchHeaderMsg(batch.Header), TxHashes: txHashBytes, Txs: batch.EncryptedTxBlob}
}

func ToBatchHeaderMsg(header *common.BatchHeader) *generated.BatchHeaderMsg {
	if header == nil {
		return nil
	}
	var headerMsg generated.BatchHeaderMsg

	baseFee := uint64(0)
	if header.BaseFee != nil {
		baseFee = header.BaseFee.Uint64()
	}
	headerMsg = generated.BatchHeaderMsg{
		ParentHash:                  header.ParentHash.Bytes(),
		Proof:                       header.L1Proof.Bytes(),
		Root:                        header.Root.Bytes(),
		TxHash:                      header.TxHash.Bytes(),
		Number:                      header.Number.Uint64(),
		SequencerOrderNo:            header.SequencerOrderNo.Uint64(),
		ReceiptHash:                 header.ReceiptHash.Bytes(),
		Extra:                       header.Extra,
		R:                           header.R.Bytes(),
		S:                           header.S.Bytes(),
		GasLimit:                    header.GasLimit,
		GasUsed:                     header.GasUsed,
		Time:                        header.Time,
		BaseFee:                     baseFee,
		CrossChainMessages:          ToCrossChainMsgs(header.CrossChainMessages),
		LatestInboundCrossChainHash: header.LatestInboundCrossChainHash.Bytes(),
	}

	if header.LatestInboundCrossChainHeight != nil {
		headerMsg.LatestInboundCrossChainHeight = header.LatestInboundCrossChainHeight.Bytes()
	}

	return &headerMsg
}

func FromExtBatchMsg(msg *generated.ExtBatchMsg) *common.ExtBatch {
	if msg.Header == nil {
		return &common.ExtBatch{
			Header: nil,
		}
	}

	// We recreate the transaction hashes.
	txHashes := make([]gethcommon.Hash, len(msg.TxHashes))
	for idx, bytes := range msg.TxHashes {
		txHashes[idx] = gethcommon.BytesToHash(bytes)
	}

	return &common.ExtBatch{
		Header:          FromBatchHeaderMsg(msg.Header),
		TxHashes:        txHashes,
		EncryptedTxBlob: msg.Txs,
	}
}

func FromBatchHeaderMsg(header *generated.BatchHeaderMsg) *common.BatchHeader {
	if header == nil {
		return nil
	}

	r := &big.Int{}
	s := &big.Int{}
	return &common.BatchHeader{
		ParentHash:                    gethcommon.BytesToHash(header.ParentHash),
		L1Proof:                       gethcommon.BytesToHash(header.Proof),
		Root:                          gethcommon.BytesToHash(header.Root),
		TxHash:                        gethcommon.BytesToHash(header.TxHash),
		Number:                        big.NewInt(int64(header.Number)),
		SequencerOrderNo:              big.NewInt(int64(header.SequencerOrderNo)),
		ReceiptHash:                   gethcommon.BytesToHash(header.ReceiptHash),
		Extra:                         header.Extra,
		R:                             r.SetBytes(header.R),
		S:                             s.SetBytes(header.S),
		GasLimit:                      header.GasLimit,
		GasUsed:                       header.GasUsed,
		Time:                          header.Time,
		BaseFee:                       big.NewInt(int64(header.BaseFee)),
		CrossChainMessages:            FromCrossChainMsgs(header.CrossChainMessages),
		LatestInboundCrossChainHash:   gethcommon.BytesToHash(header.LatestInboundCrossChainHash),
		LatestInboundCrossChainHeight: big.NewInt(0).SetBytes(header.LatestInboundCrossChainHeight),
	}
}

func ToExtRollupMsg(rollup *common.ExtRollup) generated.ExtRollupMsg {
	if rollup == nil || rollup.Header == nil {
		return generated.ExtRollupMsg{}
	}

	return generated.ExtRollupMsg{Header: ToRollupHeaderMsg(rollup.Header), BatchPayloads: rollup.BatchPayloads, BatchHeaders: rollup.BatchHeaders}
}

func ToRollupHeaderMsg(header *common.RollupHeader) *generated.RollupHeaderMsg {
	if header == nil {
		return nil
	}
	headerMsg := generated.RollupHeaderMsg{
		ParentHash:         header.ParentHash.Bytes(),
		Proof:              header.L1Proof.Bytes(),
		ProofNumber:        header.L1ProofNumber.Uint64(),
		Number:             header.Number.Uint64(),
		R:                  header.R.Bytes(),
		S:                  header.S.Bytes(),
		Time:               header.Time,
		Coinbase:           header.Coinbase.Bytes(),
		CrossChainMessages: ToCrossChainMsgs(header.CrossChainMessages),
		LastBatchSeqNo:     header.LastBatchSeqNo,
	}

	return &headerMsg
}

func FromExtRollupMsg(msg *generated.ExtRollupMsg) *common.ExtRollup {
	if msg.Header == nil {
		return &common.ExtRollup{
			Header: nil,
		}
	}

	return &common.ExtRollup{
		Header:        FromRollupHeaderMsg(msg.Header),
		BatchPayloads: msg.BatchPayloads,
		BatchHeaders:  msg.BatchHeaders,
	}
}

func FromRollupHeaderMsg(header *generated.RollupHeaderMsg) *common.RollupHeader {
	if header == nil {
		return nil
	}

	r := &big.Int{}
	s := &big.Int{}
	return &common.RollupHeader{
		ParentHash:         gethcommon.BytesToHash(header.ParentHash),
		L1Proof:            gethcommon.BytesToHash(header.Proof),
		L1ProofNumber:      big.NewInt(int64(header.ProofNumber)),
		Number:             big.NewInt(int64(header.Number)),
		R:                  r.SetBytes(header.R),
		S:                  s.SetBytes(header.S),
		Time:               header.Time,
		Coinbase:           gethcommon.BytesToAddress(header.Coinbase),
		CrossChainMessages: FromCrossChainMsgs(header.CrossChainMessages),
		LastBatchSeqNo:     header.LastBatchSeqNo,
	}
}
