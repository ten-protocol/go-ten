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

	"github.com/obscuronet/go-obscuro/go/config"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// gethRPCClient implements the EthClient interface and allows connection to a real ethereum node
type gethRPCClient struct {
	client *ethclient.Client  // the underlying eth rpc client
	id     gethcommon.Address // TODO remove the id common.Address
}

// NewEthClientFromConfig instantiates a new ethadapter.EthClient that connects to an ethereum node
// TODO Refactor this according with https://github.com/obscuronet/go-obscuro/pull/310#discussion_r910705504
func NewEthClientFromConfig(config config.HostConfig) (EthClient, error) {
	client, err := connect(config.L1NodeHost, config.L1NodeWebsocketPort, config.L1ConnectionTimeout)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the eth node - %w", err)
	}

	log.Trace("Initialized eth node connection - port: %d - id: %s", config.L1NodeWebsocketPort, config.ID)
	return &gethRPCClient{
		client: client,
		id:     config.ID,
	}, nil
}

// NewEthClient instantiates a new ethadapter.EthClient that connects to an ethereum node
func NewEthClient(ipaddress string, port uint, timeout time.Duration) (EthClient, error) {
	client, err := connect(ipaddress, port, timeout)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the eth node - %w", err)
	}

	log.Trace("Initialized eth node connection - addr: %s port: %d", ipaddress, port)
	return &gethRPCClient{
		client: client,
	}, nil
}

func (e *gethRPCClient) FetchHeadBlock() *types.Block {
	blk, err := e.client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Panic("could not fetch head block. Cause: %s", err)
	}
	return blk
}

func (e *gethRPCClient) Info() Info {
	return Info{
		ID: e.id,
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
	var availBlocks []*types.Block

	block, err := e.client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Panic("could not fetch head block. Cause: %s", err)
	}
	availBlocks = append(availBlocks, block)

	for {
		// todo set this to genesis hash
		if block.ParentHash().Hex() == "0x0000000000000000000000000000000000000000000000000000000000000000" {
			break
		}

		block, err = e.client.BlockByHash(context.Background(), block.ParentHash())
		if err != nil {
			log.Panic("could not fetch parent block with hash %s. Cause: %s", block.ParentHash().String(), err)
		}

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
	return e.client.SendTransaction(context.Background(), signedTx)
}

func (e *gethRPCClient) TransactionReceipt(hash gethcommon.Hash) (*types.Receipt, error) {
	return e.client.TransactionReceipt(context.Background(), hash)
}

func (e *gethRPCClient) Nonce(account gethcommon.Address) (uint64, error) {
	return e.client.PendingNonceAt(context.Background(), account)
}

func (e *gethRPCClient) BlockListener() chan *types.Header {
	ch := make(chan *types.Header, 1)
	// TODO this should return the subscription and cleanly Unsubscribe() when the node shutsdown
	_, err := e.client.SubscribeNewHead(context.Background(), ch)
	if err != nil {
		log.Panic("could not subscribe for new head blocks. Cause: %s", err)
	}

	return ch
}

func (e *gethRPCClient) BlockByNumber(n *big.Int) (*types.Block, error) {
	return e.client.BlockByNumber(context.Background(), n)
}

func (e *gethRPCClient) BlockByHash(hash gethcommon.Hash) (*types.Block, error) {
	return e.client.BlockByHash(context.Background(), hash)
}

func (e *gethRPCClient) CallContract(msg ethereum.CallMsg) ([]byte, error) {
	return e.client.CallContract(context.Background(), msg, nil)
}

func (e *gethRPCClient) EthClient() *ethclient.Client {
	return e.client
}

func (e *gethRPCClient) Stop() {
	e.client.Close()
}

func connect(ipaddress string, port uint, connectionTimeout time.Duration) (*ethclient.Client, error) {
	var err error
	var c *ethclient.Client
	for start := time.Now(); time.Since(start) < connectionTimeout; time.Sleep(time.Second) {
		c, err = ethclient.Dial(fmt.Sprintf("ws://%s:%d", ipaddress, port))
		if err == nil {
			break
		}
	}

	return c, err
}
