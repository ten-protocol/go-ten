package ethadapter

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ten-protocol/go-ten/contracts/generated/DataAvailabilityRegistry"

	"github.com/TwiN/gocache/v2"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	gethlog "github.com/ethereum/go-ethereum/log"

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
	client      *ethclient.Client // the underlying eth rpc client
	timeout     time.Duration     // the timeout for connecting to, or communicating with, the L1 node
	logger      gethlog.Logger
	rpcURL      string
	headerCache *gocache.Cache
	blockCache  *gocache.Cache
}

// NewEthClientFromURL instantiates a new ethadapter.EthClient that connects to an ethereum node
func NewEthClientFromURL(rpcURL string, timeout time.Duration, logger gethlog.Logger) (EthClient, error) {
	client, err := connect(rpcURL, timeout, logger)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the eth node (%s) - %w", rpcURL, err)
	}

	logger.Trace(fmt.Sprintf("Initialized eth node connection - addr: %s", rpcURL))

	return &gethRPCClient{
		client:      client,
		timeout:     timeout,
		logger:      logger,
		rpcURL:      rpcURL,
		headerCache: newFifoCache(_defaultBlockCacheSize, 5*time.Minute),
		blockCache:  newFifoCache(_defaultBlockCacheSize, 5*time.Minute),
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
	// buffer allows the subscription to receive blocks even when consumer is busy processing
	ch := make(chan *types.Header, 10)
	if len(ch) > 5 {
		e.logger.Warn("L1 block channel filling up", "buffered", len(ch))
	}
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
	if n != nil {
		if header := e.cachedHeaderByNumber(n); header != nil {
			return header, nil
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	b, err := e.client.BlockByNumber(ctx, n)
	if err != nil {
		return nil, err
	}
	header := b.Header()
	e.cacheHeader(header)
	return header, nil
}

func (e *gethRPCClient) HeaderByHash(hash gethcommon.Hash) (*types.Header, error) {
	if header := e.cachedHeaderByHash(hash); header != nil {
		return header, nil
	}

	// not in cache, fetch from RPC
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	block, err := e.client.BlockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	header := block.Header()
	e.cacheHeader(header)
	return header, nil
}

func (e *gethRPCClient) BlockByHash(hash gethcommon.Hash) (*types.Block, error) {
	if blk := e.cachedBlockByHash(hash); blk != nil {
		return blk, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	block, err := e.client.BlockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	e.cacheHeader(block.Header())
	e.cacheBlock(block)
	return block, nil
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

func (e *gethRPCClient) SupportsEventLogs() bool {
	return true
}

func (e *gethRPCClient) Stop() {
	e.headerCache.StopJanitor()
	e.blockCache.StopJanitor()
	e.client.Close()
}

func (e *gethRPCClient) FetchLastBatchSeqNo(address gethcommon.Address) (*big.Int, error) {
	contract, err := DataAvailabilityRegistry.NewDataAvailabilityRegistry(address, e.EthClient())
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

	client, err := connect(e.rpcURL, e.timeout, e.logger)
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
		e.logger.Error("Unable to fetch BlockNumber rpc endpoint - client connection is in error state. Cause: %w", err)
		return false
	}
	return true
}

func connect(rpcURL string, connectionTimeout time.Duration, logger gethlog.Logger) (*ethclient.Client, error) {
	var err error
	var c *ethclient.Client
	retryInterval := 1 * time.Second
	for start := time.Now(); time.Since(start) < connectionTimeout; time.Sleep(retryInterval) {
		c, err = ethclient.Dial(rpcURL)
		if err == nil {
			break
		}
		logger.Warn("Unable to connect to ethereum node", "ethNodeURL", rpcURL, "err", err, "retryAfter", retryInterval)
	}

	return c, err
}

func (e *gethRPCClient) cacheHeader(header *types.Header) {
	if header == nil {
		return
	}
	e.headerCache.Set(headerHashCacheKey(header.Hash()), *header)
	if header.Number != nil {
		e.headerCache.Set(headerNumberCacheKey(header.Number), *header)
	}
}

func (e *gethRPCClient) cachedHeaderByHash(hash gethcommon.Hash) *types.Header {
	if cached, found := e.headerCache.Get(headerHashCacheKey(hash)); found {
		if header, ok := cached.(types.Header); ok {
			return &header
		}
	}
	return nil
}

func (e *gethRPCClient) cachedHeaderByNumber(num *big.Int) *types.Header {
	if num == nil {
		return nil
	}
	if cached, found := e.headerCache.Get(headerNumberCacheKey(num)); found {
		if h, ok := cached.(types.Header); ok {
			return &h
		}
	}
	return nil
}

func (e *gethRPCClient) cacheBlock(block *types.Block) {
	if block == nil {
		return
	}
	e.blockCache.Set(blockHashCacheKey(block.Hash()), block)
}

func (e *gethRPCClient) cachedBlockByHash(hash gethcommon.Hash) *types.Block {
	if cached, found := e.blockCache.Get(blockHashCacheKey(hash)); found {
		if block, ok := cached.(*types.Block); ok {
			return block
		}
	}
	return nil
}

func headerHashCacheKey(hash gethcommon.Hash) string {
	return "header_hash:" + hash.Hex()
}

func headerNumberCacheKey(num *big.Int) string {
	return "header_number:" + num.String()
}

func blockHashCacheKey(hash gethcommon.Hash) string {
	return "block_hash:" + hash.Hex()
}
