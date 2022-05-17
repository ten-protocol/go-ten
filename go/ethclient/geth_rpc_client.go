package ethclient

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/wallet"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// TODO move this to a config
var connectionTimeout = 15 * time.Second

// gethRPCClient implements the EthClient interface and allows connection to a real ethereum node
// Beyond connection, EthClient requires transaction transformation to be handled (txhandle),
// chainID and transaction signage to be done (wallet)
type gethRPCClient struct {
	client  *ethclient.Client // the underlying eth rpc client
	id      common.Address    // TODO remove the id common.Address
	wallet  wallet.Wallet     // wallet containing the account information // TODO this does not need to be coupled together
	chainID int               // chainID is used to sign transactions
}

// NewEthClient instantiates a new ethclient.EthClient that connects to an ethereum node
func NewEthClient(id common.Address, ipaddress string, port uint, wallet wallet.Wallet, contractAddress *common.Address) (EthClient, error) {
	client, err := connect(ipaddress, port)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the eth node - %w", err)
	}

	// gets the next nonce to use on the account
	nonce, err := client.PendingNonceAt(context.Background(), wallet.Address())
	if err != nil {
		panic(err)
	}

	wallet.SetNonce(nonce)

	log.Trace("Initialized eth node connection - rollup contract address: %s - port: %d - wallet: %s - id: %s", contractAddress, port, wallet.Address(), id.String())
	return &gethRPCClient{
		client:  client,
		id:      id,
		wallet:  wallet, // TODO this does not need to be coupled together
		chainID: 1337,   // hardcoded for testnets // TODO this should be configured
	}, nil
}

func (e *gethRPCClient) FetchHeadBlock() *types.Block {
	blk, err := e.client.BlockByNumber(context.Background(), nil)
	if err != nil {
		panic(err)
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

	for currentBlk := lastBlock; currentBlk != nil && currentBlk.Hash() != startingBlock.Hash() && currentBlk.ParentHash() != common.HexToHash(""); {
		currentBlk, err = e.FetchBlock(currentBlk.ParentHash())
		if err != nil {
			panic(err)
		}
		blocksBetween = append(blocksBetween, currentBlk)
	}

	return blocksBetween
}

func (e *gethRPCClient) IsBlockAncestor(block *types.Block, maybeAncestor obscurocommon.L1RootHash) bool {
	if maybeAncestor == block.Hash() || maybeAncestor == obscurocommon.GenesisBlock.Hash() {
		return true
	}

	if block.Number().Int64() == int64(obscurocommon.L1GenesisHeight) {
		return false
	}

	resolvedBlock, err := e.FetchBlock(maybeAncestor)
	if err != nil {
		panic(err)
	}
	if resolvedBlock == nil {
		if resolvedBlock.Number().Int64() >= block.Number().Int64() {
			return false
		}
	}

	p, err := e.FetchBlock(block.ParentHash())
	if err != nil {
		panic(err)
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
		panic(err)
	}
	availBlocks = append(availBlocks, block)

	for {
		// todo set this to genesis hash
		if block.ParentHash().Hex() == "0x0000000000000000000000000000000000000000000000000000000000000000" {
			break
		}

		block, err = e.client.BlockByHash(context.Background(), block.ParentHash())
		if err != nil {
			panic(err)
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

func (e *gethRPCClient) SubmitTransaction(tx types.TxData) (*types.Transaction, error) {
	signedTx, err := e.wallet.SignTransaction(e.chainID, tx)
	if err != nil {
		panic(err)
	}

	return signedTx, e.IssueTransaction(signedTx)
}

func (e *gethRPCClient) IssueTransaction(signedTx *types.Transaction) error {
	return e.client.SendTransaction(context.Background(), signedTx)
}

func (e *gethRPCClient) FetchTxReceipt(hash common.Hash) (*types.Receipt, error) {
	return e.client.TransactionReceipt(context.Background(), hash)
}

func (e *gethRPCClient) BroadcastTx(tx types.TxData) {
	if _, err := e.SubmitTransaction(tx); err != nil {
		panic(err)
	}
}

func (e *gethRPCClient) BlockListener() chan *types.Header {
	ch := make(chan *types.Header, 1)
	// TODO this should return the subscription and cleanly Unsubscribe() when the node shutsdown
	_, err := e.client.SubscribeNewHead(context.Background(), ch)
	if err != nil {
		panic(err)
	}

	return ch
}

func (e *gethRPCClient) FetchBlockByNumber(n *big.Int) (*types.Block, error) {
	return e.client.BlockByNumber(context.Background(), n)
}

func (e *gethRPCClient) FetchBlock(hash common.Hash) (*types.Block, error) {
	return e.client.BlockByHash(context.Background(), hash)
}

func (e *gethRPCClient) Stop() {
	e.client.Close()
}

func connect(ipaddress string, port uint) (*ethclient.Client, error) {
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
