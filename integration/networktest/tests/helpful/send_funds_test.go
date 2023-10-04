package helpful

import (
	"context"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/networktest"
	"github.com/obscuronet/go-obscuro/integration/networktest/actions"
	"github.com/obscuronet/go-obscuro/integration/networktest/env"
	"math/big"
	"testing"
)

var _amtToSend = new(big.Int).Mul(big.NewInt(31), big.NewInt(1e18))
var _fromAccPK = ""
var _toAcc = gethcommon.HexToAddress("0x563EAc5dfDFebA3C53c2160Bf1Bd62941E3D0005")

func TestSendFunds(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"L1 funds xfer",
		t,
		env.SepoliaTestnet(),
		actions.RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
			pk, err := crypto.HexToECDSA(_fromAccPK)
			if err != nil {
				return nil, err
			}
			wal := wallet.NewInMemoryWalletFromPK(params.SepoliaChainConfig.ChainID, pk, nil)

			client, err := network.GetL1Client()
			if err != nil {
				return nil, err
			}
			nonce, err := client.Nonce(wal.Address())
			if err != nil {
				return nil, err
			}
			txData := &types.LegacyTx{
				Value: _amtToSend,
				To:    &_toAcc,
			}
			tx, err := client.PrepareTransactionToSend(txData, wal.Address(), nonce)
			if err != nil {
				return nil, err
			}
			signedTx, err := wal.SignTransaction(tx)
			if err != nil {
				return nil, err
			}
			err = client.SendTransaction(signedTx)
			if err != nil {
				return nil, err
			}

			return ctx, nil
		}),
	)
}
