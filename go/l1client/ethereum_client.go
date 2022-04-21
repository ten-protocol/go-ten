package l1client

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/obscuronet/obscuro-playground/go/l1client/wallet"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/obscuronet/obscuro-playground/go/buildhelper/buildconstants"
	"github.com/obscuronet/obscuro-playground/go/buildhelper/helpertypes"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

var connectionTimeout = 15 * time.Second
var nonceLock sync.RWMutex

type EthNode struct {
	client  *ethclient.Client
	id      common.Address // TODO remove the id common.Address
	wallet  wallet.Wallet
	chainID int
}

func (e *EthNode) FetchHeadBlock() (*types.Block, uint64) {
	//TODO implement me
	panic("implement me")
}

func (e *EthNode) IssueTx(tx obscurocommon.EncodedL1Tx) {
	//TODO implement me
	panic("implement me")
}

func (e *EthNode) Info() Info {
	//TODO implement me
	panic("implement me")
}

func (e *EthNode) BlocksBetween(block *types.Block, head *types.Block) []*types.Block {
	//TODO implement me
	panic("implement me")
}

func (e *EthNode) IsBlockAncestor(block *types.Block, proof obscurocommon.L1RootHash) bool {
	//TODO implement me
	panic("implement me")
}

// NewEthClient instantiates a new l1client.Client that connects to an ethereum node
func NewEthClient(id common.Address, ipaddress string, port uint) (Client, error) {
	client, err := connect(ipaddress, port)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the eth node - %w", err)
	}

	log.Log(fmt.Sprintf("Initializing eth node at contract: %s", buildconstants.ContractAddress))
	return &EthNode{
		client:  client,
		id:      id,
		wallet:  wallet.NewInMemoryWallet(),
		chainID: 1337,
	}, nil
}

func (e *EthNode) RPCBlockchainFeed() []*types.Block {
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

func (e *EthNode) BroadcastTx(tx *obscurocommon.L1TxData) {
	nonceLock.Lock()
	defer nonceLock.Unlock()

	fromAddr := e.wallet.Address()
	nonce, err := e.client.PendingNonceAt(context.Background(), fromAddr)
	if err != nil {
		panic(err)
	}

	ethTx := &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: big.NewInt(20000000000),
		Gas:      1024_000_000,
		To:       &buildconstants.ContractAddress,
	}

	contractABI, err := abi.JSON(strings.NewReader(buildconstants.ContractAbi))
	if err != nil {
		panic(err)
	}

	// TODO each of these cases should be a function:
	// TODO like: func createRollupTx() or func createDepositTx()
	// TODO And then eventually, these functions would be called directly, when we get rid of our special format. (we'll have to change the mock thing as well for that)
	switch tx.TxType {
	case obscurocommon.DepositTx:
		ethTx.Value = big.NewInt(int64(tx.Amount))
		data, err := contractABI.Pack("Deposit", tx.Dest)
		if err != nil {
			panic(err)
		}
		ethTx.Data = data
		log.Log(fmt.Sprintf("BROADCAST TX: Issuing DepositTx - Addr: %s deposited %d to %s ",
			fromAddr, tx.Amount, tx.Dest))

	case obscurocommon.RollupTx:
		r, err := nodecommon.DecodeRollup(tx.Rollup)
		if err != nil {
			panic(err)
		}
		zipped := helpertypes.Compress(tx.Rollup)
		encRollupData := helpertypes.EncodeToString(zipped)
		data, err := contractABI.Pack("AddRollup", encRollupData)
		if err != nil {
			panic(err)
		}

		ethTx.Data = data
		derolled, _ := nodecommon.DecodeRollup(tx.Rollup)

		msg := ethereum.CallMsg{
			From:     fromAddr,
			To:       &common.Address{},
			GasPrice: big.NewInt(20000000000),
			Gas:      1024_000_000,
			Data:     data,
		}

		estimate, err := e.client.EstimateGas(context.Background(), msg)
		if err != nil {
			panic(err)
		}
		log.Log(fmt.Sprintf("Estimated Cost of: %d\n", estimate))
		log.Log(fmt.Sprintf("BROADCAST TX - Issuing Rollup: %s - %d txs - datasize: %d - gas: %d \n", r.Hash(), len(derolled.Transactions), len(data), ethTx.Gas))

	case obscurocommon.StoreSecretTx:
		data, err := contractABI.Pack("StoreSecret", helpertypes.EncodeToString(tx.Secret))
		if err != nil {
			panic(err)
		}
		ethTx.Data = data
		log.Log(fmt.Sprintf("BROADCAST TX: Issuing StoreSecretTx: encoded as %s", helpertypes.EncodeToString(tx.Secret)))
	case obscurocommon.RequestSecretTx:
		data, err := contractABI.Pack("RequestSecret")
		if err != nil {
			panic(err)
		}
		ethTx.Data = data
		log.Log("BROADCAST TX: Issuing RequestSecret")
	}

	signedTx, err := e.wallet.SignTransaction(e.chainID, ethTx)
	if err != nil {
		panic(err)
	}

	err = e.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		panic(err)
	}
}

func (e *EthNode) BlockListener() chan *types.Header {
	ch := make(chan *types.Header, 1)
	subs, err := e.client.SubscribeNewHead(context.Background(), ch)
	if err != nil {
		panic(err)
	}
	// we should hook the subs to cleanup
	fmt.Println(subs)

	return ch
}

func (e *EthNode) FetchBlockByNumber(n *big.Int) (*types.Block, error) {
	return e.client.BlockByNumber(context.Background(), n)
}

func (e *EthNode) FetchBlock(hash common.Hash) (*types.Block, error) {
	return e.client.BlockByHash(context.Background(), hash)
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
