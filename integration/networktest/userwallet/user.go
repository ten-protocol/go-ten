package userwallet

import (
	"context"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/wallet"
)

// User - abstraction for networktest users -  two implementations initially:
// 1. AuthClientUser - a user that uses the auth client to talk to the network
// 2. GatewayUser - a user that uses the gateway to talk to the network
//
// This abstraction allows us to use the same tests for both types of users
type User interface {
	Wallet() wallet.Wallet
	SendFunds(ctx context.Context, addr gethcommon.Address, value *big.Int) (*gethcommon.Hash, error)
	AwaitReceipt(ctx context.Context, txHash *gethcommon.Hash) (*types.Receipt, error)
	NativeBalance(ctx context.Context) (*big.Int, error)
	GetPersonalTransactions(ctx context.Context, pagination common.QueryPagination) (types.Receipts, uint64, error)
}
