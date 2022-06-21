package evm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
)

// these addresses are the result of the deploying a smart contract from the hardcoded owners
// Todo - remove these in a next iteration
var WBtcContract = common.BytesToAddress(common.Hex2Bytes("f3a8bd422097bFdd9B3519Eaeb533393a1c561aC"))
var WEthContract = common.BytesToAddress(common.Hex2Bytes("9802F661d17c65527D7ABB59DAAD5439cb125a67"))

var WithdrawalAddress = common.BytesToAddress(common.Hex2Bytes("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"))

type ERC20 int

// the supported ERC20 tokens
const (
	BTC ERC20 = iota
	ETH
)

type Token struct {
	Name ERC20

	// L1Owner   wallet.Wallet
	L1Address *common.Address

	Owner     wallet.Wallet // for now the wrapped L2 version is owned by a wallet, but this will change
	L2Address *common.Address
}

type Bridge struct {
	// These are hardcoded values necessary as an intermediary step.
	SupportedTokens   map[ERC20]*Token
	WithdrawalAddress common.Address
}

func NewBridge(obscuroChainID int64, btcAddress *common.Address, ethAddress *common.Address) *Bridge {
	tokens := make(map[ERC20]*Token, 0)

	// create the ethereum wallets to be used to deploy ERC20 contracts
	// and their counterparts in the Obscuro world for the Owner versions
	// this cannot be random for now, because there is hardcoded logic in the obscuro core
	// to generate synthetic "transfer" transactions on the Owner erc20 for each erc20 deposit on ethereum
	// and these transactions need to be signed
	btc, _ := crypto.HexToECDSA("6e384a07a01263518a09a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682")
	tokens[BTC] = &Token{
		Name:      BTC,
		L1Address: btcAddress,
		Owner:     wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), btc),
		L2Address: &WBtcContract,
	}

	eth, _ := crypto.HexToECDSA("4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8")
	tokens[ETH] = &Token{
		Name:      ETH,
		L1Address: ethAddress,
		Owner:     wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), eth),
		L2Address: &WEthContract,
	}

	return &Bridge{
		SupportedTokens:   tokens,
		WithdrawalAddress: WithdrawalAddress,
	}
}

func (b *Bridge) IsWithdrawal(address common.Address) bool {
	return address == b.WithdrawalAddress
}

// L1Address - returns the L1 address of a token based on the mapping
func (b *Bridge) L1Address(l2Address *common.Address) *common.Address {
	for _, t := range b.SupportedTokens {
		if l2Address == t.L2Address {
			return t.L1Address
		}
	}
	return nil
}

func (b *Bridge) Token(l1ContractAddress *common.Address) *Token {
	for _, t := range b.SupportedTokens {
		if t.L1Address == l1ContractAddress {
			return t
		}
	}
	return nil
}
