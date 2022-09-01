package contractdeployer

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
)

func prepareEthDeployer(cfg *Config) (contractDeployerClient, error) {
	client, err := ethadapter.NewEthClient(cfg.NodeHost, cfg.NodePort, 30*time.Second, common.HexToAddress("0x0"))
	if err != nil {
		return nil, err
	}
	return &EthDeployer{client: client}, nil
}

type EthDeployer struct {
	client ethadapter.EthClient
}

func (e *EthDeployer) Nonce(address common.Address) (uint64, error) {
	return e.client.Nonce(address)
}

func (e *EthDeployer) SendTransaction(tx *types.Transaction) error {
	return e.client.SendTransaction(tx)
}

func (e *EthDeployer) TransactionReceipt(hash common.Hash) (*types.Receipt, error) {
	return e.client.TransactionReceipt(hash)
}
