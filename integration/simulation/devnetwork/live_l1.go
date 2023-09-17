package devnetwork

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"
)

type liveL1Network struct {
	deployWallet wallet.Wallet // wallet that can deploy to the live L1 network
	rpcURLs      []string

	clients          []ethadapter.EthClient
	seqWallet        wallet.Wallet
	validatorWallets []wallet.Wallet
}

func (l *liveL1Network) Prepare() {
	// nothing to do really, sanity check the L1 connection
	logger := testlog.Logger()
	l.prepareClients()
	client := l.GetClient(0)
	blockNum, err := client.BlockNumber()
	if err != nil {
		panic(fmt.Sprintf("unable to fetch head block number for live L1 network, rpc=%s err=%s",
			l.rpcURLs[0], err))
	}
	fmt.Println("Connected to L1 successfully", "currHeight", blockNum)
	logger.Info("Connected to L1 successfully", "currHeight", blockNum)

	nonce, err := client.Nonce(l.deployWallet.Address())
	if err != nil {
		panic(err)
	}
	l.deployWallet.SetNonce(nonce)
}

func (l *liveL1Network) CleanUp() {
	// nothing to clean up
}

func (l *liveL1Network) NumNodes() int {
	return 1
}

func (l *liveL1Network) GetClient(i int) ethadapter.EthClient {
	return l.clients[i%len(l.clients)]
}

func (l *liveL1Network) prepareClients() {
	l.clients = make([]ethadapter.EthClient, len(l.rpcURLs))
	for i, addr := range l.rpcURLs {
		client, err := ethadapter.NewEthClientFromURL(addr, 20*time.Second, common.HexToAddress("0x0"), testlog.Logger())
		if err != nil {
			panic(fmt.Sprintf("unable to create live L1 eth client, addr=%s err=%s", addr, err))
		}
		l.clients[i] = client
	}
}
