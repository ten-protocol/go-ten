package evm

import (
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethcore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/ten-protocol/go-ten/go/common"
)

func adjustPublishingCostGas(tx *common.L2PricedTransaction, msg *gethcore.Message, s *state.StateDB, header *types.Header, noBaseFee bool, execute func() (receipt *types.Receipt, err error)) error {
	l1cost := tx.PublishingCost
	l1Gas := big.NewInt(0)
	hasL1Cost := l1cost.Cmp(big.NewInt(0)) != 0

	// If a transaction has to be published on the l1, it will have an l1 cost
	if hasL1Cost {
		l1Gas.Div(l1cost, header.BaseFee) // TotalCost/CostPerGas = Gas
		l1Gas.Add(l1Gas, big.NewInt(1))   // Cover from leftover from the division

		// The gas limit of the transaction (evm message) should always be higher than the gas overhead
		// used to cover the l1 cost
		// todo - this check has to be added to the mempool as well
		if msg.GasLimit < l1Gas.Uint64() {
			return fmt.Errorf("%w. Want at least: %d have: %d", ErrGasNotEnoughForL1, l1Gas, msg.GasLimit)
		}

		// Remove the gas overhead for l1 publishing from the gas limit in order to define
		// the actual gas limit for execution
		msg.GasLimit -= l1Gas.Uint64()

		// Remove the l1 cost from the sender
		// and pay it to the coinbase of the batch
		s.SubBalance(msg.From, uint256.MustFromBig(l1cost), BalanceDecreaseL1Payment)
		s.AddBalance(header.Coinbase, uint256.MustFromBig(l1cost), BalanceIncreaseL1Payment)
	}

	receipt, err := execute()
	if err != nil {
		// If the transaction has l1 cost, then revert the funds exchange
		// as it will not be published on error (no receipt condition)
		if hasL1Cost {
			s.SubBalance(header.Coinbase, uint256.MustFromBig(l1cost), BalanceRevertIncreaseL1Payment)
			s.AddBalance(msg.From, uint256.MustFromBig(l1cost), BalanceRevertDecreaseL1Payment)
		}
		return err
	}

	// Synthetic transactions and ten zen are free. Do not increase the balance of the coinbase.
	isPaidProcessing := !noBaseFee

	// Do not increase the balance of zero address as it is the contract deployment address.
	// Doing so might cause weird interactions.
	if header.Coinbase.Big().Cmp(gethcommon.Big0) != 0 && isPaidProcessing {
		gasUsed := big.NewInt(0).SetUint64(receipt.GasUsed)
		executionGasCost := big.NewInt(0).Mul(gasUsed, header.BaseFee)
		// As the baseFee is burned, we add it back to the coinbase.
		// Geth should automatically add the tips.
		s.AddBalance(header.Coinbase, uint256.MustFromBig(executionGasCost), tracing.BalanceDecreaseGasBuy)
	}

	receipt.GasUsed += l1Gas.Uint64()

	return nil
}
