package params

import (
	"math/big"

	"github.com/obscuronet/obscuro-playground/go/enclave/bridge"

	"github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/obscuro-playground/go/wallet"
	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
)

// SimToken - mapping between the ERC20s on Ethereum and Obscuro. This holds both the contract addresses and the keys of the contract owners,
// because it needs to sign transactions and deploy contracts.
// Note: For now the l2 values are taken from the "bridge" inside the Obscuro core.
type SimToken struct {
	Name bridge.ERC20

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

	Tokens map[bridge.ERC20]*SimToken // The supported tokens
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
		Name:              bridge.BTC,
		L1Owner:           datagenerator.RandomWallet(ethereumChainID),
		L2Owner:           wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), bridge.WBtcOwner),
		L2ContractAddress: &bridge.WBtcContract,
	}
	eth := SimToken{
		Name:              bridge.ETH,
		L1Owner:           datagenerator.RandomWallet(ethereumChainID),
		L2Owner:           wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), bridge.WEthOnwer),
		L2ContractAddress: &bridge.WEthContract,
	}

	return &SimWallets{
		MCOwnerWallet: mcOwnerWallet,
		NodeWallets:   nodeWallets,
		SimEthWallets: simEthWallets,
		SimObsWallets: simObsWallets,
		Tokens: map[bridge.ERC20]*SimToken{
			bridge.BTC: &btc,
			bridge.ETH: &eth,
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
	addresses = append(addresses, w.Tokens[bridge.BTC].L1ContractAddress)
	addresses = append(addresses, w.Tokens[bridge.ETH].L1ContractAddress)
	return addresses
}

func (w *SimWallets) AllObsWallets() []wallet.Wallet {
	obsWallets := make([]wallet.Wallet, 0)
	for _, token := range w.Tokens {
		obsWallets = append(obsWallets, token.L2Owner)
	}
	return append(w.SimObsWallets, obsWallets...)
}
