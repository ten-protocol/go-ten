package params

import (
	"math/big"

	"github.com/ten-protocol/go-ten/integration/common/testlog"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ten-protocol/go-ten/go/wallet"
	testcommon "github.com/ten-protocol/go-ten/integration/common"
	"github.com/ten-protocol/go-ten/integration/datagenerator"
)

// SimToken - mapping between the ERC20s on Ethereum and Obscuro. This holds both the contract addresses and the keys of the contract owners,
// because it needs to sign transactions and deploy contracts.
// Note: For now the l2 values are taken from the "bridge" inside the Obscuro core.
type SimToken struct {
	Name testcommon.ERC20

	L1Owner           wallet.Wallet
	L1ContractAddress *common.Address

	L2Owner           wallet.Wallet
	L2ContractAddress *common.Address
}

type L1PrefundWallets struct {
	HOC    wallet.Wallet
	POC    wallet.Wallet
	Faucet wallet.Wallet
}

type SimWallets struct {
	ContractOwnerWallet wallet.Wallet   // owner of the contracts deployed on Ethereum
	NodeWallets         []wallet.Wallet // the keys used by the obscuro nodes to submit rollups to Eth

	SimEthWallets []wallet.Wallet // the wallets of the simulated users on the Ethereum side
	SimObsWallets []wallet.Wallet // and their equivalents on the obscuro side (with a different chainId)

	GasBridgeWallet     wallet.Wallet
	GasWithdrawalWallet wallet.Wallet

	L2FaucetWallet wallet.Wallet // the wallet of the L2 faucet
	L2FeesWallet   wallet.Wallet
	Tokens         map[testcommon.ERC20]*SimToken // The supported tokens

	PrefundedEthWallets L1PrefundWallets
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
		simObsWallets[i] = wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), simEthWallets[i].PrivateKey(), testlog.Logger())
	}

	// create the wallet to deploy the Management contract
	mcOwnerWallet := datagenerator.RandomWallet(ethereumChainID)

	// create the L2 faucet wallet
	l2FaucetPrivKey, err := crypto.HexToECDSA(testcommon.TestnetPrefundedPK)
	if err != nil {
		panic("could not initialise L2 faucet private key")
	}
	l2FaucetWallet := wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), l2FaucetPrivKey, testlog.Logger())

	GasBridgingKeys, _ := crypto.GenerateKey()
	GasWithdrawalKeys, _ := crypto.GenerateKey()
	gasWallet := wallet.NewInMemoryWalletFromPK(big.NewInt(ethereumChainID), GasBridgingKeys, testlog.Logger())
	withdrawalWallet := wallet.NewInMemoryWalletFromPK(big.NewInt(ethereumChainID), GasWithdrawalKeys, testlog.Logger())

	sequencerGasKeys, _ := crypto.GenerateKey()
	sequencerFeeWallet := wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), sequencerGasKeys, testlog.Logger())

	// create the L1 addresses of the two tokens, and connect them to the hardcoded addresses from the enclave
	hoc := SimToken{
		Name:              testcommon.HOC,
		L1Owner:           datagenerator.RandomWallet(ethereumChainID),
		L2Owner:           wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), testcommon.HOCOwner, testlog.Logger()),
		L2ContractAddress: &testcommon.HOCContract,
	}
	poc := SimToken{
		Name:              testcommon.POC,
		L1Owner:           datagenerator.RandomWallet(ethereumChainID),
		L2Owner:           wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), testcommon.POCOwner, testlog.Logger()),
		L2ContractAddress: &testcommon.POCContract,
	}

	return &SimWallets{
		ContractOwnerWallet: mcOwnerWallet,
		NodeWallets:         nodeWallets,
		SimEthWallets:       simEthWallets,
		SimObsWallets:       simObsWallets,
		L2FaucetWallet:      l2FaucetWallet,
		L2FeesWallet:        sequencerFeeWallet,
		GasBridgeWallet:     gasWallet,
		GasWithdrawalWallet: withdrawalWallet,
		Tokens: map[testcommon.ERC20]*SimToken{
			testcommon.HOC: &hoc,
			testcommon.POC: &poc,
		},
		PrefundedEthWallets: L1PrefundWallets{
			HOC:    datagenerator.RandomWallet(ethereumChainID),
			POC:    datagenerator.RandomWallet(ethereumChainID),
			Faucet: datagenerator.RandomWallet(ethereumChainID),
		},
	}
}

func (w *SimWallets) AllEthWallets() []wallet.Wallet {
	ethWallets := make([]wallet.Wallet, 0)
	for _, token := range w.Tokens {
		ethWallets = append(ethWallets, token.L1Owner)
	}
	ethWallets = append(ethWallets, w.GasBridgeWallet)
	ethWallets = append(ethWallets, w.GasWithdrawalWallet)
	ethWallets = append(ethWallets, w.PrefundedEthWallets.POC)
	ethWallets = append(ethWallets, w.PrefundedEthWallets.HOC)
	ethWallets = append(ethWallets, w.PrefundedEthWallets.Faucet)
	return append(append(append(w.NodeWallets, w.SimEthWallets...), w.ContractOwnerWallet), ethWallets...)
}

func (w *SimWallets) AllObsWallets() []wallet.Wallet {
	obsWallets := make([]wallet.Wallet, 0)
	for _, token := range w.Tokens {
		obsWallets = append(obsWallets, token.L2Owner)
	}
	return append(w.SimObsWallets, obsWallets...)
}
