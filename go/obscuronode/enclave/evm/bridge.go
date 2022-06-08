package evm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// These are hardcoded values necessary as an intermediary step.
// The assumption is that there is a single ERC20 which represents "The balance"
// Todo - this has to be changed to mapping of "supported ERC20 Ethereum address - Obscuro address" ( eg.: USDT address -> Obscuro WUSDT address)
// Todo - also on depositing, there has to be a minting step
var (
	Erc20OwnerKey, _     = crypto.HexToECDSA("6e384a07a01263518a09a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682")
	Erc20OwnerAddress    = crypto.PubkeyToAddress(Erc20OwnerKey.PublicKey)
	Erc20ContractAddress = common.BytesToAddress(common.Hex2Bytes("f3a8bd422097bFdd9B3519Eaeb533393a1c561aC"))
)

// WithdrawalAddress Custom address used for exiting Obscuro
// Todo - This should be the address of a Bridge contract.
var WithdrawalAddress = common.HexToAddress("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
