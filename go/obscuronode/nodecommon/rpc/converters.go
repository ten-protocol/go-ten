package rpc

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon/rpc/generated"
)

// Functions to convert classes that need to be sent between the host and the enclave to and from their equivalent
// Protobuf message classes.

func ToAttestationReportMsg(report obscurocommon.AttestationReport) generated.AttestationReportMsg {
	return generated.AttestationReportMsg{Owner: report.Owner.Bytes()}
}

func FromAttestationReportMsg(msg *generated.AttestationReportMsg) obscurocommon.AttestationReport {
	return obscurocommon.AttestationReport{Owner: common.BytesToAddress(msg.Owner)}
}

func ToBlockSubmissionResponseMsg(response nodecommon.BlockSubmissionResponse) generated.BlockSubmissionResponseMsg {
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

func FromBlockSubmissionResponseMsg(msg *generated.BlockSubmissionResponseMsg) nodecommon.BlockSubmissionResponse {
	return nodecommon.BlockSubmissionResponse{
		BlockHeader:           FromBlockHeaderMsg(msg.GetBlockHeader()),
		IngestedBlock:         msg.IngestedBlock,
		BlockNotIngestedCause: msg.BlockNotIngestedCause,

		ProducedRollup: FromExtRollupMsg(msg.ProducedRollup),
		FoundNewHead:   msg.IngestedNewRollup,
		RollupHead:     FromRollupHeaderMsg(msg.RollupHead),
	}
}

func ToExtRollupMsg(rollup *nodecommon.ExtRollup) generated.ExtRollupMsg {
	if rollup.Header != nil {
		txs := make([][]byte, 0)
		for _, tx := range rollup.Txs {
			txs = append(txs, tx)
		}

		return generated.ExtRollupMsg{Header: ToRollupHeaderMsg(rollup.Header), Txs: txs}
	}

	return generated.ExtRollupMsg{Header: nil}
}

func ToRollupHeaderMsg(header *nodecommon.Header) *generated.HeaderMsg {
	if header == nil {
		return nil
	}
	var headerMsg generated.HeaderMsg
	withdrawalMsgs := make([]*generated.WithdrawalMsg, 0)
	for _, withdrawal := range header.Withdrawals {
		withdrawalMsg := generated.WithdrawalMsg{Amount: withdrawal.Amount, Address: withdrawal.Address.Bytes()}
		withdrawalMsgs = append(withdrawalMsgs, &withdrawalMsg)
	}

	headerMsg = generated.HeaderMsg{
		ParentHash:  header.ParentHash.Bytes(),
		Agg:         header.Agg.Bytes(),
		Nonce:       header.Nonce,
		L1Proof:     header.L1Proof.Bytes(),
		StateRoot:   header.State.Bytes(),
		Height:      header.Number,
		Withdrawals: withdrawalMsgs,
	}

	return &headerMsg
}

func FromExtRollupMsg(msg *generated.ExtRollupMsg) nodecommon.ExtRollup {
	if msg.Header == nil {
		return nodecommon.ExtRollup{
			Header: nil,
		}
	}

	txs := make([]nodecommon.EncryptedTx, 0)
	for _, tx := range msg.Txs {
		txs = append(txs, tx)
	}

	return nodecommon.ExtRollup{
		Header: FromRollupHeaderMsg(msg.Header),
		Txs:    txs,
	}
}

func FromRollupHeaderMsg(header *generated.HeaderMsg) *nodecommon.Header {
	if header == nil {
		return nil
	}
	withdrawals := make([]nodecommon.Withdrawal, 0)
	for _, withdrawalMsg := range header.Withdrawals {
		address := common.BytesToAddress(withdrawalMsg.Address)
		withdrawal := nodecommon.Withdrawal{Amount: withdrawalMsg.Amount, Address: address}
		withdrawals = append(withdrawals, withdrawal)
	}

	return &nodecommon.Header{
		ParentHash:  common.BytesToHash(header.ParentHash),
		Agg:         common.BytesToAddress(header.Agg),
		Nonce:       header.Nonce,
		L1Proof:     common.BytesToHash(header.L1Proof),
		State:       common.BytesToHash(header.StateRoot),
		Number:      header.Height,
		Withdrawals: withdrawals,
	}
}

func FromBlockHeaderMsg(msg *generated.BlockHeaderMsg) *types.Header {
	if msg == nil {
		return nil
	}
	return &types.Header{
		ParentHash:  common.BytesToHash(msg.ParentHash),
		UncleHash:   common.BytesToHash(msg.UncleHash),
		Coinbase:    common.BytesToAddress(msg.Coinbase),
		Root:        common.BytesToHash(msg.Root),
		TxHash:      common.BytesToHash(msg.TxHash),
		ReceiptHash: common.BytesToHash(msg.ReceiptHash),
		Bloom:       types.BytesToBloom(msg.Bloom),
		Difficulty:  big.NewInt(int64(msg.Difficulty)),
		Number:      big.NewInt(int64(msg.Number)),
		GasLimit:    msg.GasLimit,
		GasUsed:     msg.GasUsed,
		Time:        msg.Time,
		Extra:       msg.Extra,
		MixDigest:   common.BytesToHash(msg.MixDigest),
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
