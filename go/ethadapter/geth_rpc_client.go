package ethadapter

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/common"

	"github.com/ethereum/go-ethereum"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// gethRPCClient implements the EthClient interface and allows connection to a real ethereum node
type gethRPCClient struct {
	client  *ethclient.Client  // the underlying eth rpc client
	l2ID    gethcommon.Address // the address of the Obscuro node this client is dedicated to
	timeout time.Duration      // the timeout for connecting to, or communicating with, the L1 node
}

// NewEthClient instantiates a new ethadapter.EthClient that connects to an ethereum node
func NewEthClient(address string, timeout time.Duration, l2ID gethcommon.Address) (EthClient, error) {
	client, err := connect(address, timeout)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the eth node - %w", err)
	}

	log.Trace("Initialized eth node connection - addr: %s", address)
	return &gethRPCClient{
		client:  client,
		l2ID:    l2ID,
		timeout: timeout,
	}, nil
}

func (e *gethRPCClient) FetchHeadBlock() *types.Block {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	blk, err := e.client.BlockByNumber(ctx, nil)
	if err != nil {
		log.Panic("could not fetch head block. Cause: %s", err)
	}
	return blk
}

func (e *gethRPCClient) Info() Info {
	return Info{
		L2ID: e.l2ID,
	}
}

func (e *gethRPCClient) BlocksBetween(startingBlock *types.Block, lastBlock *types.Block) []*types.Block {
	// TODO this should be a stream
	var blocksBetween []*types.Block
	var err error

	for currentBlk := lastBlock; currentBlk != nil && !bytes.Equal(currentBlk.Hash().Bytes(), startingBlock.Hash().Bytes()) && !bytes.Equal(currentBlk.ParentHash().Bytes(), gethcommon.HexToHash("").Bytes()); {
		currentBlk, err = e.BlockByHash(currentBlk.ParentHash())
		if err != nil {
			log.Panic("could not fetch parent block with hash %s. Cause: %s", currentBlk.ParentHash().String(), err)
		}
		blocksBetween = append(blocksBetween, currentBlk)
	}

	return blocksBetween
}

func (e *gethRPCClient) IsBlockAncestor(block *types.Block, maybeAncestor common.L1RootHash) bool {
	if bytes.Equal(maybeAncestor.Bytes(), block.Hash().Bytes()) || bytes.Equal(maybeAncestor.Bytes(), common.GenesisBlock.Hash().Bytes()) {
		return true
	}

	if block.Number().Int64() == int64(common.L1GenesisHeight) {
		return false
	}

	resolvedBlock, err := e.BlockByHash(maybeAncestor)
	if err != nil {
		log.Panic("could not fetch parent block with hash %s. Cause: %s", maybeAncestor.String(), err)
	}
	if resolvedBlock == nil {
		if resolvedBlock.Number().Int64() >= block.Number().Int64() {
			return false
		}
	}

	p, err := e.BlockByHash(block.ParentHash())
	if err != nil {
		log.Panic("could not fetch parent block with hash %s. Cause: %s", block.ParentHash().String(), err)
	}
	if p == nil {
		return false
	}

	return e.IsBlockAncestor(p, maybeAncestor)
}

func (e *gethRPCClient) RPCBlockchainFeed() []*types.Block {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	var availBlocks []*types.Block
	block, err := e.client.BlockByNumber(ctx, nil)
	if err != nil {
		log.Panic("could not fetch head block. Cause: %s", err)
	}
	availBlocks = append(availBlocks, block)

	for {
		// todo set this to genesis hash
		if block.ParentHash().Hex() == "0x0000000000000000000000000000000000000000000000000000000000000000" {
			break
		}
		block = e.getParentBlock(block)
		availBlocks = append(availBlocks, block)
	}

	// TODO double check the list is ordered [genesis, 1, 2, 3, 4, ..., last]
	// TODO It's pretty ugly but it avoids creating a new slice
	// TODO The approach of feeding all the blocks should change from all-blocks-in-memory to a stream
	for i, j := 0, len(availBlocks)-1; i < j; i, j = i+1, j-1 {
		availBlocks[i], availBlocks[j] = availBlocks[j], availBlocks[i]
	}
	return availBlocks
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
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	// this channel holds blocks that have been received from the geth network but not yet processed by the host,
	// with more than 1 capacity the buffer provides resilience in case of intermittent RPC or processing issues
	ch := make(chan *types.Header, 100)
	sub, err := e.client.SubscribeNewHead(ctx, ch)
	if err != nil {
		log.Panic("could not subscribe for new head blocks. Cause: %s", err)
	}

	return ch, sub
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

func (e *gethRPCClient) BalanceAt(gethcommon.Address, *big.Int) (*big.Int, error) {
	panic("not implemented")
}

func (e *gethRPCClient) Stop() {
	e.client.Close()
}

func connect(address string, connectionTimeout time.Duration) (*ethclient.Client, error) {
	var err error
	var c *ethclient.Client
	for start := time.Now(); time.Since(start) < connectionTimeout; time.Sleep(time.Second) {
		c, err = ethclient.Dial(address)
		if err == nil {
			break
		}
	}

	return c, err
}

func (e *gethRPCClient) getParentBlock(block *types.Block) *types.Block {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	parentBlock, err := e.client.BlockByHash(ctx, block.ParentHash())
	if err != nil {
		log.Panic("could not fetch parent block with hash %s. Cause: %s", block.ParentHash().String(), err)
	}

	return parentBlock
}
