package rpc

import (
	"errors"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
)

func SubmitTxValidate(reqParams []any, builder *CallBuilder[common.L2Tx, gethcommon.Hash], _ *EncryptionManager) error {
	txStr, ok := reqParams[0].(string)
	if !ok {
		return errors.New("unsupported format")
	}
	l2Tx, err := ExtractTx(txStr)
	if err != nil {
		builder.Err = fmt.Errorf("could not extract transaction. Cause: %w", err)
		return nil
	}
	builder.Param = l2Tx
	return nil
}

func SubmitTxExecute(builder *CallBuilder[common.L2Tx, gethcommon.Hash], rpc *EncryptionManager) error {
	if rpc.processors.Local.IsSyntheticTransaction(builder.Param) {
		builder.Err = fmt.Errorf("synthetic transaction coming from external rpc")
		return nil
	}

	if err := rpc.service.SubmitTransaction(builder.Param); err != nil {
		rpc.logger.Debug("Could not submit transaction", log.TxKey, builder.Param.Hash(), log.ErrKey, err)
		builder.Err = err
		return nil
	}
	h := builder.Param.Hash()
	fmt.Printf("Submit TX execute hash: %s\n", h.Hex())
	builder.ReturnValue = &h
	return nil
}
