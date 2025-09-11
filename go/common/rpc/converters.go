package rpc

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/kzg4844"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/rpc/generated"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Functions to convert classes that need to be sent between the host and the enclave to and from their equivalent
// Protobuf message classes.

func ToAttestationReportMsg(report *common.AttestationReport) generated.AttestationReportMsg {
	return generated.AttestationReportMsg{Report: report.Report, PubKey: report.PubKey, EnclaveID: report.EnclaveID.Bytes(), HostAddress: report.HostAddress}
}

func FromAttestationReportMsg(msg *generated.AttestationReportMsg) *common.AttestationReport {
	return &common.AttestationReport{
		Report:      msg.Report,
		PubKey:      msg.PubKey,
		EnclaveID:   gethcommon.BytesToAddress(msg.EnclaveID),
		HostAddress: msg.HostAddress,
	}
}

func ToSecretRespMsg(responses []*common.ProducedSecretResponse) []*generated.SecretResponseMsg {
	respMsgs := make([]*generated.SecretResponseMsg, len(responses))

	for i, resp := range responses {
		msg := generated.SecretResponseMsg{
			Secret:      resp.Secret,
			RequesterID: resp.RequesterID.Bytes(),
			AttesterID:  resp.AttesterID.Bytes(),
			HostAddress: resp.HostAddress,
			Signature:   resp.Signature,
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
			AttesterID:  gethcommon.BytesToAddress(msgResp.AttesterID),
			HostAddress: msgResp.HostAddress,
			Signature:   msgResp.Signature,
		}
		respList[i] = &r
	}
	return respList
}

func ToBlockSubmissionResponseMsg(response *common.BlockSubmissionResponse) (*generated.BlockSubmissionResponseMsg, error) {
	if response == nil {
		return nil, fmt.Errorf("no response that could be converted to a message")
	}

	msg := &generated.BlockSubmissionResponseMsg{
		ProducedSecretResponses: ToSecretRespMsg(response.ProducedSecretResponses),
	}
	for _, metadata := range response.RollupMetadata {
		msg.RollupMetadata = append(msg.RollupMetadata, &generated.ExtRollupMetadataResponseMsg{
			CrossChainTree: metadata.CrossChainTree,
		})
	}

	return msg, nil
}

func FromBlockSubmissionErrorMsg(msg *generated.BlockSubmissionErrorMsg) *errutil.BlockRejectError {
	if msg == nil {
		return nil
	}

	return &errutil.BlockRejectError{
		L1Head:  gethcommon.BytesToHash(msg.L1Head),
		Wrapped: errors.New(msg.Cause),
	}
}

func FromBlockSubmissionResponseMsg(msg *generated.BlockSubmissionResponseMsg) (*common.BlockSubmissionResponse, error) {
	rollupMetadata := make([]common.ExtRollupMetadata, len(msg.RollupMetadata))
	for i, metadata := range msg.RollupMetadata {
		rollupMetadata[i] = common.ExtRollupMetadata{
			CrossChainTree: metadata.CrossChainTree,
		}
	}
	return &common.BlockSubmissionResponse{
		ProducedSecretResponses: FromSecretRespMsg(msg.ProducedSecretResponses),
		RollupMetadata:          rollupMetadata,
		RejectError:             FromBlockSubmissionErrorMsg(msg.Error),
	}, nil
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
		ParentHash:       header.ParentHash.Bytes(),
		Proof:            header.L1Proof.Bytes(),
		Root:             header.Root.Bytes(),
		TxHash:           header.TxHash.Bytes(),
		Number:           header.Number.Uint64(),
		SequencerOrderNo: header.SequencerOrderNo.Uint64(),
		ReceiptHash:      header.ReceiptHash.Bytes(),
		Extra:            header.Extra,
		Signature:        header.Signature,
		GasLimit:         header.GasLimit,
		GasUsed:          header.GasUsed,
		Time:             header.Time,
		BaseFee:          baseFee,
		CrossChainRoot:   header.CrossChainRoot.Bytes(),
		Coinbase:         header.Coinbase.Bytes(),
		CrossChainTree:   header.CrossChainTree,
		PayloadHash:      header.PayloadHash.Bytes(),
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

	return &common.BatchHeader{
		ParentHash:       gethcommon.BytesToHash(header.ParentHash),
		L1Proof:          gethcommon.BytesToHash(header.Proof),
		Root:             gethcommon.BytesToHash(header.Root),
		TxHash:           gethcommon.BytesToHash(header.TxHash),
		Number:           big.NewInt(int64(header.Number)),
		SequencerOrderNo: big.NewInt(int64(header.SequencerOrderNo)),
		ReceiptHash:      gethcommon.BytesToHash(header.ReceiptHash),
		Extra:            header.Extra,
		Signature:        header.Signature,
		GasLimit:         header.GasLimit,
		GasUsed:          header.GasUsed,
		Time:             header.Time,
		CrossChainRoot:   gethcommon.BytesToHash(header.CrossChainRoot),
		BaseFee:          big.NewInt(0).SetUint64(header.BaseFee),
		Coinbase:         gethcommon.BytesToAddress(header.Coinbase),
		CrossChainTree:   header.CrossChainTree,
		PayloadHash:      gethcommon.BytesToHash(header.PayloadHash),
	}
}

func ToExtRollupMsg(rollup *common.ExtRollup) generated.ExtRollupMsg {
	if rollup == nil || rollup.Header == nil {
		return generated.ExtRollupMsg{}
	}

	return generated.ExtRollupMsg{Header: ToRollupHeaderMsg(rollup.Header), BatchPayloads: rollup.BatchPayloads, CalldataRollupHeader: rollup.CalldataRollupHeader}
}

func ToRollupHeaderMsg(header *common.RollupHeader) *generated.RollupHeaderMsg {
	if header == nil {
		return nil
	}
	headerMsg := generated.RollupHeaderMsg{
		CompressionL1Head:   header.CompressionL1Head.Bytes(),
		CompressionL1Number: header.CompressionL1Number.Bytes(),
		CrossChainRoot:      header.CrossChainRoot.Bytes(),
		FirstBatchSeqNo:     header.FirstBatchSeqNo,
		LastBatchSeqNo:      header.LastBatchSeqNo,
		LastBatchHash:       header.LastBatchHash.Bytes(),
	}

	return &headerMsg
}

func ToRollupDataMsg(rollupData *common.PublicRollupMetadata) generated.PublicRollupDataMsg {
	if rollupData == nil {
		return generated.PublicRollupDataMsg{}
	}

	return generated.PublicRollupDataMsg{StartSeq: rollupData.FirstBatchSequence.Uint64(), Timestamp: rollupData.StartTime}
}

func FromRollupDataMsg(msg *generated.PublicRollupDataMsg) (*common.PublicRollupMetadata, error) {
	if msg.Timestamp == 0 {
		return nil, fmt.Errorf("timestamp on the rollup can not be zero")
	}

	if msg.StartSeq == 0 {
		return &common.PublicRollupMetadata{
			FirstBatchSequence: nil,
			StartTime:          msg.Timestamp,
		}, nil
	}

	return &common.PublicRollupMetadata{
		FirstBatchSequence: big.NewInt(int64(msg.StartSeq)),
		StartTime:          msg.Timestamp,
	}, nil
}

func ToBlobMsgs(blobs []*kzg4844.Blob) []*generated.BlobMsg {
	if blobs == nil {
		return nil
	}
	msgs := make([]*generated.BlobMsg, len(blobs))
	for i, blob := range blobs {
		msgs[i] = &generated.BlobMsg{
			Blob: blob[:],
		}
	}
	return msgs
}

func FromBlobMsgs(msgs []*generated.BlobMsg) []*kzg4844.Blob {
	if msgs == nil {
		return nil
	}
	blobs := make([]*kzg4844.Blob, len(msgs))
	for i, msg := range msgs {
		var blob kzg4844.Blob
		copy(blob[:], msg.Blob)
		blobs[i] = &blob
	}
	return blobs
}
