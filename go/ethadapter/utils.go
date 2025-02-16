package ethadapter

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/misc/eip4844"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
)

var _retryPriceMultiplier = 1.3 // over five attempts will give multipliers of 1.3, 1.7, 2.2, 2.8, 3.7

// SetTxGasPrice takes a txData type and overrides the From, Gas and Gas Price field with current values
// it bumps the price by a multiplier for retries. retryNumber is zero on first attempt (no multiplier on price)
func SetTxGasPrice(ctx context.Context, ethClient EthClient, txData types.TxData, from gethcommon.Address, nonce uint64, retryNumber int) (types.TxData, error) {
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

	/*
		// there is no need to adjust the gasTipCap if we have to retry because we are already using the latest suggested value
		// adjust the gasTipCap if we have to retry
		// it should never happen but to avoid any risk of repeated price increases we cap the possible retry price bumps to 5
		// we apply a 100% gas price increase for each retry (retrying with similar price gets rejected by mempool)
		// Retry '0' is the first attempt, gives multiplier of 1.0
		retryMultiplier := math.Pow(_retryPriceMultiplier, float64(min(_maxTxRetryPriceIncreases, retryNumber)))
		gasTipCap = big.NewInt(0).Mul(gasTipCap, big.NewInt(int64(retryMultiplier)))
	*/

	// calculate the gas fee cap
	head, err := ethClient.HeaderByNumber(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get the latest block header: %w", err)
	}

	baseFee := head.BaseFee
	gasFeeCap := big.NewInt(0).Add(baseFee, gasTipCap)

	if blobTx, ok := txData.(*types.BlobTx); ok {
		var blobBaseFee *big.Int
		if head.ExcessBlobGas != nil {
			blobBaseFee = eip4844.CalcBlobFee(*head.ExcessBlobGas)
		} else {
			return nil, fmt.Errorf("should not happen. missing blob base fee")
		}

		// for blob transactions, increase the tip by 30% on each retry
		if retryNumber > 0 {
			multiplier := int64(math.Pow(_retryPriceMultiplier, float64(retryNumber)))
			gasTipCap = new(big.Int).Mul(gasTipCap, big.NewInt(multiplier))
		}
		gasFeeCapUpdated := new(big.Int).Add(head.BaseFee, gasTipCap)

		return &types.BlobTx{
			Nonce:      nonce,
			GasTipCap:  uint256.MustFromBig(gasFeeCapUpdated),
			GasFeeCap:  uint256.MustFromBig(gasFeeCap),
			Gas:        estimatedGas,
			To:         *to,
			Value:      uint256.MustFromBig(value),
			Data:       data,
			BlobFeeCap: uint256.MustFromBig(blobBaseFee),
			BlobHashes: blobTx.BlobHashes,
			Sidecar:    blobTx.Sidecar,
		}, nil
	}

	// For non-blob transactions, just use the latest suggested values without multiplier
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
