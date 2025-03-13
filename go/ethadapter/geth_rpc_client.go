package ethadapter

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ten-protocol/go-ten/contracts/generated/RollupContract"

	"github.com/TwiN/gocache/v2"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/retry"
)

const (
	connRetryMaxWait       = 10 * time.Minute // after this duration, we will stop retrying to connect and return the failure
	connRetryInterval      = 500 * time.Millisecond
	_defaultBlockCacheSize = 51 // enough for 50 request batch size and one for previous block

)

// gethRPCClient implements the EthClient interface and allows connection to a real ethereum node
type gethRPCClient struct {
	client     *ethclient.Client // the underlying eth rpc client
	timeout    time.Duration     // the timeout for connecting to, or communicating with, the L1 node
	logger     gethlog.Logger
	rpcURL     string
	blockCache *gocache.Cache
}

// NewEthClientFromURL instantiates a new ethadapter.EthClient that connects to an ethereum node
func NewEthClientFromURL(rpcURL string, timeout time.Duration, logger gethlog.Logger) (EthClient, error) {
	client, err := connect(rpcURL, timeout)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the eth node (%s) - %w", rpcURL, err)
	}

	logger.Trace(fmt.Sprintf("Initialized eth node connection - addr: %s", rpcURL))

	return &gethRPCClient{
		client:     client,
		timeout:    timeout,
		logger:     logger,
		rpcURL:     rpcURL,
		blockCache: newFifoCache(_defaultBlockCacheSize, 5*time.Minute),
	}, nil
}

func newFifoCache(nrElem int, ttl time.Duration) *gocache.Cache {
	cache := gocache.NewCache().WithMaxSize(nrElem).WithEvictionPolicy(gocache.FirstInFirstOut).WithDefaultTTL(ttl)
	err := cache.StartJanitor()
	if err != nil {
		panic("failed to start cache.")
	}
	return cache
}

// NewEthClient instantiates a new ethadapter.EthClient that connects to an ethereum node
func NewEthClient(ipaddress string, port uint, timeout time.Duration, logger gethlog.Logger) (EthClient, error) {
	return NewEthClientFromURL(fmt.Sprintf("ws://%s:%d", ipaddress, port), timeout, logger)
}

func (e *gethRPCClient) FetchHeadBlock() (*types.Header, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.HeaderByNumber(ctx, nil)
}

func (e *gethRPCClient) Info() Info {
	return Info{}
}

func (e *gethRPCClient) BlocksBetween(startingBlock *types.Header, lastBlock *types.Header) ([]*types.Header, error) {
	var blocksBetween []*types.Header
	var err error

	for currentBlk := lastBlock; currentBlk != nil && !bytes.Equal(currentBlk.Hash().Bytes(), startingBlock.Hash().Bytes()) && !bytes.Equal(currentBlk.ParentHash.Bytes(), gethcommon.HexToHash("").Bytes()); {
		c := currentBlk.ParentHash
		currentBlk, err = e.HeaderByHash(currentBlk.ParentHash)
		if err != nil {
			e.logger.Error(fmt.Sprintf("could not fetch parent block with hash %s.", c.String()), log.ErrKey, err)
			return nil, fmt.Errorf("could not fetch parent block with hash %s. cause: %w", c.String(), err)
		}
		blocksBetween = append(blocksBetween, currentBlk)
	}

	return blocksBetween, nil
}

func (e *gethRPCClient) IsBlockAncestor(block *types.Header, maybeAncestor common.L1BlockHash) bool {
	if bytes.Equal(maybeAncestor.Bytes(), block.Hash().Bytes()) || bytes.Equal(maybeAncestor.Bytes(), (common.L1BlockHash{}).Bytes()) {
		return true
	}

	if block.Number.Int64() == int64(common.L1GenesisHeight) {
		return false
	}

	resolvedAncestorBlock, err := e.HeaderByHash(maybeAncestor)
	if err != nil {
		e.logger.Crit(fmt.Sprintf("could not fetch parent block with hash %s.", maybeAncestor.String()), log.ErrKey, err)
	}
	if resolvedAncestorBlock.Number.Int64() >= block.Number.Int64() {
		return false
	}

	p, err := e.HeaderByHash(block.ParentHash)
	if err != nil {
		e.logger.Crit(fmt.Sprintf("could not fetch parent block with hash %s", block.ParentHash.String()), log.ErrKey, err)
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

func (e *gethRPCClient) TransactionByHash(hash gethcommon.Hash) (*types.Transaction, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.TransactionByHash(ctx, hash)
}

func (e *gethRPCClient) Nonce(account gethcommon.Address) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.PendingNonceAt(ctx, account)
}

func (e *gethRPCClient) BlockListener() (chan *types.Header, ethereum.Subscription) {
	// we do not buffer here, we expect the consumer to always be ready to receive new blocks and not fall behind
	ch := make(chan *types.Header)
	var sub ethereum.Subscription
	var err error
	err = retry.Do(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
		defer cancel()
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

func (e *gethRPCClient) HeaderByNumber(n *big.Int) (*types.Header, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	b, err := e.client.BlockByNumber(ctx, n)
	if err != nil {
		return nil, err
	}
	return b.Header(), nil
}

func (e *gethRPCClient) HeaderByHash(hash gethcommon.Hash) (*types.Header, error) {
	cachedBlock, found := e.blockCache.Get(hash.Hex())
	if found {
		h, ok := cachedBlock.(types.Header)
		if !ok {
			return nil, fmt.Errorf("should not happen. could not cast cached block to header")
		}
		return &h, nil
	}

	// not in cache, fetch from RPC
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	block, err := e.client.BlockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	e.blockCache.Set(hash.Hex(), *block.Header())
	return block.Header(), nil
}

func (e *gethRPCClient) BlockByHash(hash gethcommon.Hash) (*types.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.BlockByHash(ctx, hash)
}

func (e *gethRPCClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return e.client.SuggestGasTipCap(ctx)
}

func (e *gethRPCClient) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	return e.client.EstimateGas(ctx, msg)
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

func (e *gethRPCClient) GetLogs(q ethereum.FilterQuery) ([]types.Log, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.FilterLogs(ctx, q)
}

func (e *gethRPCClient) Stop() {
	e.blockCache.StopJanitor()
	e.client.Close()
}

func (e *gethRPCClient) FetchLastBatchSeqNo(address gethcommon.Address) (*big.Int, error) {
	contract, err := RollupContract.NewRollupContract(address, e.EthClient())
	if err != nil {
		return nil, err
	}

	return contract.LastBatchSeqNo(&bind.CallOpts{})
}

// ReconnectIfClosed closes the existing client connection and creates a new connection to the same address:port
func (e *gethRPCClient) ReconnectIfClosed() error {
	if e.Alive() {
		// connection is not closed
		return nil
	}
	e.client.Close()

	client, err := connect(e.rpcURL, e.timeout)
	if err != nil {
		return fmt.Errorf("unable to connect to the eth node (%s) - %w", e.rpcURL, err)
	}
	e.client = client
	return nil
}

// Alive tests the client
func (e *gethRPCClient) Alive() bool {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	_, err := e.client.BlockNumber(ctx)
	if err != nil {
		e.logger.Error("Unable to fetch BlockNumber rpc endpoint - client connection is in error state")
		return false
	}
	return true
}

func connect(rpcURL string, connectionTimeout time.Duration) (*ethclient.Client, error) {
	var err error
	var c *ethclient.Client
	for start := time.Now(); time.Since(start) < connectionTimeout; time.Sleep(time.Second) {
		c, err = ethclient.Dial(rpcURL)
		if err == nil {
			break
		}
	}

	return c, err
}
