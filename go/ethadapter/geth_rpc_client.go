package ethadapter

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/contracts/generated/ManagementContract"
	"github.com/obscuronet/go-obscuro/go/common/retry"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/common"

	"github.com/ethereum/go-ethereum"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	connRetryMaxWait  = 10 * time.Minute // after this duration, we will stop retrying to connect and return the failure
	connRetryInterval = 500 * time.Millisecond
)

// gethRPCClient implements the EthClient interface and allows connection to a real ethereum node
type gethRPCClient struct {
	client     *ethclient.Client  // the underlying eth rpc client
	l2ID       gethcommon.Address // the address of the Obscuro node this client is dedicated to
	timeout    time.Duration      // the timeout for connecting to, or communicating with, the L1 node
	logger     gethlog.Logger
	rpcAddress string
}

// NewEthClientFromAddress instantiates a new ethadapter.EthClient that connects to an ethereum node
func NewEthClientFromAddress(rpcAddress string, timeout time.Duration, l2ID gethcommon.Address, logger gethlog.Logger) (EthClient, error) {
	client, err := connect(rpcAddress, timeout)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the eth node (%s) - %w", rpcAddress, err)
	}

	logger.Trace(fmt.Sprintf("Initialized eth node connection - addr: %s", rpcAddress))
	return &gethRPCClient{
		client:     client,
		l2ID:       l2ID,
		timeout:    timeout,
		logger:     logger,
		rpcAddress: rpcAddress,
	}, nil
}

// NewEthClient instantiates a new ethadapter.EthClient that connects to an ethereum node
func NewEthClient(ipaddress string, port uint, timeout time.Duration, l2ID gethcommon.Address, logger gethlog.Logger) (EthClient, error) {
	rpcAddress := fmt.Sprintf("ws://%s:%d", ipaddress, port)
	client, err := connect(rpcAddress, timeout)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the eth node (%s) - %w", rpcAddress, err)
	}

	logger.Trace(fmt.Sprintf("Initialized eth node connection - addr: %s port: %d", ipaddress, port))
	return &gethRPCClient{
		client:     client,
		l2ID:       l2ID,
		timeout:    timeout,
		logger:     logger.New(log.PackageKey, "gethrpcclient"),
		rpcAddress: rpcAddress,
	}, nil
}

func (e *gethRPCClient) FetchHeadBlock() (*types.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.BlockByNumber(ctx, nil)
}

func (e *gethRPCClient) Info() Info {
	return Info{
		L2ID: e.l2ID,
	}
}

func (e *gethRPCClient) BlocksBetween(startingBlock *types.Block, lastBlock *types.Block) []*types.Block {
	var blocksBetween []*types.Block
	var err error

	for currentBlk := lastBlock; currentBlk != nil && !bytes.Equal(currentBlk.Hash().Bytes(), startingBlock.Hash().Bytes()) && !bytes.Equal(currentBlk.ParentHash().Bytes(), gethcommon.HexToHash("").Bytes()); {
		currentBlk, err = e.BlockByHash(currentBlk.ParentHash())
		if err != nil {
			e.logger.Crit(fmt.Sprintf("could not fetch parent block with hash %s.", currentBlk.ParentHash().String()), log.ErrKey, err)
		}
		blocksBetween = append(blocksBetween, currentBlk)
	}

	return blocksBetween
}

func (e *gethRPCClient) IsBlockAncestor(block *types.Block, maybeAncestor common.L1BlockHash) bool {
	if bytes.Equal(maybeAncestor.Bytes(), block.Hash().Bytes()) || bytes.Equal(maybeAncestor.Bytes(), (common.L1BlockHash{}).Bytes()) {
		return true
	}

	if block.Number().Int64() == int64(common.L1GenesisHeight) {
		return false
	}

	resolvedBlock, err := e.BlockByHash(maybeAncestor)
	if err != nil {
		e.logger.Crit(fmt.Sprintf("could not fetch parent block with hash %s.", maybeAncestor.String()), log.ErrKey, err)
	}
	if resolvedBlock == nil {
		if resolvedBlock.Number().Int64() >= block.Number().Int64() {
			return false
		}
	}

	p, err := e.BlockByHash(block.ParentHash())
	if err != nil {
		e.logger.Crit(fmt.Sprintf("could not fetch parent block with hash %s", block.ParentHash().String()), log.ErrKey, err)
	}
	if p == nil {
		return false
	}

	return e.IsBlockAncestor(p, maybeAncestor)
}

func (e *gethRPCClient) SendTransaction(signedTx *types.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.SendTransaction(ctx, signedTx)
}

func (e *gethRPCClient) TransactionReceipt(hash gethcommon.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.TransactionReceipt(ctx, hash)
}

func (e *gethRPCClient) Nonce(account gethcommon.Address) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.PendingNonceAt(ctx, account)
}

func (e *gethRPCClient) BlockListener() (chan *types.Header, ethereum.Subscription) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// we do not buffer here, we expect the consumer to always be ready to receive new blocks and not fall behind
	ch := make(chan *types.Header, 1)
	var sub ethereum.Subscription
	var err error
	err = retry.Do(func() error {
		sub, err = e.client.SubscribeNewHead(ctx, ch)
		if err != nil {
			e.logger.Warn("could not subscribe for new head blocks", log.ErrKey, err)
		}
		return err
	}, retry.NewTimeoutStrategy(connRetryMaxWait, connRetryInterval))
	if err != nil {
		// todo (#1638) - handle this scenario better. Health monitor to report L1 unavailable to node operator, be able to recover without restarting host.
		// couldn't connect after timeout period, will not continue
		e.logger.Crit("could not subscribe for new head blocks.", log.ErrKey, err)
	}

	return ch, sub
}

func (e *gethRPCClient) BlockNumber() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.BlockNumber(ctx)
}

func (e *gethRPCClient) BlockByNumber(n *big.Int) (*types.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.BlockByNumber(ctx, n)
}

func (e *gethRPCClient) BlockByHash(hash gethcommon.Hash) (*types.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.BlockByHash(ctx, hash)
}

func (e *gethRPCClient) CallContract(msg ethereum.CallMsg) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.CallContract(ctx, msg, nil)
}

func (e *gethRPCClient) EthClient() *ethclient.Client {
	return e.client
}

func (e *gethRPCClient) BalanceAt(address gethcommon.Address, blockNum *big.Int) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.BalanceAt(ctx, address, blockNum)
}

func (e *gethRPCClient) Stop() {
	e.client.Close()
}

func (e *gethRPCClient) FetchLastBatchSeqNo(address gethcommon.Address) (*big.Int, error) {
	contract, err := ManagementContract.NewManagementContract(address, e.EthClient())
	if err != nil {
		return nil, err
	}

	return contract.LastBatchSeqNo(&bind.CallOpts{})
}

// EstimateGasAndGasPrice takes a txData type and overrides the Gas and Gas Price field with estimated values
func (e *gethRPCClient) EstimateGasAndGasPrice(txData types.TxData, from gethcommon.Address) (types.TxData, error) {
	unEstimatedTx := types.NewTx(txData)
	gasPrice, err := e.EthClient().SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	gasLimit, err := e.EthClient().EstimateGas(context.Background(), ethereum.CallMsg{
		From:  from,
		To:    unEstimatedTx.To(),
		Value: unEstimatedTx.Value(),
		Data:  unEstimatedTx.Data(),
	})
	if err != nil {
		return nil, err
	}

	return &types.LegacyTx{
		Nonce:    unEstimatedTx.Nonce(),
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       unEstimatedTx.To(),
		Value:    unEstimatedTx.Value(),
		Data:     unEstimatedTx.Data(),
	}, nil
}

// Reconnect closes the existing client connection and creates a new connection to the same address:port
func (e *gethRPCClient) Reconnect() error {
	e.client.Close()

	client, err := connect(e.rpcAddress, e.timeout)
	if err != nil {
		return fmt.Errorf("unable to connect to the eth node (%s) - %w", e.rpcAddress, err)
	}
	e.client = client
	return nil
}

// Alive tests the client
func (e *gethRPCClient) Alive() bool {
	_, err := e.client.BlockNumber(context.Background())
	if err != nil {
		e.logger.Error("Unable to fetch BlockNumber rpc endpoint - client connection is in error state")
		return false
	}
	return err == nil
}

func connect(rpcAddress string, connectionTimeout time.Duration) (*ethclient.Client, error) {
	var err error
	var c *ethclient.Client
	for start := time.Now(); time.Since(start) < connectionTimeout; time.Sleep(time.Second) {
		c, err = ethclient.Dial(rpcAddress)
		if err == nil {
			break
		}
	}

	return c, err
}
