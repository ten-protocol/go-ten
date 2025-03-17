package ethadapter

import (
	"context"
	"encoding/base64"
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
	_retryPriceMultiplier     = 1.7
	_blobPriceMultiplier      = 2.0
	_maxTxRetryPriceIncreases = 5
	_baseFeeIncreaseFactor    = 2 // increase the base fee by 20%
)

var minGasTipCap = big.NewInt(2 * params.GWei)

// SetTxGasPrice takes a txData type and overrides the From, Gas and Gas Price field with current values
// it bumps the price by a multiplier for retries. retryNumber is zero on first attempt (no multiplier on price)
func SetTxGasPrice(ctx context.Context, ethClient EthClient, txData types.TxData, from gethcommon.Address, nonce uint64, retryNumber int, l1ChainCfg *params.ChainConfig, logger gethlog.Logger) (types.TxData, error) {
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

	// make sure the gas tip is always greater than the minimum
	if gasTipCap.Cmp(minGasTipCap) < 0 {
		gasTipCap = minGasTipCap
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

	// increase the baseFee by a factor
	baseFee := big.NewInt(0).SetInt64(int64(float64(head.BaseFee.Uint64()) * _baseFeeIncreaseFactor))
	gasFeeCap := big.NewInt(0).Add(baseFee, gasTipCap)

	if blobTx, ok := txData.(*types.BlobTx); ok {
		if head.ExcessBlobGas == nil {
			return nil, fmt.Errorf("should not happen. missing blob base fee")
		}
		blobBaseFee := eip4844.CalcBlobFee(l1ChainCfg, head)
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

// Base64EncodeToString encodes a byte array to a string
func Base64EncodeToString(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

// Base64DecodeFromString decodes a string to a byte array
func Base64DecodeFromString(in string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(in)
}
