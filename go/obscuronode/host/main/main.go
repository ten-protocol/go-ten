package main

import (
	"math/big"
	"os"
	"time"

	"github.com/obscuronet/obscuro-playground/go/l1client"

	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/p2p"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
)

func main() {
	setLogs()
	config := parseCLIArgs()

	nodeID := common.BytesToAddress([]byte(*config.nodeID))
	hostCfg := host.AggregatorCfg{GossipRoundDuration: time.Duration(*config.gossipRoundNanos), ClientRPCTimeoutSecs: *config.rpcTimeoutSecs}
	enclaveClient := host.NewEnclaveRPCClient(*config.enclaveAddr, host.ClientRPCTimeoutSecs*time.Second, nodeID)
	aggP2P := p2p.NewSocketP2PLayer(*config.ourP2PAddr, config.peerP2PAddrs)
	agg := host.NewObscuroAggregator(nodeID, hostCfg, l1NodeDummy{}, nil, *config.isGenesis, enclaveClient, aggP2P, ethereum_mock.NewMockTxHandler())

	agg.Start()
}

// Sets the log file.
func setLogs() {
	logFile, err := os.Create("host_logs.txt")
	if err != nil {
		panic(err)
	}
	log.SetLog(logFile)
}

// TODO - Replace this dummy once we have implemented communication with L1 nodes.
type l1NodeDummy struct{}

func (l l1NodeDummy) IssueCustomTx(tx types.TxData) (*types.Transaction, error) {
	// TODO implement me
	panic("implement me")
}

func (l l1NodeDummy) TransactionReceipt(hash common.Hash) (*types.Receipt, error) {
	// TODO implement me
	panic("implement me")
}

func (l l1NodeDummy) FetchBlockByNumber(n *big.Int) (*types.Block, error) {
	// TODO implement me
	panic("implement me")
}

func (l l1NodeDummy) FetchHeadBlock() (*types.Block, uint64) {
	// TODO implement me
	panic("implement me")
}

func (l l1NodeDummy) Info() l1client.Info {
	// TODO implement me
	panic("implement me")
}

func (l l1NodeDummy) BlocksBetween(block *types.Block, head *types.Block) []*types.Block {
	// TODO implement me
	panic("implement me")
}

func (l l1NodeDummy) IsBlockAncestor(block *types.Block, proof obscurocommon.L1RootHash) bool {
	// TODO implement me
	panic("implement me")
}

func (l l1NodeDummy) BroadcastTx(t *obscurocommon.L1TxData) {
	// TODO implement me
	panic("implement me")
}

func (l l1NodeDummy) BlockListener() chan *types.Header {
	// TODO implement me
	panic("implement me")
}

func (l l1NodeDummy) FetchBlock(hash common.Hash) (*types.Block, error) {
	// TODO implement me
	panic("implement me")
}

func (l l1NodeDummy) RPCBlockchainFeed() []*types.Block {
	return []*types.Block{}
}
