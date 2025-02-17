package core

import (
	"fmt"
	"math/big"

	"github.com/ten-protocol/go-ten/go/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/measure"
)

// VerifySignature - Checks that the L2Tx has a valid signature.
func VerifySignature(chainID int64, tx *types.Transaction) error {
	signer := types.LatestSignerForChainID(big.NewInt(chainID))
	_, err := types.Sender(signer, tx)
	return err
}

// GetAuthenticatedSender - Get sender and tx nonce from transaction
func GetAuthenticatedSender(chainID int64, tx *types.Transaction) (*gethcommon.Address, error) {
	signer := types.LatestSignerForChainID(tx.ChainId())
	sender, err := types.Sender(signer, tx)
	if err != nil {
		return nil, err
	}
	return &sender, nil
}

type DurationThresholds struct {
	Error int64
	Warn  int64
	Info  int64
	Debug int64
}

var (
	// default thresholds for quick operations
	DefaultThresholds = DurationThresholds{
		Error: 500,
		Warn:  200,
		Info:  100,
		Debug: 50,
	}

	// relaxed thresholds for known longer operations
	RelaxedThresholds = DurationThresholds{
		Error: 4000,
		Warn:  1000,
		Info:  100,
		Debug: 50,
	}
)

// LogMethodDuration - call only with "defer"
func LogMethodDuration(logger gethlog.Logger, stopWatch *measure.Stopwatch, msg string, args ...any) {
	var thresholds *DurationThresholds
	if len(args) > 0 {
		if t, ok := args[0].(*DurationThresholds); ok {
			thresholds = t
			args = args[1:] // remove thresholds from args if present
		}
	}
	if thresholds == nil {
		thresholds = &DefaultThresholds
	}

	var f func(msg string, ctx ...interface{})
	durationMillis := stopWatch.Measure().Milliseconds()

	switch {
	case durationMillis > thresholds.Error:
		f = logger.Error
	case durationMillis > thresholds.Warn:
		f = logger.Warn
	case durationMillis > thresholds.Info:
		f = logger.Info
	case durationMillis > thresholds.Debug:
		f = logger.Debug
	default:
		f = logger.Trace
	}
	newArgs := append([]any{log.DurationKey, stopWatch}, args...)
	f(fmt.Sprintf("LogMethodDuration::%s", msg), newArgs...)
}

// GetExternalTxSigner returns the address that signed a transaction
func GetExternalTxSigner(tx *types.Transaction) (gethcommon.Address, error) {
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("could not recover sender for transaction. Cause: %w", err)
	}

	return from, nil
}

func GetTxSigner(tx *common.L2PricedTransaction) (gethcommon.Address, error) {
	if tx.SystemDeployer {
		return common.MaskedSender(gethcommon.BigToAddress(big.NewInt(tx.Tx.ChainId().Int64()))), nil
	} else if tx.FromSelf {
		return common.MaskedSender(*tx.Tx.To()), nil
	}

	from, err := types.Sender(types.LatestSignerForChainID(tx.Tx.ChainId()), tx.Tx)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("could not recover sender for transaction. Cause: %w", err)
	}

	return from, nil
}
