package eth2network

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
	"github.com/stretchr/testify/assert"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

const (
	_startPort = integration.StartPortEth2NetworkTests
	// TODO ensure it works with more than 1 node
	_numTestNodes = 1
)

func TestEnsureBinariesAreAvail(t *testing.T) {
	path, err := EnsureBinariesExist()
	assert.Nil(t, err)
	t.Logf("Successfully downloaded files to %s", path)
}

func TestStartEth2Network(t *testing.T) {
	binDir, err := EnsureBinariesExist()
	assert.Nil(t, err)

	chainID := int(datagenerator.RandomUInt64())
	randomWallets := make([]wallet.Wallet, _numTestNodes)
	randomWalletAddrs := make([]string, _numTestNodes)
	for i := 0; i < _numTestNodes; i++ {
		randomWallets[i] = datagenerator.RandomWallet(int64(chainID))
		randomWalletAddrs[i] = randomWallets[i].Address().Hex()
	}

	network := NewEth2Network(
		binDir,
		_startPort,
		_startPort+integration.DefaultGethWSPortOffset,
		_startPort+integration.DefaultGethAUTHPortOffset,
		_startPort+integration.DefaultGethNetworkPortOffset,
		_startPort+integration.DefaultPrysmHTTPPortOffset,
		_startPort+integration.DefaultPrysmP2PPortOffset,
		chainID,
		_numTestNodes,
		1,
		randomWalletAddrs,
	)
	// wait until the merge has happened
	assert.Nil(t, network.Start())

	defer network.Stop() //nolint: errcheck

	// test input configurations
	t.Run("areConfigsUphold", func(t *testing.T) {
		areConfigsUphold(t, randomWallets[0].Address(), chainID)
	})

	// test number of nodes
	t.Run("numberOfNodes", func(t *testing.T) {
		numberOfNodes(t)
	})

	// txs are minted
	t.Run("txsAreMinted", func(t *testing.T) {
		txsAreMinted(t, randomWallets)
	})
}

func areConfigsUphold(t *testing.T, addr gethcommon.Address, chainID int) {
	url := fmt.Sprintf("http://127.0.0.1:%d", _startPort)
	conn, err := ethclient.Dial(url)
	assert.Nil(t, err)

	at, err := conn.BalanceAt(context.Background(), addr, nil)
	assert.Nil(t, err)
	assert.True(t, at.Cmp(big.NewInt(1)) == 1)

	id, err := conn.NetworkID(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, int64(chainID), id.Int64())
}

func numberOfNodes(t *testing.T) {
	for i := 0; i < _numTestNodes; i++ {
		url := fmt.Sprintf("http://127.0.0.1:%d", _startPort+i)

		req, err := http.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			url,
			strings.NewReader(`{"jsonrpc": "2.0", "method": "net_peerCount", "params": [], "id": 1}`),
		)
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")

		client := &http.Client{}
		response, err := client.Do(req)
		assert.Nil(t, err)

		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		var res map[string]interface{}
		err = json.Unmarshal(body, &res)
		assert.Nil(t, err)

		err = json.Unmarshal(body, &res)
		assert.Nil(t, err)

		assert.Equal(t, fmt.Sprintf("0x%x", _numTestNodes-1), res["result"])
	}
}

func txsAreMinted(t *testing.T, wallets []wallet.Wallet) {
	var err error

	ethclients := make([]ethadapter.EthClient, _numTestNodes)
	for i := 0; i < _numTestNodes; i++ {
		ethclients[i], err = ethadapter.NewEthClient("127.0.0.1", uint(_startPort+100+i), 30*time.Second, common.L2Address{}, gethlog.New())
		assert.Nil(t, err)
	}

	for i := 0; i < _numTestNodes; i++ {
		ethClient := ethclients[i]
		w := wallets[i]

		toAddr := datagenerator.RandomAddress()
		estimatedTx, err := ethClient.EstimateGasAndGasPrice(&types.LegacyTx{
			Nonce: w.GetNonceAndIncrement(),
			To:    &toAddr,
			Value: big.NewInt(100),
		}, w.Address())
		assert.Nil(t, err)

		signedTx, err := w.SignTransaction(estimatedTx)
		assert.Nil(t, err)

		err = ethClient.SendTransaction(signedTx)
		assert.Nil(t, err)

		fmt.Printf("Created Tx: %s on node %d\n", signedTx.Hash().Hex(), i)
		// make sure it's mined into a block within an acceptable time and avail in all nodes
		for j := i; j < _numTestNodes; j++ {
			fmt.Printf("Checking for tx receipt for %s on node %d\n", signedTx.Hash(), j)
			var receipt *types.Receipt
			for start := time.Now(); time.Since(start) < 30*time.Second; time.Sleep(time.Second) {
				receipt, err = ethclients[j].TransactionReceipt(signedTx.Hash())
				if err == nil {
					break
				}
				if !errors.Is(err, ethereum.NotFound) {
					t.Fatal(err)
				}
			}

			if receipt == nil {
				t.Fatalf("Did not mine the transaction after %s seconds in node %d - receipt: %+v", 30*time.Second, j, receipt)
			}
		}
	}
}
