package contractdeployer

// TODO we might merge this with the network manager package
import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/obscuronet/go-obscuro/go/enclave/rollupchain"

	"github.com/obscuronet/go-obscuro/contracts/managementcontract"
	"github.com/obscuronet/go-obscuro/integration/erc20contract"
	"github.com/obscuronet/go-obscuro/integration/guessinggame"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/wallet"
)

// The types of contracts supported by the deployer
const (
	mgmtContract         = "MGMT"
	l2Erc20Contract      = "L2ERC20"
	l1Erc20Contract      = "L1ERC20"
	guessingGameContract = "GUESS"
)

const (
	timeoutWait   = 80 * time.Second
	retryInterval = 2 * time.Second
	prealloc      = 2_050_000_000_000_000_000 // The amount preallocated to the contract deployer wallet.
)

type contractDeployer struct {
	client       ethadapter.EthClient
	wallet       wallet.Wallet
	faucetClient ethadapter.EthClient
	faucetWallet wallet.Wallet
	contractCode []byte
}

func Deploy(config *Config) error {
	deployer, err := newContractDeployer(config)
	if err != nil {
		return err
	}
	return deployer.run(config.IsL1Deployment)
}

func newContractDeployer(config *Config) (*contractDeployer, error) {
	cfgStr, _ := json.MarshalIndent(config, "", "  ")
	fmt.Printf("Preparing contract deployer with config: %s\n", cfgStr)
	wal, err := setupWallet(config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup wallet - %w", err)
	}

	client, err := getClient(config, wal)
	if err != nil {
		return nil, err
	}

	contractCode, err := getContractCode(config)
	if err != nil {
		return nil, fmt.Errorf("failed to find contract bytecode to deploy - %w", err)
	}

	deployer := &contractDeployer{
		wallet:       wal,
		client:       client,
		contractCode: contractCode,
	}

	if !config.IsL1Deployment {
		// Create the L2 faucet wallet and client.
		faucetPrivKey, err := crypto.HexToECDSA(rollupchain.FaucetPrivateKeyHex)
		if err != nil {
			panic("could not initialise L2 faucet private key")
		}
		deployer.faucetWallet = wallet.NewInMemoryWalletFromPK(config.ChainID, faucetPrivKey)

		faucetClient, err := getClient(config, deployer.faucetWallet)
		if err != nil {
			return nil, err
		}

		deployer.faucetClient = faucetClient
	}

	return deployer, nil
}

func (cd *contractDeployer) run(isL1Deployment bool) error {
	// On the L2, we need to prefund the contract deployer account.
	if !isL1Deployment {
		err := cd.prefundAccount()
		if err != nil {
			return fmt.Errorf("unable to prefund contract deployer account. Cause: %w", err)
		}
	}

	nonce, err := cd.client.Nonce(cd.wallet.Address())
	if err != nil {
		return fmt.Errorf("failed to fetch wallet nonce: %w", err)
	}
	cd.wallet.SetNonce(nonce)

	deployContractTx := types.LegacyTx{
		Nonce:    cd.wallet.GetNonceAndIncrement(),
		GasPrice: big.NewInt(2000000000),
		Gas:      1025_000_000,
		Data:     cd.contractCode,
	}

	contractAddr, err := signAndSendTxWithReceipt(cd.wallet, cd.client, &deployContractTx)
	if err != nil {
		return err
	}
	if contractAddr == nil {
		return fmt.Errorf("transaction was successful but could not retrieve address for deployed contract")
	}

	// print the contract address, to be read if necessary by the caller (important: this must be the last message output by the script)
	fmt.Print(contractAddr.Hex())

	// this is a safety sleep to make sure the output is printed
	time.Sleep(5 * time.Second)
	return nil
}

func setupWallet(cfg *Config) (wallet.Wallet, error) {
	privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("could not recover private key from hex. Cause: %w", err)
	}

	// load the wallet
	return wallet.NewInMemoryWalletFromPK(cfg.ChainID, privateKey), nil
}

func getClient(config *Config, wal wallet.Wallet) (ethadapter.EthClient, error) {
	var client ethadapter.EthClient
	var err error

	startConnectingTime := time.Now()
	// since the nodes we are connecting to may have only just started, we retry connection until it is successful
	for client == nil && time.Since(startConnectingTime) < timeoutWait {
		client, err = setupClient(config, wal)
		if err == nil {
			break // success
		}
		// if there was an error we'll retry, if we timeout the last seen error will display
		time.Sleep(retryInterval)
	}
	if client == nil {
		return nil, fmt.Errorf("failed to initialise client connection after retrying for %s, %w", timeoutWait, err)
	}

	return client, nil
}

func setupClient(cfg *Config, wal wallet.Wallet) (ethadapter.EthClient, error) {
	if cfg.IsL1Deployment {
		// return a connection to the l1
		return ethadapter.NewEthClient(cfg.NodeHost, cfg.NodePort, 30*time.Second, common.HexToAddress("0x0"))
	}
	return ethadapter.NewObscuroRPCClient(cfg.NodeHost, cfg.NodePort, wal)
}

func (cd *contractDeployer) prefundAccount() error {
	nonce, err := cd.faucetClient.Nonce(cd.faucetWallet.Address())
	if err != nil {
		return fmt.Errorf("failed to fetch faucet nonce. Cause: %w", err)
	}

	destAddr := cd.wallet.Address()
	tx := &types.LegacyTx{
		Nonce:    nonce,
		Value:    big.NewInt(prealloc),
		Gas:      uint64(1_000_000),
		GasPrice: common.Big1,
		To:       &destAddr,
	}
	_, err = signAndSendTxWithReceipt(cd.faucetWallet, cd.faucetClient, tx)
	if err != nil {
		return fmt.Errorf("failed to complete faucet transaction. Cause: %w", err)
	}

	return nil
}

func signAndSendTxWithReceipt(wallet wallet.Wallet, client ethadapter.EthClient, tx *types.LegacyTx) (*common.Address, error) {
	signedTx, err := wallet.SignTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to sign contract deploy transaction: %w", err)
	}

	err = client.SendTransaction(signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send contract deploy transaction: %w", err)
	}

	var start time.Time
	for start = time.Now(); time.Since(start) < timeoutWait; time.Sleep(retryInterval) {
		receipt, err := client.TransactionReceipt(signedTx.Hash())
		if err != nil {
			log.Info(err.Error())
		}
		if err == nil && receipt != nil {
			if receipt.Status != types.ReceiptStatusSuccessful {
				return nil, fmt.Errorf("unable to deploy contract, receipt status unsuccessful: %v", receipt)
			}
			log.Info("Contract successfully deployed to %s", receipt.ContractAddress)
			return &receipt.ContractAddress, nil
		}

		log.Info("Contract deploy tx has not been mined into a block after %s...", time.Since(start))
	}
	return nil, fmt.Errorf("failed to mine contract deploy tx %s into a block after %s. Aborting", signedTx.Hash(), time.Since(start))
}

func getContractCode(cfg *Config) ([]byte, error) {
	switch cfg.ContractName {
	case mgmtContract:
		return managementcontract.Bytecode()

	case l2Erc20Contract:
		tokenName := cfg.ConstructorParams[0]
		return erc20contract.L2BytecodeWithDefaultSupply(tokenName), nil

	case l1Erc20Contract:
		tokenName := cfg.ConstructorParams[0]
		return erc20contract.L1BytecodeWithDefaultSupply(tokenName), nil

	case guessingGameContract:
		size, err := strconv.Atoi(cfg.ConstructorParams[0])
		if err != nil {
			return nil, err
		}
		address := common.BytesToAddress(common.Hex2Bytes(cfg.ConstructorParams[1]))

		return guessinggame.Bytecode(uint8(size), address)

	default:
		return nil, fmt.Errorf("unrecognised contract %s - no bytecode configured for that contract name", cfg.ContractName)
	}
}
