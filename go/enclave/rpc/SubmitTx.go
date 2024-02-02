package rpc

import (
	"errors"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/enclave/core"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
)

// todo - do you really need authentication?
func ExtractSubmitTxRequest(reqParams []any, _ *EncryptionManager) (*UserRPCRequest1[common.L2Tx], error) {
	l2Tx, err := ExtractTx(reqParams[0].(string))
	if err != nil {
		return nil, fmt.Errorf("could not extract transaction. Cause: %w", err)
	}
	sender, err := core.GetSender(l2Tx)
	if err != nil {
		if errors.Is(err, types.ErrInvalidSig) {
			return nil, fmt.Errorf("invalid signature")
		}
		return nil, fmt.Errorf("could not recover from address. Cause: %w", err)
	}
	return &UserRPCRequest1[common.L2Tx]{&sender, l2Tx}, nil
}

func ExecuteSubmitTx(decodedParams *UserRPCRequest1[common.L2Tx], rpc *EncryptionManager) (*UserResponse[gethcommon.Hash], error) {
	if rpc.processors.Local.IsSyntheticTransaction(*decodedParams.Param1) {
		return &UserResponse[gethcommon.Hash]{nil, fmt.Errorf("synthetic transaction coming from external rpc")}, nil
	}

	if err := rpc.service.SubmitTransaction(decodedParams.Param1); err != nil {
		rpc.logger.Debug("Could not submit transaction", log.TxKey, decodedParams.Param1.Hash(), log.ErrKey, err)
		return &UserResponse[gethcommon.Hash]{nil, err}, nil
	}
	h := decodedParams.Param1.Hash()
	return &UserResponse[gethcommon.Hash]{&h, nil}, nil
}
