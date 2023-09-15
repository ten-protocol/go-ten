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
	rpcAddress   string

	client           ethadapter.EthClient
	seqWallet        wallet.Wallet
	validatorWallets []wallet.Wallet
}

func (l *liveL1Network) Prepare() {
	// nothing to do really, sanity check the L1 connection
	logger := testlog.Logger()
	client, err := ethadapter.NewEthClientFromURL(l.rpcAddress, 20*time.Second, common.HexToAddress("0x0"), logger)
	if err != nil {
		panic("unable to create live L1 eth client, err=" + err.Error())
	}
	l.client = client
	blockNum, err := client.BlockNumber()
	if err != nil {
		panic(fmt.Sprintf("unable to fetch head block number for live L1 network, rpc=%s err=%s",
			l.rpcAddress, err))
	}
	fmt.Println("Connected to L1 successfully", "currHeight", blockNum)
	logger.Info("Connected to L1 successfully", "currHeight", blockNum)

	nonce, err := l.client.Nonce(l.deployWallet.Address())
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

func (l *liveL1Network) GetClient(_ int) ethadapter.EthClient {
	return l.client
}
