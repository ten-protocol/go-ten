package contractdeployer

// TODO we might merge this with the network manager package
import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/obscuronet/go-obscuro/contracts/managementcontract"
	"github.com/obscuronet/go-obscuro/integration/erc20contract"
	"github.com/obscuronet/go-obscuro/integration/guessinggame"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/wallet"
)

// The types of contracts supported by the deployer
const (
	mgmtContract         = "MGMT"
	l2Erc20Contract      = "L2ERC20"
	l1Erc20Contract      = "L1ERC20"
	GuessingGameContract = "GUESS"
)

const (
	timeoutWait   = 80 * time.Second
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
}

// Deploy deploys the contract specified in the config, and returns its deployed address.
func Deploy(config *Config) (string, error) {
	deployer, err := newContractDeployer(config)
	if err != nil {
		return "", err
	}
	return deployer.run()
}

func newContractDeployer(config *Config) (*contractDeployer, error) {
	cfgStr, _ := json.MarshalIndent(config, "", "  ")
	fmt.Printf("Preparing contract deployer with config: %s\n", cfgStr)
	wal, err := setupWallet(config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup wallet - %w", err)
	}

	deployerClient, err := prepareDeployerClient(config, wal)
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
		GasPrice: big.NewInt(2000000000),
		Gas:      1025_000_000,
		Data:     cd.contractCode,
	}

	contractAddr, err := signAndSendTxWithReceipt(cd.wallet, cd.deployer, &deployContractTx)
	if err != nil {
		return "", err
	}
	if contractAddr == nil {
		return "", fmt.Errorf("transaction was successful but could not retrieve address for deployed contract")
	}

	return contractAddr.Hex(), nil
}

func setupWallet(cfg *Config) (wallet.Wallet, error) {
	privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("could not recover private key from hex. Cause: %w", err)
	}

	// load the wallet
	return wallet.NewInMemoryWalletFromPK(cfg.ChainID, privateKey), nil
}

func prepareDeployerClient(config *Config, wal wallet.Wallet) (contractDeployerClient, error) {
	if config.IsL1Deployment {
		return prepareEthDeployer(config)
	}
	return prepareObscuroDeployer(config, wal)
}

func signAndSendTxWithReceipt(wallet wallet.Wallet, deployer contractDeployerClient, tx *types.LegacyTx) (*common.Address, error) {
	signedTx, err := wallet.SignTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to sign contract deploy transaction: %w", err)
	}

	err = deployer.SendTransaction(signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send contract deploy transaction: %w", err)
	}

	var start time.Time
	for start = time.Now(); time.Since(start) < timeoutWait; time.Sleep(retryInterval) {
		receipt, err := deployer.TransactionReceipt(signedTx.Hash())
		if err == nil && receipt != nil {
			if receipt.Status != types.ReceiptStatusSuccessful {
				return nil, fmt.Errorf("unable to deploy contract, receipt status unsuccessful: %v", receipt)
			}
			log.Info("Contract successfully deployed to %s", receipt.ContractAddress)
			return &receipt.ContractAddress, nil
		}

		log.Info("Contract deploy tx %s has not been mined into a block after %s...", signedTx.Hash(), time.Since(start))
	}
	return nil, fmt.Errorf("failed to mine contract deploy tx %s into a block after %s. Aborting", signedTx.Hash(), time.Since(start))
}

func getContractCode(cfg *Config) ([]byte, error) {
	switch cfg.ContractName {
	case mgmtContract:
		return managementcontract.Bytecode()

	case l2Erc20Contract:
		tokenName := cfg.ConstructorParams[0]
		tokenSymbol := cfg.ConstructorParams[1]
		supply := cfg.ConstructorParams[2]
		return erc20contract.L2Bytecode(tokenName, tokenSymbol, supply), nil

	case l1Erc20Contract:
		tokenName := cfg.ConstructorParams[0]
		tokenSymbol := cfg.ConstructorParams[1]
		supply := cfg.ConstructorParams[2]
		return erc20contract.L1Bytecode(tokenName, tokenSymbol, supply), nil

	case GuessingGameContract:
		size, err := strconv.Atoi(cfg.ConstructorParams[0])
		if err != nil {
			return nil, err
		}
		address := common.HexToAddress(cfg.ConstructorParams[1])

		return guessinggame.Bytecode(uint8(size), address)

	default:
		return nil, fmt.Errorf("unrecognised contract %s - no bytecode configured for that contract name", cfg.ContractName)
	}
}
