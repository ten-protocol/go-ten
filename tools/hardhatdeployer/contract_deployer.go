package contractdeployer

// todo (@pedro) - we might merge this with the network manager package
import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common/constants"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration/erc20contract"

	gethlog "github.com/ethereum/go-ethereum/log"
)

// The types of contracts supported by the deployer
const (
	// mgmtContract        = "MGMT"
	enclaveRegistry     = "ENCLAVE_REGISTRY"
	rollup              = "ROLLUP"
	crossChain          = "CROSS_CHAIN"
	networkConfig       = "NETWORK_CONFIG"
	Layer2Erc20Contract = "Layer2ERC20"
	layer1Erc20Contract = "Layer1ERC20"
)

const (
	timeoutWait   = 120 * time.Second
	retryInterval = 2 * time.Second
	Prealloc      = 2_050_000_000_000_000_000 // The amount preallocated to the contract deployer wallet.
)

// contractDeployerClient provides a common interface for the L1 and the L2 client
type contractDeployerClient interface {
	Nonce(address common.Address) (uint64, error)
	SendTransaction(tx *types.Transaction) error
	TransactionReceipt(hash common.Hash) (*types.Receipt, error)
}

type contractDeployer struct {
	deployer     contractDeployerClient
	wallet       wallet.Wallet
	contractCode []byte
	logger       gethlog.Logger
}

// Deploy deploys the contract specified in the config, and returns its deployed address.
func Deploy(config *Config, logger gethlog.Logger) (string, error) {
	deployer, err := newContractDeployer(config, logger)
	if err != nil {
		return "", err
	}
	return deployer.run()
}

func newContractDeployer(config *Config, logger gethlog.Logger) (*contractDeployer, error) {
	cfgStr, _ := json.MarshalIndent(config, "", "  ")
	fmt.Printf("Preparing contract deployer with config: %s\n", cfgStr)
	wal, err := setupWallet(config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to setup wallet - %w", err)
	}

	deployerClient, err := prepareDeployerClient(config, wal, logger)
	if err != nil {
		return nil, err
	}

	contractCode, err := getContractCode(config)
	if err != nil {
		return nil, fmt.Errorf("failed to find contract bytecode to deploy - %w", err)
	}

	deployer := &contractDeployer{
		wallet:       wal,
		deployer:     deployerClient,
		contractCode: contractCode,
		logger:       logger,
	}

	return deployer, nil
}

func (cd *contractDeployer) run() (string, error) {
	nonce, err := cd.deployer.Nonce(cd.wallet.Address())
	if err != nil {
		return "", fmt.Errorf("failed to fetch wallet nonce: %w", err)
	}
	cd.wallet.SetNonce(nonce)

	deployContractTx := types.LegacyTx{
		Nonce:    cd.wallet.GetNonceAndIncrement(),
		GasPrice: big.NewInt(params.InitialBaseFee),
		Gas:      uint64(50_000_000),
		Data:     cd.contractCode,
	}

	contractAddr, err := cd.signAndSendTxWithReceipt(cd.wallet, cd.deployer, &deployContractTx)
	if err != nil {
		return "", err
	}
	if contractAddr == nil {
		return "", fmt.Errorf("transaction was successful but could not retrieve address for deployed contract")
	}

	return contractAddr.Hex(), nil
}

func setupWallet(cfg *Config, logger gethlog.Logger) (wallet.Wallet, error) {
	privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("could not recover private key from hex. Cause: %w", err)
	}

	// load the wallet
	return wallet.NewInMemoryWalletFromPK(cfg.ChainID, privateKey, logger), nil
}

func prepareDeployerClient(config *Config, wal wallet.Wallet, logger gethlog.Logger) (contractDeployerClient, error) {
	if config.IsL1Deployment {
		return prepareEthDeployer(config, logger)
	}
	return prepareObscuroDeployer(config, wal, logger)
}

func (cd *contractDeployer) signAndSendTxWithReceipt(wallet wallet.Wallet, deployer contractDeployerClient, tx *types.LegacyTx) (*common.Address, error) {
	signedTx, err := wallet.SignTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to sign contract deploy transaction: %w", err)
	}

	err = deployer.SendTransaction(signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send contract deploy transaction: %w", err)
	}

	cd.logger.Info(fmt.Sprintf("Waiting (up to %s) for deploy tx to be mined into a block...", timeoutWait))

	var receipt *types.Receipt
	err = retry.Do(func() error {
		receipt, err = deployer.TransactionReceipt(signedTx.Hash())
		return err
	}, retry.NewTimeoutStrategy(timeoutWait, retryInterval))
	if err != nil {
		return nil, fmt.Errorf("failed to deploy contract - %w", err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("unable to deploy contract, receipt status unsuccessful: %v", receipt)
	}

	return &receipt.ContractAddress, nil
}

func getContractCode(cfg *Config) ([]byte, error) {
	switch cfg.ContractName {
	case enclaveRegistry:
		return constants.EnclaveRegistryBytecode()
	case crossChain:
		return constants.CrossChainBytecode()
	case rollup:
		return constants.DataAvailabilityRegistryBytecode()
	case networkConfig:
		return constants.NetworkConfigBytecode()

	case Layer2Erc20Contract:
		tokenName := cfg.ConstructorParams[0]
		tokenSymbol := cfg.ConstructorParams[1]
		supply := cfg.ConstructorParams[2]
		// 0x526c84529b2b8c11f57d93d3f5537aca3aecef9b - address of the L2 message bus contract. todo (@stefan) - remove once there is a proper way to extract it.
		return erc20contract.L2Bytecode(tokenName, tokenSymbol, supply, common.HexToAddress("0x526c84529b2b8c11f57d93d3f5537aca3aecef9b")), nil

	case layer1Erc20Contract:
		tokenName := cfg.ConstructorParams[0]
		tokenSymbol := cfg.ConstructorParams[1]
		supply := cfg.ConstructorParams[2]
		mgmtAddr := common.HexToAddress(cfg.ConstructorParams[3])

		return erc20contract.L1Bytecode(tokenName, tokenSymbol, supply, mgmtAddr), nil

	default:
		return nil, fmt.Errorf("unrecognised contract %s - no bytecode configured for that contract name", cfg.ContractName)
	}
}
