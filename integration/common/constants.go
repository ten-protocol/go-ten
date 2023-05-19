package common

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/go-obscuro/go/wallet"
)

// The Contract addresses are the result of the deploying a smart contract from hardcoded owners.
// The "owners" are keys which are the de-facto "admins" of those erc20s and are able to transfer or mint tokens.
// The contracts and addresses cannot be random for now, because there is hardcoded logic in the core
// to generate synthetic "transfer" transactions for each erc20 deposit on ethereum
// and these transactions need to be signed. Which means the platform needs to "own" ERC20s.

// ERC20 - the supported ERC20 tokens. A list of made-up tokens used for testing.
type ERC20 string

const (
	HOC            ERC20 = "HOC"
	POC            ERC20 = "POC"
	HOCAddr              = "f3a8bd422097bFdd9B3519Eaeb533393a1c561aC"
	pocAddr              = "9802F661d17c65527D7ABB59DAAD5439cb125a67"
	bridgeAddr           = "deB34A740ECa1eC42C8b8204CBEC0bA34FDD27f3"
	hocOwnerKeyHex       = "6e384a07a01263518a09a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682"
	pocOwnerKeyHex       = "4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8"
)

var HOCOwner, _ = crypto.HexToECDSA(hocOwnerKeyHex)

// HOCContract - address of the deployed "hocus" erc20 on the L2
var HOCContract = gethcommon.BytesToAddress(gethcommon.Hex2Bytes(HOCAddr))

var POCOwner, _ = crypto.HexToECDSA(pocOwnerKeyHex)

// POCContract - address of the deployed "pocus" erc20 on the L2
var POCContract = gethcommon.BytesToAddress(gethcommon.Hex2Bytes(pocAddr))

// BridgeAddress - address of the virtual bridge
var BridgeAddress = gethcommon.BytesToAddress(gethcommon.Hex2Bytes(bridgeAddr))

// ERC20Mapping - maps an L1 Erc20 to an L2 Erc20 address
type ERC20Mapping struct {
	Name ERC20

	// L1Owner   wallet.Wallet
	L1Address *gethcommon.Address

	Owner     wallet.Wallet // for now the wrapped L2 version is owned by a wallet, but this will change
	L2Address *gethcommon.Address
}
