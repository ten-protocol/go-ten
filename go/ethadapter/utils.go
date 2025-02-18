package ethadapter

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/params"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/misc/eip4844"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
)

const (
	_retryPriceMultiplier     = 1.3 // over five attempts will give multipliers of 1.3, 1.7, 2.2, 2.8, 3.7
	_blobPriceMultiplier      = 2.0
	_maxTxRetryPriceIncreases = 5
)

// SetTxGasPrice takes a txData type and overrides the From, Gas and Gas Price field with current values
// it bumps the price by a multiplier for retries. retryNumber is zero on first attempt (no multiplier on price)
func SetTxGasPrice(ctx context.Context, ethClient EthClient, txData types.TxData, from gethcommon.Address, nonce uint64, retryNumber int, logger gethlog.Logger) (types.TxData, error) {
	rawTx := types.NewTx(txData)
	to := rawTx.To()
	value := rawTx.Value()
	data := rawTx.Data()
	blobHashes := rawTx.BlobHashes()

	// estimate gas
	estimatedGas, err := ethClient.EstimateGas(ctx, ethereum.CallMsg{
		From:       from,
		To:         to,
		Value:      value,
		Data:       data,
		BlobHashes: blobHashes,
	})
	if err != nil {
		return nil, fmt.Errorf("could not estimate gas - %w", err)
	}

	// calculate the gas tip
	gasTipCap, err := ethClient.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not suggest gas price - %w", err)
	}

	// adjust the gasTipCap if we have to retry
	// it should never happen but to avoid any risk of repeated price increases we cap the possible retry price bumps to 5
	// we apply a 30% gas price increase for each retry (retrying with similar price gets rejected by mempool)
	// Retry '0' is the first attempt, gives multiplier of 1.0

	retryMultiplier := calculateRetryMultiplier(_retryPriceMultiplier, retryNumber)
	gasTipCap = big.NewInt(0).SetUint64(uint64(retryMultiplier * float64(gasTipCap.Uint64())))

	// calculate the gas fee cap
	head, err := ethClient.HeaderByNumber(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get the latest block header: %w", err)
	}

	baseFee := head.BaseFee
	gasFeeCap := big.NewInt(0).Add(baseFee, gasTipCap)

	if blobTx, ok := txData.(*types.BlobTx); ok {
		if head.ExcessBlobGas == nil {
			return nil, fmt.Errorf("should not happen. missing blob base fee")
		}
		blobBaseFee := eip4844.CalcBlobFee(*head.ExcessBlobGas)
		blobMultiplier := calculateRetryMultiplier(_blobPriceMultiplier, retryNumber)
		blobFeeCap := new(uint256.Int).Mul(
			uint256.MustFromBig(blobBaseFee),
			uint256.NewInt(uint64(math.Ceil(blobMultiplier)))) // double base fee with retry multiplier,

		// even if we hit the minimum, we should still increase for retries
		if blobFeeCap.Lt(uint256.NewInt(params.GWei)) {
			blobFeeCap = new(uint256.Int).Mul(
				uint256.NewInt(params.GWei),
				uint256.NewInt(uint64(math.Ceil(blobMultiplier))),
			)
		}

		logger.Info("Sending blob tx with gas price", "retry", retryNumber, "nonce", nonce, "blobFeeCap",
			blobFeeCap, "gasTipCap", gasTipCap, "gasFeeCap", gasFeeCap, "estimatedGas", estimatedGas, "to", to)

		return &types.BlobTx{
			Nonce:      nonce,
			GasTipCap:  uint256.MustFromBig(gasTipCap),
			GasFeeCap:  uint256.MustFromBig(gasFeeCap),
			Gas:        estimatedGas,
			To:         *to,
			Value:      uint256.MustFromBig(value),
			Data:       data,
			BlobFeeCap: blobFeeCap,
			BlobHashes: blobTx.BlobHashes,
			Sidecar:    blobTx.Sidecar,
		}, nil
	}

	logger.Info("Sending tx with gas price", "retry", retryNumber, "nonce", nonce, "gasTipCap", gasTipCap, "gasFeeCap", gasFeeCap, "estimatedGas", estimatedGas, "to", to)

	// For non-blob transactions, just use the latest suggested values with multiplier
	return &types.DynamicFeeTx{
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       estimatedGas,
		To:        to,
		Value:     value,
		Data:      data,
	}, nil
}

func calculateRetryMultiplier(baseMultiplier float64, retryNumber int) float64 {
	return math.Pow(baseMultiplier, float64(min(_maxTxRetryPriceIncreases, retryNumber)))
}
