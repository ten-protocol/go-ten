package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/obscuronet/go-obscuro/go/enclave/bridge"

	"github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
)

const (
	faucetPrivateKeyHex = "8dfb8083da6275ae3e4f41e3e8a8c19d028d32c9247e24530933782f2a05035b" // The faucet's private key.
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

	L2FaucetWallet wallet.Wallet              // the wallet of the L2 faucet
	Tokens         map[bridge.ERC20]*SimToken // The supported tokens
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

	// create the L2 faucet wallet
	l2FaucetPrivKey, err := crypto.HexToECDSA(faucetPrivateKeyHex)
	if err != nil {
		panic("could not initialise L2 faucet private key")
	}
	l2FaucetWallet := wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), l2FaucetPrivKey)

	// create the L1 addresses of the two tokens, and connect them to the hardcoded addresses from the enclave
	obx := SimToken{
		Name:              bridge.OBX,
		L1Owner:           datagenerator.RandomWallet(ethereumChainID),
		L2Owner:           wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), bridge.WOBXOwner),
		L2ContractAddress: &bridge.WOBXContract,
	}
	eth := SimToken{
		Name:              bridge.ETH,
		L1Owner:           datagenerator.RandomWallet(ethereumChainID),
		L2Owner:           wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), bridge.WETHOwner),
		L2ContractAddress: &bridge.WETHContract,
	}

	return &SimWallets{
		MCOwnerWallet:  mcOwnerWallet,
		NodeWallets:    nodeWallets,
		SimEthWallets:  simEthWallets,
		SimObsWallets:  simObsWallets,
		L2FaucetWallet: l2FaucetWallet,
		Tokens: map[bridge.ERC20]*SimToken{
			bridge.OBX: &obx,
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
	addresses = append(addresses, w.Tokens[bridge.OBX].L1ContractAddress)
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
