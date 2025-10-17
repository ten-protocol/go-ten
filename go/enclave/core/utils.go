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

// validateChainID checks if the transaction has a valid (non-nil, non-zero) chain ID
func validateChainID(tx *types.Transaction) error {
	chainID := tx.ChainId()
	if chainID == nil || chainID.Int64() == 0 {
		return fmt.Errorf("transaction cannot have nil or zero chain ID")
	}
	return nil
}

// GetAuthenticatedSender - Get sender and tx nonce from transaction
func GetAuthenticatedSender(chainID int64, tx *types.Transaction) (*gethcommon.Address, error) {
	if err := validateChainID(tx); err != nil {
		return nil, err
	}
	if chainID != tx.ChainId().Int64() {
		return nil, fmt.Errorf("ten chain id does not match tx chain id. Expected %d vs %d", chainID, tx.ChainId().Int64())
	}
	signer := types.LatestSignerForChainID(tx.ChainId())
	sender, err := types.Sender(signer, tx)
	if err != nil {
		return nil, err
	}
	return &sender, nil
}

// GetExternalTxSigner returns the address that signed a transaction
func GetExternalTxSigner(tx *types.Transaction) (gethcommon.Address, error) {
	if err := validateChainID(tx); err != nil {
		return gethcommon.Address{}, fmt.Errorf("cannot get external tx signer: %w", err)
	}
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("could not recover sender for transaction. Cause: %w", err)
	}

	return from, nil
}

func GetTxSigner(tx *common.L2PricedTransaction) (gethcommon.Address, error) {
	if err := validateChainID(tx.Tx); err != nil {
		return gethcommon.Address{}, err
	}

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

type DurationThresholds struct {
	High    int64
	Medium  int64
	Low     int64
	Trivial int64
}

var (
	// default thresholds for quick operations
	DefaultThresholds = DurationThresholds{
		High:    500,
		Medium:  200,
		Low:     100,
		Trivial: 50,
	}

	// relaxed thresholds for known longer operations
	RelaxedThresholds = DurationThresholds{
		High:    4000,
		Medium:  1000,
		Low:     100,
		Trivial: 50,
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

	var label string
	switch {
	case durationMillis > thresholds.High:
		f = logger.Info
		label = "High"
	case durationMillis > thresholds.Medium:
		f = logger.Info
		label = "Medium"
	case durationMillis > thresholds.Low:
		f = logger.Debug
		label = "Low"
	case durationMillis > thresholds.Trivial:
		f = logger.Debug
		label = "Trivial"
	default:
		f = logger.Trace
		label = "Nothing"
	}
	newArgs := append([]any{log.DurationKey, stopWatch}, args...)
	f(fmt.Sprintf("Duration::%s::%s", label, msg), newArgs...)
}
