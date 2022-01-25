package wallet_mock

import "github.com/google/uuid"

// todo - use proper crypto
type Address = uuid.UUID

type Wallet struct {
	Address Address
}
