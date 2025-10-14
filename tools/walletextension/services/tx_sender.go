package services

import (
	"context"
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	tencommonrpc "github.com/ten-protocol/go-ten/go/common/rpc"
	tenrpc "github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

// dustThresholdWei is the minimum balance considered worth recovering from an expired
// session key. Set in wei for precision.
// 1_000_000_000_000 wei = 1e12 wei = 1,000 gwei = 0.000001 ETH (~$0.004 at $4,000/ETH)
// Adjust the USD intuition based on current ETH price.
var dustThresholdWei = big.NewInt(1_000_000_000_000)

// TxSender encapsulates sending ETH value transactions signed by a session key
type TxSender interface {
	// SendAllMinusGasWithSK transfers entire balance minus gas from `from` to `to`
	SendAllMinusGasWithSK(ctx context.Context, user *wecommon.GWUser, from gethcommon.Address, to gethcommon.Address) (gethcommon.Hash, error)
}

type txSender struct {
	backend   *BackendRPC
	skManager SKManager
	logger    gethlog.Logger
}

func NewTxSender(backend *BackendRPC, skManager SKManager, logger gethlog.Logger) TxSender {
	return &txSender{backend: backend, skManager: skManager, logger: logger}
}

func (s *txSender) SendAllMinusGasWithSK(ctx context.Context, user *wecommon.GWUser, from gethcommon.Address, to gethcommon.Address) (gethcommon.Hash, error) {
	// Get balance at pending block so pending txs are reflected
	pending := rpc.PendingBlockNumber
	blockNrOrHash := rpc.BlockNumberOrHash{BlockNumber: &pending}

	var balance hexutil.Big
	if err := s.withSK(ctx, user, from, func(ctx context.Context, rpcClient *tenrpc.EncRPCClient) error {
		var b hexutil.Big
		if callErr := rpcClient.CallContext(ctx, &b, tencommonrpc.ERPCGetBalance, from, blockNrOrHash); callErr != nil {
			return callErr
		}
		balance = b
		return nil
	}); err != nil {
		return gethcommon.Hash{}, fmt.Errorf("failed to get balance via RPC: %w", err)
	}

	// Check if balance is below dust threshold - skip transfer if so
	if balance.ToInt().Cmp(dustThresholdWei) <= 0 {
		return gethcommon.Hash{}, nil
	}

	gasPrice, err := s.getGasPrice(ctx)
	if err != nil {
		return gethcommon.Hash{}, err
	}

	gasLimit, err := s.estimateGas(ctx, user, from, to, &balance)
	if err != nil {
		return gethcommon.Hash{}, err
	}

	gasCost := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	amountToSend := new(big.Int).Sub(balance.ToInt(), gasCost)
	if amountToSend.Sign() <= 0 {
		return gethcommon.Hash{}, fmt.Errorf("insufficient balance for gas: balance=%s gasCost=%s", balance.ToInt().String(), gasCost.String())
	}

	dynTx := &types.DynamicFeeTx{
		To:        &to,
		Value:     amountToSend,
		Gas:       gasLimit,
		GasTipCap: gasPrice,
		GasFeeCap: gasPrice,
	}

	tx := types.NewTx(dynTx)
	if tx == nil {
		return gethcommon.Hash{}, fmt.Errorf("failed to create transaction")
	}

	signedTx, err := s.skManager.SignTx(ctx, user, from, tx)
	if err != nil {
		s.logger.Error("Failed to sign transaction with session key", "error", err, "sessionKeyAddress", from.Hex())
		return gethcommon.Hash{}, fmt.Errorf("failed to sign transaction with session key: %w", err)
	}
	blob, err := signedTx.MarshalBinary()
	if err != nil {
		s.logger.Error("Failed to marshal signed transaction", "error", err)
		return gethcommon.Hash{}, fmt.Errorf("failed to marshal signed transaction: %w", err)
	}

	txHash, err := s.sendRawTransaction(ctx, user, from, blob)
	if err != nil {
		s.logger.Error("Failed to send transaction", "error", err, "sessionKeyAddress", from.Hex())
		return gethcommon.Hash{}, fmt.Errorf("failed to send transaction: %w", err)
	}
	return txHash, nil
}

// withSK opens an encrypted RPC connection authorized by the session key at `addr`
func (s *txSender) withSK(
	ctx context.Context,
	user *wecommon.GWUser,
	addr gethcommon.Address,
	fn func(ctx context.Context, c *tenrpc.EncRPCClient) error,
) error {
	sk, ok := user.SessionKeys[addr]
	if !ok {
		return fmt.Errorf("session key not found for address %s", addr.Hex())
	}
	_, err := WithEncRPCConnection(ctx, s.backend, sk.Account, func(c *tenrpc.EncRPCClient) (*struct{}, error) {
		return &struct{}{}, fn(ctx, c)
	})
	return err
}

func (s *txSender) getGasPrice(ctx context.Context) (*big.Int, error) {
	var result hexutil.Big
	_, err := WithPlainRPCConnection(ctx, s.backend, func(client *rpc.Client) (*hexutil.Big, error) {
		err := client.CallContext(ctx, &result, tenrpc.GasPrice)
		return &result, err
	})
	if err != nil {
		s.logger.Error("Failed to get gas price", "error", err)
		return nil, fmt.Errorf("failed to get gas price via RPC: %w", err)
	}
	return result.ToInt(), nil
}

func (s *txSender) estimateGas(ctx context.Context, user *wecommon.GWUser, from, to gethcommon.Address, value *hexutil.Big) (uint64, error) {
	var result hexutil.Uint64
	err := s.withSK(ctx, user, from, func(ctx context.Context, rpcClient *tenrpc.EncRPCClient) error {
		params := map[string]interface{}{
			"from":  from.Hex(),
			"to":    to.Hex(),
			"value": value.String(),
		}
		return rpcClient.CallContext(ctx, &result, tencommonrpc.ERPCEstimateGas, params)
	})
	if err != nil {
		s.logger.Error("Failed to estimate gas", "error", err, "from", from.Hex(), "to", to.Hex(), "value", value.String())
		return 0, fmt.Errorf("failed to estimate gas via RPC: %w", err)
	}
	return uint64(result), nil
}

func (s *txSender) sendRawTransaction(ctx context.Context, user *wecommon.GWUser, sessionKeyAddr gethcommon.Address, input hexutil.Bytes) (gethcommon.Hash, error) {
	var result gethcommon.Hash
	err := s.withSK(ctx, user, sessionKeyAddr, func(ctx context.Context, rpcClient *tenrpc.EncRPCClient) error {
		return rpcClient.CallContext(ctx, &result, tencommonrpc.ERPCSendRawTransaction, input)
	})
	if err != nil {
		return gethcommon.Hash{}, err
	}
	return result, nil
}
