package params

import (
	"math/big"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/evm"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
)

type SimWallets struct {
	MCOwnerWallet wallet.Wallet   // owner of the management contract deployed on Ethereum
	NodeWallets   []wallet.Wallet // the keys used by the obscuro nodes to submit rollups to Eth

	SimEthWallets []wallet.Wallet // the wallets of the simulated users on the Ethereum side
	SimObsWallets []wallet.Wallet // and their equivalents on the obscuro side (with a different chainId)

	Erc20EthOwnerWallets []wallet.Wallet // the owners of the supported ethereum erc20 contracts
	Erc20ObsOwnerWallets []wallet.Wallet // and the owners of the respective wrapped versions on Obscuro
}

func NewSimWallets(nrSimWallets int, nNodes int, nrErc20s int, ethereumChainID int64, obscuroChainID int64) *SimWallets {
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

	// create the ethereum wallets to be used to deploy ERC20 contracts
	// and their counterparts in the Obscuro world for the wrapped versions
	if nrErc20s != 1 {
		panic("only one erc20 supported for now")
	}
	erc20EthWallets := make([]wallet.Wallet, nrErc20s)
	erc20ObsWallets := make([]wallet.Wallet, nrErc20s)
	erc20EthWallets[0] = datagenerator.RandomWallet(ethereumChainID)

	// this cannot be random for now, because there is hardcoded logic in the obscuro core
	// to generate synthetic "transfer" transactions on the wrapped erc20 for each erc20 deposit on ethereum
	// and these transactions need to be signed
	erc20ObsWallets[0] = wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), evm.Erc20OwnerKey)

	return &SimWallets{
		MCOwnerWallet:        mcOwnerWallet,
		NodeWallets:          nodeWallets,
		SimEthWallets:        simEthWallets,
		SimObsWallets:        simObsWallets,
		Erc20EthOwnerWallets: erc20EthWallets,
		Erc20ObsOwnerWallets: erc20ObsWallets,
	}
}

func (w *SimWallets) AllEthWallets() []wallet.Wallet {
	return append(append(append(w.NodeWallets, w.SimEthWallets...), w.MCOwnerWallet), w.Erc20EthOwnerWallets...)
}
