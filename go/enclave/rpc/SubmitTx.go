package rpc

import (
	"errors"
	"fmt"

	"github.com/ten-protocol/go-ten/go/enclave/core"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/responses"
)

func (rpc *EncryptionManager) SubmitTx(encryptedTxParams common.EncryptedTx) (*responses.RawTx, common.SystemError) {
	return withVKEncryption1[common.L2Tx](
		rpc,
		rpc.config.ObscuroChainID,
		encryptedTxParams,
		// extract sender and arguments
		func(reqParams []any) (*UserRPCRequest1[common.L2Tx], error) {
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
		},
		// make call and return result
		func(decodedParams *UserRPCRequest1[common.L2Tx]) (any, error, error) {
			if rpc.processors.Local.IsSyntheticTransaction(*decodedParams.Param1) {
				return nil, fmt.Errorf("synthetic transaction coming from external rpc"), nil
			}

			if err := rpc.service.SubmitTransaction(decodedParams.Param1); err != nil {
				rpc.logger.Debug("Could not submit transaction", log.TxKey, decodedParams.Param1.Hash(), log.ErrKey, err)
				return nil, err, nil
			}
			return decodedParams.Param1.Hash(), nil, nil
		},
	)
}
