package ethadapter

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/misc/eip4844"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
)

// SetTxGasPrice takes a txData type and overrides the From, Gas and Gas Price field with current values
// it bumps the price by a multiplier for retries. retryNumber is zero on first attempt (no multiplier on price)
func SetTxGasPrice(ctx context.Context, ethClient EthClient, txData types.TxData, from gethcommon.Address, nonce uint64, retryNumber int) (types.TxData, error) {
	rawTx := types.NewTx(txData)
	to := rawTx.To()
	value := rawTx.Value()
	data := rawTx.Data()

	// estimate gas
	estimatedGas, err := ethClient.EstimateGas(ctx, ethereum.CallMsg{
		From:  from,
		To:    to,
		Value: value,
		Data:  data,
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

	baseFee := head.BaseFee // Base fee per gas (EIP-1559)
	gasFeeCap := big.NewInt(0).Add(baseFee, gasTipCap)

	// Check if txData is of type *types.BlobTx
	if blobTx, ok := txData.(*types.BlobTx); ok {
		var blobBaseFee *big.Int
		if head.ExcessBlobGas != nil {
			blobBaseFee = eip4844.CalcBlobFee(*head.ExcessBlobGas)
		} else {
			return nil, fmt.Errorf("should not happen. missing blob base fee")
		}
		// blobFeeCap := calcBlobFeeCap(blobBaseFee, retryNumber)

		return &types.BlobTx{
			Nonce:      nonce,
			GasTipCap:  uint256.MustFromBig(gasTipCap),
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

// geth enforces a 1 gwei minimum for blob tx fee
var minBlobTxFee = big.NewInt(params.GWei)

// calcBlobFeeCap computes a suggested blob fee cap that is twice the current header's blob base fee
// value, with a minimum value of minBlobTxFee. It also doubles the blob fee cap based on the retry number.
func calcBlobFeeCap(blobBaseFee *big.Int, retryNumber int) *big.Int {
	// Base calculation: twice the current blob base fee
	// todo - why?
	blobFeeCap := new(big.Int).Mul(blobBaseFee, big.NewInt(2))

	// Ensure the blob fee cap is at least the minimum value
	if blobFeeCap.Cmp(minBlobTxFee) < 0 {
		blobFeeCap.Set(minBlobTxFee)
	}

	// Double the blob fee cap for each retry attempt
	if retryNumber > 0 {
		multiplier := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(retryNumber)), nil)
		blobFeeCap.Mul(blobFeeCap, multiplier)
	}

	return blobFeeCap
}
