package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/evm"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
)

type SimToken struct {
	Name              evm.ERC20
	L1Owner           wallet.Wallet
	L1ContractAddress *common.Address

	L2Owner           wallet.Wallet
	L2ContractAddress *common.Address
}

type SimWallets struct {
	MCOwnerWallet wallet.Wallet   // owner of the management contract deployed on Ethereum
	NodeWallets   []wallet.Wallet // the keys used by the obscuro nodes to submit rollups to Eth

	SimEthWallets []wallet.Wallet // the wallets of the simulated users on the Ethereum side
	SimObsWallets []wallet.Wallet // and their equivalents on the obscuro side (with a different chainId)

	Tokens map[evm.ERC20]*SimToken // The supported tokens
}

func NewSimWallets(nrSimWallets int, nNodes int, ethereumChainID int64, obscuroChainID int64) *SimWallets {
	// create the ethereum wallets to be used by the nodes
	nodeWallets := make([]wallet.Wallet, nNodes)
	for i := 0; i < nNodes; i++ {
		nodeWallets[i] = datagenerator.RandomWallet(ethereumChainID)
	}

	// create the wallets to be used by the simulated users
	// they will use the same key on both Ethereum and Obscuro, but different chainIDs
	simEthWallets := make([]wallet.Wallet, nrSimWallets)
	simObsWallets := make([]wallet.Wallet, nrSimWallets)
	for i := 0; i < nrSimWallets; i++ {
		simEthWallets[i] = datagenerator.RandomWallet(ethereumChainID)
		simObsWallets[i] = wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), simEthWallets[i].PrivateKey())
	}

	// create the wallet to deploy the Management contract
	mcOwnerWallet := datagenerator.RandomWallet(ethereumChainID)

	// create the L1 addresses of the two tokens, and connect them to the hardcoded addresses from the enclave
	btc := SimToken{
		Name:              evm.BTC,
		L1Owner:           datagenerator.RandomWallet(ethereumChainID),
		L2Owner:           wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), evm.WBtcOwner),
		L2ContractAddress: &evm.WBtcContract,
	}
	eth := SimToken{
		Name:              evm.ETH,
		L1Owner:           datagenerator.RandomWallet(ethereumChainID),
		L2Owner:           wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), evm.WEthOnwer),
		L2ContractAddress: &evm.WEthContract,
	}

	return &SimWallets{
		MCOwnerWallet: mcOwnerWallet,
		NodeWallets:   nodeWallets,
		SimEthWallets: simEthWallets,
		SimObsWallets: simObsWallets,
		Tokens: map[evm.ERC20]*SimToken{
			evm.BTC: &btc,
			evm.ETH: &eth,
		},
	}
}

func (w *SimWallets) AllEthWallets() []wallet.Wallet {
	ethWallets := make([]wallet.Wallet, 0)
	for _, token := range w.Tokens {
		ethWallets = append(ethWallets, token.L1Owner)
	}
	return append(append(append(w.NodeWallets, w.SimEthWallets...), w.MCOwnerWallet), ethWallets...)
}

func (w *SimWallets) AllEthAddresses() []*common.Address {
	addresses := make([]*common.Address, 0)
	addresses = append(addresses, w.Tokens[evm.BTC].L1ContractAddress)
	addresses = append(addresses, w.Tokens[evm.ETH].L1ContractAddress)
	return addresses
}
