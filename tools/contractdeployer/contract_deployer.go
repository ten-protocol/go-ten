package contractdeployer

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
	"github.com/obscuronet/obscuro-playground/integration/erc20contract"
)

type ContractDeployer struct {
	config *Config
}

func NewContractDeployer(config *Config) *ContractDeployer {
	return &ContractDeployer{
		config: config,
	}
}

func (cd *ContractDeployer) Run() error {
	// connect to the L1
	privateKey, err := crypto.HexToECDSA(cd.config.PrivateKey)
	if err != nil {
		return fmt.Errorf("could not recover private key from hex. Cause: %w", err)
	}

	// load the wallet
	w := wallet.NewInMemoryWalletFromPK(cd.config.ChainID, privateKey)
	// connect to the l1
	client, err := ethclient.NewEthClient(cd.config.L1NodeHost, cd.config.L1NodePort, 30*time.Second)
	if err != nil {
		return fmt.Errorf("unable to connect to the l1 host: %w", err)
	}

	// deploy the contracts
	contractAddr, err := deployContract(client, w, common.Hex2Bytes(mgmtcontractlib.MgmtContractByteCode))
	if err != nil {
		return fmt.Errorf("unable to deploy contract to the l1 host: %w", err)
	}
	erc20contractAddr, err := deployContract(client, w, common.Hex2Bytes(erc20contract.ContractByteCode))
	if err != nil {
		return fmt.Errorf("unable to deploy contract to the l1 host: %w", err)
	}

	// print the contract address
	fmt.Printf("{\"MgmtContractAddr\":\"%s\", \"ERC20ContractAddr\":\"%s\"}",
		contractAddr.Hex(),
		erc20contractAddr.Hex(),
	)
	time.Sleep(5 * time.Second)
	return nil
}

// deployContract deploys a contract (with a tremendous amount of gas)
func deployContract(c ethclient.EthClient, w wallet.Wallet, contractBytes []byte) (*common.Address, error) {
	nonce, err := c.Nonce(w.Address())
	if err != nil {
		return nil, err
	}
	w.SetNonce(nonce)

	deployContractTx := types.LegacyTx{
		Nonce:    w.GetNonceAndIncrement(),
		GasPrice: big.NewInt(2000000000),
		Gas:      1025_000_000,
		Data:     contractBytes,
	}

	signedTx, err := w.SignTransaction(&deployContractTx)
	if err != nil {
		return nil, err
	}

	err = c.SendTransaction(signedTx)
	if err != nil {
		return nil, err
	}

	var start time.Time
	var receipt *types.Receipt
	for start = time.Now(); time.Since(start) < 80*time.Second; time.Sleep(2 * time.Second) {
		receipt, err = c.TransactionReceipt(signedTx.Hash())
		if err == nil && receipt != nil {
			if receipt.Status != types.ReceiptStatusSuccessful {
				return nil, errors.New("unable to deploy contract")
			}
			log.Info("Contract successfully deployed to %s", receipt.ContractAddress)
			return &receipt.ContractAddress, nil
		}

		log.Info("Contract deploy tx has not been mined into a block after %s...", time.Since(start))
	}

	return nil, fmt.Errorf("failed to mine contract deploy tx into a block after %s. Aborting", time.Since(start))
}
