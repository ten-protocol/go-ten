package wallet_mock

import (
	"github.com/google/uuid"
	"github.com/obscuronet/obscuro-playground/go/common"
)

type Wallet struct {
	Address common.Address
}

func New() Wallet {
	return Wallet{Address: uuid.New().ID()}
}
