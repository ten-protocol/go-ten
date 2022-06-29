package rpc

import (
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/common"
	"github.com/obscuronet/obscuro-playground/go/common/rpc/generated"
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

func ToBlockSubmissionResponseMsg(response common.BlockSubmissionResponse) generated.BlockSubmissionResponseMsg {
	producedRollupMsg := ToExtRollupMsg(&response.ProducedRollup)
	return generated.BlockSubmissionResponseMsg{
		BlockHeader:           ToBlockHeaderMsg(response.BlockHeader),
		IngestedBlock:         response.IngestedBlock,
		BlockNotIngestedCause: response.BlockNotIngestedCause,
		ProducedRollup:        &producedRollupMsg,
		IngestedNewRollup:     response.FoundNewHead,
		RollupHead:            ToRollupHeaderMsg(response.RollupHead),
	}
}

func FromBlockSubmissionResponseMsg(msg *generated.BlockSubmissionResponseMsg) common.BlockSubmissionResponse {
	return common.BlockSubmissionResponse{
		BlockHeader:           FromBlockHeaderMsg(msg.GetBlockHeader()),
		IngestedBlock:         msg.IngestedBlock,
		BlockNotIngestedCause: msg.BlockNotIngestedCause,

		ProducedRollup: FromExtRollupMsg(msg.ProducedRollup),
		FoundNewHead:   msg.IngestedNewRollup,
		RollupHead:     FromRollupHeaderMsg(msg.RollupHead),
	}
}

func ToExtRollupMsg(rollup *common.ExtRollup) generated.ExtRollupMsg {
	if rollup.Header != nil {
		return generated.ExtRollupMsg{Header: ToRollupHeaderMsg(rollup.Header), Txs: rollup.EncryptedTxBlob}
	}

	return generated.ExtRollupMsg{Header: nil}
}

func ToRollupHeaderMsg(header *common.Header) *generated.HeaderMsg {
	if header == nil {
		return nil
	}
	var headerMsg generated.HeaderMsg
	withdrawalMsgs := make([]*generated.WithdrawalMsg, 0)
	for _, withdrawal := range header.Withdrawals {
		withdrawalMsg := generated.WithdrawalMsg{Amount: withdrawal.Amount, Recipient: withdrawal.Recipient.Bytes(), Contract: withdrawal.Contract.Bytes()}
		withdrawalMsgs = append(withdrawalMsgs, &withdrawalMsg)
	}

	headerMsg = generated.HeaderMsg{
		ParentHash:  header.ParentHash.Bytes(),
		Node:        header.Agg.Bytes(),
		Nonce:       header.Nonce,
		Proof:       header.L1Proof.Bytes(),
		Root:        header.Root.Bytes(),
		TxHash:      header.TxHash.Bytes(),
		Number:      header.Number.Uint64(),
		Withdrawals: withdrawalMsgs,
		Bloom:       header.Bloom.Bytes(),
		ReceiptHash: header.ReceiptHash.Bytes(),
		Extra:       header.Extra,
		R:           header.R.Bytes(),
		S:           header.S.Bytes(),
	}

	return &headerMsg
}

func FromExtRollupMsg(msg *generated.ExtRollupMsg) common.ExtRollup {
	if msg.Header == nil {
		return common.ExtRollup{
			Header: nil,
		}
	}

	return common.ExtRollup{
		Header:          FromRollupHeaderMsg(msg.Header),
		EncryptedTxBlob: msg.Txs,
	}
}

func FromRollupHeaderMsg(header *generated.HeaderMsg) *common.Header {
	if header == nil {
		return nil
	}
	withdrawals := make([]common.Withdrawal, 0)
	for _, withdrawalMsg := range header.Withdrawals {
		recipient := gethcommon.BytesToAddress(withdrawalMsg.Recipient)
		contract := gethcommon.BytesToAddress(withdrawalMsg.Contract)
		withdrawal := common.Withdrawal{Amount: withdrawalMsg.Amount, Recipient: recipient, Contract: contract}
		withdrawals = append(withdrawals, withdrawal)
	}

	r := &big.Int{}
	s := &big.Int{}
	return &common.Header{
		ParentHash:  gethcommon.BytesToHash(header.ParentHash),
		Agg:         gethcommon.BytesToAddress(header.Node),
		Nonce:       header.Nonce,
		L1Proof:     gethcommon.BytesToHash(header.Proof),
		Root:        gethcommon.BytesToHash(header.Root),
		Number:      big.NewInt(int64(header.Number)),
		TxHash:      gethcommon.BytesToHash(header.TxHash),
		Withdrawals: withdrawals,
		ReceiptHash: gethcommon.BytesToHash(header.ReceiptHash),
		Bloom:       types.BytesToBloom(header.Bloom),
		Extra:       header.Extra,
		R:           r.SetBytes(header.R),
		S:           s.SetBytes(header.S),
	}
}

func FromBlockHeaderMsg(msg *generated.BlockHeaderMsg) *types.Header {
	if msg == nil {
		return nil
	}
	return &types.Header{
		ParentHash:  gethcommon.BytesToHash(msg.ParentHash),
		UncleHash:   gethcommon.BytesToHash(msg.UncleHash),
		Coinbase:    gethcommon.BytesToAddress(msg.Coinbase),
		Root:        gethcommon.BytesToHash(msg.Root),
		TxHash:      gethcommon.BytesToHash(msg.TxHash),
		ReceiptHash: gethcommon.BytesToHash(msg.ReceiptHash),
		Bloom:       types.BytesToBloom(msg.Bloom),
		Difficulty:  big.NewInt(int64(msg.Difficulty)),
		Number:      big.NewInt(int64(msg.Number)),
		GasLimit:    msg.GasLimit,
		GasUsed:     msg.GasUsed,
		Time:        msg.Time,
		Extra:       msg.Extra,
		MixDigest:   gethcommon.BytesToHash(msg.MixDigest),
		Nonce:       types.EncodeNonce(msg.Nonce),
		BaseFee:     big.NewInt(int64(msg.BaseFee)),
	}
}

func ToBlockHeaderMsg(header *types.Header) *generated.BlockHeaderMsg {
	if header == nil {
		return nil
	}
	baseFee := uint64(0)
	if header.BaseFee != nil {
		baseFee = header.BaseFee.Uint64()
	}
	return &generated.BlockHeaderMsg{
		ParentHash:  header.ParentHash.Bytes(),
		UncleHash:   header.UncleHash.Bytes(),
		Coinbase:    header.Coinbase.Bytes(),
		Root:        header.Root.Bytes(),
		TxHash:      header.TxHash.Bytes(),
		ReceiptHash: header.ReceiptHash.Bytes(),
		Bloom:       header.Bloom.Bytes(),
		Difficulty:  header.Difficulty.Uint64(),
		Number:      header.Number.Uint64(),
		GasLimit:    header.GasLimit,
		GasUsed:     header.GasUsed,
		Time:        header.Time,
		Extra:       header.Extra,
		MixDigest:   header.MixDigest.Bytes(),
		Nonce:       header.Nonce.Uint64(),
		BaseFee:     baseFee,
	}
}
