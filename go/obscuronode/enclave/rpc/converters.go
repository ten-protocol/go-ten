package rpc

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/rpc/generated"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// Functions to convert classes that need to be sent between the host and the enclave to and from their equivalent
// Protobuf message classes.

func toAttestationReportMsg(report obscurocommon.AttestationReport) generated.AttestationReportMsg {
	return generated.AttestationReportMsg{Owner: report.Owner.Bytes()}
}

func fromAttestationReportMsg(msg *generated.AttestationReportMsg) obscurocommon.AttestationReport {
	return obscurocommon.AttestationReport{Owner: common.BytesToAddress(msg.Owner)}
}

func toBlockSubmissionResponseMsg(response enclave.BlockSubmissionResponse) generated.BlockSubmissionResponseMsg {
	withdrawalMsgs := make([]*generated.WithdrawalMsg, 0)
	for _, withdrawal := range response.Withdrawals {
		withdrawalMsg := generated.WithdrawalMsg{Amount: withdrawal.Amount, Address: withdrawal.Address.Bytes()}
		withdrawalMsgs = append(withdrawalMsgs, &withdrawalMsg)
	}

	producedRollupMsg := toExtRollupMsg(&response.ProducedRollup)

	return generated.BlockSubmissionResponseMsg{
		L1Hash:            response.L1Hash.Bytes(),
		L1Height:          response.L1Height,
		L1Parent:          response.L1Parent.Bytes(),
		L2Hash:            response.L2Hash.Bytes(),
		L2Height:          response.L2Height,
		L2Parent:          response.L2Parent.Bytes(),
		Withdrawals:       withdrawalMsgs,
		ProducedRollup:    &producedRollupMsg,
		IngestedBlock:     response.IngestedBlock,
		IngestedNewRollup: response.IngestedNewRollup,
	}
}

func fromBlockSubmissionResponseMsg(msg *generated.BlockSubmissionResponseMsg) enclave.BlockSubmissionResponse {
	withdrawals := make([]nodecommon.Withdrawal, 0)
	for _, withdrawalMsg := range msg.Withdrawals {
		address := common.BytesToAddress(withdrawalMsg.Address)
		withdrawal := nodecommon.Withdrawal{Amount: withdrawalMsg.Amount, Address: address}
		withdrawals = append(withdrawals, withdrawal)
	}

	return enclave.BlockSubmissionResponse{
		L1Hash:            common.BytesToHash(msg.L1Hash),
		L1Height:          msg.L1Height,
		L1Parent:          common.BytesToHash(msg.L1Parent),
		L2Hash:            common.BytesToHash(msg.L2Hash),
		L2Height:          msg.L2Height,
		L2Parent:          common.BytesToHash(msg.L2Parent),
		Withdrawals:       withdrawals,
		ProducedRollup:    fromExtRollupMsg(msg.ProducedRollup),
		IngestedBlock:     msg.IngestedBlock,
		IngestedNewRollup: msg.IngestedNewRollup,
	}
}

func toExtRollupMsg(rollup *nodecommon.ExtRollup) generated.ExtRollupMsg {
	var headerMsg generated.HeaderMsg
	if rollup.Header != nil {
		withdrawalMsgs := make([]*generated.WithdrawalMsg, 0)
		for _, withdrawal := range rollup.Header.Withdrawals {
			withdrawalMsg := generated.WithdrawalMsg{Amount: withdrawal.Amount, Address: withdrawal.Address.Bytes()}
			withdrawalMsgs = append(withdrawalMsgs, &withdrawalMsg)
		}

		headerMsg = generated.HeaderMsg{
			ParentHash:  rollup.Header.ParentHash.Bytes(),
			Agg:         rollup.Header.Agg.Bytes(),
			Nonce:       rollup.Header.Nonce,
			L1Proof:     rollup.Header.L1Proof.Bytes(),
			StateRoot:   rollup.Header.State,
			Withdrawals: withdrawalMsgs,
		}
	}

	txs := make([][]byte, 0)
	for _, tx := range rollup.Txs {
		txs = append(txs, tx)
	}

	return generated.ExtRollupMsg{Header: &headerMsg, Txs: txs}
}

func fromExtRollupMsg(msg *generated.ExtRollupMsg) nodecommon.ExtRollup {
	withdrawals := make([]nodecommon.Withdrawal, 0)
	for _, withdrawalMsg := range msg.Header.Withdrawals {
		address := common.BytesToAddress(withdrawalMsg.Address)
		withdrawal := nodecommon.Withdrawal{Amount: withdrawalMsg.Amount, Address: address}
		withdrawals = append(withdrawals, withdrawal)
	}

	header := nodecommon.Header{
		ParentHash:  common.BytesToHash(msg.Header.ParentHash),
		Agg:         common.BytesToAddress(msg.Header.Agg),
		Nonce:       msg.Header.Nonce,
		L1Proof:     common.BytesToHash(msg.Header.L1Proof),
		State:       msg.Header.StateRoot,
		Withdrawals: withdrawals,
	}

	txs := make([]nodecommon.EncryptedTx, 0)
	for _, tx := range msg.Txs {
		txs = append(txs, tx)
	}

	return nodecommon.ExtRollup{
		Header: &header,
		Txs:    txs,
	}
}
